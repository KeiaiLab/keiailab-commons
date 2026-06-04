// SPDX-License-Identifier: MIT

package bundle

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeValidBundle creates a temporary bundle directory with the minimal valid
// structure and returns its root path. The caller does not need to clean up —
// t.TempDir() handles it.
func makeValidBundle(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	manifestsDir := filepath.Join(root, "manifests")
	metadataDir := filepath.Join(root, "metadata")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write a minimal CSV manifest.
	csv := filepath.Join(manifestsDir, "csv.yaml")
	if err := os.WriteFile(csv, []byte("apiVersion: operators.coreos.com/v1alpha1\nkind: ClusterServiceVersion\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write a valid annotations.yaml with all required keys.
	a := NewAnnotations("test-operator", []string{"stable"}, "stable")
	var sb strings.Builder
	sb.WriteString("annotations:\n")
	for k, v := range a.Map() {
		sb.WriteString("  " + k + ": " + v + "\n")
	}
	annotationsPath := filepath.Join(metadataDir, "annotations.yaml")
	if err := os.WriteFile(annotationsPath, []byte(sb.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	return root
}

func TestValidateDir_ValidBundle(t *testing.T) {
	t.Parallel()

	root := makeValidBundle(t)
	if err := ValidateDir(root); err != nil {
		t.Errorf("ValidateDir on valid bundle returned error: %v", err)
	}
}

func TestValidateDir_MissingManifests(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	// Create metadata/ but not manifests/.
	metadataDir := filepath.Join(root, "metadata")
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrMissingManifests) {
		t.Errorf("expected ErrMissingManifests, got %v", err)
	}
}

func TestValidateDir_MissingMetadata(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	// Create manifests/ with a YAML file but no metadata/.
	manifestsDir := filepath.Join(root, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestsDir, "csv.yaml"), []byte("kind: CSV\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrMissingMetadata) {
		t.Errorf("expected ErrMissingMetadata, got %v", err)
	}
}

func TestValidateDir_MissingAnnotationsYAML(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	metadataDir := filepath.Join(root, "metadata")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestsDir, "csv.yml"), []byte("kind: CSV\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// metadata/ exists but annotations.yaml is absent.

	err := ValidateDir(root)
	if !errors.Is(err, ErrMissingAnnotations) {
		t.Errorf("expected ErrMissingAnnotations, got %v", err)
	}
}

func TestValidateDir_EmptyManifests(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	metadataDir := filepath.Join(root, "metadata")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// manifests/ exists but has no YAML files (only a .txt).
	if err := os.WriteFile(filepath.Join(manifestsDir, "readme.txt"), []byte("not yaml"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrNoManifestFiles) {
		t.Errorf("expected ErrNoManifestFiles, got %v", err)
	}
}

func TestValidateDir_ManifestsWithYMLExtension(t *testing.T) {
	t.Parallel()

	root := makeValidBundle(t)
	// The helper already creates csv.yaml. Add a .yml file to verify both
	// extensions are accepted.
	ymlPath := filepath.Join(root, "manifests", "extra.yml")
	if err := os.WriteFile(ymlPath, []byte("kind: ConfigMap\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ValidateDir(root); err != nil {
		t.Errorf("ValidateDir should accept .yml files: %v", err)
	}
}

func TestValidateDir_NonexistentPath(t *testing.T) {
	t.Parallel()

	err := ValidateDir("/nonexistent/path/that/does/not/exist")
	if !errors.Is(err, ErrMissingManifests) {
		t.Errorf("expected ErrMissingManifests for nonexistent path, got %v", err)
	}
}

func TestValidateDir_AnnotationsMissingRequiredKey(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	metadataDir := filepath.Join(root, "metadata")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestsDir, "csv.yaml"), []byte("kind: CSV\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write annotations.yaml missing the PackageKey.
	incomplete := "annotations:\n  " + MediaTypeKey + ": " + MediaTypeRegistryV1 + "\n"
	annotationsPath := filepath.Join(metadataDir, "annotations.yaml")
	if err := os.WriteFile(annotationsPath, []byte(incomplete), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if err == nil {
		t.Fatal("expected error for incomplete annotations.yaml")
	}
	if !strings.Contains(err.Error(), "missing required key") {
		t.Errorf("expected 'missing required key' in error, got %v", err)
	}
}

func TestValidateDir_ManifestsIsFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	// Create manifests as a regular file, not a directory.
	if err := os.WriteFile(filepath.Join(root, "manifests"), []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrMissingManifests) {
		t.Errorf("expected ErrMissingManifests when manifests/ is a file, got %v", err)
	}
}

func TestValidateDir_MetadataIsFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestsDir, "csv.yaml"), []byte("kind: CSV\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Create metadata as a regular file, not a directory.
	if err := os.WriteFile(filepath.Join(root, "metadata"), []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrMissingMetadata) {
		t.Errorf("expected ErrMissingMetadata when metadata/ is a file, got %v", err)
	}
}

func TestValidateDir_ManifestsOnlyNonYAMLFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Only non-YAML files in manifests/.
	if err := os.WriteFile(filepath.Join(manifestsDir, "notes.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(manifestsDir, "data.txt"), []byte("text"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := ValidateDir(root)
	if !errors.Is(err, ErrNoManifestFiles) {
		t.Errorf("expected ErrNoManifestFiles, got %v", err)
	}
}

func TestValidateDir_ManifestsSubdirsIgnored(t *testing.T) {
	t.Parallel()

	root := makeValidBundle(t)
	// Add a subdirectory whose name sorts before csv.yaml to ensure the
	// IsDir() skip branch in dirContainsYAML is exercised before any YAML
	// file is found. ("aaa-dir" < "csv.yaml" in ReadDir order.)
	subDir := filepath.Join(root, "manifests", "aaa-dir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := ValidateDir(root); err != nil {
		t.Errorf("subdirectories in manifests/ should be ignored: %v", err)
	}
}

func TestValidateDir_ManifestsUnreadable(t *testing.T) {
	t.Parallel()

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	root := t.TempDir()
	manifestsDir := filepath.Join(root, "manifests")
	if err := os.MkdirAll(manifestsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Write a YAML file first, then make dir unreadable.
	if err := os.WriteFile(filepath.Join(manifestsDir, "csv.yaml"), []byte("kind: CSV\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(manifestsDir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(manifestsDir, 0o755) })

	err := ValidateDir(root)
	if err == nil {
		t.Error("expected error for unreadable manifests/")
	}
	if strings.Contains(err.Error(), "reading manifests/") {
		// Good — the wrapped os.ReadDir error.
		return
	}
	// On some OSes, Stat on the dir might fail differently.
	t.Logf("got error (acceptable): %v", err)
}

func TestValidateDir_AnnotationsUnreadable(t *testing.T) {
	t.Parallel()

	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	root := makeValidBundle(t)
	annotationsPath := filepath.Join(root, "metadata", "annotations.yaml")
	if err := os.Chmod(annotationsPath, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(annotationsPath, 0o644) })

	err := ValidateDir(root)
	if err == nil {
		t.Error("expected error for unreadable annotations.yaml")
	}
	if strings.Contains(err.Error(), "reading annotations.yaml") {
		// Good — the wrapped os.ReadFile error.
		return
	}
	t.Logf("got error (acceptable): %v", err)
}
