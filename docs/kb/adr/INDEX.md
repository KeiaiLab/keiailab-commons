# ADR Index — operator-commons

> **English** | [한국어](INDEX.ko.md) | [日本語](INDEX.ja.md)

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-tooling-unification-adoption.md) | Tooling unification — `.golangci.yml` + `Makefile` introduction | Accepted | 2026-05-09 |
| [0003](0003-pkg-status-finalizer-adoption.md) | `pkg/status` sugar (`SetAvailable` + `SetReadyFalse`) addition | Accepted | 2026-05-09 |
| [0004](0004-pkg-version-generic-matrix.md) | `pkg/version` generic `Matrix[E]` introduction | Accepted | 2026-05-09 |
| [0005](0005-library-chart-adoption.md) | `keiailab-commons` Helm library chart (commonLabels + ServiceMonitor) | Accepted | 2026-05-09 |
| [0006](0006-networkpolicy-partials.md) | Helm chart `keiailab.networkpolicy.{dataplane, controlplane}` partial | Accepted | 2026-05-09 |
| [0007](0007-rbac-partials.md) | Helm chart `keiailab.rbac.{serviceAccount, controllerBase, workloadBase}` partial | Accepted | 2026-05-09 |
| [0008](0008-security-partials.md) | Helm chart `keiailab.security.{podSecurityContext, containerSecurityContext}` partial | Accepted | 2026-05-09 |
| [0011](0011-lefthook-config-consolidation.md) | lefthook configuration consolidation (`.lefthook.yml` → `lefthook.yml`) | Accepted | 2026-05-21 |
| [0012](0012-gha-block-hook.md) | GitHub Actions block — lefthook pre-commit hook auto-enforcement | Accepted | 2026-05-21 |
| [0014](0014-release-script-ssot.md) | `scripts/release.sh` — manual library release pipeline | Accepted | 2026-05-21 |
| [0015](0015-lefthook-augmentation.md) | lefthook augmentation — `go-licenses` + `markdown-link-check` | Accepted | 2026-05-21 |
| [0016](0016-pvc-topology-extraction.md) | `pkg/pvc` + `pkg/topology` introduction (downstream dedup) | Accepted | 2026-05-21 |

## Conventions

- File name: `NNNN-<English kebab-case slug>.md` (four digits, zero-padded; numbers are not reused).
- Location: `docs/kb/adr/`.
- Format: Nygard's five sections (Context / Decision / Status / Consequences / Refs).
- Status machine: Proposed → Accepted → (Deprecated | Superseded by ADR-XXXX).
