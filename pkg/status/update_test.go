// SPDX-License-Identifier: MIT

package status

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const (
	updateTestNS  = "ns"
	updateTestPod = "pod-1"
)

func podScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	if err := corev1.AddToScheme(s); err != nil {
		t.Fatalf("corev1 scheme: %v", err)
	}
	return s
}

func seedPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: updateTestPod, Namespace: updateTestNS},
		Status:     corev1.PodStatus{Phase: corev1.PodPending},
	}
}

// conflictingClient — 처음 conflictTimes 회의 status update 를 Conflict 로 강제하는
// fake client. calls 포인터로 status update 호출 횟수를 관측한다.
func conflictingClient(t *testing.T, conflictTimes int, calls *int) client.Client {
	t.Helper()
	return fake.NewClientBuilder().
		WithScheme(podScheme(t)).
		WithObjects(seedPod()).
		WithStatusSubresource(&corev1.Pod{}).
		WithInterceptorFuncs(interceptor.Funcs{
			SubResourceUpdate: func(ctx context.Context, inner client.Client, sub string,
				obj client.Object, opts ...client.SubResourceUpdateOption) error {
				*calls++
				if *calls <= conflictTimes {
					return apierrors.NewConflict(
						schema.GroupResource{Resource: "pods"}, obj.GetName(), nil)
				}
				return inner.SubResource(sub).Update(ctx, obj, opts...)
			},
		}).
		Build()
}

func TestUpdateWithRetry(t *testing.T) {
	tests := []struct {
		name          string
		conflictTimes int
		withMutate    bool
		wantPhase     corev1.PodPhase
		wantCalls     int
	}{
		{
			name:       "conflict 없음 — 1회 갱신으로 호출자 변경 영속",
			withMutate: true,
			wantPhase:  corev1.PodRunning,
			wantCalls:  1,
		},
		{
			name:          "conflict 1회 — refetch 후 mutate 재적용으로 호출자 변경 보존",
			conflictTimes: 1,
			withMutate:    true,
			wantPhase:     corev1.PodRunning,
			wantCalls:     2,
		},
		{
			name:          "conflict 1회 + mutate 부재 — 호출자 변경이 서버 상태로 덮임 (silent-loss 의미론)",
			conflictTimes: 1,
			withMutate:    false,
			wantPhase:     corev1.PodPending,
			wantCalls:     2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			calls := 0
			c := conflictingClient(t, tt.conflictTimes, &calls)

			pod := &corev1.Pod{}
			if err := c.Get(ctx,
				client.ObjectKey{Namespace: updateTestNS, Name: updateTestPod}, pod); err != nil {
				t.Fatalf("get: %v", err)
			}
			// 호출자 패턴 — status 변경을 1차 적용한 뒤, 동일 클로저를 재적용용으로 전달.
			mutate := func() { pod.Status.Phase = corev1.PodRunning }
			mutate()
			var args []func()
			if tt.withMutate {
				args = append(args, mutate)
			}

			if err := UpdateWithRetry(ctx, c, pod, args...); err != nil {
				t.Fatalf("UpdateWithRetry: %v", err)
			}
			if calls != tt.wantCalls {
				t.Errorf("status update %d회 호출, want %d", calls, tt.wantCalls)
			}

			got := &corev1.Pod{}
			if err := c.Get(ctx,
				client.ObjectKey{Namespace: updateTestNS, Name: updateTestPod}, got); err != nil {
				t.Fatalf("re-get: %v", err)
			}
			if got.Status.Phase != tt.wantPhase {
				t.Errorf("서버 phase = %s, want %s", got.Status.Phase, tt.wantPhase)
			}
		})
	}
}

func TestUpdateWithRetry_non_conflict_error_immediate_return(t *testing.T) {
	ctx := context.Background()
	calls := 0
	c := fake.NewClientBuilder().
		WithScheme(podScheme(t)).
		WithObjects(seedPod()).
		WithStatusSubresource(&corev1.Pod{}).
		WithInterceptorFuncs(interceptor.Funcs{
			SubResourceUpdate: func(_ context.Context, _ client.Client, _ string,
				_ client.Object, _ ...client.SubResourceUpdateOption) error {
				calls++
				return apierrors.NewBadRequest("boom")
			},
		}).
		Build()

	err := UpdateWithRetry(ctx, c, seedPod())
	if err == nil || !apierrors.IsBadRequest(err) {
		t.Fatalf("BadRequest 즉시 전파 기대, got %v", err)
	}
	if calls != 1 {
		t.Errorf("비-conflict 에러는 재시도 없이 1회 호출 기대, got %d", calls)
	}
}
