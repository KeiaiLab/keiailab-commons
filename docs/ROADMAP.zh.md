# ROADMAP — keiailab-commons

> [English](ROADMAP.md) | [한국어](ROADMAP.ko.md) | [日本語](ROADMAP.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本 ROADMAP 沿三条轴线追踪库的演进：*API
stability tier*、*v1.0.0 graduation criteria* 和 *per-package
follow-up items*。本项目不维护基于时间的截止日期 —
本库根据下游 consumer 的需求演进。

## Checkbox 含义

| Marker | Meaning |
|---|---|
| `[x]` | 代码 + 测试均存在。下游 import 可用。 |
| `[~]` | 部分实现（helper 存在，验证尚未完成）。 |
| `[ ]` | 尚未开始。 |

## API stability tier (current v0.9.x candidate)

| Package | Tier | Promotion criterion |
|---|---|---|
| `pkg/finalizer` | **Stable** | v1 entry (no additional work). |
| `pkg/labels` | **Stable** | v1 entry (no additional work). |
| `pkg/status` | **Stable** | v1 entry (no additional work). |
| `pkg/storageclass` | **Stable** | Trivial validation surface (regex + nil check). |
| `pkg/version` (incl. `Matrix`) | Beta | Generic `Matrix[E]` cross-repo verify. |
| `pkg/monitoring` | Beta | `ServiceMonitor` cross-downstream equivalence e2e. |
| `pkg/networkpolicy` | Beta | 4-direction (ingress / egress × TCP / UDP) verify. |
| `pkg/security` | Beta | Restricted PSA guard across downstream. |
| `pkg/events` | Beta | Downstream live adoption + reconciliation regression 0. |
| `pkg/pvc` | Beta | Downstream PVC expansion live adoption. |
| `pkg/topology` | Beta | Downstream topology spread live adoption. |
| `pkg/webhook` | **Experimental** | Multi-downstream adoption + stabilization. |
| `pkg/probes` | **Experimental** | 2+ downstream adoption → Beta. |

**Tier 语义**：

- **Stable** — patch / minor release 中无 BREAKING CHANGE。使用
  deprecation：标记、保留 2 个 minor release、然后移除。
- **Beta** — minor release 中允许 BREAKING CHANGE（必须出现在
  CHANGELOG 中）。API 形状基本稳定。
- **Experimental** — 任何 release 都可能 BREAKING CHANGE。调用方
  自担风险。

## v1.0.0 graduation criteria（checklist）

- [ ] 所有包达到 **Stable** tier。
- [ ] 连续 6 个或更多 minor release 中 BREAKING CHANGE 为零。
- [ ] godoc 覆盖率 ≥ 80 %（`go doc -all ./...` basis）。
- [ ] CHANGELOG.md v0.x 演进历史 + v1.0.0 release notes。
- [ ] CITATION.cff + DOI（学术引用）。
- [ ] 下游 consumer 在 regression 0 下导入 v1.0.0 commons。
- [x] `go vet ./... && go test ./...` clean（覆盖率 96.3 % > 85 %
  阈值）。
- [x] API stability 承诺文档 — [STABILITY.md](STABILITY.zh.md)。
- Verify：下游 consumer CI 通过针对
  `keiailab-commons v1.0.0` 的所有 e2e 测试。

## Per-package follow-up

### pkg/finalizer (Stable)

- [x] `Add`、`Remove`、`Contains` helper — `pkg/finalizer/finalizer.go`。
- [x] 避免 controller-runtime（仅 stdlib `slices`）。
- [x] 单元测试 — `pkg/finalizer/finalizer_test.go`。
- [x] 多 finalizer 排序辅助 — `pkg/finalizer/order.go`
  `EnsureOrder`。
- Verify：下游 finalizer regression 0。

### pkg/labels (Stable)

- [x] Recommended Kubernetes labels（`app.kubernetes.io/*`）—
  `pkg/labels/labels.go`。
- [x] Component / instance / part-of 映射。
- [x] 单元测试 — `pkg/labels/labels_test.go`。
- [x] Recommended labels v2 (K8s 1.30+) — `pkg/labels/v2.go` `AllV2`。
- Verify：下游 `metadata.labels` 一致性。

### pkg/status (Stable)

- [x] Condition 目录 helper — `pkg/status/conditions.go`。
- [x] `SetAvailable` sugar。
- [x] 单元测试。
- [x] Condition reason 目录文档 —
  `pkg/status/REASONS.md`。
- Verify：下游 `kubectl get <kind> -o yaml`
  `.status.conditions` parity。

### pkg/version (Beta)

- [x] `Matrix[E]` 泛型 — `pkg/version/matrix.go`。
- [x] `SetAvailable` sugar。
- [x] semver-aware version 比较 — `pkg/version/version.go`。
- [x] 跨版本兼容性测试 —
  `pkg/version/api_stability_test.go`。
- [x] `Matrix[E]` 序列化器（JSON / YAML）—
  `pkg/version/serializer.go`。
- [ ] **Tier promotion** → Stable。
- Verify：下游 version 校验 parity。

### pkg/monitoring (Beta)

- [x] Prometheus ServiceMonitor builder — `pkg/monitoring/monitoring.go`。
- [x] 单元测试。
- [x] PrometheusRule builder（alert / recording 共享）—
  `pkg/monitoring/rule.go`。
- [x] OpenTelemetry exporter helper — `pkg/monitoring/otel.go`。
- [ ] 下游 equivalence e2e — 相同输入 → 相同 manifest 输出。
- [ ] **Tier promotion** → Stable。
- Verify：`monitoring_test.go` golden file diff = 0。

### pkg/networkpolicy (Beta)

- [x] NetworkPolicy builder — `pkg/networkpolicy/networkpolicy.go`。
- [x] Default-deny + 显式规则 helper。
- [x] 单元测试。
- [x] 4-direction 验证 —
  `pkg/networkpolicy/four_dir_test.go`。
- [x] CIDR + namespace + pod selector combo —
  `pkg/networkpolicy/combo.go`。
- [ ] **Tier promotion** → Stable。
- Verify：kind 集群应用 NetworkPolicy 并观察到的 deny /
  allow 路径与预期匹配。

### pkg/security (Beta)

- [x] SecurityContext helper（restricted PSA-compliant）—
  `pkg/security/security.go`。
- [x] RBAC helper。
- [x] 单元测试。
- [x] Restricted PSA regression guard —
  `pkg/security/psa_guard_test.go`。
- [x] Pod / Container SecurityContext 分离 —
  `pkg/security/split.go`。
- [x] seccompProfile 默认 helper — `pkg/security/seccomp.go`。
- [ ] **Tier promotion** → Stable。
- Verify：在 `kubectl label ns <ns>
  pod-security.kubernetes.io/enforce=restricted` 之后，下游 pod 到达
  Ready 状态。

### pkg/webhook (Experimental)

- [x] Webhook 工具基础 — `pkg/webhook/webhook.go`。
- [x] 单元测试。
- [x] Conversion webhook helper — `pkg/webhook/conversion.go`。
- [x] Validation webhook 模式 —
  `pkg/webhook/validation_patterns.go`。
- [ ] 多下游 live adoption → 稳定化。
- [ ] **Tier promotion** → Beta → Stable。
- Verify：2 个或更多下游 consumer 使用同一 helper 且
  regression 0。

### pkg/storageclass (Stable)

- [x] DNS-1123 subdomain 校验 —
  `pkg/storageclass/validator.go`。
- [x] Normalize / MustNormalize — 空值 → nil（cluster default）+ trim
  + 指针返回。
- [x] 12 个单元测试（100 % 覆盖率）—
  `pkg/storageclass/validator_test.go`。
- [ ] 下游 live adoption + regression 0。

### pkg/events (Beta)

- [x] Recorder interface — 不导入即可与 `client-go`
  `record.EventRecorder` 兼容。
- [x] 9 个 Reason 常量。
- [x] Emit / Emitf / EmitWarning / EmitWarningf / WrappedError — 全部
  nil-safe。
- [x] 单元测试（100 % 覆盖率）— `pkg/events/events_test.go`。
- [ ] 下游 live adoption — Event reason 在 reconcile 路径中
  统一。
- [ ] **Tier promotion** → Stable。
- Verify：下游 Reconcile 路径使用 commons reason 常量且
  regression 0。

### pkg/probes (Experimental)

- [x] Fluent builder — HTTP / HTTPS / TCP / Exec handler。
- [x] kubelet 默认值（Period = 10 s / Timeout = 1 s /
  SuccessThreshold = 1 / FailureThreshold = 3）。
- [x] InitialDelay / Period / Timeout 负值夹紧为 0。
- [x] `Build()` 在未设置 handler 时 panic（fail-fast 契约）。
- [x] 单元测试（100 % 覆盖率）— `pkg/probes/builder_test.go`。
- [ ] 2+ 下游 live adoption（Beta 标准）。
- [ ] **Tier promotion** → Beta → Stable。

### pkg/pvc (Beta)

- [x] PVC expansion helper — `pkg/pvc/expansion.go`。
- [x] 单元测试 — `pkg/pvc/expansion_test.go`。
- [ ] 下游 live adoption 且 PVC resize regression 0。
- [ ] **Tier promotion** → Stable。

### pkg/topology (Beta)

- [x] PVC topology spread helper — `pkg/topology/spread.go`。
- [x] 单元测试 — `pkg/topology/spread_test.go`。
- [ ] 下游 live adoption 且 spread constraint 验证。
- [ ] **Tier promotion** → Stable。

## 依赖政策

- **仅 Kubernetes API** — `k8s.io/api`、`k8s.io/apimachinery`、
  `k8s.io/utils`。controller-runtime 依赖*绝不在 leaf 包中
  添加*。
- **仅 permissive-license-compatible license** — 每次添加依赖都
  需要 ADR。
- **完整 godoc** — 每个新增 public API 都需要 godoc。

## 治理 / 追踪

- **CHANGELOG.md** — 由 `git-cliff` 自动生成。严格语义
  版本规范。
- **CITATION.cff** — 学术引用。DOI 在 v1.0.0 发放。
- **ADR** — `docs/kb/adr/` 追踪每个设计决策。
- **AGENTS.md** — AI 协作 runbook。

## Non-Goals（有意排除）

- ❌ **controller-runtime 依赖** — leaf 包必须保持
  controller-runtime free。
- ❌ **下游特定逻辑** — operator 特定代码位于
  调用方 repo 中。本库仅提供*共享* helper。
- ❌ **基于时间的 roadmap** — 使用 feature checklist 加上完成
  百分比。
- ❌ **GitHub Actions release gate** — 委托给本地四
  层。
- ❌ **Plugin / extension SDK 定位** — 这是 library，不是
  framework。
- ❌ **过早的 v1.0.0** — 在满足 graduation
  标准之前保持 v0.x。

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
