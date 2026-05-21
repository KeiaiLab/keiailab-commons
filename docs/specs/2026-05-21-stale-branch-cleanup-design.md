# S2 — operator-commons stale 브랜치 정리 + 레거시 lefthook 통합 (Design)

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Proposed |
| 작성자 | keiailab — superpowers cycle (post-supercycle drill-down) |
| 범위 | `operator-commons` 단일 저장소 (Go library SSOT) |
| Supercycle | [`2026-05-21-portfolio-cleanup-supercycle-design.md`](../superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md) 의 **Sub-Project S2** |
| 후속 | implementation plan (`docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md` 의 *archive 처리 + lefthook 통합* 부분에 적용) |
| Commit 작성자 | `TaeHwan Park <eightynine01@gmail.com>` + `Signed-off-by:` + `Co-Authored-By: Claude` |
| 우선순위 | High (portfolio supercycle G5 stale-branch=0 + Wave 2 잔여 부채 해소 게이트) |

## 1. 배경 (Background)

### 1.1 현재 상태 스냅샷 (2026-05-21 13:00 KST)

직접 측정 (`git ls-remote origin` + `git tag -l` + `gh pr list` + `ls -la`):

| 항목 | 측정값 | 비고 |
|---|---|---|
| `origin/main` | `bdb659c` (S2 plan 머지 직후) | clean fast-forward |
| `origin/archive/main-13-commits-merge-style-2026-05-21` | `910042a` | **본 spec 의 1차 정리 대상** |
| `origin/feat/i18n-ja-native-2026-05-21` | `be72bae` | in-flight (native reviewer 후속) — 본 spec OOS |
| `origin/feat/i18n-zh-native-2026-05-21` | `a64cf05` | in-flight (native reviewer 후속) — 본 spec OOS |
| 머지된 stale 브랜치 (origin) | 0개 | 사용자가 사전 정리 완료 (브리프 §컨텍스트) |
| open PR | 0건 | `gh pr list` = empty |
| open issue | 0건 | `gh issue list` = empty |
| 기존 tag (`archive/*`) | 2건 | `archive/fix-lefthook-config-merge-2026-05-21` + `archive/main-13-commits-merge-style-2026-05-21-tag` |
| `.lefthook.yml` (레거시) | 65 LOC | RFC-0002 §1 SSOT 정합 minimal 버전 (Go 단일 stack) |
| `lefthook.yml` (신규) | 407 LOC | ai-dev SSOT template + Go append (multi-stack 가드) |

### 1.2 두 lefthook 파일 공존 = Wave 2 잔여 부채

분석 (`diff <(head -20 .lefthook.yml) <(head -20 lefthook.yml)` + 양 파일 정독):

- **레거시 `.lefthook.yml`** (Go-only minimal):
  - pre-commit: gofmt + govet + golangci-lint
  - pre-push: go test + govulncheck + gitleaks + go-mod-tidy drift
  - commit-msg: Conventional Commits + DCO sign-off
- **신규 `lefthook.yml`** (ai-dev template + Go append):
  - `assert_lefthook_installed: true` + `no_tty: true` (CI 안전 옵션)
  - pre-commit: py / js / rust / helm-default-falsy-toggle / helm-cycle-26-hardening / adr-phantom-check / hetzner-mention-block / **gofmt 미포함** (Go-only stack 에서 누락)
  - pre-push: dirty-merged-worktree-check / gitlab-mcp-bypass-block / **go-vet + go-test + golangci-lint** (Go 가드 정상)
  - commit-msg: Conventional Commits + DCO sign-off + ...
  - post-merge: deps-log + RFC-0046 auto-branch-cleanup

**lefthook 동작 규약**: lefthook CLI 는 `lefthook.yml` 을 우선 로드하고 `.lefthook.yml` 은 *fallback override* 로 해석한다. 양 파일이 동시 존재하면 신규 `lefthook.yml` 이 활성화되고 레거시 `.lefthook.yml` 의 *Go pre-commit 가드 (gofmt 포함)* 는 무시되어 **부분적 회귀** 가능. 이미 신규 파일에 go-vet / go-test / golangci-lint 가 포함되었으나 *gofmt write-staged* (`gofmt -l -w {staged_files}`) 는 누락.

### 1.3 archive 브랜치의 *의도된 보존* 의미

- `archive/main-13-commits-merge-style-2026-05-21` = 2026-05-21 오전 main 의 merge-style commit 13건을 squash-flatten 하기 전 *원본 보존* 백업. 이미 동일 작업의 squash 결과가 main 에 있고, archive-only commit 들은 main 의 head 와 *file-content equivalent*.
- 그러나 사용자가 본 spec 의 결정 사항 1번 에서 명시: "annotated tag 후 브랜치 삭제". → 즉 *영구 archive 형태* 는 `tag` 로 변환하고, *브랜치 ref 형태* 는 폐기.

### 1.4 글로벌 거버넌스 정합 (CLAUDE.md)

| 조항 | 본 spec 의 적용 |
|---|---|
| §2 한국어 | 본 spec / 후속 ADR / commit message 본문 모두 한국어. tag annotation 도 한국어. |
| §2 테스트 없는 기능 = 존재 불가 | 본 spec 은 *문서* 만 작성. 후속 implementation 의 §6 success criteria bash 5건이 검증 게이트. |
| §2 GitHub Actions 영구 금지 (RFC-0002) | operator-commons `.github/workflows/` = 0건 유지 (현 상태 보존). 본 spec 은 lefthook 통합 = *로컬 게이트 강화* 방향성과 정합. |
| §5 git-flow 미사용 = 실패 | `spec/stale-branch-cleanup-2026-05-21` branch → squash merge → branch 삭제. atomic 1 spec = 1 PR. |
| §5 범위 외 수정 = 실패 | OOS 5건 명시 (§7). 추가 발견 시 별 sub-spec. |
| §8 atomic 정책 | Phase 1 (tag+삭제) / Phase 2 (lefthook 통합) / Phase 3 (검증) = 각 1 commit. cleanup 은 Phase 4 verify 100% PASS 후. |

## 2. 목표 (Goals) + 비목표 (Non-Goals)

### 2.1 Goals (사용자 시점 시나리오)

| ID | 목표 | 사용자 시점 검증 | 측정 명령 |
|---|---|---|---|
| G1 | archive 브랜치 → annotated tag 변환 + 브랜치 삭제 | "GitHub Branches 탭에 archive 브랜치 없음 + Tags 탭에 보존본 1건" | `git ls-remote --heads origin \| grep -c archive = 0` + `git ls-remote --tags origin \| grep -c archive-merge-style-v0.7.0 = 1` |
| G2 | 레거시 `.lefthook.yml` 삭제 + 신규 `lefthook.yml` 의 Go 가드 보강 (gofmt 누락분 add) | "개발자가 `git commit` 시 gofmt drift 자동 차단" | `test ! -f .lefthook.yml` + `grep -q 'gofmt' lefthook.yml` |
| G3 | ADR 작성 (`docs/kb/adr/0009-lefthook-config-consolidation.md`) | "사후 추적: 왜 통합했는가" | `test -f docs/kb/adr/0009-*.md` |
| G4 | 검증 (`lefthook run pre-push` 통과) | "통합 후 pre-push 무회귀" | `lefthook run pre-push` exit 0 |
| G5 | stale 브랜치 0 (main + in-flight i18n native 2건만) | "원격 브랜치 = main + feat/i18n-ja + feat/i18n-zh" | `git ls-remote --heads origin \| wc -l = 3` |

### 2.2 Non-Goals (본 S2 의 OOS)

- ❌ `feat/i18n-ja-native-2026-05-21` / `feat/i18n-zh-native-2026-05-21` 처리 — native reviewer 후속 wave (S4 마스터 spec 의 Phase 7 진입점)
- ❌ docs/i18n/glossary-{ko,ja,zh}.md 의 placeholder → 본문 보강 — S4 마스터 spec Phase 1 의 범위
- ❌ README.{ja,zh}.md 본문 확장 — S4 Phase 4 우선순위 1
- ❌ `docs/family.md` 의 v0.7.0 → v0.8.0 drift 갱신 — 별 PR (S2 plan 의 task 4) 으로 분리
- ❌ Wave 3 branding 잔여 — S3 sub-spec

## 3. 아키텍처 (Architecture)

### 3.1 Phase 모델 (4 단계 atomic)

```
Phase 0 사전확인 → Phase 1 archive tag+삭제 → Phase 2 lefthook 통합+ADR → Phase 3 검증
                                                                          └→ (PASS 시) cleanup
```

각 Phase = 1 commit + 1 git ref 변경. Phase 별 git push (또는 tag push) 후 즉시 검증.

### 3.2 git ref 변환 다이어그램

```
[ before ]
 origin/main                            ← bdb659c
 origin/archive/main-13-commits-...     ← 910042a (branch)
 origin/feat/i18n-ja-native-...         ← be72bae (OOS, 보존)
 origin/feat/i18n-zh-native-...         ← a64cf05 (OOS, 보존)
 tags: v0.1.0 .. v0.8.0
       archive/fix-lefthook-config-merge-2026-05-21
       archive/main-13-commits-merge-style-2026-05-21-tag  (이전 cycle 잔존 tag)

[ after ]
 origin/main                            ← <Phase 2 머지 commit>
 origin/feat/i18n-ja-native-...         ← be72bae (보존)
 origin/feat/i18n-zh-native-...         ← a64cf05 (보존)
 tags: ... + archive-merge-style-v0.7.0 (NEW, annotated)
       (기존 archive/main-13-commits-merge-style-2026-05-21-tag 는 *이미* 존재하므로
        실제 동작은 §4 Phase 1 결정 분기 참조)
```

### 3.3 lefthook 통합 분기

- 옵션 A (채택): 레거시 `.lefthook.yml` 의 Go-stack 가드 (gofmt write-staged) 를 신규 `lefthook.yml` 에 *merge in* 후 `.lefthook.yml` 삭제 → ADR 작성
- 옵션 B (기각): 양 파일 보존 — lefthook CLI 의 fallback 동작상 *실제 사용 anti-pattern*. ADR 부재면 §5 실패
- 옵션 C (기각): 신규 `lefthook.yml` 폐기 + 레거시 `.lefthook.yml` 유지 — ai-dev SSOT (multi-stack 가드 + RFC-0027 helm + RFC-0046 auto-cleanup + DCO sign-off + cycle 26 hardening) 손실 발생

## 4. Phase 상세

### 4.1 Phase 0 — 사전확인 (read-only)

**목적**: 본 spec 작성 시점의 원격 상태가 implementation 시점에도 동일한지 검증.

**Run**:
```bash
git fetch --prune origin
git ls-remote --heads origin | awk '{print $2}' | sort > /tmp/s2-heads-now.txt
diff /tmp/s2-heads-now.txt - <<EOF
refs/heads/archive/main-13-commits-merge-style-2026-05-21
refs/heads/feat/i18n-ja-native-2026-05-21
refs/heads/feat/i18n-zh-native-2026-05-21
refs/heads/main
EOF
gh pr list --state open --json number,title --jq 'length'  # 0
gh issue list --state open --json number,title --jq 'length'  # 0 또는 1 (#0 = 자동 닫힘)
test -f .lefthook.yml && test -f lefthook.yml  # 양 파일 공존 확인
```

**Expected**: diff exit 0 + open PR/issue 0 + 양 lefthook 파일 존재.

**Gate**: 1건이라도 불일치면 Phase 1 진입 차단. 새 상태로 본 spec 의 §1.1 표 갱신 후 재진입.

### 4.2 Phase 1 — archive 브랜치 → annotated tag + 브랜치 삭제

**목적**: §3.2 의 git ref 변환 (좌→우) 의 archive 브랜치 부분 적용.

**4.2.1 신규 tag 명명 결정**

- 기존: `archive/main-13-commits-merge-style-2026-05-21-tag` (이미 존재, slash prefix + suffix `-tag`)
- 신규: `archive-merge-style-v0.7.0` (사용자 결정 사항)
- 충돌 검증: `git ls-remote --tags origin | grep -c 'archive-merge-style-v0.7.0' = 0` (= 신규)

**4.2.2 명령**

```bash
# 1. 신규 annotated tag 생성 (브랜치 tip 보존)
git tag -a archive-merge-style-v0.7.0 \
  origin/archive/main-13-commits-merge-style-2026-05-21 \
  -m "archive: main v0.7.0 merge-style commit 13건 사전-squash 백업 (S2 cleanup 2026-05-21)

본 tag 는 origin/archive/main-13-commits-merge-style-2026-05-21 브랜치 삭제 전 보존본.
원본 13 commit 은 main 의 squash merge 결과와 file-content equivalent.
참조: docs/specs/2026-05-21-stale-branch-cleanup-design.md §4.2"

# 2. tag push
git push origin archive-merge-style-v0.7.0

# 3. 원격 브랜치 삭제
git push origin --delete archive/main-13-commits-merge-style-2026-05-21

# 4. 검증
git ls-remote --tags origin | grep archive-merge-style-v0.7.0  # 1줄
git ls-remote --heads origin | grep -c archive  # 0
```

**Expected**:
- tag push 출력: `[new tag] archive-merge-style-v0.7.0 -> archive-merge-style-v0.7.0`
- 브랜치 삭제 출력: `- [deleted] archive/main-13-commits-merge-style-2026-05-21`
- 검증 1줄 + 0개

**Commit**: 본 Phase 는 *git ref* 변경만 (working tree 변경 없음) → repo commit 없음. PR 의 description 에 명령 출력 인용.

**Risk**: tag push 권한 부족 시 (`gh auth status` 의 token scope) → 사전 `git push --dry-run origin archive-merge-style-v0.7.0` 로 사전 검증.

### 4.3 Phase 2 — 레거시 `.lefthook.yml` 통합 + 삭제 + ADR

**목적**: §3.3 옵션 A 적용. 신규 `lefthook.yml` 에 누락된 Go 가드 (gofmt write-staged) 보강 후 레거시 파일 삭제.

**4.3.1 gofmt 가드 보강**

신규 `lefthook.yml` 의 `pre-commit.commands` 블록에 추가 (line 11 직후, py-lint 보다 위):

```yaml
    gofmt:
      glob: "*.go"
      run: gofmt -l -w {staged_files}
      stage_fixed: true
```

**근거**: 레거시 `.lefthook.yml` 의 Go pre-commit 3종 (gofmt + govet + golangci-lint) 중 govet / golangci-lint 는 이미 신규 파일의 pre-push 에 존재. gofmt 만 누락 — *write-staged* 동작은 commit 자체를 막지 않고 자동 수정 후 commit 진입을 보장하므로 pre-commit 위치가 적합.

**4.3.2 레거시 파일 삭제 + ADR 작성**

```bash
git rm .lefthook.yml
```

ADR 신규 파일: `docs/kb/adr/0009-lefthook-config-consolidation.md`

ADR 본문 골격:
```markdown
# ADR-0009: lefthook 설정 통합 (.lefthook.yml → lefthook.yml)

## Status
Accepted (2026-05-21)

## Context
operator-commons 는 2026-05-20 시점에 ai-dev SSOT template 의 `lefthook.yml` (407 LOC) 을
도입하면서 기존 RFC-0002 §1 정합 minimal 버전 `.lefthook.yml` (65 LOC) 을 *중복 보존* 했다.
lefthook CLI 는 `lefthook.yml` 을 우선 로드하고 `.lefthook.yml` 은 fallback override 로
해석하므로, 양 파일 공존은 Go pre-commit 가드 (gofmt) 무시 가능성을 내포한다.

## Decision
1. 신규 `lefthook.yml` 에 누락된 gofmt write-staged 가드를 add.
2. 레거시 `.lefthook.yml` 을 `git rm` 으로 삭제.
3. 단일 SSOT 파일 = `lefthook.yml` 로 통일.

## Consequences
+ Go 가드 (gofmt) 가 pre-commit 시점에 강제 적용 (회귀 차단).
+ 설정 파일 1개 = drift 발생 가능성 0.
+ ai-dev SSOT 와 일관 — 향후 sync drift seal (RFC-0029 §6.5) 정합.
- 레거시 minimal 정의는 git 이력으로만 추적 (commit hash + 본 ADR cross-link).

## Alternatives
- (B) 양 파일 보존 — lefthook fallback 의 anti-pattern. ADR 부재 시 §5 실패.
- (C) 신규 파일 폐기 + 레거시 유지 — multi-stack 가드 + RFC-0027/0046 자산 손실.

## Cross-link
- spec: docs/specs/2026-05-21-stale-branch-cleanup-design.md §4.3
- 관련 RFC: RFC-0002 (GHA 금지 → 로컬 4 계층 일원화)
```

**4.3.3 Commit**

```bash
git add lefthook.yml docs/kb/adr/0009-lefthook-config-consolidation.md
git rm .lefthook.yml  # 이미 위에서 수행 시 add 만
git commit -s -m "$(cat <<'EOF'
chore(lefthook): 레거시 .lefthook.yml 삭제 + 신규 lefthook.yml 에 gofmt 가드 보강

- 신규 lefthook.yml 에 pre-commit gofmt write-staged 추가 (Go drift 차단)
- 레거시 .lefthook.yml (65 LOC) 삭제 — 신규 lefthook.yml (407 LOC) 로 단일화
- ADR-0009 사후 추적: lefthook config consolidation
- RFC-0002 §1 로컬 4계층 정합 + ai-dev SSOT (RFC-0029 §6.5) 정합

Refs: docs/specs/2026-05-21-stale-branch-cleanup-design.md §4.3

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

**Note**: `-s` 로 DCO sign-off trailer 자동 추가. Conventional Commits `chore(lefthook):` prefix 로 lefthook commit-msg 통과.

### 4.4 Phase 3 — 검증 + push + PR

**4.4.1 로컬 검증**

```bash
# lefthook 통합 후 무회귀 확인
lefthook run pre-commit
lefthook run pre-push  # go-vet + go-test + golangci-lint + gitleaks ...

# 원격 ref 정합
git ls-remote --heads origin  # main + feat/i18n-ja + feat/i18n-zh (3건)
git ls-remote --tags origin | tail -3  # archive-merge-style-v0.7.0 포함
```

**Expected**:
- `lefthook run pre-push` exit 0 (모든 명령 PASS)
- 원격 heads = 3건 (in-flight i18n 2건 보존)
- 신규 tag 존재

**4.4.2 push + PR**

```bash
git push -u origin spec/stale-branch-cleanup-2026-05-21
gh pr create \
  --title "spec(S2): operator-commons stale 브랜치 정리 — design doc" \
  --body "$(cat <<'EOF'
## 요약

S2 sub-project 의 design spec. 머지 후 implementation plan 에서 본 spec 의 §4 Phase 1-3 을 실행한다.

본 PR 자체는 *문서 1건* (`docs/specs/2026-05-21-stale-branch-cleanup-design.md`) 만 추가하며, Phase 1/2 의 git ref 변경 (archive tag+삭제) 및 코드 변경 (lefthook 통합) 은 별 PR 로 분리된다.

## 범위

- 신규 spec 파일 1개 (docs/specs/ 하위)
- 코드 변경 없음, lefthook 변경 없음, ADR 신규 없음 (본 PR 한정)

## 후속

- implementation PR: `chore(lefthook): consolidate + ADR-0009` + tag/branch 변환 명령 실행
- 관련: portfolio supercycle spec G5 (stale-branch=0) 정합

## 검증

- [x] 한국어 작성
- [x] 사용자 시점 시나리오 명세 (§2.1 G1~G5)
- [x] 범위 = 의도 (§7 OOS 명시)
- [x] context7 MCP 미사용 (외부 라이브러리 사용 없음)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

## 5. 리스크 + 완화 (Risks)

| ID | 리스크 | 영향 | 완화 |
|---|---|---|---|
| R1 | tag push 권한 부족 (gh token scope) | Phase 1 차단 | `git push --dry-run origin <tag>` 사전 검증 + `gh auth refresh -s write:packages` 등 scope 확인 |
| R2 | 신규-레거시 lefthook 차이로 회귀 (gofmt 외 가드 누락 미발견) | pre-commit drift | Phase 2 적용 후 `lefthook run pre-commit` + 의도적 *bad-format* Go 파일 stage 테스트 |
| R3 | archive 브랜치 의 commit 이 *완전히 file-content equivalent 가 아닐* 가능성 | 보존 손실 | `git cherry origin/main origin/archive/...` 로 `-` (equivalent) vs `+` (only-in-archive) 사전 정밀 식별. `+` commit 의 *file diff* 가 main 의 squash 결과와 동일함을 확인. tag 보존으로 *모든 경우* 안전. |
| R4 | 양 lefthook 동시 존재 기간 (Phase 2 commit 전) 의 hooks 동작 불확실성 | 짧은 윈도우 (~수분) drift | Phase 2 commit 을 빠르게 처리 + spec 머지 직후 implementation 진입 |
| R5 | in-flight feat/i18n-{ja,zh}-native 브랜치 의 작업이 S4 마스터 spec 의 Phase 7 와 충돌 | 작업 중복 | S2 OOS 명시 (§2.2). S4 spec 작성 시 본 브랜치 처리 정책 명시 (D5) |

## 6. 성공 조건 (Success Criteria)

implementation 완료 시 다음 5건 모두 PASS:

```bash
# SC1: archive 브랜치 부재
test "$(git ls-remote --heads origin | grep -c archive)" = "0"

# SC2: 신규 tag 존재
git ls-remote --tags origin | grep -q archive-merge-style-v0.7.0

# SC3: 레거시 lefthook 부재 + 신규 lefthook 의 gofmt 존재
test ! -f .lefthook.yml && grep -q "gofmt" lefthook.yml

# SC4: ADR 작성
test -f docs/kb/adr/0009-lefthook-config-consolidation.md

# SC5: 원격 heads = 3 (main + i18n in-flight 2)
test "$(git ls-remote --heads origin | wc -l | tr -d ' ')" = "3"
```

5건 모두 exit 0 시 implementation 완료. 1건이라도 실패 시 §6 검증 의무 미이행 = §5 실패.

## 7. 범위 외 (Out-of-Scope)

본 S2 spec 은 *operator-commons 단일 저장소의 stale 브랜치 + lefthook 통합* 만 다룬다. 다음 항목은 별 sub-spec 또는 후속 wave 의 범위:

- **S3 (브랜딩 통일)**: valkey + commons + forgewise 의 BRANDING.md / docs/family.md 정합 — 별 spec
- **S4 (다국어 4-lang 마스터)**: 본 작업 2 의 산출물. operator-commons 가 i18n SSOT 로서의 정책 + 5 repo 배포 + 자동 번역 파이프라인
- **S5 (operator-commons 공통화)**: 3 operators 의 중복 코드를 commons 패키지화 (예: probes / storageclass / events 외 추가 추출)
- **S6 (forgewise 정합화)**: 완료 (Wave 2/3/4 머지 완료, 본 spec 시점 기준 0 부채)
- **S7 (postgres + mongodb GHA 제거)**: 각 저장소의 `.github/workflows/` 삭제 — 별 sub-project
- **in-flight i18n 브랜치 처리** (`feat/i18n-{ja,zh}-native-2026-05-21`): native reviewer 후속 wave (S4 Phase 7 진입점)

## 8. 변경 이력

| 날짜 | 변경 | 상태 |
|---|---|---|
| 2026-05-21 | 초안 작성 (portfolio supercycle S2 drill-down) | Proposed |
