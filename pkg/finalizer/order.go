// SPDX-License-Identifier: Apache-2.0

package finalizer

import "slices"

// EnsureOrder 는 finalizers slice 가 주어진 desiredOrder 에 따라 정렬되도록
// 보장한다. desiredOrder 에 없는 finalizer 는 *뒤에* 안정 정렬로 유지.
//
// 호출 시 *반드시* result 를 다시 할당해야 한다 (in-place 정렬 + 빈 슬롯 보장).
//
// 사용 예 (multi-finalizer 순서 강제):
//
//	desired := []string{"order-1.example.com", "order-2.example.com"}
//	obj.Finalizers = finalizer.EnsureOrder(obj.Finalizers, desired)
//
// 본 helper 는 pkg/status + pkg/finalizer 표준 후속 ROADMAP "다중 finalizer 순서 보장" ()
// 의 표준 구현.
func EnsureOrder(finalizers []string, desiredOrder []string) []string {
	if len(finalizers) == 0 || len(desiredOrder) == 0 {
		return finalizers
	}
	// rank lookup: lower rank == earlier in desiredOrder
	rank := make(map[string]int, len(desiredOrder))
	for i, name := range desiredOrder {
		rank[name] = i
	}
	out := slices.Clone(finalizers)
	slices.SortStableFunc(out, func(a, b string) int {
		ra, aOK := rank[a]
		rb, bOK := rank[b]
		switch {
		case aOK && bOK:
			return ra - rb
		case aOK && !bOK:
			return -1 // a is ordered, b is not — a first
		case !aOK && bOK:
			return 1 // b is ordered, a is not — b first
		default:
			return 0 // both unordered — stable
		}
	})
	return out
}
