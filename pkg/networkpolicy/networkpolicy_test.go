// SPDX-License-Identifier: MIT

package networkpolicy

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

func TestNew_DenyByDefault(t *testing.T) {
	t.Parallel()
	np := New("test", "default", map[string]string{"app": "myapp"})
	if np.Name != "test" || np.Namespace != "default" {
		t.Errorf("metadata mismatch: %+v", np.ObjectMeta)
	}
	if got := np.Spec.PodSelector.MatchLabels["app"]; got != "myapp" {
		t.Errorf("podSelector mismatch: %v", np.Spec.PodSelector.MatchLabels)
	}
	if len(np.Spec.PolicyTypes) != 1 || np.Spec.PolicyTypes[0] != networkingv1.PolicyTypeIngress {
		t.Errorf("PolicyTypes default = [Ingress], got %v", np.Spec.PolicyTypes)
	}
	if len(np.Spec.Ingress) != 0 {
		t.Errorf("default Ingress should be empty (deny-all), got %d rules", len(np.Spec.Ingress))
	}
}

func TestWithLabels(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"},
		WithLabels(map[string]string{"team": "data", "env": "prod"}))
	if np.Labels["team"] != "data" || np.Labels["env"] != "prod" {
		t.Errorf("labels mismatch: %v", np.Labels)
	}
}

func TestWithDenyEgress(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"}, WithDenyEgress())
	hasEgress := false
	for _, t2 := range np.Spec.PolicyTypes {
		if t2 == networkingv1.PolicyTypeEgress {
			hasEgress = true
		}
	}
	if !hasEgress {
		t.Errorf("WithDenyEgress should add Egress to PolicyTypes, got %v", np.Spec.PolicyTypes)
	}
}

func TestWithDenyEgress_Idempotent(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"}, WithDenyEgress(), WithDenyEgress())
	count := 0
	for _, t2 := range np.Spec.PolicyTypes {
		if t2 == networkingv1.PolicyTypeEgress {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Egress in PolicyTypes should appear once, got %d", count)
	}
}

func TestWithSelfIngress(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"app": "myapp"},
		WithSelfIngress([]int32{6379, 16379}))
	if len(np.Spec.Ingress) != 1 {
		t.Fatalf("expected 1 ingress rule, got %d", len(np.Spec.Ingress))
	}
	rule := np.Spec.Ingress[0]
	if len(rule.From) != 1 {
		t.Fatalf("expected 1 from peer, got %d", len(rule.From))
	}
	peer := rule.From[0]
	if peer.PodSelector == nil || peer.PodSelector.MatchLabels["app"] != "myapp" {
		t.Errorf("self-peer podSelector mismatch: %+v", peer.PodSelector)
	}
	if len(rule.Ports) != 2 {
		t.Errorf("expected 2 ports (6379, 16379), got %d", len(rule.Ports))
	}
	for _, p := range rule.Ports {
		if p.Protocol == nil || *p.Protocol != corev1.ProtocolTCP {
			t.Errorf("port protocol must be TCP, got %v", p.Protocol)
		}
	}
}

func TestWithSelfIngress_EmptyPorts(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"}, WithSelfIngress(nil))
	if len(np.Spec.Ingress) != 0 {
		t.Errorf("empty ports → no rule, got %d rules", len(np.Spec.Ingress))
	}
}

func TestWithIngressFromPeers(t *testing.T) {
	t.Parallel()
	peers := []Peer{
		{PodSelector: map[string]string{"role": "frontend"}},
		{NamespaceSelector: map[string]string{"team": "data"}},
		{
			PodSelector:       map[string]string{"role": "monitor"},
			NamespaceSelector: map[string]string{"kubernetes.io/metadata.name": "observability"},
		},
	}
	np := New("t", "ns", map[string]string{"app": "v"},
		WithIngressFromPeers(peers, []int32{6379}))
	if len(np.Spec.Ingress) != 1 {
		t.Fatalf("expected 1 ingress rule")
	}
	if len(np.Spec.Ingress[0].From) != 3 {
		t.Errorf("expected 3 peers, got %d", len(np.Spec.Ingress[0].From))
	}
	// 3번째 peer 의 PodSelector + NamespaceSelector 모두 set 검증.
	third := np.Spec.Ingress[0].From[2]
	if third.PodSelector == nil || third.NamespaceSelector == nil {
		t.Errorf("3rd peer should have both selectors: %+v", third)
	}
}

func TestWithIngressFromPeers_EmptySkipped(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"},
		WithIngressFromPeers(nil, []int32{6379}),
		WithIngressFromPeers([]Peer{{PodSelector: map[string]string{"x": "y"}}}, nil))
	if len(np.Spec.Ingress) != 0 {
		t.Errorf("empty peers / empty ports → no rule, got %d", len(np.Spec.Ingress))
	}
}

func TestWithEgressToPeers(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"},
		WithDenyEgress(),
		WithEgressToPeers([]Peer{
			{PodSelector: map[string]string{"role": "db"}},
			{NamespaceSelector: map[string]string{"kubernetes.io/metadata.name": "kube-system"}},
		}, []int32{5432, 53}))
	if len(np.Spec.Egress) != 1 {
		t.Fatalf("expected 1 egress rule")
	}
	if len(np.Spec.Egress[0].To) != 2 {
		t.Fatalf("expected 2 egress peers")
	}
	if np.Spec.Egress[0].To[0].PodSelector.MatchLabels["role"] != "db" {
		t.Errorf("egress peer 1 podSelector mismatch")
	}
	if np.Spec.Egress[0].To[1].NamespaceSelector == nil ||
		np.Spec.Egress[0].To[1].NamespaceSelector.MatchLabels["kubernetes.io/metadata.name"] != "kube-system" {
		t.Errorf("egress peer 2 namespaceSelector mismatch")
	}
}

func TestWithEgressToPeers_EmptySkipped(t *testing.T) {
	t.Parallel()
	np := New("t", "ns", map[string]string{"k": "v"},
		WithEgressToPeers(nil, []int32{1}),
		WithEgressToPeers([]Peer{{PodSelector: map[string]string{"x": "y"}}}, nil))
	if len(np.Spec.Egress) != 0 {
		t.Errorf("empty → no egress rule, got %d", len(np.Spec.Egress))
	}
}

func TestCombination_FullPolicy(t *testing.T) {
	t.Parallel()
	np := New("myapp-cluster", "data", map[string]string{"app": "myapp"},
		WithLabels(map[string]string{"managed-by": "downstream-operator"}),
		WithSelfIngress([]int32{6379, 16379}),
		WithIngressFromPeers(
			[]Peer{{NamespaceSelector: map[string]string{"kubernetes.io/metadata.name": "observability"}}},
			[]int32{9121}),
		WithDenyEgress(),
	)
	// Ingress rule = 2 (self + observability).
	if len(np.Spec.Ingress) != 2 {
		t.Errorf("expected 2 ingress rules, got %d", len(np.Spec.Ingress))
	}
	// PolicyTypes = [Ingress, Egress].
	if len(np.Spec.PolicyTypes) != 2 {
		t.Errorf("expected 2 PolicyTypes, got %v", np.Spec.PolicyTypes)
	}
	// labels.
	if np.Labels["managed-by"] != "downstream-operator" {
		t.Errorf("labels missing")
	}
}
