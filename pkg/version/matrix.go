// SPDX-License-Identifier: Apache-2.0

// matrix.go — generic Matrix[E MatrixEntry] 구조 (ADR-0004).
//
// downstream operator 의 internal/version/matrix.go 의 Combo / Channel /
// FeatureGate 같은 *struct entry* 패턴을 표준화. 단순 []string 화이트리스트
// (기존 List) 와 별개로 *rich entry* 표현 필요 시 사용.
//
// 사용 예 (downstream operator version matrix source):
//
//	type Combo struct {
//	    Major   string  // PG major version (예: "16")
//	    Image   string  // 컨테이너 이미지
//	    Channel string  // alpha | beta | stable
//	    Gates   []string
//	}
//
//	func (c Combo) PrimaryKey() string { return c.Major }
//
//	var Supported = version.MustMatrix(
//	    Combo{Major: "16", Image: "...", Channel: "stable", Gates: nil},
//	    Combo{Major: "17", Image: "...", Channel: "beta", Gates: []string{"experimental"}},
//	)
//
// 호환성:
//   - 기존 `List` API 무변경 — 단순 string 화이트리스트는 List 그대로.
//   - Matrix[E] 는 *추가* — semver / channel / feature gate 등 메타데이터
//     필요한 operator 만 사용.

package version

// MatrixEntry — Matrix[E] 의 element 가 구현해야 하는 interface.
//
// PrimaryKey() 는 matrix 내 unique 식별자. duplicate 는 MustMatrix 가
// init-time panic. 빈 string PrimaryKey 도 panic — 누락된 const / typo
// 가드.
type MatrixEntry interface {
	PrimaryKey() string
}

// Matrix — generic 화이트리스트. element 는 MatrixEntry 구현 필수.
type Matrix[E MatrixEntry] struct {
	entries []E
}

// MustMatrix — init-time 화이트리스트 선언. 빈 리스트 / 빈 PrimaryKey /
// duplicate PrimaryKey 시 panic.
//
// 각 operator 의 internal/version/*.go 또는 api/v1alpha*/*_types.go 에서
// var ... = version.MustMatrix(...) 형태 사용.
func MustMatrix[E MatrixEntry](entries ...E) Matrix[E] {
	if len(entries) == 0 {
		panic("version.MustMatrix: 빈 리스트는 허용되지 않음 — 최소 1개 이상 명시")
	}
	seen := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		k := e.PrimaryKey()
		if k == "" {
			panic("version.MustMatrix: 빈 PrimaryKey — 누락된 const 또는 typo 추정")
		}
		if _, ok := seen[k]; ok {
			panic("version.MustMatrix: 중복 PrimaryKey: " + k)
		}
		seen[k] = struct{}{}
	}
	cp := make([]E, len(entries))
	copy(cp, entries)
	return Matrix[E]{entries: cp}
}

// Find — PrimaryKey 로 entry 검색. 없으면 (zero E, false).
func (m Matrix[E]) Find(key string) (E, bool) {
	for _, e := range m.entries {
		if e.PrimaryKey() == key {
			return e, true
		}
	}
	var zero E
	return zero, false
}

// IsSupported — key 가 matrix 에 존재하는지.
func (m Matrix[E]) IsSupported(key string) bool {
	_, ok := m.Find(key)
	return ok
}

// Entries — 모든 entry 반환 (방어 복사).
func (m Matrix[E]) Entries() []E {
	cp := make([]E, len(m.entries))
	copy(cp, m.entries)
	return cp
}

// Keys — 모든 PrimaryKey 반환. webhook 의 field.NotSupported 호환.
func (m Matrix[E]) Keys() []string {
	keys := make([]string, len(m.entries))
	for i, e := range m.entries {
		keys[i] = e.PrimaryKey()
	}
	return keys
}

// Len — entry 수.
func (m Matrix[E]) Len() int {
	return len(m.entries)
}
