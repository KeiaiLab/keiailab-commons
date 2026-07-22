# Branding Guide — `keiailab-commons`

> **English** | [한국어](BRANDING.ko.md) | [日本語](BRANDING.ja.md) | [中文](BRANDING.zh.md)

> Visual identity, voice, and tone for the `keiailab-commons` library.

This document is the canonical reference for `keiailab-commons` branding
decisions. It applies to the README, release notes, and any external
communication about the project.

## 1. Identity

**Organization**: [keiailab](https://keiailab.com).

**Project**: `keiailab-commons` — a shared Go library for Kubernetes operator
scaffolding (finalizer / labels / status / version / security / monitoring
partials).

The library is published as a Go module at
`github.com/keiailab/keiailab-commons` and as a Helm library chart
(`charts/keiailab-commons`). It is consumed by downstream operator
implementations via standard Go module import — no specific consumer is
named or endorsed here.

### 1.1 Family charter (this document is the SSOT)

This Branding Guide is the **single source of truth** for the visual grammar
shared by the keiailab operator family. The color palette (§3), the ring
grammar and glyph-swap procedure (§2b), the per-repo glyph registry (§2b), the
README header and footer templates (§6 / §7), and the badge order (§8) are
defined **here** and copied — not re-derived — by each family repository.
`keiailab-commons` owns the charter because it is the shared dependency every
operator already imports.

The family is four sister operators plus this shared library:

| Project | Focus | Repository |
|---|---|---|
| `mongodb-operator` | MongoDB 8.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ | https://github.com/keiailab/valkey-operator |
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `qdrant-operator` | Qdrant vector database | https://github.com/keiailab/qdrant-operator |
| `keiailab-commons` | Shared Go library | https://github.com/keiailab/keiailab-commons |

When a translation of this guide lags, the English text here stays canonical.

## 2. Logo and visual assets

| Asset | URL | Usage |
|---|---|---|
| Current primary logo | `docs/branding/symbol.png` | README header, slide decks |
| Keiailab base symbol | `docs/branding/base-symbol.png` | Source reference for the outer rotating-arrow mark |
| Current favicon | `https://keiailab.com/favicon.ico` | Favicon, social cards |
| Planned SVG kit | `https://keiailab.com/assets/{logo,mark,wordmark}.svg` | Future replacement after URLs return 200 |

**Logo placement**: Top-center of README, width 96 px. Always link to
`https://keiailab.com`.

**Clear space**: Minimum padding around the logo equals 25 % of the logo
width.

**Do not**:

- Recolor the logo
- Add drop shadows or filters
- Place the logo on backgrounds with insufficient contrast
- Combine with other logos without keiailab brand approval

## 2b. Ring grammar & glyph-swap procedure

The family shares one outer symbol — the three-color rotating-arrow ring in
`docs/branding/base-symbol.png` (460×460; dark-navy → blue → green arrows around
a center sphere). **Every** family repository vendors this exact file
byte-for-byte; the ring is never recolored, redrawn, or regenerated.

A project's own `docs/branding/symbol.png` is produced by keeping the ring and
swapping **only the center glyph**:

1. Start from the shared `base-symbol.png` — do not touch the ring.
2. Alpha-composite the product glyph into the central clear zone
   (~185 px diameter); nothing overlaps the ring.
3. Export to `docs/branding/symbol.png` at the source resolution.
4. Reference it from the README header at width 96 px (§6) and, for an operator,
   as the Helm chart `icon`.

Keeping a shared ring and swapping only the glyph is what makes the marks read
as one family while staying per-product recognizable. Recoloring the ring or
drawing a new outer mark is prohibited.

### Glyph registry

| Repository | Center glyph |
|---|---|
| `mongodb-operator` | leaf |
| `valkey-operator` | hexagon + keyhole |
| `postgres-operator` | elephant |
| `qdrant-operator` | kNN constellation |
| `keiailab-commons` | node grid |

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
- **Library-focused** — `keiailab-commons` is a *library*. Controller-runtime,
  CRDs, and reconcilers are the responsibility of the downstream consumer,
  not of this library.
- **License-aware** — MIT only. The charter goal is zero AGPL / BUSL
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
  <img src="docs/branding/symbol.png" alt="keiailab" width="96"/>
</p>

# keiailab-commons

> **Shared Go library for Kubernetes operator scaffolding — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/<license placeholder>-blue.svg" alt="License"/></a>
  <!-- additional shield.io badges, in the §8 order -->
</p>

<p align="center">
  <b>English</b> |
  <a href="README.ko.md">한국어</a> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>
```

> **Template placeholder**: replace `<license placeholder>` with the copying
> repository's actual SPDX badge segment — for example `License-MIT` for MIT
> repositories (most of the family), or the matching `License-<SPDX-id>` for a
> differently licensed member. A single hardcoded license in this shared
> template mis-badges any member whose license differs, so the token stays
> generic here and is filled in per repo.

## 7. README footer standard

Every README and root-level `.md` file ends with a single attribution line:

```markdown
---

<p align="center">© 2026 keiailab · <license placeholder> · <a href="https://keiailab.com">keiailab.com</a></p>
```

> **Template placeholder**: replace `<license placeholder>` with the copying
> repository's actual license name — the same per-repo substitution rule as
> the §6 badge token. `keiailab-commons`' own actual value: MIT.

No additional cross-link block. Keep the footer minimal so the document
remains self-contained.

## 8. Badge order

Badges split into **MUST** (always present) and **SHOULD** (present only when
the backing infrastructure is live). Forcing all eight when a workflow or
registry does not yet exist leaves a badge perpetually red, so the family
standard is an honest 4 + 4.

**MUST (4)** — every family repository, left → right:

1. License
2. Go Version
3. Product (MongoDB / Valkey / PostgreSQL / Qdrant)
4. Kubernetes Version

**SHOULD (4)** — add each only when its backing surface is active:

5. Container image (`ghcr.io/keiailab/<repo>`)
6. Helm chart (Artifact Hub)
7. OpenSSF Scorecard
8. GitHub Discussions

> **Library exception**: `keiailab-commons` ships neither a container image, a
> Helm application chart, nor a Kubernetes workload, so its MUST set is
> License / Go Version / **Go Reference** (pkg.go.dev) in place of the product
> and Kubernetes badges, and its SHOULD set is OpenSSF Scorecard + GitHub
> Discussions. The container / Helm / Kubernetes badges belong on a downstream
> operator that ships an image or chart.

## 9. Discussions, issues, PR templates

- **Discussions**: `https://github.com/keiailab/keiailab-commons/discussions` — package API questions, integration patterns, new helper proposals.
- **Issues**: bug reports plus concrete feature requests; include the use case and downstream consumer impact when relevant.
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits + user-facing scenario + verification command output.

## 10. Social and external

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/keiailab-commons>

## 11. License and attribution

- License: [MIT](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attributions: see [NOTICE](../NOTICE)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
