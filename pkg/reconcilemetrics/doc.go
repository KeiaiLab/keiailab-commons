// SPDX-License-Identifier: MIT

// Package reconcilemetrics 는 downstream operator 의 공통 reconcile SLO
// Prometheus metrics trio (reconcile_total / reconcile_duration_seconds /
// reconcile_errors_total) 와 CR 삭제 시 cardinality 정리 (DeleteFor) 를
// 제공한다.
//
// 3 operator (mongodb / postgres / valkey) 가 subsystem 상수만 바꿔 복붙하던
// near-identical 정의 (~120 LOC) 를 흡수 — Help 문자열·bucket 경계·라벨
// 구조를 원본과 byte-동일하게 보존해 기존 시계열 이름
// (<subsystem>_reconcile_total 등) 이 그대로 유지된다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator 의 *internal/controller/metrics.go* 중복
// reconcile trio 를 추출하여 신규 도입되었다. downstream consumer 동시 적용
// 회귀 통과 후 Stable Tier 격상 예정 (자세한 격상 조건은
// docs/ROADMAP.md §API Stability Tier).
//
// # 의존성 정책
//
// github.com/prometheus/client_golang 을 직접 의존한다 — controller-runtime
// 경유 transitive 로 모든 downstream operator 가 이미 보유한 표준 의존이므로
// 신규 의존 부담은 0.
//
//   - 등록은 MustRegister(reg prometheus.Registerer) 로 호출자가 주입 —
//     controller-runtime metrics.Registry 글로벌 hard-coupling 을 피하고
//     테스트 격리 (test 별 NewRegistry) 를 보장한다.
//   - 도메인 메트릭 (mongodb query/replication, postgres WAL lag,
//     valkey cluster state 등) 은 본 패키지 scope 외 — 각 repo 잔류.
//
// # 설계 원칙
//
//   - 시계열 호환 절대 제약: New(subsystem) 의 metric 이름/Help/bucket/라벨은
//     원본 3 operator 정의와 byte-동일. 기존 Grafana dashboard +
//     PrometheusRule 이 변경 없이 동작한다.
//   - 라벨 cardinality 제어: namespace/name (+result/+component) 만 —
//     shard/pod 레벨 라벨은 의도적으로 제외 (대규모 cluster 시 폭발 방지).
//   - DeleteFor 는 trio *만* 삭제 — 도메인 메트릭 정리는 각 repo 의 기존
//     DeleteMetricsFor wrapper 가 DeleteFor 호출 후 이어서 수행한다.
//   - 원본 호출 시점 보존: Total 은 Reconcile *진입* 시 (IncTotal),
//     Latency 는 *종료* defer 시 (ObserveReconcile) — 단일 합성 메서드로
//     시점을 합치지 않는다 (in-flight 카운트 의미 변화 방지).
//   - struct 필드 (Total/Latency/Errors) 는 export — 각 repo 가
//     `var MetricReconcileTotal = reconMetrics.Total` alias 선언으로
//     controller 콜사이트 0 변경 마이그레이션이 가능하다.
package reconcilemetrics
