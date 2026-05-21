// SPDX-License-Identifier: Apache-2.0

package security

import corev1 "k8s.io/api/core/v1"

// RuntimeDefaultSeccompProfile — `RuntimeDefault` seccomp profile pointer.
//
// PodSecurity restricted 의 의무 seccomp 설정. Pod 와 Container 양쪽에
// 적용 가능.
//
// 사용 예:
//
//	pod.Spec.SecurityContext.SeccompProfile = security.RuntimeDefaultSeccompProfile()
//	container.SecurityContext.SeccompProfile = security.RuntimeDefaultSeccompProfile()
//
// Refs: docs/ROADMAP.md 'seccompProfile 기본값 helper'
func RuntimeDefaultSeccompProfile() *corev1.SeccompProfile {
	return &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeRuntimeDefault}
}

// LocalhostSeccompProfile — `Localhost` profile (custom seccomp.json path).
//
// 사용 시 노드의 `/var/lib/kubelet/seccomp/<localhostProfile>` 경로에
// 프로파일 파일이 *반드시* 존재해야 한다 (kubelet 이 로드).
func LocalhostSeccompProfile(localhostProfile string) *corev1.SeccompProfile {
	return &corev1.SeccompProfile{
		Type:             corev1.SeccompProfileTypeLocalhost,
		LocalhostProfile: &localhostProfile,
	}
}

// UnconfinedSeccompProfile — `Unconfined` profile.
//
// 주의: restricted PodSecurity 정책에서 *거부됨*. 디버깅 / 특수 워크로드
// 한정. 운영 클러스터에서는 사용 금지.
func UnconfinedSeccompProfile() *corev1.SeccompProfile {
	return &corev1.SeccompProfile{Type: corev1.SeccompProfileTypeUnconfined}
}
