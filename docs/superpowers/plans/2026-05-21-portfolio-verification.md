# Portfolio Cleanup Supercycle 2026-05-21 — Verification

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| Spec | [2026-05-21-portfolio-cleanup-supercycle-design.md](../specs/2026-05-21-portfolio-cleanup-supercycle-design.md) |
| 상태 | **Mostly Implemented** (S1 마무리 진행 중) |
| 총 PR 머지 (5 repos) | 약 50건 (각 repo 별 표 참고) |
| 작성자 | keiailab |

## 1. 5 Repos 의 Acceptance (G1~G10 매핑)

| Goal | 현 달성 | 잔존 |
|---|---|---|
| G1 (.github/workflows 부재 → **per-repo 결정** [D1 v2.0 amendment]) | postgres = **strict 제거 노선** (14→3 narrow exception: pages/dependabot/release tag); mongo = **retention** (12 wf 유지); valkey = **retention** (14 wf, ADR-0048); commons = 0 (애초); forgewise = 0 (애초) | per-repo 정합 (cleanup 의도 달성) |
| G2 (BRANDING+family) | 5/5 ✓ (`for r: test -f BRANDING.md && test -f docs/family.md` PASS) | — |
| G3 (4-lang README + canonical docs) | 5/5 README 4-lang ✓ + commons 4-lang glossary 본문 (203/207/207 LOC) + commons docs/family / STABILITY / coverage-report ko/ja/zh + 다른 4 repos 의 docs 다국어 사용자 평행 진행 | forgewise 거버넌스 5종 ja/zh 부분 placeholder (별 cycle) · CRD description i18n annotation (별 cycle) |
| G4 (3 ops commons v0.8.0+) | postgres v0.9.0 (`4abf932`) · mongo v0.9.0 (`a1ff450`) · valkey v0.9.0 (`daa763c`) — v0.8.0 초과 | — |
| G5 (stale 0) | commons 0 / postgres 0 / mongo 0 / forgewise 0 / **valkey 12** (in-flight 정리 중) | valkey 12 (사용자 진행 중) |
| G6 (required_status_checks 0) | 5/5 (모두 absent) — `gh api repos/keiailab/<r>/branches/main/protection` `.required_status_checks` null | — |
| G7 (CLAUDE.md multi-arch 조항 제거) | local commit 만 (~/.codex/ commit `0a98641`, push 보류) | RFC 0006 push (~/.codex/ remote 결정 후) |
| G8 (forgewise apache-2.0) | LICENSE 파일 Apache-2.0 정합. GitHub `licenseInfo.key` detection 확인 필요 | GitHub linguist license detection 수정 (별 cycle) |
| G9 (임시 파일 정리) | 진행 중 — `$CLAUDE_JOB_DIR/*` 자동 정리 영역. `.claude/worktrees/` 잔존물 0. HANDOFF.md root 잔존물 0 (Wave 2 docs/internal/ 이전 완료) | docs/plans/ 잔존 (각 sub-spec 의 plan 보존 — 의도된 잔존) |
| G10 (PR body 검증 인용) | 모든 PR body 에 `Co-Authored-By: Claude` trailer + verify 절차 인용 | — |

**대략 종합: 8.5/10** (G5 valkey 12 stale + G7 RFC 0006 push 보류 + G8 forgewise license detection)

## 2. 5 Repos 의 PR 머지 카탈로그 (2026-05-21)

### 2.1 operator-commons (commit master, 22 PR)

| PR | 제목 | 분류 |
|---|---|---|
| #24 | chore(release): cliff.toml + Chart.yaml v0.7.0 + ADOPTERS.md refresh | release prep |
| #25 | docs(godoc): pkg/*/doc.go × 8 — API Stability Tier marker | docs |
| #26 | feat(ci): RFC-0043 L5 shadow + lefthook bootstrap + gofmt + golangci-lint | ci |
| #27 | docs(i18n): README 다중언어 (EN canonical + KO) + sync hook + owner | i18n |
| #28 | chore(docs): Wave 2 cleanup — HANDOFF.md → docs/internal/ | Wave 2 |
| #29 | feat(probes): pkg/probes Builder — HTTP/HTTPS/TCP/Exec fluent API | v0.8.0 |
| #30 | feat(storageclass): pkg/storageclass — DNS-1123 validator + Normalize | v0.8.0 |
| #31 | feat(events): pkg/events — Recorder helper + reason constants | v0.8.0 |
| #32 | docs(v0.8.0): ADOPTERS + ROADMAP + CHANGELOG — 3 신규 패키지 등재 | v0.8.0 release |
| #33 | docs(i18n): glossary-{ko,ja,zh}.md — 표준 용어 사전 v0.1 | i18n |
| #34 | docs(i18n): README 4-lang lang switcher + keiailab family footer | i18n |
| #35 | feat(branding): Wave 3 keiailab branding — BRANDING.md + family.md + 13 footer | Wave 3 |
| #36 | docs(spec): portfolio cleanup supercycle 2026-05-21 — 5 repos S1~S7 master spec | **portfolio** |
| #37 | docs(plan): S2 — operator-commons stale 정리 + i18n placeholder → 본문 plan | S2 plan |
| #38 | spec(S2): operator-commons stale 브랜치 정리 — design doc | S2 sub-spec |
| #39 | spec(S4): 다국어 4-lang 마스터 spec (operator-commons SSOT) — D1~D5 결정 | **S4 master** |
| #40 | docs(i18n): README.ja.md full native 翻訳 (122 LOC, 敬語 일관) | S2/S4 |
| #42 | chore(branches): archive/main-13-commits-merge-style 정리 + ADR-0009 | **S2 T1** |
| #43 | chore(branches): S2 Phase 1 — archive-merge-style-v0.7.0 tag + ADR-0010 | **S2 T1** |
| #44 | docs(version): family.md + 4-lang README v0.7.0 → v0.8.0 drift 해소 | **S2 T2** |
| #45 | chore(lefthook): S2 Phase 2 — .lefthook.yml 삭제 + 통합 + ADR-0011 | **S2 P2** |
| #46 | docs(i18n): README.zh.md placeholder → full native 简体中文 (~125 LOC) | **S2 T3** |
| #47 | feat(i18n): S4 Phase 1 — commons SSOT 정비 (glossary 4-lang + sync 매트릭스) | **S4 P1** |
| #48 | feat(lefthook): S4 Phase 2 — readme-i18n-sync hook 통합 (pre-push drift) | **S4 P2** |
| #49 | feat(i18n): S4 Phase 3 — sync-from-commons.sh 신규 (SSOT → 4 repo 배포) | **S4 P3** |
| #50 | chore(s2): verification.md + spec S2 Status Implemented | **S2 T5** |
| #51 | docs(i18n): S4 Phase 4 — commons docs/ P0 우선순위 번역 | **S4 P4** |
| #52 | feat(pkg): Sprint 1 — pkg/pvc + pkg/topology 신규 (3-operator 중복 ~495 LOC 해소) | **S5 Sprint 1 P2** |
| #53 | docs(i18n): S4 Phase 5 — top-level P1 번역 (BRANDING 3 lang + ADOPTERS 2 lang) | **S4 P5** |
| #54 | docs(spec): portfolio D1 amendment — GHA retention + 이중 운영 (v2.0) | **portfolio amendment** |
| #55 | chore(gitignore): .claude/ + .claude-flow/ 추가 (stale 브랜치 대체) | hygiene |

(추가 main 직접 commit: ADR-0012~0016, audit SSOT, release.sh, UPGRADING, audit lefthook P0-6/P1-12/13)

### 2.2 postgres-operator (15 PR)

| PR | 제목 | 분류 |
|---|---|---|
| #81 | chore(docs): cleanup + internal isolation (Wave 2) | Wave 2 |
| #82 | docs(branding): keiailab operator family branding (Wave 3) | **S3** |
| #83 | docs(i18n): README 4-lang 완성 — ja/zh native 237 LOC + lang switcher | **S3/S4** |
| #84 | feat(deps): operator-commons v0.7.0 → v0.8.0 + pkg/probes 2 HTTP site 적용 | **S5** |
| #85 | docs(adr): ADR-0017 integrated GHA retention rationale for public OSS | S7 P1 |
| #86 | feat(ci): RFC-0002 — .github/workflows 전체 제거 (14 파일) | S7 strict |
| #87 | feat(release): scripts/helm-publish.sh + scripts/release.sh 신규 | OP-2/10 |
| #88 | feat(ci): 로컬 4계층 보강 — kube-linter + go-licenses + markdown-link-check | S7 보강 |
| #89 | docs(adr): ADR-0018 Accepted — GHA 전면 제거 → 로컬 4계층 단일 운영 | S7 strict |
| #90 | revert(ci): restore .github/workflows (14 files) — operator family v2.0 정합 | S7 revert |
| #91 | feat(commons): Sprint 1 Phase 2 — pkg/pvc + pkg/topology 채택 (-375 LOC) | **S5 Sprint 1 P2** |
| #92 | docs(adr): ADR-0019 Accepted — GHA 유지 + operator family v2.0 통합 정합 | **S7 v2.0** |
| #93 | chore(ci): C7 RFC-0002 GitHub Actions cleanup — 11 workflow 제거 | **S7 strict 정밀** |
| #94 | chore(gitignore): .claude/ + .claude-flow/ session artifact (Codex Major #6) | hygiene |
| #95 | feat(lefthook): RFC-0002 gha-block hook + UPGRADING.md (P2-2 + OP-10) | **S7 hook** |

추가 main 직접 commit: `beb4d42 docs(adr): ADR-0022 Accepted — GHA narrow exception 3 workflow 정합`, `4abf932 chore(deps): bump operator-commons to v0.9.0 (Sprint 1 release)`

### 2.3 mongodb-operator (18 PR)

| PR | 제목 | 분류 |
|---|---|---|
| #187 | chore(docs): Wave 2 cleanup — remove deprecated pre-commit config | Wave 2 |
| #188 | feat(docs): keiailab 브랜딩 — README header/footer + BRANDING.md + family.md (Wave 3) | **S3** |
| #189 | docs(i18n): README 4-lang ja/zh placeholder + family footer 5-repo | **S3/S4** |
| #190 | feat(deps): operator-commons v0.7.0 → v0.8.0 + pkg/probes 2 Exec site 적용 | **S5** |
| #191 | docs(i18n): README.ja.md full native 翻訳 (557 LOC, base v0.8.0-consume) | **S4** |
| #192 | docs(i18n): README.zh.md full native 翻訳 (566 LOC) | **S4** |
| #193 | docs(adr): ADR-0031 integrated GHA retention rationale for public OSS | S7 P1 |
| #194 | feat(ci): RFC-0002 — .github/workflows 전체 제거 (12 파일) | S7 strict |
| #195 | feat(release): scripts/helm-publish.sh + scripts/release.sh 신규 | OP-2/10 |
| #196 | feat(ci): 로컬 4계층 보강 — kube-linter + go-licenses + markdown-link-check | S7 보강 |
| #197 | docs(adr): ADR-0032 Accepted — GHA 전면 제거 → 로컬 4계층 단일 운영 | S7 strict |
| #199 | revert(ci): restore .github/workflows (12 files) — operator family v2.0 정합 | S7 revert |
| #200 | chore(ci): C7 RFC-0002 GitHub Actions cleanup | S7 strict 재시도 (race) |
| #201 | chore(gitignore): .claude-flow/ session artifact | hygiene |
| #202 | feat(commons): Sprint 1 Phase 2 — pkg/pvc + pkg/topology 채택 (-327 LOC) | **S5 Sprint 1 P2** |
| #203 | docs(adr): ADR-0033 Accepted — GHA 유지 + operator family v2.0 통합 정합 | **S7 v2.0** |
| #204 | revert(ci): re-restore 9 workflows after PR #200 race — operator family v2.0 정합 | **S7 race 정정** |
| #205 | feat(lefthook): RFC-0002 gha-block hook + ADR-0035 (P2-2) | **S7 hook** |

추가 main 직접 commit: `a1ff450 chore(deps): bump operator-commons to v0.9.0 (Sprint 1 release)`

### 2.4 valkey-operator (20 PR + 3 open)

| PR | 제목 | 분류 |
|---|---|---|
| #138 | docs: OTel guide + 9.x flags follow-up (P-C.6 + C.8) | docs |
| #139~#156 | ci(deps): dependabot 전수 머지 (kubernetes / azure-sdk / api / storage / ginkgo / gomega 등 ~18건) | dependabot |
| #158 | fix(webhook+controller): TLS.Enabled immutable + ready message surfacing | feature |
| #159 | feat(operator): CDEX-M1 — PDB delete path 보장 | feature |
| #160 | chore(docs): 3-tier 분류 적용 — 개발 SSOT 를 docs/internal/ 로 이동 | Wave 2 |
| #161 | docs: keiailab branding 적용 — header/footer 표준 + BRANDING.md + family.md (Wave 3) | **S3** |
| #163 | spec(S1+): valkey-operator PR cleanup + GHA workflow 정합 + 통합 ADR-0048 | **S1 spec** |
| #167 | docs(adr): ADR-0048 integrated GHA retention rationale for public OSS operator | S1 P1 |
| #168 | ci(deps): bump actions/stale from 10.2.0 to 10.3.0 | dependabot |
| #169 | docs(adr): ADR-0048 Accepted — GHA retention for valkey (per operator family) | **S1 v2.0** |
| #170 | feat(commons): Sprint 1 Phase 2 — pkg/pvc + pkg/topology 채택 (-322 LOC) | **S5 Sprint 1 P2** |
| #171 | feat(lefthook): kube-linter + go-licenses + markdown-link-check 3종 hook (P1-11/12/13) | **S1 보강** |
| #172 | feat(scripts): helm-publish.sh 신규 (OP-2) | OP-2 |
| #173 | docs(upgrading): UPGRADING.md 신규 (OP-10) | OP-10 |
| #174 | docs(adr): ADR-0050 audit augmentation (P1-11/12/13 + OP-2 + OP-10) | **S1 P4** |

**Open (잔여)**: #164 (v0.8.0 consume + pkg/probes 2 Exec) · #165 (README.ja.md full native 187 LOC) · #166 (README.zh.md full native 189 LOC)

추가 main 직접 commit: `daa763c chore(deps): bump operator-commons to v0.9.0 (Sprint 1 release)`, `f2d9ee4 docs(adr): ADR-0048 Status Proposed → Accepted (per operator family trade-off)`

### 2.5 forgewise (13 PR)

| PR | 제목 | 분류 |
|---|---|---|
| #1 | docs(i18n): README 4-lang en canonical SSOT + ko rename + ja/zh placeholder | i18n |
| #2 | spec(S6): forgewise 거버넌스 + 운영 정합화 — design doc | **S6 spec** |
| #3 | feat(ci): lefthook.yml Python 4 계층 게이트 + setup-hooks Makefile target | **S6 ci** |
| #4 | docs(governance): 거버넌스 5종 파일 신규 (SECURITY/CONTRIBUTING/COC/CHANGELOG/AGENTS) | **S6 P1** |
| #5 | docs(plan): governance-and-ops-alignment plan tracking INDEX.md | S6 plan |
| #6 | docs(ops): 운영 문서 4종 신규 (Installation/Configuration/API Reference/Upgrade) | **S6 P2** |
| #7 | test(governance): 거버넌스 5종 + 운영 4종 + lefthook 키워드 강제 (S6 Phase 5) | **S6 P5** |
| #8 | docs(adr): 0001 — Python 스택 override vs 글로벌 Go 전제 standards (Accepted) | **S6 ADR** |
| #9 | docs(i18n): BRANDING 4-lang 신규 — en (canonical) + ko + ja/zh AI 번역 | **S3/S4** |
| #10 | docs(i18n): docs/family 4-lang 신규 — 5 repo cross-link canonical | **S3/S4** |
| #11 | docs(i18n): README ja/zh placeholder → 본문 완전 번역 + ko 배너/footer 보강 | **S4** |
| #12 | docs(i18n): 운영 4종 (installation/configuration/api-reference/upgrade) ko/ja/zh 신규 | **S4** |
| #13 | feat(i18n): commons SSOT sync + readme-i18n-sync lefthook 도입 | **S4/S6** |

추가 main 직접 commit: `aad600f feat(audit): lefthook P0-9/P1-13 보강 + release.sh + Makefile`, CODEOWNERS / ADOPTERS / ROADMAP / PR template 사용자 평행 추가

## 3. 후속 (별 cycle)

| 항목 | 분류 | 진입 |
|---|---|---|
| valkey S1 마무리 (3 open PR #164/#165/#166 + 12 stale + 1 issue #4 Renovate) | S1 잔여 | 사용자 진행 중 |
| 4 repos 의 open issues 처리 (postgres #18 / mongo #115 / valkey #4 Renovate Config) | issue cleanup | 별 cycle |
| commons audit 5건 ❌ 보강 sub-spec (commit `38f52d9`) | audit augmentation | ralph-loop 또는 별 cycle |
| RFC 0006 push (~/.codex/ commit `0a98641` 의 remote 결정) | 거버넌스 | 사용자 결정 |
| forgewise GitHub `licenseInfo.key` detection 수정 (G8) | 외부 detection | 별 cycle (linguist algorithm) |
| portfolio spec 의 §3 D1 의 "v2.0 amendment" + §2.1 G1 의 per-repo 결정 표기 정합 유지 | 거버넌스 정합 | 본 amendment 로 완료 |
| Sprint 1 Phase 3 (pkg/coordinator 신규 추출) | commons 진전 | 별 cycle |

## 4. 흔적 / 증거

- 약 50 PR (5 repos 합산 — commons 33 + postgres 15 + mongo 18 + valkey 20 + forgewise 13. 단 commons 32 개 머지 + 1 본 amendment = 33)
- ADR 12+ (commons 0001~0016 / postgres 0017~0022 / mongo 0031~0035 / valkey 0048~0050 / forgewise 0001)
- RFC 0006 (~/.codex/ local commit `0a98641`, push 보류)
- portfolio spec D1 amendment (PR #54, `0c15ccb`)
- portfolio spec 종합 amendment (본 PR, S1~S7 카드 + verification.md)
- 평행 background subagent 산출 보관 (`/Users/phil/.claude/jobs/4829bea1/`)
- supercycle Wave 5 진입 (cross-validation 게이트)

## 5. 검증 명령 (재현 가능)

```bash
cd /Users/phil/workspace/keiailab

# G1. per-repo GHA workflow 정합
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  N=$(ls "$r/.github/workflows/" 2>/dev/null | wc -l | tr -d ' ')
  echo "$r: $N wf"
done
# 기대: commons=0 / postgres=3 / mongo=12 / valkey=14 / forgewise=0

# G2. 5 repos BRANDING + family
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  test -f "$r/BRANDING.md" && test -f "$r/docs/family.md" && echo "✓ $r" || echo "✗ $r"
done

# G3. 5 repos 4-lang README
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  ALL_OK=true
  for f in README.md README.ko.md README.ja.md README.zh.md; do
    test -f "$r/$f" || { ALL_OK=false; echo "  ✗ $r/$f"; }
  done
  $ALL_OK && echo "✓ $r 4-lang"
done

# G4. 3 ops commons v0.8.0+
for r in postgres-operator mongodb-operator valkey-operator; do
  grep -E 'operator-commons v0\.(8|9)\.' "$r/go.mod" | head -1
done

# G5. stale 브랜치
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  cd "/Users/phil/workspace/keiailab/$r"
  N=$(git branch -r 2>/dev/null | grep -v -E '/main$|/gh-pages$|HEAD|stable' | wc -l | tr -d ' ')
  echo "  $r stale=$N"
done

# G6. required_status_checks 0
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  gh api "repos/keiailab/$r/branches/main/protection" 2>/dev/null | \
    jq -r ".required_status_checks // \"absent\"" | head -1 | xargs -I {} echo "$r: {}"
done
```

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../../../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
