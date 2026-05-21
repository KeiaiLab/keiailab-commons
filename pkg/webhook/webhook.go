// SPDX-License-Identifier: Apache-2.0

// Package webhook — admission validation 보조 helpers.
//
// downstream operator 의 webhook validation 공통 패턴:
//
//	if v.Spec.Version.Version != "" && !IsSupported(v.Spec.Version.Version) {
//	    errs = append(errs, field.NotSupported(path, value, allowedList))
//	}
//
// 본 패키지는 위 패턴을 *one-liner* 로 통합. 빈 문자열은 *defaulter 책임* —
// validation 이 skip (defaulter 가 채운 후 다시 검증).
//
// 사용 예 (valkey webhook):
//
//	import commonswebhook "github.com/keiailab/operator-commons/pkg/webhook"
//	import commonsversion "github.com/keiailab/operator-commons/pkg/version"
//
//	var supportedValkey = commonsversion.MustList("8.0.9", "8.1.6", "9.0.4")
//
//	func validateValkeySpec(v *Valkey) field.ErrorList {
//	    var errs field.ErrorList
//	    p := field.NewPath("spec")
//	    if err := commonswebhook.ValidateAllowedVersion(
//	        p.Child("version", "version"), v.Spec.Version.Version, supportedValkey,
//	    ); err != nil {
//	        errs = append(errs, err)
//	    }
//	    return errs
//	}
package webhook

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/keiailab/operator-commons/pkg/version"
)

// ValidateAllowedVersion — value 가 list 에 정확 매칭하는지 검증.
//
//   - 빈 문자열 → nil (defaulter 가 처리한다는 가정).
//   - 허용 → nil.
//   - 거부 → *field.NotSupported (allowed = list.Strings()).
//
// semver-prefix 매칭 (8.3.1 → 8.3) 등 *exact match 외* 검증은 본 함수 외부 —
// caller 가 자체 predicate 적용 후 ValidateWithPredicate 사용.
func ValidateAllowedVersion(path *field.Path, value string, list version.List) *field.Error {
	if value == "" {
		return nil
	}
	if list.IsSupported(value) {
		return nil
	}
	return field.NotSupported(path, value, list.Strings())
}

// ValidateWithPredicate — caller 가 정의한 predicate 함수 + 허용 목록.
// semver-prefix 매칭 같은 비-exact 검증에 사용.
//
// 예 (mongodb):
//
//		commonswebhook.ValidateWithPredicate(path, value,
//		    IsSupportedMongoDBVersion, // major.minor 추출 후 matching
//		    SupportedMongoDBVersions,  // 외부 노출 슬라이스
//		)
//
//	  - 빈 문자열 → nil.
//	  - predicate(value) == true → nil.
//	  - false → *field.NotSupported.
func ValidateWithPredicate(
	path *field.Path,
	value string,
	predicate func(string) bool,
	allowed []string,
) *field.Error {
	if value == "" {
		return nil
	}
	if predicate(value) {
		return nil
	}
	return field.NotSupported(path, value, allowed)
}
