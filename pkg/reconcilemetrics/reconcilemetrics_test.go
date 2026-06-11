// SPDX-License-Identifier: MIT

package reconcilemetrics_test

import (
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"

	"github.com/keiailab/keiailab-commons/pkg/reconcilemetrics"
)

// 테스트 공용 fixture 라벨 — goconst 회피 + 의도 명시.
const (
	testNS   = "default"
	testName = "sample"
	compSTS  = "statefulset"
	compSVC  = "service"
)

// gatherFamilies 는 trio 에 시계열을 1개씩 생성한 뒤 PedanticRegistry 로
// Gather 해 family name → MetricFamily 맵을 돌려준다 (테스트 공용 헬퍼).
func gatherFamilies(t *testing.T, m *reconcilemetrics.ReconcileMetrics) map[string]*dto.MetricFamily {
	t.Helper()

	reg := prometheus.NewPedanticRegistry()
	m.MustRegister(reg)

	m.IncTotal(testNS, testName)
	m.ObserveReconcile(testNS, testName, reconcilemetrics.ResultSuccess, 0.01)
	m.IncError(testNS, testName, compSTS)

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather 실패: %v", err)
	}
	got := make(map[string]*dto.MetricFamily, len(families))
	for _, f := range families {
		got[f.GetName()] = f
	}
	return got
}

// TestNew_MetricNamesPerSubsystem 은 3 operator 의 기존 subsystem 상수를
// 그대로 전달했을 때 기존 시계열 이름이 byte-동일하게 보존되는지 검증한다
// — 시계열 호환 절대 제약 (설계 입력 revised).
//
// AAA:
//
//	Arrange — 각 operator 의 실제 subsystem 상수로 New
//	Act — 시계열 생성 후 Gather
//	Assert — family 이름 3종이 원본 metrics.go 노출 이름과 정확 일치
func TestNew_MetricNamesPerSubsystem(t *testing.T) {
	cases := []struct {
		subsystem string // 원본 repo 의 metricSubsystem 상수
		wantNames []string
	}{
		{
			subsystem: "mongodb", // mongodb-operator internal/controller/metrics.go
			wantNames: []string{
				"mongodb_reconcile_duration_seconds",
				"mongodb_reconcile_errors_total",
				"mongodb_reconcile_total",
			},
		},
		{
			subsystem: "postgrescluster", // postgres-operator internal/controller/metrics.go
			wantNames: []string{
				"postgrescluster_reconcile_duration_seconds",
				"postgrescluster_reconcile_errors_total",
				"postgrescluster_reconcile_total",
			},
		},
		{
			subsystem: "valkey_cluster", // valkey-operator internal/controller/metrics.go
			wantNames: []string{
				"valkey_cluster_reconcile_duration_seconds",
				"valkey_cluster_reconcile_errors_total",
				"valkey_cluster_reconcile_total",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.subsystem, func(t *testing.T) {
			families := gatherFamilies(t, reconcilemetrics.New(tc.subsystem))

			gotNames := make([]string, 0, len(families))
			for name := range families {
				gotNames = append(gotNames, name)
			}
			sort.Strings(gotNames)

			if len(gotNames) != len(tc.wantNames) {
				t.Fatalf("family 수 불일치: got %v, want %v", gotNames, tc.wantNames)
			}
			for i, want := range tc.wantNames {
				if gotNames[i] != want {
					t.Errorf("metric 이름 [%d]: got %q, want %q", i, gotNames[i], want)
				}
			}
		})
	}
}

// TestNew_HelpStringsMatchOriginals 는 Help 문자열 3종이 원본과 byte-동일한지
// 검증한다 — Help 변경은 Prometheus 상 metric 재정의로 간주될 수 있어
// 호환 절대 제약에 포함된다.
func TestNew_HelpStringsMatchOriginals(t *testing.T) {
	families := gatherFamilies(t, reconcilemetrics.New("test"))

	cases := []struct {
		family   string
		wantHelp string // 원본 3 repo metrics.go 와 byte-동일
	}{
		{"test_reconcile_total", "Total Reconcile invocations"},
		{"test_reconcile_duration_seconds", "Reconcile function wall-clock duration in seconds"},
		{"test_reconcile_errors_total", "Total Reconcile component failures"},
	}
	for _, tc := range cases {
		t.Run(tc.family, func(t *testing.T) {
			f, ok := families[tc.family]
			if !ok {
				t.Fatalf("family %q 부재", tc.family)
			}
			if got := f.GetHelp(); got != tc.wantHelp {
				t.Errorf("Help 불일치: got %q, want %q", got, tc.wantHelp)
			}
		})
	}
}

// TestNew_LatencyBucketsMatchOriginals 는 Histogram bucket 경계 12개가 원본
// (5ms ~ 30s) 과 정확히 일치하는지 검증한다 — bucket 변경 = 기존
// histogram_quantile dashboard 의 단절.
func TestNew_LatencyBucketsMatchOriginals(t *testing.T) {
	want := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0}

	families := gatherFamilies(t, reconcilemetrics.New("test"))
	f, ok := families["test_reconcile_duration_seconds"]
	if !ok {
		t.Fatal("duration family 부재")
	}
	hist := f.GetMetric()[0].GetHistogram()

	got := make([]float64, 0, len(want))
	for _, b := range hist.GetBucket() {
		if math.IsInf(b.GetUpperBound(), +1) {
			continue // +Inf 는 정의 bucket 아님 (포맷 표현)
		}
		got = append(got, b.GetUpperBound())
	}

	if len(got) != len(want) {
		t.Fatalf("bucket 수 불일치: got %d (%v), want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("bucket [%d]: got %v, want %v", i, got[i], want[i])
		}
	}
}

// TestCounters_ExpositionGolden 은 counter 2종의 exposition text 를 원본
// 노출 형태와 byte-수준으로 고정한다 (이름 + HELP + TYPE + 라벨 구조).
func TestCounters_ExpositionGolden(t *testing.T) {
	m := reconcilemetrics.New("mongodb")
	m.IncTotal(testNS, testName)
	m.IncError(testNS, testName, compSTS)

	cases := []struct {
		metric prometheus.Collector
		name   string
		golden string
	}{
		{
			metric: m.Total,
			name:   "mongodb_reconcile_total",
			golden: `
# HELP mongodb_reconcile_total Total Reconcile invocations
# TYPE mongodb_reconcile_total counter
mongodb_reconcile_total{name="sample",namespace="default"} 1
`,
		},
		{
			metric: m.Errors,
			name:   "mongodb_reconcile_errors_total",
			golden: `
# HELP mongodb_reconcile_errors_total Total Reconcile component failures
# TYPE mongodb_reconcile_errors_total counter
mongodb_reconcile_errors_total{component="statefulset",name="sample",namespace="default"} 1
`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := testutil.CollectAndCompare(tc.metric, strings.NewReader(tc.golden), tc.name); err != nil {
				t.Errorf("exposition golden 불일치: %v", err)
			}
		})
	}
}

// TestObserveReconcile_RecordsLatencyOnly 는 ObserveReconcile 이 Latency 만
// 기록하고 Total 은 건드리지 않음을 검증한다 — 원본 시점 분리 (Total 은
// 진입 IncTotal, Latency 는 종료 defer) 시맨틱 보존.
func TestObserveReconcile_RecordsLatencyOnly(t *testing.T) {
	m := reconcilemetrics.New("test")

	m.ObserveReconcile(testNS, testName, reconcilemetrics.ResultError, 1.5)

	if got := testutil.CollectAndCount(m.Latency); got != 1 {
		t.Errorf("Latency 시계열 수: got %d, want 1", got)
	}
	if got := testutil.CollectAndCount(m.Total); got != 0 {
		t.Errorf("Total 은 미기록이어야 함: got %d 시계열", got)
	}

	// sum/count 로 관측 값 자체 검증.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(m.Latency)
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather 실패: %v", err)
	}
	hist := families[0].GetMetric()[0].GetHistogram()
	if hist.GetSampleCount() != 1 {
		t.Errorf("sample count: got %d, want 1", hist.GetSampleCount())
	}
	if hist.GetSampleSum() != 1.5 {
		t.Errorf("sample sum: got %v, want 1.5", hist.GetSampleSum())
	}
}

// TestIncTotalAndIncError_IncrementCounters 는 진입/실패 카운터 증가를
// 검증한다 (원본 콜사이트: Reconcile 첫머리 + applyErrorCondition).
func TestIncTotalAndIncError_IncrementCounters(t *testing.T) {
	m := reconcilemetrics.New("test")

	m.IncTotal(testNS, testName)
	m.IncTotal(testNS, testName)
	m.IncError(testNS, testName, compSTS)

	if got := testutil.ToFloat64(m.Total.WithLabelValues(testNS, testName)); got != 2 {
		t.Errorf("Total: got %v, want 2", got)
	}
	if got := testutil.ToFloat64(m.Errors.WithLabelValues(testNS, testName, compSTS)); got != 1 {
		t.Errorf("Errors: got %v, want 1", got)
	}
}

// TestResultFor 는 reconcile 반환 에러 → result 라벨 변환을 검증한다 —
// 3 operator defer closure 의 success/error 판정 복붙 흡수.
func TestResultFor(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"nil 에러는 success", nil, reconcilemetrics.ResultSuccess},
		{"non-nil 에러는 error", errors.New("boom"), reconcilemetrics.ResultError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := reconcilemetrics.ResultFor(tc.err); got != tc.want {
				t.Errorf("ResultFor(%v): got %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}

// TestDeleteFor_RemovesOnlyTargetCR 은 CR 삭제 시 해당 namespace/name 의
// trio 시계열만 제거되고 다른 CR 시계열은 보존됨을 검증한다 — cardinality
// 정리의 핵심 행동 (원본 DeleteMetricsFor trio 부분과 동일 거동).
func TestDeleteFor_RemovesOnlyTargetCR(t *testing.T) {
	type cr struct{ ns, name string }
	target := cr{"ns-a", "cr-a"}
	other := cr{"ns-b", "cr-b"}

	m := reconcilemetrics.New("test")
	for _, c := range []cr{target, other} {
		m.IncTotal(c.ns, c.name)
		m.ObserveReconcile(c.ns, c.name, reconcilemetrics.ResultSuccess, 0.01)
		m.ObserveReconcile(c.ns, c.name, reconcilemetrics.ResultError, 0.02)
		m.IncError(c.ns, c.name, compSTS)
		m.IncError(c.ns, c.name, compSVC)
	}

	// 사전 상태: CR 2개 × (Total 1 / Latency result 2 / Errors component 2).
	pre := []struct {
		label string
		col   prometheus.Collector
		want  int
	}{
		{"Total", m.Total, 2},
		{"Latency", m.Latency, 4},
		{"Errors", m.Errors, 4},
	}
	for _, p := range pre {
		if got := testutil.CollectAndCount(p.col); got != p.want {
			t.Fatalf("사전 %s 시계열 수: got %d, want %d", p.label, got, p.want)
		}
	}

	m.DeleteFor(target.ns, target.name)

	post := []struct {
		label string
		col   prometheus.Collector
		want  int
	}{
		{"Total", m.Total, 1},
		{"Latency", m.Latency, 2},
		{"Errors", m.Errors, 2},
	}
	for _, p := range post {
		if got := testutil.CollectAndCount(p.col); got != p.want {
			t.Errorf("사후 %s 시계열 수: got %d, want %d", p.label, got, p.want)
		}
	}

	// 잔존 시계열이 other CR 의 것인지 값으로 확인 (target 부활 방지를 위해
	// target 라벨은 재조회하지 않는다 — WithLabelValues 는 시계열을 재생성).
	if got := testutil.ToFloat64(m.Total.WithLabelValues(other.ns, other.name)); got != 1 {
		t.Errorf("other CR Total 보존 실패: got %v, want 1", got)
	}
}

// TestConcurrentUse_RaceFree 는 복수 reconcile goroutine 의 기록과 CR 삭제
// (DeleteFor) 가 동시 진행돼도 안전함을 -race 로 검증한다 — controller
// worker 동시성 + 삭제 cleanup 의 실사용 패턴.
func TestConcurrentUse_RaceFree(t *testing.T) {
	m := reconcilemetrics.New("race")

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			ns := testNS
			if worker%2 == 0 {
				ns = "ns-even"
			}
			for j := range 50 {
				m.IncTotal(ns, testName)
				m.ObserveReconcile(ns, testName, reconcilemetrics.ResultFor(nil), 0.001)
				m.IncError(ns, testName, compSTS)
				if j%10 == 0 {
					m.DeleteFor(ns, testName)
				}
			}
		}(i)
	}
	wg.Wait()
}
