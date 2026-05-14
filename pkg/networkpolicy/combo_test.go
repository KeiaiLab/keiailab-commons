package networkpolicy

import (
	"testing"
)

func TestComboPeer_ToNetworkPolicyPeer(t *testing.T) {
	t.Run("cidr only", func(t *testing.T) {
		p := ComboPeer{CIDR: "10.0.0.0/8"}
		out := p.ToNetworkPolicyPeer()
		if out.IPBlock == nil || out.IPBlock.CIDR != "10.0.0.0/8" {
			t.Errorf("CIDR not set")
		}
		if out.NamespaceSelector != nil || out.PodSelector != nil {
			t.Errorf("selectors should be nil")
		}
	})
	t.Run("ns + pod selector", func(t *testing.T) {
		p := ComboPeer{
			NamespaceSelector: map[string]string{"trusted": "true"},
			PodSelector:       map[string]string{"app": "monitoring"},
		}
		out := p.ToNetworkPolicyPeer()
		if out.NamespaceSelector == nil || out.NamespaceSelector.MatchLabels["trusted"] != "true" {
			t.Errorf("ns selector wrong")
		}
		if out.PodSelector == nil || out.PodSelector.MatchLabels["app"] != "monitoring" {
			t.Errorf("pod selector wrong")
		}
	})
	t.Run("full combo with except", func(t *testing.T) {
		p := ComboPeer{
			CIDR:              "10.0.0.0/8",
			Except:            []string{"10.0.0.1/32"},
			NamespaceSelector: map[string]string{"a": "b"},
			PodSelector:       map[string]string{"c": "d"},
		}
		out := p.ToNetworkPolicyPeer()
		if out.IPBlock.CIDR != "10.0.0.0/8" || len(out.IPBlock.Except) != 1 {
			t.Errorf("ipblock wrong")
		}
		if out.NamespaceSelector == nil || out.PodSelector == nil {
			t.Errorf("selectors should be set")
		}
	})
}
