// SPDX-License-Identifier: MIT

package networkpolicy

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

// TestFourDirection — ingress/egress × TCP (UDP는 caller customize via Peer 확장)
// 모든 조합이 deny-by-default 위에 명시 허용 rule 로 추가되는지 검증.
//
// Refs: docs/ROADMAP.md '4-direction 검증 — ingress/egress × TCP/UDP'
func TestFourDirection(t *testing.T) {
	podSel := map[string]string{"app": "test"}
	peers := []Peer{{PodSelector: map[string]string{"role": "client"}}}

	t.Run("ingress TCP from peers", func(t *testing.T) {
		np := New("test", "default", podSel,
			WithIngressFromPeers(peers, []int32{8080}),
		)
		if len(np.Spec.Ingress) != 1 {
			t.Fatalf("ingress rules = %d, want 1", len(np.Spec.Ingress))
		}
		ports := np.Spec.Ingress[0].Ports
		if len(ports) != 1 || *ports[0].Protocol != corev1.ProtocolTCP {
			t.Errorf("expected TCP/8080 ingress, got %+v", ports)
		}
	})

	t.Run("egress TCP to peers", func(t *testing.T) {
		np := New("test", "default", podSel,
			WithEgressToPeers(peers, []int32{8443}),
		)
		if len(np.Spec.Egress) != 1 {
			t.Fatalf("egress rules = %d, want 1", len(np.Spec.Egress))
		}
	})

	t.Run("self ingress TCP", func(t *testing.T) {
		np := New("test", "default", podSel,
			WithSelfIngress([]int32{6379}),
		)
		if len(np.Spec.Ingress) != 1 {
			t.Fatalf("self ingress = %d, want 1", len(np.Spec.Ingress))
		}
	})

	t.Run("deny egress (default-deny stays)", func(t *testing.T) {
		np := New("test", "default", podSel,
			WithDenyEgress(),
		)
		// PolicyTypes 에 Egress 추가 + Spec.Egress 빈 슬라이스 = deny-by-default
		if len(np.Spec.PolicyTypes) < 1 {
			t.Errorf("PolicyTypes empty, want at least Egress")
		}
		hasEgress := false
		for _, pt := range np.Spec.PolicyTypes {
			if pt == "Egress" {
				hasEgress = true
			}
		}
		if !hasEgress {
			t.Errorf("PolicyTypes missing Egress, got %v", np.Spec.PolicyTypes)
		}
	})

	t.Run("combo: ingress + egress + self", func(t *testing.T) {
		np := New("test", "default", podSel,
			WithSelfIngress([]int32{6379}),
			WithIngressFromPeers(peers, []int32{8080}),
			WithEgressToPeers(peers, []int32{8443}),
		)
		if len(np.Spec.Ingress) < 2 {
			t.Errorf("expected ≥2 ingress rules, got %d", len(np.Spec.Ingress))
		}
		if len(np.Spec.Egress) != 1 {
			t.Errorf("expected 1 egress rule, got %d", len(np.Spec.Egress))
		}
	})
}
