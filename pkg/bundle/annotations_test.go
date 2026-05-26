// SPDX-License-Identifier: Apache-2.0

package bundle

import (
	"testing"
)

func TestNewAnnotations(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("my-operator", []string{"stable", "candidate"}, "stable")

	if a.PackageName != "my-operator" {
		t.Errorf("PackageName = %q, want %q", a.PackageName, "my-operator")
	}
	if len(a.Channels) != 2 || a.Channels[0] != "stable" || a.Channels[1] != "candidate" {
		t.Errorf("Channels = %v, want [stable candidate]", a.Channels)
	}
	if a.DefaultChannel != "stable" {
		t.Errorf("DefaultChannel = %q, want %q", a.DefaultChannel, "stable")
	}
}

func TestAnnotations_Map_AllKeys(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("my-operator", []string{"stable", "candidate"}, "stable")
	m := a.Map()

	requiredKeys := []string{
		MediaTypeKey,
		ManifestsKey,
		MetadataKey,
		PackageKey,
		ChannelsKey,
		DefaultChannelKey,
	}

	for _, key := range requiredKeys {
		if _, ok := m[key]; !ok {
			t.Errorf("Map() missing required key %q", key)
		}
	}

	if len(m) != 6 {
		t.Errorf("Map() returned %d keys, want 6", len(m))
	}
}

func TestAnnotations_Map_Values(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("my-operator", []string{"stable", "candidate"}, "stable")
	m := a.Map()

	cases := []struct {
		key  string
		want string
	}{
		{MediaTypeKey, MediaTypeRegistryV1},
		{ManifestsKey, "manifests/"},
		{MetadataKey, "metadata/"},
		{PackageKey, "my-operator"},
		{ChannelsKey, "stable,candidate"},
		{DefaultChannelKey, "stable"},
	}

	for _, tc := range cases {
		if got := m[tc.key]; got != tc.want {
			t.Errorf("Map()[%q] = %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestAnnotations_Map_SingleChannel(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("pkg", []string{"alpha"}, "alpha")
	m := a.Map()

	if got := m[ChannelsKey]; got != "alpha" {
		t.Errorf("single channel: ChannelsKey = %q, want %q", got, "alpha")
	}
}

func TestAnnotations_Map_EmptyChannels(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("pkg", nil, "")
	m := a.Map()

	if got := m[ChannelsKey]; got != "" {
		t.Errorf("empty channels: ChannelsKey = %q, want empty", got)
	}
	if got := m[DefaultChannelKey]; got != "" {
		t.Errorf("empty default: DefaultChannelKey = %q, want empty", got)
	}
	// Even with empty inputs, all 6 keys must be present.
	if len(m) != 6 {
		t.Errorf("Map() returned %d keys, want 6", len(m))
	}
}

func TestAnnotations_DockerLabels_MatchesMap(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("my-operator", []string{"stable", "fast"}, "stable")
	m := a.Map()
	dl := a.DockerLabels()

	if len(dl) != len(m) {
		t.Fatalf("DockerLabels() len=%d, Map() len=%d — must match", len(dl), len(m))
	}

	for k, v := range m {
		if got, ok := dl[k]; !ok {
			t.Errorf("DockerLabels() missing key %q", k)
		} else if got != v {
			t.Errorf("DockerLabels()[%q] = %q, Map()[%q] = %q", k, got, k, v)
		}
	}
}

func TestAnnotations_DockerLabels_MultipleChannels(t *testing.T) {
	t.Parallel()

	a := NewAnnotations("op", []string{"a", "b", "c"}, "a")
	dl := a.DockerLabels()

	if got := dl[ChannelsKey]; got != "a,b,c" {
		t.Errorf("DockerLabels channels = %q, want %q", got, "a,b,c")
	}
}
