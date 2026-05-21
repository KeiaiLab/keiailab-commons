# Governance

> [English](GOVERNANCE.md) | [한국어](GOVERNANCE.ko.md) | [日本語](GOVERNANCE.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文档定义 `keiailab/operator-commons` 库的
决策流程。本库在多个下游 consumer operator 之间
被共同导入，因此对 *public API* 的变更影响
下游兼容性。

## 原则

1. **开放性** — 每个决策都在公开渠道（GitHub
   Issue / PR / ADR）中作出。
2. **Lazy Consensus** — 例行变更在没有明确异议的情况下
   继续推进。
3. **Explicit Consensus** — public-API breaking 变更、新包
   引入和 license 变更需要 maintainer 的 **2/3 supermajority**
   *外加*至少一位下游 consumer maintainer 的 LGTM。
4. **共享责任** — maintainer 共同承担库稳定性、
   下游运营影响和安全性的责任。

## 决策类别

### 例行变更（Lazy Consensus）

- bug 修复、文档改进、追加测试、minor /
  patch 依赖升级、不更改 public API 的
  内部 refactor。
- 流程：PR → 至少一位 maintainer LGTM → merge。
- 窗口：无 — 当本地 gate 通过即可 merge。本
  项目不使用 GitHub Actions；所有质量 gate 均由
  本地四层（`lefthook.yml`、`Makefile`、reviewer evidence、
  ADR coverage）强制。

### 中等变更（Explicit Consensus）

- 添加新的 public-API 函数或类型、major 依赖升级、
  引入新的 `pkg/<sub>` 包。
- 流程：issue 或 ADR 提案 → 7 天评论窗口 → 多数
  maintainer LGTM → merge。
- 一个或多个异议意见触发 maintainer 讨论。

### Public-API breaking 变更（需要 ADR）

- 函数 signature 变更、类型删除、module-path 变更、license
  变更。
- 流程：
  1. 提交 `docs/kb/adr/NNNN-<slug>.md`。
  2. 14 天评论窗口。
  3. 2/3 maintainer supermajority 加上至少一位下游 consumer
     LGTM。
  4. 将 ADR `Status: Draft → Accepted` 并 merge
     实现 PR。

## 安全决策

CVE 报告遵循 [SECURITY.md](../SECURITY.zh.md)。报告以私有方式
处理。Embargo 一直保持到下游 consumer 能够发布
修复为止。

## Release 决策

- **v0.x**：单一 maintainer 可在 Lazy Consensus 下标记
  minor / patch release。
- **v1.0+ (stable)**：严格 SemVer — major bump 需要 ADR 加上
  2/3 supermajority。

## 变更历史

| Date | Change |
|---|---|
| 2026-05-09 | 文档引入 — governance baseline。 |

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
