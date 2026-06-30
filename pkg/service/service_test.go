// SPDX-License-Identifier: MIT

package service_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/keiailab/keiailab-commons/pkg/service"
)

func port(name string, p int32) corev1.ServicePort {
	return corev1.ServicePort{Name: name, Port: p, TargetPort: intstr.FromInt(int(p)), Protocol: corev1.ProtocolTCP}
}

func TestBuild_Headless(t *testing.T) {
	t.Parallel()
	got := service.Build(service.Params{
		Name: "db-headless", Namespace: "ns",
		Labels: map[string]string{"k": "v"}, Selector: map[string]string{"s": "1"},
		Ports: []corev1.ServicePort{port("client", 6379)}, Headless: true,
		Type: corev1.ServiceTypeNodePort, // Headless 시 무시돼야 함
	})
	if got.Spec.ClusterIP != corev1.ClusterIPNone {
		t.Fatalf("headless ClusterIP=%q, want None", got.Spec.ClusterIP)
	}
	if !got.Spec.PublishNotReadyAddresses {
		t.Fatal("headless 는 PublishNotReadyAddresses=true 여야 함")
	}
	if got.Spec.Type != "" {
		t.Fatalf("headless 는 Type 미설정이어야 함, got %q", got.Spec.Type)
	}
	if got.Name != "db-headless" || got.Labels["k"] != "v" || got.Spec.Selector["s"] != "1" {
		t.Fatalf("meta/selector mismatch: %+v", got.ObjectMeta)
	}
}

func TestBuild_ClientDefaultsToClusterIP(t *testing.T) {
	t.Parallel()
	got := service.Build(service.Params{Name: "db", Ports: []corev1.ServicePort{port("client", 27017)}})
	if got.Spec.Type != corev1.ServiceTypeClusterIP {
		t.Fatalf("기본 Type=%q, want ClusterIP", got.Spec.Type)
	}
	if got.Spec.ClusterIP == corev1.ClusterIPNone {
		t.Fatal("non-headless 는 ClusterIP None 이면 안 됨")
	}
	if got.Spec.PublishNotReadyAddresses {
		t.Fatal("non-headless 기본 PublishNotReadyAddresses=false")
	}
}

func TestBuild_TypeOverrideAndIPFamilies(t *testing.T) {
	t.Parallel()
	pol := corev1.IPFamilyPolicyPreferDualStack
	got := service.Build(service.Params{
		Name: "db", Type: corev1.ServiceTypeLoadBalancer,
		IPFamilyPolicy: &pol, IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
		Annotations: map[string]string{"a": "b"},
	})
	if got.Spec.Type != corev1.ServiceTypeLoadBalancer {
		t.Fatalf("Type override 무시: %q", got.Spec.Type)
	}
	if got.Spec.IPFamilyPolicy == nil || *got.Spec.IPFamilyPolicy != corev1.IPFamilyPolicyPreferDualStack {
		t.Fatalf("IPFamilyPolicy 전달 실패: %v", got.Spec.IPFamilyPolicy)
	}
	if len(got.Spec.IPFamilies) != 2 {
		t.Fatalf("IPFamilies 전달 실패: %v", got.Spec.IPFamilies)
	}
	if got.Annotations["a"] != "b" {
		t.Fatalf("annotations 전달 실패: %v", got.Annotations)
	}
}

func TestBuild_NonHeadlessPublishNotReady(t *testing.T) {
	t.Parallel()
	got := service.Build(service.Params{Name: "x", PublishNotReadyAddresses: true})
	if !got.Spec.PublishNotReadyAddresses {
		t.Fatal("명시 PublishNotReadyAddresses=true 가 무시됨")
	}
}

// 회귀 가드: valkey headless(ClusterIP None + PublishNotReady + IPFamilies) 등가성.
func TestBuild_ValkeyHeadlessEquivalent(t *testing.T) {
	t.Parallel()
	pol := corev1.IPFamilyPolicySingleStack
	got := service.Build(service.Params{
		Name: "cache-headless", Namespace: "ns", Headless: true,
		Selector:       map[string]string{"app.kubernetes.io/instance": "cache"},
		Ports:          []corev1.ServicePort{port("client", 6379), port("cluster-bus", 16379)},
		IPFamilyPolicy: &pol,
	})
	if got.Spec.ClusterIP != corev1.ClusterIPNone || !got.Spec.PublishNotReadyAddresses {
		t.Fatal("valkey headless 등가성 실패")
	}
	if len(got.Spec.Ports) != 2 || got.Spec.IPFamilyPolicy == nil {
		t.Fatal("ports/IPFamilyPolicy 전달 실패")
	}
}
