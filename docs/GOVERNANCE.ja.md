# Governance

> [English](GOVERNANCE.md) | [한국어](GOVERNANCE.ko.md) | **日本語** | [中文](GOVERNANCE.zh.md)

> ⚠️ This translation is AI-generated and pending native review.

本ドキュメントは `keiailab/operator-commons` ライブラリの意思決定
プロセスを定義します。本ライブラリは downstream の consumer operator から
共通で import されるため、*public API* への変更は downstream 互換性に
影響を与えます。

## 原則

1. **オープン性** — すべての決定は public channel (GitHub Issue / PR /
   ADR) で行われます。
2. **Lazy Consensus** — 通常変更は明示的な反対がない限り進行します。
3. **Explicit Consensus** — public-API breaking change、新規パッケージ
   導入、ライセンス変更には maintainer の **2/3 supermajority** に加え、
   downstream consumer maintainer 最低 1 名の LGTM が必要です。
4. **責任の共有** — maintainer はライブラリの安定性、downstream 運用への
   影響、セキュリティに対して責任を共有します。

## 決定カテゴリ

### 通常変更 (Lazy Consensus)

- バグ修正、ドキュメント改善、テスト追加、minor / patch 依存性アップ
  グレード、public API を変えない内部リファクタ。
- プロセス: PR → maintainer 1 名以上の LGTM → merge。
- ウィンドウ: なし — ローカルゲートが pass すれば変更を merge 可能。
  本プロジェクトは GitHub Actions を使用しません。すべての品質ゲートは
  ローカル 4 階層 (`lefthook.yml`、`Makefile`、reviewer による証拠確認、
  ADR coverage) によって強制されます。

### 中間的な変更 (Explicit Consensus)

- 新規 public-API 関数または型の追加、メジャー依存性アップグレード、
  新規 `pkg/<sub>` パッケージの導入。
- プロセス: issue または ADR proposal → 7 日間コメント期間 →
  maintainer 過半数 LGTM → merge。
- 反対意見が 1 つ以上あれば maintainer 議論を起動します。

### Public-API breaking change (ADR 必須)

- 関数シグネチャの変更、型の削除、モジュールパス変更、ライセンス変更。
- プロセス:
  1. `docs/kb/adr/NNNN-<slug>.md` を提出。
  2. 14 日間コメント期間。
  3. maintainer 2/3 supermajority に加え downstream consumer 最低 1 名の
     LGTM。
  4. ADR を `Status: Draft → Accepted` に変更し、実装 PR を merge。

## セキュリティ決定

CVE 報告は [SECURITY.md](../SECURITY.ja.md) に従います。報告は private に
扱われます。downstream consumer が修正をリリースできるまで embargo を
維持します。

## リリース決定

- **v0.x**: maintainer 1 名で Lazy Consensus 下に minor / patch リリース
  を tag 可能。
- **v1.0+ (stable)**: 厳格な SemVer — major bump には ADR および 2/3
  supermajority が必要。

## 変更履歴

| 日付 | 変更 |
|---|---|
| 2026-05-09 | ドキュメント導入 — ガバナンス基盤。 |

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
