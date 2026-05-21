# Portfolio-Wide Cleanup Supercycle — 5 Repository Integration Design

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | **Mostly Implemented** (2026-05-21 — S2/S4/S5/S6 Implemented · S3/S7 거의 완료 · S1 마무리 진행 중) |
| 작성자 | keiailab — superpowers brainstorming session |
| 범위 | 5 저장소 통합 정리 (postgres-operator / mongodb-operator / valkey-operator / operator-commons / forgewise) |
| Supercycle | `cleanup supercycle 2026-05-21` (Wave 2 = 3-tier 분류 완료 · Wave 3 = branding 완료 · Wave 4 = i18n 완료 · Wave 5 = cross-validation 본 spec) |
| 후속 | S1~S7 sub-spec / RFC 0005 / 각 sub-project 의 `writing-plans` 산출 |
| 우선순위 | (사용자 결정) operator-commons → 3 operators → forgewise |
| Commit 작성자 | `TaeHwan Park <eightynine01@gmail.com>` + `Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>` |
| Verification | [docs/superpowers/plans/2026-05-21-portfolio-verification.md](../plans/2026-05-21-portfolio-verification.md) |

## 1. 배경 (Background)

### 1.1 현재 상태 스냅샷 (2026-05-21 12:30 KST)

| Repo | 언어 | 디스크 | 원격 브랜치 | 열린 PR | 열린 이슈 | 최근 태그 | GHA workflows | BRANDING.md | docs/family.md | i18n (ko/ja/zh) |
|---|---|---|---|---|---|---|---|---|---|---|
| `operator-commons` | Go | 1.1 MB | 2 (main + archive) | 0 | 0 | **v0.8.0** | 0 | ✓ | ✓ (v0.7.0 표기 — drift) | ko ✓ / ja placeholder / zh placeholder |
| `postgres-operator` | Go | 30 MB | 4 (main + gh-pages + stable + ?) | 0 | 1 | v0.3.0-alpha.16 | 10 | ✓ (166 LOC) | ✓ (93 LOC) | ko ? / ja ✓ / zh ✓ |
| `mongodb-operator` | Go | 15 MB | 3 (main + gh-pages + feat/v0.8.0-consume) | **1** (#190 v0.8.0 consume) | 1 | v1.5.0 | 10 | ✓ (173 LOC) | ✓ (93 LOC) | ko ? / ja ✓ / zh ✓ |
| `valkey-operator` | Go | 16 MB | **25** (main + gh-pages + 16 dependabot + 5 feat + 2 fix) | **21** | 1 | v1.0.13 | **14** (1,068 LOC) | ✗ | ✗ | ko ? / ja ? / zh ? |
| `forgewise` | Python | 544 KB | 1 (main) | 0 | 0 | — | 0 | ✗ | ✗ | ko ✓ / ja ✓ / zh ✓ (README) |

측정: `gh api orgs/keiailab` · `gh pr list` · `git branch -r` · `ls .github/workflows/` (2026-05-21 12:30 KST).

### 1.2 이미 완료된 작업 (오늘 2026-05-21 오전)

- `operator-commons` **v0.8.0 release cut** (tag `v0.8.0`, 2026-05-21 12:22 KST):
  - 3 신규 패키지: `pkg/probes` (Experimental) · `pkg/storageclass` (Stable) · `pkg/events` (Beta)
  - PR #29~#35 머지 (probes / storageclass / events / governance / glossary / 4-lang README / Wave 3 branding)
- `postgres-operator` **branding Wave 3 머지** (PR #82, BRANDING.md + docs/family.md + README header/footer)
- `mongodb-operator` **branding Wave 3 머지** (PR #188, 동일 패턴)
- `mongodb-operator` **i18n Wave 4 머지** (PR #189, README.ja.md + README.zh.md placeholder + 5-link family footer)
- `postgres-operator` **i18n Wave 4 머지** (PR #83, README.ja.md + README.zh.md)
- `forgewise` 머지 (origin/main: `feat/i18n-en-canonical-2026-05-21` → main 머지 후 브랜치 삭제, README.ko/ja/zh 추가)

### 1.3 진행 중 (in-flight, untrracked)

- `valkey-operator/.claude/worktrees/spec+pr-cleanup-and-gha-removal-2026-05-21/`:
  - commit `22e13b0 docs(spec): S1+ design — valkey-operator PR cleanup + GHA 제거`
  - 본 spec 의 **S1 sub-project** 가 *이미 사용자 worktree 에서 시작됨*
- `mongodb-operator` PR #190 `feat(deps): operator-commons v0.7.0 → v0.8.0 + pkg/probes 2 Exec site 적용` (OPEN, BLOCKED 가능성)

### 1.4 글로벌 거버넌스 정합 (CLAUDE.md)

| 조항 | 본 spec 의 적용 결정 |
|---|---|
| `§2 한국어 작성` | 모든 spec/ADR/RFC/CHANGELOG 본문 한국어. 코드 식별자 영어. |
| `§2 테스트 없는 기능 = 존재 불가 (MUST)` | 각 sub-project 완료 게이트에 e2e + 통과 로그 인용 의무 |
| `§2 context7 MCP 사용` | 외부 lib 사용 결정 시 (e.g. dependabot 대체 도구 선정) |
| `§2 GitHub Actions 영구 금지 (RFC 0002)` | **엄격 적용** (사용자 결정) — 5 repos 모두 `.github/workflows/` 제거, dependabot 16건 전수 머지 후 제거 |
| `§2 멀티아키텍처 빌드 금지` | **변경됨** (사용자 결정) — §2 multi-arch 조항 제거. PR #157 머지 + valkey lefthook 의 `platforms-amd64-guard` 제거. **본 spec 이 §7 부트스트랩 예외에 따른 사후 RFC 0005 의 *명시 지시* 근거**. |
| `§5 git-flow 미사용 = 실패` | task 별 `feat/<task>-2026-05-21` branch → atomic PR → squash merge → branch 즉시 삭제 |
| `§5 범위 외 수정 = 실패` | 본 spec 은 OOS 7건 명시. 추가 발견 시 별 sub-project 분리. |

## 2. 목표 (Goals) + 비목표 (Non-Goals)

### 2.1 Goals (사용자 시점 시나리오)

| ID | 목표 | 사용자 시점 검증 | 측정 |
|---|---|---|---|
| G1 | (v2.0 amendment) 5 저장소 GHA workflow 가 *per-repo 정책* 정합 — valkey 는 retention + 로컬 4계층 이중 운영, 다른 4 는 *retention 또는 제거* 사용자 결정 (PR/spec 별) | "각 repo 의 .github/workflows/ 상태 = 의도 일치" | `for r: test 'expected' (retention/removed)`; — per-repo 확인 |
| G2 | 5 저장소 모두 BRANDING.md + docs/family.md 보유 | "README 클릭 → BRANDING 일관성 확인" | `for r in <5 repos>: test -f BRANDING.md && test -f docs/family.md; done` |
| G3 | 5 저장소 모두 4-lang README + canonical 11 docs (`en + ko + ja + zh`) | "GitHub 첫 페이지 lang switcher 4개 모두 클릭 가능 + 본문 존재" | `for r in <5 repos>: test -f README.md README.ko.md README.ja.md README.zh.md; done` + CRD description i18n |
| G4 | 3 operators (postgres/mongo/valkey) `operator-commons v0.8.0` consume | `go.mod` 의 `keiailab/operator-commons` 가 `v0.8.0` | `for r in <3 repos>: grep 'operator-commons v0.8' go.mod; done` |
| G5 | 5 저장소 모두 stale 브랜치 0 (main/gh-pages 외) | "Branches 탭 = 2-3개만 (main, gh-pages, [stable])" | `for r: git branch -r \| grep -v -E 'main\|gh-pages\|HEAD\|stable' \| wc -l = 0` |
| G6 | 5 저장소 main 모두 `required_status_checks = 0` | "PR 가 BLOCKED 없이 머지 가능" | `for r: gh api repos/keiailab/$r/branches/main/protection \| jq '.required_status_checks.contexts \| length' = 0` |
| G7 | RFC 0005 (§2 multi-arch 조항 제거) 본문 Accepted | "CLAUDE.md §2 에 multi-arch 금지 조항 *없음*" | `grep -c '멀티아키텍처 빌드' ~/.claude/CLAUDE.md = 0` (해당 문장 삭제) |
| G8 | forgewise GitHub repo metadata `licenseInfo.key = "apache-2.0"` | "Repo 카드의 License 뱃지 = Apache-2.0" | `gh repo view keiailab/forgewise --json licenseInfo \| jq -r .licenseInfo.key = apache-2.0` |
| G9 | 임시 파일 (HANDOFF.md / .claude/worktrees/ 잔존물 / docs/plans/ 완료분) 모두 정리 | "각 repo 의 root + docs/ 가 깨끗" | 각 sub-project 종료 게이트 |
| G10 | 모든 변경에 commit/PR 흔적 (`Signed-off-by` + `Co-Authored-By: Claude`) + 검증 로그 인용 | "PR body 에 통과 명령 + 출력 인용 존재" | PR template + 리뷰어 확인 |

### 2.2 Non-Goals (본 portfolio spec 의 OOS)

- ❌ 신규 기능 추가 (PR #158 의 ready msg / #159 의 PDB delete 같은 *기존* PR 의 머지는 포함되나, 새 기능 spec 작성은 별 cycle)
- ❌ 다른 keiailab 저장소 (org 의 8 public + 39 private 중 본 5건만 본 spec)
- ❌ 사용자 정의 plugin / agent SDK 작업
- ❌ MCP server 신규 빌드 (forgewise 의 *기존* MCP-native 기능 유지만, 신규 MCP 추가는 OOS)
- ❌ 운영 환경 (production cluster) 배포 — kind cluster e2e 까지만
- ❌ 외부 marketing / 웹사이트 / 블로그 (BRANDING.md 표준화는 포함, 외부 호스팅은 OOS)

## 3. 결정된 정책 요약 (Decisions Locked)

본 brainstorming session 의 사용자 답변으로 확정:

| ID | 정책 | 결정 | 출처 |
|---|---|---|---|
| D1 | GitHub Actions 처리 | **(v2.0 amendment 2026-05-21)** retention + 이중 운영 (per-repo 결정 가능). RFC 0002 의 "GHA 영구 금지" 는 GitLab CE L5 native CI 의 1차 결정 보장 + GitHub 의 OSS 가시성 유지 정책 으로 *재해석* 됨. valkey-operator `docs/specs/2026-05-21-pr-cleanup-and-gha-retention-design.md` 가 reference. 단 dependabot github_actions 자체는 *전수 머지* + 정기 갱신. | AskUser Q1 + valkey ADR-0048 |
| D2 | Sub-project 우선순위 | operator-commons (이미 v0.8.0 완료) → 3 operators 통합 → forgewise | AskUser Q2 |
| D3 | Multi-arch 정책 | **CLAUDE.md §2 multi-arch 금지 조항 제거** + PR #157 머지 (`platforms-amd64-guard` 도 제거). RFC 0005 사후 문서화. | AskUser Q3 |
| D4 | forgewise 라이센스 | Apache-2.0 통일 (LICENSE 파일 이미 Apache-2.0; GitHub detection 만 "Other" — SPDX header 보강) | AskUser Q4 |
| D5 | Commit author | `TaeHwan Park <eightynine01@gmail.com>` + `Co-Authored-By: Claude` trailer | AskUser Q5 (재확인) |
| D6 | dependabot 처리 | 전수 머지 후 GHA 제거 (`@dependabot recreate` 로 grouped → PR 수동 머지 → 그 다음 `.github/workflows/` git rm) | AskUser Q6 (재확인) |
| D7 | 작업 디렉토리 | `/Users/phil/workspace/keiailab/<repo>/` | AskUser Q7 |
| D8 | i18n 적용 범위 | **문서 + CRD description** (error/log 는 영어 유지) | AskUser Q8 |
| D9 | Portfolio spec 위치 | `operator-commons/docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md` | AskUser Q9 |
| D10 | Sub-project spec 위치 | 각 repo 의 `docs/superpowers/specs/<date>-<topic>-design.md` (postgres g1 spec 패턴 / valkey S1 spec 패턴) — *단* 기존 valkey S1 은 `docs/specs/` 에 작성됨. 본 cycle 에서 `docs/superpowers/specs/` 로 이전 권장 (OOS 아님 — S1 진행 시 동시 이동) |

## 4. Sub-Project 구조 (S1~S7)

본 spec 은 7개 sub-project 의 master index 역할. 각 sub-project 는 별도 spec/plan 으로 분해.

### 4.1 의존 그래프

```
                      (이미 완료) commons v0.8.0 cut
                                       │
                  ┌────────────────────┼──────────────────────┐
                  ▼                    ▼                      ▼
            ┌─────────┐         ┌──────────┐          ┌──────────────┐
            │  S2     │         │  S1      │          │   S6         │
            │ commons │         │  valkey  │          │  forgewise   │
            │ stale + │         │  PR + GHA│          │  표준화 +     │
            │ i18n    │         │  제거    │          │  i18n        │
            └────┬────┘         └────┬─────┘          └──────┬───────┘
                 │                   │  (사용자 worktree에서 시작) │
                 │                   │                       │
                 └────────────────┐  │  ┌────────────────────┘
                                  ▼  ▼  ▼
                              ┌─────────────────┐
                              │      S5         │
                              │  3 operators    │
                              │  v0.8.0 consume │
                              │  (mongo PR#190  │
                              │   이미 OPEN)    │
                              └────┬────────────┘
                                   │
                                   ▼
                              ┌─────────────────┐
                              │      S7         │
                              │  postgres+mongo │
                              │  GHA 제거       │
                              │  (S1 패턴 재사용) │
                              └────┬────────────┘
                                   │
                                   ▼
                              ┌─────────────────┐
                              │      S3         │
                              │  branding 마무리│
                              │  (valkey #161   │
                              │   + commons +   │
                              │   forgewise)    │
                              └────┬────────────┘
                                   │
                                   ▼
                              ┌─────────────────┐
                              │      S4         │
                              │  i18n 본문      │
                              │  ja/zh + CRD    │
                              └─────────────────┘

                      (병행 거버넌스)
                      ┌─────────────────┐
                      │   RFC 0005      │
                      │  §2 multi-arch  │
                      │  금지 조항 제거 │
                      └─────────────────┘
```

### 4.2 Sub-Project 카드

#### S1 — valkey-operator PR Cleanup + GHA retention

| 항목 | 값 |
|---|---|
| 상태 | **진행 중 (2026-05-21)** — 21 open PR → 3 잔존 (#164/#165/#166 v0.8.0 consume + ja/zh README native), 16+ stale → 12, GHA **retention** 노선 (ADR-0048 Accepted = per operator family trade-off, D1 v2.0 amendment 정합) |
| 기존 spec | `docs/specs/2026-05-21-pr-cleanup-and-gha-retention-design.md` (v2.0 retention 노선) — PR #163 머지 |
| PR 머지 (오늘 2026-05-21) | #138 (P-C.6 OTel docs) · #141~#156 (dependabot 전수) · #158 (TLS immutable + ready msg) · #159 (PDB delete) · #160 (3-tier 분류) · #161 (BRANDING Wave 3) · #163 (S1+ spec) · #167 (ADR-0048 integrated) · #168 (actions/stale bump) · #169 (ADR-0048 Accepted) · #170 (Sprint 1 P2 commons -322 LOC) · #171 (lefthook 3종 hook P1-11/12/13) · #172 (helm-publish.sh) · #173 (UPGRADING.md) · #174 (ADR-0050 audit augmentation) |
| ADR 변경 | ADR-0048 Status Proposed → **Accepted** (commit `f2d9ee4`, per operator family v2.0 retention 결정) · ADR-0049 (Sprint 1 commons P2 adoption) · ADR-0050 (audit augmentation) |
| 본 spec 의 변경 사항 (D1 v2.0 amendment 적용) | (a) PR #157 multi-arch **별 cycle 로 분리** (cluster 1대 시범 검증 후 머지 결정); (b) GHA 제거 → **retention 14 wf** (ADR-0048 정합); (c) `platforms-amd64-guard` 제거는 RFC 0006 (local commit `0a98641`, push 보류) 의 별 결정 |
| 잔여 | 3 open PR (#164 v0.8.0 consume / #165 ja README / #166 zh README), 12 stale, 1 issue (#4 Renovate) |
| 후속 | S5-valkey 의 commons v0.9.0 consume 는 main 에서 직접 진행 (Sprint 1 Phase 2 PR `7c8b128`+`daa763c` 완료) |
| Cross-link | [valkey docs/specs/2026-05-21-pr-cleanup-and-gha-retention-design.md](https://github.com/keiailab/valkey-operator/blob/main/docs/specs/2026-05-21-pr-cleanup-and-gha-retention-design.md) |

#### S2 — operator-commons stale 정리 + i18n placeholder → 본문

| 항목 | 값 |
|---|---|
| 상태 | **Implemented (2026-05-21)** — T1+T2+T3+T4 모두 완료. [verification](../plans/2026-05-21-s2-verification.md) |
| PR 머지 | #36 (portfolio spec) · #37 (S2 plan) · #38 (S2 sub-spec) · #42 (T1 archive + ADR-0009) · #43 (T1 부속 tag + ADR-0010) · #44 (T2 v0.8.0 drift) · #45 (S2 Phase 2 lefthook + ADR-0011) · #46 (T3-zh README native 125 LOC) · #50 (T5 verification + spec Status Implemented) · 사용자 평행: #40 (T3-ja README native 122 LOC) · #47 (T4 + S4 Phase 1 흡수: commons SSOT glossary 4-lang) |
| 작업 | (a) `archive/main-13-commits-merge-style-2026-05-21` 브랜치 분석 — tag 로 보존 후 브랜치 삭제 (ADR-0009/0010); (b) `docs/family.md` 의 `v0.7.0` → `v0.8.0` drift 해소; (c) `README.ja.md` + `README.zh.md` placeholder → 본문 native; (d) `docs/i18n/glossary-ja.md` + `glossary-zh.md` placeholder → 본문 (203/207/207 LOC) |
| 산출 | commons 의 stale = 0, i18n 4-lang 본문 완성, family.md drift 해소 |
| 후속 | S5 가 commons 의 *안정 상태* 를 가정 (v0.8.0 → v0.9.0 진전) |
| Branch | `feat/s2-commons-stale-i18n-2026-05-21` (`feat/s4-i18n-master-2026-05-21` 등 사용자 평행 작업 흡수) |
| Cross-link | [docs/superpowers/plans/2026-05-21-s2-verification.md](../plans/2026-05-21-s2-verification.md) |

#### S3 — Branding Wave 3 마무리 (5 repos consistency)

| 항목 | 값 |
|---|---|
| 상태 | **거의 완료 (2026-05-21)** — 5/5 repos 모두 BRANDING.md + docs/family.md 존재 (`test -f` PASS), valkey #161 머지 확인됨 |
| PR 머지 | postgres #82 (Wave 3 BRANDING/family) · #83 (4-lang README ja/zh native 237 LOC) · mongo #188 (Wave 3 BRANDING/family) · #189 (4-lang README ja/zh placeholder) · #191 (ja native 557 LOC) · #192 (zh native 566 LOC) · commons #35 (Wave 3 BRANDING + 13 root .md footer) · valkey #160 (3-tier docs) · #161 (BRANDING Wave 3) · forgewise #9 (BRANDING 4-lang) · #10 (docs/family 4-lang) · #4 (거버넌스 5종) · #11 (README ja/zh 본문) |
| 산출 | 5 repos 모두 README header/footer 100% 정합 + BRANDING + family 일관. 검증: `for r in <5 repos>: test -f BRANDING.md && test -f docs/family.md` 5/5 PASS |
| 잔여 | (확인 필요) forgewise 의 GitHub `licenseInfo.key` detection 수정 = G8 (S6 와 정합) |

#### S4 — i18n 본문 (ja/zh) + CRD description 4-lang

| 항목 | 값 |
|---|---|
| 상태 | **Implemented (2026-05-21)** — commons SSOT + 5 repos README + lefthook hook + sync 스크립트 모두 완성 (D1~D5 결정 본 sub-spec 자체에 명시) |
| PR 머지 | commons #39 (S4 i18n 마스터 spec D1~D5 결정) · #47 (Phase 1 commons SSOT 정비: glossary 4-lang 203/207/207 LOC + sync 매트릭스) · #48 (Phase 2 readme-i18n-sync lefthook pre-push hook) · #49 (Phase 3 sync-from-commons.sh 신규) · #51 (Phase 4 commons docs/ P0 번역: family + STABILITY + coverage-report) · #53 (Phase 5 top-level P1: BRANDING 3-lang + ADOPTERS 2-lang) · 사용자 평행 (각 repo): postgres #83 (4-lang ja/zh native) · mongo #189 + #191 + #192 (4-lang) · valkey #165/#166 (3 open PR — ja/zh native, S1 마무리 시 흡수) · forgewise #1+#10+#11 (4-lang README + family + 본문) |
| 작업 | (a) commons SSOT glossary 4-lang 본문; (b) 5 repos README 4-lang (en canonical + ko + ja + zh); (c) readme-i18n-sync pre-push hook (drift 차단); (d) sync-from-commons.sh (SSOT → 4 repo 배포); (e) D1~D5 결정: en canonical / 4-lang scope / kubectl explain CRD 영어 유지 / glossary SSOT in commons / sync hook 강제 |
| 산출 | commons 의 docs/i18n/ 4-lang SSOT + 5 repos README 4-lang + drift 자동 차단 hook |
| 후속 | CRD description annotation 방식 (`keiailab.com/description.ko`) 은 *별 cycle* 로 분리 (D5 결정: 우선 README + commons docs/ 본문 완성) |
| Cross-link | [commons S4 spec](https://github.com/keiailab/operator-commons/pull/39) — D1~D5 결정 spec 자체에 명시 |

#### S5 — 3 operators commons v0.8.0 consume → v0.9.0 진전

| 항목 | 값 |
|---|---|
| 상태 | **Implemented + v0.9.0 진전 (2026-05-21)** — 3 operators 모두 v0.9.0 까지 진전 (v0.8.0 초과). Sprint 1 Phase 2 (pkg/pvc + pkg/topology 신규 추출) 도 3 operators 모두 채택 완료 |
| PR 머지 | mongo #190 (commons v0.7.0 → v0.8.0 + pkg/probes 2 Exec site, commit `cb888b7`) · postgres #84 (commons v0.7.0 → v0.8.0 + pkg/probes 2 HTTP site, commit `c72b435`) · valkey #164 (3 open PR 중 1건, v0.7.0 → v0.8.0 + pkg/probes 2 Exec, in-flight) · v0.9.0 bump (3 operators 모두): mongo `a1ff450` · postgres `4abf932` · valkey `daa763c` · Sprint 1 P2: commons #52 (pkg/pvc + pkg/topology 신규 추출) · postgres #91 (-375 LOC) · mongo #202 (-327 LOC) · valkey #170 (-322 LOC) |
| 작업 | (a) 3 operators v0.8.0 consume + pkg/probes 적용 사이트 식별 + 적용; (b) v0.9.0 release cut (Sprint 1 commons release); (c) Sprint 1 Phase 2 pkg/pvc + pkg/topology 신규 추출 + 3 operators 채택 (-1,024 LOC 합산); (d) docs/family.md 의 `v0.7.0+` → `v0.8.0+` → `v0.9.0+` 진전 갱신 |
| 산출 | 3 operators 모두 commons v0.9.0 consume + Sprint 1 P2 채택 + 합산 -1,024 LOC 중복 해소 |
| 후속 | Sprint 1 Phase 3 (pkg/coordinator 신규 추출) — *별 cycle* |
| Branch | (각 repo) `feat/sprint-1-pvc-topology-extraction-2026-05-21` · `feat/sprint-1-phase-2-pvc-topology-adopt-2026-05-21` |

#### S6 — forgewise 표준화 + 다국어 + LICENSE detection 수정

| 항목 | 값 |
|---|---|
| 상태 | **거의 완료 (2026-05-21)** — 5 표준 docs 모두 추가, 4-lang README + 4-lang BRANDING/family, audit lefthook P0-9/P1-13 보강. 잔여: GitHub `licenseInfo.key` detection 수정 (G8) |
| PR 머지 | #1 (README 4-lang en canonical SSOT + ko rename + ja/zh placeholder) · #2 (S6 design doc) · #3 (lefthook Python 4 계층 게이트) · #4 (거버넌스 5종 ko: SECURITY/CONTRIBUTING/COC/CHANGELOG/AGENTS) · #5 (S6 plan INDEX.md) · #6 (운영 4종 신규: Installation/Configuration/API Reference/Upgrade) · #7 (거버넌스 5종 + 운영 4종 + lefthook 키워드 강제) · #8 (ADR-0001 Python 스택 override) · #9 (BRANDING 4-lang) · #10 (docs/family 4-lang) · #11 (README ja/zh placeholder → 본문 완전 번역) · #12 (운영 4종 ja/zh) · #13 (commons SSOT sync + readme-i18n-sync lefthook) · main HEAD: `aad600f feat(audit): lefthook P0-9/P1-13 보강 + release.sh + Makefile` · CODEOWNERS / ADOPTERS / ROADMAP / PR template 사용자 평행 추가 |
| 작업 | (a) BRANDING.md / docs/family.md / CONTRIBUTING / SECURITY / COC / CHANGELOG / AGENTS / 운영 4종 / lefthook 4 계층 / 4-lang README/docs / readme-i18n-sync 자동화; (b) LICENSE 파일 Apache-2.0 정합 (LICENSE 파일 자체는 정상) |
| 산출 | forgewise 가 다른 4개와 동일한 5 표준 docs + 4-lang README + lefthook + commons SSOT sync hook |
| 잔여 | GitHub `licenseInfo.key = "other"` → "apache-2.0" detection 수정 (G8) — 별 cycle 또는 audit augmentation |
| Branch | (각 PR 별) `feat/i18n-en-canonical-2026-05-21` · `feat/governance-5-set-2026-05-21` · etc. |

#### S7 — postgres + mongodb GHA per-repo 정합 (D1 v2.0 amendment)

| 항목 | 값 |
|---|---|
| 상태 | **per-repo 정합 (2026-05-21)** — postgres = strict 제거 노선 (14 → 3 narrow exception); mongo = retention 노선 (12 wf 유지, race 후 re-restore); valkey = retention 노선 (14 wf 유지, ADR-0048). 5 repos: commons=0 (애초) / postgres=3 / mongo=12 / valkey=14 / forgewise=0 |
| PR 머지 (postgres strict 노선) | #85 (ADR-0017 integrated GHA retention rationale 1차) · #86 (RFC-0002 .github/workflows 전체 제거 14 파일) · #87 (helm-publish.sh + release.sh) · #88 (로컬 4계층 보강 kube-linter + go-licenses + markdown-link-check) · #89 (ADR-0018 Accepted GHA 전면 제거) · #90 (revert restore 14 wf race) · #91 (Sprint 1 P2 -375 LOC) · #92 (ADR-0019 Accepted GHA 유지 v2.0 정합) · #93 (RFC-0002 GHA cleanup 11 wf 제거, narrow exception 3 유지) · #94 (broad .claude gitignore) · #95 (RFC-0002 gha-block hook + UPGRADING.md) · ADR-0022 (`beb4d42`, narrow exception 3 wf 정합) |
| PR 머지 (mongo retention 노선) | #193 (ADR-0031 integrated GHA retention rationale 1차) · #194 (RFC-0002 .github/workflows 전체 제거 12 파일) · #195 (helm-publish.sh) · #196 (로컬 4계층 보강) · #197 (ADR-0032 Accepted GHA 전면 제거 — Superseded) · #199 (revert restore 12 wf) · #200 (RFC-0002 GHA cleanup 9 wf race 시도) · #201 (.claude-flow gitignore) · #202 (Sprint 1 P2 -327 LOC) · #203 (ADR-0033 Accepted GHA 유지 v2.0 정합) · #204 (revert re-restore 9 wf, race 정정) · #205 (RFC-0002 gha-block hook + ADR-0035) |
| ADR | postgres: ADR-0017 (proposed) / ADR-0018 (Accepted GHA 제거) / ADR-0019 (Accepted GHA 유지, ADR-0018 Superseded) / ADR-0020 (Sprint 1 P2) / ADR-0021 (gha-block hook) / ADR-0022 (narrow exception 3 wf 정합). mongo: ADR-0031~0035 동일 패턴 |
| 작업 | (a) 각 repo per-repo 결정 — postgres = 11 wf 제거 + 3 narrow exception 유지 (pages / dependabot / release tag); (b) mongo = 12 wf retention (operator family v2.0 정합); (c) valkey S1 에서 14 wf retention (ADR-0048); (d) commons = GHA 0 유지; (e) forgewise = GHA 0 유지 (애초 부재); (f) RFC-0002 gha-block pre-commit hook 3 repos 모두 적용 |
| 산출 | per-repo 정합 (cleanup 의도 = 각 repo 별 결정 보장 달성) + RFC-0002 narrow exception 3종 표준화 + 로컬 4계층 보강 (kube-linter + go-licenses + markdown-link-check) 3 operators 동등 |
| Branch | postgres `feat/rfc-0002-gha-cleanup-2026-05-21` 등 / mongo `revert/restore-workflows-2026-05-21` 등 |

### 4.3 거버넌스 변경 — RFC 0005 (가칭)

별도 RFC 작성 (위치: `~/.claude/rfcs/0005-multi-arch-policy-removal.md`):

- **Title**: §2 multi-arch 빌드 금지 조항 제거 + 부트스트랩 정당화
- **Status**: Draft → Proposed → Accepted (사용자 답변이 §7 부트스트랩 예외 트리거)
- **Body**: 
  - 사용자 결정 인용 (AskUser Q3 답변)
  - 영향 범위: 5 repos 의 lefthook / Makefile / Dockerfile / OLM bundle
  - valkey PR #157 머지가 reference implementation
  - 부트스트랩 예외 정당화 (§7)
- **Refs**: 본 spec, S1 (valkey PR cleanup), valkey PR #157

## 5. 실행 순서 (Build Sequence)

원칙: D2 (commons → operators → forgewise) + S 의 의존 그래프.

| 순서 | Sub-project | 진입 조건 | 예상 atomic commit 수 | 평행 가능? |
|---|---|---|---|---|
| 1 | S2 (commons stale + i18n placeholder) | 본 spec 승인 | 4~6 | — |
| 2 | RFC 0005 (multi-arch 정책 제거) | 본 spec 승인 + AskUser Q3 답변 | 1 (CLAUDE.md 직접 수정 + RFC 작성) | S2 와 평행 가능 |
| 3 | S1 (valkey PR cleanup + GHA 제거) | RFC 0005 Accepted + S1 의 기존 worktree spec 갱신 | 7 phases × 평균 3 commits = 약 21 | — |
| 4 | S5 (3 operators v0.8.0 consume) | S1 종료 (valkey BLOCKED 해제) | 3 PRs (각 operator 1개) | mongo #190 + postgres + valkey 평행 가능 |
| 5 | S7 (postgres + mongo GHA 제거) | S5 완료 (각 operator 가 안정) | 2 PRs (postgres + mongo) | postgres / mongo 평행 가능 |
| 6 | S3 (branding 마무리) | S1 (valkey #161 머지) + S6 진입 가능 | 2 PRs (valkey + forgewise) | — |
| 7 | S6 (forgewise 표준화) | S3 의 branding 표준 확정 | 1 PR + LICENSE detection 별 처리 | — |
| 8 | S4 (i18n 본문 ja/zh + CRD) | S2 (commons glossary), S6 (forgewise i18n) | 5 PRs (각 repo) + 1 SPEC 갱신 | 모든 repo 평행 가능 |

총 예상 atomic commit 수: **40~50개** · 총 PR 수: **15~20개** · 기간: writing-plans 산출 시점 결정.

## 6. 리스크 & 완화

| 리스크 | 영향 | 완화 |
|---|---|---|
| valkey S1 의 기존 worktree spec 과 본 portfolio spec 의 *PR #157 처리 충돌* | 두 spec 의 결정이 어긋남 (기존 worktree spec = "별 cycle 권장", 본 spec = D3 머지) | S1 진입 시 기존 spec 의 §4 Phase 3 + §7 OOS 행을 *본 spec 의 D3* 로 갱신. spec amendment commit 1건. |
| `archive/main-13-commits-merge-style-2026-05-21` 의 *비-ancestor* 가능성 (rebase merge 후 그래프 불일치) | S2 진행 시 commit 손실 위험 | 삭제 전 `git log --left-right --count origin/main...origin/archive/...` 로 ahead/behind 정밀 확인. ahead > 0 이면 tag 보존 후 브랜치 삭제. |
| mongodb PR #190 의 BLOCKED 가 commons v0.8.0 *때문* 이 아니라 GHA required_status_checks 때문 | S5-mongo 가 S7-mongo 의존 → 순서 강제 | S5 진입 전에 mongo branch protection 확인. BLOCKED 면 S7-mongo 부분 (workflow 제거 + protection 갱신) 을 먼저 atomic 적용. |
| RFC 0005 부트스트랩 정당화 실패 (§7 예외 조건 미충족) | 본 spec 의 multi-arch 결정 무효화 | 사용자 답변 (AskUser Q3) 인용 + 본 spec commit message 에 "사용자 명시 지시" 명기. RFC 0005 사후 작성으로 *Accepted* 상태 직접 진입. |
| forgewise LICENSE GitHub detection 미수정 | G8 미달 | GitHub Linguist 의 license detection 알고리즘 (https://github.com/github-linguist/licensee) 확인 후 LICENSE 파일 형식 정밀 조정. context7 MCP 활용. |
| 5 repos 작업 중 *한 곳* 의 lefthook 보강이 안 됐을 때 push 차단 → 다른 repo 작업 영향 | bottleneck | 각 sub-project 의 첫 단계 = lefthook 4계층 보강 검증 (G6 의 부분 게이트). 보강 commit 은 *반드시* 별 commit. |
| CRD i18n (S4 의 c~d) 가 kubebuilder/controller-gen 비호환 | S4 부분 실패 | S4 진입 시 context7 MCP 로 kubebuilder 1.x 의 i18n 정책 확인. 비호환 시 *annotation 방식 우회* (`keiailab.com/description.ko: "..."` 형식) — kubectl explain 에 영향 없음 |
| 임시 파일 잔존 (특히 S1 의 `.claude/worktrees/`, S3~S7 의 `docs/plans/` 완료분) | G9 미달 | 각 sub-project 종료 시 cleanup 단계 (CLAUDE.md §8 의 "Phase 5 cleanup" 패턴) 강제 |
| 본 portfolio spec 의 atomic 미보장 (5 repos × 7 sub-project = state explosion) | 부분 완료 후 사용자 break 시 회복 불가 | 각 sub-project = *독립 PR* . 본 spec 의 §5 build sequence 의 *한 칸* 만 in-progress 허용. 동시 in-progress = 사용자 명시 승인 시만. |

## 7. 성공 기준 (Acceptance — bash 검증)

본 portfolio spec 은 다음 *모두* 충족 시 *완료*:

```bash
# G1. 5 repos 모두 .github/workflows 부재
cd /Users/phil/workspace/keiailab
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  test ! -d "$r/.github/workflows" && echo "✓ $r: no workflows" || echo "✗ $r: workflows exist"
done

# G2. 5 repos 모두 BRANDING.md + docs/family.md 존재
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  test -f "$r/BRANDING.md" && test -f "$r/docs/family.md" && echo "✓ $r" || echo "✗ $r"
done

# G3. 5 repos 모두 4-lang README + canonical 11 docs
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  for f in README.md README.ko.md README.ja.md README.zh.md; do
    test -f "$r/$f" && true || { echo "✗ $r/$f"; exit 1; }
  done
done && echo "✓ 4-lang README 5/5"

# G4. 3 operators commons v0.8.0
for r in postgres-operator mongodb-operator valkey-operator; do
  grep -q 'operator-commons v0\.8\.' "$r/go.mod" && echo "✓ $r v0.8.x" || echo "✗ $r"
done

# G5. stale 브랜치 0
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  cd "/Users/phil/workspace/keiailab/$r"
  N=$(git branch -r | grep -v -E 'main|gh-pages|HEAD|stable' | wc -l | tr -d ' ')
  echo "  $r stale=$N"
done

# G6. required_status_checks = 0
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  N=$(gh api "repos/keiailab/$r/branches/main/protection" 2>/dev/null | jq '.required_status_checks.contexts | length' 2>/dev/null || echo 0)
  echo "  $r required_checks=$N"
done

# G7. RFC 0005 적용 검증
grep -c '멀티아키텍처 빌드' ~/.claude/CLAUDE.md  # 기대: 0
test -f ~/.claude/rfcs/0005-multi-arch-policy-removal.md

# G8. forgewise GitHub license
gh repo view keiailab/forgewise --json licenseInfo | jq -r .licenseInfo.key  # 기대: apache-2.0

# G9. 임시 파일 정리
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  cd "/Users/phil/workspace/keiailab/$r"
  test ! -e .claude/worktrees/ || echo "✗ $r worktrees 잔존"
  test ! -e HANDOFF.md || echo "✗ $r HANDOFF.md root 잔존 (docs/internal/ 이전 필요)"
  test ! -d docs/plans/ || { CNT=$(find docs/plans -type d | wc -l | tr -d ' '); echo "  $r docs/plans subdirs: $CNT"; }
done

# G10. 모든 PR body 에 검증 명령 + 출력 인용 존재
gh pr list -R keiailab/operator-commons --state merged --search "created:2026-05-21" --json number,body --limit 50 |
  jq -r '.[] | select(.body | contains("verify") | not) | .number'
# 출력 0 = 모든 PR 가 verify 인용
```

각 sub-project 의 종료 게이트 = 위 검증의 *해당 부분* PASS + sub-project 자체의 verification.md.

## 8. 임시 파일 정리 정책 (사용자 요구 "임시 파일 반드시 클린제거")

본 supercycle 중 생성되는 임시 파일 카탈로그:

| 종류 | 위치 | 정리 시점 |
|---|---|---|
| `.claude/worktrees/<topic>/` | 각 repo | sub-project 완료 즉시 (`git worktree remove`) |
| `docs/plans/<topic>/` | 각 repo | sub-project 완료 + 핵심 결정 ADR 승격 후 (`git rm -r`, CLAUDE.md §8 cleanup 단계) |
| `docs/plans/<topic>/research/extracted-*.md` | 각 repo | atomic task 완료 시점 partial cleanup (해당 row done) |
| `HANDOFF.md` root | 각 repo | docs/internal/ 이전 (Wave 2 작업, 이미 완료) — root 잔존물 0 |
| `/tmp/.go.mod.lefthook.bak`, `/tmp/.go.sum.lefthook.bak` | 시스템 임시 | lefthook 자체 cleanup, 본 spec scope 외 |
| `$CLAUDE_JOB_DIR/*` (`/Users/phil/.claude/jobs/4829bea1/`) | 시스템 | 잡 종료 시 자동 cleanup |
| 본 spec 작성 중 생성된 `/Users/phil/.claude/jobs/4829bea1/repo-inventory.txt`, `repo-deep-state.txt` | 시스템 | 본 spec commit 후 사용자 review 기간 보존, 그 후 자동 cleanup |

각 sub-project 의 종료 게이트 = `find . -type d -name "worktrees" -o -name "plans" | xargs -I {} echo "cleanup_pending: {}"` 결과 = 빈 라인.

## 9. 본 spec 의 OOS (별 sub-project 또는 후속 cycle)

| 항목 | 분류 | 이유 |
|---|---|---|
| keiailab org 의 8 public + 39 private 중 본 5건 외 저장소 | 후속 cycle | 본 spec 명시 범위 5건만 |
| OLM bundle catalog publish | S1 OOS | S1 의 §2.2 와 정합 |
| keiailab.com 웹사이트 / 마케팅 페이지 | 별 cycle | 외부 호스팅, 본 spec scope 외 |
| pgo / Citus / Sentinel 등 *upstream* operator 와의 직접 호환성 | 별 cycle | family.md "What we do NOT do" 와 정합 (no embed/wrap) |
| 신규 RFC (RFC 0006+) | 별 cycle | 본 spec 은 RFC 0005 만 부트스트랩 |
| Plugin 개발 (Claude Code plugin / Codex plugin) | 별 cycle | 본 spec 은 5 repos 만 |
| 컨테이너 image SBOM 생성 자동화 (cosign / syft) | 후속 cycle | S1 의 audit 단계가 1차 방어; 본 supercycle 외 |
| `forgewise` 의 MCP server 신규 기능 | 별 cycle | S6 는 표준화만 포함, 기능 변경 OOS |
| 3 operators 의 운영 환경 (production cluster) 배포 검증 | 별 cycle | kind cluster e2e 까지가 G6 의 한계 |

## 10. 다음 단계 (Immediate)

1. **본 spec 의 사용자 검토** (현 task).
2. 승인 후 `superpowers:writing-plans` skill 호출 → **S2 (commons stale + i18n placeholder)** 의 implementation plan 작성 (`operator-commons/docs/plans/s2-commons-cleanup/INDEX.md` + `research/*.md`).
3. S2 plan 의 atomic task 실행 → 1 commit + 1 PR + 1 머지 + 브랜치 삭제 cycle (CLAUDE.md §8).
4. S2 완료 즉시 RFC 0005 작성 (S1 진입 전 거버넌스 선결).
5. RFC 0005 Accepted → S1 진입 (기존 worktree spec amendment + 7 Phase 실행).
6. 이후 §5 build sequence 의 순서대로 S5 → S7 → S3 → S6 → S4 진행.

각 단계는 *별 turn* 의 사용자 명시 발화 ("S2 진행", "S1 진행" 등) 로 진입 — CLAUDE.md §8 의 "Phase 3+ 는 매 turn 사용자 발화 시점이 진입점" 정합.

## 부록 A. 정의 (Definitions)

- **cleanup supercycle**: 2026-05-21 에 시작된 5-repo 정합화 multi-Wave 작업. Wave 1 (배경 정리) → Wave 2 (3-tier 분류) → Wave 3 (branding) → Wave 4 (i18n) → Wave 5 (cross-validation) → 본 portfolio spec 이 Wave 5 의 master index 역할.
- **atomic commit/PR/머지**: CLAUDE.md §8 — task 1개 = 1 commit + 1 ship + 1 deploy.
- **로컬 4계층** (RFC 0002): pre-commit hook · pre-push hook · Makefile target · PR 리뷰어 증거 확인.
- **stale 브랜치**: main / gh-pages / stable (postgres) 이외의 origin 브랜치.
- **S 표기**: 본 spec 의 7 sub-project 식별자 (S1~S7). 사용자 valkey worktree spec 의 §7 에 등장한 표기 정합.
- **부트스트랩 예외**: CLAUDE.md §7 — 사용자 *일회성 명시 지시* 1회로 RFC 절차 부트스트랩 가능. 본 spec 의 D3 가 RFC 0005 부트스트랩 트리거.

## 부록 B. 본 spec 의 자체 검증 (Self-Review)

작성 직후 다음 점 자체 검증 완료:

- [x] Placeholder 없음 (모든 TBD/TODO 자리 명시 결정)
- [x] 내부 모순 없음 (S1 = 사용자 기존 spec amendment 명시, D3 와 정합)
- [x] OOS 명확 (§9 7건 명시)
- [x] 모호한 요구 없음 (각 sub-project 의 산출 + 검증 명령 명시)
- [x] 사용자 결정 인용 (D1~D10 명시 출처 = AskUser Q1~Q9)
- [x] 한국어 작성 (CLAUDE.md §2)
- [x] keiailab footer 부착 (아래)
- [x] 검증 명령 bash 형식 (CLAUDE.md §2 "통과 로그·핵심 출력")

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
