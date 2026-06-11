// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/keiailab/keiailab-commons/pkg/finalizer"
)

// HandleFinalizerCleanup 는 deletionTimestamp 가 설정된 객체에 대해
// (1) cleanup() 콜백 실행 → (2) finalizer 제거 → (3) client.Update 패턴을 처리한다.
//
//   - obj 에 finalizerName 이 없으면 no-op (이미 처리됨 / 무관 객체).
//   - cleanup 이 nil 이면 finalizer 제거만 수행.
//   - cleanup 에러 시 finalizer 를 유지한 채 에러 반환 — 다음 reconcile 에서 재시도
//     가능 (외부 자원 정리의 at-least-once 보장).
//
// in-memory finalizer 조작은 pkg/finalizer 위임 — 본 함수는 그 위의 orchestration
// (cleanup 주입 + Update 영속화) 만 담당한다. deletionTimestamp 검사는 호출자 책임
// (controller Reconcile 진입부의 표준 분기).
//
//nolint:unparam // controller-runtime 표준 (ctrl.Result, error) 시그니처 보존 — 호출자 일관성 (원본 양 판 정합).
func HandleFinalizerCleanup(
	ctx context.Context,
	c client.Client,
	obj client.Object,
	finalizerName string,
	cleanup func(context.Context) error,
) (ctrl.Result, error) {
	if !finalizer.Has(obj, finalizerName) {
		return ctrl.Result{}, nil
	}

	if cleanup != nil {
		if err := cleanup(ctx); err != nil {
			return ctrl.Result{}, fmt.Errorf("finalizer cleanup: %w", err)
		}
	}

	finalizer.Remove(obj, finalizerName)
	if err := c.Update(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
