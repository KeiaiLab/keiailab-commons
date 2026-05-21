# ADR-0013: audit-production-grade.sh — 5 repo SSOT 측정 자동화

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | CLAUDE.md §7 (v3.x-stable 정의), ADR-0012 (gha-block hook), standards/enforcement.md (P0/P1/P2) |

## Context

CLAUDE.md §7: "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 — `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable* 선언."

5 repo (postgres / mongodb / valkey / commons / forgewise) 의 *상용 제품 수준* 측정이 그동안 *수동 의지* 에 의존. 본 ADR 은 다음을 정합:

- 5 축 정의: P0 (기본 안전) + P1 (품질 게이트) + P2 (거버넌스) + OP (운영) + C (커뮤니티)
- 측정 자동화: `scripts/audit-production-grade.sh`
- 기준 SSOT: `docs/quality/production-grade-checklist.md` (8.5 KB, ~50 항목)
- Make target: `make audit-quality`

본 audit 결과는 후속 sub-cycle 의 *우선순위 입력*. 모든 항목 ✅ = v3.x-stable 선언 조건.

## Decision

operator-commons 가 5 repo audit 의 SSOT:
- `scripts/audit-production-grade.sh` — bash 스크립트, 5 repo parent dir (`/Users/phil/Workspace/keiailab`) parameter, 표 형식 markdown 출력
- `docs/quality/production-grade-checklist.md` — 5 축 ~50 항목 정의, 측정 명령, 임계값
- `Makefile: audit-quality` — local 측정 entry point
- 4 repo (postgres / mongodb / valkey / forgewise) 는 `scripts/sync-from-commons.sh` 통해 audit 스크립트 sync (별 sub-cycle)

## Consequences

- ✅ v3.x-stable 선언 조건의 *측정 자동화* — 사람 판단 의존 제거
- ✅ 후속 sub-cycle (S5 공통화 / S4-A/B/C/D 다국어 / S9 GHA block hook 4 repo / S10 upgrade guide / S8 최종 audit) 의 *우선순위 입력*
- ⚠️ audit 항목 추가 시 ADR 갱신 (예: P3 성능 게이트, P4 DR 게이트)
- ⚠️ N/A (`—`) 항목 (예: forgewise/commons 에서 kube-linter, OLM bundle) 의 정합 — checklist.md 의 "(operator only)" 마킹

## 사용

```bash
# 5 repo 측정
make audit-quality
# 또는 직접
bash scripts/audit-production-grade.sh
# 또는 다른 parent dir
bash scripts/audit-production-grade.sh /path/to/parent
```

출력 (markdown 표):
- P0 / P1 / P2 / OP / C 5 섹션
- 각 섹션: 항목 × 5 repo 매트릭스
- ✅ / ❌ / `—` (N/A) / `N/M` (multi-language count)

상세 저장: `/tmp/audit-production-grade-<timestamp>.md`.

## 첫 baseline (2026-05-21 14:03)

5 cycle 완료 시점 측정:
- commons: P0 ✅×7 / P1 ✅×2 (Go) / P2 ✅×7 (P2-2 ✅) / OP ✅×4 + ❌×2 / C ✅ (4/4 BRANDING, 4/4 README)
- forgewise: P0 ✅×4 ❌×3 (gitleaks/Conventional/mod-drift 부재 — Python 환경 차이) / P2 ✅×2 ❌×4 (BRANDING/family/CODEOWNERS/PR template) — S4-D + 별 sub-cycle 대기
- postgres + mongodb: WF 0 (S7 결과, v2.0 revert 대기) → ADR-0019/0033 신설 후 ✅
- valkey: WF 14 유지 (ralph-loop 관리), P1-11/12/13 lefthook 3종 보강 대기
