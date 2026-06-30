// SPDX-License-Identifier: MIT

package pdb

import (
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Params 는 Build 의 입력이다. CR 타입 비의존을 위해 min/max 는 호출자가 자신의
// spec 에서 추출한 *intstr.IntOrString 로 전달한다.
type Params struct {
	// Name / Namespace / Labels — 생성될 PDB 의 ObjectMeta.
	Name      string
	Namespace string
	Labels    map[string]string

	// Selector — 대상 Pod 의 matchLabels (보통 SelectorLabels).
	Selector map[string]string

	// Replicas — 기본 minAvailable 계산용 (min/max 둘 다 nil 일 때만 사용).
	Replicas int32

	// MinAvailable / MaxUnavailable — 호출자 spec 에서 추출. 둘 다 nil 이면
	// 기본 정책(아래 DefaultFloor)을 적용한다. 둘 다 설정 시 MinAvailable 우선
	// (webhook 이 동시 설정을 reject 해야 한다).
	MinAvailable   *intstr.IntOrString
	MaxUnavailable *intstr.IntOrString

	// DefaultFloor — min/max 미지정 시 minAvailable = max(Replicas-1, DefaultFloor).
	// valkey=1(primary 항상 보존), mongo=0.
	DefaultFloor int
}

// Build 는 Params 로 PodDisruptionBudget 을 조립한다.
//
// MinAvailable / MaxUnavailable 우선순위:
//   - MinAvailable 설정 → 그 값.
//   - 아니면 MaxUnavailable 설정 → 그 값.
//   - 둘 다 nil → minAvailable = max(Replicas-1, DefaultFloor).
func Build(p Params) *policyv1.PodDisruptionBudget {
	out := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels:    p.Labels,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: p.Selector},
		},
	}
	switch {
	case p.MinAvailable != nil:
		out.Spec.MinAvailable = p.MinAvailable
	case p.MaxUnavailable != nil:
		out.Spec.MaxUnavailable = p.MaxUnavailable
	default:
		minAvail := max(int(p.Replicas)-1, p.DefaultFloor)
		v := intstr.FromInt(minAvail)
		out.Spec.MinAvailable = &v
	}
	return out
}
