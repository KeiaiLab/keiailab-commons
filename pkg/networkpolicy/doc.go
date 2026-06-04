// SPDX-License-Identifier: MIT

// Stability: Stable.
//
// 본 패키지는 Beta → Stable 격상 완료 (PR #19 4-direction 검증 통과).
// deny-by-default NetworkPolicy builder + 4 functional option
// (WithSelfIngress / WithIngressFromPeers / WithDenyEgress / WithEgressToPeers).
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package networkpolicy
