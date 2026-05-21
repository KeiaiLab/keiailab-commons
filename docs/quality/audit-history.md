# keiailab Production-Grade Audit History

본 문서는 5 repo (postgres-operator / mongodb-operator / valkey-operator /
operator-commons / forgewise) 의 `commons/scripts/audit-production-grade.sh`
측정 결과 시계열. CLAUDE.md §7 의 *v3.x-stable* 선언 조건 (P0+P1+P2+OP+C
모두 ✅) 의 진척 추적.

## 측정 명령

```bash
make audit-quality
# 또는
bash commons/scripts/audit-production-grade.sh /path/to/keiailab/parent
```

## 시계열

### 2026-05-21 13:30 — baseline (5 cycle 시작 전)

❌ count: **38** (38/~60 ❌)

주요 ❌:
- 5 repo 의 거버넌스 4종 다수 부재
- valkey 21 open PR + 6 stale branch + 14 GHA workflow
- BRANDING 4-lang: 모든 repo 1/4 (또는 0)
- gha-block hook: 5 repo 모두 ❌

### 2026-05-21 14:03 — 5 cycle 일부 완료 후

❌ count: **16** (22 ✅ 추가)

주요 진척:
- S6 forgewise: 거버넌스 5종 + 운영 4종 + lefthook + README 4-lang
- S7 postgres + mongodb: workflow 0 + lefthook 보강 + postgres-ADR/0018/0032
- S2 commons: archive tag + lefthook 통합
- S4 commons SSOT: glossary 4-lang + sync hook + 자체 docs 번역

### 2026-05-21 14:20 — S7 revert 진행 + commons 보강 + S4-D 진행 중

❌ count: **15** → **11** → **9** → **7**

주요 변경:
- commons 의 6 commit (gha-block + audit SSOT + UPGRADING + release + S5/S8 spec + lefthook P1-12/13)
- audit script v2.0 정합 갱신 (P0-6 의 ADR 인정)
- audit script keyword 정밀화 (P0-2 detect-secrets, P0-4 commitlint, P0-9 uv lock)
- forgewise 4 신규 (CODEOWNERS / PR template / ADOPTERS / ROADMAP)
- forgewise lefthook 보강 (uv-lock-drift + markdown-link-check) + release.sh
- S7 revert subagent: postgres + mongodb workflow 복원 + postgres-ADR/0019/0033 Accepted

### 2026-05-21 14:40 — S4-D 완료 + S9 subagent 진행 중

❌ count: **7**

주요 진척:
- S4-D forgewise: 5 PR (BRANDING/family/README ja+zh/ops 4종/i18n hook) + 거버넌스 5종 i18n
- forgewise C-1/C-2 4/4 ✅
- 약 50 신규 파일 (BRANDING 4 + family 4 + README 4 + ops 12 + 거버넌스 15 + commons SSOT 6 + lefthook 1)

### 남은 7 ❌ 분류 (2026-05-21 14:40)

| 분류 | 항목 | 책임 |
|---|---|---|
| **3 operator P2-2 GHA block hook** | postgres ❌ mongodb ❌ valkey ❌ | S9 sub-cycle 진행 중 |
| **valkey P1-11/12/13 + OP-2** | kube-linter + go-licenses + markdown-link-check + helm-publish | ralph-loop 관할 |
| **postgres OP-10 upgrade guide** | UPGRADING.md push 실패 | S9 sub-cycle 의 postgres push fix |
| **valkey OP-10** | UPGRADING.md 미작성 | ralph-loop 관할 |

## 정합 충돌 사례 (자동화 간)

### postgres workflow 변동
- PR #86 (S7 cycle): 14 workflow 제거 (RFC-0002 strict)
- PR #89 (S7 cycle): postgres-ADR/0018 Accepted (RFC-0002 strict)
- PR #90 (S7 revert): 14 workflow 복원 (사용자 v2.0 결정)
- PR #92 (S7 revert): postgres-ADR/0019 Accepted (GHA 유지)
- PR #93 (다른 자동화): 11 workflow 제거 (§7 narrow exception 3종 보존)
- → 현 main 상태: 3 workflow (helm-publish + release + scorecard)
- → postgres-ADR/0019 (GHA 14 유지) 와 실 상태 (3 유지) 일관성 깨짐 — *별 ADR 또는 0019 갱신 필요*

### mongodb workflow 변동
- PR #194-#197 (S7 cycle): 12 workflow 제거 + mongodb-ADR/0032
- PR #199 (S7 revert): 12 workflow 복원
- PR #200 (다른 자동화): 9 workflow 재제거 (§7 narrow exception)
- PR #203 (S7 revert): mongodb-ADR/0033 Accepted (12 유지)
- PR #204 (S7 revert): 9 workflow 재복원 (race recovery)
- → 현 main 상태: 12 workflow 유지 (정합)

### 교훈
- **자동화 정합성**: 동일 저장소에 *여러 자동화* (본 thread + S7 + ralph-loop + 다른 cron-style automation) 가 *동시 작업* 시 정책 충돌 발생. 후속 cycle 의 *동시성 제어* RFC 필요.
- **권장 패턴**: 1 저장소 = 1 자동화 + 사용자 명시 결정에 의한 직접 commit 만 추가.

## 다음 단계 (v3.x-stable 선언까지)

1. **S9 sub-cycle 완료** — postgres + mongodb P2-2 ✅ + postgres OP-10 push 해결
2. **valkey ralph-loop 결과** — P1-11/12/13 + OP-2 + OP-10 ralph-loop 가 보강 (또는 본 thread 가 ralph-loop 중지 후 직접)
3. **postgres postgres-ADR/0019 ↔ 실 상태 (3 workflow) 일관성** — 별 ADR 또는 postgres-ADR/0019 갱신
4. **audit ❌ 0 시점** = **S8 v3.x-stable 선언 cycle** 진입
5. **5 repo vX.Y.Z release tag** (S8 의 Phase 4)

## 자동 갱신

본 문서는 audit script 출력의 *수동 정리*. 후속 cycle 에서 자동화 가능:
- `make audit-quality > docs/quality/audit-latest.md`
- 시계열 자동 append
- 진척 그래프 (mermaid)
