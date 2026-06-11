// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const testSecretName = "keyfile-secret"

func newOwner() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "owner-cm", Namespace: testNS, UID: "owner-uid-1"},
	}
}

func buildTestSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: testSecretName, Namespace: testNS},
		Data:       map[string][]byte{"key": []byte("generated")},
	}
}

func TestSecretIfNotExists_creates_with_ownerRef(t *testing.T) {
	ctx := context.Background()
	owner := newOwner()
	scheme := testScheme(t)
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(owner).Build()

	buildCalled := false
	err := SecretIfNotExists(ctx, c, scheme, owner, testSecretName, func() *corev1.Secret {
		buildCalled = true
		return buildTestSecret()
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !buildCalled {
		t.Error("부재 시 build 호출 기대")
	}

	got := &corev1.Secret{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: testNS, Name: testSecretName}, got); err != nil {
		t.Fatalf("secret 생성 기대: %v", err)
	}
	if string(got.Data["key"]) != "generated" {
		t.Errorf("data = %q, want generated", got.Data["key"])
	}
	if len(got.OwnerReferences) != 1 {
		t.Fatalf("ownerRef 1건 기대, got %d", len(got.OwnerReferences))
	}
	if ref := got.OwnerReferences[0]; ref.UID != owner.UID || ref.Controller == nil || !*ref.Controller {
		t.Errorf("controller ownerRef 기대, got %+v", ref)
	}
}

func TestSecretIfNotExists_idempotent_when_exists(t *testing.T) {
	ctx := context.Background()
	owner := newOwner()
	scheme := testScheme(t)
	existing := buildTestSecret()
	existing.Data = map[string][]byte{"key": []byte("original")}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(owner, existing).Build()

	buildCalled := false
	err := SecretIfNotExists(ctx, c, scheme, owner, testSecretName, func() *corev1.Secret {
		buildCalled = true
		return buildTestSecret()
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if buildCalled {
		t.Error("이미 존재 — build 미호출 기대")
	}

	got := &corev1.Secret{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: testNS, Name: testSecretName}, got); err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(got.Data["key"]) != "original" {
		t.Errorf("기존 데이터 무변경 기대, got %q", got.Data["key"])
	}
}

func TestSecretIfNotExists_already_exists_race_tolerated(t *testing.T) {
	// Get NotFound → Create 사이에 병렬 reconcile 이 먼저 생성한 race 시뮬레이션
	// (mongodb iteration 41 canonical — AlreadyExists 는 성공으로 흡수).
	ctx := context.Background()
	owner := newOwner()
	scheme := testScheme(t)
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(owner).
		WithInterceptorFuncs(interceptor.Funcs{
			Create: func(_ context.Context, _ client.WithWatch,
				obj client.Object, _ ...client.CreateOption) error {
				return apierrors.NewAlreadyExists(
					schema.GroupResource{Resource: "secrets"}, obj.GetName())
			},
		}).Build()

	if err := SecretIfNotExists(ctx, c, scheme, owner, testSecretName, buildTestSecret); err != nil {
		t.Fatalf("AlreadyExists race 흡수 기대, got %v", err)
	}
}

func TestSecretIfNotExists_get_error_propagated(t *testing.T) {
	ctx := context.Background()
	owner := newOwner()
	scheme := testScheme(t)
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(owner).
		WithInterceptorFuncs(interceptor.Funcs{
			Get: func(_ context.Context, _ client.WithWatch, _ client.ObjectKey,
				_ client.Object, _ ...client.GetOption) error {
				return apierrors.NewServiceUnavailable("apiserver down")
			},
		}).Build()

	err := SecretIfNotExists(ctx, c, scheme, owner, testSecretName, buildTestSecret)
	if err == nil || !apierrors.IsServiceUnavailable(err) {
		t.Fatalf("비-NotFound Get 에러 전파 기대, got %v", err)
	}
}
