<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **Shared Go library for Kubernetes operator scaffolding — finalizer / labels / status / version / security / monitoring partials.**
>
> **English** | [한국어](README.ko.md) | [日本語](README.ja.md) | [中文](README.zh.md)

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>

<p align="center">
  <b>English</b> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>

---

A reusable Go library that removes scaffolding drift from Kubernetes operator
codebases — PodSecurity restricted contexts, supported-version allowlists,
NetworkPolicy templates, ServiceMonitor builders, finalizer / status helpers,
and Helm library chart partials, packaged behind a small, stable API surface.

> Status: **v0.x — API may break.** v1.0 onwards SemVer stable.

## Why

Operator authors repeatedly implement identical scaffolding — restricted
PodSecurity contexts, supported-version matrices, default-deny NetworkPolicies,
ServiceMonitor builders, finalizer helpers, status condition catalogs.
Independent re-implementation produces silent inconsistencies between similar
reconcilers and gradually drifts apart on minor revisions. `operator-commons`
is the single source of truth for that scaffolding: import the helper, get the
canonical implementation, and stop re-inventing it in every repository.

## Packages

| Package | Tier | Purpose |
|---|---|---|
| `pkg/finalizer` | Stable | Finalizer helpers — `Add` / `Remove` / `Has` / `EnsureOrder` (stdlib `slices` only, no controller-runtime dependency). |
| `pkg/labels` | Stable | Recommended Kubernetes labels (`app.kubernetes.io/*`) builder — `Set`, `All()`, `Selector()`, plus v2 mapping (`AllV2`). |
| `pkg/status` | Stable | Four standard Condition Types + six Reason catalog + helpers (`SetReady`, `SetAvailable`, `SetReadyFalse`). |
| `pkg/storageclass` | Stable | DNS-1123 storageClass validator + `Normalize` / `MustNormalize` (empty → cluster default pointer). |
| `pkg/version` | Beta | Version allowlist convention (`MustList`, `IsSupported`, `Strings`, `Default`) + generic `Matrix[E MatrixEntry]` + serializer. |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` and `PrometheusRule` builders (unstructured — CRD-soft). |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy builder + functional options (`WithSelfIngress`, `WithIngressFromPeers`, `WithDenyEgress`, `WithEgressToPeers`, `ComboPeer`). |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext builder + Pod / Container split + seccomp profile pointers. |
| `pkg/events` | Beta | Minimal `Recorder` interface + nine standard `Reason` constants + `Emit` / `EmitWarning` / `WrappedError` (nil-safe). |
| `pkg/probes` | Experimental | `corev1.Probe` fluent builder — HTTP / HTTPS / TCP / Exec with kubelet defaults and clamp. |
| `pkg/webhook` | Experimental | Admission validation helpers — `ValidateAllowedVersion`, `ValidateWithPredicate`, conversion registry. |

[docs/STABILITY.md](docs/STABILITY.md) defines the tier promise.
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) covers the package surface and
design invariants. [docs/ROADMAP.md](docs/ROADMAP.md) tracks the tier promotion
criteria and the v1.0 graduation checklist.

## Usage

```go
import (
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/version"
    corev1 "k8s.io/api/core/v1"
)

var supportedVersions = version.MustList("1.0", "1.1", "1.2")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

Per-package examples live in the corresponding `pkg/<name>/doc.go` package
documentation (`go doc github.com/keiailab/operator-commons/pkg/<name>`).

## Versioning and release

- **v0.x**: API breaking changes are allowed. Each tag (`v0.N.M`) bumps either
  a package's public API or a meaningful behavioural change. Consumers pin a
  specific version via `go.mod`.
- **v1.0 onwards**: Semantic Versioning. Breaking changes require an ADR
  (`docs/kb/adr/`).
- A local `replace` directive is acceptable for cross-repo development; release
  tags always carry the canonical module path.

## Community

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — package API questions, integration patterns, new helper proposals.
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — bugs and concrete feature requests.
- **Security**: see [SECURITY.md](SECURITY.md) for the private disclosure process.
- **Contributing**: see [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow.

## License

Apache-2.0 — see [LICENSE](LICENSE). Zero AGPL / BUSL transitive dependency
goal (audited per minor release).

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
