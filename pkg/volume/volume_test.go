// SPDX-License-Identifier: MIT

package volume_test

import (
	"testing"

	"github.com/keiailab/keiailab-commons/pkg/volume"
)

func TestTLSSecretMount(t *testing.T) {
	t.Parallel()
	vol, mount := volume.TLSSecretMount("tls", "db-tls", "/tls")

	if vol.Name != "tls" || mount.Name != "tls" {
		t.Fatalf("volume/mount 이름 불일치: %s / %s", vol.Name, mount.Name)
	}
	if vol.Secret == nil || vol.Secret.SecretName != "db-tls" {
		t.Fatalf("secret 이름 불일치: %+v", vol.Secret)
	}
	if vol.Secret.DefaultMode == nil || *vol.Secret.DefaultMode != 0o400 {
		t.Fatalf("DefaultMode=%v, want 0o400 (cert 불변식)", vol.Secret.DefaultMode)
	}
	if mount.MountPath != "/tls" {
		t.Fatalf("mountPath=%q, want /tls", mount.MountPath)
	}
	if !mount.ReadOnly {
		t.Fatal("cert mount 는 ReadOnly=true 여야 함")
	}
}

func TestTLSSecretMount_CustomName(t *testing.T) {
	t.Parallel()
	// mongo 식 이름/경로도 동일 불변식.
	vol, mount := volume.TLSSecretMount("tls-server", "mdb-tls", "/etc/ssl/mongo")
	if vol.Name != "tls-server" || mount.MountPath != "/etc/ssl/mongo" {
		t.Fatalf("custom 이름/경로 전달 실패: %s / %s", vol.Name, mount.MountPath)
	}
	if *vol.Secret.DefaultMode != volume.CertDefaultMode || !mount.ReadOnly {
		t.Fatal("불변식(0o400/readonly) 미적용")
	}
}
