# Contributing to operator-commons

> [English](CONTRIBUTING.md) | [한국어](CONTRIBUTING.ko.md) | [日本語](CONTRIBUTING.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

`keiailab/operator-commons` 是被下游 Kubernetes operators 导入的
Go library。所有贡献遵循
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) 和
[docs/GOVERNANCE.md](docs/GOVERNANCE.zh.md)。

## 贡献流程

1. **通过 issue 或 ADR 表明意图**（针对非琐碎变更）。
2. **Fork + feature branch** — 使用 `feat/<slug>`、`fix/<slug>`、
   `docs/<slug>` 或 `refactor/<slug>`。
3. **验证本地 gate**：
   - `lefthook install --force`
   - `lefthook run pre-commit --all-files`
   - `make lint test`
4. **开启 PR** — Conventional Commits 格式；正文可使用英文或韩文。
5. **Review SLA**：maintainer 在 24 小时内回复。

## PR 检查清单（作者）

- [ ] PR 标题：Conventional Commits 格式（`feat`、`fix`、`docs`、
  `refactor`、`test`、`chore`）。
- [ ] PR 正文：变更摘要 + 验证命令 + 引用输出。
- [ ] 添加或更新单元测试（`pkg/<sub>` 内任何变更均强制）。
- [ ] 当 public API 变更时更新 godoc。
- [ ] 对于 **public-API breaking changes**：链接 ADR 以及
  下游 consumer 影响分析。
- [ ] `go.mod` / `go.sum` drift = 0（运行 `go mod tidy` 无变化）。
- [ ] 新依赖：在 PR 正文中引用 license 和 CVE 审查。

## 本地开发（与下游 consumer 的跨切关联工作）

当变更同时涉及 `operator-commons` 和下游 operator 时：

```fish
# 1. 在 consumer operator 的 go.mod 中添加 replace directive
#    （仅本地；不要提交它）
# go.mod 尾部：
#   replace github.com/keiailab/operator-commons => ../operator-commons

# 2. 双方均编辑 + 各自运行 `go test ./...`

# 3. 拆分 PR：
#    - operator-commons 侧：merge + tag（例如 v0.9.0）
#    - consumer 侧：bump require directive（删除 replace）
```

## 添加新的 `pkg/<sub>` 包

遵循
[docs/GOVERNANCE.md](docs/GOVERNANCE.zh.md)中的"中等变更"流程：

1. 开启 issue 或 ADR，解释*为什么*它属于 commons 以及
   *哪个*下游 consumer 将使用它。
2. 等待 7 天评论窗口。
3. 在多个 maintainer LGTM 后合并。

## Release

- v0.x SemVer：每次 minor bump 表示 public-API 变更或
  有意义的行为变更。
- Release 流程：
  1. `git tag v0.X.Y`（annotated）。
  2. `git-cliff` 重新生成 CHANGELOG PR。
  3. `git push origin v0.X.Y`。
  4. 在每个下游 consumer 中开启后续 PR 以 bump
     `require` directive。

## 安全漏洞

使用 [SECURITY.md](SECURITY.zh.md) 中的私有披露流程。公开
issue 不是合适的渠道。

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
