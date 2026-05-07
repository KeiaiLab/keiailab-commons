package monitoring

import (
	"testing"
)

func TestNewServiceMonitor_RequiredFields(t *testing.T) {
	t.Parallel()
	sm := NewServiceMonitor(ServiceMonitorParams{
		Name:      "mongodb-metrics",
		Namespace: "default",
		Selector:  map[string]string{"app": "mongodb"},
		Port:      "metrics",
	})
	if sm.GetAPIVersion() != "monitoring.coreos.com/v1" {
		t.Errorf("APIVersion = %q, want monitoring.coreos.com/v1", sm.GetAPIVersion())
	}
	if sm.GetKind() != "ServiceMonitor" {
		t.Errorf("Kind = %q, want ServiceMonitor", sm.GetKind())
	}
	if sm.GetName() != "mongodb-metrics" {
		t.Errorf("Name = %q", sm.GetName())
	}
	if sm.GetNamespace() != "default" {
		t.Errorf("Namespace = %q", sm.GetNamespace())
	}
}

func TestNewServiceMonitor_AllOptional(t *testing.T) {
	t.Parallel()
	sm := NewServiceMonitor(ServiceMonitorParams{
		Name:              "vk-metrics",
		Namespace:         "data",
		Labels:            map[string]string{"team": "data"},
		Selector:          map[string]string{"app.kubernetes.io/name": "valkey"},
		NamespaceSelector: []string{"data", "data-staging"},
		Port:              "metrics",
		Path:              "/custom-metrics",
		Interval:          "15s",
		ScrapeTimeout:     "10s",
		Scheme:            "https",
		HonorLabels:       true,
	})
	if got := sm.GetLabels()["team"]; got != "data" {
		t.Errorf("custom label not set: %v", sm.GetLabels())
	}
	endpoints, found, err := unstructuredField(sm.Object, "spec", "endpoints")
	if err != nil || !found {
		t.Fatalf("spec.endpoints not found: %v / %v", found, err)
	}
	first := endpoints.([]any)[0].(map[string]any)
	if first["port"] != "metrics" {
		t.Errorf("port = %v, want metrics", first["port"])
	}
	if first["path"] != "/custom-metrics" {
		t.Errorf("path missing")
	}
	if first["interval"] != "15s" {
		t.Errorf("interval missing")
	}
	if first["scrapeTimeout"] != "10s" {
		t.Errorf("scrapeTimeout missing")
	}
	if first["scheme"] != "https" {
		t.Errorf("scheme missing")
	}
	if first["honorLabels"] != true {
		t.Errorf("honorLabels missing")
	}
	// NamespaceSelector 검증
	nsField, _, _ := unstructuredField(sm.Object, "spec", "namespaceSelector")
	matchNames := nsField.(map[string]any)["matchNames"].([]any)
	if len(matchNames) != 2 || matchNames[0] != "data" || matchNames[1] != "data-staging" {
		t.Errorf("namespaceSelector.matchNames mismatch: %v", matchNames)
	}
}

func TestNewServiceMonitor_OptionalDefaults(t *testing.T) {
	t.Parallel()
	sm := NewServiceMonitor(ServiceMonitorParams{
		Name: "n", Namespace: "ns",
		Selector: map[string]string{"k": "v"},
		Port:     "p",
	})
	endpoints, _, _ := unstructuredField(sm.Object, "spec", "endpoints")
	first := endpoints.([]any)[0].(map[string]any)
	for _, key := range []string{"path", "interval", "scrapeTimeout", "scheme", "honorLabels"} {
		if _, ok := first[key]; ok {
			t.Errorf("optional %q should be absent when zero-value", key)
		}
	}
	if sm.GetLabels() != nil && len(sm.GetLabels()) > 0 {
		t.Errorf("labels should be nil/empty when not specified, got %v", sm.GetLabels())
	}
	// NamespaceSelector 미지정 시 spec 에 나타나면 안됨.
	if _, ok, _ := unstructuredField(sm.Object, "spec", "namespaceSelector"); ok {
		t.Error("namespaceSelector should be absent when slice is empty/nil")
	}
}

// unstructuredField — 작은 helper, depth 2 path 만.
func unstructuredField(obj map[string]any, path ...string) (any, bool, error) {
	cur := any(obj)
	for _, p := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false, nil
		}
		cur, ok = m[p]
		if !ok {
			return nil, false, nil
		}
	}
	return cur, true, nil
}
