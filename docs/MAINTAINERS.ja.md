# Maintainers

> [English](MAINTAINERS.md) | [한국어](MAINTAINERS.ko.md) | **日本語** | [中文](MAINTAINERS.zh.md)

> ⚠️ This translation is AI-generated and pending native review.

本ドキュメントは `keiailab/operator-commons` の意思決定権限を有する
maintainer を記録します。

## 現在の maintainer

| 名前 / チーム | GitHub | Role | スコープ |
|---|---|---|---|
| keiailab maintainers | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | Lead | 全領域 |

GitHub team `@keiailab/maintainers` がライブラリの全領域における merge
および release-tag 権限を保持します。

## Maintainer 資格

downstream consumer operator の maintainer、*または* 以下の基準を
6 ヶ月以上満たした contributor:

- merge 済 PR ≥ 10 件 (ライブラリの PR 頻度は一般的な operator より低い
  ため、基準は概ね半分)。
- review 済 PR ≥ 20 件 (downstream consumer PR も含めて可)。
- 1 つ以上の `pkg/` 領域への深い精通 (security、labels、webhook、
  monitoring、networkpolicy、version、status、finalizer、storageclass、
  events、probes、pvc、topology)。

## 追加手順

1. 既存 maintainer (または候補者本人) が issue または ADR を作成。
2. `@keiailab/maintainers` チームが lazy consensus を適用 (7 日間
   コメント期間)。
3. 反対なしで候補者を GitHub team に追加し、本ファイルを PR で更新。

## 非アクティブ maintainer

6 ヶ月連続で活動がない maintainer は emeritus に移動 (権限剥奪、honorary
roll に名前は保持)。

## Cross-repo 合意

*public-API breaking change* には ADR 段階で downstream consumer の
maintainer による LGTM が必要 — [GOVERNANCE.md](GOVERNANCE.ja.md) 参照。

## i18n ドキュメントオーナー

| 言語 | オーナー | ファイル | 責任 |
|---|---|---|---|
| English (canonical) | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | `README.md` および canonical 文書 | Source of truth |
| Korean | TaeHwan Park ([@eightynine01](https://github.com/eightynine01)) | `README.ko.md` および `*.ko.md` | EN canonical sync + 翻訳レビュー |
| Japanese | (募集中 — issue で立候補) | `*.ja.md` | AI 翻訳 + native レビュー |
| Chinese | (募集中 — issue で立候補) | `*.zh.md` | AI 翻訳 + native レビュー |

**Drift 検証**: `bash scripts/check-readme-sync.sh` — ファイル存在、
セクションヘッダー数の一致、言語別閾値以下の行数差、cross-link の
双方向性を確認します。lefthook の `pre-push` hook `readme-i18n-sync` が
自動的に強制します。

## Emeritus

(まだなし)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
