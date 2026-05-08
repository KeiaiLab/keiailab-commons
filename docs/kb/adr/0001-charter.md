# ADR-0001: operator-commons charter

- Date: 2026-05-07
- Status: Accepted
- Authors: @eightynine01

## Context

3 keiailab Kubernetes operators (`mongodb-operator`, `valkey-operator`,
`postgresql-operator`) were independently re-implementing identical
scaffolding code:

- PodSecurity *restricted* SecurityContext builders (each with 4–5 inline
  occurrences across StatefulSet / Deployment / init-container code paths).
- "Supported DB versions" allowlist + webhook validation (3 different
  function names: `IsSupported`, `IsSupportedValkeyVersion`, none in
  postgres).
- NetworkPolicy deny-by-default templates.
- ServiceMonitor builders with shared label conventions.

Drift between repositories was already manifesting:
- valkey adopted `SupportedValkeyVersions` as a Go slice constant; mongodb
  did not have explicit validation; postgres had `IsSupported()` only for a
  custom semver type. **Pattern fragmentation P0** flagged in iteration 8 plan.
- mongodb's `buildKeyfileInitContainerSecurityContext` was added only after
  a PodSecurity admission incident (2026-05-07) — the same pattern was
  *not* yet hardened in valkey and postgres.

## Decision

Create `github.com/keiailab/operator-commons` as a **separate repository**
(not a monorepo module).

Rationale (from iteration 8 plan AskUserQuestion):
- Independent SemVer release + replace-directive workflow.
- Cross-cut changes still require 3 PRs (one per consuming operator) but
  this enforces explicit review of each adoption.
- ghcr publish pattern matches the consuming operators.
- ADR-0011 (mongodb pre-commit) precedent confirms keiailab tolerates
  per-repo policy with shared base.

### v0.x charter

- Apache-2.0 license, **zero AGPL/BUSL transitive dependency** goal.
  Audited per minor release.
- API breaking allowed in v0.x; v1.0 onwards SemVer stable.
- Each tag (`v0.N.M`) bumps either pkg, public API, or significant
  behavior — no silent rewrites.
- 3 operators pin via `go.mod` `require`; `replace` directive permitted
  during local cross-cut development.

## Consequences

### Positive
- Single source of truth for PodSecurity / version allowlist / NetworkPolicy.
- v0.1.0 MVP (security + version) immediately consumable by 3 operators
  in iteration 8.
- Future packages (networkpolicy, monitoring, labels, webhook) ship as
  v0.2.0+ without breaking v0.1.0 consumers.

### Negative
- 3 PRs needed for any cross-cut breaking change — same overhead as
  3 separate operators today, but with explicit visibility.
- Local development across operator-commons + 3 consumers requires
  Go workspace or `replace` directive — extra setup vs monorepo.

### Trade-offs
- *Independent release / explicit cross-cut review* (this ADR) vs
  *single PR for cross-cut / monorepo* (rejected in iteration 8 because
  a monorepo migration would itself be a separate iteration).

## Alternatives Considered

1. **Monorepo (`keiailab/keiai-platform`)** — rejected: monorepo bootstrap
   is itself a separate iteration; operator-commons MVP urgent now.
2. **Vendor commons into `mongodb-operator/internal/commons`** — rejected:
   delays the deduplication benefit; valkey + postgres still on stale
   inline copies; promotes mongodb as de-facto canonical without explicit
   governance.

## Refs

- iteration 8 plan: `~/.claude/plans/iridescent-squishing-locket.md`
- mongodb-operator ADR-0011 (pre-commit policy)
- iteration 6 PodSecurity incident (HANDOFF.md 2026-05-07)
