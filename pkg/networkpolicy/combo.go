package networkpolicy

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ComboPeer — CIDR + NamespaceSelector + PodSelector 조합 peer.
//
// 사용 예 (cross-cluster CIDR + 같은 namespace 안 특정 pod):
//
//	peer := networkpolicy.ComboPeer{
//	    CIDR:              "10.0.0.0/8",
//	    Except:            []string{"10.0.0.1/32"},
//	    NamespaceSelector: map[string]string{"trusted": "true"},
//	    PodSelector:       map[string]string{"app": "monitoring"},
//	}
//
// Refs: ROADMAP.md 'CIDR + namespace selector + pod selector 조합 helper'
//       (P-B.8.2)
type ComboPeer struct {
	CIDR              string
	Except            []string
	NamespaceSelector map[string]string
	PodSelector       map[string]string
}

// ToNetworkPolicyPeer — ComboPeer 를 NetworkPolicy peer 로 변환.
func (c ComboPeer) ToNetworkPolicyPeer() networkingv1.NetworkPolicyPeer {
	out := networkingv1.NetworkPolicyPeer{}
	if c.CIDR != "" {
		out.IPBlock = &networkingv1.IPBlock{CIDR: c.CIDR, Except: c.Except}
	}
	if len(c.NamespaceSelector) > 0 {
		out.NamespaceSelector = &metav1.LabelSelector{MatchLabels: c.NamespaceSelector}
	}
	if len(c.PodSelector) > 0 {
		out.PodSelector = &metav1.LabelSelector{MatchLabels: c.PodSelector}
	}
	return out
}

// WithComboIngressFromPeers — ComboPeer 리스트 + TCP ports 로 Ingress rule 추가.
func WithComboIngressFromPeers(peers []ComboPeer, tcpPorts []int32) Option {
	return func(b *builder) {
		if len(peers) == 0 || len(tcpPorts) == 0 {
			return
		}
		policyPeers := make([]networkingv1.NetworkPolicyPeer, 0, len(peers))
		for _, p := range peers {
			policyPeers = append(policyPeers, p.ToNetworkPolicyPeer())
		}
		b.ingress = append(b.ingress, networkingv1.NetworkPolicyIngressRule{
			From:  policyPeers,
			Ports: tcpPortsToPolicy(tcpPorts),
		})
	}
}
