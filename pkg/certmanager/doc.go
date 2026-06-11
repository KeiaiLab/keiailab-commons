// SPDX-License-Identifier: MIT

// Package certmanager 는 cert-manager (cert-manager.io/v1) 의 Certificate /
// Issuer CR 을 unstructured 로 생성하는 공통 builder 를 제공한다. downstream
// operator 3종 (mongodb / postgres / valkey) 의 near-identical TLS 빌더 중복
// (~600 LOC — CertificateGVK 선언 verbatim ×3 + spec map 조립 + 4단 FQDN SAN
// 확장 클로저 복붙) 을 흡수한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator 의 *tls.go / certificate.go* 중복 코드를
// 추출하여 신규 도입되었다. downstream consumer 동시 적용 회귀 통과 후 Stable
// Tier 격상 예정 (자세한 격상 조건은 docs/ROADMAP.md §API Stability Tier).
//
// # 의존성 정책
//
// cert-manager Go SDK 를 의존하지 않는다. Certificate / Issuer 는
// *unstructured.Unstructured* 로 반환하고 호출자가 client 로 apply 한다 —
// pkg/monitoring 의 ServiceMonitor unstructured 패턴과 동일. 따라서
// cert-manager CRD 미설치 cluster 에서도 호출자 코드가 컴파일·동작하며,
// apply 시점의 NoMatchError 는 호출자가 fail-soft 처리할 수 있다.
// import 는 apimachinery 만 — client / 외부 SDK 의존 0.
//
// # 설계 원칙
//
//   - 전 필드 호출자 공급 + 숨은 default 최소화 — repo 별 기존 출력의
//     byte-identical 보존이 목표 (default 통일 금지 → 운영 cert 재발급
//     트리거 0). 유일한 fallback 은 IssuerKind 빈 값 → "Issuer"
//     (mongo/postgres 공통 패턴 — valkey 의 ClusterIssuer fallback 은
//     호출자 resolveIssuer 책임으로 잔류).
//   - Duration / RenewBefore 빈 값 = spec 필드 자체 생략 → cert-manager
//     default (90d / 15d) 위임.
//   - unstructured.SetNestedField 는 []any 만 deep copy 가능 — []string 을
//     넘기면 panic (postgres tls.go 원본 주석이 명시한 실제 footgun).
//     []string → []any 변환을 builder 내부에서 1회 흡수해 호출자에서 제거.
//   - SAN 목록 *조립* (shard ordinal 열거 등) 은 호출자 도메인 책임 — 본
//     패키지는 Service 1개의 4단 FQDN 확장 (ServiceSANs) 만 제공한다.
package certmanager
