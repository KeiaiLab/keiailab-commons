// SPDX-License-Identifier: MIT

package version

import (
	"testing"
)

func TestMustList_Empty_Panics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty list")
		}
	}()
	MustList()
}

func TestMustList_EmptyString_Panics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty string entry")
		}
	}()
	MustList("8.0", "", "8.2")
}

func TestMustList_Valid(t *testing.T) {
	t.Parallel()
	l := MustList("8.0", "8.2", "8.3")
	if !l.IsSupported("8.0") {
		t.Error("8.0 should be supported")
	}
	if !l.IsSupported("8.3") {
		t.Error("8.3 should be supported")
	}
	if l.IsSupported("9.0") {
		t.Error("9.0 should not be supported")
	}
	if l.IsSupported("") {
		t.Error("empty string should not be supported")
	}
}

func TestList_Strings_DefensiveCopy(t *testing.T) {
	t.Parallel()
	l := MustList("a", "b")
	out := l.Strings()
	out[0] = "MUTATED"
	if l.IsSupported("MUTATED") {
		t.Error("List.Strings must return defensive copy — internal state mutated")
	}
	if !l.IsSupported("a") {
		t.Error("original 'a' lost after caller mutation")
	}
}

func TestList_Default(t *testing.T) {
	t.Parallel()
	l := MustList("8.0", "8.2")
	if got := l.Default(); got != "8.0" {
		t.Errorf("Default = %q, want 8.0", got)
	}
}

func TestList_MutationIsolation(t *testing.T) {
	t.Parallel()
	src := []string{"x", "y"}
	l := MustList(src...)
	src[0] = "MUTATED"
	if l.IsSupported("MUTATED") {
		t.Error("MustList must defensive-copy varargs")
	}
}
