// SPDX-License-Identifier: Apache-2.0

package monitoring

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"sigs.k8s.io/yaml"
)

// updateGolden — `go test ./pkg/monitoring -run TestGolden -update` 로 골든 파일을
// 재생성한다. builder 의 *의도된* 출력 변경 시에만 사용하며, 평소 회귀 가드는
// 골든 파일과의 byte-diff = 0 을 강제한다.
var updateGolden = flag.Bool("update", false, "골든 매니페스트 파일 재생성")

// goldenCase — 단일 downstream-equivalence 케이스. obj 를 YAML 로 직렬화한 결과가
// testdata/<file> 과 정확히 일치해야 한다 (same input → same manifest output).
type goldenCase struct {
	// name — 하위 테스트 이름.
	name string
	// file — testdata 하위 골든 파일명.
	file string
	// build — builder 호출. 동일 입력은 항상 동일 출력이어야 한다.
	build func() any
}

// goldenCases — operator-commons 골든 회귀 코퍼스. 각 케이스는 downstream operator 가
// 실제로 생성하는 대표 매니페스트를 모사한다. ServiceMonitor (필수/전체옵션) +
// PrometheusRule (alert/recording 혼합) 를 커버해 builder 출력 drift 를 봉인한다.
func goldenCases() []goldenCase {
	return []goldenCase{
		{
			// 필수 필드만 — Prometheus Operator default 에 위임하는 최소 형태.
			name: "service_monitor_minimal",
			file: "service_monitor_minimal.yaml",
			build: func() any {
				return NewServiceMonitor(ServiceMonitorParams{
					Name:      "mongodb-metrics",
					Namespace: "default",
					Selector:  map[string]string{"app.kubernetes.io/name": "mongodb"},
					Port:      "metrics",
				})
			},
		},
		{
			// 전체 옵션 — endpoint 모든 필드 + namespaceSelector + custom labels.
			name: "service_monitor_full",
			file: "service_monitor_full.yaml",
			build: func() any {
				return NewServiceMonitor(ServiceMonitorParams{
					Name:              "valkey-metrics",
					Namespace:         "data",
					Labels:            map[string]string{"app.kubernetes.io/part-of": "valkey", "team": "data"},
					Selector:          map[string]string{"app.kubernetes.io/name": "valkey"},
					NamespaceSelector: []string{"data", "data-staging"},
					Port:              "metrics",
					Path:              "/custom-metrics",
					Interval:          "15s",
					ScrapeTimeout:     "10s",
					Scheme:            "https",
					HonorLabels:       true,
				})
			},
		},
		{
			// alert + recording rule 혼합 — RuleGroup 직렬화 순서/형태 봉인.
			name: "prometheus_rule_mixed",
			file: "prometheus_rule_mixed.yaml",
			build: func() any {
				return NewPrometheusRule(
					"valkey-alerts", "data",
					map[string]string{"app.kubernetes.io/instance": "valkey-cluster"},
					RuleGroup{
						Name:     "valkey.rules",
						Interval: "30s",
						Alerts: []AlertRule{
							{
								Alert:       "ValkeyDown",
								Expr:        `up{job="valkey"} == 0`,
								For:         "5m",
								Labels:      map[string]string{"severity": "critical"},
								Annotations: map[string]string{"summary": "Valkey instance down"},
							},
						},
						Records: []RecordingRule{
							{
								Record: "valkey:up:ratio",
								Expr:   `avg(up{job="valkey"})`,
								Labels: map[string]string{"team": "data"},
							},
						},
					},
				)
			},
		},
	}
}

// TestGolden_DownstreamEquivalence — builder 출력이 체크인된 골든 매니페스트와 정확히
// 일치함을 강제한다 (docs/ROADMAP.md pkg/monitoring 'Downstream equivalence e2e —
// same input → same manifest output' / 'Verify: golden file diff = 0').
//
// sigs.k8s.io/yaml 은 key 정렬 결정론적 출력을 보장하므로 byte-diff 비교가 안정적이다.
func TestGolden_DownstreamEquivalence(t *testing.T) {
	t.Parallel()
	for _, tc := range goldenCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			obj := tc.build()
			got, err := yaml.Marshal(obj)
			if err != nil {
				t.Fatalf("yaml.Marshal(%s) error: %v", tc.name, err)
			}
			goldenPath := filepath.Join("testdata", tc.file)
			if *updateGolden {
				if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
					t.Fatalf("골든 파일 쓰기 실패 %s: %v", goldenPath, err)
				}
				t.Logf("골든 갱신: %s", goldenPath)
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("골든 파일 읽기 실패 %s (최초 생성은 -update): %v", goldenPath, err)
			}
			if string(got) != string(want) {
				t.Errorf("manifest drift for %s\n--- got ---\n%s\n--- want (%s) ---\n%s",
					tc.name, got, goldenPath, want)
			}
		})
	}
}

// TestGolden_Deterministic — 동일 입력을 두 번 빌드·직렬화했을 때 byte-동일함을
// 검증한다. map 순회 비결정성이 출력에 새지 않음을 보장 (golden diff 안정성 전제).
func TestGolden_Deterministic(t *testing.T) {
	t.Parallel()
	for _, tc := range goldenCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			first, err := yaml.Marshal(tc.build())
			if err != nil {
				t.Fatalf("first marshal error: %v", err)
			}
			for i := range 16 {
				next, err := yaml.Marshal(tc.build())
				if err != nil {
					t.Fatalf("repeat marshal error: %v", err)
				}
				if string(first) != string(next) {
					t.Fatalf("비결정적 출력 (iter %d) for %s:\nfirst=%s\nnext=%s",
						i, tc.name, first, next)
				}
			}
		})
	}
}
