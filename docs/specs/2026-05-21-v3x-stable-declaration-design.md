# S8 Design Draft: 최종 audit + v3.x-stable 선언

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Draft (5 cycle + S5 완료 후 진입) |
| 위치 | `~/.codex/rfcs/0005-v3x-stable-declaration.md` (글로벌 RFC) + 각 repo ADR |
| 의존 | S1+, S2, S4, S6, S7, S5 모두 완료 |

## 1. 목적

CLAUDE.md §7: "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 — `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable* 선언."

S8 = 그 *선언* cycle. 모든 항목 ✅ 확인 → ADR + RFC 작성 → main commit + tag.

## 2. 입력

- `production-grade-checklist.md` (5축, 약 50 항목)
- `audit-production-grade.sh` (자동 측정)
- 5 cycle + S5 의 산출

## 3. Phase

### Phase 0 — 사전 확인
- 모든 cycle 완료 확인 (task list)
- audit script 실행 → 첫 baseline 측정
- gap 식별

### Phase 1 — gap 해소 sub-cycle
- 측정 결과의 ❌ 항목 별 sub-cycle dispatch
- 예시:
  - `P2-7 거버넌스 4종` 미충족 repo → 해당 repo 에 작성 (4-lang 포함)
  - `OP-5 ADOPTERS.md` 미충족 → 작성
  - `C-1 README 4-lang` 미충족 (valkey) → S4-C sub-cycle 결과 적용

### Phase 2 — 최종 audit
- audit script 재실행
- 모든 항목 ✅ 확인
- 결과 표 → `docs/quality/audit-2026-05-21.md` (commons SSOT) 저장

### Phase 3 — RFC + ADR 작성
- `~/.codex/rfcs/0005-v3x-stable-declaration.md` (Status: Accepted)
  - 본문: 5 repo 의 P0+P1+P2 자동화 충족 측정 결과 인용
  - audit script 명령 + 결과 인용
- 각 repo `docs/kb/adr/<N>-v3x-stable-baseline.md` (Status: Accepted)

### Phase 4 — release tag
- 5 repo 각각 `vX.Y.Z` tag (semver 정합)
- 예시:
  - postgres-operator v1.0.0
  - mongodb-operator v1.4.4 → v1.5.0
  - valkey-operator v0.x.y → v1.0.0
  - operator-commons v0.9.0 → v1.0.0
  - forgewise v0.1.0
- `make release` 실행 (S1+/S7 Phase 2.5 의 로컬 스크립트)
- GitHub Release 본문 자동 생성 (cliff)

### Phase 5 — 가족 메타 문서 갱신
- commons `docs/quality/production-grade-checklist.md` (`$CLAUDE_JOB_DIR` 에서 commons 로 이전)
- commons `docs/family.md` (5 repo 의 v1.0 + audit 통과 표시)
- 다른 4 repo 의 `docs/family.md` sync (commons SSOT)

### Phase 6 — CLAUDE.md §7 수정
- `~/.codex/CLAUDE.md` 의 §7 의 "v3.x-stable 선언" 부분:
  - 현재: "P0+P1+P2 자동화 모두 충족 시 *v3.x-stable* 선언"
  - 갱신 후: "**2026-05-21 v3.0.1-stable 선언** (RFC-0005). audit 결과: `commons/docs/quality/audit-2026-05-21.md`."
- commit + push (글로벌 standards repo)

### Phase 7 — 5 repo 의 main README 에 v3.x-stable 배지
- 각 README 상단에 `[![keiailab v3.x-stable](https://...)](commons/quality/audit-2026-05-21.md)` 배지
- 4-lang README 모두 동일

## 4. Success Criteria

```bash
# 1. audit script ✅ 100% (—제외)
./audit-production-grade.sh | grep -cE '❌' 
# 결과 = 0

# 2. RFC-0005 Accepted
grep -A1 'Status:' ~/.codex/rfcs/0005-v3x-stable-declaration.md | grep -q Accepted

# 3. 5 repo 모두 v1.0 tag (또는 정합 semver)
for repo in ...; do
  git -C $repo tag -l 'v*' | tail -1 | grep -q '^v[1-9]' || exit 1
done

# 4. 가족 메타 sync
# (commons/docs/family.md vs 다른 4 repo)

# 5. CLAUDE.md §7 갱신
grep -q 'v3.0.1-stable 선언' ~/.codex/CLAUDE.md
```

## 5. 후속 (v3.1+)

- v3.x-stable 유지 cycle (월 1회 audit)
- P3 (성능 게이트) / P4 (DR) / P5 (커뮤니티 KPI) 신설 RFC
- v3.1-stable → v3.2-stable 진화
