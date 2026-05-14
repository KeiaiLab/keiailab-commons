package webhook

import (
	"strings"
	"testing"
)

func TestValidateStorageClassDNS1123(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty (default)", "", false},
		{"valid simple", "fast-ssd", false},
		{"valid with dot", "ceph.rbd-fast", false},
		{"invalid uppercase", "Fast-SSD", true},
		{"invalid leading dash", "-fast", true},
		{"invalid trailing dash", "fast-", true},
		{"invalid special", "fast_ssd", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStorageClassDNS1123(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStorageClassDNS1123(%q) err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateReplicaCountLowerBound(t *testing.T) {
	t.Run("above min", func(t *testing.T) {
		if err := ValidateReplicaCountLowerBound(5, 3, "HA"); err != nil {
			t.Errorf("err = %v, want nil", err)
		}
	})
	t.Run("below min", func(t *testing.T) {
		err := ValidateReplicaCountLowerBound(2, 3, "HA")
		if err == nil || !strings.Contains(err.Error(), "HA") {
			t.Errorf("expected HA error, got %v", err)
		}
	})
	t.Run("equal min", func(t *testing.T) {
		if err := ValidateReplicaCountLowerBound(3, 3, "HA"); err != nil {
			t.Errorf("err = %v, want nil", err)
		}
	})
}

func TestValidateTopologySpreadKeyAllowed(t *testing.T) {
	t.Run("allowed zone", func(t *testing.T) {
		if err := ValidateTopologySpreadKeyAllowed("topology.kubernetes.io/zone"); err != nil {
			t.Errorf("zone err = %v", err)
		}
	})
	t.Run("disallowed custom", func(t *testing.T) {
		err := ValidateTopologySpreadKeyAllowed("custom-label")
		if err == nil {
			t.Errorf("expected error for custom key")
		}
	})
}
