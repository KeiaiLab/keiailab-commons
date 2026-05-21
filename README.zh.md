<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab 系列 operator 共享的 Go 库 — finalizer / labels / status / version / security / monitoring 等共用部件**
>
> [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | **中文**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <b>中文</b>
</p>

---

**keiailab** 系列 Kubernetes operator (`mongodb-operator`、`valkey-operator`、`postgresql-operator`) 共享的 Go 库。

> 状态: **v0.x — API 可能发生破坏性变更 (breaking change)**。v1.0 之后采用 SemVer (语义化版本) 保持稳定。

## 为什么需要

3 个 operator 各自独立实现了相同的基础骨架 (PodSecurity restricted 上下文、版本 allowlist、NetworkPolicy 模板、ServiceMonitor 构建器)。各仓库之间的维护漂移 (drift) 已经开始产生不一致 — 本库即是它们的唯一正本 (single source of truth)。

## 软件包列表 (v0.7.0)

| 软件包 | 用途 |
|---|---|
| `pkg/version` | 支持的 DB 版本 allowlist 约定 (`MustList`、`IsSupported`、`Strings`、`Default`) + 泛型 `Matrix[E MatrixEntry]`。 |
| `pkg/security` | 带 functional option 的 PodSecurity *restricted* SecurityContext 构建器。 |
| `pkg/labels` | 推荐的 Kubernetes 标签 (`app.kubernetes.io/*`) 构建器 — `Set`、`All()`、`Selector()` (版本感知 split)。 |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 构建器 (unstructured — CRD-soft)。 |
| `pkg/networkpolicy` | NetworkPolicy 构建器 — deny-by-default + functional option (`WithSelfIngress`、`WithIngressFromPeers`、`WithDenyEgress`、`WithEgressToPeers`)。 |
| `pkg/webhook` | Admission 校验辅助函数 — `ValidateAllowedVersion` (精确匹配)、`ValidateWithPredicate` (调用方提供的 matcher,例如 semver-prefix)。 |
| `pkg/finalizer` | Finalizer 辅助函数 — `Add` / `Remove` / `Has` (规避对 controller-runtime 的依赖,仅使用标准库 `slices`)。 |
| `pkg/status` | 4 种标准 Condition Type + 6 种 Reason 目录 + 辅助函数 (`SetReady`、`SetAvailable`、`SetReadyFalse`)。 |

`pkg/conditions` *推荐直接使用上游的 `k8s.io/apimachinery/pkg/api/meta.SetStatusCondition`* (决定不在 commons 中新增 — 边界 (boundary) 分析的结果,详情请参考 mongodb-operator HANDOFF iteration 32)。

## 采用矩阵 (3 个 operator)

| Operator | sec | ver | lab | mon | np | wh | 采用率 |
|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | **4/6 (67%)** |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **6/6 (100%)** 🎉 |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | **3/6 (50%)** |

valkey 是 *首个 100% 采用* — 作为其他 operator 的 carbon-copy 参照基准。适用案例 commits:

- `pkg/security`: it8 (3 个 operator cross-cut) — `23fd3da` mongodb / `a0be4cf` valkey / `ac2e647` postgres
- `pkg/version`: mongodb it9 `a8db040`、valkey it8
- `pkg/labels`: mongodb it27 `ebc5803`、postgres it28 `c68b451`、valkey it29 `e8428b1`
- `pkg/monitoring`: valkey it23 `1765b54`
- `pkg/networkpolicy`: valkey it25 `97162b5`、mongodb it26 `ca0ec27`
- `pkg/webhook`: valkey it31 `14be0db`、postgres it34 `1d8fa17`

⏳ 区域因 *伴随功能新增* (例如 mongodb 的 webhook server / ServiceMonitor reconciler) 或 *更适合独立抽象* (postgres 的 `version matrix.go` 中的 `Combo` struct 比 `commons.MustList` 更丰富,不适合委派) 而 *暂缓深度集成*。

## 使用方法

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

## 版本控制与发布

- v0.x: 允许 API 的 breaking change。每个 tag (`v0.N.M`) 都会 bump 软件包、公共 API 或者重要的行为之一。
- 各 consumer operator 通过 `go.mod` 的 `require` 进行 pin — 在本仓库与 3 个 operator 之间进行本地开发时,允许使用 `replace` 指令。
- v1.0 之后: 采用语义化版本控制 (Semantic Versioning)。Breaking change 必须经过 RFC。

## 社区

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — pkg API 的问题、集成案例、新辅助函数的提议
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — Bug 报告 / API 请求
- **下游 (Downstream)**: 3 个 operator (mongodb-operator / postgres-operator / valkey-operator) — 通过 `go.mod replace` 或直接 `require` 使用
- **稳定性矩阵**: `pkg/labels`、`pkg/security`、`pkg/version`、`pkg/webhook` (v0.5+ 已 stable) / `pkg/networkpolicy`、`pkg/monitoring` (experimental)

## 许可证

Apache-2.0 — 请参阅 [LICENSE](./LICENSE)。每个 minor 版本审计零 AGPL/BUSL 传递依赖为目标。

## 参考

- [English README](README.md) — canonical SSOT (正本)
- [한국어 README](README.ko.md) — 韩文版
- [日本語 README](README.ja.md) — 日文版
- [中文术语表](docs/i18n/glossary-zh.md) — 标准术语表 (本仓库)

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
