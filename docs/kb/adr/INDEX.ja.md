# ADR Index — operator-commons

> [English](INDEX.md) | [한국어](INDEX.ko.md) | **日本語**

> ⚠️ This translation is AI-generated and pending native review.

| ID | タイトル | ステータス | 日付 |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-tooling-unification-adoption.md) | Tooling 統一 — `.golangci.yml` + `Makefile` 導入 | Accepted | 2026-05-09 |
| [0003](0003-pkg-status-finalizer-adoption.md) | `pkg/status` sugar (`SetAvailable` + `SetReadyFalse`) 追加 | Accepted | 2026-05-09 |
| [0004](0004-pkg-version-generic-matrix.md) | `pkg/version` generic `Matrix[E]` 導入 | Accepted | 2026-05-09 |
| [0005](0005-library-chart-adoption.md) | `keiailab-commons` Helm library chart (commonLabels + ServiceMonitor) | Accepted | 2026-05-09 |
| [0006](0006-networkpolicy-partials.md) | Helm chart `keiailab.networkpolicy.{dataplane, controlplane}` partial | Accepted | 2026-05-09 |
| [0007](0007-rbac-partials.md) | Helm chart `keiailab.rbac.{serviceAccount, controllerBase, workloadBase}` partial | Accepted | 2026-05-09 |
| [0008](0008-security-partials.md) | Helm chart `keiailab.security.{podSecurityContext, containerSecurityContext}` partial | Accepted | 2026-05-09 |
| [0011](0011-lefthook-config-consolidation.md) | lefthook 設定統合 (`.lefthook.yml` → `lefthook.yml`) | Accepted | 2026-05-21 |
| [0012](0012-gha-block-hook.md) | GitHub Actions ブロック — lefthook pre-commit hook 自動強制 | Accepted | 2026-05-21 |
| [0014](0014-release-script-ssot.md) | `scripts/release.sh` — 手動 library release パイプライン | Accepted | 2026-05-21 |
| [0015](0015-lefthook-augmentation.md) | lefthook augmentation — `go-licenses` + `markdown-link-check` | Accepted | 2026-05-21 |
| [0016](0016-pvc-topology-extraction.md) | `pkg/pvc` + `pkg/topology` 導入 (downstream 重複排除) | Accepted | 2026-05-21 |

## 規約

- ファイル名: `NNNN-<English kebab-case slug>.md` (4 桁、ゼロ埋め。
  番号は再利用しない)。
- 場所: `docs/kb/adr/`。
- 形式: Nygard の 5 セクション (Context / Decision / Status /
  Consequences / Refs)。
- ステータスマシン: Proposed → Accepted → (Deprecated | Superseded by
  ADR-XXXX)。
