# i18n — keiailab-commons 多语言政策

> [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文档定义 `keiailab-commons` 文档如何在多种
语言中维护。本项目保持英文 canonical 加上
韩文 / 日文 / 中文翻译。

## §1 政策

### 1.1 核心原则

- **Canonical 语言**：英文（所有源的 ground truth）。
- **翻译语言（三种额外）**：Korean (`ko`)、Japanese
  (`ja`)、Chinese (`zh`，Simplified)。
- **四语言骨架**：`README.md`（英文）+ `README.ko.md` +
  `README.ja.md` + `README.zh.md`。
- **不接受这四种之外的语言** — 添加 `es` / `fr` / `de` / …
  需要单独的 ADR。

### 1.2 范围内 vs. 范围外

**范围内**（面向用户或外部文档）：

- `README.{md, ko, ja, zh}.md`
- 顶层 governance / branding：`BRANDING.md`、`STABILITY.md`、
  `coverage-report.md` 等四语言矩阵。
- `docs/i18n/glossary-{ko, ja, zh}.md`（本 SSOT）。

**范围外**（内部、治理或法律）：

- `docs/kb/adr/*.md`（决策记录 — 仅英文 canonical）。
- `docs/kb/deps/*.md`（依赖审计 — 生成的）。
- `LICENSE`、`NOTICE`、`CITATION.cff`、`.gitignore`（法律或配置）。

## §2 术语表

### 2.1 位置

SSOT：[`glossary-ko.md`](glossary-ko.md) / [`glossary-ja.md`](glossary-ja.md)
/ [`glossary-zh.md`](glossary-zh.md)。

### 2.2 优先级

1. 术语表定义优先。
2. 官方 Kubernetes 韩文 / 日文 / 中文指南。
3. 通用词典。
4. **代码标识符保持英文**（不翻译）— `pkg/probes`、
   `Reconciler`、`kubectl` 等。

### 2.3 一致性规则（所有四种语言）

- 不要在单一文档中混合 formal 与 informal register。
- 面向用户的文档使用 formal register。
- 内部文档（AGENTS）可以使用 informal register。
- 在首次出现时，包含英文原文加上
  带括号的翻译；后续出现可以仅使用
  翻译。

## §3 自动翻译 SOP

### 3.1 引擎

**默认**：Claude direct（AI subagent 读取源并
翻译）。

理由：韩文质量和 prompt caching 成本。日文和
中文翻译之后由 native review 跟进。

### 3.2 命令

```bash
# 手动单文件翻译
# subagent 读取源、翻译并写入输出
# 带有 AI 翻译警告横幅。

# 自动化（脚本化）
./scripts/i18n-translate.sh <source.md> --lang all --engine claude

# Dry run
./scripts/i18n-translate.sh README.md --dry-run
```

### 3.3 结果 marker

每个 AI 翻译的文件都必须包含警告横幅：

```markdown
> ⚠️ This translation is AI-generated and pending native review.
```

附加 marker：`[needs review]`（review 待定）、`[reviewed]`
（由 native reviewer 通过 PR 提升）。

## §4 Native review SOP

### 4.1 Reviewer 招募

默认政策：AI 翻译 + `[needs review]` marker + 警告
横幅。在招募到 native reviewer 之前 placeholder 保留。

### 4.2 Promotion 标准

将 `[needs review]` → `[reviewed]` 提升时，native reviewer
确认：

1. **Idiomatic vs. literal 平衡** — 先求自然阅读，
   但技术术语保持精确。
2. **术语表一致性** — 每个术语表词条应用
   一致。
3. **Register 一致性** — 文档内 formal / informal register
   统一。
4. **链接完整性** — Cross-link 按语言解析。
5. **代码标识符不翻译**。

### 4.3 Promotion PR 格式

```
title: docs(i18n): <lang> <doc> native reviewer promotion
body:
  ## Reviewer
  - <name / handle>

  ## Scope reviewed
  - <file path>

  ## Changes
  - [needs review] → [reviewed] marker change
  - Warning banner removed (or kept as partial)
  - <other idiomatic edits>

  ## Verification
  - [x] Glossary consistency
  - [x] Register consistency
  - [x] Cross-link integrity
```

## §5 Drift 控制

### 5.1 lefthook hook

`lefthook.yml` 在 pre-push 时强制 `readme-i18n-sync`：

```yaml
readme-i18n-sync:
  glob: "README*.md"
  run: bash scripts/check-readme-sync.sh
```

pre-push 位置使工作流摩擦保持较低。

### 5.2 阈值

| Lang | Line-diff 阈值 | 原因 |
|---|---|---|
| ko | ≤ 15 % | 韩文与英文 LOC 大致相同。 |
| ja | ≤ 25 % | 日文混合假名与汉字 — 比英文短约 15–20 %。 |
| zh | ≤ 30 % | 中文完全是汉字 — 比英文短约 25 %。 |

这些是推荐值；实证测量后可以调整。

### 5.3 4-lang 矩阵

`scripts/check-readme-sync.sh` 自动检查四语言
矩阵（EN ↔ {ko, ja, zh}）：

- 目标语言文件不存在 → 跳过（处理部分骨架）。
- 提供每种语言的 bypass
  （`SKIP_CHECK_README_SYNC_JA=1` 等）。

## §6 参考

- 术语表：[`glossary-ko.md`](glossary-ko.md) /
  [`glossary-ja.md`](glossary-ja.md) / [`glossary-zh.md`](glossary-zh.md)。
- 检查脚本：[`../../scripts/check-readme-sync.sh`](../../scripts/check-readme-sync.sh)。
- 翻译脚本：[`../../scripts/i18n-translate.sh`](../../scripts/i18n-translate.sh)。
