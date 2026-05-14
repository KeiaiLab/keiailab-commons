package security

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestRestrictedPodSecurityContext(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		psc := RestrictedPodSecurityContext()
		if psc.RunAsNonRoot == nil || !*psc.RunAsNonRoot {
			t.Errorf("RunAsNonRoot must be true")
		}
		if psc.SeccompProfile == nil || psc.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
			t.Errorf("SeccompProfile.Type must be RuntimeDefault")
		}
	})
	t.Run("with fsgroup", func(t *testing.T) {
		psc := RestrictedPodSecurityContext(WithPodFSGroup(1000))
		if psc.FSGroup == nil || *psc.FSGroup != 1000 {
			t.Errorf("FSGroup = %v, want 1000", psc.FSGroup)
		}
	})
	t.Run("with run as", func(t *testing.T) {
		psc := RestrictedPodSecurityContext(WithPodRunAsUser(999), WithPodRunAsGroup(999))
		if *psc.RunAsUser != 999 || *psc.RunAsGroup != 999 {
			t.Errorf("RunAsUser/Group not set")
		}
	})
}

func TestSeccompProfiles(t *testing.T) {
	rd := RuntimeDefaultSeccompProfile()
	if rd.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Errorf("RuntimeDefault type wrong")
	}
	local := LocalhostSeccompProfile("custom.json")
	if local.Type != corev1.SeccompProfileTypeLocalhost {
		t.Errorf("Localhost type wrong")
	}
	if local.LocalhostProfile == nil || *local.LocalhostProfile != "custom.json" {
		t.Errorf("LocalhostProfile = %v, want custom.json", local.LocalhostProfile)
	}
	unc := UnconfinedSeccompProfile()
	if unc.Type != corev1.SeccompProfileTypeUnconfined {
		t.Errorf("Unconfined type wrong")
	}
}
