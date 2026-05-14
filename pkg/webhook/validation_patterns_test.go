package webhook

import "testing"

func TestValidateStorageClassDNS1123(t *testing.T) {
	cases := []struct{ in string; wantErr bool }{
		{"", false}, {"fast-ssd", false}, {"Fast-SSD", true}, {"-x", true}, {"x_y", true},
	}
	for _, c := range cases {
		err := ValidateStorageClassDNS1123(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("ValidateStorageClassDNS1123(%q) err=%v wantErr=%v", c.in, err, c.wantErr)
		}
	}
}

func TestValidateReplicaCountLowerBound(t *testing.T) {
	if err := ValidateReplicaCountLowerBound(5, 3, "HA"); err != nil {
		t.Errorf("err=%v", err)
	}
	if ValidateReplicaCountLowerBound(2, 3, "HA") == nil {
		t.Errorf("expected error")
	}
}

func TestValidateTopologySpreadKeyAllowed(t *testing.T) {
	if err := ValidateTopologySpreadKeyAllowed("topology.kubernetes.io/zone"); err != nil {
		t.Errorf("zone err=%v", err)
	}
	if ValidateTopologySpreadKeyAllowed("custom") == nil {
		t.Errorf("expected error")
	}
}
