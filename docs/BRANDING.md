# Branding Guide — `operator-commons`

> **English** | [한국어](BRANDING.ko.md) | [日本語](BRANDING.ja.md) | [中文](BRANDING.zh.md)

> Visual identity, voice, and tone for the `operator-commons` library.

This document is the canonical reference for `operator-commons` branding
decisions. It applies to the README, release notes, and any external
communication about the project.

## 1. Identity

**Organization**: [keiailab](https://keiailab.com).

**Project**: `operator-commons` — a shared Go library for Kubernetes operator
scaffolding (finalizer / labels / status / version / security / monitoring
partials).

The library is published as a Go module at
`github.com/keiailab/operator-commons` and as a Helm library chart
(`charts/keiailab-commons`). It is consumed by downstream operator
implementations via standard Go module import — no specific consumer is
named or endorsed here.

## 2. Logo and visual assets

| Asset | URL | Usage |
|---|---|---|
| Primary logo (SVG) | `https://keiailab.com/assets/logo.svg` | README header, slide decks |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon, social cards |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer, dark backgrounds |

**Logo placement**: Top-center of README, width 120 px. Always link to
`https://keiailab.com`.

**Clear space**: Minimum padding around the logo equals 25 % of the logo
width.

**Do not**:

- Recolor the logo
- Add drop shadows or filters
- Place the logo on backgrounds with insufficient contrast
- Combine with other logos without keiailab brand approval

## 3. Color palette

| Role | Hex | Usage |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | Headers, primary actions, links |
| Secondary (deep navy) | `#0F172A` | Dark backgrounds, code blocks |
| Accent (warm amber) | `#F59E0B` | Highlights, badge accents |
| Neutral grey | `#64748B` | Body text on light backgrounds |
| Background light | `#F8FAFC` | Documentation page background |
| Background dark | `#020617` | Dark-mode code editor theme |

GitHub README shield.io badges should use the same hex values.

## 4. Typography

- **Headings**: System default (GitHub default: `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: System default (GitHub-native consistency)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...` (GitHub default monospace)

No external web font is used — keep rendering identical to native GitHub.

## 5. Voice and tone

**Audience**: Kubernetes platform engineers, DBAs, SREs, and Go library
consumers.

**Voice principles**:

- **Direct** — prefer bullet points over paragraphs where possible.
- **Evidence-based** — claims include a benchmark, an SLA, or a link.
- **Library-focused** — `operator-commons` is a *library*. Controller-runtime,
  CRDs, and reconcilers are the responsibility of the downstream consumer,
  not of this library.
- **License-aware** — Apache-2.0 only. The charter goal is zero AGPL / BUSL
  transitive dependencies (see `docs/kb/adr/0001-charter.md`).

**Avoid**:

- Marketing superlatives ("blazing fast", "revolutionary", "best-in-class").
- Vague comparisons ("enterprise-grade quality") — qualify each claim with a
  specific metric or benchmark.
- Time-based deadlines in the roadmap — use the feature checklist in
  [ROADMAP.md](ROADMAP.md).

## 6. README header standard

Every README's first block follows this layout:

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **Shared Go library for Kubernetes operator scaffolding — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <!-- additional shield.io badges -->
</p>

<p align="center">
  <b>English</b> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>
```

## 7. README footer standard

Every README and root-level `.md` file ends with a single attribution line:

```markdown
---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
```

No additional cross-link block. Keep the footer minimal so the document
remains self-contained.

## 8. Badge order

The shield.io badges in the README appear in this order (left → right):

1. License (Apache-2.0)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `operator-commons` is a *library*, so container image, Helm
> chart, or Kubernetes deployment badges do not belong here — they belong on
> the README of a downstream operator that ships an image or chart.

## 9. Discussions, issues, PR templates

- **Discussions**: `https://github.com/keiailab/operator-commons/discussions` — package API questions, integration patterns, new helper proposals.
- **Issues**: bug reports plus concrete feature requests; include the use case and downstream consumer impact when relevant.
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits + user-facing scenario + verification command output.

## 10. Social and external

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/operator-commons>

## 11. License and attribution

- License: [Apache-2.0](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attributions: see [NOTICE](../NOTICE)

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
