# API Stability Promise

> operator-commons API stability tiers and breaking-change policy.

## Tiers

operator-commons uses 3-tier stability:

| Tier | Compatibility | Breaking change policy |
|---|---|---|
| **Stable** | Backwards-compatible across minor releases. Only patch fixes. | Requires major version bump (`v2.0.0`) + 6+ month deprecation notice |
| **Beta** | Compatible across patch releases. Minor releases may include source-compatible API improvements with at least 1-release deprecation window | Minor version bump permitted with deprecation entry in `CHANGELOG.md` |
| **Experimental** | No compatibility promise. May break across any release. | Any release — must be flagged in `CHANGELOG.md` "BREAKING" section |

## Current tier matrix

Per `ROADMAP.md` "API Stability Tier" table:

| Package | Tier | Promotion criterion |
|---|---|---|
| `pkg/finalizer` | Stable | (v1 entry — no additional work) |
| `pkg/labels` | Stable | (v1 entry) |
| `pkg/status` | Stable | (v1 entry) |
| `pkg/version` | Beta | Generic `Matrix[E]` 3-repo verify |
| `pkg/monitoring` | Beta | ServiceMonitor 3-repo e2e |
| `pkg/networkpolicy` | Beta | 4-direction TCP/UDP verify |
| `pkg/security` | Beta | restricted PSA 3-repo guard |
| `pkg/webhook` | Experimental | Multi-repo adoption + stabilize |

## Promotion process

1. Sub-task PR opens with proposed promotion (e.g., `feat(pkg/X): promote to Stable`)
2. Promotion criterion (per ROADMAP) verified via CI:
   - Cross-repo import passes (3 operators)
   - godoc coverage on package ≥80%
   - Unit + integration test coverage ≥85%
   - No `// TODO` or `// FIXME` in exported API
3. ROADMAP.md tier table updated in same PR
4. `CHANGELOG.md` entry under "Changed"

## Breaking-change policy

A **breaking change** is any of:
- Removed exported identifier (function / type / constant / variable)
- Changed exported signature (parameters, return types)
- Removed package
- Behavior change that requires caller code modification

### For each tier:

- **Stable**: forbidden until v2.0.0. Use deprecation: add `// Deprecated: ...` godoc + new alternative, keep old for 6+ months
- **Beta**: permitted with 1-release deprecation. Must appear in `CHANGELOG.md` "Deprecated" → "Removed" pipeline
- **Experimental**: permitted any release, must appear in `CHANGELOG.md` "BREAKING" section

## Semantic versioning

`vMAJOR.MINOR.PATCH` per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html):
- **MAJOR**: breaking changes to Stable tier — requires v2.0.0 graduation review
- **MINOR**: new features, Beta tier additions, non-breaking Stable improvements
- **PATCH**: bug fixes only, no API surface change

## v1.0.0 graduation

Requires *all* of:
1. 8/8 packages reach Stable tier
2. 0 BREAKING CHANGE across 6+ consecutive minor releases (v0.8 → v0.13)
3. godoc coverage ≥80% (this doc + per-package — verified by `scripts/godoc-coverage.sh`)
4. CITATION.cff + Zenodo DOI
5. 3-repo import e2e verification on `v1.0.0-rc.N`
6. `go vet ./... && go test ./...` clean with ≥85% coverage
7. CHANGELOG.md v0.x evolution history + v1.0.0 release notes
8. This `docs/STABILITY.md` (you are here)
9. `pkg/finalizer` multi-finalizer order guarantee
10. `pkg/labels` K8s 1.30+ v2 mapping
11. `pkg/status` Condition reason catalog docs

Tracking: `~/.claude/plans/2026-05-14-4-operators-100pct/P-B.md` (29 sub-tasks).

## Caller responsibilities

Callers (mongodb-operator, valkey-operator, postgres-operator):
- Pin to `vMAJOR.MINOR.PATCH` in `go.mod` until v1.0.0
- Subscribe to `CHANGELOG.md` for deprecation warnings
- Test against `v1.0.0-rc.N` before GA

## References

- `ROADMAP.md` — tier table + graduation criteria
- `CHANGELOG.md` — version history
- `CITATION.cff` — academic citation
- `ADOPTERS.md` — 3-repo adoption matrix
- [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html)
