// SPDX-License-Identifier: Apache-2.0

// Package version — operator 들이 공통으로 사용하는 *지원 DB 버전 화이트리스트* 패턴.
//
// 사용 예 (각 operator 의 api/v1alpha1/*_types.go):
//
//	var SupportedMongoDBVersions = version.MustList("8.0", "8.2", "8.3")
//
// 그리고 webhook validation 에서:
//
//	if !version.IsSupported(SupportedMongoDBVersions, v.Spec.Version) {
//	    errs = append(errs, field.NotSupported(...))
//	}
//
// 디자인 결정:
//   - 단순 []string 컨벤션 — semver 파싱 없이 정확 매칭. RDB / WAL / oplog format
//     호환성은 *각 operator 가 책임*.
//   - MustList 는 빈 리스트 또는 빈 문자열 포함 시 panic — init time 가드.
package version

import "slices"

// List — 지원 버전 화이트리스트의 immutable 표현.
type List struct {
	versions []string
}

// MustList — 컴파일 타임 화이트리스트 선언. 빈 리스트 / 빈 문자열 포함 시 panic.
// 각 operator 의 *_types.go 에서 var ... = version.MustList(...) 형태 사용.
func MustList(versions ...string) List {
	if len(versions) == 0 {
		panic("version.MustList: 빈 리스트는 허용되지 않음 — 최소 1개 이상 명시")
	}
	for _, v := range versions {
		if v == "" {
			panic("version.MustList: 빈 문자열 항목 — 누락된 const 또는 typo 추정")
		}
	}
	cp := make([]string, len(versions))
	copy(cp, versions)
	return List{versions: cp}
}

// IsSupported — v 가 화이트리스트에 정확 매칭하는지 검사. semver range 가 아닌 string equality.
func (l List) IsSupported(v string) bool {
	return slices.Contains(l.versions, v)
}

// Strings — webhook 의 field.NotSupported 에 전달할 슬라이스 (방어 복사).
func (l List) Strings() []string {
	cp := make([]string, len(l.versions))
	copy(cp, l.versions)
	return cp
}

// Default — 화이트리스트의 첫 항목을 default 로 사용. operator 의 defaulting webhook 에서 호출.
// 빈 리스트는 MustList 가드로 사전 차단되므로 panic 가능성 없음.
func (l List) Default() string {
	return l.versions[0]
}
