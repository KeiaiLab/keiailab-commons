<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **Shared Go library for keiailab operators — finalizer / labels / status / version / security / monitoring partials**
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

Shared Go library for **keiailab** Kubernetes operators (`mongodb-operator`,
`valkey-operator`, `postgresql-operator`).

> Status: **v0.x — API may break**. v1.0 onwards SemVer stable.

## Why

3 operators independently implemented identical scaffolding (PodSecurity
restricted contexts, version allowlists, NetworkPolicy templates, ServiceMonitor
builders). Maintenance drift between repos was already producing inconsistencies
— this library is the single source of truth.

## Packages (v0.8.0)

| Package | Purpose |
|---|---|
| `pkg/version` | Supported DB version allowlist convention (`MustList`, `IsSupported`, `Strings`, `Default`) + generic `Matrix[E MatrixEntry]`. |
| `pkg/security` | PodSecurity *restricted* SecurityContext builder with functional options. |
| `pkg/labels` | Recommended Kubernetes labels (`app.kubernetes.io/*`) builder — `Set`, `All()`, `Selector()` (version-aware split). |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` builder (unstructured — CRD-soft). |
| `pkg/networkpolicy` | NetworkPolicy builder — deny-by-default + functional options (`WithSelfIngress`, `WithIngressFromPeers`, `WithDenyEgress`, `WithEgressToPeers`). |
| `pkg/webhook` | Admission validation helpers — `ValidateAllowedVersion` (exact match), `ValidateWithPredicate` (caller-supplied matcher e.g. semver-prefix). |
| `pkg/finalizer` | Finalizer helpers — `Add`/`Remove`/`Has` (controller-runtime 의존 회피, std `slices` 만 사용). |
| `pkg/status` | 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼 (`SetReady`, `SetAvailable`, `SetReadyFalse`). |

`pkg/conditions` 는 *upstream `k8s.io/apimachinery/pkg/api/meta.SetStatusCondition` 활용 권장* (commons 미추가 결정 — boundary 분석 결과, 자세히는 mongodb-operator HANDOFF iteration 32 참조).

## Adoption Matrix (3 operator)

| Operator | sec | ver | lab | mon | np | wh | 채택률 |
|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | **4/6 (67%)** |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **6/6 (100%)** 🎉 |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | **3/6 (50%)** |

valkey 가 *first 100% 채택* — 다른 operator 의 carbon-copy reference 역할. 적용
사례 commits:
- `pkg/security`: it8 (3 operator cross-cut) — `23fd3da` mongodb / `a0be4cf` valkey / `ac2e647` postgres
- `pkg/version`: mongodb it9 `a8db040`, valkey it8
- `pkg/labels`: mongodb it27 `ebc5803`, postgres it28 `c68b451`, valkey it29 `e8428b1`
- `pkg/monitoring`: valkey it23 `1765b54`
- `pkg/networkpolicy`: valkey it25 `97162b5`, mongodb it26 `ca0ec27`
- `pkg/webhook`: valkey it31 `14be0db`, postgres it34 `1d8fa17`

⏳ 영역은 *기능 추가 동반* (예: mongodb webhook server / ServiceMonitor reconciler)
또는 *별 추상화 적합* (postgres version matrix.go의 Combo struct 가 commons.MustList
보다 풍부 — 위임 부적합) 으로 *deepening 보류* 상태.

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

## Community

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — pkg API 질문, integration 사례, 새 helper 제안
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — 버그 / API 요청
- **Downstream**: 3 operator (mongodb-operator / postgres-operator / valkey-operator) — `go.mod replace` 또는 직접 `require` 로 사용
- **Stability matrix**: `pkg/labels`, `pkg/security`, `pkg/version`, `pkg/webhook` (stable v0.5+) / `pkg/networkpolicy`, `pkg/monitoring` (experimental)

## License

Apache-2.0 — see [LICENSE](./LICENSE). Zero AGPL/BUSL transitive dependency
goal (audited per minor release).

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a> ·
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/forgewise">forgewise</a>
</p>

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
