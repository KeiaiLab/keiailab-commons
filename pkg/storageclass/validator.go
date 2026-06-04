// SPDX-License-Identifier: MIT

// Package storageclass provides DNS-1123 subdomain validation and
// nil-pointer normalization for K8s StorageClass name fields.
//
// Replaces the duplicated storageClassPtr() pattern across keiailab
// downstream operator, downstream operator, and downstream operator builders.
package storageclass

import (
	"errors"
	"regexp"
	"strings"
)

const (
	// MaxLength is the DNS-1123 subdomain maximum length per RFC 1123 §2.1.
	MaxLength = 253
)

// dns1123Subdomain matches the K8s DNS-1123 subdomain form used by StorageClass
// names: lowercase alphanumeric segments separated by hyphens or dots,
// no leading/trailing hyphens within a segment.
//
// Reference: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
var dns1123Subdomain = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

// ErrInvalidStorageClassName signals that the input is not a valid DNS-1123
// subdomain. Returned by Validate for non-empty invalid input.
var ErrInvalidStorageClassName = errors.New(
	"storage class name must conform to DNS-1123 subdomain " +
		"(lowercase alphanumeric, hyphens, dots; max 253 chars)",
)

// IsValid reports whether s is a valid DNS-1123 subdomain.
//
// Empty input is NOT valid per DNS-1123 rules — IsValid returns false for
// empty strings. Use Normalize when empty input signals "use cluster default
// storage class" intent.
func IsValid(s string) bool {
	if len(s) == 0 || len(s) > MaxLength {
		return false
	}
	return dns1123Subdomain.MatchString(s)
}

// Validate returns nil if s is empty (treated as "cluster default" intent) or
// is a valid DNS-1123 subdomain. Returns ErrInvalidStorageClassName otherwise.
//
// Use in webhook validation paths where empty input must be allowed but
// non-empty input must conform to K8s naming rules.
func Validate(s string) error {
	if s == "" {
		return nil
	}
	if !IsValid(s) {
		return ErrInvalidStorageClassName
	}
	return nil
}

// Normalize returns nil if s is empty after trimming surrounding whitespace
// (signals cluster default StorageClass), or a pointer to the trimmed string
// if non-empty.
//
// Use in PersistentVolumeClaim builder paths:
//
//	pvc.Spec.StorageClassName = storageclass.Normalize(spec.StorageClass)
//
// Format validation is not performed — call Validate separately if the input
// originates from an untrusted source.
func Normalize(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// MustNormalize behaves like Normalize for valid (or empty) input, and panics
// for non-empty invalid input. Intended for tests and bootstrap code where
// invalid input represents a programming error, not user input.
func MustNormalize(s string) *string {
	if err := Validate(s); err != nil {
		panic("storageclass.MustNormalize: " + err.Error() + ": " + s)
	}
	return Normalize(s)
}
