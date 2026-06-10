# Security Policy

> [English](SECURITY.md) | [한국어](SECURITY.ko.md) | [日本語](SECURITY.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

`keiailab/keiailab-commons` 被下游 Kubernetes operators
导入。本库的漏洞会直接影响这些下游 consumer 的
运营安全。

## 漏洞报告

**请勿提交 public issue。**

### 渠道

使用下列私有渠道之一：

1. **GitHub Security Advisory**（首选）：
   `https://github.com/keiailab/keiailab-commons/security/advisories/new`
2. **Email**：`security@keiailab.com`（PGP 可选）：
   - PGP fingerprint：
     `89A4 0947 6828 CB99 2338  C378 651E 51AF 520B CB78`。

### 应包含内容

- 受影响版本（release tag 或 commit SHA）。
- 受影响包（`pkg/security`、`pkg/webhook` 等）。
- 重现步骤（如可能提供最小化重现；当重现依赖
  下游环境时请声明）。
- 影响评估 — 下游 consumer 的影响范围。
- 如果可用，自评的 CVSS 分数。

## 响应 SLA

| 阶段 | 时间 |
|---|---|
| 初始响应（确认） | 72 小时内 |
| 严重性评估 | 7 天内 |
| Patch release | 取决于严重性（Critical：14 天，High：30，Medium：60） |
| 公开披露 | patch 之后 14 天（或下游 consumer 能发布修复的最早时间点） |

## Embargo 处理

涉及 public API 的漏洞将受 embargo 直到下游
consumer 能够并行发布修复。Maintainer 在披露之前与下游
maintainer 共享私有 advisory。

## 支持的版本

| 版本 | 支持 |
|---------|-----------|
| 0.x（alpha） | ✅ 仅最新 minor |
| 1.0+（stable） | TBD — 在首个稳定 release 之后更新 |

本库当前处于 v0.x。Public API 可能 break；安全
patch 仅针对最新 minor 发布。

## 依赖安全

当添加或升级依赖时，PR 正文中需引用 license
和 CVE 审查。Dependabot / Renovate 自动更新 PR 会被
优先 review。

## License / 供应链

本库**仅 MIT**，charter 目标是零 AGPL /
BUSL 传递依赖（`docs/kb/adr/0001-charter.md`）。每次 minor release
均运行 license audit。

## 下游 consumer 最佳实践

导入本库的 operator 应该：

1. **使用 `pkg/security`** — 调用 restricted PodSecurity
   SecurityContext builder 而不是自行实现。
2. **使用 `pkg/webhook`** — 不要重新实现 version 校验。
3. **使用 `pkg/networkpolicy`** — deny-by-default NetworkPolicy builder。
4. 在 `go.mod` 中追踪
   `github.com/keiailab/keiailab-commons` 的最新 patch（Renovate
   自动 PR）。

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
