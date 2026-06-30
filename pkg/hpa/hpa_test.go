// SPDX-License-Identifier: MIT

package hpa_test

import (
	"testing"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/keiailab/keiailab-commons/pkg/hpa"
)

func TestBuild_ScaleTargetRefAndMeta(t *testing.T) {
	t.Parallel()
	got := hpa.Build(hpa.Params{
		Name: "cache", Namespace: "ns", Labels: map[string]string{"k": "v"},
		TargetKind: "StatefulSet", TargetName: "cache-sts",
		MinReplicas: 3, MaxReplicas: 5, MinFloor: 2,
		Metrics: []autoscalingv2.MetricSpec{hpa.CPUUtilization(70)},
	})
	if got.Spec.ScaleTargetRef.APIVersion != "apps/v1" {
		t.Fatalf("default apiVersion = %q, want apps/v1", got.Spec.ScaleTargetRef.APIVersion)
	}
	if got.Spec.ScaleTargetRef.Kind != "StatefulSet" || got.Spec.ScaleTargetRef.Name != "cache-sts" {
		t.Fatalf("ScaleTargetRef mismatch: %+v", got.Spec.ScaleTargetRef)
	}
	if got.Name != "cache" || got.Namespace != "ns" || got.Labels["k"] != "v" {
		t.Fatalf("ObjectMeta mismatch: %+v", got.ObjectMeta)
	}
}

func TestBuild_MinReplicasClampFloor(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name             string
		min, max, floor  int32
		wantMin, wantMax int32
	}{
		{"valkey floor2 위로 보정", 1, 5, 2, 2, 5}, // max(1,2)=2
		{"valkey floor2 미보정", 3, 5, 2, 3, 5},
		{"mongo floor1", 0, 4, 1, 1, 4}, // max(0,1)=1
		{"max<min 보정", 2, 1, 1, 2, 2},   // maxR=max(1,2)=2
		{"floor가 max도 끌어올림", 1, 1, 4, 4, 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := hpa.Build(hpa.Params{
				TargetKind: "StatefulSet", TargetName: "x",
				MinReplicas: c.min, MaxReplicas: c.max, MinFloor: c.floor,
			})
			if got.Spec.MinReplicas == nil || *got.Spec.MinReplicas != c.wantMin {
				t.Fatalf("minReplicas=%v, want %d", got.Spec.MinReplicas, c.wantMin)
			}
			if got.Spec.MaxReplicas != c.wantMax {
				t.Fatalf("maxReplicas=%d, want %d", got.Spec.MaxReplicas, c.wantMax)
			}
		})
	}
}

func TestBuild_CustomAPIVersion(t *testing.T) {
	t.Parallel()
	got := hpa.Build(hpa.Params{TargetAPIVersion: "apps/v1beta1", TargetKind: "Deployment", TargetName: "d"})
	if got.Spec.ScaleTargetRef.APIVersion != "apps/v1beta1" {
		t.Fatalf("apiVersion override 무시: %q", got.Spec.ScaleTargetRef.APIVersion)
	}
}

func TestCPUUtilization(t *testing.T) {
	t.Parallel()
	m := hpa.CPUUtilization(65)
	if m.Type != autoscalingv2.ResourceMetricSourceType {
		t.Fatalf("type=%v", m.Type)
	}
	if m.Resource == nil || m.Resource.Name != corev1.ResourceCPU {
		t.Fatalf("resource name mismatch: %+v", m.Resource)
	}
	if m.Resource.Target.Type != autoscalingv2.UtilizationMetricType {
		t.Fatalf("target type=%v", m.Resource.Target.Type)
	}
	if m.Resource.Target.AverageUtilization == nil || *m.Resource.Target.AverageUtilization != 65 {
		t.Fatalf("avgUtil=%v, want 65", m.Resource.Target.AverageUtilization)
	}
}

func TestMemoryUtilization(t *testing.T) {
	t.Parallel()
	m := hpa.MemoryUtilization(80)
	if m.Resource == nil || m.Resource.Name != corev1.ResourceMemory {
		t.Fatalf("resource name mismatch: %+v", m.Resource)
	}
	if m.Resource.Target.AverageUtilization == nil || *m.Resource.Target.AverageUtilization != 80 {
		t.Fatalf("avgUtil=%v, want 80", m.Resource.Target.AverageUtilization)
	}
}

// 회귀 가드: valkey 의 (CPU 70 + Mem optional) metric 조립이 헬퍼로 등가 재현되는지.
func TestBuild_ValkeyEquivalentMetrics(t *testing.T) {
	t.Parallel()
	metrics := []autoscalingv2.MetricSpec{hpa.CPUUtilization(70), hpa.MemoryUtilization(75)}
	got := hpa.Build(hpa.Params{
		TargetKind: "StatefulSet", TargetName: "v", MinReplicas: 2, MaxReplicas: 4, MinFloor: 2,
		Metrics: metrics,
	})
	if len(got.Spec.Metrics) != 2 {
		t.Fatalf("metric 수=%d, want 2", len(got.Spec.Metrics))
	}
	if got.Spec.Metrics[0].Resource.Name != corev1.ResourceCPU ||
		got.Spec.Metrics[1].Resource.Name != corev1.ResourceMemory {
		t.Fatal("metric 순서/종류 불일치")
	}
}
