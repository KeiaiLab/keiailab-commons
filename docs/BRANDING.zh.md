# 品牌指南 — `keiailab-commons`

> [English](BRANDING.md) | [한국어](BRANDING.ko.md) | [日本語](BRANDING.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.
>
> `keiailab-commons` 库的 visual identity、voice、tone。

本文档是 `keiailab-commons` 品牌决策的 canonical reference。适用于 README、
发布说明、以及关于本项目的外部沟通。

## 1. Identity

**Organization**: [keiailab](https://keiailab.com)。

**Project**: `keiailab-commons` —— 用于 Kubernetes operator 通用 scaffolding
(finalizer / labels / status / version / security / monitoring partial) 的
Go 库。

本库以 Go 模块 `github.com/keiailab/keiailab-commons` 与 Helm library chart
(`charts/keiailab-commons`) 形式发布。downstream operator 通过标准 Go module
import 使用 —— 这里不指名或背书任何具体 consumer。

### 1.1 家族 charter (本文档为 SSOT)

本 Branding Guide 是 keiailab operator 家族共享视觉文法的 **单一事实来源
(SSOT)**。调色板 (§3)、环形文法与字形替换流程 (§2b)、各 repo 字形注册表 (§2b)、
README 页眉 / footer 模板 (§6 / §7)、徽章顺序 (§8) 均在 **此处** 定义,由各家族
仓库复制而非重新推导。`keiailab-commons` 拥有该 charter,因为它是每个 operator
都已 import 的共享依赖。

家族由 4 个 sister operator 加本共享库组成:

| Project | Focus | Repository |
|---|---|---|
| `mongodb-operator` | MongoDB 8.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ | https://github.com/keiailab/valkey-operator |
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `qdrant-operator` | Qdrant vector database | https://github.com/keiailab/qdrant-operator |
| `keiailab-commons` | Shared Go library | https://github.com/keiailab/keiailab-commons |

当本指南的翻译滞后时,此处的英文文本保持 canonical。

## 2. 标志与视觉资源

| 资源 | URL | 用途 |
|---|---|---|
| Primary 标志 | `docs/branding/symbol.png` | README header、幻灯片 |
| Keiailab base symbol | `docs/branding/base-symbol.png` | 外圈旋转箭头标记的源参考 |
| Favicon | `https://keiailab.com/favicon.ico` | Favicon、social card |
| 计划 SVG kit | `https://keiailab.com/assets/{logo,mark,wordmark}.svg` | URL 返回 200 后替换 |

**标志位置**: README 顶部居中,width 96 px。始终链接到
`https://keiailab.com`。

**Clear space**: 标志周围的最小 padding 等于标志 width 的 25 %。

**禁止**:

- 修改标志颜色
- 添加 drop shadow / filter
- 放置在对比度不足的背景上
- 未经 keiailab 品牌批准与其他标志组合

## 2b. 环形文法与字形替换流程

家族共享一个外圈符号 —— `docs/branding/base-symbol.png` 中的三色旋转箭头环
(460×460; 深海军蓝 → 蓝 → 绿箭头,中心球体)。**每个** 家族仓库都按 byte 逐一
vendor 此文件,环永不重新着色、重绘或重新生成。

各项目自己的 `docs/branding/symbol.png` 通过保留环并 **仅替换中心字形** 生成:

1. 从共享 `base-symbol.png` 开始 —— 不要触碰环。
2. 将产品字形 alpha-composite 到中心 clear zone (~185 px 直径); 不与环重叠。
3. 以源分辨率 export 到 `docs/branding/symbol.png`。
4. 在 README 页眉以 width 96 px 引用 (§6),operator 还作为 Helm chart `icon`
   引用。

保留共享环并仅替换字形,正是让这些标记读起来是一个家族、同时保持各产品可识别的
关键。重新着色环或绘制新的外圈标记是禁止的。

### 字形注册表

| Repository | Center glyph |
|---|---|
| `mongodb-operator` | leaf |
| `valkey-operator` | hexagon + keyhole |
| `postgres-operator` | elephant |
| `qdrant-operator` | kNN constellation |
| `keiailab-commons` | node grid |

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
- **Library-focused** — `keiailab-commons` 是 *库*。controller-runtime、CRD、
  reconciler 是 downstream consumer 的责任,而非本库的责任。
- **License-aware** — 仅 MIT。charter 目标为 AGPL / BUSL transitive
  依赖 0 件 (`docs/kb/adr/0001-charter.md`)。

**Avoid**:

- 营销最高级 ("blazing fast"、"revolutionary"、"best-in-class")。
- 模糊比较 ("enterprise-grade quality") —— 用具体的 metric 或 benchmark 限定。
- 路线图中基于时间的截止日期 —— 使用 [ROADMAP.md](ROADMAP.md) 的功能清单。

## 6. README Header 标准

每个 README 的第一块遵循以下格式:

```markdown
<p align="center">
  <img src="docs/branding/symbol.png" alt="keiailab" width="96"/>
</p>

# keiailab-commons

> **用于 Kubernetes operator 通用 scaffolding 的 Go 共享库 — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/<license placeholder>-blue.svg" alt="License"/></a>
  <!-- 按 §8 顺序的其他 shield.io 徽章 -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <b>中文</b>
</p>
```

> **模板 placeholder**: 将 `<license placeholder>` 替换为复制仓库的实际 SPDX
> 徽章片段 —— 例如 MIT 仓库 (家族大多数) 用 `License-MIT`,许可证不同的成员用
> 对应的 `License-<SPDX-id>`。在共享模板中硬编码单一许可证会给许可证不同的成员
> 再生产错误徽章,因此 token 在此保持 generic,由各仓库填入。

## 7. README Footer 标准

每个 README 与根级 `.md` 文件以单行 attribution 结尾:

```markdown
---

<p align="center">© 2026 keiailab · <license placeholder> · <a href="https://keiailab.com">keiailab.com</a></p>
```

> **模板 placeholder**: 将 `<license placeholder>` 替换为复制仓库的实际
> 许可证名称 —— 与 §6 徽章 token 相同的按仓库替换规则。`keiailab-commons`
> 自身的实际值: MIT。

不添加额外 cross-link 区块。footer 保持最小化,使文档 self-contained。

## 8. 徽章顺序

徽章分为 **MUST** (始终显示) 与 **SHOULD** (仅当后端基础设施上线时显示)。当
工作流或注册表尚不存在却强制 8 个时,徽章会长期显示 red,因此家族标准是诚实的
4 + 4。

**MUST (4)** — 每个家族仓库,左→右:

1. License
2. Go Version
3. Product (MongoDB / Valkey / PostgreSQL / Qdrant)
4. Kubernetes Version

**SHOULD (4)** — 仅当其后端表面处于活跃时各自添加:

5. Container image (`ghcr.io/keiailab/<repo>`)
6. Helm chart (Artifact Hub)
7. OpenSSF Scorecard
8. GitHub Discussions

> **库例外**: `keiailab-commons` 既不出货 container image、Helm application
> chart,也不出货 Kubernetes 工作负载,因此其 MUST 集用 License / Go Version /
> **Go Reference** (pkg.go.dev) 代替 product 与 Kubernetes 徽章,SHOULD 集为
> OpenSSF Scorecard + GitHub Discussions。container / Helm / Kubernetes 徽章放在
> 出货 image 或 chart 的 downstream operator 上。

## 9. Discussions / Issues / PR Template

- **Discussions**: `https://github.com/keiailab/keiailab-commons/discussions` — 包 API 提问、integration 案例、新 helper 提案。
- **Issues**: bug 报告与包含 use case 的具体 feature request。相关时明示 downstream consumer 影响。
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits + 用户场景 + 验证命令输出引用。

## 10. Social and External

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/keiailab-commons>

## 11. License and Attribution

- License: [MIT](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: 参见 [NOTICE](../NOTICE)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
