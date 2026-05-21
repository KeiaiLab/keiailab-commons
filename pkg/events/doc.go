// SPDX-License-Identifier: Apache-2.0

// Stability: Beta.
//
// 본 패키지는 Kubernetes Event 생성 helper 를 통일합니다.
// 표준 Reason constants (Created / Updated / Deleted / Reconciled / Failed /
// Provisioning / Ready) + minimal Recorder interface (client-go record.EventRecorder
// 와 structurally compatible — caller 가 client-go 의존 추가 불필요).
//
// 진본 동기: downstream operator 결정 Phase 2 events.EventRecorder 적용
// (commit 1494ff6) + downstream operator 동일 진행 + downstream operator sister 적용
// candidate. downstream operator 라이브 적용 후 Stable Tier 격상.
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package events
