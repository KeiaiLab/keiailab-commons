<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# keiailab オペレーターファミリー

> ⚠️ This translation is AI-generated and pending native review. — 本翻訳は Claude による機械翻訳結果です。母語話者による検証完了まで `[検証必要]` 状態です。

> 共通基盤の上に構築された 4 つの姉妹 Kubernetes オペレーター — `operator-commons` (Go ライブラリ) + Helm partial + Apache-2.0 スタック。

本ページは **`operator-commons`** リポジトリから読み出されています。本ページはファミリー全体の canonical cross-link です。

## ファミリー概要

| プロジェクト | データベース | 状態 | リポジトリ |
|---|---|---|---|
| **`postgres-operator`** | PostgreSQL 18+ | active | https://github.com/keiailab/postgres-operator |
| **`mongodb-operator`** | MongoDB 7.0+ | active | https://github.com/keiailab/mongodb-operator |
| **`valkey-operator`** | Valkey 8.0+ (Redis fork, BSD-3) | active | https://github.com/keiailab/valkey-operator |
| **`operator-commons`** | 共有 Go ライブラリ | **v0.8.0** (現在のページ) | https://github.com/keiailab/operator-commons |

## 共有しているもの

4 プロジェクトすべてが同一の運用 primitive に収束しています:

- **Apache-2.0** end-to-end — SSPL 無し、SaaS サーフェスの copyleft 無し
- **`operator-commons`** 共有 Go ライブラリ (v0.8.0+) — ファイナライザー、ラベル、状態シュガー、security context ビルダー、NetworkPolicy / ServiceMonitor partial
- **Helm chart skeleton** — RFC-0027 `default` falsy-toggle 防止、RFC-0026 コンポーネント kkey の values、cycle 26 hardening 6 marker (priorityClassName / lifecycle / SA / minReadySeconds / automount / revisionHistoryLimit)
- **OLM bundle parity** — scorecard v1alpha3 6-test matrix
- **i18n** — README + canonical docs を英語 / 한국어 / 日本語 / 中文 (cleanup supercycle 2026-05-21 の Wave 4)

## ファミリー内での `operator-commons` の役割

本リポジトリは **共有 Go ライブラリ** です — コントローラーでは*ありません*. 以下を提供します:

| パッケージ | 目的 | Tier |
|---|---|---|
| `pkg/finalizer` | std `slices` のみを使用する `Add` / `Remove` / `Has` ファイナライザーヘルパー | **Stable** |
| `pkg/labels` | 推奨 K8s ラベルビルダー — `Set`, `All()`, `Selector()` | **Stable** |
| `pkg/status` | 4 つの標準 Condition Type + 6 つの Reason カタログ + ヘルパー | **Stable** |
| `pkg/version` | DB バージョン allowlist convention + generic `Matrix[E MatrixEntry]` | Beta |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` ビルダー (unstructured) | Beta |
| `pkg/networkpolicy` | Deny-by-default NetworkPolicy ビルダー + functional option | Beta |
| `pkg/security` | PodSecurity *restricted* SecurityContext ビルダー | Beta |
| `pkg/webhook` | Admission validation ヘルパー | Experimental |

設計 invariant: **leaf パッケージは stdlib + k8s API 型のみ**. controller-runtime 無し、logr 無し、operator-sdk leak 無し。

詳細なパッケージ surface は [ARCHITECTURE.md](ARCHITECTURE.md), tier 昇格基準は [ROADMAP.md](ROADMAP.md) を参照。

## 我々が*しない*こと

- ❌ **upstream オペレーターの埋め込み・wrap** (PGO, CloudNativePG, MongoDB Community Operator, Sentinel) — license-clean, copyleft 義務無し
- ❌ **GitHub Actions を release gate として使用** — ローカル 4 階層 + GitLab CI L5 (RFC-0002, RFC-0043 参照)
- ❌ **時間ベースのロードマップ締切** — 機能チェックリスト + 完了パーセント (`standards/roadmap.md §1.1` 参照)
- ❌ **Bitnami chart / image** — registry deprecation リスク、Broadcom 買収 (ADR-0136 / ADR-0057 参照)
- ❌ **本リポジトリの CRD / Reconciler** — consumer operator が当該責務を保有

## どこから始めるか

| 作業 | 入口 |
|---|---|
| オペレーターで `operator-commons` をインポート | [README.md](../README.md) Usage セクション |
| アーキテクチャを読む | [ARCHITECTURE.md](ARCHITECTURE.md) |
| イシューや機能要望を提出 | https://github.com/keiailab/operator-commons/issues |
| 設計やロードマップを議論 | https://github.com/keiailab/operator-commons/discussions |
| コード貢献 | [CONTRIBUTING.md](../CONTRIBUTING.md) |
| セキュリティイシュー報告 | [SECURITY.md](../SECURITY.md) |
| ブランド / ボイス学習 | [BRANDING.md](BRANDING.md) |
| アダプター追跡 / 利用者確認 | [ADOPTERS.md](ADOPTERS.md) |
| メンテナーを探す | [MAINTAINERS.md](MAINTAINERS.md) |
| ガバナンスモデル検討 | [GOVERNANCE.md](GOVERNANCE.md) |
| 今後の作業確認 | [ROADMAP.md](ROADMAP.md) |
| API 安定性の約束を確認 | [docs/STABILITY.md](STABILITY.md) |

## ファミリー間互換性

3 つのデータベースオペレーターすべてが `github.com/keiailab/operator-commons` を一致するバージョン (現在 `v0.8.0+`) でインポートします:

```go
import (
    "github.com/keiailab/operator-commons/pkg/version"
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/labels"
    "github.com/keiailab/operator-commons/pkg/monitoring"
    "github.com/keiailab/operator-commons/pkg/finalizer"
    "github.com/keiailab/operator-commons/pkg/status"
)
```

`operator-commons` の breaking change は 3 データベースオペレーターすべてで同期 bump が必要 — supercycle Wave 5 の `make cross-validation` ターゲットで検証。

ライブ consumer matrix (3 operator × 8 package × 採用率 %) は [ADOPTERS.md](ADOPTERS.md) 参照。

## i18n

本ページ (および全ての canonical プロジェクト文書) は 4 言語で提供されます:

- **English** (canonical, 原本ファイル)
- [한국어](family.ko.md)
- [日本語](family.ja.md) (現在のファイル)
- [中文](family.zh.md)

疑わしい場合、技術内容は英語版が権威であり、ローカライズ版は同じ決定を母語表現で反映します。

---

<p align="center">
  <b>keiailab オペレーターファミリー</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
