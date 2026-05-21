# ARCHITECTURE — operator-commons

> Single-page architecture description for the keiailab operator shared library. Updated when package surface / tier / adoption matrix changes.

## Overview

- **Purpose**: Shared Go library that eliminates scaffolding drift across `mongodb-operator`, `valkey-operator`, `postgres-operator` (PodSecurity restricted contexts, version allowlists, NetworkPolicy templates, ServiceMonitor builders, finalizer/status helpers).
- **Scope**: Pure-Go helper package — no CRDs of its own, no controller-runtime hard dependency at the leaf-package level.
- **Stability tier**: v0.x (API may break). v1.0+ SemVer-stable per `ROADMAP.md` graduation criteria.
- **Latest release**: v0.7.0 (2026-05-11)
- **License**: Apache-2.0
- **Module path**: `github.com/keiailab/operator-commons`

## Package surface (8 packages)

| Package | Tier | Purpose | Controller-runtime dep? |
|---|---|---|---|
| `pkg/finalizer` | **Stable** | `Add` / `Remove` / `Has` finalizer helpers using std `slices` only | No |
| `pkg/labels` | **Stable** | Recommended K8s labels builder — `Set`, `All()`, `Selector()` (version-aware split) | No |
| `pkg/status` | **Stable** | 4 standard Condition Types + 6 Reason catalog + helpers (`SetReady`, `SetAvailable`, `SetReadyFalse`) | No |
| `pkg/version` | Beta | DB version allowlist convention (`MustList`, `IsSupported`, `Strings`, `Default`) + generic `Matrix[E MatrixEntry]` | No |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` builder (unstructured — CRD-soft) | No |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy builder + functional options (`WithSelfIngress`, `WithIngressFromPeers`, `WithDenyEgress`, `WithEgressToPeers`) | No |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext builder with functional options | No |
| `pkg/webhook` | Experimental | Admission validation helpers — `ValidateAllowedVersion` (exact match), `ValidateWithPredicate` (caller-supplied matcher) | No |

Design invariant: **leaf packages are stdlib + k8s API types only**. No controller-runtime, no logr, no operator-sdk leakage. This keeps the library usable by any operator framework or even plain `client-go` projects.

## No CRDs, no reconciler

operator-commons is a **library**, not a controller. There is no:
- CRD definitions
- Reconciler / manager
- RBAC manifests
- Admission webhook server
- ArgoCD application

Callers (the 3 operators below) own those concerns. operator-commons provides building blocks they assemble.

## Adoption matrix

| Operator | sec | ver | lab | mon | np | wh | fin | sta | Tier % |
|---|---|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | ✅ | ✅ | 6/8 (75%) |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **8/8 (100%)** |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | ✅ | ✅ | 5/8 (63%) |

valkey-operator is the **carbon-copy reference** — first repo to reach 100% adoption. Deferred (`⏳`) cells correspond either to feature additions still in flight in the host operator (mongodb webhook server, ServiceMonitor reconciler) or to *separately abstracted* needs (postgres `version.Combo` struct is richer than `commons.MustList`, so delegation is not appropriate).

Live evidence of each promotion is in `README.md` "Adoption Matrix" section with per-iteration commit SHAs.

## ADR cross-link (8 ADRs)

| ADR | Title | Status |
|---|---|---|
| ADR-0001 | Charter | Accepted |
| ADR-0002 | RFC-0017 tooling unification adoption | Accepted |
| ADR-0003 | RFC-0018 `pkg/status` + `pkg/finalizer` adoption | Accepted |
| ADR-0004 | `pkg/version` generic `Matrix[E]` | Accepted |
| ADR-0005 | RFC-0019 library chart adoption | Accepted |
| ADR-0006 | RFC-0019 NetworkPolicy partials | Accepted |
| ADR-0007 | RFC-0019 RBAC partials | Accepted |
| ADR-0008 | RFC-0019 Security partials | Accepted |

All under `docs/kb/adr/`. `INDEX.md` is the navigation entry.

## Build / test surface

- **Go**: 1.25+ (per `go.mod`)
- **Test command**: `go test ./...`
- **Lint command**: `golangci-lint run` (config `.golangci.yml`)
- **Custom linter**: `.custom-gcl.yml` for project-specific rules
- **Pre-commit**: `lefthook.yml` (DCO + commitlint + golangci-lint)
- **Coverage target**: ≥85% per package (v1.0 graduation criterion)
- **Renovate**: `renovate.json` for dependency updates

## Helm chart partials

`templates/` and `charts/` directories ship helm *library chart* partials for cross-repo helm template reuse (RFC-0019):

- `templates/observability/_servicemonitor.tpl` — ServiceMonitor partial
- `templates/security/_pod_security_context.tpl` — restricted PSA partial
- `templates/networkpolicy/_default_deny.tpl` — deny-by-default partial

Callers `helm dependency` on this chart and `include` the partials.

## v1.0.0 graduation roadmap

Per `ROADMAP.md` §"v1.0.0 graduation criteria":

1. All 8 packages reach Stable tier (currently 3 Stable / 4 Beta / 1 Experimental)
2. 0 BREAKING CHANGE across 6+ consecutive minor releases
3. godoc coverage ≥80%
4. CITATION.cff + Zenodo DOI
5. 3-repo import e2e verification on `v1.0.0-rc`
6. `go vet && go test ./...` clean with ≥85% coverage
7. `docs/STABILITY.md` formal API stability promise
8. `CHANGELOG.md` v0.x evolution history + v1.0.0 release notes
9. `pkg/finalizer` multi-finalizer order guarantee
10. `pkg/labels` K8s 1.30+ v2 mapping
11. `pkg/status` Condition reason catalog docs

Tracking: `~/.claude/plans/2026-05-14-4-operators-100pct/P-B.md` (29 sub-tasks → 0).

## Non-goals

- ❌ CRD definitions (callers own)
- ❌ Reconciler runtime (callers own)
- ❌ Operator-framework abstraction (we serve any operator-sdk / kubebuilder / homegrown)
- ❌ Controller-runtime hard dependency at leaf-package level
- ❌ K8s version older than current supported set per `pkg/version`

## References

- `README.md` — installation, usage examples, package badges
- `ROADMAP.md` — checklist + tier promotion criteria
- `CHANGELOG.md` — versioned history
- `ADOPTERS.md` — external adopters list
- `CONTRIBUTING.md` — DCO, lefthook, branch naming
- `GOVERNANCE.md` — Lazy Consensus + RFC process
- `HANDOFF.md` — current session context
- `MAINTAINERS.md` — maintainer roster
- `SECURITY.md` — vulnerability disclosure policy
- `AGENTS.md` — AI-assistant runbook (Claude Code / Cursor)
- `CITATION.cff` — academic citation
- `docs/kb/adr/` — 8 architecture decision records
- `docs/kb/deps/` — monthly dependency audit logs

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
