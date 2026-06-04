// SPDX-License-Identifier: MIT

package storageclass_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/keiailab/operator-commons/pkg/storageclass"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty", "", false},
		{"single lowercase", "a", true},
		{"alphanumeric", "fast-ssd-2", true},
		{"with dot segment", "ceph.rbd", true},
		{"multi dot", "rook.ceph.block", true},
		{"uppercase invalid", "Fast-SSD", false},
		{"leading hyphen invalid", "-fast", false},
		{"trailing hyphen invalid", "fast-", false},
		{"leading dot invalid", ".rbd", false},
		{"underscore invalid", "fast_ssd", false},
		{"length 253 boundary", strings.Repeat("a", 253), true},
		{"length 254 exceeds max", strings.Repeat("a", 254), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := storageclass.IsValid(tc.input)
			if got != tc.want {
				t.Errorf("IsValid(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidate_emptyAllowed(t *testing.T) {
	t.Parallel()

	if err := storageclass.Validate(""); err != nil {
		t.Errorf("Validate(\"\") = %v, want nil (empty treated as cluster default)", err)
	}
}

func TestValidate_validReturnsNil(t *testing.T) {
	t.Parallel()

	if err := storageclass.Validate("fast-ssd"); err != nil {
		t.Errorf("Validate(\"fast-ssd\") = %v, want nil", err)
	}
}

func TestValidate_invalidReturnsErr(t *testing.T) {
	t.Parallel()

	err := storageclass.Validate("Fast-SSD")
	if err == nil {
		t.Fatal("Validate(\"Fast-SSD\") = nil, want ErrInvalidStorageClassName")
	}
	if !errors.Is(err, storageclass.ErrInvalidStorageClassName) {
		t.Errorf("error sentinel mismatch: got=%v", err)
	}
}

func TestNormalize_empty(t *testing.T) {
	t.Parallel()

	if got := storageclass.Normalize(""); got != nil {
		t.Errorf("Normalize(\"\") = %v, want nil", got)
	}
}

func TestNormalize_whitespaceOnly(t *testing.T) {
	t.Parallel()

	if got := storageclass.Normalize("   \t\n  "); got != nil {
		t.Errorf("Normalize(whitespace-only) = %v, want nil", got)
	}
}

func TestNormalize_validReturnsPtr(t *testing.T) {
	t.Parallel()

	got := storageclass.Normalize("ceph-rbd")
	if got == nil {
		t.Fatal("Normalize(\"ceph-rbd\") = nil, want pointer")
	}
	if *got != "ceph-rbd" {
		t.Errorf("Normalize: *got = %q, want %q", *got, "ceph-rbd")
	}
}

func TestNormalize_trimsWhitespace(t *testing.T) {
	t.Parallel()

	got := storageclass.Normalize("  fast-ssd  ")
	if got == nil {
		t.Fatal("Normalize trimmed result nil")
	}
	if *got != "fast-ssd" {
		t.Errorf("Normalize: *got = %q, want trimmed %q", *got, "fast-ssd")
	}
}

func TestMustNormalize_emptyOK(t *testing.T) {
	t.Parallel()

	if got := storageclass.MustNormalize(""); got != nil {
		t.Errorf("MustNormalize(\"\") = %v, want nil", got)
	}
}

func TestMustNormalize_validOK(t *testing.T) {
	t.Parallel()

	got := storageclass.MustNormalize("fast-ssd")
	if got == nil || *got != "fast-ssd" {
		t.Errorf("MustNormalize(\"fast-ssd\"): got=%v, want=\"fast-ssd\"", got)
	}
}

func TestMustNormalize_invalidPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustNormalize on invalid input must panic")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value type: want=string got=%T", r)
		}
		if !strings.Contains(msg, "storageclass.MustNormalize") {
			t.Errorf("panic message missing prefix: %q", msg)
		}
	}()

	_ = storageclass.MustNormalize("Invalid_Name")
}

func TestStorageClassPtrPattern(t *testing.T) {
	t.Parallel()

	// Replaces: downstream operator storageClassPtr() helper pattern.
	// Pattern: empty → nil (cluster default), non-empty → &string.

	cases := []struct {
		input   string
		wantPtr bool
		wantVal string
	}{
		{"", false, ""},
		{"fast-ssd", true, "fast-ssd"},
		{"rook-ceph-block", true, "rook-ceph-block"},
	}

	for _, tc := range cases {
		got := storageclass.Normalize(tc.input)
		if (got != nil) != tc.wantPtr {
			t.Errorf("input=%q: ptr presence mismatch (got=%v, want=%v)", tc.input, got != nil, tc.wantPtr)
			continue
		}
		if got != nil && *got != tc.wantVal {
			t.Errorf("input=%q: value mismatch (got=%q, want=%q)", tc.input, *got, tc.wantVal)
		}
	}
}
