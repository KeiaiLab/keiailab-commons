// SPDX-License-Identifier: Apache-2.0

package bundle

import "encoding/json"

// FBC (File-Based Catalog) schema name constants.
const (
	SchemaPackage     = "olm.package"
	SchemaChannel     = "olm.channel"
	SchemaBundle      = "olm.bundle"
	SchemaDeprecation = "olm.deprecations"
)

// Package represents an olm.package entry in a File-Based Catalog.
type Package struct {
	Schema         string `json:"schema"`
	Name           string `json:"name"`
	DefaultChannel string `json:"defaultChannel"`
	Description    string `json:"description,omitempty"`
	Icon           *Icon  `json:"icon,omitempty"`
}

// Icon holds the optional icon data for a Package.
type Icon struct {
	Data      string `json:"base64data"`
	MediaType string `json:"mediatype"`
}

// NewPackage creates a Package with the olm.package schema preset.
func NewPackage(name, defaultChannel, description string) *Package {
	return &Package{
		Schema:         SchemaPackage,
		Name:           name,
		DefaultChannel: defaultChannel,
		Description:    description,
	}
}

// Channel represents an olm.channel entry in a File-Based Catalog.
type Channel struct {
	Schema  string         `json:"schema"`
	Package string         `json:"package"`
	Name    string         `json:"name"`
	Entries []ChannelEntry `json:"entries"`
}

// ChannelEntry describes a single entry (version) within a Channel.
type ChannelEntry struct {
	Name      string   `json:"name"`
	Replaces  string   `json:"replaces,omitempty"`
	Skips     []string `json:"skips,omitempty"`
	SkipRange string   `json:"skipRange,omitempty"`
}

// NewChannel creates a Channel with the olm.channel schema preset.
func NewChannel(packageName, channelName string, entries ...ChannelEntry) *Channel {
	return &Channel{
		Schema:  SchemaChannel,
		Package: packageName,
		Name:    channelName,
		Entries: entries,
	}
}

// Bundle represents an olm.bundle entry in a File-Based Catalog.
type Bundle struct {
	Schema        string         `json:"schema"`
	Package       string         `json:"package"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Properties    []Property     `json:"properties,omitempty"`
	RelatedImages []RelatedImage `json:"relatedImages,omitempty"`
}

// Property is a typed key-value property attached to a Bundle.
type Property struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// RelatedImage references a container image related to the Bundle.
type RelatedImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NewBundle creates a Bundle with the olm.bundle schema preset.
func NewBundle(packageName, name, image string) *Bundle {
	return &Bundle{
		Schema:  SchemaBundle,
		Package: packageName,
		Name:    name,
		Image:   image,
	}
}
