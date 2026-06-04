// SPDX-License-Identifier: MIT

package security

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestRestrictedContainer_Defaults(t *testing.T) {
	t.Parallel()
	sc := RestrictedContainer()
	if sc.RunAsNonRoot == nil || !*sc.RunAsNonRoot {
		t.Error("RunAsNonRoot must be true")
	}
	if sc.AllowPrivilegeEscalation == nil || *sc.AllowPrivilegeEscalation {
		t.Error("AllowPrivilegeEscalation must be false")
	}
	if sc.Capabilities == nil || len(sc.Capabilities.Drop) != 1 || sc.Capabilities.Drop[0] != "ALL" {
		t.Errorf("Capabilities.Drop must be [ALL], got %v", sc.Capabilities)
	}
	if sc.SeccompProfile == nil || sc.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Errorf("SeccompProfile must be RuntimeDefault, got %v", sc.SeccompProfile)
	}
	// Defaults 에서는 RunAsUser / RunAsGroup / ReadOnlyRootFilesystem 미설정.
	if sc.RunAsUser != nil {
		t.Errorf("RunAsUser should be nil by default, got %v", *sc.RunAsUser)
	}
	if sc.RunAsGroup != nil {
		t.Errorf("RunAsGroup should be nil by default, got %v", *sc.RunAsGroup)
	}
}

func TestRestrictedContainer_WithRunAsUser(t *testing.T) {
	t.Parallel()
	sc := RestrictedContainer(WithRunAsUser(999))
	if sc.RunAsUser == nil || *sc.RunAsUser != 999 {
		t.Errorf("RunAsUser = %v, want 999", sc.RunAsUser)
	}
}

func TestRestrictedContainer_WithRunAsGroup(t *testing.T) {
	t.Parallel()
	sc := RestrictedContainer(WithRunAsGroup(1000))
	if sc.RunAsGroup == nil || *sc.RunAsGroup != 1000 {
		t.Errorf("RunAsGroup = %v, want 1000", sc.RunAsGroup)
	}
}

func TestRestrictedContainer_WithReadOnlyRootFilesystem(t *testing.T) {
	t.Parallel()
	sc := RestrictedContainer(WithReadOnlyRootFilesystem(true))
	if sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
		t.Error("ReadOnlyRootFilesystem must be true when option specified")
	}
	// Invariants 유지 검증.
	if !*sc.RunAsNonRoot || *sc.AllowPrivilegeEscalation {
		t.Error("invariants broken by option application")
	}
}

func TestRestrictedContainer_AllOptionsCombined(t *testing.T) {
	t.Parallel()
	sc := RestrictedContainer(
		WithRunAsUser(999),
		WithRunAsGroup(999),
		WithReadOnlyRootFilesystem(false),
	)
	if *sc.RunAsUser != 999 || *sc.RunAsGroup != 999 || *sc.ReadOnlyRootFilesystem {
		t.Errorf("options not applied correctly: %+v", sc)
	}
}

func TestRestrictedPod_NoFSGroup(t *testing.T) {
	t.Parallel()
	psc := RestrictedPod(nil)
	if psc.RunAsNonRoot == nil || !*psc.RunAsNonRoot {
		t.Error("RunAsNonRoot must be true")
	}
	if psc.SeccompProfile == nil || psc.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Error("SeccompProfile must be RuntimeDefault")
	}
	if psc.FSGroup != nil {
		t.Errorf("FSGroup should be nil when fsGroup arg nil, got %v", *psc.FSGroup)
	}
}

func TestRestrictedPod_WithFSGroup(t *testing.T) {
	t.Parallel()
	psc := RestrictedPod(ptr.To[int64](2000))
	if psc.FSGroup == nil || *psc.FSGroup != 2000 {
		t.Errorf("FSGroup = %v, want 2000", psc.FSGroup)
	}
}
