# ADR-0001: operator-commons charter

- Date: 2026-05-07
- Status: Accepted
- Authors: @eightynine01

## Context

Multiple Kubernetes operator codebases were independently re-implementing
identical scaffolding:

- PodSecurity *restricted* SecurityContext builders (inline occurrences
  across StatefulSet / Deployment / init-container code paths).
- "Supported database versions" allowlist plus webhook validation (with
  divergent function names across implementations).
- NetworkPolicy deny-by-default templates.
- ServiceMonitor builders with shared label conventions.

Drift between separate operator repositories was already producing
inconsistencies — for example, hardening that was introduced after a
production PodSecurity admission incident in one operator had not yet been
propagated to others, and version-allowlist signatures diverged.

## Decision

Create `github.com/keiailab/operator-commons` as a **separate Go module
repository** (not a monorepo module).

Rationale:

- Independent SemVer release plus `replace`-directive workflow.
- Cross-cut changes still require a PR per downstream consumer, which
  enforces explicit review of each adoption.
- Helm library chart partials can ride alongside the Go module.
- Per-repo policy with a shared base matches the existing development
  workflow.

### v0.x charter

- Apache-2.0 license, **zero AGPL / BUSL transitive dependency** goal,
  audited per minor release.
- API breaking changes are allowed in v0.x; v1.0 onwards SemVer stable.
- Each tag (`v0.N.M`) bumps either a package, a public API, or a meaningful
  behavior — no silent rewrites.
- Downstream consumers pin via `go.mod` `require`; `replace` directive
  permitted during local cross-cut development.

## Consequences

### Positive

- Single source of truth for PodSecurity / version allowlist /
  NetworkPolicy.
- A v0.1.0 MVP (security + version) becomes immediately consumable by
  downstream operators.
- Future packages (networkpolicy, monitoring, labels, webhook, …) ship as
  v0.2.0+ without breaking v0.1.0 consumers.

### Negative

- A cross-cut breaking change requires one PR per downstream consumer —
  the same overhead as separate repos today, but with explicit visibility.
- Local development across `operator-commons` plus downstream consumers
  requires either a Go workspace or a `replace` directive — slightly more
  setup than a monorepo.

### Trade-offs

- *Independent release plus explicit cross-cut review* (this ADR) vs.
  *single PR plus monorepo* — rejected because a monorepo migration is
  itself a separate large effort.

## Alternatives Considered

1. **Monorepo (`keiailab/keiai-platform`)** — rejected: monorepo bootstrap
   is itself a separate effort, and the `operator-commons` MVP is needed
   sooner.
2. **Vendor commons into a single operator's `internal/commons`** —
   rejected: delays the deduplication benefit and elevates one downstream
   as the de-facto canonical implementation without explicit governance.

## Refs

- Local Helm library chart: `charts/keiailab-commons/`.
- Format: Nygard 5 sections.
