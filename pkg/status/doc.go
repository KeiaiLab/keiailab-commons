// SPDX-License-Identifier: MIT

// Stability: Stable.
//
// 본 패키지는 v0.5+ 부터 BREAKING CHANGE 없이 안정 운영 중이며,
// v1.0.0 진입 시 자동으로 Stable Tier 됩니다.
// 4 표준 Condition Type + 6 Reason 카탈로그 (Ready/Progressing/Degraded/Available).
// 단 update.go 의 UpdateWithRetry 는 신규 표면 — Beta (downstream consumer 동시 적용
// 회귀 통과 후 패키지 Tier 합류).
//
// # 의존성 정책
//
// 파일 단위 격리 — condition 카탈로그 (conditions.go) 는 k8s.io/apimachinery 만
// 의존하는 순수 함수로 유지하고, status subresource 영속화 (update.go 의
// UpdateWithRetry) 만 controller-runtime/pkg/client + k8s.io/client-go (util/retry)
// 를 의존한다.
//
//   - Go 의존은 패키지 단위이므로 본 패키지 import 시 client 의존이 빌드 그래프에
//     포함되지만, controller-runtime client 는 downstream operator 가 이미 보유한
//     표준 의존이라 신규 부담은 0 (pkg/pvc §의존성 정책 sister).
//   - 순수 함수 코드 (conditions.go) 에 client 호출을 추가하지 않는 경계는 파일
//     분리로 유지한다 — client 의존 API 는 update.go 에만 둔다.
//
// 자세한 격상 조건 / API 안정성 매트릭스는 docs/ROADMAP.md §API Stability Tier 참조.
package status
