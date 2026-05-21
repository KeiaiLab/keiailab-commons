# 品牌指南 — `operator-commons`

> ⚠️ This translation is AI-generated and pending native review. — 本翻译为 Claude 机器翻译结果。

> keiailab 操作器家族的 visual identity, voice, tone.

本文档是 `operator-commons` 品牌决策的 canonical reference. 适用于 README, release note, 营销材料和任何代表项目的第三方沟通。

## 1. Identity

**Organization**: [keiailab](https://keiailab.com) — Kubernetes-native 数据平台操作器 (Apache-2.0, license-clean, vanilla-upstream 兼容).

**Project**: `operator-commons` — keiailab 操作器的共享 Go 库 — finalizer / labels / status / version / security / monitoring partial.

**Family**: 共享 [`operator-commons`](https://github.com/keiailab/operator-commons) 共享库的 4 个姊妹操作器之一:

| 项目 | 数据库 | 仓库 |
|---|---|---|
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `mongodb-operator` | MongoDB 7.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ (Redis fork, BSD-3) | https://github.com/keiailab/valkey-operator |
| `operator-commons` | 共享 Go 库 | https://github.com/keiailab/operator-commons |

## 2. 标志 & 视觉资源

| 资源 | URL | 用途 |
|---|---|---|
| Primary 标志 (SVG) | `https://keiailab.com/assets/logo.svg` | README header, 幻灯片 |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon, social card |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer, dark background |

**标志放置**: README 顶部中央, width 120px. 始终链接到 https://keiailab.com.

**Clear space**: 标志周围最小 padding = 标志 width 的 25%.

**禁止**:
- 标志颜色变更
- 添加 drop shadow / filter
- 放置在对比度不足的背景
- 未经 keiailab 品牌批准与其他标志组合

## 3. 色彩调色板

| 角色 | Hex | 用途 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | 标题, primary action, 链接 |
| Secondary (deep navy) | `#0F172A` | dark 背景, 代码块 |
| Accent (warm amber) | `#F59E0B` | 强调, 徽章 accent |
| Neutral grey | `#64748B` | light 背景的 body text |
| Background light | `#F8FAFC` | 文档页面背景 |
| Background dark | `#020617` | 代码编辑器主题, dark mode |

GitHub README 的 shield.io badge 推荐使用上述 hex.

## 4. 字体

- **Heading**: 系统默认 (GitHub 的 default `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: 相同 (GitHub-native 一致)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...`

不使用单独 webfont (GitHub README rendering 一致)。

## 5. Voice & Tone

**Audience**: Kubernetes 平台工程师 / DBA / SRE / Go 库 consumer.

**Voice 原则**:
- **Direct** — 可能时使用 bullet point 替代段落
- **Evidence-based** — claim 包含 benchmark / SLA / link
- **Library-focused** — `operator-commons` 是*库* — controller-runtime / CRD / reconciler 是 *consumer operator 的责任*
- **License-aware** — 仅 Apache-2.0, AGPL/BUSL transitive 依赖 0 件目标 (ADR-0001 charter)

**Avoid**:
- 营销最高级 ("blazing fast", "revolutionary", "best-in-class")
- 模糊比较 ("X-class quality") — *用具体的 metric 或 benchmark 限定*
- 路线图中基于时间的截止日期 (`standards/roadmap.md §1.1` — 用功能清单替代)

## 6. README Header 标准

所有 README 的第一段使用以下格式 (Wave 3 标准):

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab 操作器的共享 Go 库 — finalizer / labels / status / version / security / monitoring partial**

<p align="center">
  <a href="../LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
</p>

<p align="center">
  <a href="../README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <b>中文</b>
</p>
```

## 7. README Footer 标准

所有 README + root-level .md 文件末尾附加以下 footer (Wave 3 标准):

```markdown
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
```

## 8. Badge 标准顺序

README 的 shield.io badge 顺序 (左→右):

1. License (Apache-2.0)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `operator-commons` 是*库*, 所以 Kubernetes / Container Image / Helm Chart badge 放在 *consumer operator README 中*. 本库以 `pkg.go.dev` + `OpenSSF Scorecard` 为中心.

## 9. Discussions / Issues / PR 模板

- **Discussions**: `https://github.com/keiailab/operator-commons/discussions` — pkg API 提问, integration 案例, 新 helper 提案
- **Issues**: bug report + 有 use case 的具体 feature request (推荐明示对 consumer operator 端的影响)
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` 标准 (用户场景 + 验证命令引用义务, `standards/checklist.md §3`)

## 10. Social & External

- **Website**: https://keiailab.com
- **GitHub Org**: https://github.com/keiailab
- **pkg.go.dev**: https://pkg.go.dev/github.com/keiailab/operator-commons

## 11. License & Attribution

- License: [Apache-2.0](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: 参见 [NOTICE](../NOTICE) (如适用)

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
