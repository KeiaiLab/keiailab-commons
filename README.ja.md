# operator-commons

> [English](README.md) | [한국어](README.ko.md) | **日本語** (placeholder) | [中文](README.zh.md) (placeholder)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg)](https://pkg.go.dev/github.com/keiailab/operator-commons)

> **注意 (Notice)**: この日本語 README はプレースホルダーです。正本は [README.md](README.md) (English) を参照してください。ネイティブレビュアーによる完訳は今後の作業項目です (RFC-0025 `[~]` partial marker)。

## 概要 (Overview)

**keiailab** Kubernetes operator (`mongodb-operator`, `valkey-operator`, `postgresql-operator`) が共有する Go ライブラリです。

> ステータス: **v0.x — API は変更される可能性があります**。v1.0 以降は SemVer (セマンティックバージョニング) で安定化されます。

詳細な英語版本文は [README.md](README.md) を参照してください。

## パッケージ一覧 (Packages, v0.7.0)

| Package | 用途 |
|---|---|
| `pkg/version` | サポート対象 DB バージョン allowlist 規約 |
| `pkg/security` | PodSecurity *restricted* SecurityContext ビルダー |
| `pkg/labels` | 推奨 Kubernetes ラベル (`app.kubernetes.io/*`) ビルダー |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` ビルダー |
| `pkg/networkpolicy` | NetworkPolicy ビルダー — deny-by-default + functional options |
| `pkg/webhook` | Admission 検証ヘルパー |
| `pkg/finalizer` | Finalizer ヘルパー (controller-runtime 非依存) |
| `pkg/status` | 4 標準 Condition Type + 6 Reason カタログ + ヘルパー |

詳細な API シグネチャと使用例は [README.md](README.md) の `Packages`, `Usage` セクションを参照してください。

## ライセンス (License)

Apache-2.0 — [LICENSE](./LICENSE) を参照。minor リリースごとに監査される AGPL/BUSL 推移的依存ゼロを目標としています。

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a> ·
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/forgewise">forgewise</a>
</p>

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
