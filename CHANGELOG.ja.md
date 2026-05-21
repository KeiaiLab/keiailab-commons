# Changelog

> [English](CHANGELOG.md) | [한국어](CHANGELOG.ko.md) | **日本語**

> ⚠️ This translation is AI-generated and pending native review.

本ライブラリのすべての注目すべき変更は、このファイルに記録されます。
形式: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)。
バージョニング: [Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

自動生成: `git-cliff` — changelog は release-tag 時点に PR として
再生成されます。

## [Unreleased]

### Added (v0.9.x candidate)

- `pkg/probes` (Experimental tier) — `corev1.Probe` fluent builder。HTTP /
  HTTPS / TCP / Exec handler + kubelet デフォルト値 (Period = 10 s /
  Timeout = 1 s / SuccessThreshold = 1 / FailureThreshold = 3) +
  InitialDelay / Period / Timeout 負値を 0 に clamp + handler 未設定時
  `Build()` panic (fail-fast)。カバレッジ 100 %、lint ゼロ。
- `pkg/storageclass` (Stable tier 即時昇格) — DNS-1123 subdomain
  validator。`IsValid` / `Validate` (+ `ErrInvalidStorageClassName`
  sentinel) / `Normalize` (empty → nil cluster default + trim + pointer
  return) / `MustNormalize`。カバレッジ 100 %、lint ゼロ。
- `pkg/events` (Beta tier) — Kubernetes Event recorder helper および
  9 個の標準 `Reason` 定数 (Created / Updated / Deleted / Reconciled /
  ReconcileError / Provisioning / Ready / Degraded / Failed)。最小
  `Recorder` interface (`client-go` `record.EventRecorder` を import せず
  互換)。`Emit` / `Emitf` / `EmitWarning` /
  `EmitWarningf` (nil-safe) + `WrappedError`。カバレッジ 100 %、lint ゼロ。
- `pkg/pvc` (Beta tier) — PVC expansion helper および安全な in-place
  update パス。
- `pkg/topology` (Beta tier) — PVC topology spread helper および
  zone-aware affinity。

### Added (v1.0.0 graduation track)

- `scripts/godoc-coverage.sh` — パッケージ別 + 全体の godoc カバレッジ
  計測。v1.0 で要求される 80 % 閾値を検証。
- `docs/STABILITY.md` — 3-tier API 安定性の約束および昇格基準、
  breaking-change ポリシー。
- `pkg/status/REASONS.md` — Reason × Type × Status 使用マトリクス。
- `pkg/finalizer.EnsureOrder` — 複数 finalizer のための順序保証 helper。
  `desiredOrder` に対する stable sort、リスト未掲載の finalizer は末尾保持。
- `pkg/labels.AllV2` + `V2` struct — Kubernetes 1.30+ Recommended labels
  v2 マッピング。
- `pkg/version.AsMap` + `MarshalJSON` — `Matrix[E]` serializer、
  安定かつ JSON / YAML 互換のキー順序。
- `pkg/version/api_stability_test.go` — public-API-surface ガード。
- `pkg/networkpolicy.ComboPeer` + `WithComboIngressFromPeers` — CIDR +
  NamespaceSelector + PodSelector 複合 peer helper。
- `pkg/security.RestrictedPodSecurityContext` + option 群
  (`WithPodFSGroup`, `WithPodRunAsUser`, `WithPodRunAsGroup`) — Pod-level
  restricted SecurityContext。
- `pkg/security.RuntimeDefaultSeccompProfile` +
  `LocalhostSeccompProfile` + `UnconfinedSeccompProfile` — seccomp
  profile ポインタ helper。

## [0.7.0] — 2026-05-09

### Added

- `pkg/version`: generic `Matrix[E MatrixEntry]` — caller が提供する entry
  型をライブラリに委譲可能。
- `docs/kb/adr/0004-*` — `Matrix` generic 決定を記録した ADR。

## [0.6.0] — 2026-05-09

### Added

- `pkg/status`: `SetAvailable` + `SetReadyFalse` sugar helper。
- `docs/kb/adr/0003-*` — status sugar helper 決定を記録した ADR。
- `.codecov.yml` — 以前は consumer 横断で統一されていたカバレッジ
  下限値。

## [0.5.0] — 2026-05-09

### Added

- ガバナンス文書: AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY /
  MAINTAINERS / CODE_OF_CONDUCT。
- `pkg/status/`: 4 つの標準 Condition Type + 6 個の Reason カタログ +
  helper 群。外部依存: `k8s.io/apimachinery` のみ。
- `pkg/finalizer/`: `Add` / `Remove` / `Has` helper。controller-runtime は
  回避し、stdlib `slices` パッケージが唯一の依存。
- `templates/observability/_servicemonitor.tpl`: Helm chart partial
  named template `keiailab.observability.serviceMonitor`。
- `templates/observability/README.md`: メトリック命名規約および共有
  アラート推奨、consumer chart 利用法。
- `Makefile` (lint / test / audit / cover / tidy / tag)。
- `.golangci.yml` + `.custom-gcl.yml`。
- `CHANGELOG.md` + `docs/kb/deps/2026-05.md` 依存性 audit log。
- `docs/kb/adr/0002-tooling-unification-adoption.md` +
  `docs/kb/adr/INDEX.md`。
- `NOTICE` (Apache-2.0 §4(d) 遵守)。
- `CODEOWNERS`。
- README バッジ — License / Go / pkg.go.dev / OpenSSF Scorecard /
  Discussions および Community セクション。
- `renovate.json`。
- `lefthook.yml` (library 最小構成)。
- DCO Signed-off-by warn-only commit-msg gate。

### Changed

- ADR ディレクトリ移動: `docs/adr/` → `docs/kb/adr/`。
- `go` directive `1.25.0` → `1.25.7`。
- `pkg/finalizer` lint fix: 手書き `for` + `==` → `slices.Contains` /
  `slices.Index` / `slices.Delete` (modernize linter)。
- `pkg/status/conditions.go` の `SetReady` シグネチャを multi-line 化
  (`lll` 通過)。

## [0.4.0] — 2026-05-07

それ以前の履歴は git tag および release notes で追跡されています。本
`CHANGELOG.md` は audit cycle 中に作成されました。

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
