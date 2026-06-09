# Security Policy

> [English](SECURITY.md) | [한국어](SECURITY.ko.md) | **日本語** | [中文](SECURITY.zh.md)

> ⚠️ This translation is AI-generated and pending native review.

`keiailab/operator-commons` は downstream の Kubernetes operator から
import されます。本ライブラリの脆弱性は、それらの downstream consumer の
運用セキュリティに直接影響を与える可能性があります。

## 脆弱性の報告

**public issue を作成しないでください。**

### チャネル

以下のいずれかの private チャネルを使用してください:

1. **GitHub Security Advisory** (推奨):
   `https://github.com/keiailab/operator-commons/security/advisories/new`
2. **Email**: `security@keiailab.com` (PGP optional):
   - PGP fingerprint:
     `89A4 0947 6828 CB99 2338  C378 651E 51AF 520B CB78`。

### 含める情報

- 影響を受けるバージョン (release tag または commit SHA)。
- 影響を受けるパッケージ (`pkg/security`、`pkg/webhook` 等)。
- 再現手順 (可能なら最小再現コード。downstream 環境に依存する場合は
  明記)。
- 影響評価 — downstream consumer への影響範囲。
- 自己評価による CVSS スコア (可能な場合)。

## レスポンス SLA

| 段階 | 時間 |
|---|---|
| 初動応答 (受領確認) | 72 時間以内 |
| 重大度評価 | 7 日以内 |
| パッチリリース | 重大度依存 (Critical: 14 日、High: 30、Medium: 60) |
| 公開開示 | パッチ後 14 日 (または downstream consumer が修正をリリースできる最早時点) |

## Embargo 取り扱い

public API に影響する脆弱性は、downstream consumer が同時に修正を
リリースできるまで embargo されます。Maintainer は開示に先立ち
downstream maintainer と private advisory を共有します。

## サポート対象バージョン

| バージョン | サポート |
|---------|-----------|
| 0.x (alpha) | ✅ 最新 minor のみ |
| 1.0+ (stable) | TBD — 初回 stable release 後に更新 |

本ライブラリは現在 v0.x です。public API は壊れる可能性があり、
セキュリティパッチは最新 minor に対してのみリリースされます。

## 依存性のセキュリティ

依存性の追加または更新時には、PR 本文にライセンスおよび CVE レビューを
引用します。Dependabot / Renovate の自動更新 PR は優先的にレビューされ
ます。

## ライセンス / supply chain

本ライブラリは **MIT only** であり、AGPL / BUSL の transitive
依存をゼロにする charter goal を持ちます (`docs/kb/adr/0001-charter.md`)。
すべての minor release でライセンス監査が実行されます。

## downstream consumer のベストプラクティス

本ライブラリを import する operator は以下を行うべきです:

1. **`pkg/security` を使用** — 独自実装せず、restricted PodSecurity
   SecurityContext builder を呼び出す。
2. **`pkg/webhook` を使用** — version validation を再実装しない。
3. **`pkg/networkpolicy` を使用** — deny-by-default NetworkPolicy
   builder。
4. `go.mod` 内の `github.com/keiailab/operator-commons` の最新パッチを
   追跡 (Renovate 自動 PR)。

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
