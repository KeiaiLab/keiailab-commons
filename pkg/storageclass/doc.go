// SPDX-License-Identifier: Apache-2.0

// Stability: Stable.
//
// 본 패키지는 K8s PersistentVolumeClaim.Spec.StorageClassName 의 DNS-1123
// subdomain validation + nil-pointer normalization 을 통일합니다.
// downstream operator (postgres / mongodb / valkey) 의 storageClassPtr()
// pattern (empty → nil, non-empty → &string) 의 공통 source.
//
// API surface 가 trivial (regex + nil check) 하므로 Stable Tier 즉시.
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package storageclass
