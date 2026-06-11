# Changelog

> [English](CHANGELOG.md) | [한국어](CHANGELOG.ko.md) | [日本語](CHANGELOG.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文件记录本库的所有重要变更。
格式：[Keep a Changelog](https://keepachangelog.com/en/1.1.0/)。
版本规范：[Semantic Versioning](https://semver.org/spec/v2.0.0.html)。

自动生成：`git-cliff` — 在 release-tag 时点以 PR 形式重新生成 changelog。

## [Unreleased]

### Added

- `pkg/apply`（Beta tier）— ConfigMap / Secret / Service / StatefulSet /
  Deployment / NetworkPolicy / PodDisruptionBudget /
  HorizontalPodAutoscaler 的 idempotent apply 辅助。Service ClusterIP /
  IPFamilies create-only 守卫、StatefulSet immutable 字段保留 +
  `RetryOnConflict`、Deployment server-default 保留 + revision annotation
  保留、`preserveReplicas` 选项（避免与 HPA 冲突）。消除三个 operator 间
  约 530 LOC 的重复。依赖 controller-runtime。
- `pkg/reconcile`（Beta tier）— `Statusable` interface（`client.Object` +
  `GetConditions` + `SetPhase`）以及 `ApplyErrorCondition` +
  `HandleFinalizerCleanup` + `SecretIfNotExists`。依赖 controller-runtime。
- `pkg/certmanager`（Beta tier）— `CertParams` + `BuildCertificate` +
  `BuildSelfSignedIssuer` + `ServiceSANs`。基于 unstructured；对
  cert-manager CRD 的 Go 依赖为零。
- `pkg/reconcilemetrics`（Beta tier）— `ReconcileMetrics`（Total /
  Latency / Errors）+ `New(subsystem)` + `MustRegister` + `IncTotal` /
  `ObserveReconcile` / `IncError` / `DeleteFor` + `ResultFor`。通过注入
  subsystem 保留各 operator 现有的 Prometheus 时间序列名称。
- `pkg/status` `UpdateWithRetry`（Stable 包内的 Beta surface）— refetch +
  mutate + `RetryOnConflict` 持久化 status subresource。

### Changed

- 新增直接依赖：`github.com/prometheus/client_golang` v1.23.2
  （由 `pkg/reconcilemetrics` 引入）。

## [0.10.0] — 2026-06-11

### Added

- `pkg/bundle`（Experimental tier）— OLM v1 operator bundle metadata
  辅助：`Annotations`（+ `DockerLabels`）、FBC `Package` / `Channel` /
  `Bundle` builder、`ValidateDir`。
- Helm chart partial 命名模板 `keiailab.secrets.externalSecret`
  （`charts/keiailab-commons/templates/_externalsecret.tpl`）。
- GitLab CI shadow pipeline。

### Changed

- **BREAKING**：module path 变更 —
  `github.com/keiailab/operator-commons` →
  `github.com/keiailab/keiailab-commons`。所有 consumer 都需更新 import
  path；因处于 0.x 阶段，以 minor bump 发布。
- 许可证标准化为 MIT — 所有 `.go` 文件添加 SPDX 头。
- README 以准确性为准重写。
- self-contained 重构 — 移除外部服务引用（B1～B14）。

## [0.9.0] — 2026-05-21

### Added

- `pkg/pvc`（Beta tier）— PVC expansion 辅助以及安全的 in-place
  update 路径。
- `pkg/topology`（Beta tier）— PVC topology spread 辅助以及 zone-aware
  affinity。
- `scripts/release.sh` — 库的手动 release pipeline（ADR-0014）。
- `docs/UPGRADING.md` — semver 政策以及三个 operator 的迁移指南。
- i18n S4 Phase 1～5 — 完成 4 语言 glossary 以及翻译 sync hook。

### Changed

- keiailab branding Wave 3 — README header / footer 以及 `BRANDING.md`、
  `docs/family.md`。

## [0.8.0] — 2026-05-21

### Added

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
- `pkg/monitoring.NewPrometheusRule` + `AlertRule` / `RecordingRule` /
  `RuleGroup` — PrometheusRule（`monitoring.coreos.com/v1`）manifest
  builder。
- `pkg/webhook.ConversionRegistry` — CRD version pair 转换函数
  registry（`Register` / `Convert` / `HasPair`）。
- `pkg/networkpolicy.ComboPeer` + `WithComboIngressFromPeers` — CIDR +
  NamespaceSelector + PodSelector 组合 peer 辅助。
- `pkg/security.RestrictedPodSecurityContext` + options
  (`WithPodFSGroup`, `WithPodRunAsUser`, `WithPodRunAsGroup`) — Pod-level
  restricted SecurityContext。
- `pkg/security.RuntimeDefaultSeccompProfile` +
  `LocalhostSeccompProfile` + `UnconfinedSeccompProfile` — seccomp
  profile 指针辅助。
- `pkg/version.AsMap` + `MarshalJSON` — `Matrix[E]` 序列化器，提供
  稳定的、JSON / YAML 兼容的 key 排序。
- `pkg/version/api_stability_test.go` — 公共 API surface 守卫。
- `pkg/finalizer.EnsureOrder` — 多 finalizer 的顺序保证辅助；
  以 `desiredOrder` 为基准进行稳定排序，未列出的 finalizer 保留在尾部。
- `pkg/labels.AllV2` + `V2` struct — Kubernetes 1.30+ Recommended labels
  v2 映射。
- `pkg/status/REASONS.md` — Reason × Type × Status 使用矩阵。
- `docs/STABILITY.md` — 三层 API 稳定性承诺以及 graduation
  criteria 与 breaking-change 政策。
- `scripts/godoc-coverage.sh` — per-package + total godoc 覆盖率
  测量。验证 v1.0 所需的 80 % 阈值。
- `docs/ARCHITECTURE.md` — 单页架构说明。
- README 4 语言 i18n 启动 — 英文 canonical + 韩文翻译、4 语言 switcher
  以及日文 / 中文 placeholder。

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

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
