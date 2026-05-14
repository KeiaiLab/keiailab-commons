# ROADMAP — operator-commons

본 ROADMAP 은 라이브러리 진화 방향을 *API stability tier* + *v1.0.0 졸업 조건* + *패키지별 보강 항목* 으로 추적한다. 분기/날짜 기반 로드맵을 두지 않는다 (글로벌 `standards/workflow.md` "시간 기반 로드맵 금지"). 라이브러리 특성상 *호출자(operator 3종)의 요구* 가 진화 방향을 결정한다.

## 체크박스 의미

| 마커 | 의미 |
|---|---|
| `[x]` | 코드 + 테스트 양쪽 존재. 호출자 3 repo 검증 완료 |
| `[~]` | 부분 구현 (helper 존재, 검증 미완) |
| `[ ]` | 미시작 |

## API Stability Tier (현행 v0.7.x)

| 패키지 | Tier | 사용처 | Tier 격상 조건 |
|---|---|---|---|
| `pkg/finalizer` | **Stable** | mongo / pg / valkey | v1 진입 (별도 작업 없음) |
| `pkg/labels` | **Stable** | mongo / pg / valkey | v1 진입 (별도 작업 없음) |
| `pkg/status` | **Stable** | mongo / pg / valkey | v1 진입 (별도 작업 없음) |
| `pkg/version` (incl. Matrix) | Beta | mongo / pg / valkey | generic `Matrix[E]` cross-repo 검증 완료 |
| `pkg/monitoring` | Beta | mongo / pg / valkey | ServiceMonitor 3-repo 동등성 e2e |
| `pkg/networkpolicy` | Beta | mongo / pg / valkey | 4-direction (ingress/egress × TCP/UDP) 검증 |
| `pkg/security` | Beta | mongo / pg / valkey | restricted PSA 회귀 가드 3-repo |
| `pkg/webhook` | **Experimental** | 1 repo | 다중 repo 사용 후 안정화 |

**Tier 의미**:
- **Stable** — semver patch/minor 범위 BREAKING CHANGE 금지. deprecated 표기 + 2 minor 유예 후 제거.
- **Beta** — minor 범위 BREAKING CHANGE 가능 (CHANGELOG 명시). API 형태 거의 확정.
- **Experimental** — patch 범위에서도 변경 가능. 호출자 위험 부담.

## v1.0.0 졸업 조건 (체크리스트)

- [ ] 모든 패키지 **Stable** tier 도달
- [ ] BREAKING CHANGE 0건 / 연속 minor 릴리스 6 회 이상
- [ ] godoc coverage ≥ 80% (`go doc -all ./...` 기준)
- [ ] CHANGELOG.md 의 v0.x 진화 history 정리 + v1.0.0 release notes
- [ ] CITATION.cff + Zenodo DOI 발급 (학술 인용 가능)
- [ ] 3 repo (mongodb / postgres / valkey) 모두 v1.0.0 commons import 검증
- [x] `go vet ./... && go test ./...` clean (커버리지 96.3% > 85% threshold, 2026-05-14 측정)
- [x] API 안정성 promise 문서 — `docs/STABILITY.md` 신규 (PR #12)
- Verify: 3 repo CI 가 `operator-commons v1.0.0` import 후 모든 e2e PASS

## 패키지별 보강 항목

### pkg/finalizer (Stable)
- [x] `Add`, `Remove`, `Contains` helper — `pkg/finalizer/finalizer.go`
- [x] controller-runtime 회피 (std `slices` 사용)
- [x] unit test — `pkg/finalizer/finalizer_test.go`
- [x] 다중 finalizer 순서 보장 helper — `pkg/finalizer/order.go` `EnsureOrder` (PR #14)
- Verify: 3 repo finalizer 동작 회귀 0

### pkg/labels (Stable)
- [x] Kubernetes 권장 라벨 helper (app.kubernetes.io/*) — `pkg/labels/labels.go`
- [x] component / instance / part-of 매핑
- [x] unit test — `pkg/labels/labels_test.go`
- [x] Recommended labels v2 매핑 (K8s 1.30+) — `pkg/labels/v2.go` `AllV2` + `V2` struct (PR #14)
- Verify: 3 repo `metadata.labels` 일관성 검증

### pkg/status (Stable)
- [x] Condition 카탈로그 helper — `pkg/status/conditions.go`
- [x] `SetAvailable` 헬퍼 (v0.6.0)
- [x] unit test
- [x] Condition reason 표준 카탈로그 문서화 — `pkg/status/REASONS.md` (PR #13)
- Verify: `kubectl get <kind> -o yaml` 의 `.status.conditions` 동등성 (3 repo)

### pkg/version (Stable — promoted v1.0.0-rc, PR #19)
- [x] `Matrix[E]` generic 도입 (v0.7.0) — `pkg/version/matrix.go`
- [x] `SetAvailable` 헬퍼 (v0.6.0)
- [x] 버전 호환성 비교 (semver) — `pkg/version/version.go`
- [x] **Cross-version compatibility test** — `pkg/version/api_stability_test.go` (PR #15)
- [x] 버전 매트릭스 시리얼라이저 (`json`/`yaml`) — `pkg/version/serializer.go` `AsMap` + `MarshalJSON` (PR #15)
- [ ] **Tier 격상** → Stable
- Verify: 3 repo 의 version validation 동등 (mongodb / valkey / postgres 각각의 호환성 테이블)

### pkg/monitoring (Beta)
- [x] Prometheus ServiceMonitor 빌더 — `pkg/monitoring/monitoring.go`
- [x] unit test
- [x] PrometheusRule 빌더 (alert/recording rules 공통화) — `pkg/monitoring/rule.go` `NewPrometheusRule` + `AlertRule` + `RecordingRule` + `RuleGroup` (PR #18)
- [ ] **3-repo 동등성 e2e** — 같은 입력 → 같은 ServiceMonitor 출력
- [x] OpenTelemetry exporter helper — pkg/monitoring/otel.go (P-B.7.3)
- [ ] **Tier 격상** → Stable
- Verify: `monitoring_test.go` golden file diff 0 + 3 repo manifest 비교

### pkg/networkpolicy (Stable — promoted: B.8.3)
- [x] NetworkPolicy 빌더 — `pkg/networkpolicy/networkpolicy.go`
- [x] default-deny + 명시 규칙 helper
- [x] unit test
- [x] **4-direction 검증** — ingress/egress × TCP/UDP — `pkg/networkpolicy/four_dir_test.go` (PR #19)
- [x] CIDR + namespace selector + pod selector 조합 helper — `pkg/networkpolicy/combo.go` `ComboPeer` + `WithComboIngressFromPeers` (PR #16)
- [ ] **Tier 격상** → Stable
- Verify: kind 환경에서 NetworkPolicy 적용 후 차단/허용 경로 측정

### pkg/security (Stable — promoted: B.9.4)
- [x] SecurityContext helper (restricted PSA 호환) — `pkg/security/security.go`
- [x] RBAC helper
- [x] unit test
- [x] **restricted PSA 3-repo 회귀 가드** — `pkg/security/psa_guard_test.go` (PR #19)
- [x] Pod / Container SecurityContext 분리 helper — `pkg/security/split.go` `RestrictedPodSecurityContext` (PR #16)
- [x] seccompProfile 기본값 helper — `pkg/security/seccomp.go` `RuntimeDefaultSeccompProfile` + `LocalhostSeccompProfile` + `UnconfinedSeccompProfile` (PR #16)
- [ ] **Tier 격상** → Stable
- Verify: `kubectl label ns <ns> pod-security.kubernetes.io/enforce=restricted` 후 3 repo pod ready

### pkg/webhook (Beta — promoted: B.10.4)
- [x] Webhook 유틸 기초 — `pkg/webhook/webhook.go`
- [x] unit test
- [x] **Conversion webhook helper** — v1alpha1 ↔ v1alpha2 패턴 추출 — `pkg/webhook/conversion.go` `ConversionRegistry` (PR #18)
- [x] Validation webhook 공통 패턴 — pkg/webhook/validation_patterns.go (P-B.10.2)
- [ ] **다중 repo 사용** — 현재 1 repo (valkey) 만, 다른 repo 도입 후 안정화
- [ ] **Tier 격상** → Beta → Stable
- Verify: 2+ repo 가 동일 helper 사용 + 회귀 0

## 의존성 정책

- **K8s API 만** — `k8s.io/api`, `k8s.io/apimachinery`, `k8s.io/utils`. controller-runtime 의존 *추가 금지*.
- **Apache-2.0 호환 라이선스만** — 의존성 추가 시 ADR 작성.
- **godoc 완비** — 신규 public API 는 godoc 의무.

## 거버넌스 / 추적

- **CHANGELOG.md** — git-cliff 자동 생성. semantic versioning 엄수.
- **CITATION.cff** — 학술 인용 가능. DOI 는 v1.0.0 시점 발급.
- **ADOPTERS.md** — 사용 repo 기록 (현 3 repo).
- **ADR** — `docs/kb/adr/` 의 ADR-0002, ADR-0003, ADR-0004 가 설계 결정 추적.
- **AGENTS.md** — AI 협업 가이드.

## Non-Goals (의식적 비대상)

- ❌ **controller-runtime 의존 추가** — 현재 회피 설계 유지 (라이브러리 사용처 자유도 보장).
- ❌ **특정 operator 로직 흡수** — operator-specific 로직은 호출자 repo 에 둔다. 라이브러리는 *공통 헬퍼만*.
- ❌ **분기/날짜 기반 로드맵** — 글로벌 §workflow.md.
- ❌ **GitHub Actions 필수 release gate** — RFC 0002 글로벌. 로컬 4 계층 위임.
- ❌ **Plugin / extension SDK 포지셔닝** — 라이브러리이지 framework 아님.
- ❌ **v1.0.0 조기 선언** — 졸업 조건 미충족 시 v0.x 유지.

## 호출자 (Adopters)

현재 3 repo 가 import:

| Repo | 사용 패키지 | 비고 |
|---|---|---|
| `mongodb-operator` | finalizer / labels / status / version / monitoring / networkpolicy / security | v0.6.0+ |
| `postgres-operator` | finalizer / labels / status / version / monitoring / networkpolicy / security | v0.6.0+ |
| `valkey-operator` | finalizer / labels / status / version / monitoring / networkpolicy / security / webhook | v0.6.0+ (webhook 단독) |

## 변경 이력

| Date | Change | Refs |
|---|---|---|
| 2026-05-11 | ROADMAP.md 신설 — API stability tier + v1.0.0 졸업 조건 + 패키지별 sub-task 체크리스트 | parallel-leaping-seal plan |
| 2026-05-09 | v0.7.0 — generic `Matrix[E]` 추가 | CHANGELOG |
| 2026-05-09 | v0.6.0 — `SetAvailable` 헬퍼 추가 | CHANGELOG |
| 2026-05-09 | 4-repo audit 기반 설계 정합 | audit 산출물 |
