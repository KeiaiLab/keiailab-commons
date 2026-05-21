# ブランドガイド — `operator-commons`

> ⚠️ This translation is AI-generated and pending native review. — 本翻訳は Claude による機械翻訳結果です。

> keiailab オペレーターファミリーの visual identity, voice, tone.

本文書は `operator-commons` ブランディング決定の canonical reference です。README, release note, マーケティング資料, プロジェクトを代表する third-party コミュニケーションに適用されます。

## 1. Identity

**Organization**: [keiailab](https://keiailab.com) — Kubernetes-native データプラットフォームオペレーター (Apache-2.0, license-clean, vanilla-upstream 互換).

**Project**: `operator-commons` — keiailab オペレーターの共有 Go ライブラリ — finalizer / labels / status / version / security / monitoring partial.

**Family**: [`operator-commons`](https://github.com/keiailab/operator-commons) 共有ライブラリを共有する 4 つの姉妹オペレーターの 1 つ:

| プロジェクト | データベース | リポジトリ |
|---|---|---|
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `mongodb-operator` | MongoDB 7.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ (Redis fork, BSD-3) | https://github.com/keiailab/valkey-operator |
| `operator-commons` | 共有 Go ライブラリ | https://github.com/keiailab/operator-commons |

## 2. ロゴ & ビジュアル資産

| 資産 | URL | 用途 |
|---|---|---|
| Primary ロゴ (SVG) | `https://keiailab.com/assets/logo.svg` | README header, スライド |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon, social card |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer, dark background |

**ロゴ配置**: README 上部中央, width 120px. 常に https://keiailab.com にリンク。

**Clear space**: ロゴ周りの最小 padding = ロゴ width の 25%。

**禁止事項**:
- ロゴの色変更
- drop shadow / filter 追加
- コントラスト不足の背景に配置
- keiailab ブランド承認無しに他ロゴと結合

## 3. カラーパレット

| 役割 | Hex | 用途 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | ヘッダー, primary action, リンク |
| Secondary (deep navy) | `#0F172A` | dark 背景, コードブロック |
| Accent (warm amber) | `#F59E0B` | 強調, バッジ accent |
| Neutral grey | `#64748B` | light 背景の body text |
| Background light | `#F8FAFC` | 文書ページ背景 |
| Background dark | `#020617` | コードエディターテーマ, dark mode |

GitHub README の shield.io badge は上記 hex 使用推奨。

## 4. タイポグラフィ

- **Heading**: システムデフォルト (GitHub の default `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: 同一 (GitHub-native 整合)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...`

別途 webfont 不使用 (GitHub README rendering 整合)。

## 5. Voice & Tone

**Audience**: Kubernetes プラットフォームエンジニア / DBA / SRE / Go ライブラリ consumer.

**Voice 原則**:
- **Direct** — 可能ならば段落より bullet point
- **Evidence-based** — claim は benchmark / SLA / link を含む
- **Library-focused** — `operator-commons` は *ライブラリ* である — controller-runtime / CRD / reconciler は *consumer operator の責任*
- **License-aware** — Apache-2.0 only, AGPL/BUSL transitive 依存性 0 件目標 (ADR-0001 charter)

**Avoid**:
- マーケティング最上級 ("blazing fast", "revolutionary", "best-in-class")
- 曖昧な比較 ("X-class quality") — *具体的 metric または benchmark で qualify*
- ロードマップの時間ベース締切 (`standards/roadmap.md §1.1` — 機能チェックリスト代用)

## 6. README Header 標準

すべての README の最初の段落は次の形式 (Wave 3 標準):

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab オペレーターの共有 Go ライブラリ — finalizer / labels / status / version / security / monitoring partial**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <b>日本語</b> |
  <a href="README.zh.md">中文</a>
</p>
```

## 7. README Footer 標準

すべての README + root-level .md ファイルの末尾に次の footer 付着 (Wave 3 標準):

```markdown
---

<p align="center">
  <b>keiailab オペレーターファミリー</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
```

## 8. Badge 標準順序

README の shield.io badge 順序 (左→右):

1. License (Apache-2.0)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `operator-commons` は *ライブラリ* なので Kubernetes / Container Image / Helm Chart badge は *consumer operator 側 README に配置*. 本ライブラリは `pkg.go.dev` + `OpenSSF Scorecard` 中心.

## 9. Discussions / Issues / PR テンプレート

- **Discussions**: `https://github.com/keiailab/operator-commons/discussions` — pkg API 質問, integration 事例, 新 helper 提案
- **Issues**: bug report + use case ありの具体的 feature request (consumer operator 側影響明示推奨)
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` 標準 (ユーザーシナリオ + 検証コマンド引用義務, `standards/checklist.md §3`)

## 10. Social & External

- **Website**: https://keiailab.com
- **GitHub Org**: https://github.com/keiailab
- **pkg.go.dev**: https://pkg.go.dev/github.com/keiailab/operator-commons

## 11. License & Attribution

- License: [Apache-2.0](LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: [NOTICE](NOTICE) 参照 (該当時)

---

<p align="center">
  <b>keiailab オペレーターファミリー</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
