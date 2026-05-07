// Package security — PodSecurity "restricted" 정책 충족 SecurityContext builder.
//
// PodSecurity restricted (https://kubernetes.io/docs/concepts/security/pod-security-standards/)
// 의 모든 *항상 적용* 가드를 단일 진실원으로 통합:
//
//   - capabilities.drop = ["ALL"]
//   - seccompProfile.type = RuntimeDefault
//   - allowPrivilegeEscalation = false
//   - runAsNonRoot = true
//
// 옵션 (operator 별 차이 — RunAsUser/Group, ReadOnlyRootFilesystem) 은 functional
// options 패턴으로 노출.
//
// 사용 예:
//
//	sc := security.RestrictedContainer(security.WithRunAsUser(999), security.WithRunAsGroup(999))
package security

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// Option — RestrictedContainer 의 functional option.
type Option func(*corev1.SecurityContext)

// WithRunAsUser — runAsUser 명시. 미지정 시 image 의 USER 디렉티브 따름 (단,
// runAsNonRoot=true 이므로 root 는 admission 거부).
func WithRunAsUser(uid int64) Option {
	return func(sc *corev1.SecurityContext) {
		sc.RunAsUser = ptr.To(uid)
	}
}

// WithRunAsGroup — runAsGroup 명시.
func WithRunAsGroup(gid int64) Option {
	return func(sc *corev1.SecurityContext) {
		sc.RunAsGroup = ptr.To(gid)
	}
}

// WithReadOnlyRootFilesystem — true 시 root fs 쓰기 금지 (PodSecurity 에는
// 포함 안 되나 hardening 권장).
func WithReadOnlyRootFilesystem(readOnly bool) Option {
	return func(sc *corev1.SecurityContext) {
		sc.ReadOnlyRootFilesystem = ptr.To(readOnly)
	}
}

// RestrictedContainer — PodSecurity "restricted" 충족 *최소* SecurityContext.
// 모든 invariant (cap drop / seccomp / allowPrivilegeEscalation / runAsNonRoot)
// 는 항상 적용. 변경 가능 옵션은 functional option 으로.
func RestrictedContainer(opts ...Option) *corev1.SecurityContext {
	sc := &corev1.SecurityContext{
		RunAsNonRoot:             ptr.To(true),
		AllowPrivilegeEscalation: ptr.To(false),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

// RestrictedPod — Pod 레벨 SecurityContext (fsGroup / supplementalGroups 등은
// pod-scope only). containerOptions 와 별개로 pod-level invariant 만 보장.
func RestrictedPod(fsGroup *int64) *corev1.PodSecurityContext {
	psc := &corev1.PodSecurityContext{
		RunAsNonRoot: ptr.To(true),
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
	if fsGroup != nil {
		psc.FSGroup = fsGroup
	}
	return psc
}
