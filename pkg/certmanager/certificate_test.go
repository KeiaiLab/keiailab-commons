// SPDX-License-Identifier: MIT

package certmanager

import (
	"reflect"
	"testing"
)

// composeDNS 는 원본 mongo addSvc 클로저 / postgres 인라인 루프와 동일한
// 방식으로 호출자 측 SAN 목록을 조립한다 (head = bare cluster name).
func composeDNS(head string, ns string, wildcard bool, svcs ...string) []string {
	dns := make([]string, 0, 1+len(svcs)*5)
	dns = append(dns, head)
	for _, svc := range svcs {
		dns = append(dns, ServiceSANs(svc, ns, wildcard)...)
	}
	return dns
}

// wantIssuerRef — 원본 mongo/postgres 의 issuerRef 블록 golden
// (group 은 항상 cert-manager.io).
func wantIssuerRef(name, kind string) map[string]any {
	return map[string]any{
		"name":  name,
		"kind":  kind,
		"group": "cert-manager.io",
	}
}

// wantUsages — 3 operator 공통 usages golden.
func wantUsages() []any {
	return []any{"server auth", "client auth"}
}

// wantECDSAPrivateKey — mongo/postgres 공통 privateKey 블록 golden.
func wantECDSAPrivateKey() map[string]any {
	return map[string]any{
		"algorithm":      "ECDSA",
		"size":           int64(256),
		"rotationPolicy": "Always",
	}
}

// TestBuildCertificate_Golden 은 추출 전 3 operator call site 의 unstructured
// 출력을 raw map 리터럴로 고정하고, 추출 후 BuildCertificate 가 동일 출력을
// 내는지 검증한다 — default 통일 금지 (운영 cert 재발급 트리거 0) 봉인.
func TestBuildCertificate_Golden(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params CertParams
		want   map[string]any
	}{
		{
			// mongodb-operator buildCertificate (sharded, shards=1) 형 —
			// IssuerKind 빈 값 fallback + duration/renewBefore 명시.
			name: "mongo_sharded",
			params: CertParams{
				Name:      "mycl-tls",
				Namespace: "data",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "mongodbsharded",
					"app.kubernetes.io/instance":   "mycl",
					"app.kubernetes.io/managed-by": "keiailab-mongodb-operator",
					"mongodb.keiailab.com/role":    "server-tls",
				},
				SecretName: "mycl-tls",
				CommonName: "mycl",
				DNSNames: composeDNS("mycl", "data", false,
					"mycl-mongos", "mycl-cfg-headless", "mycl-shard-0-headless"),
				IssuerName:      "cluster-ca",
				IssuerKind:      "", // → "Issuer" fallback (원본 동일)
				Duration:        "2160h",
				RenewBefore:     "360h",
				ECDSAPrivateKey: true,
			},
			want: map[string]any{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]any{
					"name":      "mycl-tls",
					"namespace": "data",
					"labels": map[string]any{
						"app.kubernetes.io/name":       "mongodbsharded",
						"app.kubernetes.io/instance":   "mycl",
						"app.kubernetes.io/managed-by": "keiailab-mongodb-operator",
						"mongodb.keiailab.com/role":    "server-tls",
					},
				},
				"spec": map[string]any{
					"secretName": "mycl-tls",
					"commonName": "mycl",
					"dnsNames": []any{
						"mycl",
						"mycl-mongos",
						"mycl-mongos.data",
						"mycl-mongos.data.svc",
						"mycl-mongos.data.svc.cluster.local",
						"mycl-cfg-headless",
						"mycl-cfg-headless.data",
						"mycl-cfg-headless.data.svc",
						"mycl-cfg-headless.data.svc.cluster.local",
						"mycl-shard-0-headless",
						"mycl-shard-0-headless.data",
						"mycl-shard-0-headless.data.svc",
						"mycl-shard-0-headless.data.svc.cluster.local",
					},
					"issuerRef":   wantIssuerRef("cluster-ca", "Issuer"),
					"usages":      wantUsages(),
					"privateKey":  wantECDSAPrivateKey(),
					"duration":    "2160h",
					"renewBefore": "360h",
				},
			},
		},
		{
			// mongodb-operator buildRSCertificate (ReplicaSet) 형 —
			// wildcard SAN 포함. 원본은 bare name 이 head + addSvc(name) 으로
			// 중복 등장 — builder 는 dedupe 없이 passthrough (원본 보존).
			name: "mongo_replicaset_wildcard",
			params: CertParams{
				Name:      "rs1-tls",
				Namespace: "db",
				Labels: map[string]string{
					"app.kubernetes.io/name":     "mongodb",
					"app.kubernetes.io/instance": "rs1",
				},
				SecretName:      "rs1-tls",
				CommonName:      "rs1",
				DNSNames:        composeDNS("rs1", "db", true, "rs1-headless", "rs1"),
				IssuerName:      "rs-ca",
				IssuerKind:      "Issuer",
				ECDSAPrivateKey: true,
			},
			want: map[string]any{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]any{
					"name":      "rs1-tls",
					"namespace": "db",
					"labels": map[string]any{
						"app.kubernetes.io/name":     "mongodb",
						"app.kubernetes.io/instance": "rs1",
					},
				},
				"spec": map[string]any{
					"secretName": "rs1-tls",
					"commonName": "rs1",
					"dnsNames": []any{
						"rs1",
						"rs1-headless",
						"rs1-headless.db",
						"rs1-headless.db.svc",
						"rs1-headless.db.svc.cluster.local",
						"*.rs1-headless.db.svc.cluster.local",
						"rs1",
						"rs1.db",
						"rs1.db.svc",
						"rs1.db.svc.cluster.local",
						"*.rs1.db.svc.cluster.local",
					},
					"issuerRef":  wantIssuerRef("rs-ca", "Issuer"),
					"usages":     wantUsages(),
					"privateKey": wantECDSAPrivateKey(),
				},
			},
		},
		{
			// postgres-operator buildCertificate 형 — duration/renewBefore
			// spec field 자체 부재 (cert-manager default 위임).
			name: "postgres_cluster",
			params: CertParams{
				Name:      "pg1-tls",
				Namespace: "infra",
				Labels: map[string]string{
					"app.kubernetes.io/name":     "postgrescluster",
					"app.kubernetes.io/instance": "pg1",
					"postgres.keiailab.io/role":  "server-tls",
				},
				SecretName:      "pg1-tls",
				CommonName:      "pg1",
				DNSNames:        composeDNS("pg1", "infra", false, "pg1-shard-0"),
				IssuerName:      "pg-ca",
				IssuerKind:      "",
				ECDSAPrivateKey: true,
			},
			want: map[string]any{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]any{
					"name":      "pg1-tls",
					"namespace": "infra",
					"labels": map[string]any{
						"app.kubernetes.io/name":     "postgrescluster",
						"app.kubernetes.io/instance": "pg1",
						"postgres.keiailab.io/role":  "server-tls",
					},
				},
				"spec": map[string]any{
					"secretName": "pg1-tls",
					"commonName": "pg1",
					"dnsNames": []any{
						"pg1",
						"pg1-shard-0",
						"pg1-shard-0.infra",
						"pg1-shard-0.infra.svc",
						"pg1-shard-0.infra.svc.cluster.local",
					},
					"issuerRef":  wantIssuerRef("pg-ca", "Issuer"),
					"usages":     wantUsages(),
					"privateKey": wantECDSAPrivateKey(),
				},
			},
		},
		{
			// valkey-operator BuildCertificateForValkey 골격 채택 형 —
			// privateKey 블록 부재 + ClusterIssuer 명시 + valkey 자체 SAN 형태
			// (ServiceSANs 미사용 — 호출자 조립). 원본 대비 issuerRef.group
			// 필드가 추가 emit 되는 의도적 deviation (cert-manager default 와
			// semantic 동일 — doc.go 설계 원칙 / 설계 입력 "골격만" 채택 정합).
			name: "valkey_skeleton",
			params: CertParams{
				Name:      "vk1-tls",
				Namespace: "cache",
				Labels: map[string]string{
					"app.kubernetes.io/instance": "vk1",
				},
				SecretName: "vk1-tls",
				CommonName: "vk1.cache.svc",
				DNSNames: []string{
					"vk1.cache.svc",
					"vk1.cache.svc.cluster.local",
					"vk1-headless.cache.svc",
					"*.vk1-headless.cache.svc",
				},
				IssuerName: "org-ca",
				IssuerKind: "ClusterIssuer",
				Duration:   "720h",
			},
			want: map[string]any{
				"apiVersion": "cert-manager.io/v1",
				"kind":       "Certificate",
				"metadata": map[string]any{
					"name":      "vk1-tls",
					"namespace": "cache",
					"labels": map[string]any{
						"app.kubernetes.io/instance": "vk1",
					},
				},
				"spec": map[string]any{
					"secretName": "vk1-tls",
					"commonName": "vk1.cache.svc",
					"dnsNames": []any{
						"vk1.cache.svc",
						"vk1.cache.svc.cluster.local",
						"vk1-headless.cache.svc",
						"*.vk1-headless.cache.svc",
					},
					"issuerRef": wantIssuerRef("org-ca", "ClusterIssuer"),
					"usages":    wantUsages(),
					"duration":  "720h",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BuildCertificate(tt.params)
			if got == nil {
				t.Fatal("BuildCertificate = nil")
			}
			if !reflect.DeepEqual(got.Object, tt.want) {
				t.Errorf("Object mismatch:\n got = %#v\nwant = %#v", got.Object, tt.want)
			}
		})
	}
}

// TestBuildCertificate_IssuerKindFallback 은 IssuerKind 빈 값 → "Issuer"
// fallback (mongo/postgres 패턴) 과 명시 값 보존을 검증한다.
func TestBuildCertificate_IssuerKindFallback(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		kind     string
		wantKind string
	}{
		{name: "empty_falls_back_to_Issuer", kind: "", wantKind: "Issuer"},
		{name: "explicit_Issuer_preserved", kind: "Issuer", wantKind: "Issuer"},
		{name: "explicit_ClusterIssuer_preserved", kind: "ClusterIssuer", wantKind: "ClusterIssuer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BuildCertificate(CertParams{
				Name: "c", Namespace: "n", IssuerName: "i", IssuerKind: tt.kind,
			})
			ref := got.Object["spec"].(map[string]any)["issuerRef"].(map[string]any)
			if ref["kind"] != tt.wantKind {
				t.Errorf("issuerRef.kind = %v, want %v", ref["kind"], tt.wantKind)
			}
		})
	}
}

// TestBuildCertificate_OptionalFields 는 Duration / RenewBefore /
// ECDSAPrivateKey 의 빈 값 시 spec 필드 생략 (cert-manager default 위임) 을
// 검증한다.
func TestBuildCertificate_OptionalFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		duration        string
		renewBefore     string
		ecdsa           bool
		wantDuration    bool
		wantRenewBefore bool
		wantPrivateKey  bool
	}{
		{name: "all_empty_all_omitted"},
		{name: "duration_only", duration: "2160h", wantDuration: true},
		{name: "renewBefore_only", renewBefore: "360h", wantRenewBefore: true},
		{name: "ecdsa_only", ecdsa: true, wantPrivateKey: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BuildCertificate(CertParams{
				Name: "c", Namespace: "n", IssuerName: "i",
				Duration: tt.duration, RenewBefore: tt.renewBefore, ECDSAPrivateKey: tt.ecdsa,
			})
			spec := got.Object["spec"].(map[string]any)
			checks := []struct {
				key  string
				want bool
			}{
				{key: "duration", want: tt.wantDuration},
				{key: "renewBefore", want: tt.wantRenewBefore},
				{key: "privateKey", want: tt.wantPrivateKey},
			}
			for _, c := range checks {
				if _, ok := spec[c.key]; ok != c.want {
					t.Errorf("spec[%q] presence = %v, want %v", c.key, ok, c.want)
				}
			}
		})
	}
}

// TestBuildCertificate_NilLabels 는 Labels nil 시 metadata.labels 가 아예
// 설정되지 않음을 검증한다.
func TestBuildCertificate_NilLabels(t *testing.T) {
	t.Parallel()
	got := BuildCertificate(CertParams{Name: "c", Namespace: "n", IssuerName: "i"})
	meta := got.Object["metadata"].(map[string]any)
	if _, ok := meta["labels"]; ok {
		t.Errorf("metadata.labels 가 설정됨 (nil Labels 는 미설정이어야 함): %v", meta["labels"])
	}
}

// TestBuildCertificate_DeepCopySafe 는 []string → []any 변환이 SetNestedField
// / DeepCopy 의 panic footgun 을 흡수했는지 검증한다 — dnsNames 가 []string
// 으로 저장되면 DeepCopy 가 panic 한다 (postgres tls.go 원본 주석의 사고).
func TestBuildCertificate_DeepCopySafe(t *testing.T) {
	t.Parallel()
	got := BuildCertificate(CertParams{
		Name: "c", Namespace: "n", IssuerName: "i",
		DNSNames: []string{"a", "b"},
	})
	// DeepCopy 가 panic 하면 본 테스트가 즉시 실패 — 회귀 가드.
	cp := got.DeepCopy()
	dns, ok := cp.Object["spec"].(map[string]any)["dnsNames"].([]any)
	if !ok {
		t.Fatalf("spec.dnsNames 타입이 []any 가 아님: %T",
			cp.Object["spec"].(map[string]any)["dnsNames"])
	}
	if !reflect.DeepEqual(dns, []any{"a", "b"}) {
		t.Errorf("dnsNames = %v, want [a b]", dns)
	}
}

// TestServiceSANs 는 4단 FQDN 확장 (+ wildcard 5단) 의 정확한 순서·내용을
// 검증한다 — 원본 mongo addSvc / postgres 인라인과 동일.
func TestServiceSANs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		svc      string
		ns       string
		wildcard bool
		want     []string
	}{
		{
			name: "four_forms",
			svc:  "web", ns: "apps",
			want: []string{
				"web",
				"web.apps",
				"web.apps.svc",
				"web.apps.svc.cluster.local",
			},
		},
		{
			name: "wildcard_adds_fifth_form",
			svc:  "web", ns: "apps", wildcard: true,
			want: []string{
				"web",
				"web.apps",
				"web.apps.svc",
				"web.apps.svc.cluster.local",
				"*.web.apps.svc.cluster.local",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ServiceSANs(tt.svc, tt.ns, tt.wildcard)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceSANs = %v, want %v", got, tt.want)
			}
		})
	}
}
