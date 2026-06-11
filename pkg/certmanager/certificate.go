// SPDX-License-Identifier: MIT

package certmanager

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// certManagerGroup — cert-manager API group (Certificate / Issuer 공통).
const certManagerGroup = "cert-manager.io"

// CertificateGVK 는 cert-manager Certificate CR 의 GroupVersionKind.
// 3 operator 가 verbatim 으로 중복 보유하던 선언의 단일 진실원.
var CertificateGVK = schema.GroupVersionKind{
	Group:   certManagerGroup,
	Version: "v1",
	Kind:    "Certificate",
}

// CertParams 는 BuildCertificate 의 입력. 전 필드 호출자 공급 — 숨은 default
// 는 IssuerKind 빈 값 → "Issuer" 1건뿐 (출력 byte-identical 보존 원칙).
type CertParams struct {
	// Name — Certificate metadata.name (3 operator 관례: "<cluster>-tls").
	Name string
	// Namespace — Certificate metadata.namespace.
	Namespace string
	// Labels — metadata.labels. nil 또는 빈 map 이면 미설정.
	Labels map[string]string
	// SecretName — spec.secretName. cert-manager 가 발급 결과를 저장할 Secret.
	SecretName string
	// CommonName — spec.commonName.
	CommonName string
	// DNSNames — spec.dnsNames (SAN). 목록 조립은 호출자 도메인 책임 —
	// Service 1개의 표준 FQDN 확장은 ServiceSANs 참조.
	DNSNames []string
	// IssuerName — spec.issuerRef.name.
	IssuerName string
	// IssuerKind — spec.issuerRef.kind. 빈 값이면 "Issuer" fallback
	// (mongo/postgres 패턴 — valkey 의 ClusterIssuer fallback 해석은
	// 호출자 resolveIssuer 책임으로 잔류).
	IssuerKind string
	// Duration — spec.duration (예: "2160h"). 빈 값이면 필드 생략 →
	// cert-manager default (90d) 위임.
	Duration string
	// RenewBefore — spec.renewBefore (예: "360h"). 빈 값이면 필드 생략 →
	// cert-manager default (15d) 위임.
	RenewBefore string
	// ECDSAPrivateKey — true 시 spec.privateKey 에
	// {algorithm: ECDSA, size: 256, rotationPolicy: Always} 블록 emit
	// (mongo/postgres 공통값). false 시 필드 생략 (valkey 패턴 —
	// cert-manager default 위임).
	ECDSAPrivateKey bool
}

// BuildCertificate 는 cert-manager Certificate CR 의 unstructured 표현을
// 반환한다. 호출자가 controller-runtime client 로 apply (CreateOrUpdate 등).
//
// spec.issuerRef.group 은 항상 "cert-manager.io" — cert-manager 는
// Issuer/ClusterIssuer 만 지원하며 빈 group 의 default 와 의미 동일.
// spec.usages 는 ["server auth", "client auth"] 고정 (3 operator 공통값).
//
// DNSNames 는 내부에서 []any 로 변환된다 — unstructured.SetNestedField (및
// 이후의 DeepCopy) 는 []string 에서 panic 하므로 본 변환이 호출자 footgun 을
// 1회 흡수한다 (postgres tls.go 원본 주석 참조).
func BuildCertificate(p CertParams) *unstructured.Unstructured {
	kind := p.IssuerKind
	if kind == "" {
		kind = "Issuer"
	}

	// []string → []any 변환: SetNestedField 의 DeepCopyJSONValue 는 []any 만 허용.
	dnsNames := make([]any, 0, len(p.DNSNames))
	for _, d := range p.DNSNames {
		dnsNames = append(dnsNames, d)
	}

	cert := &unstructured.Unstructured{}
	cert.SetGroupVersionKind(CertificateGVK)
	cert.SetName(p.Name)
	cert.SetNamespace(p.Namespace)
	if len(p.Labels) > 0 {
		cert.SetLabels(p.Labels)
	}

	spec := map[string]any{
		"secretName": p.SecretName,
		"commonName": p.CommonName,
		"dnsNames":   dnsNames,
		"issuerRef": map[string]any{
			"name":  p.IssuerName,
			"kind":  kind,
			"group": certManagerGroup,
		},
		"usages": []any{"server auth", "client auth"},
	}
	if p.ECDSAPrivateKey {
		spec["privateKey"] = map[string]any{
			"algorithm":      "ECDSA",
			"size":           int64(256),
			"rotationPolicy": "Always",
		}
	}
	if p.Duration != "" {
		spec["duration"] = p.Duration
	}
	if p.RenewBefore != "" {
		spec["renewBefore"] = p.RenewBefore
	}
	if err := unstructured.SetNestedField(cert.Object, spec, "spec"); err != nil {
		// programming error — spec 은 JSON-safe 단순 map 만 포함, 도달 불가
		// (mongo/postgres tls.go 원본과 동일 가드).
		return nil
	}
	return cert
}

// ServiceSANs 는 Kubernetes Service 1개의 4단 FQDN SAN 목록을 반환한다.
// mongo addSvc 클로저 ×2 + postgres 인라인 ×1 의 3 call site 복붙 흡수:
//
//	<svc>
//	<svc>.<ns>
//	<svc>.<ns>.svc
//	<svc>.<ns>.svc.cluster.local
//
// wildcard=true 시 StatefulSet per-pod DNS 커버용 와일드카드 1건 추가
// (mongo ReplicaSet 패턴):
//
//	*.<svc>.<ns>.svc.cluster.local
//
// 반환 슬라이스를 누적 (append) 해 CertParams.DNSNames 를 조립한다.
func ServiceSANs(svc, ns string, wildcard bool) []string {
	sans := []string{
		svc,
		fmt.Sprintf("%s.%s", svc, ns),
		fmt.Sprintf("%s.%s.svc", svc, ns),
		fmt.Sprintf("%s.%s.svc.cluster.local", svc, ns),
	}
	if wildcard {
		sans = append(sans, fmt.Sprintf("*.%s.%s.svc.cluster.local", svc, ns))
	}
	return sans
}
