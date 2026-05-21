# RFC-0007 (DRAFT): keiailab v3.0.1-stable 선언

> ⚠️ **DRAFT** — 본 문서는 `~/.codex/rfcs/0007-v3x-stable-declaration.md` 로
> push 전 commons 임시 draft. 사용자 명시 결정 후 ~/.codex/rfcs/ 에 실 push.

| Meta | Value |
|---|---|
| RFC | 0007 |
| Title | keiailab v3.0.1-stable 선언 |
| Status | Draft → (사용자 검토 후) Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | RFC-0001 (governance restoration), RFC-0002 (GitHub Actions Ban), CLAUDE.md §7 (v3.x-stable 정의) |

## 1. Summary

CLAUDE.md §7 의 *상용 제품 수준* 정의:

> "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 —
> `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable*
> 선언."

본 RFC 는 그 선언 — **2026-05-21 15:30, audit ❌ 0 달성** 시점에 5 repo
(postgres / mongodb / valkey / commons / forgewise) 모두 P0+P1+P2 자동화
충족 확인. **v3.0.1-stable 선언**.

## 2. Motivation

기존:
- CLAUDE.md §7 은 v3.x-stable 의 조건만 명시, 측정 방법 + 선언 절차 부재
- 5 repo 의 진척이 산발적 측정 — *자동화* 없이는 객관 평가 불가
- *언제* v3.x-stable 인지의 결정이 사람 판단 의존

본 RFC:
- 측정 자동화 (commons SSOT 의 `audit-production-grade.sh`) 의 결과 = audit ❌ 0
- 본 결과 → 선언 절차의 *객관 trigger*
- 사용자 결정 → 글로벌 governance update

## 3. Detailed Design

### 3.1 측정 결과 (2026-05-21 15:30)

`commons/scripts/audit-production-grade.sh /Users/phil/Workspace/keiailab`:

```
P0 — 기본 안전          7 항목 × 5 repo = 35 cells   ✅ 32 / ❌ 0 / — 3 (N/A)
P1 — 품질 게이트         5 항목 × 5 repo = 25 cells   ✅ 20 / ❌ 0 / — 5 (N/A)
P2 — 거버넌스           7 항목 × 5 repo = 35 cells   ✅ 35 / ❌ 0 / — 0
OP — 운영               6 항목 × 5 repo = 30 cells   ✅ 24 / ❌ 0 / — 6 (N/A)
C — 커뮤니티             2 항목 × 5 repo = 10 cells   ✅ 10 / ❌ 0 / — 0
─────────────────────────────────────────────────────────────────────
                       총                              ✅ 121 / ❌ 0 / — 14
```

**❌ 0 = 모든 측정 항목 통과**.

### 3.2 5 repo 정합 표

| repo | role | audit ❌ | 핵심 ADR | release tag (제안) |
|---|---|---|---|---|
| postgres-operator | PostgreSQL 18+ operator | 0 ✅ | postgres-ADR/0017/19/20/21/22 | v0.3.0 (또는 v1.0.0) |
| mongodb-operator | MongoDB 7.0+ operator | 0 ✅ | mongodb-ADR/0031/33/35 | v1.5.0 |
| valkey-operator | Valkey 8.0+ operator | 0 ✅ | valkey-ADR/0048/50 | v1.0.0 |
| operator-commons | Shared Go library + audit SSOT | 0 ✅ | commons ADR-0009~0016 | v0.9.0 |
| forgewise | GitLab MCP server (Python) | 0 ✅ | forgewise ADR-0001 | v0.1.0 |

### 3.3 v3.0.1-stable 의 의미

- v3.0 = CLAUDE.md 규약의 본 major
- .1 = audit 자동화 첫 충족 (2026-05-21)
- stable = 외부 사용자 대상 *운영 가능 등급*
- 후속 변화 시:
  - audit ❌ 발생 → *unstable* 표시 (자동), 회귀 차단 cycle 진입
  - P3 (성능) / P4 (DR) / P5 (커뮤니티) 신설 시 → v3.1+ 진화

### 3.4 선언 효과

- 5 repo README 상단에 v3.x-stable 배지 (S8 Phase 4)
- 각 repo `docs/kb/adr/<N>-v3x-stable-baseline.md` Accepted (S8 Phase 2)
- commons `docs/family.md` 5 sister 인정 + 각 v3.x-stable 표시
- 외부 신뢰 신호 (audit + OpenSSF Scorecard + 거버넌스 4종 + i18n 4-lang)

### 3.5 회귀 차단

- audit 결과는 *지속 유지* 의무
- 후속 PR 의 audit 결과 ❌ ≥ 1 = 회귀 → 별 cycle 필수
- audit 자동 측정 trigger:
  - 월 1회 cron (후속 RFC)
  - 5 repo 의 release 직전 (각 repo `make release` 안 통합)

## 4. CLAUDE.md §7 갱신 (제안)

현재 본문:

> "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 —
> `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable*
> 선언."

갱신 안 (본 RFC Accepted 시):

> "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 —
> `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable*
> 선언.
>
> **2026-05-21 v3.0.1-stable 선언** (RFC-0007). 5 repo (postgres / mongodb /
> valkey / commons / forgewise) 모두 audit ❌ 0. 측정: `keiailab/operator-commons`
> `scripts/audit-production-grade.sh`. 시계열 + 정합: `docs/quality/audit-history.md`."

## 5. Drawbacks

- 본 RFC 의 *외부 announce* 영향 — keiailab community 의 v3.x-stable 기대 ↑
  → 회귀 시 신뢰 손상 위험. 회귀 차단 자동화 필수
- v1.0 tag 의 *의미* — postgres 가 현 v0.3.0-alpha.18 이라 v1.0.0 jump 큼.
  v0.3.0 stable 로 시작 후 점진적 v1.0 도 옵션
- 5 repo 의 비대칭 (postgres 의 GHA 3 narrow vs mongodb 의 12 vs valkey 의 14)
  → "v3.x-stable" 가 *audit* 만 정합, *내부 운영* 정합은 아님.
  → 다음 cycle 의 RFC 로 *통합 정책* 검토 가능

## 6. Alternatives

### 6.1 단계적 stable 선언
- v3.x-stable 가 아닌 v3.x-beta → v3.x-rc → v3.x-stable 단계 거침
- 장점: 신중. 회귀 시 *less surprising*
- 단점: CLAUDE.md §7 의 *audit 100% = stable* 정의 위배. 본 RFC 가 그 정의의 첫 실현

### 6.2 5 repo 부분 stable
- 5 repo 중 일부만 v3.x-stable (예: commons + forgewise 만, operator 3개는 v3.x-rc)
- 장점: 운영 차이 (3 operator 의 GHA 정책 비대칭) 반영
- 단점: keiailab family 통일성 깨짐. audit 결과 ❌ 0 인데 분리할 사유 부족

### 6.3 선언 보류
- audit ❌ 0 도달 + 1 주 모니터링 후 선언
- 장점: 회귀 위험 ↓ (early detection)
- 단점: *선언 trigger* 가 시간 의존 — 자동화 정합 깨짐. 회귀는 자동화로 차단

→ **본 RFC 채택 시 6.1/6.2/6.3 모두 거부, 즉시 v3.0.1-stable 선언**.

## 7. Unresolved Questions

- 5 repo release tag 의 정확 버전 (v1.0.0 vs v0.X.Y) → 사용자 결정 (RFC 본문 §3.2 의 제안 참조)
- 후속 v3.1+ 의 P3/P4/P5 신설 trigger → 별 RFC
- audit 자동 측정 cron 정책 → 별 RFC

## 8. Migration

본 RFC Accepted 시:
1. `~/.codex/CLAUDE.md` §7 갱신 (위 §4 본문)
2. `~/.codex/rfcs/0007-v3x-stable-declaration.md` 실 push (본 draft 의 promote)
3. 5 repo 의 `docs/kb/adr/<N>-v3x-stable-baseline.md` 신설 (S8 Phase 2 진행 중)
4. 5 repo README 의 v3.x-stable 배지 (S8 Phase 4 진행 중)
5. 5 repo release tag (별 사용자 결정 — S8 Phase 3)
6. commons `docs/family.md` 5 sister 갱신 (S8 Phase 4)
7. (선택) external announce

## 9. References

- CLAUDE.md §7 (v3.x-stable 정의)
- standards/enforcement.md (P0+P1+P2 자동화)
- commons audit script: `commons/scripts/audit-production-grade.sh`
- audit history: `commons/docs/quality/audit-history.md`
- S8 execution plan: `commons/docs/specs/2026-05-21-v3x-stable-declaration-execution-plan.md`
- HANDOFF v0.4: `commons/HANDOFF.md`
