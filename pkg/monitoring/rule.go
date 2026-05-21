package monitoring

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PrometheusRuleGVR — kube-prometheus-stack PrometheusRule.
var PrometheusRuleGVR = schema.GroupVersionResource{
	Group:    "monitoring.coreos.com",
	Version:  "v1",
	Resource: "prometheusrules",
}

// AlertRule — Prometheus alerting rule (alert + expr + for + labels + annotations).
type AlertRule struct {
	Alert       string            `json:"alert"`
	Expr        string            `json:"expr"`
	For         string            `json:"for,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// RecordingRule — Prometheus recording rule (record + expr + labels).
type RecordingRule struct {
	Record string            `json:"record"`
	Expr   string            `json:"expr"`
	Labels map[string]string `json:"labels,omitempty"`
}

// RuleGroup — Prometheus RuleGroup (name + interval + rules).
type RuleGroup struct {
	Name     string `json:"name"`
	Interval string `json:"interval,omitempty"`
	Alerts   []AlertRule
	Records  []RecordingRule
}

// NewPrometheusRule — kube-prometheus-stack PrometheusRule unstructured object.
//
// 사용 예:
//
//	rule := monitoring.NewPrometheusRule("vk-alerts", "default",
//	    map[string]string{"app.kubernetes.io/instance": "vk-cluster"},
//	    monitoring.RuleGroup{
//	        Name:     "valkey.rules",
//	        Interval: "30s",
//	        Alerts: []monitoring.AlertRule{
//	            {
//	                Alert:       "ValkeyDown",
//	                Expr:        "up{job=\"valkey\"} == 0",
//	                For:         "5m",
//	                Labels:      map[string]string{"severity": "critical"},
//	                Annotations: map[string]string{"summary": "Valkey instance down"},
//	            },
//	        },
//	    },
//	)
//
// Refs: docs/ROADMAP.md 'PrometheusRule 빌더 (alert/recording rules 공통화)'
func NewPrometheusRule(
	name, namespace string,
	labels map[string]string,
	groups ...RuleGroup,
) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "monitoring.coreos.com",
		Version: "v1",
		Kind:    "PrometheusRule",
	})
	obj.SetName(name)
	obj.SetNamespace(namespace)
	if len(labels) > 0 {
		obj.SetLabels(labels)
	}

	rawGroups := make([]any, 0, len(groups))
	for _, g := range groups {
		rules := make([]any, 0, len(g.Alerts)+len(g.Records))
		for _, a := range g.Alerts {
			r := map[string]any{
				"alert": a.Alert,
				"expr":  a.Expr,
			}
			if a.For != "" {
				r["for"] = a.For
			}
			if len(a.Labels) > 0 {
				r["labels"] = anyMap(a.Labels)
			}
			if len(a.Annotations) > 0 {
				r["annotations"] = anyMap(a.Annotations)
			}
			rules = append(rules, r)
		}
		for _, rec := range g.Records {
			r := map[string]any{
				"record": rec.Record,
				"expr":   rec.Expr,
			}
			if len(rec.Labels) > 0 {
				r["labels"] = anyMap(rec.Labels)
			}
			rules = append(rules, r)
		}
		gm := map[string]any{
			"name":  g.Name,
			"rules": rules,
		}
		if g.Interval != "" {
			gm["interval"] = g.Interval
		}
		rawGroups = append(rawGroups, gm)
	}
	_ = unstructured.SetNestedField(obj.Object, rawGroups, "spec", "groups")
	return obj
}

func anyMap(m map[string]string) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
