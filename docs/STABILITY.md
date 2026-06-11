# API Stability Promise

> **English** | [한국어](STABILITY.ko.md) | [日本語](STABILITY.ja.md) | [中文](STABILITY.zh.md)

> `keiailab-commons` API stability tiers and breaking-change policy.

## Tiers

`keiailab-commons` uses a three-tier stability model:

| Tier | Compatibility | Breaking change policy |
|---|---|---|
| **Stable** | Backwards-compatible across minor releases. Only patch fixes. | Requires major version bump (`v2.0.0`) plus a 6-month deprecation notice. |
| **Beta** | Compatible across patch releases. Minor releases may include source-compatible API improvements with at least one release of deprecation. | Minor version bump permitted with a deprecation entry in `CHANGELOG.md`. |
| **Experimental** | No compatibility promise. May break across any release. | Any release; must be flagged in `CHANGELOG.md` "BREAKING" section. |

## Current tier matrix

The current tier of each package mirrors [ROADMAP.md](ROADMAP.md) "API
Stability Tier":

| Package | Tier | Promotion criterion |
|---|---|---|
| `pkg/finalizer` | Stable | (v1 entry — no additional work). |
| `pkg/labels` | Stable | (v1 entry). |
| `pkg/status` | Stable | (v1 entry). `update.go` (`UpdateWithRetry`) is a Beta surface inside the Stable package. |
| `pkg/storageclass` | Stable | Trivial validation surface (regex + nil check). |
| `pkg/version` | Beta | Generic `Matrix[E]` cross-repo verify. |
| `pkg/monitoring` | Beta | ServiceMonitor downstream e2e. |
| `pkg/networkpolicy` | Beta | 4-direction (ingress / egress × TCP / UDP) verify. |
| `pkg/security` | Beta | Restricted PSA guard across downstream. |
| `pkg/events` | Beta | Downstream Reconcile path adoption + Event reason consistency. |
| `pkg/pvc` | Beta | Downstream PVC expansion live adoption. |
| `pkg/topology` | Beta | Downstream topology spread live adoption. |
| `pkg/apply` | Beta | Downstream idempotent apply live adoption. |
| `pkg/reconcile` | Beta | Downstream reconcile-loop helper live adoption. |
| `pkg/certmanager` | Beta | Downstream Certificate / Issuer builder live adoption. |
| `pkg/reconcilemetrics` | Beta | Downstream adoption + Prometheus time-series name preservation verify. |
| `pkg/webhook` | Experimental | Multi-downstream adoption + stabilization. |
| `pkg/probes` | Experimental | 2+ downstream live adoption. |
| `pkg/bundle` | Experimental | 2+ downstream bundle adoption. |

## Promotion process

1. A PR opens a proposed promotion (e.g. `feat(pkg/X): promote to Stable`).
2. Promotion criterion (per ROADMAP) is verified via the local quality
   gates:
   - Downstream import passes.
   - godoc coverage on the package ≥ 80 %.
   - Unit + integration test coverage ≥ 85 %.
   - No `// TODO` or `// FIXME` in exported API.
3. [ROADMAP.md](ROADMAP.md) tier table is updated in the same PR.
4. [CHANGELOG.md](../CHANGELOG.md) gets a "Changed" entry.

## Breaking-change policy

A **breaking change** is any of:

- Removed exported identifier (function, type, constant, variable).
- Changed exported signature (parameters, return types).
- Removed package.
- Behavior change that requires caller code modification.

### Per tier

- **Stable**: forbidden until v2.0.0. Use deprecation — add `// Deprecated: …`
  godoc plus a new alternative; keep the old API for at least 6 months.
- **Beta**: permitted with at least one release of deprecation. Must appear
  in the `CHANGELOG.md` "Deprecated" → "Removed" pipeline.
- **Experimental**: permitted at any release; must appear in the
  `CHANGELOG.md` "BREAKING" section.

## Semantic versioning

`vMAJOR.MINOR.PATCH` per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html):

- **MAJOR**: breaking changes to a Stable tier — requires v2.0.0 graduation
  review.
- **MINOR**: new features, Beta tier additions, non-breaking Stable
  improvements.
- **PATCH**: bug fixes only, no API surface change.

## v1.0.0 graduation

Requires *all* of:

1. All Stable-candidate packages reach Stable tier.
2. Zero BREAKING CHANGE across six or more consecutive minor releases (v0.9
   → v0.14).
3. godoc coverage ≥ 80 % (this document + per-package) — verified by
   `scripts/godoc-coverage.sh`.
4. CITATION.cff + DOI.
5. Downstream import end-to-end verification on `v1.0.0-rc.N`.
6. `go vet ./... && go test ./...` clean with ≥ 85 % coverage.
7. [CHANGELOG.md](../CHANGELOG.md) v0.x evolution history + v1.0.0 release
   notes.
8. This `docs/STABILITY.md` (you are here).
9. `pkg/finalizer` multi-finalizer order guarantee.
10. `pkg/labels` K8s 1.30+ v2 mapping.
11. `pkg/status` Condition reason catalog docs.

## Caller responsibilities

Downstream consumers should:

- Pin to a specific `vMAJOR.MINOR.PATCH` in `go.mod` until v1.0.0.
- Subscribe to [CHANGELOG.md](../CHANGELOG.md) for deprecation warnings.
- Test against `v1.0.0-rc.N` releases before GA.

## References

- [ROADMAP.md](ROADMAP.md) — tier table + graduation criteria.
- [CHANGELOG.md](../CHANGELOG.md) — version history.
- [CITATION.cff](../CITATION.cff) — academic citation.
- [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html).
