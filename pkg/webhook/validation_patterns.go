// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
	"strings"
)

// ValidateStorageClassDNS1123 checks RFC 1123 subdomain compliance.
func ValidateStorageClassDNS1123(name string) error {
	if name == "" {
		return nil
	}
	if len(name) > 253 {
		return fmt.Errorf("storageClassName length %d > 253", len(name))
	}
	for _, c := range name {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '-' && c != '.' {
			return fmt.Errorf("storageClassName invalid char %q", c)
		}
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("storageClassName must not start/end with '-'")
	}
	return nil
}

// ValidateReplicaCountLowerBound checks replicas >= minimum.
func ValidateReplicaCountLowerBound(actual, minimum int32, reason string) error {
	if actual < minimum {
		return fmt.Errorf("replicas (%d) below minimum %d (%s)", actual, minimum, reason)
	}
	return nil
}

// ValidateTopologySpreadKeyAllowed restricts to well-known topology keys.
func ValidateTopologySpreadKeyAllowed(key string) error {
	allowed := map[string]bool{
		"topology.kubernetes.io/zone":      true,
		"topology.kubernetes.io/region":    true,
		"kubernetes.io/hostname":           true,
		"node.kubernetes.io/instance-type": true,
	}
	if !allowed[key] {
		return fmt.Errorf("topology key %q not in allowed list", key)
	}
	return nil
}
