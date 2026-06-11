# Changelog

> [English](CHANGELOG.md) | **한국어** | [日本語](CHANGELOG.ja.md) | [中文](CHANGELOG.zh.md)

본 라이브러리의 모든 주요 변경은 본 파일에 기록한다.
형식: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
버저닝: [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

자동 생성: `git-cliff` — release tag 시점에 PR 자동 갱신 후속.

## [Unreleased]

### Added

- `pkg/apply` (Beta tier) — ConfigMap / Secret / Service / StatefulSet /
  Deployment / NetworkPolicy / PodDisruptionBudget /
  HorizontalPodAutoscaler idempotent apply helper. Service ClusterIP /
  IPFamilies create-only 가드 + StatefulSet immutable 필드 보존 +
  `RetryOnConflict` + Deployment server-default 보존 + revision
  annotation 보존 + `preserveReplicas` 옵션 (HPA 충돌 방지). 3-operator
  중복 ~530 LOC 해소. controller-runtime 의존.
- `pkg/reconcile` (Beta tier) — `Statusable` interface (`client.Object` +
  `GetConditions` + `SetPhase`) + `ApplyErrorCondition` +
  `HandleFinalizerCleanup` + `SecretIfNotExists`. controller-runtime 의존.
- `pkg/certmanager` (Beta tier) — `CertParams` + `BuildCertificate` +
  `BuildSelfSignedIssuer` + `ServiceSANs` — unstructured 기반, cert-manager
  CRD Go 의존 0.
- `pkg/reconcilemetrics` (Beta tier) — `ReconcileMetrics` (Total / Latency /
  Errors) + `New(subsystem)` + `MustRegister` + `IncTotal` /
  `ObserveReconcile` / `IncError` / `DeleteFor` + `ResultFor` — subsystem
  주입으로 기존 operator Prometheus 시계열 이름 보존.
- `pkg/status` `UpdateWithRetry` (Stable 패키지 내 Beta 표면) — refetch +
  mutate + `RetryOnConflict` 로 status subresource 영속화.

### Changed

- 신규 직접 의존: `github.com/prometheus/client_golang` v1.23.2
  (`pkg/reconcilemetrics` 가 도입).

## [0.10.0] - 2026-06-11

### Added

- `pkg/bundle` (Experimental tier) — OLM v1 operator bundle metadata
  helper: `Annotations` (+ `DockerLabels`) + FBC `Package` / `Channel` /
  `Bundle` builder + `ValidateDir`.
- helm chart partial named template `keiailab.secrets.externalSecret`
  신규 (`charts/keiailab-commons/templates/_externalsecret.tpl`).
- GitLab CI shadow pipeline 신규.

### Changed

- **BREAKING**: module path 변경 —
  `github.com/keiailab/operator-commons` →
  `github.com/keiailab/keiailab-commons`. 모든 consumer 의 import path
  갱신 필요 — 0.x 단계라 minor bump 으로 진행.
- 라이선스 MIT 표준화 — 전 `.go` 파일 SPDX 헤더 추가.
- README 정확성 기준 재작성.
- self-contained 재구성 — 외부 서비스 참조 제거 (B1~B14).

## [0.9.0] - 2026-05-21

### Added

- `pkg/pvc` (Beta tier) — PVC expansion helper + 안전한 in-place update.
- `pkg/topology` (Beta tier) — PVC topology spread helper + zone-aware
  affinity.
- `scripts/release.sh` — 라이브러리 수동 release pipeline (ADR-0014).
- `docs/UPGRADING.md` — semver 정책 + 3 operator 마이그레이션 가이드.
- i18n S4 Phase 1~5 — glossary 4-lang 완성 + 번역 sync hook.

### Changed

- keiailab branding Wave 3 — README header / footer + `BRANDING.md` +
  `docs/family.md`.

## [0.8.0] - 2026-05-21

### Added

- `pkg/probes` (Experimental tier) — corev1.Probe fluent builder. HTTP /
  HTTPS / TCP / Exec 4 handlers + kubelet default (Period = 10 s /
  Timeout = 1 s / SuccessThreshold = 1 / FailureThreshold = 3) + InitialDelay
  / Period / Timeout 음수 → 0 clamp + Build() handler 미설정 시 panic
  (fail-fast). 100 % coverage + 0 lint.
- `pkg/storageclass` (Stable tier 즉시) — DNS-1123 subdomain validator.
  `IsValid` / `Validate` (+ `ErrInvalidStorageClassName` sentinel) /
  `Normalize` (empty → nil cluster default + trim + pointer return) /
  `MustNormalize`. 100 % coverage + 0 lint.
- `pkg/events` (Beta tier) — Kubernetes Event 생성 helper + 9 표준 Reason
  상수 (Created / Updated / Deleted / Reconciled / ReconcileError /
  Provisioning / Ready / Degraded / Failed) + minimal `Recorder` interface +
  Emit / Emitf / EmitWarning / EmitWarningf (nil-safe) + WrappedError.
  100 % coverage + 0 lint.
- `pkg/monitoring.NewPrometheusRule` + `AlertRule` / `RecordingRule` /
  `RuleGroup` — PrometheusRule (`monitoring.coreos.com/v1`) manifest
  builder.
- `pkg/webhook.ConversionRegistry` — CRD version pair 변환 함수 registry
  (`Register` / `Convert` / `HasPair`).
- `pkg/networkpolicy.ComboPeer` + `WithComboIngressFromPeers` — CIDR +
  NamespaceSelector + PodSelector 조합 helper.
- `pkg/security.RestrictedPodSecurityContext` + 옵션 (`WithPodFSGroup`,
  `WithPodRunAsUser`, `WithPodRunAsGroup`) — Pod-level restricted
  SecurityContext.
- `pkg/security.RuntimeDefaultSeccompProfile` + `LocalhostSeccompProfile` +
  `UnconfinedSeccompProfile` — seccomp profile pointer helpers.
- `pkg/version.AsMap` + `MarshalJSON` — `Matrix[E]` 시리얼라이저. JSON / YAML
  호환, key 정렬 stable output.
- `pkg/version/api_stability_test.go` — public API surface 가드.
- `pkg/finalizer.EnsureOrder` — 다중 finalizer 순서 보장 helper. desiredOrder
  안정 정렬 + 미지정 finalizer 후미 유지.
- `pkg/labels.AllV2` + `V2` struct — K8s 1.30+ Recommended labels v2 매핑.
- `pkg/status/REASONS.md` — Reason × Type × Status 사용 매트릭스.
- `docs/STABILITY.md` — 3-tier API stability promise + graduation criteria +
  breaking change policy.
- `scripts/godoc-coverage.sh` — per-package + total godoc coverage 계산.
  v1.0 80 % threshold 검증.
- `docs/ARCHITECTURE.md` — single-page 아키텍처 설명.
- README 4-lang i18n 시작 — EN canonical + KO 번역 + 4-lang switcher +
  ja / zh placeholder.

## [0.7.0] - 2026-05-09

### Added

- `pkg/version`: generic `Matrix[E MatrixEntry]` 추가 — caller-supplied entry
  type 위임 지원.
- `docs/kb/adr/0004-*` 신규 — `Matrix` generic 도입 결정 근거.

## [0.6.0] - 2026-05-09

### Added

- `pkg/status`: `SetAvailable` + `SetReadyFalse` 슈가 헬퍼 추가.
- `docs/kb/adr/0003-*` 신규 — status 슈가 헬퍼 결정 근거.
- `.codecov.yml` 신규 — 코드 커버리지 floor 통일.

## [0.5.0] - 2026-05-09

### Added

- 거버넌스 doc 신설 (AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY /
  MAINTAINERS / CODE_OF_CONDUCT).
- `pkg/status/`: 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼. 외부 의존:
  `k8s.io/apimachinery` 만.
- `pkg/finalizer/`: `Add` / `Remove` / `Has` 헬퍼. controller-runtime 의존
  회피, stdlib `slices` 만 사용.
- `templates/observability/_servicemonitor.tpl`: helm chart partial named
  template `keiailab.observability.serviceMonitor`.
- `templates/observability/README.md`: 메트릭 명명 규약 + 공통 alert 권장 +
  consumer chart 사용법.
- `Makefile` 신규 (lint / test / audit / cover / tidy / tag).
- `.golangci.yml` + `.custom-gcl.yml` 신규.
- `CHANGELOG.md` 신설 + `docs/kb/deps/2026-05.md` 의존성 audit log 신규.
- `docs/kb/adr/0002-tooling-unification-adoption.md` + `docs/kb/adr/INDEX.md`
  신규.
- `NOTICE` 신설 (Apache-2.0 §4(d) 정합).
- `CODEOWNERS` 신설.
- README badges 신설 — License / Go / pkg.go.dev / OpenSSF Scorecard /
  Discussions + Community 섹션.
- `renovate.json` 신설.
- `lefthook.yml` 신설 (라이브러리 minimal).
- DCO Signed-off-by warn-only commit-msg gate.

### Changed

- ADR 디렉토리 이전: `docs/adr/` → `docs/kb/adr/`.
- `go` directive `1.25.0` → `1.25.7`.
- `pkg/finalizer` lint fix: 수동 for + == → `slices.Contains` /
  `slices.Index` / `slices.Delete` (modernize linter).
- `pkg/status/conditions.go` `SetReady` 함수 시그니처 multi-line (lll 통과).

## [0.4.0] - 2026-05-07

(이전 버전 history 는 git tag log 또는 release notes 참조 — 본 CHANGELOG.md
는 audit 시점에 신설)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
