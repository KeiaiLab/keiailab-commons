# ARCHITECTURE — operator-commons

> Single-page architecture description for the `operator-commons` shared
> library. Updated when the package surface, tier, or design invariants
> change.

## Overview

- **Purpose**: a shared Go library that eliminates scaffolding drift across
  Kubernetes operator implementations — PodSecurity restricted contexts,
  version allowlists, NetworkPolicy templates, ServiceMonitor builders,
  finalizer / status helpers, PVC topology helpers.
- **Scope**: pure-Go helper packages — no CRDs, no controller-runtime hard
  dependency at the leaf-package level.
- **Stability tier**: v0.x (API may break). v1.0+ SemVer-stable per
  [ROADMAP.md](ROADMAP.md) graduation criteria.
- **License**: Apache-2.0.
- **Module path**: `github.com/keiailab/operator-commons`.

## Package surface

| Package | Tier | Purpose | controller-runtime dependency |
|---|---|---|---|
| `pkg/finalizer` | **Stable** | `Add` / `Remove` / `Has` / `EnsureOrder` (stdlib `slices` only). | No |
| `pkg/labels` | **Stable** | Recommended Kubernetes labels (`app.kubernetes.io/*`) builder — `Set`, `All()`, `Selector()`, plus `AllV2` for K8s 1.30+. | No |
| `pkg/status` | **Stable** | Four standard Condition Types + six Reason catalog + helpers (`SetReady`, `SetAvailable`, `SetReadyFalse`). | No |
| `pkg/storageclass` | **Stable** | DNS-1123 storageClass validator + `Normalize` / `MustNormalize`. | No |
| `pkg/version` | Beta | Version allowlist convention + generic `Matrix[E MatrixEntry]` + serializer. | No |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` and `PrometheusRule` builders (unstructured — CRD-soft). | No |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy builder + functional options + `ComboPeer`. | No |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext builder + Pod / Container split + seccomp profile pointers. | No |
| `pkg/events` | Beta | Minimal `Recorder` interface + nine standard `Reason` constants + `Emit` / `EmitWarning` / `WrappedError` (nil-safe). | No |
| `pkg/pvc` | Beta | PVC expansion helpers — comparison + safe in-place update. | Yes (a single package; see ADR-0016). |
| `pkg/topology` | Beta | PVC topology spread helpers + zone-aware affinity. | No |
| `pkg/probes` | Experimental | `corev1.Probe` fluent builder — HTTP / HTTPS / TCP / Exec, kubelet defaults + clamp. | No |
| `pkg/webhook` | Experimental | Admission validation helpers — `ValidateAllowedVersion`, `ValidateWithPredicate`, conversion registry. | No |
| `pkg/bundle` | **Experimental** | OLM v1 bundle metadata — annotations, FBC schema types, directory validation. | No |

Design invariant: **leaf packages depend on the Kubernetes API types and
stdlib only**. No controller-runtime, no logr, no operator-sdk leakage.
This keeps the library usable by any operator framework, including plain
`client-go` projects. `pkg/pvc` is the documented single exception
(ADR-0016).

## No CRDs, no reconciler

`operator-commons` is a **library**, not a controller. It deliberately
does not provide:

- CRD definitions.
- A reconciler or manager.
- RBAC manifests.
- An admission webhook server (only validation helpers).
- An ArgoCD application or similar deployment artifact.

Downstream consumers own those concerns. `operator-commons` provides the
building blocks they assemble.

## ADR index

Architecture decisions are tracked in [`docs/kb/adr/`](kb/adr/INDEX.md).
The index covers the project charter, tooling unification, the `pkg/status`
and `pkg/finalizer` sugar adoption, the generic `Matrix[E]` decision, the
Helm library chart adoption, the lefthook consolidation, the GitHub
Actions block hook, the manual release script, and the PVC / topology
extraction.

## Build and test surface

- **Go**: 1.26+ (per `go.mod`).
- **Test command**: `go test ./...`.
- **Lint command**: `golangci-lint run` (config `.golangci.yml`).
- **Custom linter**: `.custom-gcl.yml` for project-specific rules.
- **Pre-commit / pre-push**: `lefthook.yml` (DCO + Conventional Commits +
  gofmt + go vet + go test + govulncheck + go-mod-tidy drift).
- **Coverage target**: ≥ 85 % per package (v1.0 graduation criterion).
- **Renovate**: `renovate.json` for dependency updates.

## Helm chart partials

`charts/keiailab-commons/` ships a Helm *library chart* with partial
named templates for cross-repo Helm template reuse:

- `templates/_security.tpl` — restricted PSA partial.
- `templates/_networkpolicy.tpl` — deny-by-default partial.
- `templates/_rbac.tpl` — ServiceAccount + workload-base partial.
- `templates/_helpers.tpl` — shared label / naming helpers.

Consumers add the chart as a `helm dependency` and `include` the partials
in their own chart templates. The observability ServiceMonitor partial
lives in `templates/observability/_servicemonitor.tpl`.

## v1.0.0 graduation roadmap

Per [ROADMAP.md](ROADMAP.md) "v1.0.0 졸업 조건":

1. All Stable-candidate packages reach Stable tier.
2. Zero BREAKING CHANGE across six or more consecutive minor releases.
3. godoc coverage ≥ 80 %.
4. CITATION.cff + DOI.
5. Downstream import end-to-end verification on `v1.0.0-rc.N`.
6. `go vet && go test ./...` clean with ≥ 85 % coverage.
7. [docs/STABILITY.md](STABILITY.md) formal API stability promise.
8. CHANGELOG.md v0.x evolution history + v1.0.0 release notes.

## Non-goals

- CRD definitions (downstream consumers own them).
- Reconciler runtime (downstream consumers own it).
- Operator-framework abstraction (the library serves any operator-sdk /
  kubebuilder / homegrown reconciler).
- Controller-runtime hard dependency at the leaf-package level.
- Kubernetes versions older than the current supported set in
  `pkg/version`.

## References

- [README.md](../README.md) — installation, usage, package summary, badges.
- [ROADMAP.md](ROADMAP.md) — checklist + tier promotion criteria.
- [STABILITY.md](STABILITY.md) — API stability promise.
- [UPGRADING.md](UPGRADING.md) — version-bump migration notes.
- [CHANGELOG.md](../CHANGELOG.md) — versioned history.
- [CONTRIBUTING.md](../CONTRIBUTING.md) — DCO, lefthook, branch naming.
- [GOVERNANCE.md](GOVERNANCE.md) — Lazy Consensus + ADR process.
- [MAINTAINERS.md](MAINTAINERS.md) — maintainer roster.
- [SECURITY.md](../SECURITY.md) — vulnerability disclosure.
- [AGENTS.md](AGENTS.md) — AI-assistant runbook.
- [CITATION.cff](../CITATION.cff) — academic citation.
- [docs/kb/adr/](kb/adr/INDEX.md) — architecture decision records.
- [docs/kb/deps/](kb/deps/) — dependency audit logs.

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
