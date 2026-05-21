// SPDX-License-Identifier: Apache-2.0

package monitoring

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewPrometheusRule(t *testing.T) {
	r := NewPrometheusRule("test", "default", map[string]string{"k": "v"},
		RuleGroup{
			Name:     "test.rules",
			Interval: "30s",
			Alerts: []AlertRule{{
				Alert:       "TestAlert",
				Expr:        "up == 0",
				For:         "5m",
				Labels:      map[string]string{"severity": "critical"},
				Annotations: map[string]string{"summary": "test"},
			}},
			Records: []RecordingRule{{
				Record: "test:up:sum",
				Expr:   "sum(up)",
			}},
		},
	)
	if r.GetName() != "test" {
		t.Errorf("name = %q, want test", r.GetName())
	}
	if r.GetKind() != "PrometheusRule" {
		t.Errorf("kind = %q, want PrometheusRule", r.GetKind())
	}

	groups, found, err := unstructured.NestedSlice(r.Object, "spec", "groups")
	if err != nil || !found {
		t.Fatalf("groups not found: found=%v err=%v", found, err)
	}
	if len(groups) != 1 {
		t.Fatalf("groups len = %d, want 1", len(groups))
	}
	g, ok := groups[0].(map[string]any)
	if !ok {
		t.Fatalf("group[0] is not map[string]any")
	}
	if g["name"] != "test.rules" {
		t.Errorf("group name = %v, want test.rules", g["name"])
	}
	rules, ok := g["rules"].([]any)
	if !ok {
		t.Fatalf("rules is not []any")
	}
	if len(rules) != 2 {
		t.Errorf("rules len = %d, want 2 (1 alert + 1 record)", len(rules))
	}
}
