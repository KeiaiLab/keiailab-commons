# ROADMAP — operator-commons

> [English](ROADMAP.md) | **한국어**

본 ROADMAP 은 라이브러리의 진화 방향을 *API stability tier* + *v1.0.0 졸업
조건* + *패키지별 보강 항목* 으로 추적합니다. 분기 / 날짜 기반 마감일은 두지
않으며 (시간 기반 로드맵 금지), 라이브러리 특성상 *downstream consumer 의
요구* 가 진화 방향을 결정합니다.

## 체크박스 의미

| 마커 | 의미 |
|---|---|
| `[x]` | 코드 + 테스트 양쪽 존재. downstream import 가능. |
| `[~]` | 부분 구현 (helper 존재, 검증 미완). |
| `[ ]` | 미시작. |

## API Stability Tier (현행 v0.9.x candidate)

| 패키지 | Tier | Tier 격상 조건 |
|---|---|---|
| `pkg/finalizer` | **Stable** | v1 진입 (별도 작업 없음). |
| `pkg/labels` | **Stable** | v1 진입 (별도 작업 없음). |
| `pkg/status` | **Stable** | v1 진입 (별도 작업 없음). |
| `pkg/storageclass` | **Stable** | trivial validation surface (regex + nil check) — 즉시 Stable. |
| `pkg/version` (incl. `Matrix`) | Beta | generic `Matrix[E]` 의 cross-repo 검증 완료. |
| `pkg/monitoring` | Beta | `ServiceMonitor` 의 downstream 동등성 e2e. |
| `pkg/networkpolicy` | Beta | 4-direction (ingress / egress × TCP / UDP) 검증. |
| `pkg/security` | Beta | restricted PSA 회귀 가드 downstream. |
| `pkg/events` | Beta | downstream 라이브 적용 + reconciliation 회귀 0. |
| `pkg/pvc` | Beta | PVC expansion downstream 라이브 적용. |
| `pkg/topology` | Beta | topology spread downstream 라이브 적용. |
| `pkg/webhook` | **Experimental** | 다중 downstream 사용 후 안정화. |
| `pkg/probes` | **Experimental** | 2+ downstream 라이브 적용 후 Beta. |

**Tier 의미**:

- **Stable** — semver patch / minor 범위 BREAKING CHANGE 금지. deprecated
  표기 + 2 minor 유예 후 제거.
- **Beta** — minor 범위 BREAKING CHANGE 가능 (CHANGELOG 명시). API 형태 거의
  확정.
- **Experimental** — patch 범위에서도 변경 가능. 호출자가 위험을 부담합니다.

## v1.0.0 졸업 조건 (체크리스트)

- [ ] 모든 패키지 **Stable** tier 도달.
- [ ] BREAKING CHANGE 0 건 / 연속 minor 릴리스 6 회 이상.
- [ ] godoc coverage ≥ 80 % (`go doc -all ./...` 기준).
- [ ] CHANGELOG.md 의 v0.x 진화 history 정리 + v1.0.0 release notes.
- [ ] CITATION.cff + Zenodo DOI 발급 (학술 인용 가능).
- [ ] downstream consumer 가 v1.0.0 commons import 후 회귀 0.
- [x] `go vet ./... && go test ./...` clean (커버리지 96.3 % > 85 % threshold).
- [x] API 안정성 promise 문서 — [STABILITY.md](STABILITY.md).
- Verify: downstream consumer CI 가 `operator-commons v1.0.0` import 후 모든
  e2e PASS.

## 패키지별 보강 항목

### pkg/finalizer (Stable)

- [x] `Add`, `Remove`, `Contains` helper — `pkg/finalizer/finalizer.go`.
- [x] controller-runtime 회피 (stdlib `slices` 사용).
- [x] unit test — `pkg/finalizer/finalizer_test.go`.
- [x] 다중 finalizer 순서 보장 helper — `pkg/finalizer/order.go` `EnsureOrder`.
- Verify: downstream finalizer 동작 회귀 0.

### pkg/labels (Stable)

- [x] Kubernetes 권장 라벨 helper (`app.kubernetes.io/*`) — `pkg/labels/labels.go`.
- [x] component / instance / part-of 매핑.
- [x] unit test — `pkg/labels/labels_test.go`.
- [x] Recommended labels v2 매핑 (K8s 1.30+) — `pkg/labels/v2.go` `AllV2`.
- Verify: downstream `metadata.labels` 일관성.

### pkg/status (Stable)

- [x] Condition 카탈로그 helper — `pkg/status/conditions.go`.
- [x] `SetAvailable` 헬퍼.
- [x] unit test.
- [x] Condition reason 표준 카탈로그 문서화 — `pkg/status/REASONS.md`.
- Verify: `kubectl get <kind> -o yaml` 의 `.status.conditions` downstream 동등성.

### pkg/version (Beta)

- [x] `Matrix[E]` generic 도입 — `pkg/version/matrix.go`.
- [x] `SetAvailable` 헬퍼.
- [x] 버전 호환성 비교 (semver) — `pkg/version/version.go`.
- [x] cross-version compatibility test — `pkg/version/api_stability_test.go`.
- [x] 버전 매트릭스 시리얼라이저 (`json` / `yaml`) — `pkg/version/serializer.go`.
- [ ] **Tier 격상** → Stable.
- Verify: downstream version validation 동등성.

### pkg/monitoring (Beta)

- [x] Prometheus ServiceMonitor 빌더 — `pkg/monitoring/monitoring.go`.
- [x] unit test.
- [x] PrometheusRule 빌더 (alert / recording 규칙 공통화) — `pkg/monitoring/rule.go`.
- [x] OpenTelemetry exporter helper — `pkg/monitoring/otel.go`.
- [ ] downstream 동등성 e2e — 같은 입력 → 같은 manifest 출력.
- [ ] **Tier 격상** → Stable.
- Verify: `monitoring_test.go` golden file diff 0.

### pkg/networkpolicy (Beta)

- [x] NetworkPolicy 빌더 — `pkg/networkpolicy/networkpolicy.go`.
- [x] default-deny + 명시 규칙 helper.
- [x] unit test.
- [x] 4-direction 검증 — `pkg/networkpolicy/four_dir_test.go`.
- [x] CIDR + namespace selector + pod selector 조합 helper — `pkg/networkpolicy/combo.go`.
- [ ] **Tier 격상** → Stable.
- Verify: kind 환경에서 NetworkPolicy 적용 후 차단 / 허용 경로 측정.

### pkg/security (Beta)

- [x] SecurityContext helper (restricted PSA 호환) — `pkg/security/security.go`.
- [x] RBAC helper.
- [x] unit test.
- [x] restricted PSA 회귀 가드 — `pkg/security/psa_guard_test.go`.
- [x] Pod / Container SecurityContext 분리 helper — `pkg/security/split.go`.
- [x] seccompProfile 기본값 helper — `pkg/security/seccomp.go`.
- [ ] **Tier 격상** → Stable.
- Verify: `kubectl label ns <ns> pod-security.kubernetes.io/enforce=restricted`
  적용 후 downstream pod ready.

### pkg/webhook (Experimental)

- [x] Webhook 유틸 기초 — `pkg/webhook/webhook.go`.
- [x] unit test.
- [x] Conversion webhook helper — `pkg/webhook/conversion.go`.
- [x] Validation webhook 공통 패턴 — `pkg/webhook/validation_patterns.go`.
- [ ] 다중 downstream 라이브 적용 후 안정화.
- [ ] **Tier 격상** → Beta → Stable.
- Verify: 2 이상의 downstream 이 동일 helper 사용 + 회귀 0.

### pkg/storageclass (Stable)

- [x] DNS-1123 subdomain validation — `pkg/storageclass/validator.go`.
- [x] Normalize / MustNormalize — empty → nil (cluster default) + trim +
  pointer return.
- [x] unit test 12 cases — `pkg/storageclass/validator_test.go` (100 % coverage).
- [ ] downstream 라이브 적용 + 회귀 0.

### pkg/events (Beta)

- [x] Recorder interface — client-go `record.EventRecorder` 구조 정합
  (client-go 의존 회피).
- [x] 9 Reason constants (Created / Updated / Deleted / Reconciled /
  ReconcileError / Provisioning / Ready / Degraded / Failed).
- [x] Emit / Emitf / EmitWarning / EmitWarningf / WrappedError — nil-safe.
- [x] unit test — `pkg/events/events_test.go` (100 % coverage).
- [ ] downstream 라이브 적용 후 reconciliation Event reason 통일.
- [ ] **Tier 격상** → Stable.
- Verify: downstream Reconcile path 의 Event reason 이 commons constants 사용 +
  회귀 0.

### pkg/probes (Experimental)

- [x] Builder fluent API — HTTP / HTTPS / TCP / Exec 4 handlers.
- [x] kubelet default (Period=10 s / Timeout=1 s / SuccessThreshold=1 /
  FailureThreshold=3).
- [x] InitialDelay / Period / Timeout 음수 → 0 clamp.
- [x] Build() handler 미설정 시 panic (fail-fast contract).
- [x] unit test — `pkg/probes/builder_test.go` (100 % coverage).
- [ ] 2+ downstream 라이브 적용 (Beta 격상 조건).
- [ ] **Tier 격상** → Beta → Stable.

### pkg/pvc (Beta)

- [x] PVC expansion helper — `pkg/pvc/expansion.go`.
- [x] unit test — `pkg/pvc/expansion_test.go`.
- [ ] downstream 라이브 적용 후 PVC resize 회귀 0.
- [ ] **Tier 격상** → Stable.

### pkg/topology (Beta)

- [x] PVC topology spread helper — `pkg/topology/spread.go`.
- [x] unit test — `pkg/topology/spread_test.go`.
- [ ] downstream 라이브 적용 후 spread constraint 검증.
- [ ] **Tier 격상** → Stable.

## 의존성 정책

- **Kubernetes API 만** — `k8s.io/api`, `k8s.io/apimachinery`, `k8s.io/utils`.
  controller-runtime 의존 *추가 금지*.
- **Apache-2.0 호환 라이선스만** — 의존성 추가 시 ADR 작성.
- **godoc 완비** — 신규 public API 는 godoc 의무.

## 거버넌스 / 추적

- **CHANGELOG.md** — git-cliff 자동 생성. semantic versioning 엄수.
- **CITATION.cff** — 학술 인용 가능. DOI 는 v1.0.0 시점 발급.
- **ADR** — `docs/kb/adr/` 가 설계 결정을 추적.
- **AGENTS.md** — AI 협업 가이드.

## Non-Goals (의식적 비대상)

- ❌ **controller-runtime 의존 추가** — leaf 패키지의 회피 설계를 유지.
- ❌ **downstream-specific 로직 흡수** — operator-specific 코드는 호출자
  repo 에 둡니다. 라이브러리는 *공통 헬퍼만*.
- ❌ **분기 / 날짜 기반 로드맵** — 기능 체크리스트 + 진행률.
- ❌ **GitHub Actions 필수 release gate** — 로컬 4 계층 위임.
- ❌ **Plugin / extension SDK 포지셔닝** — 라이브러리이지 framework 가
  아닙니다.
- ❌ **v1.0.0 조기 선언** — 졸업 조건 미충족 시 v0.x 유지.

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
