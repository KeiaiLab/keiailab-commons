// SPDX-License-Identifier: MIT

// Stability: Experimental.
//
// Package bundle provides build-time helpers for OLM v1 (Operator Lifecycle
// Manager v1) operator bundle metadata.
//
// It includes:
//   - Annotation constants and a builder for the six required registry+v1 bundle
//     annotation keys (used in annotations.yaml and Dockerfile LABELs).
//   - File-Based Catalog (FBC) schema types: Package, Channel, Bundle, and their
//     constructors — suitable for programmatic generation of FBC catalog fragments.
//   - Bundle directory validation (manifests/, metadata/, annotations.yaml presence
//     and minimal structural checks).
//
// This package intentionally has zero external dependencies beyond the Go standard
// library. It does NOT import any operator-framework or controller-runtime
// packages — it is a pure build-time / code-generation helper.
//
// See also: docs/ROADMAP.md §API Stability Tier for promotion criteria.
package bundle
