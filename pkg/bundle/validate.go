// SPDX-License-Identifier: MIT

package bundle

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Sentinel errors returned by ValidateDir.
var (
	ErrMissingManifests  = errors.New("bundle: manifests/ directory not found")
	ErrMissingMetadata   = errors.New("bundle: metadata/ directory not found")
	ErrMissingAnnotations = errors.New("bundle: metadata/annotations.yaml not found")
	ErrNoManifestFiles   = errors.New("bundle: manifests/ contains no .yaml or .yml files")
)

// ValidateDir checks that path is a valid OLM operator bundle directory.
//
// A valid bundle directory must contain:
//   - manifests/ subdirectory with at least one .yaml or .yml file
//   - metadata/ subdirectory
//   - metadata/annotations.yaml file with the required annotation keys
//
// ValidateDir returns the first error encountered. All errors are one of the
// typed sentinel errors (ErrMissing*, ErrNoManifestFiles) or a wrapped variant
// with additional context.
func ValidateDir(path string) error {
	// Check manifests/ directory.
	manifestsDir := filepath.Join(path, "manifests")
	info, err := os.Stat(manifestsDir)
	if err != nil || !info.IsDir() {
		return ErrMissingManifests
	}

	// Check that manifests/ contains at least one YAML file.
	hasYAML, err := dirContainsYAML(manifestsDir)
	if err != nil {
		return fmt.Errorf("bundle: reading manifests/: %w", err)
	}
	if !hasYAML {
		return ErrNoManifestFiles
	}

	// Check metadata/ directory.
	metadataDir := filepath.Join(path, "metadata")
	info, err = os.Stat(metadataDir)
	if err != nil || !info.IsDir() {
		return ErrMissingMetadata
	}

	// Check metadata/annotations.yaml existence.
	annotationsPath := filepath.Join(metadataDir, "annotations.yaml")
	if _, err := os.Stat(annotationsPath); err != nil {
		return ErrMissingAnnotations
	}

	// Validate that annotations.yaml contains the required keys.
	if err := validateAnnotationsFile(annotationsPath); err != nil {
		return err
	}

	return nil
}

// dirContainsYAML reports whether dir contains at least one file with a .yaml
// or .yml extension.
func dirContainsYAML(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".yaml" || ext == ".yml" {
			return true, nil
		}
	}
	return false, nil
}

// requiredAnnotationKeys lists the keys that must be present in annotations.yaml.
var requiredAnnotationKeys = []string{
	MediaTypeKey,
	ManifestsKey,
	MetadataKey,
	PackageKey,
	ChannelsKey,
}

// validateAnnotationsFile reads annotations.yaml and checks that each required
// annotation key appears as a substring. This is intentionally simple — it does
// NOT use a YAML parser to keep the package free of non-stdlib dependencies.
// For build-time validation this approach is sufficient: the keys are unique
// enough that a substring check is reliable.
func validateAnnotationsFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("bundle: reading annotations.yaml: %w", err)
	}
	content := string(data)
	for _, key := range requiredAnnotationKeys {
		if !strings.Contains(content, key) {
			return fmt.Errorf("bundle: annotations.yaml missing required key %q", key)
		}
	}
	return nil
}
