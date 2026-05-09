# ADR Index — operator-commons

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-charter.md) | operator-commons charter | Accepted | 2026-05-07 |
| [0002](0002-rfc-0017-tooling-unification-adoption.md) | RFC-0017 operator tooling unification 채택 (.golangci.yml + Makefile 신규) | Proposed | 2026-05-09 |
| [0003](0003-rfc-0018-pkg-status-finalizer-adoption.md) | RFC-0018 채택 — pkg/status 슈가 (SetAvailable + SetReadyFalse) 추가, pkg/finalizer 변경 없음 | Accepted | 2026-05-09 |

## 작성 규약

- 파일명: `NNNN-<영어 kebab-case slug>.md` (4자리 0-padded, 재사용 금지)
- 위치: `docs/kb/adr/` (3 operator repo 표준 일치 — 본 디렉토리는 2026-05-09 `docs/adr/` 에서 이전됨)
- 형식: standards/adr.md §3 (Nygard 5섹션)
- 상태 머신: Proposed → Accepted → (Deprecated | Superseded by ADR-XXXX)
