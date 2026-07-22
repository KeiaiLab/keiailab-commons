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

### 1.1 ファミリー charter (本ドキュメントが SSOT)

本 Branding Guide は keiailab operator ファミリーが共有するビジュアル文法の
**単一の真実 (SSOT)** です。カラーパレット (§3)、リンググラマーとグリフ入れ替え
手順 (§2b)、repo 別グリフレジストリ (§2b)、README ヘッダー / footer テンプレート
(§6 / §7)、バッジ順序 (§8) は **ここで** 定義し、各ファミリーリポジトリが
再導出せずコピーします。`keiailab-commons` が charter を所有する理由は、すべての operator が
既に import する共有依存だからです。

ファミリーは 4 つの sister operator と本共有ライブラリで構成されます:

| Project | Focus | Repository |
|---|---|---|
| `mongodb-operator` | MongoDB 8.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ | https://github.com/keiailab/valkey-operator |
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `qdrant-operator` | Qdrant vector database | https://github.com/keiailab/qdrant-operator |
| `keiailab-commons` | Shared Go library | https://github.com/keiailab/keiailab-commons |

本ガイドの翻訳が遅れている場合、ここの英語テキストが canonical のまま維持されます。

## 2. ロゴとビジュアル資産

| 資産 | URL | 用途 |
|---|---|---|
| Primary ロゴ | `docs/branding/symbol.png` | README ヘッダー、スライド |
| Keiailab base symbol | `docs/branding/base-symbol.png` | 外側の回転矢印マークのソースリファレンス |
| Favicon | `https://keiailab.com/favicon.ico` | Favicon、ソーシャルカード |
| 予定 SVG kit | `https://keiailab.com/assets/{logo,mark,wordmark}.svg` | URL が 200 を返した後に置換予定 |

**ロゴ配置**: README 上部中央、width 96 px。常に `https://keiailab.com`
にリンク。

**Clear space**: ロゴ周辺の最小 padding はロゴ width の 25 %。

**禁止事項**:

- ロゴの色変更
- drop shadow / filter の追加
- コントラスト不足な背景への配置
- keiailab ブランド承認なしで他ロゴと結合

## 2b. リンググラマーとグリフ入れ替え手順

ファミリーは 1 つの外側シンボル — `docs/branding/base-symbol.png` の 3 色回転
矢印リング (460×460; ダークネイビー → ブルー → グリーンの矢印、中央に球体) —
を共有します。**すべての** ファミリーリポジトリはこのファイルを byte 単位で
同一に vendor し、リングは決して再着色・再描画・再生成しません。

各プロジェクトの `docs/branding/symbol.png` は、リングを保持し **中央のグリフ
のみ** を入れ替えて作成します:

1. 共有 `base-symbol.png` から開始 — リングには触れません。
2. 製品グリフを中央の clear zone (~185 px 直径) に alpha-composite — リングと
   重なりません。
3. ソース解像度で `docs/branding/symbol.png` に export。
4. README ヘッダーで width 96 px で参照 (§6)、operator は Helm chart `icon`
   としても参照。

共有リング + グリフ入れ替えが、マークを 1 つのファミリーとして読ませつつ製品ごと
の識別性を保つ核心です。リングの再着色や新しい外側マークの作図は禁止です。

### グリフレジストリ

| Repository | Center glyph |
|---|---|
| `mongodb-operator` | leaf |
| `valkey-operator` | hexagon + keyhole |
| `postgres-operator` | elephant |
| `qdrant-operator` | kNN constellation |
| `keiailab-commons` | node grid |

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
  <img src="docs/branding/symbol.png" alt="keiailab" width="96"/>
</p>

# keiailab-commons

> **Kubernetes operator 共通 scaffolding のための Go ライブラリ — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/<license placeholder>-blue.svg" alt="License"/></a>
  <!-- §8 の順序に従う追加 shield.io バッジ -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <b>日本語</b> |
  <a href="README.zh.md">中文</a>
</p>
```

> **テンプレート placeholder**: `<license placeholder>` を、コピー先リポジトリの
> 実際の SPDX バッジセグメントに置き換えます — 例: MIT リポジトリ (ファミリーの
> 大半) は `License-MIT`、ライセンスが異なるメンバーは対応する
> `License-<SPDX-id>`。共有テンプレートにライセンスをハードコードすると、
> ライセンスが異なるメンバーに誤バッジが再生産されるため、トークンはここでは
> generic のままにし、リポジトリごとに埋めます。

## 7. README Footer 標準

すべての README とルートレベル `.md` ファイルは次の単一行 attribution で
終わります:

```markdown
---

<p align="center">© 2026 keiailab · <license placeholder> · <a href="https://keiailab.com">keiailab.com</a></p>
```

> **テンプレート placeholder**: `<license placeholder>` を、コピー先
> リポジトリの実際のライセンス名に置き換えます — §6 のバッジトークンと
> 同じリポジトリ別置換ルールです。`keiailab-commons` 自身の実際の値: MIT。

追加の cross-link ブロックは置きません。Footer を最小化して文書を
self-contained に保ちます。

## 8. バッジ順序

バッジは **MUST** (常に表示) と **SHOULD** (バックエンドインフラがライブのときのみ
表示) に分かれます。ワークフローやレジストリがまだ無いのに 8 個を強制すると
バッジが常時 red のままになるため、ファミリー標準は正直な 4 + 4 です。

**MUST (4)** — すべてのファミリーリポジトリ、左→右:

1. License
2. Go Version
3. Product (MongoDB / Valkey / PostgreSQL / Qdrant)
4. Kubernetes Version

**SHOULD (4)** — バックエンド表面がアクティブなときのみ各々追加:

5. Container image (`ghcr.io/keiailab/<repo>`)
6. Helm chart (Artifact Hub)
7. OpenSSF Scorecard
8. GitHub Discussions

> **ライブラリ例外**: `keiailab-commons` は container image、Helm application
> chart、Kubernetes ワークロードのいずれも出荷しないため、MUST セットは product と
> Kubernetes バッジの代わりに License / Go Version / **Go Reference** (pkg.go.dev)、
> SHOULD セットは OpenSSF Scorecard + GitHub Discussions です。container / Helm /
> Kubernetes バッジは、イメージや chart を出荷する downstream operator に置きます。

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
