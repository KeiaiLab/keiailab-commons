// SPDX-License-Identifier: MIT

// Package networkpolicy — Kubernetes NetworkPolicy builder.
//
// downstream operator 의 공통 패턴:
//   - deny-by-default (PolicyTypes 명시)
//   - 같은 인스턴스 pod 간 ingress 허용 (self-peer)
//   - 추가 peers (PodSelector / NamespaceSelector matchLabels)
//   - 명시 ports (TCP)
//
// 사용 예:
//
//	np := networkpolicy.New("vk-cluster", "default",
//	    map[string]string{"app.kubernetes.io/instance": "vk-cluster"},
//	    networkpolicy.WithLabels(myLabels),
//	    networkpolicy.WithSelfIngress([]int32{6379, 16379}),
//	    networkpolicy.WithIngressFromPeers(extraPeers, []int32{6379}),
//	)
package networkpolicy

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Peer — NetworkPolicy ingress / egress peer (PodSelector + NamespaceSelector).
// 두 selector 모두 nil 인 peer 는 caller 책임 (의미 없음).
type Peer struct {
	PodSelector       map[string]string
	NamespaceSelector map[string]string
}

// builder — internal accumulator. New() 가 초기화, 옵션이 mutate.
type builder struct {
	np      *networkingv1.NetworkPolicy
	ingress []networkingv1.NetworkPolicyIngressRule
	egress  []networkingv1.NetworkPolicyEgressRule
}

// Option — New() 의 functional option.
type Option func(*builder)

// New — deny-by-default NetworkPolicy. 옵션으로 ingress / egress / labels 추가.
//
// 기본 동작 (옵션 미적용 시): PolicyTypes=[Ingress], Ingress=[] (=== deny all
// ingress). Egress 는 미명시 (Kubernetes default = allow all egress).
//
// WithDenyEgress() 옵션 추가 시 PolicyTypes=[Ingress, Egress] + Egress=[]
// (=== deny all egress except WithEgress* 로 명시).
func New(name, namespace string, podSelector map[string]string, opts ...Option) *networkingv1.NetworkPolicy {
	b := &builder{
		np: &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{MatchLabels: podSelector},
				PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
			},
		},
	}
	for _, opt := range opts {
		opt(b)
	}
	b.np.Spec.Ingress = b.ingress
	b.np.Spec.Egress = b.egress
	return b.np
}

// WithLabels — NetworkPolicy metadata.labels 설정.
func WithLabels(labels map[string]string) Option {
	return func(b *builder) {
		b.np.Labels = labels
	}
}

// WithDenyEgress — PolicyTypes 에 Egress 추가 (deny-by-default egress).
// 명시 egress 는 WithEgressTo* 옵션으로.
func WithDenyEgress() Option {
	return func(b *builder) {
		hasEgress := false
		for _, t := range b.np.Spec.PolicyTypes {
			if t == networkingv1.PolicyTypeEgress {
				hasEgress = true
			}
		}
		if !hasEgress {
			b.np.Spec.PolicyTypes = append(b.np.Spec.PolicyTypes, networkingv1.PolicyTypeEgress)
		}
	}
}

// WithIngressFromPeers — 명시 peers 로부터의 TCP ports ingress 허용.
func WithIngressFromPeers(peers []Peer, tcpPorts []int32) Option {
	return func(b *builder) {
		if len(peers) == 0 || len(tcpPorts) == 0 {
			return
		}
		from := make([]networkingv1.NetworkPolicyPeer, 0, len(peers))
		for _, p := range peers {
			peer := networkingv1.NetworkPolicyPeer{}
			if len(p.PodSelector) > 0 {
				peer.PodSelector = &metav1.LabelSelector{MatchLabels: p.PodSelector}
			}
			if len(p.NamespaceSelector) > 0 {
				peer.NamespaceSelector = &metav1.LabelSelector{MatchLabels: p.NamespaceSelector}
			}
			from = append(from, peer)
		}
		b.ingress = append(b.ingress, networkingv1.NetworkPolicyIngressRule{
			From:  from,
			Ports: tcpPortsToPolicy(tcpPorts),
		})
	}
}

// WithSelfIngress — 같은 PodSelector 가 가리키는 pod 들 간 ingress 허용 (TCP ports).
func WithSelfIngress(tcpPorts []int32) Option {
	return func(b *builder) {
		if len(tcpPorts) == 0 {
			return
		}
		selfPeer := networkingv1.NetworkPolicyPeer{
			PodSelector: &metav1.LabelSelector{MatchLabels: b.np.Spec.PodSelector.MatchLabels},
		}
		b.ingress = append(b.ingress, networkingv1.NetworkPolicyIngressRule{
			From:  []networkingv1.NetworkPolicyPeer{selfPeer},
			Ports: tcpPortsToPolicy(tcpPorts),
		})
	}
}

// WithEgressToPeers — 명시 peers 로의 TCP ports egress 허용. WithDenyEgress 와
// 함께 사용 권장 (그렇지 않으면 PolicyTypes 에 Egress 가 없어 무시됨).
func WithEgressToPeers(peers []Peer, tcpPorts []int32) Option {
	return func(b *builder) {
		if len(peers) == 0 || len(tcpPorts) == 0 {
			return
		}
		to := make([]networkingv1.NetworkPolicyPeer, 0, len(peers))
		for _, p := range peers {
			peer := networkingv1.NetworkPolicyPeer{}
			if len(p.PodSelector) > 0 {
				peer.PodSelector = &metav1.LabelSelector{MatchLabels: p.PodSelector}
			}
			if len(p.NamespaceSelector) > 0 {
				peer.NamespaceSelector = &metav1.LabelSelector{MatchLabels: p.NamespaceSelector}
			}
			to = append(to, peer)
		}
		b.egress = append(b.egress, networkingv1.NetworkPolicyEgressRule{
			To:    to,
			Ports: tcpPortsToPolicy(tcpPorts),
		})
	}
}

func tcpPortsToPolicy(ports []int32) []networkingv1.NetworkPolicyPort {
	tcp := corev1.ProtocolTCP
	out := make([]networkingv1.NetworkPolicyPort, 0, len(ports))
	for _, p := range ports {
		v := intstr.FromInt(int(p))
		out = append(out, networkingv1.NetworkPolicyPort{
			Protocol: &tcp,
			Port:     &v,
		})
	}
	return out
}
