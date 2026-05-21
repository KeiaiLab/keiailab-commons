<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab Operator 系列共享的 Go 库 — finalizer / labels / status / version / security / monitoring partials**
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

> **注意 (Notice)**: 本中文 README 由机器翻译生成,处于 *partial* (RFC-0025 `[~]`) 状态。技术内容以 [README.md](README.md) (英文) 为准。母语审阅者的完整审定计划在后续周期进行。

**keiailab** Kubernetes Operator (`mongodb-operator`、`valkey-operator`、`postgresql-operator`) 共享的 Go 库。

> 状态: **v0.x — API 可能发生破坏性变更**。v1.0 之后将采用 SemVer (语义化版本) 保持稳定。

## 为什么 (Why)

3 个 Operator 各自独立地实现了相同的脚手架代码 (PodSecurity restricted 安全上下文、版本 allowlist、NetworkPolicy 模板、ServiceMonitor 构建器)。仓库之间的维护漂移已经开始产生不一致 — 本库即唯一可信来源 (single source of truth)。

## 软件包列表 (Packages, v0.8.0)

| 软件包 | 用途 |
|---|---|
| `pkg/version` | 支持的 DB 版本 allowlist 规约 (`MustList`、`IsSupported`、`Strings`、`Default`) + 泛型 `Matrix[E MatrixEntry]`。 |
| `pkg/security` | 带 functional option 的 PodSecurity *restricted* SecurityContext 构建器。 |
| `pkg/labels` | 推荐 Kubernetes 标签 (`app.kubernetes.io/*`) 构建器 — `Set`、`All()`、`Selector()` (版本感知分支)。 |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 构建器 (unstructured — CRD-soft)。 |
| `pkg/networkpolicy` | NetworkPolicy 构建器 — deny-by-default + functional option (`WithSelfIngress`、`WithIngressFromPeers`、`WithDenyEgress`、`WithEgressToPeers`)。 |
| `pkg/webhook` | Admission 验证助手 — `ValidateAllowedVersion` (精确匹配)、`ValidateWithPredicate` (调用方提供的 matcher,例如 semver-prefix)。 |
| `pkg/finalizer` | Finalizer 助手 — `Add` / `Remove` / `Has` (避免依赖 controller-runtime,仅使用标准库 `slices`)。 |
| `pkg/status` | 4 个标准 Condition Type + 6 个 Reason 目录 + 助手函数 (`SetReady`、`SetAvailable`、`SetReadyFalse`)。 |

`pkg/conditions` *推荐使用上游 `k8s.io/apimachinery/pkg/api/meta.SetStatusCondition`* (commons 不添加的决定 — 经过 boundary 分析得出。详见 mongodb-operator HANDOFF iteration 32)。

## 采用矩阵 (Adoption Matrix, 3 个 Operator)

| Operator | sec | ver | lab | mon | np | wh | 采用率 |
|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | **4/6 (67%)** |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **6/6 (100%)** 🎉 |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | **3/6 (50%)** |

valkey 是 *首个达成 100% 采用* 的 Operator — 承担其他 Operator 的 carbon-copy 参考职责。应用案例 commit:

- `pkg/security`: it8 (3 个 Operator cross-cut) — `23fd3da` mongodb / `a0be4cf` valkey / `ac2e647` postgres
- `pkg/version`: mongodb it9 `a8db040`、valkey it8
- `pkg/labels`: mongodb it27 `ebc5803`、postgres it28 `c68b451`、valkey it29 `e8428b1`
- `pkg/monitoring`: valkey it23 `1765b54`
- `pkg/networkpolicy`: valkey it25 `97162b5`、mongodb it26 `ca0ec27`
- `pkg/webhook`: valkey it31 `14be0db`、postgres it34 `1d8fa17`

⏳ 区域处于 *推迟深入* (deepening 保留) 状态 — 要么 *伴随功能新增* (例如 mongodb webhook server / ServiceMonitor reconciler),要么 *其他抽象更合适* (postgres `version matrix.go` 中的 `Combo` struct 比 `commons.MustList` 更丰富,不适合委托)。

## 使用方法 (Usage)

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

## 版本管理与发布 (Versioning + Release)

- v0.x: 允许 API 的破坏性变更。每个 tag (`v0.N.M`) 表示软件包、公开 API、或重要行为之一发生 bump。
- 每个 consumer Operator 通过 `go.mod` 的 `require` 进行版本固定 — 在本仓库与 3 个 Operator 之间进行本地开发时,也允许使用 `replace` 指令。
- v1.0 之后: 语义化版本 (Semantic Versioning)。破坏性变更必须经过 RFC 流程。

## 社区 (Community)

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — pkg API 的疑问、集成案例、新助手函数提案
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — bug 报告 / API 请求
- **Downstream**: 3 个 Operator (mongodb-operator / postgres-operator / valkey-operator) — 通过 `go.mod replace` 或直接 `require` 使用
- **稳定性矩阵**: `pkg/labels`、`pkg/security`、`pkg/version`、`pkg/webhook` (自 v0.5+ stable) / `pkg/networkpolicy`、`pkg/monitoring` (experimental)

## 许可证 (License)

Apache-2.0 — 详见 [LICENSE](./LICENSE)。以零 AGPL/BUSL 传递依赖为目标,每个 minor 版本进行审计。

## 参考 (References)

- [English README](README.md) — canonical SSOT (规范正本)
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
