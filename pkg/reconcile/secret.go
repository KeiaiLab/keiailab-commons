// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SecretIfNotExists 는 owner 의 namespace 에 secretName Secret 이 존재하지 않으면
// build() 결과로 생성한다. 이미 존재하면 무변경 no-op (멱등) — build 는 호출되지
// 않는다.
//
// keyfile / password Secret 처럼 *immutable + create-once* 자원에 사용한다
// (생성 후 값을 갱신하면 안 되는 credential 류).
//
//   - owner reference + scheme 은 GC (owner 삭제 시 Secret cascade) 를 위해 필수 —
//     controllerutil.SetControllerReference 로 controller=true ownerRef 부여.
//   - race-tolerant create: 병렬 reconcile 의 *Get NotFound → Create AlreadyExists*
//     race 를 흡수해 nil 을 반환한다 (mongodb iteration 41 cross-cut canonical —
//     valkey 판은 guard 부재로 race 시 에러 전파하던 drift).
//   - build() 가 반환하는 Secret 의 namespace 는 owner 와 동일해야 한다 (cross-ns
//     ownerRef 는 K8s 가 금지).
func SecretIfNotExists(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	secretName string,
	build func() *corev1.Secret,
) error {
	existing := &corev1.Secret{}
	err := c.Get(ctx, client.ObjectKey{Name: secretName, Namespace: owner.GetNamespace()}, existing)
	if err == nil {
		return nil // 이미 존재 — 무변경 (멱등)
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	secret := build()
	if err := controllerutil.SetControllerReference(owner, secret, scheme); err != nil {
		return fmt.Errorf("set owner ref: %w", err)
	}
	if err := c.Create(ctx, secret); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}
