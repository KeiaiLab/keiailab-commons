# i18n — operator-commons 多言語ポリシー

> [English](README.md) | [한국어](README.ko.md) | **日本語**

> ⚠️ This translation is AI-generated and pending native review.

本ドキュメントは `operator-commons` ドキュメントを複数言語でどのように
維持するかを定義します。本プロジェクトは英語 canonical に加え、韓国語 /
日本語 / 中国語の翻訳を保持します。

## §1 ポリシー

### 1.1 基本原則

- **Canonical 言語**: 英語 (すべての source の ground truth)。
- **翻訳対象言語 (3 つ追加)**: 韓国語 (`ko`)、日本語 (`ja`)、中国語
  (`zh`、簡体字)。
- **4 言語骨格**: `README.md` (英語) + `README.ko.md` +
  `README.ja.md` + `README.zh.md`。
- **この 4 言語以外不可** — `es` / `fr` / `de` / … の追加には別途 ADR
  が必要。

### 1.2 翻訳対象 vs 対象外

**対象** (ユーザー向けまたは外部文書):

- `README.{md, ko, ja, zh}.md`
- トップレベルのガバナンス / ブランディング: `BRANDING.md`、
  `STABILITY.md`、`coverage-report.md` などの 4 言語マトリクス。
- `docs/i18n/glossary-{ko, ja, zh}.md` (この SSOT)。

**対象外** (内部、ガバナンス、または法的):

- `docs/kb/adr/*.md` (決定記録 — 英語 canonical のみ)。
- `docs/kb/deps/*.md` (依存性監査 — 自動生成)。
- `LICENSE`、`NOTICE`、`CITATION.cff`、`.gitignore` (法的または設定)。

## §2 用語集 (Glossary)

### 2.1 場所

SSOT: [`glossary-ko.md`](glossary-ko.md) / [`glossary-ja.md`](glossary-ja.md)
/ [`glossary-zh.md`](glossary-zh.md)。

### 2.2 優先順位

1. Glossary の定義が最優先。
2. Kubernetes 公式の韓国語 / 日本語 / 中国語ガイド。
3. 一般的な辞書。
4. **コード識別子は英語のまま** (翻訳しない) — `pkg/probes`、
   `Reconciler`、`kubectl` 等。

### 2.3 一貫性ルール (4 言語共通)

- 1 つの文書内で敬体と常体を混在させない。
- ユーザー向け文書は敬体を使用。
- 内部文書 (AGENTS) は常体可。
- 初出時は英語原文に加え括弧で訳語を併記。以降は訳語単独可。

## §3 自動翻訳 SOP

### 3.1 エンジン

**デフォルト**: Claude direct (AI subagent が source を読んで翻訳)。

根拠: 韓国語品質と prompt caching コスト。日本語および中国語の翻訳には
native レビューが後続します。

### 3.2 コマンド

```bash
# 1 ファイル手動翻訳
# subagent が source を読んで翻訳し、AI 翻訳 warning バナー付きで
# 出力ファイルを書き出します。

# 自動化 (スクリプト)
./scripts/i18n-translate.sh <source.md> --lang all --engine claude

# Dry run
./scripts/i18n-translate.sh README.md --dry-run
```

### 3.3 結果マーカー

すべての AI 翻訳ファイルには警告バナーを含める必要があります:

```markdown
> ⚠️ This translation is AI-generated and pending native review.
```

追加マーカー: `[needs review]` (レビュー待ち)、`[reviewed]` (native
reviewer による PR 経由昇格済)。

## §4 Native Review SOP

### 4.1 reviewer 確保

デフォルトポリシー: AI 翻訳 + `[needs review]` マーカー + 警告バナー。
native reviewer が確保されるまで placeholder を維持。

### 4.2 昇格基準

`[needs review]` → `[reviewed]` へ昇格する際、native reviewer が確認:

1. **意訳 vs 直訳のバランス** — 自然な読みを優先しつつ、技術用語の
   正確性を保つ。
2. **Glossary 一貫性** — すべての glossary 用語が同一に適用されている。
3. **敬体 / 常体の一貫性** — 1 文書内で敬体 / 常体が統一されている。
4. **リンク整合性** — Cross-link が言語別に解決される。
5. **コード識別子は翻訳されていない**。

### 4.3 昇格 PR フォーマット

```
title: docs(i18n): <lang> <doc> native reviewer 昇格
body:
  ## Reviewer
  - <name / handle>

  ## Scope reviewed
  - <file path>

  ## Changes
  - [needs review] → [reviewed] marker 変更
  - 警告バナー削除 (または partial で保持)
  - <その他の意訳修正>

  ## Verification
  - [x] Glossary 一貫性
  - [x] 敬体 / 常体の統一
  - [x] Cross-link 整合性
```

## §5 Drift Control

### 5.1 lefthook hook

`lefthook.yml` は pre-push で `readme-i18n-sync` を強制:

```yaml
readme-i18n-sync:
  glob: "README*.md"
  run: bash scripts/check-readme-sync.sh
```

`pre-push` の位置により、作業フローへの摩擦が低く抑えられます。

### 5.2 閾値

| 言語 | 行差閾値 | 根拠 |
|---|---|---|
| ko | ≤ 15 % | 韓国語の LOC は英語とほぼ同じ。 |
| ja | ≤ 25 % | 日本語は仮名と漢字混在 — 英語より ~15-20 % 短い。 |
| zh | ≤ 30 % | 中国語は漢字 100 % — 英語より ~25 % 短い。 |

上記は推奨値。実測後に調整してください。

### 5.3 4 言語マトリクス

`scripts/check-readme-sync.sh` は 4 言語マトリクス (EN ↔ {ko, ja, zh})
を自動チェック:

- target 言語ファイル不在時は skip (部分骨格に対応)。
- 言語別バイパス可能
  (`SKIP_CHECK_README_SYNC_JA=1` 等)。

## §6 参照

- Glossary: [`glossary-ko.md`](glossary-ko.md) /
  [`glossary-ja.md`](glossary-ja.md) / [`glossary-zh.md`](glossary-zh.md)。
- Check script: [`../../scripts/check-readme-sync.sh`](../../scripts/check-readme-sync.sh)。
- Translation script: [`../../scripts/i18n-translate.sh`](../../scripts/i18n-translate.sh)。
