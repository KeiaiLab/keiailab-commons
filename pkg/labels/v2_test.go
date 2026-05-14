package labels

import "testing"

func TestAllV2(t *testing.T) {
	s := Set{
		Name:      "test",
		Instance:  "test-1",
		Version:   "1.0.0",
		Component: "primary",
		PartOf:    "test-app",
	}

	t.Run("empty v2 (only v1)", func(t *testing.T) {
		got := s.AllV2(V2{})
		if _, ok := got[LabelCreatedBy]; ok {
			t.Errorf("expected no %q with empty V2.CreatedBy", LabelCreatedBy)
		}
	})

	t.Run("full v2", func(t *testing.T) {
		got := s.AllV2(V2{CreatedBy: "ctrl", Tier: "cache", Owner: "team"})
		if got[LabelCreatedBy] != "ctrl" {
			t.Errorf("CreatedBy = %q, want %q", got[LabelCreatedBy], "ctrl")
		}
		if got[LabelTier] != "cache" {
			t.Errorf("Tier = %q, want %q", got[LabelTier], "cache")
		}
		if got[LabelOwner] != "team" {
			t.Errorf("Owner = %q, want %q", got[LabelOwner], "team")
		}
	})

	t.Run("partial v2", func(t *testing.T) {
		got := s.AllV2(V2{CreatedBy: "ctrl"})
		if got[LabelCreatedBy] != "ctrl" {
			t.Errorf("CreatedBy missing")
		}
		if _, ok := got[LabelTier]; ok {
			t.Errorf("Tier should be absent")
		}
	})
}
