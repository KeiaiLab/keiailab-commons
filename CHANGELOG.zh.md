# Changelog

> [English](CHANGELOG.md) | [한국어](CHANGELOG.ko.md) | [日本語](CHANGELOG.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文件记录本库的所有重要变更。
格式：[Keep a Changelog](https://keepachangelog.com/en/1.1.0/)。
版本规范：[Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

自动生成：`git-cliff` — 在 release-tag 时点以 PR 形式重新生成 changelog。

## [Unreleased]

### Added (v0.9.x candidate)

- `pkg/probes`（Experimental tier）— `corev1.Probe` fluent builder。HTTP /
  HTTPS / TCP / Exec handler + kubelet 默认值（Period = 10 s /
  Timeout = 1 s / SuccessThreshold = 1 / FailureThreshold = 3）+
  InitialDelay / Period / Timeout 负值自动夹紧为 0 + 未设置 handler 时
  `Build()` panic（fail-fast）。100 % 覆盖率，零 lint。
- `pkg/storageclass`（立即 Stable tier）— DNS-1123 subdomain
  校验器。`IsValid` / `Validate`（+ `ErrInvalidStorageClassName`
  sentinel）/ `Normalize`（空值 → nil cluster default + trim + 指针
  返回）/ `MustNormalize`。100 % 覆盖率，零 lint。
- `pkg/events`（Beta tier）— Kubernetes Event recorder 辅助以及 9 个
  标准 `Reason` 常量（Created / Updated / Deleted / Reconciled /
  ReconcileError / Provisioning / Ready / Degraded / Failed）。极简
  `Recorder` interface（不导入即可与 `client-go` `record.EventRecorder`
  兼容）。`Emit` / `Emitf` / `EmitWarning` /
  `EmitWarningf`（nil-safe）+ `WrappedError`。100 % 覆盖率，零 lint。
- `pkg/pvc`（Beta tier）— PVC expansion 辅助以及安全的 in-place
  update 路径。
- `pkg/topology`（Beta tier）— PVC topology spread 辅助以及 zone-aware
  affinity。

### Added (v1.0.0 graduation track)

- `scripts/godoc-coverage.sh` — per-package + total godoc 覆盖率
  测量。验证 v1.0 所需的 80 % 阈值。
- `docs/STABILITY.md` — 三层 API 稳定性承诺以及 graduation
  criteria 与 breaking-change 政策。
- `pkg/status/REASONS.md` — Reason × Type × Status 使用矩阵。
- `pkg/finalizer.EnsureOrder` — 多 finalizer 的顺序保证辅助；
  以 `desiredOrder` 为基准进行稳定排序，未列出的 finalizer 保留在尾部。
- `pkg/labels.AllV2` + `V2` struct — Kubernetes 1.30+ Recommended labels
  v2 映射。
- `pkg/version.AsMap` + `MarshalJSON` — `Matrix[E]` 序列化器，提供
  稳定的、JSON / YAML 兼容的 key 排序。
- `pkg/version/api_stability_test.go` — 公共 API surface 守卫。
- `pkg/networkpolicy.ComboPeer` + `WithComboIngressFromPeers` — CIDR +
  NamespaceSelector + PodSelector 组合 peer 辅助。
- `pkg/security.RestrictedPodSecurityContext` + options
  (`WithPodFSGroup`, `WithPodRunAsUser`, `WithPodRunAsGroup`) — Pod-level
  restricted SecurityContext。
- `pkg/security.RuntimeDefaultSeccompProfile` +
  `LocalhostSeccompProfile` + `UnconfinedSeccompProfile` — seccomp
  profile 指针辅助。

## [0.7.0] — 2026-05-09

### Added

- `pkg/version`：泛型 `Matrix[E MatrixEntry]` — 调用方提供的 entry
  type 可以委托给本库。
- `docs/kb/adr/0004-*` — 记录 `Matrix` 泛型决策的 ADR。

## [0.6.0] — 2026-05-09

### Added

- `pkg/status`：`SetAvailable` + `SetReadyFalse` sugar helper。
- `docs/kb/adr/0003-*` — 记录 status sugar helper 决策的 ADR。
- `.codecov.yml` — 以前在 consumers 之间统一的覆盖率底线。

## [0.5.0] — 2026-05-09

### Added

- 治理文档：AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY /
  MAINTAINERS / CODE_OF_CONDUCT。
- `pkg/status/`：4 个标准 Condition Type + 6-Reason 目录 +
  辅助。外部依赖：仅 `k8s.io/apimachinery`。
- `pkg/finalizer/`：`Add` / `Remove` / `Has` 辅助。回避了
  controller-runtime；stdlib `slices` 包是唯一依赖。
- `templates/observability/_servicemonitor.tpl`：Helm chart partial
  命名模板 `keiailab.observability.serviceMonitor`。
- `templates/observability/README.md`：metric 命名约定以及
  共享 alert 推荐以及 consumer chart 使用方法。
- `Makefile`（lint / test / audit / cover / tidy / tag）。
- `.golangci.yml` + `.custom-gcl.yml`。
- `CHANGELOG.md` + `docs/kb/deps/2026-05.md` 依赖审计日志。
- `docs/kb/adr/0002-tooling-unification-adoption.md` +
  `docs/kb/adr/INDEX.md`。
- `NOTICE`（Apache-2.0 §4(d) compliance）。
- `CODEOWNERS`。
- README badge — License / Go / pkg.go.dev / OpenSSF Scorecard /
  Discussions 以及 Community section。
- `renovate.json`。
- `lefthook.yml`（library minimal）。
- DCO Signed-off-by warn-only commit-msg gate。

### Changed

- ADR 目录迁移：`docs/adr/` → `docs/kb/adr/`。
- `go` directive `1.25.0` → `1.25.7`。
- `pkg/finalizer` lint 修复：手动 `for` + `==` → `slices.Contains` /
  `slices.Index` / `slices.Delete`（modernize linter）。
- `pkg/status/conditions.go` `SetReady` signature 现在是多行
  （通过 `lll`）。

## [0.4.0] — 2026-05-07

更早的历史通过 git tags 和 release notes 追踪；本
`CHANGELOG.md` 在 audit cycle 期间创建。

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
