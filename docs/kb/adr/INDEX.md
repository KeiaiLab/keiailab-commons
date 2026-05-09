# ADR Index — operator-commons

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-rfc-0017-tooling-unification-adoption.md) | RFC-0017 operator tooling unification 채택 (.golangci.yml + Makefile 신규) | Proposed | 2026-05-09 |
| [0003](0003-rfc-0018-pkg-status-finalizer-adoption.md) | RFC-0018 채택 — pkg/status 슈가 (SetAvailable + SetReadyFalse) 추가, pkg/finalizer 변경 없음 | Accepted | 2026-05-09 |
| [0004](0004-pkg-version-generic-matrix.md) | pkg/version Generic Matrix[E] 추가 (postgres PR-B3 prerequisite, Plan §2 D12) | Accepted | 2026-05-09 |
| [0005](0005-rfc-0019-library-chart-adoption.md) | RFC-0019 §3.1 채택 — keiailab-commons Helm library chart 신설 (commonLabels + ServiceMonitor, PR-B2 first cut) | Accepted | 2026-05-09 |
| [0006](0006-rfc-0019-networkpolicy-partials.md) | RFC-0019 §3.2 채택 — keiailab.networkpolicy.{dataplane,controlplane} partials (chart v0.2.0, PR-B6) | Accepted | 2026-05-09 |
| [0007](0007-rfc-0019-rbac-partials.md) | RFC-0019 §3.5 채택 — keiailab.rbac.{serviceAccount,controllerBase,workloadBase} partials (chart v0.3.0, PR-C1) | Accepted | 2026-05-09 |

## 작성 규약

- 파일명: `NNNN-<영어 kebab-case slug>.md` (4자리 0-padded, 재사용 금지)
- 위치: `docs/kb/adr/` (3 operator repo 표준 일치 — 본 디렉토리는 2026-05-09 `docs/adr/` 에서 이전됨)
- 형식: standards/adr.md §3 (Nygard 5섹션)
- 상태 머신: Proposed → Accepted → (Deprecated | Superseded by ADR-XXXX)
