# 品牌指南 — `operator-commons`

> [English](BRANDING.md) | [한국어](BRANDING.ko.md) | [日本語](BRANDING.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.
>
> `operator-commons` 库的 visual identity、voice、tone。

本文档是 `operator-commons` 品牌决策的 canonical reference。适用于 README、
发布说明、以及关于本项目的外部沟通。

## 1. Identity

**Organization**: [keiailab](https://keiailab.com)。

**Project**: `operator-commons` —— 用于 Kubernetes operator 通用 scaffolding
(finalizer / labels / status / version / security / monitoring partial) 的
Go 库。

本库以 Go 模块 `github.com/keiailab/operator-commons` 与 Helm library chart
(`charts/keiailab-commons`) 形式发布。downstream operator 通过标准 Go module
import 使用 —— 这里不指名或背书任何具体 consumer。

## 2. 标志与视觉资源

| 资源 | URL | 用途 |
|---|---|---|
| Primary 标志 (SVG) | `https://keiailab.com/assets/logo.svg` | README header、幻灯片 |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon、social card |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer、深色背景 |

**标志位置**: README 顶部居中,width 120 px。始终链接到
`https://keiailab.com`。

**Clear space**: 标志周围的最小 padding 等于标志 width 的 25 %。

**禁止**:

- 修改标志颜色
- 添加 drop shadow / filter
- 放置在对比度不足的背景上
- 未经 keiailab 品牌批准与其他标志组合

## 3. 色彩调色板

| 角色 | Hex | 用途 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | 标题、primary action、链接 |
| Secondary (deep navy) | `#0F172A` | 深色背景、代码块 |
| Accent (warm amber) | `#F59E0B` | 强调、徽章 accent |
| Neutral grey | `#64748B` | 浅色背景上的 body text |
| Background light | `#F8FAFC` | 文档页面背景 |
| Background dark | `#020617` | 深色模式代码编辑器主题 |

GitHub README shield.io 徽章使用相同的 hex 值。

## 4. 字体

- **Heading**: 系统默认 (GitHub 默认 `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: 系统默认 (GitHub-native 一致)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...` (GitHub 默认 monospace)

不使用外部 web 字体 —— 保持与 GitHub 原生 rendering 完全一致。

## 5. Voice and Tone

**Audience**: Kubernetes 平台工程师、DBA、SRE、Go 库 consumer。

**Voice 原则**:

- **Direct** — 尽可能使用 bullet point 而非段落。
- **Evidence-based** — claim 包含 benchmark、SLA 或 link。
- **Library-focused** — `operator-commons` 是 *库*。controller-runtime、CRD、
  reconciler 是 downstream consumer 的责任,而非本库的责任。
- **License-aware** — 仅 Apache-2.0。charter 目标为 AGPL / BUSL transitive
  依赖 0 件 (`docs/kb/adr/0001-charter.md`)。

**Avoid**:

- 营销最高级 ("blazing fast"、"revolutionary"、"best-in-class")。
- 模糊比较 ("enterprise-grade quality") —— 用具体的 metric 或 benchmark 限定。
- 路线图中基于时间的截止日期 —— 使用 [ROADMAP.md](ROADMAP.md) 的功能清单。

## 6. README Header 标准

每个 README 的第一块遵循以下格式:

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **用于 Kubernetes operator 通用 scaffolding 的 Go 共享库 — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <!-- 其他 shield.io 徽章 -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <b>中文</b>
</p>
```

## 7. README Footer 标准

每个 README 与根级 `.md` 文件以单行 attribution 结尾:

```markdown
---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
```

不添加额外 cross-link 区块。footer 保持最小化,使文档 self-contained。

## 8. 徽章顺序

README 的 shield.io 徽章按以下顺序 (左→右):

1. License (Apache-2.0)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `operator-commons` 是 *库*,因此 container image、Helm chart、
> Kubernetes deployment 徽章不放在本库 —— 放在出货 image 或 chart 的
> downstream operator README 上。

## 9. Discussions / Issues / PR Template

- **Discussions**: `https://github.com/keiailab/operator-commons/discussions` — 包 API 提问、integration 案例、新 helper 提案。
- **Issues**: bug 报告与包含 use case 的具体 feature request。相关时明示 downstream consumer 影响。
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits + 用户场景 + 验证命令输出引用。

## 10. Social and External

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/operator-commons>

## 11. License and Attribution

- License: [Apache-2.0](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: 参见 [NOTICE](../NOTICE)

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
