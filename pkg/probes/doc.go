// SPDX-License-Identifier: MIT

// Stability: Experimental.
//
// 본 패키지는 keiailab downstream operator 의
// corev1.Probe 구성 패턴 (총 9 sites, 50-55 LOC 중복) 을 단일 fluent builder
// 로 통일합니다. HTTP / HTTPS / TCP / Exec 4 핸들러 + kubelet default
// (PeriodSeconds=10 / TimeoutSeconds=1 / SuccessThreshold=1 / FailureThreshold=3)
// 를 제공하며, 2+ repo 라이브 적용 evidence 확보 시 Beta 격상 candidate.
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package probes
