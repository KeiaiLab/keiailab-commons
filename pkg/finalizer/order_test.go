package finalizer

import (
	"reflect"
	"testing"
)

func TestEnsureOrder(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		desired []string
		want    []string
	}{
		{
			name:    "empty input",
			input:   nil,
			desired: []string{"a"},
			want:    nil,
		},
		{
			name:    "empty desired",
			input:   []string{"a", "b"},
			desired: nil,
			want:    []string{"a", "b"},
		},
		{
			name:    "reorder to match desired",
			input:   []string{"b", "a"},
			desired: []string{"a", "b"},
			want:    []string{"a", "b"},
		},
		{
			name:    "unordered after ordered (stable)",
			input:   []string{"x", "b", "y", "a"},
			desired: []string{"a", "b"},
			want:    []string{"a", "b", "x", "y"},
		},
		{
			name:    "all unordered (stable)",
			input:   []string{"y", "x", "z"},
			desired: []string{"a", "b"},
			want:    []string{"y", "x", "z"},
		},
		{
			name:    "single element",
			input:   []string{"a"},
			desired: []string{"a"},
			want:    []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnsureOrder(tt.input, tt.desired)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnsureOrder(%v, %v) = %v, want %v", tt.input, tt.desired, got, tt.want)
			}
		})
	}
}
