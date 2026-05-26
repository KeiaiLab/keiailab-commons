// SPDX-License-Identifier: Apache-2.0

package security

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

// TestRestrictedPSACompliance — PodSecurity "restricted" 정책의 의무 가드를
// RestrictedContainer helper 가 충족하는지 검증.
//
// 다수 downstream repo 가 본 helper 를 직접 사용하므로
// regression 가드 = downstream 정합 가드.
//
// PSA restricted container-level 의무 (K8s 1.25+):
//   - capabilities.drop = ["ALL"]
//   - allowPrivilegeEscalation = false
//   - runAsNonRoot = true
//   - seccompProfile.type = RuntimeDefault
//
// (Pod-level 가드는 batch-5 의 split.go + seccomp.go 머지 후 별 PR 추가)
//
// Refs: docs/ROADMAP.md 'restricted PSA downstream 회귀 가드'
func TestRestrictedPSACompliance(t *testing.T) {
	t.Run("container restricted defaults", func(t *testing.T) {
		sc := RestrictedContainer()
		if sc.Capabilities == nil {
			t.Fatalf("Capabilities must be set")
		}
		dropFound := false
		for _, c := range sc.Capabilities.Drop {
			if c == "ALL" {
				dropFound = true
			}
		}
		if !dropFound {
			t.Errorf("Capabilities.Drop must contain 'ALL', got %v", sc.Capabilities.Drop)
		}
		if sc.AllowPrivilegeEscalation == nil || *sc.AllowPrivilegeEscalation {
			t.Errorf("AllowPrivilegeEscalation must be false")
		}
		if sc.RunAsNonRoot == nil || !*sc.RunAsNonRoot {
			t.Errorf("RunAsNonRoot must be true")
		}
		if sc.SeccompProfile == nil || sc.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			t.Errorf("SeccompProfile must be RuntimeDefault, got %v", sc.SeccompProfile)
		}
	})

	t.Run("container with runAsUser option", func(t *testing.T) {
		sc := RestrictedContainer(WithRunAsUser(999))
		// 추가 option 적용 시에도 PSA restricted 의무 유지 검증
		if sc.AllowPrivilegeEscalation == nil || *sc.AllowPrivilegeEscalation {
			t.Errorf("AllowPrivilegeEscalation must remain false with options")
		}
		if sc.RunAsUser == nil || *sc.RunAsUser != 999 {
			t.Errorf("RunAsUser = %v, want 999", sc.RunAsUser)
		}
	})
}
