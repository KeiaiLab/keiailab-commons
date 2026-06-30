// SPDX-License-Identifier: MIT

package hpa

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DefaultTargetAPIVersion — ScaleTargetRef 의 기본 apiVersion.
const DefaultTargetAPIVersion = "apps/v1"

// Params 는 Build 의 입력이다. metric 은 호출자가 조립해 전달한다 (CR 타입 비의존).
type Params struct {
	// Name / Namespace / Labels — 생성될 HPA 의 ObjectMeta.
	Name      string
	Namespace string
	Labels    map[string]string

	// TargetAPIVersion — ScaleTargetRef apiVersion. "" 이면 DefaultTargetAPIVersion("apps/v1").
	TargetAPIVersion string
	// TargetKind / TargetName — scale 대상 (예: "StatefulSet", <sts-name>).
	TargetKind string
	TargetName string

	// MinReplicas / MaxReplicas — 원하는 범위.
	MinReplicas int32
	MaxReplicas int32
	// MinFloor — MinReplicas clamp 하한 (valkey=2, mongo=1).
	MinFloor int32

	// Metrics — 호출자가 조립한 metric 목록 (CPUUtilization/MemoryUtilization 헬퍼 활용).
	Metrics []autoscalingv2.MetricSpec
}

// Build 는 Params 로 HorizontalPodAutoscaler 를 조립한다.
//
// clamp: MinReplicas = max(MinReplicas, MinFloor),
// MaxReplicas = max(MaxReplicas, 보정된 MinReplicas).
func Build(p Params) *autoscalingv2.HorizontalPodAutoscaler {
	apiVersion := p.TargetAPIVersion
	if apiVersion == "" {
		apiVersion = DefaultTargetAPIVersion
	}
	minR := max(p.MinReplicas, p.MinFloor)
	maxR := max(p.MaxReplicas, minR)
	return &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels:    p.Labels,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: apiVersion,
				Kind:       p.TargetKind,
				Name:       p.TargetName,
			},
			MinReplicas: ptr.To(minR),
			MaxReplicas: maxR,
			Metrics:     p.Metrics,
		},
	}
}

// CPUUtilization — CPU 평균 사용률(%) 타깃 Resource MetricSpec.
func CPUUtilization(percent int32) autoscalingv2.MetricSpec {
	return resourceUtilization(corev1.ResourceCPU, percent)
}

// MemoryUtilization — Memory 평균 사용률(%) 타깃 Resource MetricSpec.
func MemoryUtilization(percent int32) autoscalingv2.MetricSpec {
	return resourceUtilization(corev1.ResourceMemory, percent)
}

func resourceUtilization(name corev1.ResourceName, percent int32) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: name,
			Target: autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: ptr.To(percent),
			},
		},
	}
}
