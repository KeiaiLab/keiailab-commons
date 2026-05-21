# 5 Background Cycle 완료 후 즉시 dispatch sub-cycle 계획

| 메타 | 값 |
|---|---|
| 입력 | audit-baseline-1341.md (현재 gap 38건) |
| 산출 | 후속 sub-cycle 4개의 dispatch prompt + 우선순위 |
| 사용 | 5 background subagent 완료 알림 즉시 main thread 가 일괄 dispatch |

## 우선순위 분류

| 단계 | sub-cycle | 의존 | 예상 LOC | 병렬? |
|---|---|---|---|---|
| **Wave A** (병렬 가능) | S9 P2-2 GHA block hook (5 repo) | 5 cycle 완료 | 작음 | ✅ |
| Wave A | S10 OP-10 upgrade guide (4 repo) | 5 cycle 완료 | 중간 | ✅ |
| Wave A | S11 audit script 5 repo 배포 | 5 cycle 완료 | 작음 | ✅ |
| **Wave B** (Wave A 후) | S4-A postgres 다국어 | S4 SSOT 완료 + Wave A | 큼 | ✅ (4개 동시) |
| Wave B | S4-B mongodb 다국어 | 동일 | 큼 | ✅ |
| Wave B | S4-C valkey 다국어 | 동일 + S1+ 완료 | 큼 | ✅ |
| Wave B | S4-D forgewise 다국어 | 동일 + S6 완료 | 작음 | ✅ |
| **Wave C** (Wave B 후) | **S5 공통화** (s5-spec-draft 진입) | Wave B 모두 + 코드 stable | 매우 큼 | sequential |
| **Wave D** | **S8 v3.x-stable 선언** | S5 + Wave B + audit script ≥ 100% ✅ | 작음 | sequential |

## Sub-cycle dispatch prompt — Wave A

### S9 — P2-2 GHA block hook (5 repo, 병렬 dispatch)

각 repo 별 sub-subagent 1개씩 = 5 subagent. 또는 1개가 5 repo 처리. 5 repo 의 lefthook 패턴이 유사하므로 1개로 충분.

```
keiailab 5 repo (postgres / mongodb / valkey / commons / forgewise) 의 lefthook 에
RFC-0002 강제 pre-commit hook 추가.

각 repo 별 branch + commit + push + PR + merge.

hook 내용 (lefthook.yml 또는 .lefthook.yml 의 pre-commit commands 에 추가):

  gha-block:
    glob: ".github/workflows/*.{yml,yaml}"
    run: |
      echo "❌ RFC-0002 위반: .github/workflows/ 디렉토리 추가 금지"
      echo "   대체: 로컬 4계층 (lefthook + Makefile + 리뷰어 증거)"
      echo "   예외 신청: ADR 작성 후 사용자 명시 승인"
      exit 1

검증: 
1. echo "test" > .github/workflows/test.yml; git add .github/workflows/test.yml; lefthook run pre-commit
2. 출력에 "RFC-0002 위반" 포함 + exit 1
3. rm .github/workflows/test.yml; rmdir .github/workflows (시작 상태 복원)

5 repo 각각:
- postgres-operator
- mongodb-operator (S7 후 lefthook.yml 안정 시)
- valkey-operator (S1+ 후 안정 시)
- operator-commons (이미 lefthook.yml 통합 완료)
- forgewise (S6 후 안정 시)

각 PR title: "feat(ci): RFC-0002 — gha-block pre-commit hook 강제 (P2-2)"
```

### S10 — OP-10 upgrade guide (4 repo, 병렬)

```
keiailab 4 repo (postgres / valkey / commons / forgewise) 의 docs/UPGRADING.md 작성.
mongodb 는 이미 있음 (audit 결과 ✅).

각 repo:
- branch: docs/upgrading-guide-2026-05-21
- docs/UPGRADING.md 작성:
  * "## Upgrading [latest version]" + "## Upgrading from vN-1 to vN" 양식
  * 현재 stable 버전 + 마이그레이션 단계 + breaking change 안내
  * mongodb 의 UPGRADING.md 를 reference (이미 존재)
- 영어 canonical + ko/ja/zh 는 placeholder (S4-* sub-cycle 에서 번역)
- commit + push + PR + merge

PR title: "docs(upgrade): UPGRADING.md 추가 (OP-10) — semver 마이그레이션 가이드"
```

### S11 — audit script 5 repo 배포

```
$CLAUDE_JOB_DIR/audit-production-grade.sh 를 operator-commons 의 SSOT 으로 이전 +
4 repo 에 sync.

1. commons:
   - branch: feat/audit-production-grade-ssot-2026-05-21
   - scripts/audit-production-grade.sh 추가 (job dir 의 파일)
   - docs/quality/production-grade-checklist.md 추가 ($CLAUDE_JOB_DIR 의 파일)
   - Makefile audit-quality target 추가
   - commit + PR + merge

2. 4 repo (postgres/mongodb/valkey/forgewise) 에:
   - scripts/sync-from-commons.sh 통해 audit-production-grade.sh 동기화
   - 또는 각 repo 가 자체 사본 (drift 위험) → SSOT cp 권장
   - Makefile audit-quality target 추가 (commons 패턴 동일)
   - commit + PR + merge

PR title: "feat(audit): audit-production-grade.sh SSOT + 5 repo 배포"
```

## Sub-cycle dispatch prompt — Wave B (S4 sub-cycle)

각 repo 별 다국어 적용. S4 spec PR #39 (commons SSOT 완료 후) 기반.

### S4-A postgres-operator 다국어

```
postgres-operator 의 다국어 4-lang 적용. commons SSOT 의 glossary + sync hook 활용.

branch: docs/i18n-4lang-postgres-2026-05-21

번역 대상 (audit 결과 + spec PR #39 분석):
- README.md (EN canonical 이미) → README.ko/ja/zh.md 자동 번역 (Claude 직접)
- BRANDING.md → BRANDING.{ko,ja,zh}.md 4-lang 작성
- docs/family.md → family.{ko,ja,zh}.md
- 운영 문서 (operations/runbooks/getting-started 등) 17건 → 4-lang
- 각 파일 상단에 ⚠️ `> This translation is AI-generated and pending native review. See [glossary](https://github.com/keiailab/operator-commons/blob/main/docs/i18n/) for terminology.`

번역 절차:
1. 영어 canonical 읽기
2. glossary 의 ja/zh 용어 lookup
3. Claude 가 직접 번역 (현재 세션)
4. ⚠️ 마킹 + warning 배너 추가
5. commit 단위: 1 파일 = 1 commit 또는 5-10 파일 = 1 commit (atomic 양식)

검증:
- commons 의 check-readme-sync.sh 4-lang 통과
- markdown-link-check 통과
- lefthook pre-push 통과

commit + push + PR (큰 변경 → 분할 머지 가능) + merge
```

### S4-B mongodb-operator 다국어 (S4-A 와 동일 패턴)
### S4-C valkey-operator 다국어 (동일, 단 valkey 만 README 2/4 → 4/4 보강 필수)
### S4-D forgewise 다국어 (이미 README 4/4 골격 — 본문 확장만)

## Sub-cycle dispatch prompt — Wave C (S5)

본 cycle 의 dispatch 는 *5 cycle + Wave A + Wave B* 모두 완료 후. spec 입력: `$CLAUDE_JOB_DIR/s5-spec-draft.md`.

```
S5 공통화 — s5-spec-draft.md 의 4 Phase + 5 Wave 실행.

먼저 commons 에 spec 머지 (job dir 파일 → docs/specs/2026-05-21-consolidation-design.md 로 이전).
사용자 검토 1 round (5 default 결정 확인).

Phase 1-4 = commons 패키지 추출 (pkg/reconcile + pkg/resources + pkg/storageclass 확장 + pkg/apis/common).
Wave 2 = commons v0.9.0 release.
Wave 3 = 3 operator sequential adoption.
Wave 4 = 9 미사용 helper 적용.
Wave 5 = audit script 재실행 → 결과 인용.

각 Wave 종료 시 audit script 실행 + 결과 인용. P1/P2/OP 항목 진척 측정.
```

## Sub-cycle dispatch prompt — Wave D (S8)

본 cycle 의 dispatch 는 *S5 완료 + audit script 100% ✅* 시점.

```
S8 최종 audit + v3.x-stable 선언 — s8-spec-draft.md 의 7 Phase 실행.

핵심: audit script ≥ 100% ✅ 검증 → RFC-0005 작성 → CLAUDE.md §7 갱신 → 5 repo v1.0 tag.

산출:
- ~/.codex/rfcs/0005-v3x-stable-declaration.md (Accepted)
- 5 repo docs/kb/adr/<N>-v3x-stable-baseline.md (Accepted)
- 5 repo vX.0.0 tag
- ~/.codex/CLAUDE.md §7: "2026-05-DD v3.0.1-stable 선언"
- commons docs/family.md 갱신 (v3.x-stable 배지)
```

## 5 cycle 완료 알림 시 즉시 main thread 작업

1. audit script 재실행 → 5 cycle 후 baseline
2. 비교: 이전 baseline vs 후속 baseline
3. Wave A 3 sub-cycle (S9 + S10 + S11) 병렬 dispatch
4. Wave A 완료 알림 시 Wave B (S4-A/B/C/D) 4 sub-cycle 병렬 dispatch
5. Wave B 완료 알림 시 Wave C (S5) dispatch
6. S5 완료 알림 시 Wave D (S8) dispatch
7. 모든 ✅ 시점 = goal 달성 = stop hook condition 충족
