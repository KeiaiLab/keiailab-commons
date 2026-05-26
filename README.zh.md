<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **用于 Kubernetes operator 通用 scaffolding 的 Go 共享库 — finalizer / labels / status / version / security / monitoring partials.**
>
> [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | **中文**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>


> ⚠️ This translation is AI-generated and pending native review.

---

可复用的 Go 库,用于消除 Kubernetes operator 代码库中的 scaffolding 漂移 ——
PodSecurity restricted context、支持版本 allowlist、NetworkPolicy 模板、
ServiceMonitor 构建器、finalizer / status 帮助函数,以及 Helm library chart
partial,封装在一个小而稳定的 API 表面之后。

> 状态: **v0.x — API 可能变更.** v1.0 起遵循 SemVer stable.

## Why

Operator 作者反复实现相同的 scaffolding —— restricted PodSecurity context、
支持版本矩阵、default-deny NetworkPolicy、ServiceMonitor 构建器、finalizer
帮助函数、status condition 目录。各自独立重新实现会在相似 reconciler 之间产生
隐性不一致,并随着 minor 修订逐渐分叉。`operator-commons` 是该 scaffolding 的
单一来源 —— 导入帮助函数,获得 canonical 实现,无需在每个仓库重新发明。

## 包

| 包 | Tier | 用途 |
|---|---|---|
| `pkg/finalizer` | Stable | Finalizer 帮助 — `Add` / `Remove` / `Has` / `EnsureOrder` (仅 stdlib `slices`,无 controller-runtime 依赖)。 |
| `pkg/labels` | Stable | Kubernetes 推荐标签 (`app.kubernetes.io/*`) 构建器 — `Set`、`All()`、`Selector()`,以及 v2 映射 (`AllV2`)。 |
| `pkg/status` | Stable | 4 个标准 Condition Type + 6 个 Reason 目录 + 帮助函数 (`SetReady`、`SetAvailable`、`SetReadyFalse`)。 |
| `pkg/storageclass` | Stable | DNS-1123 storageClass 验证器 + `Normalize` / `MustNormalize` (empty → cluster default 指针)。 |
| `pkg/version` | Beta | 版本 allowlist 约定 (`MustList`、`IsSupported`、`Strings`、`Default`) + 泛型 `Matrix[E MatrixEntry]` + 序列化器。 |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` 与 `PrometheusRule` 构建器 (unstructured — CRD-soft)。 |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy 构建器 + functional options (`WithSelfIngress`、`WithIngressFromPeers`、`WithDenyEgress`、`WithEgressToPeers`、`ComboPeer`)。 |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext 构建器 + Pod / Container 分离 + seccomp profile 指针。 |
| `pkg/events` | Beta | 最小 `Recorder` 接口 + 9 个标准 `Reason` 常量 + `Emit` / `EmitWarning` / `WrappedError` (nil-safe)。 |
| `pkg/pvc` | Beta | PVC 扩展助手 — 比较 + 安全就地更新 (controller-runtime 依赖 — ADR-0016)。 |
| `pkg/topology` | Beta | TopologySpreadConstraints HA 默认值 + 区域感知亲和性构建器。 |
| `pkg/probes` | Experimental | `corev1.Probe` fluent 构建器 — HTTP / HTTPS / TCP / Exec,kubelet 默认值 + clamp。 |
| `pkg/webhook` | Experimental | Admission validation 帮助 — `ValidateAllowedVersion`、`ValidateWithPredicate`、conversion registry。 |
| `pkg/bundle` | Experimental | OLM v1 捆绑包元数据助手 — 注解、FBC模式类型、目录验证 (ADR-0017)。 |

[docs/STABILITY.md](docs/STABILITY.md) 定义 tier 承诺。
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) 涵盖包表面与设计不变量。
[docs/ROADMAP.md](docs/ROADMAP.md) 跟踪 tier 晋升标准与 v1.0 毕业清单。

## 使用

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

每个包的示例位于对应的 `pkg/<name>/doc.go` 包文档中
(`go doc github.com/keiailab/operator-commons/pkg/<name>`)。

## 版本与发布

- **v0.x**: 允许 API breaking。每个 tag (`v0.N.M`) 会改变某个包的公开 API
  或具有意义的行为。消费者通过 `go.mod` 固定特定版本。
- **v1.0 起**: Semantic Versioning。Breaking change 需要 ADR (`docs/kb/adr/`)。
- 本地 `replace` 指令在 cross-repo 开发中可接受; 发布 tag 始终保留 canonical
  module path。

## 社区

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — 包 API 提问、integration 案例、新 helper 提案。
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — bug 与具体特性请求。
- **Security**: 私密上报流程见 [SECURITY.md](SECURITY.md)。
- **Contributing**: 开发流程见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可证

Apache-2.0 — 见 [LICENSE](LICENSE)。AGPL / BUSL transitive 依赖 0 件目标
(每个 minor 发布审计一次)。

---


<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
