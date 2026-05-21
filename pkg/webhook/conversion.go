// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
)

// ConvertFunc — 단일 방향 변환 함수 (e.g. v1alpha1 → v1alpha2).
type ConvertFunc func(src any) (dst any, err error)

// ConversionRegistry — version pair → ConvertFunc 매핑.
//
// 사용 예 (downstream operator):
//
//	reg := webhook.NewConversionRegistry()
//	reg.Register("valkey.io/v1alpha1", "valkey.io/v1alpha2", convertV1A1ToV1A2)
//	reg.Register("valkey.io/v1alpha2", "valkey.io/v1alpha1", convertV1A2ToV1A1)
//	dst, err := reg.Convert("valkey.io/v1alpha1", "valkey.io/v1alpha2", src)
//
// 본 helper 는 controller-runtime 의존을 회피하고 callback 패턴으로 변환
// 로직을 위임. caller 는 별도 admission webhook server 에서 호출.
//
// Refs: docs/ROADMAP.md 'Conversion webhook helper — v1alpha1 ↔ v1alpha2 패턴 추출'
type ConversionRegistry struct {
	convs map[string]ConvertFunc
}

// NewConversionRegistry — 빈 registry 생성.
func NewConversionRegistry() *ConversionRegistry {
	return &ConversionRegistry{convs: make(map[string]ConvertFunc)}
}

// Register — fromVersion → toVersion 변환 등록.
func (r *ConversionRegistry) Register(fromVersion, toVersion string, fn ConvertFunc) {
	r.convs[key(fromVersion, toVersion)] = fn
}

// Convert — src 를 fromVersion → toVersion 으로 변환. 등록되지 않은 pair 는
// error.
func (r *ConversionRegistry) Convert(fromVersion, toVersion string, src any) (any, error) {
	fn, ok := r.convs[key(fromVersion, toVersion)]
	if !ok {
		return nil, fmt.Errorf("no conversion registered: %s → %s", fromVersion, toVersion)
	}
	return fn(src)
}

// HasPair — 등록 여부 확인.
func (r *ConversionRegistry) HasPair(fromVersion, toVersion string) bool {
	_, ok := r.convs[key(fromVersion, toVersion)]
	return ok
}

func key(from, to string) string { return from + "→" + to }
