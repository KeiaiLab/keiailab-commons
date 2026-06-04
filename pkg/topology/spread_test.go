// SPDX-License-Identifier: MIT

package topology

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestDefaulted_user_provided_preserved(t *testing.T) {
	user := []corev1.TopologySpreadConstraint{{
		MaxSkew:     2,
		TopologyKey: "rack",
	}}
	got := Defaulted(user, 5, map[string]string{"app": "x"})
	if len(got) != 1 || got[0].TopologyKey != "rack" {
		t.Errorf("user TSC overridden: %v", got)
	}
}

func TestDefaulted_replicas_below_min_no_inject(t *testing.T) {
	got := Defaulted(nil, 1, map[string]string{"app": "x"})
	if got != nil {
		t.Errorf("replicas < min → nil expected, got %v", got)
	}
}

func TestDefaulted_replicas_at_min_injects(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"})
	if len(got) != 2 {
		t.Fatalf("expected 2 default TSCs (zone + hostname), got %d", len(got))
	}
	if got[0].TopologyKey != TopologyKeyZone {
		t.Errorf("first TSC: %q want %q", got[0].TopologyKey, TopologyKeyZone)
	}
	if got[1].TopologyKey != TopologyKeyHostname {
		t.Errorf("second TSC: %q want %q", got[1].TopologyKey, TopologyKeyHostname)
	}
	for _, c := range got {
		if c.MaxSkew != DefaultMaxSkew {
			t.Errorf("MaxSkew: %d want %d", c.MaxSkew, DefaultMaxSkew)
		}
		if c.WhenUnsatisfiable != corev1.ScheduleAnyway {
			t.Errorf("WhenUnsatisfiable: %q", c.WhenUnsatisfiable)
		}
	}
}

func TestDefaulted_additional_replicas_pattern_WithMinReplicas_1(t *testing.T) {
	got := Defaulted(nil, 1, map[string]string{"app": "myapp"}, WithMinReplicas(1))
	if len(got) != 2 {
		t.Fatalf("additional replicas pattern (min=1, replicas=1): expected 2 TSCs, got %d", len(got))
	}
	got = Defaulted(nil, 0, map[string]string{"app": "myapp"}, WithMinReplicas(1))
	if got != nil {
		t.Errorf("additional replicas pattern (min=1, replicas=0): nil expected, got %v", got)
	}
}

func TestDefaulted_label_selector_matches(t *testing.T) {
	selector := map[string]string{
		"app.kubernetes.io/name":     "myapp",
		"app.kubernetes.io/instance": "x",
	}
	got := Defaulted(nil, 3, selector)
	for _, c := range got {
		for k, v := range selector {
			if c.LabelSelector.MatchLabels[k] != v {
				t.Errorf("TSC label selector missing %q=%q", k, v)
			}
		}
	}
}

func TestDefaulted_WithTopologyKeys_custom(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"},
		WithTopologyKeys("rack"))
	if len(got) != 1 {
		t.Fatalf("expected 1 TSC, got %d", len(got))
	}
	if got[0].TopologyKey != "rack" {
		t.Errorf("custom key: %q want rack", got[0].TopologyKey)
	}
}

func TestDefaulted_WithTopologyKeys_three_axes(t *testing.T) {
	got := Defaulted(nil, 3, map[string]string{"app": "x"},
		WithTopologyKeys("rack", TopologyKeyZone, TopologyKeyHostname))
	if len(got) != 3 {
		t.Fatalf("expected 3 TSCs, got %d", len(got))
	}
}

func TestDefaulted_WithTopologyKeys_empty_ignored(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"}, WithTopologyKeys())
	if len(got) != 2 {
		t.Fatalf("empty keys should fall back to default 2 axes, got %d", len(got))
	}
}

func TestDefaulted_WithMaxSkew(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"}, WithMaxSkew(2))
	for i, c := range got {
		if c.MaxSkew != 2 {
			t.Errorf("TSC[%d] MaxSkew: %d want 2", i, c.MaxSkew)
		}
	}
}

func TestDefaulted_WithMaxSkew_zero_ignored(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"}, WithMaxSkew(0))
	for i, c := range got {
		if c.MaxSkew != DefaultMaxSkew {
			t.Errorf("TSC[%d] MaxSkew with 0 input: %d want default %d", i, c.MaxSkew, DefaultMaxSkew)
		}
	}
}

func TestDefaulted_WithWhenUnsatisfiable(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"},
		WithWhenUnsatisfiable(corev1.DoNotSchedule))
	for i, c := range got {
		if c.WhenUnsatisfiable != corev1.DoNotSchedule {
			t.Errorf("TSC[%d] WhenUnsatisfiable: %q want DoNotSchedule", i, c.WhenUnsatisfiable)
		}
	}
}

func TestDefaulted_WithWhenUnsatisfiable_empty_ignored(t *testing.T) {
	got := Defaulted(nil, 2, map[string]string{"app": "x"},
		WithWhenUnsatisfiable(""))
	for i, c := range got {
		if c.WhenUnsatisfiable != corev1.ScheduleAnyway {
			t.Errorf("TSC[%d] WhenUnsatisfiable with empty input: %q want default ScheduleAnyway", i, c.WhenUnsatisfiable)
		}
	}
}

func TestDefaulted_label_selector_semantic_match(t *testing.T) {
	selector := map[string]string{"app": "x"}
	got := Defaulted(nil, 2, selector)
	if len(got) < 2 {
		t.Fatalf("expected >=2 TSCs, got %d", len(got))
	}
	for k, v := range selector {
		if got[0].LabelSelector.MatchLabels[k] != v {
			t.Errorf("first TSC selector mismatch on %s", k)
		}
		if got[1].LabelSelector.MatchLabels[k] != v {
			t.Errorf("second TSC selector mismatch on %s", k)
		}
	}
}

func TestDefaultTopologyKeys_order(t *testing.T) {
	keys := DefaultTopologyKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 default keys, got %d", len(keys))
	}
	if keys[0] != TopologyKeyZone {
		t.Errorf("first key: %q want %q", keys[0], TopologyKeyZone)
	}
	if keys[1] != TopologyKeyHostname {
		t.Errorf("second key: %q want %q", keys[1], TopologyKeyHostname)
	}
}

func TestDefaulted_combined_options(t *testing.T) {
	got := Defaulted(nil, 1, map[string]string{"app": "myapp"},
		WithMinReplicas(1),
		WithTopologyKeys("rack", TopologyKeyZone, TopologyKeyHostname),
		WithMaxSkew(2),
		WithWhenUnsatisfiable(corev1.DoNotSchedule),
	)
	if len(got) != 3 {
		t.Fatalf("expected 3 TSCs, got %d", len(got))
	}
	for i, c := range got {
		if c.MaxSkew != 2 {
			t.Errorf("TSC[%d] MaxSkew: %d", i, c.MaxSkew)
		}
		if c.WhenUnsatisfiable != corev1.DoNotSchedule {
			t.Errorf("TSC[%d] WhenUnsatisfiable: %q", i, c.WhenUnsatisfiable)
		}
	}
	if got[0].TopologyKey != "rack" {
		t.Errorf("first key: %q want rack", got[0].TopologyKey)
	}
}
