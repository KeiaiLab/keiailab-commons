# ブランドガイド — `keiailab-commons`

> [English](BRANDING.md) | [한국어](BRANDING.ko.md) | **日本語** | [中文](BRANDING.zh.md)

> ⚠️ This translation is AI-generated and pending native review.
>
> `keiailab-commons` ライブラリの visual identity、voice、tone。

本ドキュメントは `keiailab-commons` ブランディング決定の canonical reference
です。README、リリースノート、プロジェクトに関する外部コミュニケーションに
適用されます。

## 1. Identity

**Organization**: [keiailab](https://keiailab.com)。

**Project**: `keiailab-commons` — Kubernetes operator 共通 scaffolding
(finalizer / labels / status / version / security / monitoring partial) の
ための Go ライブラリです。

本ライブラリは Go モジュール `github.com/keiailab/keiailab-commons` と
Helm library chart (`charts/keiailab-commons`) として公開されています。
標準 Go モジュール import によって downstream operator が利用しますが、
特定の consumer を指名・推奨することはありません。

## 2. ロゴとビジュアル資産

| 資産 | URL | 用途 |
|---|---|---|
| Primary ロゴ (SVG) | `https://keiailab.com/assets/logo.svg` | README ヘッダー、スライド |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon、ソーシャルカード |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer、ダーク背景 |

**ロゴ配置**: README 上部中央、width 120 px。常に `https://keiailab.com`
にリンク。

**Clear space**: ロゴ周辺の最小 padding はロゴ width の 25 %。

**禁止事項**:

- ロゴの色変更
- drop shadow / filter の追加
- コントラスト不足な背景への配置
- keiailab ブランド承認なしで他ロゴと結合

## 3. カラーパレット

| 役割 | Hex | 用途 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | ヘッダー、primary action、リンク |
| Secondary (deep navy) | `#0F172A` | ダーク背景、コードブロック |
| Accent (warm amber) | `#F59E0B` | 強調、バッジ accent |
| Neutral grey | `#64748B` | light 背景の body text |
| Background light | `#F8FAFC` | ドキュメントページ背景 |
| Background dark | `#020617` | ダークモードコードエディタテーマ |

GitHub README shield.io バッジは同じ hex 値を使用します。

## 4. タイポグラフィ

- **Heading**: システムデフォルト (GitHub の `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: システムデフォルト (GitHub-native 整合)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...` (GitHub デフォルト monospace)

別 web フォント未使用 — GitHub ネイティブ rendering を維持します。

## 5. Voice and Tone

**Audience**: Kubernetes プラットフォームエンジニア、DBA、SRE、Go ライブラリ
consumer。

**Voice 原則**:

- **Direct** — 可能な限り段落より bullet point。
- **Evidence-based** — claim は benchmark、SLA、link を伴います。
- **Library-focused** — `keiailab-commons` は *ライブラリ* です。
  controller-runtime、CRD、reconciler は downstream consumer の責任であり、
  本ライブラリの責任ではありません。
- **License-aware** — MIT only。charter の目標は AGPL / BUSL
  transitive 依存ゼロ件です (`docs/kb/adr/0001-charter.md`)。

**Avoid**:

- マーケティング最上級 ("blazing fast"、"revolutionary"、"best-in-class")。
- 曖昧な比較 ("enterprise-grade quality") — 各 claim を具体的 metric か
  benchmark で qualify。
- ロードマップの時間ベース締切 — [ROADMAP.md](ROADMAP.md) の機能チェック
  リストを使用。

## 6. README Header 標準

すべての README の最初のブロックは次のフォーマットに従います:

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# keiailab-commons

> **Kubernetes operator 共通 scaffolding のための Go ライブラリ — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <!-- 追加 shield.io バッジ -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <b>日本語</b> |
  <a href="README.zh.md">中文</a>
</p>
```

## 7. README Footer 標準

すべての README とルートレベル `.md` ファイルは次の単一行 attribution で
終わります:

```markdown
---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
```

追加の cross-link ブロックは置きません。Footer を最小化して文書を
self-contained に保ちます。

## 8. バッジ順序

README の shield.io バッジは次の順序 (左→右):

1. License (MIT)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `keiailab-commons` は *ライブラリ* なので、container image、
> Helm chart、Kubernetes deployment バッジは本ライブラリに付けません —
> イメージや chart を出荷する downstream operator の README に置きます。

## 9. Discussions / Issues / PR Template

- **Discussions**: `https://github.com/keiailab/keiailab-commons/discussions` — パッケージ API 質問、integration 事例、新 helper 提案。
- **Issues**: バグ報告および use case ありの具体的 feature request。関連時に downstream consumer 影響を明示。
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits + ユーザーシナリオ + 検証コマンド出力引用。

## 10. Social and External

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/keiailab-commons>

## 11. License and Attribution

- License: [MIT](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: [NOTICE](../NOTICE) 参照

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
