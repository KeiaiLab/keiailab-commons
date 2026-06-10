# Maintainers

> [English](MAINTAINERS.md) | [한국어](MAINTAINERS.ko.md) | [日本語](MAINTAINERS.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文档记录对 `keiailab/keiailab-commons` 拥有
决策权限的 maintainer。

## 当前 maintainer

| 姓名 / 团队 | GitHub | 角色 | 范围 |
|---|---|---|---|
| keiailab maintainers | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | Lead | All |

GitHub team `@keiailab/maintainers` 拥有本库每个领域的
merge 和 release-tag 权限。

## Maintainer 资格

要么是下游 consumer operator 的 maintainer，*要么*是
满足下列条件至少六个月的贡献者：

- ≥ 10 个 merged PR（library PR cadence 低于典型 operator，
  因此门槛约为其一半）。
- ≥ 20 个 reviewed PR（下游 consumer PR 可计入）。
- 至少熟悉一个 `pkg/` 领域（security、labels、
  webhook、monitoring、networkpolicy、version、status、finalizer、
  storageclass、events、probes、pvc、topology）。

## 添加流程

1. 现有 maintainer（或候选人本人）开启 issue
   或 ADR。
2. `@keiailab/maintainers` team 应用 lazy consensus（7 天
   评论窗口）。
3. 无异议时，候选人被添加到 GitHub team 并通过 PR
   更新本文件。

## 不活跃 maintainer

连续六个月不活跃的 maintainer 将
移至 emeritus（权限撤销，姓名保留在荣誉名册）。

## 跨仓库协议

*public-API breaking 变更*在 ADR 阶段需要下游 consumer
maintainer 的 LGTM — 参见 [GOVERNANCE.md](GOVERNANCE.zh.md)。

## i18n 文档负责人

| Language | Owner | Files | Responsibility |
|---|---|---|---|
| English (canonical) | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | `README.md` 和 canonical docs | Source of truth |
| Korean | TaeHwan Park ([@eightynine01](https://github.com/eightynine01)) | `README.ko.md` 和 `*.ko.md` | EN canonical sync + translation review |
| Japanese | (recruiting — volunteer via an issue) | `*.ja.md` | AI translation + native review |
| Chinese | (recruiting — volunteer via an issue) | `*.zh.md` | AI translation + native review |

**Drift 验证**：`bash scripts/check-readme-sync.sh` — 检查
文件是否存在、section header count 是否匹配、line count
之差是否小于每种语言的阈值，以及 cross-link 是否
双向。lefthook `pre-push` hook `readme-i18n-sync` 自动
强制此项。

## Emeritus

(none yet)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
