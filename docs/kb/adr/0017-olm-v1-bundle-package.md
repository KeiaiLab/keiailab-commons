# ADR-0017: `pkg/bundle` — OLM v1 bundle metadata helpers

- Date: 2026-05-26
- Status: Accepted
- Authors: @eightynine01

## Context

OLM v1 (Operator Lifecycle Manager v1) is the next-generation operator
lifecycle management system for Kubernetes. It replaces OLM v0's
Subscription / CSV / InstallPlan model with a simplified ClusterExtension
API and File-Based Catalog (FBC) format.

Downstream operators that distribute via OLM v1 require:

1. **registry+v1 bundle metadata** — six required annotation keys in
   `metadata/annotations.yaml` and Docker LABELs.
2. **FBC catalog entries** — `olm.package`, `olm.channel`, `olm.bundle`
   schema definitions for catalog images.
3. **Bundle directory structure validation** — `manifests/` + `metadata/`
   layout compliance.

Without a shared library, each downstream operator reimplements these
metadata structures independently, producing silent inconsistencies.

## Decision

Add `pkg/bundle` (Experimental tier) with three concerns:

1. **Annotations** — constants for the six required bundle annotation keys
   plus a builder (`NewAnnotations` → `Map()` / `DockerLabels()`).
2. **FBC types** — Go structs for `olm.package`, `olm.channel`,
   `olm.bundle`, `olm.deprecations` with JSON serialization.
3. **Validate** — `ValidateDir(path)` checks bundle directory structure
   compliance (manifests/ + metadata/ + annotations.yaml).

Design invariants preserved:

- **Zero external dependency** — no `operator-framework` imports. Pure
  Go stdlib + `encoding/json` only.
- **Build-time helper** — no runtime controller or CRD dependency.
- **Leaf-package principle** — consistent with the existing
  `pkg/finalizer`, `pkg/labels` pattern.

## Consequences

### Positive

- Downstream operators share a single source of truth for OLM v1 bundle
  metadata — annotation keys, FBC schema, and directory validation.
- No new external dependency introduced.
- Experimental tier allows API iteration before stabilization.

### Negative

- API surface +1 package — justified by the shared need across
  downstream consumers.
- FBC types may drift from upstream `operator-registry` `declcfg`
  package — mitigated by schema constants matching the official spec.

## Alternatives Considered

1. **Import `operator-framework/operator-registry/alpha/declcfg`** —
   rejected. Adds a transitive dependency on the operator-framework
   module graph, violating the zero-external-dep principle.
2. **Makefile-only `opm` wrapper** — rejected. Does not provide Go-level
   type safety or reusable library patterns.

## References

- [OLM v1 documentation](https://operator-framework.github.io/operator-controller/)
- [FBC specification](https://olm.operatorframework.io/docs/reference/file-based-catalogs/)
- [registry+v1 bundle format](https://github.com/operator-framework/operator-registry/blob/v1.16.1/docs/design/operator-bundle.md)
- [ARCHITECTURE.md](../../ARCHITECTURE.md) — leaf-package design invariant.
