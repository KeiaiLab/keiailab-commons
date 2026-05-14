package webhook

import (
	"fmt"
	"strings"
)

// ValidateStorageClassDNS1123 — storageClassName 이 RFC 1123 subdomain 정합인지
// 검사. CRD CEL 보다 명확한 에러 메시지 + 일관 caller path.
//
// 사용 예 (valkey-operator):
//
//	if err := webhook.ValidateStorageClassDNS1123(spec.Storage.StorageClassName); err != nil {
//	    return field.Invalid(field.NewPath("spec", "storage", "storageClassName"), spec.Storage.StorageClassName, err.Error())
//	}
//
// Refs: ROADMAP.md 'Validation webhook 공통 패턴 (RBD storageClass, topology spread, replicaCount lower bound)'
//       (P-B.10.2)
func ValidateStorageClassDNS1123(name string) error {
	if name == "" {
		return nil // empty = use default (not invalid)
	}
	if len(name) > 253 {
		return fmt.Errorf("storageClassName length %d exceeds 253 (RFC 1123 subdomain)", len(name))
	}
	for _, c := range name {
		if !isAlphanumeric(c) && c != '-' && c != '.' {
			return fmt.Errorf("storageClassName contains invalid character %q (RFC 1123 subdomain: [a-z0-9.-])", c)
		}
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("storageClassName must not start/end with '-'")
	}
	return nil
}

// ValidateReplicaCountLowerBound — replicaCount 가 최소값 이상인지 검사.
// 호출자가 *클러스터 / 토폴로지별 최소값* 을 결정.
//
// 사용 예:
//
//	if err := webhook.ValidateReplicaCountLowerBound(spec.Replicas, 3, "HA topology"); err != nil {
//	    return field.Invalid(...)
//	}
func ValidateReplicaCountLowerBound(actual, minimum int32, reason string) error {
	if actual < minimum {
		return fmt.Errorf("replicas (%d) below minimum %d (%s)", actual, minimum, reason)
	}
	return nil
}

// ValidateTopologySpreadConstraints — TopologySpreadConstraint 키 화이트리스트
// 검사. K8s 의 well-known topology key 만 허용.
func ValidateTopologySpreadKeyAllowed(key string) error {
	allowed := map[string]bool{
		"topology.kubernetes.io/zone":   true,
		"topology.kubernetes.io/region": true,
		"kubernetes.io/hostname":        true,
		"node.kubernetes.io/instance-type": true,
	}
	if !allowed[key] {
		return fmt.Errorf("topology key %q not in allowed list: zone / region / hostname / instance-type", key)
	}
	return nil
}

func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
