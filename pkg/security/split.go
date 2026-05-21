package security

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// RestrictedPodSecurityContext — Pod-level PodSecurity restricted SecurityContext.
//
// Container-level (RestrictedContainer) 와 분리된 Pod-level 필드:
//   - RunAsNonRoot = true
//   - SeccompProfile.Type = RuntimeDefault
//   - FSGroup / FSGroupChangePolicy (선택)
//
// 사용 예:
//
//	pod.Spec.SecurityContext = security.RestrictedPodSecurityContext(
//	    security.WithPodFSGroup(1000),
//	)
//
// Refs: docs/ROADMAP.md 'Pod / Container SecurityContext 분리 helper' (P-B.9.2)
func RestrictedPodSecurityContext(opts ...PodOption) *corev1.PodSecurityContext {
	psc := &corev1.PodSecurityContext{
		RunAsNonRoot: ptr.To(true),
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
	for _, o := range opts {
		o(psc)
	}
	return psc
}

// PodOption — RestrictedPodSecurityContext 의 functional option.
type PodOption func(*corev1.PodSecurityContext)

// WithPodFSGroup — fsGroup 설정. PVC 의 group ownership 적용.
func WithPodFSGroup(gid int64) PodOption {
	return func(psc *corev1.PodSecurityContext) {
		psc.FSGroup = ptr.To(gid)
	}
}

// WithPodRunAsUser — Pod-level runAsUser (모든 container 기본값).
func WithPodRunAsUser(uid int64) PodOption {
	return func(psc *corev1.PodSecurityContext) {
		psc.RunAsUser = ptr.To(uid)
	}
}

// WithPodRunAsGroup — Pod-level runAsGroup.
func WithPodRunAsGroup(gid int64) PodOption {
	return func(psc *corev1.PodSecurityContext) {
		psc.RunAsGroup = ptr.To(gid)
	}
}
