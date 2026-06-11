<p align="center">
  <img src="docs/branding/symbol.png" alt="keiailab" width="96"/>
</p>

# keiailab-commons

> **Kubernetes operator 共通 scaffolding のための Go ライブラリ — finalizer / labels / status / version / security / monitoring partials.**
>
> [English](README.md) | [한국어](README.ko.md) | **日本語** | [中文](README.zh.md)

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-0EA5A8.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/keiailab-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/keiailab-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/keiailab-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/keiailab-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/keiailab-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/keiailab-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>


> ⚠️ This translation is AI-generated and pending native review.

---

Kubernetes operator コードベースの scaffolding drift を取り除く再利用可能な
Go ライブラリです — PodSecurity restricted context、サポートバージョンの
allowlist、NetworkPolicy テンプレート、ServiceMonitor ビルダー、finalizer /
status ヘルパー、Helm library chart partial を、小さく安定した API 表面の
背後にまとめて提供します。

> ステータス: **v0.x — API 変更可能.** v1.0 から SemVer stable.

## Why

Operator 作成者は同じ scaffolding を繰り返し実装します — restricted
PodSecurity context、サポートバージョンマトリクス、default-deny NetworkPolicy、
ServiceMonitor ビルダー、finalizer ヘルパー、status condition カタログ。
これらを独立に再実装すると、似た reconciler の間で静かな不整合が生じ、
マイナーリビジョンを重ねるごとに少しずつ分岐します。`keiailab-commons` は
その scaffolding の単一ソースです — ヘルパーを import すれば canonical な
実装が手に入り、すべてのリポジトリで再発明する必要がなくなります。

## パッケージ

| パッケージ | Tier | 目的 |
|---|---|---|
| `pkg/finalizer` | Stable | Finalizer ヘルパー — `Add` / `Remove` / `Has` / `EnsureOrder` (stdlib `slices` のみ、controller-runtime 依存なし)。 |
| `pkg/labels` | Stable | Kubernetes 推奨ラベル (`app.kubernetes.io/*`) ビルダー — `Set`、`All()`、`Selector()`、v2 マッピング (`AllV2`)。 |
| `pkg/status` | Stable | 4 標準 Condition Type + 6 Reason カタログ + ヘルパー (`SetReady`、`SetAvailable`、`SetReadyFalse`)。 |
| `pkg/storageclass` | Stable | DNS-1123 storageClass バリデータ + `Normalize` / `MustNormalize` (empty → cluster default ポインタ)。 |
| `pkg/version` | Beta | バージョン allowlist 規約 (`MustList`、`IsSupported`、`Strings`、`Default`) + ジェネリック `Matrix[E MatrixEntry]` + シリアライザ。 |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` / `PrometheusRule` ビルダー (unstructured — CRD-soft)。 |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy ビルダー + functional options (`WithSelfIngress`、`WithIngressFromPeers`、`WithDenyEgress`、`WithEgressToPeers`、`ComboPeer`)。 |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext ビルダー + Pod / Container 分割 + seccomp profile ポインタ。 |
| `pkg/events` | Beta | 最小 `Recorder` インターフェース + 9 標準 `Reason` 定数 + `Emit` / `EmitWarning` / `WrappedError` (nil-safe)。 |
| `pkg/pvc` | Beta | PVC 拡張ヘルパー — 比較 + 安全なインプレース更新 (controller-runtime 依存 — ADR-0016)。 |
| `pkg/topology` | Beta | TopologySpreadConstraints HA デフォルト + ゾーン認識アフィニティビルダー。 |
| `pkg/probes` | Experimental | `corev1.Probe` fluent ビルダー — HTTP / HTTPS / TCP / Exec、kubelet default + clamp。 |
| `pkg/webhook` | Experimental | Admission validation ヘルパー — `ValidateAllowedVersion`、`ValidateWithPredicate`、conversion registry。 |
| `pkg/bundle` | Experimental | OLM v1 バンドルメタデータヘルパー — アノテーション、FBCスキーマタイプ、ディレクトリ検証 (ADR-0017)。 |

[docs/STABILITY.md](docs/STABILITY.md) が tier の約束を定義します。
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) はパッケージ表面と設計不変条件を
カバーします。[docs/ROADMAP.md](docs/ROADMAP.md) は tier 昇格基準と
v1.0 卒業チェックリストを追跡します。

## 使い方

```go
import (
    "github.com/keiailab/keiailab-commons/pkg/security"
    "github.com/keiailab/keiailab-commons/pkg/version"
    corev1 "k8s.io/api/core/v1"
)

var supportedVersions = version.MustList("1.0", "1.1", "1.2")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

各パッケージごとの例は `pkg/<name>/doc.go` のパッケージドキュメントにあります
(`go doc github.com/keiailab/keiailab-commons/pkg/<name>`)。

## バージョニングとリリース

- **v0.x**: API breaking 許容。各タグ (`v0.N.M`) はパッケージの公開 API
  または意味ある動作変更を伴います。コンシューマは `go.mod` で特定バージョンに
  ピンします。
- **v1.0 以降**: Semantic Versioning。Breaking change は ADR (`docs/kb/adr/`) 必須。
- ローカルの `replace` ディレクティブは cross-repo 開発に許容; リリースタグは常に
  canonical module path を保持します。

## コミュニティ

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/keiailab-commons/discussions) — パッケージ API 質問、integration 事例、新 helper 提案。
- **Issues**: [GitHub Issues](https://github.com/keiailab/keiailab-commons/issues) — バグおよび具体的な機能要望。
- **Security**: 非公開報告手順は [SECURITY.md](SECURITY.md) を参照。
- **Contributing**: 開発ワークフローは [CONTRIBUTING.md](CONTRIBUTING.md) を参照。

## ライセンス

MIT — [LICENSE](LICENSE) 参照。AGPL / BUSL transitive 依存 0 件目標
(各 minor リリースで監査)。

---


<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
