// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/keiailab/operator-commons/pkg/version"
)

func TestValidateAllowedVersion_EmptyValueSkipped(t *testing.T) {
	t.Parallel()
	list := version.MustList("8.0.9", "9.0.4")
	if err := ValidateAllowedVersion(field.NewPath("spec", "v"), "", list); err != nil {
		t.Errorf("empty value should be skipped, got %v", err)
	}
}

func TestValidateAllowedVersion_Allowed(t *testing.T) {
	t.Parallel()
	list := version.MustList("8.0.9", "9.0.4")
	for _, v := range []string{"8.0.9", "9.0.4"} {
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			if err := ValidateAllowedVersion(field.NewPath("spec", "v"), v, list); err != nil {
				t.Errorf("%q allowed but got error: %v", v, err)
			}
		})
	}
}

func TestValidateAllowedVersion_Rejected(t *testing.T) {
	t.Parallel()
	list := version.MustList("8.0.9", "9.0.4")
	err := ValidateAllowedVersion(field.NewPath("spec", "version"), "7.0.0", list)
	if err == nil {
		t.Fatal("expected NotSupported error, got nil")
	}
	if err.Type != field.ErrorTypeNotSupported {
		t.Errorf("Type = %v, want ErrorTypeNotSupported", err.Type)
	}
	if !strings.Contains(err.Detail, "8.0.9") || !strings.Contains(err.Detail, "9.0.4") {
		t.Errorf("Detail should list allowed values, got %q", err.Detail)
	}
	if err.Field != "spec.version" {
		t.Errorf("Field = %q, want spec.version", err.Field)
	}
}

func TestValidateWithPredicate_EmptyValueSkipped(t *testing.T) {
	t.Parallel()
	called := false
	predicate := func(string) bool { called = true; return false }
	err := ValidateWithPredicate(field.NewPath("v"), "", predicate, []string{"a", "b"})
	if err != nil {
		t.Errorf("empty → nil expected, got %v", err)
	}
	if called {
		t.Error("predicate should NOT be called for empty value (short-circuit)")
	}
}

func TestValidateWithPredicate_Allowed(t *testing.T) {
	t.Parallel()
	// semver-prefix 매칭 시뮬레이션 (major.minor pattern).
	majorMinor := func(v string) bool {
		// 단순화: "8.3" prefix 만 허용.
		return strings.HasPrefix(v, "8.3")
	}
	for _, v := range []string{"8.3", "8.3.0", "8.3.5"} {
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			err := ValidateWithPredicate(field.NewPath("v"), v, majorMinor, []string{"8.3"})
			if err != nil {
				t.Errorf("%q allowed but got error: %v", v, err)
			}
		})
	}
}

func TestValidateWithPredicate_Rejected(t *testing.T) {
	t.Parallel()
	majorMinor := func(v string) bool { return strings.HasPrefix(v, "8.3") }
	err := ValidateWithPredicate(field.NewPath("spec", "v"), "9.0.0", majorMinor, []string{"8.3"})
	if err == nil {
		t.Fatal("expected NotSupported error")
	}
	if err.Type != field.ErrorTypeNotSupported {
		t.Errorf("Type = %v, want ErrorTypeNotSupported", err.Type)
	}
	if !strings.Contains(err.Detail, "8.3") {
		t.Errorf("Detail should list allowed, got %q", err.Detail)
	}
}
