// SPDX-License-Identifier: Apache-2.0

package bundle

import "strings"

// Annotation key constants for the six required registry+v1 bundle annotations.
// These keys appear in both metadata/annotations.yaml and as Dockerfile LABELs.
const (
	MediaTypeKey      = "operators.operatorframework.io.bundle.mediatype.v1"
	ManifestsKey      = "operators.operatorframework.io.bundle.manifests.v1"
	MetadataKey       = "operators.operatorframework.io.bundle.metadata.v1"
	PackageKey        = "operators.operatorframework.io.bundle.package.v1"
	ChannelsKey       = "operators.operatorframework.io.bundle.channels.v1"
	DefaultChannelKey = "operators.operatorframework.io.bundle.channel.default.v1"

	// MediaTypeRegistryV1 is the standard media type value for registry+v1 bundles.
	MediaTypeRegistryV1 = "registry+v1"
)

// Annotations holds the metadata needed to generate a registry+v1 bundle's
// annotations.yaml and Dockerfile LABELs.
type Annotations struct {
	// PackageName is the operator package name (e.g. "my-operator").
	PackageName string

	// Channels lists the channel names this bundle belongs to (e.g. ["stable", "candidate"]).
	Channels []string

	// DefaultChannel is the default channel for the package (e.g. "stable").
	DefaultChannel string
}

// NewAnnotations creates an Annotations value with the required fields.
func NewAnnotations(packageName string, channels []string, defaultChannel string) *Annotations {
	return &Annotations{
		PackageName:    packageName,
		Channels:       channels,
		DefaultChannel: defaultChannel,
	}
}

// Map returns the full annotation map suitable for annotations.yaml.
// All six required keys are always present.
func (a *Annotations) Map() map[string]string {
	return map[string]string{
		MediaTypeKey:      MediaTypeRegistryV1,
		ManifestsKey:      "manifests/",
		MetadataKey:       "metadata/",
		PackageKey:        a.PackageName,
		ChannelsKey:       strings.Join(a.Channels, ","),
		DefaultChannelKey: a.DefaultChannel,
	}
}

// DockerLabels returns the annotation map formatted as Dockerfile LABEL keys.
// The returned map is identical to Map() — callers can iterate it to emit
// LABEL directives (e.g. LABEL "key"="value").
func (a *Annotations) DockerLabels() map[string]string {
	return a.Map()
}
