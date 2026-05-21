# ADR Index — operator-commons

> [English](INDEX.md) | [한국어](INDEX.ko.md) | [日本語](INDEX.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-tooling-unification-adoption.md) | 工具链统一 — `.golangci.yml` + `Makefile` 引入 | Accepted | 2026-05-09 |
| [0003](0003-pkg-status-finalizer-adoption.md) | `pkg/status` sugar (`SetAvailable` + `SetReadyFalse`) 添加 | Accepted | 2026-05-09 |
| [0004](0004-pkg-version-generic-matrix.md) | `pkg/version` 泛型 `Matrix[E]` 引入 | Accepted | 2026-05-09 |
| [0005](0005-library-chart-adoption.md) | `keiailab-commons` Helm library chart (commonLabels + ServiceMonitor) | Accepted | 2026-05-09 |
| [0006](0006-networkpolicy-partials.md) | Helm chart `keiailab.networkpolicy.{dataplane, controlplane}` partial | Accepted | 2026-05-09 |
| [0007](0007-rbac-partials.md) | Helm chart `keiailab.rbac.{serviceAccount, controllerBase, workloadBase}` partial | Accepted | 2026-05-09 |
| [0008](0008-security-partials.md) | Helm chart `keiailab.security.{podSecurityContext, containerSecurityContext}` partial | Accepted | 2026-05-09 |
| [0011](0011-lefthook-config-consolidation.md) | lefthook 配置统一（`.lefthook.yml` → `lefthook.yml`） | Accepted | 2026-05-21 |
| [0012](0012-gha-block-hook.md) | GitHub Actions 阻断 — lefthook pre-commit hook 自动强制 | Accepted | 2026-05-21 |
| [0014](0014-release-script-ssot.md) | `scripts/release.sh` — 手动 library release pipeline | Accepted | 2026-05-21 |
| [0015](0015-lefthook-augmentation.md) | lefthook 增强 — `go-licenses` + `markdown-link-check` | Accepted | 2026-05-21 |
| [0016](0016-pvc-topology-extraction.md) | `pkg/pvc` + `pkg/topology` 引入（下游 dedup） | Accepted | 2026-05-21 |

## 约定

- 文件名：`NNNN-<English kebab-case slug>.md`（四位数字，零填充；编号不重复使用）。
- 位置：`docs/kb/adr/`。
- 格式：Nygard 的五个 section（Context / Decision / Status / Consequences / Refs）。
- 状态机：Proposed → Accepted → (Deprecated | Superseded by ADR-XXXX)。
