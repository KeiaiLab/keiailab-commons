# S8 Execution Plan — v3.x-stable 선언 실행 명세

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Ready (audit ❌ 0 시점 즉시 진입) |
| 입력 | `docs/specs/2026-05-21-v3x-stable-declaration-design.md` (S8 spec) |
| 의존 | audit ❌ 0 (현 6 = valkey 만 — S-valkey subagent 완료 + valkey P2-2 추가 후 도달) |

## 1. 진입 조건 (자동 게이트)

```bash
# 본 명령이 exit 0 시 S8 진입 가능
test $(bash /Users/phil/Workspace/keiailab/operator-commons/scripts/audit-production-grade.sh /Users/phil/Workspace/keiailab 2>&1 | grep -c "❌") -eq 0
```

## 2. 실 명령 시퀀스 (Phase 별)

### Phase 0 — 사전 확인 (~ 5 min)

```bash
# 1. audit ❌ 0 확인 (재실행)
bash /Users/phil/Workspace/keiailab/operator-commons/scripts/audit-production-grade.sh \
  /Users/phil/Workspace/keiailab > /tmp/audit-pre-s8.md
grep -c "❌" /tmp/audit-pre-s8.md  # 결과 = 0

# 2. 5 repo 의 main 상태 확인 (clean + ahead 0)
for repo in postgres-operator mongodb-operator valkey-operator operator-commons forgewise; do
  cd /Users/phil/Workspace/keiailab/$repo
  git fetch origin main
  git status --short
  git rev-list --left-right --count HEAD...origin/main
done

# 3. 모든 cycle 의 task list completed 확인
# (본 thread 의 task list — TaskList tool)
```

### Phase 1 — RFC-0005 작성 (글로벌 standards repo)

```bash
# ~/.codex/rfcs/0005-v3x-stable-declaration.md 신규 작성
# 본문 ~/.codex/CLAUDE.md §7 의 "v3.x-stable 선언 조건" 충족 측정 결과 인용
# Status: Accepted (audit ❌ 0 + 5 repo 정합 + ADR 다수)

# 사용자 명시 결정 필요 — ~/.codex/ 직접 변경은 RFC 절차
# 또는 PR 양식으로 사용자 검토 후 머지
```

본 Phase 는 사용자 명시 결정 필요 영역 — `~/.codex/` 는 글로벌 standards repository.

### Phase 2 — 각 repo ADR 신설

각 repo `docs/kb/adr/` 에 *v3.x-stable baseline* ADR 작성:
- postgres-operator: `0023-v3x-stable-baseline.md` (ADR-0022 이후 다음 번호)
- mongodb-operator: `0036-v3x-stable-baseline.md`
- valkey-operator: `0050-v3x-stable-baseline.md` (ADR-0049 이후)
- operator-commons: `0017-v3x-stable-baseline.md` (ADR-0016 Sprint 1 이후)
- forgewise: `0002-v3x-stable-baseline.md` (ADR-0001 이후)

각 ADR 본문:
- Status: Accepted
- Date: (선언 일자)
- Context: CLAUDE.md §7 v3.x-stable 조건 충족 + audit ❌ 0
- Decision: 본 repo v1.0.0 (또는 적절 semver) tag
- Consequences: ✅ 외부 사용자 인정 가능 운영 등급, ⚠️ 후속 v3.1+ 진화

### Phase 3 — 5 repo release tag

```bash
# postgres-operator
cd /Users/phil/Workspace/keiailab/postgres-operator
# CHANGELOG.md 갱신 (v0.3.0-alpha.18 → v1.0.0?)
# 또는 alpha 단계 인정: v0.3.0
make release VERSION=v0.3.0  # 또는 v1.0.0

# mongodb-operator (현 v1.4.x)
cd /Users/phil/Workspace/keiailab/mongodb-operator
make release VERSION=v1.5.0  # minor bump (v3.x-stable 표시)

# valkey-operator (현 v0.x)
cd /Users/phil/Workspace/keiailab/valkey-operator
make release VERSION=v1.0.0  # 또는 적절 semver

# operator-commons (v0.8.0)
cd /Users/phil/Workspace/keiailab/operator-commons
make release VERSION=v0.9.0  # Sprint 1 결과 + audit SSOT 통합

# forgewise (현 unreleased)
cd /Users/phil/Workspace/keiailab/forgewise
make release VERSION=v0.1.0  # 첫 stable release
```

각 release.sh 가 자동: tag + push + gh release create + chart .tgz / wheel 첨부.

### Phase 4 — 5 repo README + family.md 갱신

```bash
# 각 repo README 상단에 v3.x-stable 배지 추가:
# [![keiailab v3.x-stable](https://img.shields.io/badge/keiailab-v3.x--stable-success)](https://github.com/keiailab/operator-commons/blob/main/docs/quality/audit-history.md)

# commons/docs/family.md: forgewise 추가 + 5 sister 인정 + 각 repo v1.0 tag 명시
# commons/docs/family.{ko,ja,zh}.md: 4-lang sync
```

### Phase 5 — CLAUDE.md §7 갱신 (글로벌)

```bash
# ~/.codex/CLAUDE.md 의 §7 의 "P0+P1+P2 자동화 모두 충족 시 *v3.x-stable* 선언"
# → "**2026-05-21 v3.0.1-stable 선언** (RFC-0005). audit: commons/docs/quality/audit-history.md."

# 사용자 명시 결정 영역 (글로벌 standards 변경)
```

### Phase 6 — release announce

```bash
# (선택) 외부 announce:
# - GitHub Release body (각 5 repo 의 gh release create 시 자동)
# - Twitter / blog / community channel
# - Adopters 에 notify

# (필수) 본 thread 의 HANDOFF.md 갱신:
# - Status: v3.x-stable Declared
# - Next session: v3.1+ planning
```

## 3. 검증 (선언 후)

```bash
# 1. audit ❌ 0 유지
bash commons/scripts/audit-production-grade.sh /path/to/parent | grep -c "❌"  # = 0

# 2. 5 repo 모두 v1.X.Y tag
for repo in postgres-operator mongodb-operator valkey-operator operator-commons forgewise; do
  cd ~/Workspace/keiailab/$repo
  git tag -l "v*" | tail -1 | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$' || echo "❌ $repo 부재"
done

# 3. RFC-0005 Accepted
grep -q "Status.*Accepted" ~/.codex/rfcs/0005-v3x-stable-declaration.md

# 4. CLAUDE.md §7 갱신
grep -q "v3.0.1-stable 선언" ~/.codex/CLAUDE.md

# 5. 5 repo ADR baseline
for repo in postgres-operator mongodb-operator valkey-operator operator-commons forgewise; do
  ls ~/Workspace/keiailab/$repo/docs/kb/adr/ | grep -q "v3x-stable-baseline" || echo "❌ $repo"
done
```

## 4. 자동화 정합 (post-cycle)

audit script 의 자동 트리거:
- `make audit-quality` 가 v3.x-stable 선언 후에도 *유지* — 회귀 차단
- 후속 cycle (v3.1+) 에서 P3/P4/P5 신설 시 audit script 갱신

## 5. v3.1+ 준비 (post-S8)

S8 완료 후 즉시 시작 가능한 v3.1 후보:
- P3 성능 게이트 (benchmark + budget)
- P4 DR 게이트 (backup + restore + chaos)
- P5 커뮤니티 KPI (이슈 응답 SLA + adopter 성장)
- 5 repo audit 의 자동 측정 cron (월 1회)
- audit-history.md 의 자동 갱신

## 6. 본 plan 의 진입 시점

본 plan 의 *실 실행* 은 다음 모두 충족 시:

```bash
# Gate 1: audit ❌ 0
[[ $(bash commons/scripts/audit-production-grade.sh /path | grep -c '❌') -eq 0 ]]

# Gate 2: 5 repo open PR ≤ 3 (각 repo)
for r in ...; do [[ $(gh pr list --repo keiailab/$r --state open --json number | jq length) -le 3 ]]; done

# Gate 3: 5 repo stale branch 정합 (각 ≤ 2)
for r in ...; do [[ $(cd $r && git branch -r | grep -vE 'main|HEAD|gh-pages' | wc -l) -le 2 ]]; done
```

현재 (2026-05-21 15:00):
- Gate 1: ❌ (6 잔여, valkey 만)
- Gate 2: postgres 1, mongodb 0, valkey 3, commons 1, forgewise 0 — 통과
- Gate 3: valkey 8 (ralph-loop 진행 중), 다른 ≤ 3 — valkey 만 미통과

→ valkey audit 6 ❌ 해소 시 = S8 진입 가능.
