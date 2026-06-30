// SPDX-License-Identifier: MIT

package secrethash

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// Hash 는 Secret data 의 결정적 SHA256 hex digest 를 반환한다.
//
// keys 가 주어지면 *그 키만, 인자 순서대로* 누적한다. 워크로드 rollout
// 트리거를 위해 특정 키(예: tls.crt / tls.key / ca.crt)의 변경만 추적할 때
// 사용한다. 누락된 키는 빈 값으로 취급되어 digest 에 키 이름만 기여한다.
//
// keys 가 비어 있으면 data 의 *모든 키를 정렬* 하여 누적한다. 전체 Secret
// 변경 추적용이며, map 순회 비결정성을 회피한다.
//
// 각 항목은 키 바이트 → 값 바이트 순으로 누적되므로, 키 추가/삭제나 값 변경이
// 모두 digest 에 반영된다. data 가 nil/빈 맵이고 keys 도 없으면 빈 입력의
// SHA256(고정값)을 반환한다.
func Hash(data map[string][]byte, keys ...string) string {
	h := sha256.New()
	if len(keys) == 0 {
		keys = make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}
	for _, k := range keys {
		_, _ = h.Write([]byte(k))
		_, _ = h.Write(data[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}
