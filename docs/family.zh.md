<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# keiailab 操作器家族

> ⚠️ This translation is AI-generated and pending native review. — 本翻译为 Claude 机器翻译结果。母语审阅者校验前为 `[待校验]` 状态。

> 基于共同基础构建的 4 个姊妹 Kubernetes 操作器 — `operator-commons` (Go 库) + Helm partial + Apache-2.0 技术栈。

本页面正从 **`operator-commons`** 仓库读取。本页是整个家族的 canonical cross-link。

## 家族概览

| 项目 | 数据库 | 状态 | 仓库 |
|---|---|---|---|
| **`postgres-operator`** | PostgreSQL 18+ | active | https://github.com/keiailab/postgres-operator |
| **`mongodb-operator`** | MongoDB 7.0+ | active | https://github.com/keiailab/mongodb-operator |
| **`valkey-operator`** | Valkey 8.0+ (Redis fork, BSD-3) | active | https://github.com/keiailab/valkey-operator |
| **`operator-commons`** | 共享 Go 库 | **v0.8.0** (当前页面) | https://github.com/keiailab/operator-commons |

## 我们共享什么

4 个项目都收敛到相同的运维 primitive:

- **Apache-2.0** end-to-end — 无 SSPL,SaaS 表面无 copyleft
- **`operator-commons`** 共享 Go 库 (v0.8.0+) — 终结器、标签、状态语法糖、security context 构建器、NetworkPolicy / ServiceMonitor partial
- **Helm chart 骨架** — RFC-0027 `default` falsy-toggle 防护、RFC-0026 组件 key 的 values、cycle 26 hardening 6 marker (priorityClassName / lifecycle / SA / minReadySeconds / automount / revisionHistoryLimit)
- **OLM bundle parity** — scorecard v1alpha3 6-test matrix
- **i18n** — README + canonical docs 提供 英文 / 한국어 / 日本語 / 中文 (cleanup supercycle 2026-05-21 的 Wave 4)

## `operator-commons` 在家族中的角色

本仓库是**共享 Go 库** — *不是*控制器. 提供:

| 包 | 用途 | Tier |
|---|---|---|
| `pkg/finalizer` | 仅使用 std `slices` 的 `Add` / `Remove` / `Has` 终结器辅助 | **Stable** |
| `pkg/labels` | 推荐 K8s 标签构建器 — `Set`, `All()`, `Selector()` | **Stable** |
| `pkg/status` | 4 个标准 Condition Type + 6 个 Reason catalog + 辅助 | **Stable** |
| `pkg/version` | DB 版本 allowlist convention + generic `Matrix[E MatrixEntry]` | Beta |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 构建器 (unstructured) | Beta |
| `pkg/networkpolicy` | Deny-by-default NetworkPolicy 构建器 + functional option | Beta |
| `pkg/security` | PodSecurity *restricted* SecurityContext 构建器 | Beta |
| `pkg/webhook` | Admission validation 辅助 | Experimental |

设计 invariant: **leaf 包只能使用 stdlib + k8s API 类型**. 无 controller-runtime,无 logr,无 operator-sdk 泄漏。

详细包 surface 见 [ARCHITECTURE.md](../ARCHITECTURE.md),tier 升级标准见 [ROADMAP.md](../ROADMAP.md)。

## 我们*不*做的事

- ❌ **嵌入或封装上游操作器** (PGO, CloudNativePG, MongoDB Community Operator, Sentinel) — license-clean,无 copyleft 义务
- ❌ **使用 GitHub Actions 作为发布门控** — 本地 4 层 + GitLab CI L5 (见 RFC-0002, RFC-0043)
- ❌ **基于时间的路线图截止日期** — 功能清单 + 完成百分比 (见 `standards/roadmap.md §1.1`)
- ❌ **Bitnami chart / image** — registry deprecation 风险,Broadcom 收购 (见 ADR-0136 / ADR-0057)
- ❌ **本仓库中的 CRD / Reconciler** — 由 consumer operator 拥有相关职责

## 从哪里开始

| 任务 | 入口 |
|---|---|
| 在操作器中导入 `operator-commons` | [README.md](../README.md) Usage 部分 |
| 阅读架构 | [ARCHITECTURE.md](../ARCHITECTURE.md) |
| 提交 issue 或功能请求 | https://github.com/keiailab/operator-commons/issues |
| 讨论设计或路线图 | https://github.com/keiailab/operator-commons/discussions |
| 贡献代码 | [CONTRIBUTING.md](../CONTRIBUTING.md) |
| 报告安全问题 | [SECURITY.md](../SECURITY.md) |
| 学习品牌 / 风格 | [BRANDING.md](../BRANDING.md) |
| 追踪采用者 / 谁在使用 | [ADOPTERS.md](../ADOPTERS.md) |
| 寻找维护者 | [MAINTAINERS.md](../MAINTAINERS.md) |
| 审查治理模型 | [GOVERNANCE.md](../GOVERNANCE.md) |
| 检查即将开展的工作 | [ROADMAP.md](../ROADMAP.md) |
| 审查 API 稳定性承诺 | [docs/STABILITY.md](STABILITY.md) |

## 跨家族兼容性

所有 3 个数据库操作器以匹配版本 (当前 `v0.8.0+`) 导入 `github.com/keiailab/operator-commons`:

```go
import (
    "github.com/keiailab/operator-commons/pkg/version"
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/labels"
    "github.com/keiailab/operator-commons/pkg/monitoring"
    "github.com/keiailab/operator-commons/pkg/finalizer"
    "github.com/keiailab/operator-commons/pkg/status"
)
```

`operator-commons` 中的 breaking change 需要 3 个数据库操作器同步 bump — supercycle Wave 5 的 `make cross-validation` 目标进行验证。

实时 consumer matrix (3 operator × 8 package × 采用率 %) 见 [ADOPTERS.md](../ADOPTERS.md)。

## i18n

本页面 (以及所有 canonical 项目文档) 提供 4 种语言:

- **English** (canonical, 原本文件)
- [한국어](family.ko.md)
- [日本語](family.ja.md)
- [中文](family.zh.md) (当前文件)

如有疑问,技术内容以英文版为权威;本地化版本以母语表达反映相同决策。

---

<p align="center">
  <b>keiailab 操作器家族</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
