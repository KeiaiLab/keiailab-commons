// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"errors"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const testFinalizer = "test.keiailab.com/finalizer"

func newConfigMap(finalizers ...string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "cm-1",
			Namespace:  testNS,
			Finalizers: finalizers,
		},
	}
}

// markDeleting — fake client 에 Delete 를 보내 deletionTimestamp 를 설정한다
// (finalizer 보유 객체는 즉시 삭제되지 않음). cm 을 최신 서버 상태로 갱신.
func markDeleting(t *testing.T, c client.Client, cm *corev1.ConfigMap) {
	t.Helper()
	ctx := context.Background()
	if err := c.Delete(ctx, cm); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := c.Get(ctx, client.ObjectKeyFromObject(cm), cm); err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if cm.DeletionTimestamp.IsZero() {
		t.Fatal("deletionTimestamp 설정 기대")
	}
}

func TestHandleFinalizerCleanup(t *testing.T) {
	cleanupBoom := errors.New("cleanup boom")
	tests := []struct {
		name          string
		finalizers    []string
		nilCleanup    bool
		cleanupErr    error
		wantErr       error
		wantCleanup   bool
		wantGone      bool // finalizer 제거로 객체 삭제 완료 기대
		wantFinalizer bool // 잔존 객체의 testFinalizer 보유 기대
	}{
		{
			name:       "finalizer 부재 — no-op (cleanup 미호출 + 객체 무변경)",
			finalizers: []string{"other.keiailab.com/finalizer"},
		},
		{
			name:       "cleanup nil — finalizer 제거만 수행, 객체 삭제 완료",
			finalizers: []string{testFinalizer},
			nilCleanup: true,
			wantGone:   true,
		},
		{
			name:        "cleanup 성공 — 호출 후 finalizer 제거, 객체 삭제 완료",
			finalizers:  []string{testFinalizer},
			wantCleanup: true,
			wantGone:    true,
		},
		{
			name:          "cleanup 실패 — finalizer 유지 (재시도 가능) + 에러 wrap 전파",
			finalizers:    []string{testFinalizer},
			cleanupErr:    cleanupBoom,
			wantErr:       cleanupBoom,
			wantCleanup:   true,
			wantFinalizer: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cm := newConfigMap(tt.finalizers...)
			c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(cm).Build()
			markDeleting(t, c, cm)

			called := false
			var cleanup func(context.Context) error
			if !tt.nilCleanup {
				cleanup = func(context.Context) error {
					called = true
					return tt.cleanupErr
				}
			}

			_, err := HandleFinalizerCleanup(ctx, c, cm, testFinalizer, cleanup)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("cleanup 에러 wrap 전파 기대, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("err: %v", err)
			}
			if called != tt.wantCleanup {
				t.Errorf("cleanup 호출 = %v, want %v", called, tt.wantCleanup)
			}

			got := &corev1.ConfigMap{}
			getErr := c.Get(ctx, client.ObjectKeyFromObject(cm), got)
			if tt.wantGone {
				if !apierrors.IsNotFound(getErr) {
					t.Errorf("finalizer 제거 후 객체 삭제 완료 기대, got err=%v finalizers=%v",
						getErr, got.Finalizers)
				}
				return
			}
			if getErr != nil {
				t.Fatalf("객체 잔존 기대: %v", getErr)
			}
			if hasFin := slices.Contains(got.Finalizers, testFinalizer); hasFin != tt.wantFinalizer {
				t.Errorf("testFinalizer 보유 = %v, want %v (finalizers=%v)",
					hasFin, tt.wantFinalizer, got.Finalizers)
			}
		})
	}
}
