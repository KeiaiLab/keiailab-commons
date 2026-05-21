# Upgrading operator-commons

> [English](UPGRADING.md) | [한국어](UPGRADING.ko.md) | **日本語**

> ⚠️ This translation is AI-generated and pending native review.

本ドキュメントは `github.com/keiailab/operator-commons` Go モジュールの
minor または major バージョンを bump する際に必要となる移行手順を集約
します。downstream consumer の共通エントリポイントです。

## 0. バージョンポリシー (SemVer)

| 変更の種類 | SemVer bump | 例 |
|---|---|---|
| 新規パッケージ追加 | minor (v0.X → v0.X+1) | `pkg/events`、`pkg/storageclass` 導入 |
| 既存 API シグネチャ変更 (breaking) | major (v0.X → v1.0 / v1.X → v2.0) | `pkg/status.SetReady()` シグネチャ変更 |
| パッケージ内部の挙動変更 (非破壊) | patch (v0.X.Y → v0.X.Y+1) | バグ修正 |
| ADR からの逸脱 | major + Deprecated 通知 | API 安定性 tier 変更 |

API 安定性 tier (`pkg/<name>/doc.go` マーカー):

- **Stable** — minor release を通じて後方互換。
- **Beta** — 次の minor で変更される可能性あり。
- **Experimental** — いつでも変更される可能性あり。

## 1. v0.7.x → v0.8.x

### Helm library chart 利用者

```bash
helm dep update charts/<your-operator>
helm template <your-operator> charts/<your-operator>
```

`keiailab-commons` chart v0.8.0 の partial (`_servicemonitor.tpl`、
`_rbac.tpl`、`_networkpolicy.tpl`) には追加作業は不要です。

### Go モジュール利用者

```bash
go get github.com/keiailab/operator-commons@v0.8.0
go mod tidy
```

追加作業なし — 後方互換。

## 2. v0.8.x → v0.9.x

### 新規パッケージ (minor bump)

| パッケージ | 目的 | Tier |
|---|---|---|
| `pkg/pvc` | PVC expansion helper | Beta |
| `pkg/topology` | PVC topology spread + zone-aware affinity | Beta |

### 移行

downstream operator で import を追加:

```go
import (
    "github.com/keiailab/operator-commons/pkg/pvc"
    "github.com/keiailab/operator-commons/pkg/topology"
)
```

### 後方互換性

- 既存パッケージ (`pkg/status`、`pkg/finalizer`、`pkg/networkpolicy`、
  `pkg/monitoring`、`pkg/probes`、`pkg/labels`、`pkg/storageclass`、
  `pkg/webhook`、`pkg/events`、`pkg/security`、`pkg/version`) は
  シグネチャを維持します。
- `keiailab-commons` chart の `_security.tpl` および
  `_servicemonitor.tpl` partial は *opt-in* — 既存の inline 定義を
  そのままにしても影響はありません。

### 推奨移行手順

```bash
# 1. 依存性を bump
go get github.com/keiailab/operator-commons@v0.9.0
go mod tidy

# 2. 検証
make verify  # lint + test + build

# 3. e2e (kind)
kind create cluster
helm install <operator> charts/<operator>
kubectl apply -f config/samples/
kubectl get <CR> -A  # reconciliation を観測
```

## 3. v0.9.x → v1.0.0

v1.0.0 昇格基準 ([STABILITY.md](STABILITY.ja.md) の「v1.0.0 graduation」
参照) が満たされたときに進行します:

- 全パッケージが Stable tier に到達。
- v0.x → v1.0 は *naming* 変更 — セマンティクスは変更なし (breaking
  change なし)。

## 4. 汎用移行チェックリスト

アップグレード前:

- [ ] `go mod tidy` で変更なし (drift = 0)。
- [ ] `make audit` pass (govulncheck CVE = 0)。
- [ ] 既存 e2e suite pass。

アップグレード後:

- [ ] downstream import path 更新 (`go get -u` または pinned バージョン)。
- [ ] `make verify` pass。
- [ ] e2e pass。
- [ ] Helm chart `charts/<operator>` の `dependencies:` 更新。

## 5. Breaking-change 通知ポリシー

- **Deprecation**: 新しい minor で `// Deprecated:` コメント追加。
  2 minor 後に削除。
- **Breaking**: major bump + 本ファイルに専用セクション + ADR。
- **silent breaking change なし**: すべての breaking change には先行する
  minor で最低 1 回の deprecation を伴います。

## 参照

- ADR index: [`docs/kb/adr/INDEX.md`](kb/adr/INDEX.ja.md)。
- API 安定性: `pkg/<name>/doc.go` の tier マーカー。
- i18n: [`docs/i18n/README.md`](i18n/README.ja.md) (多言語ポリシー)。

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
