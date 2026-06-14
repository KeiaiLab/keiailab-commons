// SPDX-License-Identifier: MIT

// Stability: Beta.
//
// 본 패키지는 Kubernetes Event 생성 helper 를 통일합니다.
// 표준 Reason constants (Created / Updated / Deleted / Reconciled / Failed /
// Provisioning / Ready / Degraded) + minimal Recorder interface — 신 events API
// (client-go events.EventRecorder, RFC-0023 Phase 2) 와 structurally compatible.
// caller 가 controller-runtime Manager 에서 받은 신식 recorder 를 adapter 없이
// 그대로 전달한다 (commons 는 client-go tools 패키지 직접 의존 회피).
//
// 진본 동기: downstream operator 3종(mongodb/postgres/valkey) 가 RFC-0023 Phase 2
// 로 신 events.EventRecorder 마이그레이션 완료(2026-05-11). 본 패키지도 동일 신 API
// 형으로 정합(레거시 record.EventRecorder 형 폐기) — pkg/reconcile.EventRecorder 와
// 단일 interface 로 통합되어 commons 의 canonical Event interface + Reason SSOT.
// downstream operator 라이브 채택 후 Stable Tier 격상.
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package events
