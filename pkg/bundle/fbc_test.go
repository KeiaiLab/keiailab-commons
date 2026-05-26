// SPDX-License-Identifier: Apache-2.0

package bundle

import (
	"encoding/json"
	"testing"
)

// --- Schema constants ---

func TestSchemaConstants(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		got  string
		want string
	}{
		{"SchemaPackage", SchemaPackage, "olm.package"},
		{"SchemaChannel", SchemaChannel, "olm.channel"},
		{"SchemaBundle", SchemaBundle, "olm.bundle"},
		{"SchemaDeprecation", SchemaDeprecation, "olm.deprecations"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

// --- Constructors ---

func TestNewPackage(t *testing.T) {
	t.Parallel()

	p := NewPackage("my-operator", "stable", "An example operator")

	if p.Schema != SchemaPackage {
		t.Errorf("Schema = %q, want %q", p.Schema, SchemaPackage)
	}
	if p.Name != "my-operator" {
		t.Errorf("Name = %q, want %q", p.Name, "my-operator")
	}
	if p.DefaultChannel != "stable" {
		t.Errorf("DefaultChannel = %q, want %q", p.DefaultChannel, "stable")
	}
	if p.Description != "An example operator" {
		t.Errorf("Description = %q, want %q", p.Description, "An example operator")
	}
	if p.Icon != nil {
		t.Errorf("Icon should be nil by default, got %v", p.Icon)
	}
}

func TestNewChannel(t *testing.T) {
	t.Parallel()

	entries := []ChannelEntry{
		{Name: "my-operator.v0.1.0"},
		{Name: "my-operator.v0.2.0", Replaces: "my-operator.v0.1.0"},
	}
	ch := NewChannel("my-operator", "stable", entries...)

	if ch.Schema != SchemaChannel {
		t.Errorf("Schema = %q, want %q", ch.Schema, SchemaChannel)
	}
	if ch.Package != "my-operator" {
		t.Errorf("Package = %q, want %q", ch.Package, "my-operator")
	}
	if ch.Name != "stable" {
		t.Errorf("Name = %q, want %q", ch.Name, "stable")
	}
	if len(ch.Entries) != 2 {
		t.Fatalf("Entries len = %d, want 2", len(ch.Entries))
	}
	if ch.Entries[1].Replaces != "my-operator.v0.1.0" {
		t.Errorf("Entries[1].Replaces = %q, want %q", ch.Entries[1].Replaces, "my-operator.v0.1.0")
	}
}

func TestNewChannel_NoEntries(t *testing.T) {
	t.Parallel()

	ch := NewChannel("pkg", "alpha")

	if ch.Entries != nil {
		t.Errorf("Entries should be nil when no entries provided, got %v", ch.Entries)
	}
}

func TestNewBundle(t *testing.T) {
	t.Parallel()

	b := NewBundle("my-operator", "my-operator.v0.1.0", "quay.io/example/my-operator-bundle:v0.1.0")

	if b.Schema != SchemaBundle {
		t.Errorf("Schema = %q, want %q", b.Schema, SchemaBundle)
	}
	if b.Package != "my-operator" {
		t.Errorf("Package = %q, want %q", b.Package, "my-operator")
	}
	if b.Name != "my-operator.v0.1.0" {
		t.Errorf("Name = %q, want %q", b.Name, "my-operator.v0.1.0")
	}
	if b.Image != "quay.io/example/my-operator-bundle:v0.1.0" {
		t.Errorf("Image = %q, want %q", b.Image, "quay.io/example/my-operator-bundle:v0.1.0")
	}
	if b.Properties != nil {
		t.Errorf("Properties should be nil by default, got %v", b.Properties)
	}
	if b.RelatedImages != nil {
		t.Errorf("RelatedImages should be nil by default, got %v", b.RelatedImages)
	}
}

// --- JSON Marshaling ---

func TestPackage_MarshalJSON(t *testing.T) {
	t.Parallel()

	p := NewPackage("my-operator", "stable", "desc")
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded["schema"] != SchemaPackage {
		t.Errorf("schema = %v, want %q", decoded["schema"], SchemaPackage)
	}
	if decoded["name"] != "my-operator" {
		t.Errorf("name = %v, want %q", decoded["name"], "my-operator")
	}
	if decoded["defaultChannel"] != "stable" {
		t.Errorf("defaultChannel = %v, want %q", decoded["defaultChannel"], "stable")
	}
}

func TestPackage_MarshalJSON_OmitsEmptyFields(t *testing.T) {
	t.Parallel()

	p := NewPackage("op", "ch", "")
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if _, ok := decoded["description"]; ok {
		t.Error("description should be omitted when empty")
	}
	if _, ok := decoded["icon"]; ok {
		t.Error("icon should be omitted when nil")
	}
}

func TestChannel_MarshalJSON(t *testing.T) {
	t.Parallel()

	ch := NewChannel("pkg", "stable",
		ChannelEntry{Name: "pkg.v1.0.0"},
		ChannelEntry{Name: "pkg.v1.1.0", Replaces: "pkg.v1.0.0", SkipRange: ">=1.0.0 <1.1.0"},
	)
	data, err := json.Marshal(ch)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded["schema"] != SchemaChannel {
		t.Errorf("schema = %v, want %q", decoded["schema"], SchemaChannel)
	}
	entries, ok := decoded["entries"].([]interface{})
	if !ok || len(entries) != 2 {
		t.Fatalf("entries len = %v, want 2", decoded["entries"])
	}
}

func TestBundle_MarshalJSON(t *testing.T) {
	t.Parallel()

	b := NewBundle("pkg", "pkg.v1.0.0", "quay.io/example/pkg:v1.0.0")
	b.Properties = []Property{
		{Type: "olm.gvk", Value: json.RawMessage(`{"group":"example.com","version":"v1","kind":"Foo"}`)},
	}
	b.RelatedImages = []RelatedImage{
		{Name: "manager", Image: "quay.io/example/manager:v1.0.0"},
	}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded["schema"] != SchemaBundle {
		t.Errorf("schema = %v, want %q", decoded["schema"], SchemaBundle)
	}
	if decoded["image"] != "quay.io/example/pkg:v1.0.0" {
		t.Errorf("image = %v, want %q", decoded["image"], "quay.io/example/pkg:v1.0.0")
	}

	props, ok := decoded["properties"].([]interface{})
	if !ok || len(props) != 1 {
		t.Fatalf("properties len = %v, want 1", decoded["properties"])
	}
	related, ok := decoded["relatedImages"].([]interface{})
	if !ok || len(related) != 1 {
		t.Fatalf("relatedImages len = %v, want 1", decoded["relatedImages"])
	}
}

func TestBundle_MarshalJSON_OmitsEmptyCollections(t *testing.T) {
	t.Parallel()

	b := NewBundle("pkg", "pkg.v1.0.0", "img")
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if _, ok := decoded["properties"]; ok {
		t.Error("properties should be omitted when nil")
	}
	if _, ok := decoded["relatedImages"]; ok {
		t.Error("relatedImages should be omitted when nil")
	}
}

func TestChannelEntry_MarshalJSON_OmitsEmpty(t *testing.T) {
	t.Parallel()

	e := ChannelEntry{Name: "pkg.v1.0.0"}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if _, ok := decoded["replaces"]; ok {
		t.Error("replaces should be omitted when empty")
	}
	if _, ok := decoded["skips"]; ok {
		t.Error("skips should be omitted when nil")
	}
	if _, ok := decoded["skipRange"]; ok {
		t.Error("skipRange should be omitted when empty")
	}
}
