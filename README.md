# operator-commons

Shared Go library for **keiailab** Kubernetes operators (`mongodb-operator`,
`valkey-operator`, `postgresql-operator`).

> Status: **v0.x — API may break**. v1.0 onwards SemVer stable.

## Why

3 operators independently implemented identical scaffolding (PodSecurity
restricted contexts, version allowlists, NetworkPolicy templates, ServiceMonitor
builders). Maintenance drift between repos was already producing inconsistencies
— this library is the single source of truth.

## Packages (v0.1.0 MVP)

| Package | Purpose |
|---|---|
| `pkg/version` | Supported DB version allowlist convention (`MustList`, `IsSupported`, `Strings`, `Default`). |
| `pkg/security` | PodSecurity *restricted* SecurityContext builder with functional options. |

Planned (v0.2.0): `pkg/networkpolicy`, `pkg/monitoring`, `pkg/labels`, `pkg/webhook`.

## Usage

```go
import (
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/version"
)

var SupportedMongoDBVersions = version.MustList("8.0", "8.2", "8.3")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

## Versioning + Release

- v0.x: API breaking allowed. Each tag (`v0.N.M`) bumps either pkg, public-API, or
  significant behavior.
- Each consuming operator pins via `go.mod` `require` — `replace` directive
  is acceptable during local development across this repo + the 3 operators.
- v1.0 onwards: Semantic Versioning. Breaking changes require RFC.

## License

Apache-2.0 — see [LICENSE](./LICENSE). Zero AGPL/BUSL transitive dependency
goal (audited per minor release).
