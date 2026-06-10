# Upgrading operator-commons

> **English** | [한국어](UPGRADING.ko.md) | [日本語](UPGRADING.ja.md) | [中文](UPGRADING.zh.md)

This document collects the migration steps needed when bumping a minor
or major version of the `github.com/keiailab/operator-commons` Go
module. It is the common entry point for downstream consumers.

## 0. Version policy (SemVer)

| Change type | SemVer bump | Example |
|---|---|---|
| New package added | minor (v0.X → v0.X+1) | `pkg/events`, `pkg/storageclass` introduced |
| Existing API signature change (breaking) | major (v0.X → v1.0 / v1.X → v2.0) | `pkg/status.SetReady()` signature change |
| Package-internal behaviour change (non-breaking) | patch (v0.X.Y → v0.X.Y+1) | bug fix |
| ADR deviation | major + Deprecated notice | API stability tier change |

API stability tier (`pkg/<name>/doc.go` marker):

- **Stable** — backwards-compatible across a minor release.
- **Beta** — may change in the next minor.
- **Experimental** — may change at any time.

## 1. v0.7.x → v0.8.x

### Helm library chart consumers

```bash
helm dep update charts/<your-operator>
helm template <your-operator> charts/<your-operator>
```

The `keiailab-commons` chart v0.8.0 adds
`keiailab.secrets.externalSecret`. It is opt-in and does not render
unless the downstream chart calls the helper.

Recommended values shape for downstream charts:

```yaml
externalSecrets:
  enabled: false
  clusterSecretStore: infisical
  refreshInterval: 1h
```

Use the helper only for generated `ExternalSecret` resources. Workload
CRDs should continue to reference the materialized Kubernetes `Secret`
name so existing clusters do not see an API breaking change.

### Go module consumers

```bash
go get github.com/keiailab/operator-commons@v0.8.0
go mod tidy
```

No additional work — backwards-compatible.

## 2. v0.8.x → v0.9.x

### New packages (minor bump)

| Package | Purpose | Tier |
|---|---|---|
| `pkg/pvc` | PVC expansion helpers | Beta |
| `pkg/topology` | PVC topology spread + zone-aware affinity | Beta |

### Migration

Add the imports in your downstream operator:

```go
import (
    "github.com/keiailab/operator-commons/pkg/pvc"
    "github.com/keiailab/operator-commons/pkg/topology"
)
```

### Backwards compatibility

- Existing packages (`pkg/status`, `pkg/finalizer`, `pkg/networkpolicy`,
  `pkg/monitoring`, `pkg/probes`, `pkg/labels`, `pkg/storageclass`,
  `pkg/webhook`, `pkg/events`, `pkg/security`, `pkg/version`) keep
  their signatures.
- The `keiailab-commons` chart's `_security.tpl` and
  `_servicemonitor.tpl` partials are *opt-in*; leaving the existing
  inline definitions alone has no effect.

### Recommended migration procedure

```bash
# 1. bump the dependency
go get github.com/keiailab/operator-commons@v0.9.0
go mod tidy

# 2. verify
make verify  # lint + test + build

# 3. e2e (kind)
kind create cluster
helm install <operator> charts/<operator>
kubectl apply -f config/samples/
kubectl get <CR> -A  # observe reconciliation
```

## 3. v0.9.x → v0.10.x

### New packages (minor bump)

| Package | Purpose | Tier |
|---|---|---|
| `pkg/bundle` | OLM v1 bundle metadata helpers | Experimental |

### Migration

Add the import in your downstream operator:

```go
import (
    "github.com/keiailab/operator-commons/pkg/bundle"
)
```

### Backwards compatibility

- All existing packages keep their signatures.
- `pkg/bundle` is opt-in — no existing code is affected.

### Recommended migration procedure

```bash
# 1. bump the dependency
go get github.com/keiailab/operator-commons@v0.10.0
go mod tidy

# 2. verify
make verify  # lint + test + build

# 3. generate bundle metadata
annotations := bundle.NewAnnotations("my-operator", []string{"stable"}, "stable")
annotationMap := annotations.Map()
```

## 4. v0.9.x → v1.0.0

Proceeds when the v1.0.0 graduation criteria (see
[STABILITY.md](STABILITY.md) "v1.0.0 graduation") are satisfied:

- All packages reach Stable tier.
- v0.x → v1.0 is a *naming* change — semantics are unchanged (no
  breaking change).

## 5. General migration checklist

Before upgrade:

- [ ] `go mod tidy` produces no change (drift = 0).
- [ ] `make audit` passes (govulncheck CVE = 0).
- [ ] Existing e2e suite passes.

After upgrade:

- [ ] Downstream import path updated (`go get -u` or pinned version).
- [ ] `make verify` passes.
- [ ] e2e passes.
- [ ] Helm chart `charts/<operator>` `dependencies:` updated.

## 6. Breaking-change notice policy

- **Deprecation**: add `// Deprecated:` comment in the new minor; remove
  two minors later.
- **Breaking**: major bump + dedicated section in this file + ADR.
- **No silent breaking changes**: every breaking change carries at least
  one minor of prior deprecation.

## References

- ADR index: [`docs/kb/adr/INDEX.md`](kb/adr/INDEX.md).
- API stability: `pkg/<name>/doc.go` tier marker.
- i18n: [`docs/i18n/README.md`](i18n/README.md) (multilingual policy).

---

<p align="center">© 2026 keiailab</p>
