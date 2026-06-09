# Contributing to operator-commons

> [English](CONTRIBUTING.md) | [한국어](CONTRIBUTING.ko.md) | **日本語** | [中文](CONTRIBUTING.zh.md)

> ⚠️ This translation is AI-generated and pending native review.

`keiailab/operator-commons` は downstream の Kubernetes operator から
import される Go ライブラリです。すべての貢献は
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) および
[docs/GOVERNANCE.md](docs/GOVERNANCE.ja.md) に従います。

## 貢献フロー

1. **issue または ADR で意図を共有** (非自明な変更の場合)。
2. **fork + feature ブランチ** — `feat/<slug>`、`fix/<slug>`、
   `docs/<slug>`、または `refactor/<slug>` を使用。
3. **ローカルゲートを検証**:
   - `lefthook install --force`
   - `lefthook run pre-commit --all-files`
   - `make lint test`
4. **PR を作成** — Conventional Commits 形式。本文は英語または韓国語
   どちらでも可。
5. **レビュー SLA**: maintainer は 24 時間以内に返信します。

## PR チェックリスト (作成者)

- [ ] PR タイトル: Conventional Commits 形式 (`feat`、`fix`、`docs`、
  `refactor`、`test`、`chore`)。
- [ ] PR 本文: 変更概要 + 検証コマンド + 出力引用。
- [ ] unit test 追加または更新 (`pkg/<sub>` 内の変更は必須)。
- [ ] public API 変更時は godoc を更新。
- [ ] **public-API breaking change** の場合: ADR リンクおよび
  downstream consumer への影響分析。
- [ ] `go.mod` / `go.sum` の drift = 0 (`go mod tidy` を実行しても変更
  なし)。
- [ ] 新規依存性: PR 本文にライセンスおよび CVE レビューを引用。

## ローカル開発 (downstream consumer と横断する作業)

`operator-commons` と downstream operator を同時に変更する場合:

```fish
# 1. consumer operator の go.mod に replace directive を追加
#    (local only。commit しないこと)
# go.mod 末尾:
#   replace github.com/keiailab/operator-commons => ../operator-commons

# 2. 両側を編集 + 各々で `go test ./...` を実行

# 3. PR を分割:
#    - operator-commons 側: merge + tag (例: v0.9.0)
#    - consumer 側: require directive を bump (replace 削除)
```

## 新しい `pkg/<sub>` パッケージの追加

[docs/GOVERNANCE.md](docs/GOVERNANCE.ja.md) の「中間的な変更」プロセスに
従います:

1. *なぜ* commons に属するかと *どの* downstream consumer が利用するかを
   説明する issue または ADR を作成。
2. 7 日間のコメント期間を待つ。
3. 複数 maintainer の LGTM 後に merge。

## リリース

- v0.x SemVer: minor bump ごとに public-API 変更または意味のある挙動
  変更を表します。
- リリース手順:
  1. `git tag v0.X.Y` (annotated)。
  2. `git-cliff` が CHANGELOG PR を再生成。
  3. `git push origin v0.X.Y`。
  4. 各 downstream consumer で `require` directive を bump する
     follow-up PR を作成。

## セキュリティ脆弱性

[SECURITY.md](SECURITY.ja.md) の private 開示プロセスを使用してください。
public issue は適切なチャネルではありません。

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
