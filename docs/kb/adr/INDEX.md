# ADR Index — operator-commons

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-rfc-0017-tooling-unification-adoption.md) | Tooling unification — `.golangci.yml` + `Makefile` 도입 | Accepted | 2026-05-09 |
| [0003](0003-rfc-0018-pkg-status-finalizer-adoption.md) | `pkg/status` 슈가 (`SetAvailable` + `SetReadyFalse`) 추가 결정 | Accepted | 2026-05-09 |
| [0004](0004-pkg-version-generic-matrix.md) | `pkg/version` 의 generic `Matrix[E]` 도입 | Accepted | 2026-05-09 |
| [0005](0005-rfc-0019-library-chart-adoption.md) | `keiailab-commons` Helm library chart 신설 (`commonLabels` + ServiceMonitor) | Accepted | 2026-05-09 |
| [0006](0006-rfc-0019-networkpolicy-partials.md) | Helm chart `keiailab.networkpolicy.{dataplane, controlplane}` partial | Accepted | 2026-05-09 |
| [0007](0007-rfc-0019-rbac-partials.md) | Helm chart `keiailab.rbac.{serviceAccount, controllerBase, workloadBase}` partial | Accepted | 2026-05-09 |
| [0008](0008-rfc-0019-security-partials.md) | Helm chart `keiailab.security.{podSecurityContext, containerSecurityContext}` partial | Accepted | 2026-05-09 |
| [0011](0011-lefthook-config-consolidation.md) | lefthook 설정 통합 (`.lefthook.yml` → `lefthook.yml`) | Accepted | 2026-05-21 |
| [0012](0012-rfc-0002-gha-block-hook.md) | GitHub Actions 차단 — lefthook pre-commit hook 자동 강제 | Accepted | 2026-05-21 |
| [0014](0014-release-script-ssot.md) | `scripts/release.sh` — 라이브러리 수동 release pipeline | Accepted | 2026-05-21 |
| [0015](0015-lefthook-p1-12-13-augmentation.md) | lefthook 보강 — `go-licenses` + `markdown-link-check` | Accepted | 2026-05-21 |
| [0016](0016-sprint-1-pvc-topology-extraction.md) | `pkg/pvc` + `pkg/topology` 신규 추출 (downstream 중복 해소) | Accepted | 2026-05-21 |

## 작성 규약

- 파일명: `NNNN-<영어 kebab-case slug>.md` (4자리 0-padded, 재사용 금지).
- 위치: `docs/kb/adr/`.
- 형식: Nygard 5섹션 (Context / Decision / Status / Consequences / Refs).
- 상태 머신: Proposed → Accepted → (Deprecated | Superseded by ADR-XXXX).
