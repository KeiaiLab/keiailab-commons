// SPDX-License-Identifier: Apache-2.0

// Package pvc 는 downstream operator 의 공통 PVC
// 처리 헬퍼를 제공한다. 현재 단일 기능: online expansion.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator 의 *pvc_resize.go* 중복
// 코드 (~360 LOC) 를 추출하여 신규 도입되었다. downstream consumer 동시 적용 회귀 통과
// 후 Stable Tier 격상 예정 (자세한 격상 조건은 docs/ROADMAP.md §API Stability Tier).
//
// # 의존성 정책
//
// commons 의 *순수 데이터 변환* 패키지 (finalizer, status, security 등) 와
// 달리 본 패키지는 K8s API 호출을 직접 수행해야 한다. 따라서 단 한 곳
// (`pkg/pvc`) 에서만 controller-runtime/pkg/client 를 의존한다.
//
//   - List + Patch + Get 의 3-step 워크플로우 — 호출자가 imperative 한
//     라이프사이클을 직접 관리하는 monitoring / security 와 다름.
//   - controller-runtime client.Client 는 downstream operator 가 이미 보유한 표준
//     의존이므로 신규 의존 부담은 0.
//
// # 설계 원칙
//
//   - StatefulSet.spec.volumeClaimTemplates 는 immutable → *기존 PVC 직접 patch*.
//   - StorageClass.AllowVolumeExpansion 사전 검증 (false 시 noop).
//   - CSI online resize 미지원 시 PVC.status 가 FileSystemResizePending 으로 남고
//     다음 pod restart 시 완료 — 본 패키지는 *patch 만* 하고 폴링하지 않음.
//   - 부분 실패 best-effort: 1 PVC patch 실패가 다른 PVC 차단하지 않음.
//   - 로그는 controller-runtime/pkg/log.FromContext — 호출자가 ctx 에 주입한
//     logger 가 자동 사용된다.
package pvc
