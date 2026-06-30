// SPDX-License-Identifier: MIT

// Package hpa 는 downstream operator 공통의 HorizontalPodAutoscaler 빌더와
// 표준 metric 헬퍼를 제공한다. apps/v1 ScaleTargetRef + MinReplicas clamp +
// Resource-Utilization metric 조립 보일러플레이트를 단일 진실원으로 고정한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 valkey-operator (`resources/hpa.go` BuildHorizontalPodAutoscaler)
// 와 mongodb-operator (`resources/builder.go` buildHPAForTarget + buildHPAMetrics)
// 의 평행 재구현을 추출하여 신규 도입되었다. downstream 동시 적용 회귀 통과 후
// Stable 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
// valkey 와 mongo 의 HPA 조립은 다음이 달랐다:
//
//   - MinReplicas clamp floor: valkey=2 / mongo=1. → Params.MinFloor 로 파라미터화.
//   - metric 모델: valkey 는 (CPU%, Mem%) 고정 2종, mongo 는 cpu/memory/custom
//     타입 리스트. → commons 는 ScaleTargetRef + clamp 의 *envelope* 와 표준
//     Resource-Utilization 헬퍼(CPUUtilization/MemoryUtilization)만 제공하고,
//     도메인별 metric 조립(특히 mongo 의 custom Pods metric)은 호출자에 잔류한다.
//   - ptr 관용구: valkey 원본은 `new(x)` 를 썼으나 commons 는 `ptr.To` 로 통일
//     (golangci modernize.newexpr 비활성 정합).
//
// # 의존성 정책
//
// k8s.io/api/autoscaling/v2 + core/v1 + apimachinery + k8s.io/utils/ptr 만
// 의존하는 *순수 데이터 변환* 패키지. K8s client 의존 없음 (생성된 객체의
// persist 는 pkg/apply.HPA 책임 — 레이어 분리).
//
// # 설계 원칙
//
//   - CR 타입 비의존: Metrics 는 호출자가 조립한 []autoscalingv2.MetricSpec 로 전달.
//   - clamp: MinReplicas = max(MinReplicas, MinFloor), MaxReplicas = max(MaxReplicas, MinReplicas).
//   - opt-in(enabled) 판정은 호출자 책임 — 비활성 시 호출자가 nil 처리.
package hpa
