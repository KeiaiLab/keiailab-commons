// SPDX-License-Identifier: MIT

package reconcilemetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// result 라벨 표준 값 — 원본 3 operator 의 "result: success | error" 정합.
// DeleteFor 가 본 2 값만 순회하므로 호출자는 임의 result 문자열 대신
// ResultFor 또는 본 상수를 사용해야 시계열 누수가 없다.
const (
	// ResultSuccess — reconcile 이 에러 없이 종료.
	ResultSuccess = "success"
	// ResultError — reconcile 이 에러로 종료.
	ResultError = "error"
)

// 라벨 이름 상수 — trio 전체가 공유 (cardinality: namespace/name 기본).
const (
	labelNamespace = "namespace"
	labelName      = "name"
	labelResult    = "result"
	labelComponent = "component"
)

// latencyBuckets — 5ms ~ 30s 12개. typical reconcile + STS/PVC API roundtrip
// 범위 커버 (하위 = 빠른 steady reconcile, 상위 = init/scale/upgrade 분기).
// 원본 3 operator 와 byte-동일 (시계열 호환 절대 제약 — 변경 금지).
var latencyBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0,
}

// ReconcileMetrics 는 operator SLO reconcile metrics trio 를 보유한다.
// 필드는 export — 각 repo 가 기존 패키지 var 를 alias 로 선언해
// controller 콜사이트 무변경 마이그레이션이 가능하다.
type ReconcileMetrics struct {
	// Total — Reconcile 호출 횟수. 라벨: namespace, name.
	Total *prometheus.CounterVec
	// Latency — Reconcile wall-clock duration (초). SLO p50/p95/p99 를
	// PromQL histogram_quantile 로 산출. 라벨: namespace, name, result.
	Latency *prometheus.HistogramVec
	// Errors — component 별 reconcile 실패 횟수.
	// 라벨: namespace, name, component.
	Errors *prometheus.CounterVec
}

// New 는 subsystem 을 받아 reconcile metrics trio 를 생성한다.
// 노출 이름은 <subsystem>_reconcile_total /
// <subsystem>_reconcile_duration_seconds /
// <subsystem>_reconcile_errors_total — 기존 operator 의 subsystem 상수
// (mongodb / postgrescluster / valkey_cluster) 를 그대로 전달하면 기존
// 시계열 이름이 보존된다 (Namespace 필드 미사용도 원본 정합).
func New(subsystem string) *ReconcileMetrics {
	return &ReconcileMetrics{
		Total: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      "reconcile_total",
				Help:      "Total Reconcile invocations",
			},
			[]string{labelNamespace, labelName},
		),
		Latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      "reconcile_duration_seconds",
				Help:      "Reconcile function wall-clock duration in seconds",
				Buckets:   latencyBuckets,
			},
			[]string{labelNamespace, labelName, labelResult},
		),
		Errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      "reconcile_errors_total",
				Help:      "Total Reconcile component failures",
			},
			[]string{labelNamespace, labelName, labelComponent},
		),
	}
}

// MustRegister 는 주어진 Registerer 에 trio 를 등록한다. 각 operator 는
// init() 또는 main 에서 controller-runtime 의 metrics.Registry 를 전달한다.
// 중복 등록 시 prometheus 표준대로 panic 한다.
func (m *ReconcileMetrics) MustRegister(reg prometheus.Registerer) {
	reg.MustRegister(m.Total, m.Latency, m.Errors)
}

// IncTotal 은 Reconcile *진입* 시점에 호출 횟수를 1 증가시킨다 — 원본
// 3 operator 모두 Reconcile 첫머리에서 Inc (in-flight 도 카운트).
func (m *ReconcileMetrics) IncTotal(namespace, name string) {
	m.Total.WithLabelValues(namespace, name).Inc()
}

// ObserveReconcile 은 Reconcile *종료* (defer) 시점에 wall-clock duration 을
// 기록한다. result 는 ResultSuccess | ResultError — 반환 에러 기반 판정은
// ResultFor 사용. Total 은 진입 시점 IncTotal 책임 (시점 분리 — 원본 정합).
func (m *ReconcileMetrics) ObserveReconcile(namespace, name, result string, seconds float64) {
	m.Latency.WithLabelValues(namespace, name, result).Observe(seconds)
}

// IncError 는 component 별 reconcile 실패를 1 증가시킨다 (예: component =
// "statefulset" / "service" — 호출자 도메인이 결정).
func (m *ReconcileMetrics) IncError(namespace, name, component string) {
	m.Errors.WithLabelValues(namespace, name, component).Inc()
}

// ResultFor 는 reconcile 반환 에러를 result 라벨 값으로 변환한다 —
// 3 operator 가 defer closure 안에 복붙하던 success/error 판정의 흡수.
func ResultFor(err error) string {
	if err != nil {
		return ResultError
	}
	return ResultSuccess
}

// DeleteFor 는 CR 삭제 시 trio 의 해당 namespace/name 시계열을 제거해
// cardinality 누적을 방지한다. 원본 DeleteMetricsFor 의 trio 부분과 동일:
// Total 은 정확 라벨 삭제, Errors 는 component 차원이 있어 partial-match,
// Latency 는 result 라벨 (success/error) 순회 삭제.
//
// 도메인 메트릭은 각 repo 의 DeleteMetricsFor wrapper 가 본 메서드 호출
// 후 이어서 삭제한다.
func (m *ReconcileMetrics) DeleteFor(namespace, name string) {
	m.Total.DeleteLabelValues(namespace, name)
	m.Errors.DeletePartialMatch(prometheus.Labels{
		labelNamespace: namespace, labelName: name,
	})
	for _, r := range []string{ResultSuccess, ResultError} {
		m.Latency.DeleteLabelValues(namespace, name, r)
	}
}
