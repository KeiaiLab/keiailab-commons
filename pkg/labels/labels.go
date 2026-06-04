// SPDX-License-Identifier: MIT

// Package labels — recommended Kubernetes labels (app.kubernetes.io/*) builder.
//
// downstream operator 가 동일 convention 으로 label 부여:
//
//	commonsabels.New("myapp", instance, "controller-manager", "1.4.6", "downstream operator")
//	→ {
//	    "app.kubernetes.io/name":       "myapp",
//	    "app.kubernetes.io/instance":   instance,
//	    "app.kubernetes.io/component":  "controller-manager",
//	    "app.kubernetes.io/version":    "1.4.6",
//	    "app.kubernetes.io/managed-by": "downstream operator",
//	    "app.kubernetes.io/part-of":    "downstream operator",
//	}
//
// SelectorLabels — Deployment / Service selector 용 *최소 부분* (immutable
// invariants 만). version 은 selector 미포함 (rolling update 시 변경).
package labels

// Set — Kubernetes 권장 label 집합 (app.kubernetes.io/*).
// 모든 필드 non-empty 가정 — caller 가 검증.
type Set struct {
	Name      string
	Instance  string
	Component string
	Version   string
	ManagedBy string
	PartOf    string
}

// All — 6개 label 모두 포함하는 map. Deployment/STS/Pod metadata.labels 용.
func (s Set) All() map[string]string {
	m := map[string]string{
		"app.kubernetes.io/name":       s.Name,
		"app.kubernetes.io/instance":   s.Instance,
		"app.kubernetes.io/managed-by": s.ManagedBy,
	}
	if s.Component != "" {
		m["app.kubernetes.io/component"] = s.Component
	}
	if s.Version != "" {
		m["app.kubernetes.io/version"] = s.Version
	}
	if s.PartOf != "" {
		m["app.kubernetes.io/part-of"] = s.PartOf
	}
	return m
}

// Selector — Deployment/Service.spec.selector.matchLabels 용. version 제외 —
// rolling update 시 selector 변경 불가 (k8s immutable field).
func (s Set) Selector() map[string]string {
	m := map[string]string{
		"app.kubernetes.io/name":     s.Name,
		"app.kubernetes.io/instance": s.Instance,
	}
	if s.Component != "" {
		m["app.kubernetes.io/component"] = s.Component
	}
	return m
}

// New — convenience constructor. Component / Version / PartOf 가 빈 문자열이면
// All() / Selector() 에서 누락. PartOf 미지정 시 ManagedBy 와 동일 사용 권장.
func New(name, instance, component, version, managedBy string) Set {
	return Set{
		Name:      name,
		Instance:  instance,
		Component: component,
		Version:   version,
		ManagedBy: managedBy,
		PartOf:    managedBy,
	}
}
