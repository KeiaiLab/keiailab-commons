// SPDX-License-Identifier: Apache-2.0

// Package monitoring — Prometheus Operator ServiceMonitor builder.
//
// downstream operator 가 동일 convention 으로 ServiceMonitor 생성:
//
//	sm := monitoring.NewServiceMonitor(monitoring.ServiceMonitorParams{
//	    Name:      "myapp-metrics",
//	    Namespace: ns,
//	    Selector:  labels.New("example-app", instance, "metrics", "", "downstream-operator").Selector(),
//	    Port:      "metrics",
//	    Interval:  "30s",
//	    Scheme:    "https",
//	})
//
// CRD 종속 회피: ServiceMonitor 는 monitoring.coreos.com/v1 (Prometheus Operator)
// 의 typed API. 본 패키지는 *unstructured.Unstructured* 반환 — 호출자가 client
// 로 apply. 이는 Prometheus Operator 미설치 환경에서도 코드 컴파일 가능.
package monitoring

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ServiceMonitorParams — ServiceMonitor 생성 매개변수.
type ServiceMonitorParams struct {
	// Name — ServiceMonitor metadata.name.
	Name string
	// Namespace — ServiceMonitor metadata.namespace.
	Namespace string
	// Labels — ServiceMonitor metadata.labels (기존 labels 패키지 .All() 등).
	Labels map[string]string
	// Selector — spec.selector.matchLabels (Service 가 가진 label 매칭).
	Selector map[string]string
	// NamespaceSelector — spec.namespaceSelector.matchNames (target Service 의
	// namespace 제한). 빈 슬라이스 시 미적용 (Prometheus default = 모든 namespace).
	NamespaceSelector []string
	// Port — spec.endpoints[0].port (Service의 port name).
	Port string
	// Path — spec.endpoints[0].path (default /metrics).
	Path string
	// Interval — spec.endpoints[0].interval (default 30s).
	Interval string
	// ScrapeTimeout — spec.endpoints[0].scrapeTimeout. 빈 문자열 = Prometheus default.
	ScrapeTimeout string
	// Scheme — http | https (default http).
	Scheme string
	// HonorLabels — spec.endpoints[0].honorLabels.
	HonorLabels bool
}

// NewServiceMonitor — Prometheus Operator ServiceMonitor (monitoring.coreos.com/v1)
// 의 unstructured 표현. 호출자가 controller-runtime client 로 Apply.
//
// 빈 필드는 적용 안 됨 (Prometheus Operator 의 default 적용).
func NewServiceMonitor(p ServiceMonitorParams) *unstructured.Unstructured {
	endpoint := map[string]any{
		"port": p.Port,
	}
	if p.Path != "" {
		endpoint["path"] = p.Path
	}
	if p.Interval != "" {
		endpoint["interval"] = p.Interval
	}
	if p.ScrapeTimeout != "" {
		endpoint["scrapeTimeout"] = p.ScrapeTimeout
	}
	if p.Scheme != "" {
		endpoint["scheme"] = p.Scheme
	}
	if p.HonorLabels {
		endpoint["honorLabels"] = true
	}

	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("monitoring.coreos.com/v1")
	obj.SetKind("ServiceMonitor")
	obj.SetName(p.Name)
	obj.SetNamespace(p.Namespace)
	if len(p.Labels) > 0 {
		obj.SetLabels(p.Labels)
	}

	spec := map[string]any{
		"selector":  map[string]any{"matchLabels": stringMapToAny(p.Selector)},
		"endpoints": []any{endpoint},
	}
	if len(p.NamespaceSelector) > 0 {
		ns := make([]any, len(p.NamespaceSelector))
		for i, n := range p.NamespaceSelector {
			ns[i] = n
		}
		spec["namespaceSelector"] = map[string]any{"matchNames": ns}
	}
	_ = unstructured.SetNestedMap(obj.Object, spec, "spec")
	return obj
}

func stringMapToAny(in map[string]string) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
