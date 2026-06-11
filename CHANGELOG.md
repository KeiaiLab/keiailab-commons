# Changelog

> **English** | [한국어](CHANGELOG.ko.md) | [日本語](CHANGELOG.ja.md) | [中文](CHANGELOG.zh.md)

All notable changes to this library are recorded in this file.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Versioning: [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Automated generation: `git-cliff` — the changelog is regenerated as a PR
at release-tag time.

## [Unreleased]

### Added

- `pkg/apply` (Beta tier) — idempotent apply helpers for ConfigMap /
  Secret / Service / StatefulSet / Deployment / NetworkPolicy /
  PodDisruptionBudget / HorizontalPodAutoscaler. Service ClusterIP /
  IPFamilies create-only guard, StatefulSet immutable-field preservation
  plus `RetryOnConflict`, Deployment server-default and
  revision-annotation preservation, `preserveReplicas` option (avoids
  fighting the HPA). Resolves ~530 LOC duplicated across the three
  operators. Depends on controller-runtime.
- `pkg/reconcile` (Beta tier) — `Statusable` interface (`client.Object`
  + `GetConditions` + `SetPhase`) plus `ApplyErrorCondition` +
  `HandleFinalizerCleanup` + `SecretIfNotExists`. Depends on
  controller-runtime.
- `pkg/certmanager` (Beta tier) — `CertParams` + `BuildCertificate` +
  `BuildSelfSignedIssuer` + `ServiceSANs`. Built on unstructured
  objects; zero Go dependency on cert-manager CRDs.
- `pkg/reconcilemetrics` (Beta tier) — `ReconcileMetrics` (Total /
  Latency / Errors) + `New(subsystem)` + `MustRegister` + `IncTotal` /
  `ObserveReconcile` / `IncError` / `DeleteFor` + `ResultFor`. The
  injected subsystem keeps each operator's existing Prometheus
  time-series names intact.
- `pkg/status` `UpdateWithRetry` (Beta surface inside a Stable package)
  — refetch + mutate + `RetryOnConflict` persistence of the status
  subresource.

### Changed

- New direct dependency: `github.com/prometheus/client_golang` v1.23.2
  (introduced by `pkg/reconcilemetrics`).

## [0.10.0] — 2026-06-11

### Added

- `pkg/bundle` (Experimental tier) — OLM v1 operator bundle metadata
  helpers: `Annotations` (+ `DockerLabels`), FBC `Package` / `Channel` /
  `Bundle` builders, `ValidateDir`.
- Helm chart partial named template `keiailab.secrets.externalSecret`
  (`charts/keiailab-commons/templates/_externalsecret.tpl`).
- GitLab CI shadow pipeline.

### Changed

- **BREAKING**: module path renamed —
  `github.com/keiailab/operator-commons` →
  `github.com/keiailab/keiailab-commons`. Every consumer must update its
  import paths; shipped as a minor bump as permitted during the 0.x
  phase.
- License standardized to MIT — SPDX headers added to every `.go` file.
- README rewritten for accuracy.
- Self-contained restructuring — external service references removed
  (B1–B14).

## [0.9.0] — 2026-05-21

### Added

- `pkg/pvc` (Beta tier) — PVC expansion helper plus a safe in-place
  update path.
- `pkg/topology` (Beta tier) — PVC topology spread helper plus zone-aware
  affinity.
- `scripts/release.sh` — manual release pipeline for the library
  (ADR-0014).
- `docs/UPGRADING.md` — semver policy plus the three-operator migration
  guide.
- i18n S4 Phase 1–5 — 4-language glossary completed plus the translation
  sync hook.

### Changed

- keiailab branding Wave 3 — README header / footer plus `BRANDING.md`
  and `docs/family.md`.

## [0.8.0] — 2026-05-21

### Added

- `pkg/probes` (Experimental tier) — `corev1.Probe` fluent builder. HTTP /
  HTTPS / TCP / Exec handlers + kubelet defaults (Period = 10 s /
  Timeout = 1 s / SuccessThreshold = 1 / FailureThreshold = 3) +
  InitialDelay / Period / Timeout negative-clamp to zero + `Build()`
  panics when no handler is set (fail-fast). 100 % coverage, zero lint.
- `pkg/storageclass` (Stable tier immediately) — DNS-1123 subdomain
  validator. `IsValid` / `Validate` (+ `ErrInvalidStorageClassName`
  sentinel) / `Normalize` (empty → nil cluster default + trim + pointer
  return) / `MustNormalize`. 100 % coverage, zero lint.
- `pkg/events` (Beta tier) — Kubernetes Event recorder helper plus nine
  standard `Reason` constants (Created / Updated / Deleted / Reconciled /
  ReconcileError / Provisioning / Ready / Degraded / Failed). Minimal
  `Recorder` interface (compatible with `client-go` `record.EventRecorder`
  without importing it). `Emit` / `Emitf` / `EmitWarning` /
  `EmitWarningf` (nil-safe) + `WrappedError`. 100 % coverage, zero lint.
- `pkg/monitoring.NewPrometheusRule` + `AlertRule` / `RecordingRule` /
  `RuleGroup` — PrometheusRule (`monitoring.coreos.com/v1`) manifest
  builder.
- `pkg/webhook.ConversionRegistry` — CRD version-pair conversion
  function registry (`Register` / `Convert` / `HasPair`).
- `pkg/networkpolicy.ComboPeer` + `WithComboIngressFromPeers` — CIDR +
  NamespaceSelector + PodSelector composite peer helper.
- `pkg/security.RestrictedPodSecurityContext` + options
  (`WithPodFSGroup`, `WithPodRunAsUser`, `WithPodRunAsGroup`) — Pod-level
  restricted SecurityContext.
- `pkg/security.RuntimeDefaultSeccompProfile` +
  `LocalhostSeccompProfile` + `UnconfinedSeccompProfile` — seccomp
  profile pointer helpers.
- `pkg/version.AsMap` + `MarshalJSON` — `Matrix[E]` serializer with
  stable, JSON / YAML-compatible key ordering.
- `pkg/version/api_stability_test.go` — public-API-surface guard.
- `pkg/finalizer.EnsureOrder` — ordering guarantee helper for multiple
  finalizers; stable sort against `desiredOrder`, finalizers not listed
  are kept at the tail.
- `pkg/labels.AllV2` + `V2` struct — Kubernetes 1.30+ Recommended labels
  v2 mapping.
- `pkg/status/REASONS.md` — Reason × Type × Status usage matrix.
- `docs/STABILITY.md` — three-tier API stability promise plus graduation
  criteria and breaking-change policy.
- `scripts/godoc-coverage.sh` — per-package + total godoc coverage
  measurement. Verifies the 80 % threshold required for v1.0.
- `docs/ARCHITECTURE.md` — single-page architecture description.
- README 4-language i18n started — English canonical plus a Korean
  translation, with a 4-language switcher and Japanese / Chinese
  placeholders.

## [0.7.0] — 2026-05-09

### Added

- `pkg/version`: generic `Matrix[E MatrixEntry]` — caller-supplied entry
  types can be delegated to the library.
- `docs/kb/adr/0004-*` — ADR documenting the `Matrix` generic decision.

## [0.6.0] — 2026-05-09

### Added

- `pkg/status`: `SetAvailable` + `SetReadyFalse` sugar helpers.
- `docs/kb/adr/0003-*` — ADR documenting the status sugar helper
  decision.
- `.codecov.yml` — coverage floor previously unified across consumers.

## [0.5.0] — 2026-05-09

### Added

- Governance documents: AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY /
  MAINTAINERS / CODE_OF_CONDUCT.
- `pkg/status/`: four standard Condition Types + six-Reason catalog +
  helpers. External dependency: `k8s.io/apimachinery` only.
- `pkg/finalizer/`: `Add` / `Remove` / `Has` helpers. controller-runtime
  is avoided; the stdlib `slices` package is the only dependency.
- `templates/observability/_servicemonitor.tpl`: Helm chart partial
  named template `keiailab.observability.serviceMonitor`.
- `templates/observability/README.md`: metric naming convention plus
  shared alert recommendations plus consumer chart usage.
- `Makefile` (lint / test / audit / cover / tidy / tag).
- `.golangci.yml` + `.custom-gcl.yml`.
- `CHANGELOG.md` + `docs/kb/deps/2026-05.md` dependency audit log.
- `docs/kb/adr/0002-tooling-unification-adoption.md` +
  `docs/kb/adr/INDEX.md`.
- `NOTICE` (Apache-2.0 §4(d) compliance).
- `CODEOWNERS`.
- README badges — License / Go / pkg.go.dev / OpenSSF Scorecard /
  Discussions plus the Community section.
- `renovate.json`.
- `lefthook.yml` (library minimal).
- DCO Signed-off-by warn-only commit-msg gate.

### Changed

- ADR directory moved: `docs/adr/` → `docs/kb/adr/`.
- `go` directive `1.25.0` → `1.25.7`.
- `pkg/finalizer` lint fix: manual `for` + `==` → `slices.Contains` /
  `slices.Index` / `slices.Delete` (modernize linter).
- `pkg/status/conditions.go` `SetReady` signature is now multi-line
  (passes `lll`).

## [0.4.0] — 2026-05-07

Earlier history is tracked via git tags and release notes; this
`CHANGELOG.md` was created during the audit cycle.

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
