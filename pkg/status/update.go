// SPDX-License-Identifier: MIT

package status

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateWithRetry 는 obj 의 status subresource 를 conflict 발생 시 재시도하며 갱신한다.
//
// 배경 — r.Status().Update(ctx, obj) 직접 호출 패턴의 두 약점 (downstream operator 실측):
//  1. 한 reconcile 안에서 객체를 여러 번 mutate 한 뒤 마지막에 한 번만 Update 하므로,
//     도중에 외부 (다른 컨트롤러 / 사용자) 가 객체를 수정하면 Update 가 conflict 로
//     실패하고 reconcile 이 처음부터 다시 시작된다 — 무한 requeue 루프 가능.
//  2. conflict 시 latest 를 refetch 하면 호출자가 obj 에 설정해 둔 status mutation 이
//     서버 상태로 덮어쓰여 silent 하게 유실된다 (mutate-after-get 부재).
//
// 본 헬퍼는 retry.RetryOnConflict 로 conflict 를 흡수한다. 정상 경로 (conflict 없음) 는
// 호출자의 in-memory 사본을 그대로 영속화한다. conflict 발생 시에는 최신 ResourceVersion
// 으로 refetch 한 뒤, 가변 인자로 받은 mutate 콜백을 다시 적용해 호출자의 status 변경을
// 잃지 않고 재 Update 한다.
//
// 주의 — mutate 미전달 (0-인자) 호출의 conflict 경로 의미론: refetch 가 호출자의 in-memory
// status 변경을 *서버 상태로 덮어쓴 채* 재 Update 하고 nil (성공) 을 반환한다 = 호출자
// 변경의 silent-loss. 기존 0-인자 호출자의 점진 마이그레이션 호환을 위한 거동이며, *신규
// 호출자는 status 변경을 mutate 클로저로 전달할 것* 을 강력 권장한다 (valkey-operator
// 의 refetch 없는 naive retry 판 — stale ResourceVersion 재시도로 사실상 무효 — 은
// 본 헬퍼 채택 시 폐기 대상, 패키지 doc.go §의존성 정책 + pkg/reconcile doc.go 기록 참조).
//
// 본 헬퍼는 status subresource 만 다루므로 spec field 손실을 일으키지 않는다.
func UpdateWithRetry(ctx context.Context, c client.Client, obj client.Object, mutate ...func()) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := c.Status().Update(ctx, obj)
		if err == nil {
			return nil
		}
		if !apierrors.IsConflict(err) {
			// conflict 가 아닌 에러는 즉시 반환 (NotFound 등은 retry 의미 없음).
			return err
		}
		// conflict: 호출자가 보유한 ResourceVersion 이 stale — 최신본으로 refetch.
		if getErr := c.Get(ctx, client.ObjectKeyFromObject(obj), obj); getErr != nil {
			return getErr
		}
		// refetch 가 호출자의 status 변경을 덮어썼으므로 mutate 콜백으로 재적용 후 다시
		// Update. 콜백이 없으면 (기존 호출자) 갱신된 ResourceVersion 으로 한 번 더 재시도.
		for _, m := range mutate {
			m()
		}
		return c.Status().Update(ctx, obj)
	})
}
