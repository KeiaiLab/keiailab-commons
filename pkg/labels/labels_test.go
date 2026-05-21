package labels

import (
	"reflect"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	t.Parallel()
	s := New("example-app", "primary", "controller-manager", "1.4.6", "downstream-component")
	if s.Name != "example-app" || s.Instance != "primary" || s.Component != "controller-manager" {
		t.Errorf("unexpected: %+v", s)
	}
	if s.Version != "1.4.6" || s.ManagedBy != "downstream-component" {
		t.Errorf("unexpected version/managedBy: %+v", s)
	}
	if s.PartOf != "downstream-component" {
		t.Errorf("PartOf default = ManagedBy expected, got %q", s.PartOf)
	}
}

func TestSet_All_AllFields(t *testing.T) {
	t.Parallel()
	s := New("example-app", "primary", "controller-manager", "1.4.6", "downstream-component")
	want := map[string]string{
		"app.kubernetes.io/name":       "example-app",
		"app.kubernetes.io/instance":   "primary",
		"app.kubernetes.io/component":  "controller-manager",
		"app.kubernetes.io/version":    "1.4.6",
		"app.kubernetes.io/managed-by": "downstream-component",
		"app.kubernetes.io/part-of":    "downstream-component",
	}
	got := s.All()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("All() mismatch:\n got: %v\nwant: %v", got, want)
	}
}

func TestSet_All_OptionalEmpty(t *testing.T) {
	t.Parallel()
	s := Set{
		Name:      "valkey",
		Instance:  "cache",
		ManagedBy: "downstream operator",
		// Component/Version/PartOf 미지정.
	}
	got := s.All()
	if got["app.kubernetes.io/name"] != "valkey" {
		t.Error("name missing")
	}
	for _, key := range []string{
		"app.kubernetes.io/component",
		"app.kubernetes.io/version",
		"app.kubernetes.io/part-of",
	} {
		if _, ok := got[key]; ok {
			t.Errorf("optional label %q should be absent when empty", key)
		}
	}
}

func TestSet_Selector_VersionExcluded(t *testing.T) {
	t.Parallel()
	s := New("example-app", "primary", "controller-manager", "1.4.6", "downstream-component")
	got := s.Selector()
	want := map[string]string{
		"app.kubernetes.io/name":      "example-app",
		"app.kubernetes.io/instance":  "primary",
		"app.kubernetes.io/component": "controller-manager",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Selector() mismatch:\n got: %v\nwant: %v", got, want)
	}
	if _, ok := got["app.kubernetes.io/version"]; ok {
		t.Error("version must NOT be in Selector (k8s immutable selector field — rolling update 차단)")
	}
}

func TestSet_Selector_ComponentEmpty(t *testing.T) {
	t.Parallel()
	s := Set{Name: "n", Instance: "i", ManagedBy: "m"}
	got := s.Selector()
	if _, ok := got["app.kubernetes.io/component"]; ok {
		t.Error("component should be absent when empty")
	}
}
