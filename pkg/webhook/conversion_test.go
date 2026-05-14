package webhook

import (
	"strings"
	"testing"
)

func TestConversionRegistry(t *testing.T) {
	r := NewConversionRegistry()
	r.Register("v1alpha1", "v1alpha2", func(src any) (any, error) {
		s, ok := src.(string)
		if !ok {
			return nil, nil
		}
		return strings.ToUpper(s), nil
	})

	t.Run("registered conversion", func(t *testing.T) {
		dst, err := r.Convert("v1alpha1", "v1alpha2", "hello")
		if err != nil {
			t.Errorf("Convert err = %v", err)
		}
		if dst != "HELLO" {
			t.Errorf("dst = %v, want HELLO", dst)
		}
	})

	t.Run("unregistered pair", func(t *testing.T) {
		_, err := r.Convert("v1alpha2", "v1beta1", "hello")
		if err == nil {
			t.Errorf("expected error for unregistered pair")
		}
	})

	t.Run("HasPair", func(t *testing.T) {
		if !r.HasPair("v1alpha1", "v1alpha2") {
			t.Errorf("HasPair returns false for registered")
		}
		if r.HasPair("v1beta1", "v1") {
			t.Errorf("HasPair returns true for unregistered")
		}
	})
}
