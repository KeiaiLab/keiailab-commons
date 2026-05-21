<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab オペレーター群が共有する Go ライブラリ — finalizer / labels / status / version / security / monitoring partials**
>
> [English](README.md) | [한국어](README.ko.md) | **日本語** | [中文](README.zh.md) (placeholder)

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <a href="README.ko.md">한국어</a> |
  <b>日本語</b> |
  <a href="README.zh.md">中文</a>
</p>

---

**keiailab** Kubernetes オペレーター (`mongodb-operator`、`valkey-operator`、`postgresql-operator`) が共有する Go ライブラリです。

> ステータス: **v0.x — API は変更される可能性があります**。v1.0 以降は SemVer (セマンティックバージョニング) で安定化されます。

## なぜ必要か

3 つのオペレーターがそれぞれ独立して同じ足回り (PodSecurity restricted コンテキスト、バージョン allowlist、NetworkPolicy テンプレート、ServiceMonitor ビルダー) を実装していました。リポジトリ間のメンテナンスドリフトが既に不整合を生み出し始めていました — 本ライブラリはその唯一の正本 (single source of truth) です。

## パッケージ一覧 (v0.7.0)

| パッケージ | 用途 |
|---|---|
| `pkg/version` | サポート対象 DB バージョン allowlist 規約 (`MustList`、`IsSupported`、`Strings`、`Default`) + ジェネリック `Matrix[E MatrixEntry]`。 |
| `pkg/security` | functional option を伴う PodSecurity *restricted* SecurityContext ビルダー。 |
| `pkg/labels` | 推奨 Kubernetes ラベル (`app.kubernetes.io/*`) ビルダー — `Set`、`All()`、`Selector()` (バージョン対応 split)。 |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` ビルダー (unstructured — CRD-soft)。 |
| `pkg/networkpolicy` | NetworkPolicy ビルダー — deny-by-default + functional option (`WithSelfIngress`、`WithIngressFromPeers`、`WithDenyEgress`、`WithEgressToPeers`)。 |
| `pkg/webhook` | Admission validation ヘルパー — `ValidateAllowedVersion` (完全一致)、`ValidateWithPredicate` (呼び出し側提供 matcher、例: semver-prefix)。 |
| `pkg/finalizer` | Finalizer ヘルパー — `Add` / `Remove` / `Has` (controller-runtime への依存を回避し、標準 `slices` のみ使用)。 |
| `pkg/status` | 4 つの標準 Condition Type + 6 つの Reason カタログ + ヘルパー (`SetReady`、`SetAvailable`、`SetReadyFalse`)。 |

`pkg/conditions` は *upstream の `k8s.io/apimachinery/pkg/api/meta.SetStatusCondition` の活用を推奨* します (commons には追加しない決定 — boundary 分析の結果。詳細は mongodb-operator HANDOFF iteration 32 を参照)。

## 採用マトリクス (3 オペレーター)

| Operator | sec | ver | lab | mon | np | wh | 採用率 |
|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | **4/6 (67%)** |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **6/6 (100%)** 🎉 |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | **3/6 (50%)** |

valkey が *最初の 100% 採用* — 他のオペレーターの carbon-copy リファレンス役を担います。適用事例の commit:

- `pkg/security`: it8 (3 オペレーター cross-cut) — `23fd3da` mongodb / `a0be4cf` valkey / `ac2e647` postgres
- `pkg/version`: mongodb it9 `a8db040`、valkey it8
- `pkg/labels`: mongodb it27 `ebc5803`、postgres it28 `c68b451`、valkey it29 `e8428b1`
- `pkg/monitoring`: valkey it23 `1765b54`
- `pkg/networkpolicy`: valkey it25 `97162b5`、mongodb it26 `ca0ec27`
- `pkg/webhook`: valkey it31 `14be0db`、postgres it34 `1d8fa17`

⏳ の領域は *機能追加を伴う* (例: mongodb webhook server / ServiceMonitor reconciler) か、*別の抽象化が適している* (postgres の `version matrix.go` の `Combo` struct は `commons.MustList` よりも豊富で委譲には不向き) ため、*深堀りを保留* している状態です。

## 使い方

```go
import (
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/version"
)

var SupportedMongoDBVersions = version.MustList("8.0", "8.2", "8.3")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

## バージョニングとリリース

- v0.x: API の breaking change を許容します。各タグ (`v0.N.M`) は、パッケージ・公開 API・もしくは重要な振る舞いのいずれかをバンプします。
- 各 consumer オペレーターは `go.mod` の `require` でピン留めします — このリポジトリと 3 つのオペレーター間でローカル開発を行う際は `replace` ディレクティブも許容されます。
- v1.0 以降: セマンティックバージョニング。Breaking change には RFC が必須となります。

## コミュニティ

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — pkg API への質問、統合事例、新しいヘルパーの提案
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — バグ報告 / API リクエスト
- **Downstream**: 3 つのオペレーター (mongodb-operator / postgres-operator / valkey-operator) — `go.mod replace` または直接 `require` で使用
- **安定性マトリクス**: `pkg/labels`、`pkg/security`、`pkg/version`、`pkg/webhook` (v0.5+ で stable) / `pkg/networkpolicy`、`pkg/monitoring` (experimental)

## ライセンス

Apache-2.0 — [LICENSE](./LICENSE) を参照してください。AGPL/BUSL の推移的依存ゼロを目標とし、minor リリースごとに監査しています。

## 参照

- [English README](README.md) — canonical SSOT (正本)
- [한국어 README](README.ko.md) — 韓国語版
- [日本語用語集](docs/i18n/glossary-ja.md) — 標準用語集 (本リポジトリ)

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
