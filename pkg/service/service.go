// SPDX-License-Identifier: MIT

package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Params 는 Build 의 입력이다. 포트/셀렉터/옵션은 호출자가 도메인 규칙대로 조립해 전달한다.
type Params struct {
	// Name / Namespace / Labels / Annotations — 생성될 Service 의 ObjectMeta.
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string

	// Selector — 대상 Pod matchLabels.
	Selector map[string]string
	// Ports — 노출 포트(호출자가 도메인 포트 규칙으로 조립).
	Ports []corev1.ServicePort

	// Headless=true → ClusterIP "None" + PublishNotReadyAddresses=true (StatefulSet
	// stable DNS). 이 경우 Type 은 무시된다.
	Headless bool
	// Type — non-headless Service 종류. "" → ClusterIP.
	Type corev1.ServiceType
	// PublishNotReadyAddresses — non-headless 에서도 미준비 주소 노출이 필요한 경우.
	// Headless=true 면 자동 true.
	PublishNotReadyAddresses bool

	// IPFamilyPolicy / IPFamilies — dual-stack 옵션(선택).
	IPFamilyPolicy *corev1.IPFamilyPolicy
	IPFamilies     []corev1.IPFamily
}

// Build 는 Params 로 corev1.Service 를 조립한다.
func Build(p Params) *corev1.Service {
	spec := corev1.ServiceSpec{
		Selector:                 p.Selector,
		Ports:                    p.Ports,
		IPFamilyPolicy:           p.IPFamilyPolicy,
		IPFamilies:               p.IPFamilies,
		PublishNotReadyAddresses: p.PublishNotReadyAddresses,
	}
	if p.Headless {
		spec.ClusterIP = corev1.ClusterIPNone
		spec.PublishNotReadyAddresses = true
	} else {
		spec.Type = p.Type
		if spec.Type == "" {
			spec.Type = corev1.ServiceTypeClusterIP
		}
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        p.Name,
			Namespace:   p.Namespace,
			Labels:      p.Labels,
			Annotations: p.Annotations,
		},
		Spec: spec,
	}
}
