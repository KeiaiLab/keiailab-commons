# S4 — 다국어 4-lang 마스터 spec (operator-commons SSOT) — Design

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Proposed (사용자 결정 5건 미해결 — §7 참조) |
| 작성자 | keiailab — superpowers cycle (post-supercycle drill-down) |
| 범위 | **마스터 spec** — 5 저장소 (commons + 3 operators + forgewise) 의 i18n 정책 통일 + SSOT 자산 마련 + 자동 번역 파이프라인 + lefthook 통합 |
| Supercycle | [`2026-05-21-portfolio-cleanup-supercycle-design.md`](../superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md) 의 **Sub-Project S4** |
| 후속 | 본 spec Accepted 후: S4-A (postgres) / S4-B (mongodb) / S4-C (valkey) / S4-D (commons) / S4-E (forgewise) 의 *실제 번역 실행 sub-spec* 5건. 본 spec 은 마스터 정책만. |
| Commit 작성자 | `TaeHwan Park <eightynine01@gmail.com>` + `Signed-off-by:` + `Co-Authored-By: Claude` |
| 우선순위 | High (portfolio supercycle G3 4-lang 완성 게이트 + 5 repo 통합 i18n 정합) |

## 1. 배경 (Background)

### 1.1 5 저장소 i18n 현황 스냅샷 (2026-05-21 13:15 KST)

직접 측정 (`find docs -name '*.md' | wc -l` + LOC 합산 + `ls README*.md`):

| Repo | README 4-lang 골격 | docs 번역대상 (count) | docs 번역대상 LOC | top-level 번역대상 (count, LOC) | 기존 ko 번역 | 기존 ja 번역 | 기존 zh 번역 |
|---|---|---|---|---|---|---|---|
| `operator-commons` | ✓ (en+ko+ja+zh) | 6 | 653 | 12, 1,240 | glossary-ko (203 LOC) | glossary-ja placeholder (102 LOC) | glossary-zh placeholder (102 LOC) |
| `postgres-operator` | ✓ (en+ko+ja+zh) | 30 | 5,991 | 13, 2,302 | 0 (README.ko 만) | 0 (README.ja 만) | 0 (README.zh 만) |
| `mongodb-operator` | ✓ (en+ko+ja+zh) | 32 | 11,916 | 17, 4,134 | 0 (README.ko 만) | 0 (README.ja 만) | 0 (README.zh 만) |
| `valkey-operator` | **✗** (en+ko, ja/zh 부재) | 33 | 5,070 | 18, 2,501 | **12개 docs/operations/\*.ko.md** | 0 | 0 |
| `forgewise` | ✓ (en+ko+ja+zh) | 5 | 297 | 1, 162 | 0 (README.ko 만) | 0 (README.ja 만) | 0 (README.zh 만) |
| **합계** | 4/5 | **106** | **23,927** | **61, 10,339** | 13 (commons glossary + valkey ops) | 1 placeholder | 1 placeholder |

**번역 대상 총량**: docs 106 + top-level 61 = **167 파일** × 3 추가 언어 = **501 신규 파일**. 평균 LOC 약 200 → 약 **100,000 LOC 신규 번역**.

(브리프 §컨텍스트 의 "138건 × 3 = 414" 추정치는 *번역 후보 식별 시점* 의 보수 추산. 본 spec 의 직접 측정 결과 = 167 × 3 = 501 파일 으로 갱신.)

### 1.2 번역 대상 vs 번역 비대상

**번역 대상** (사용자/외부 대상 문서):
- README.{md,ko,ja,zh}.md (5 repo)
- top-level governance/branding 문서 (BRANDING / GOVERNANCE / SECURITY / CONTRIBUTING / CODE_OF_CONDUCT / MAINTAINERS / ADOPTERS / ROADMAP / ARCHITECTURE / CHANGELOG / STABILITY 등)
- `docs/family.md` (cross-link)
- `docs/getting-started.md`, `docs/UPGRADING.md`, `docs/index.md`
- `docs/advanced/*.md` (mongodb 5건)
- `docs/comparison/*.md` (mongodb 3건)
- `docs/developers/*.md` (mongodb 4건) — *사용자 가시 문서로 판단*
- `docs/operations/*.md` (valkey 11건)
- `docs/i18n/glossary-{ko,ja,zh}.md` (commons SSOT)

**번역 비대상** (내부/governance/legal):
- `docs/kb/adr/*.md` (ADR — *결정 추적*. 영문 canonical 만 유지)
- `docs/kb/rfc/*.md` (RFC — 동일)
- `docs/kb/deps/*.md` (deps log — 자동 생성)
- `docs/internal/*.md` (HANDOFF / README / TASKS — 내부 운영자 전용)
- `docs/superpowers/*` (cycle artifacts — 일시적)
- `docs/specs/*.md` (design specs — 본 spec 자신 포함)
- `docs/plans/*.md` (implementation plans)
- `LICENSE`, `NOTICE`, `CITATION.cff`, `.gitignore` 등 (legal/config — 영문 유일)

### 1.3 기존 SSOT 자산 (operator-commons)

본 spec 작성 시점 (commit `bdb659c`) 에 commons 가 보유한 i18n SSOT:

- `docs/i18n/glossary-ko.md` (203 LOC) — 본문 완성
- `docs/i18n/glossary-ja.md` (102 LOC) — `[~]` placeholder (RFC-0025 §1.2 부분 구현 marker)
- `docs/i18n/glossary-zh.md` (102 LOC) — `[~]` placeholder
- `scripts/check-readme-sync.sh` — *EN ↔ KO 한정* drift check (section header 1:1 + line diff ≤ 15% + 양방향 cross-link)
- README header 의 4-lang switcher (4 repo 모두 4-lang 골격 — valkey 제외)
- `docs/family.md` footer 의 4-lang cross-link
- lefthook 에는 *현재* readme-i18n-sync hook 미통합 (`scripts/check-readme-sync.sh` 는 수동 실행만)

### 1.4 자동 번역 도구 흔적

본 spec 작성 시점에 5 repo 의 `scripts/` / `tools/` / 어디에도 자동 번역 파이프라인 부재. 모든 기존 번역은 *수동 작성* 또는 *placeholder*.

### 1.5 글로벌 거버넌스 정합 (CLAUDE.md)

| 조항 | 본 spec 의 적용 |
|---|---|
| §2 한국어 | 본 spec / 모든 sub-spec / commit message / ADR 본문 한국어. 자동 번역 결과물의 *원본* 은 영문 canonical → 4 언어 번역. |
| §2 테스트 없는 기능 = 존재 불가 | 자동 번역 파이프라인 (`scripts/i18n-translate.sh`) 에 단위 테스트 + drift 검증 hook 필수. |
| §2 context7 MCP 사용 | 자동 번역 엔진 선택 시 (DeepL / OpenAI / Claude / Google) 각 API 의 최신 공식 문서 조회 의무. |
| §2 GitHub Actions 영구 금지 (RFC-0002) | 자동 번역은 *로컬 4계층* 으로만 — `Makefile target` + `lefthook post-commit` (또는 pre-push). GHA 사용 절대 금지. |
| §5 git-flow 미사용 = 실패 | `spec/i18n-4lang-master-2026-05-21` branch → squash merge → branch 삭제. atomic 1 spec = 1 PR. |
| §5 범위 외 수정 = 실패 | OOS 5건 명시 (§8). 본 spec 은 *마스터 정책* 만; 각 repo 의 실제 번역 실행은 별 sub-spec. |
| §8 atomic 정책 | Phase 별 atomic commit. 단 본 spec 자체는 *문서 1건* 만 → PR 1건 commit 1건. |

## 2. 목표 (Goals) + 비목표 (Non-Goals)

### 2.1 Goals (사용자 시점 시나리오)

| ID | 목표 | 사용자 시점 검증 | 측정 |
|---|---|---|---|
| G1 | `docs/i18n/glossary-{ko,ja,zh}.md` 4-lang 완성 + native reviewer placeholder `[~]` → `[x]` 승격 SOP 정의 | "용어 사전 3개 모두 본문 완성 + 검증 절차 문서화" | `grep -c '^| ' docs/i18n/glossary-ja.md ≥ 30` + `grep -c '^| ' docs/i18n/glossary-zh.md ≥ 30` + `test -f docs/i18n/SOP.md` |
| G2 | `scripts/check-readme-sync.sh` → 4-lang 매트릭스 (EN ↔ KO/JA/ZH) 확장 + 5 repo 동기 배포 | "lefthook pre-push 시 4-lang drift 자동 차단" | `grep -q 'README.ja.md' scripts/check-readme-sync.sh` + `grep -q 'README.zh.md' scripts/check-readme-sync.sh` + 5 repo 동일 script 존재 |
| G3 | 자동 번역 파이프라인 (`scripts/i18n-translate.sh`) 신규 구축 + glossary forced injection + native review gate | "신규 EN 문서 추가 시 `make i18n-translate` 명령 1번으로 3개 언어 초안 생성" | `test -x scripts/i18n-translate.sh` + `make -n i18n-translate` exit 0 |
| G4 | 5 repo README + BRANDING + family + 운영 문서 4-lang 보장 (우선순위 1 tier) | "GitHub 첫 페이지 lang switcher 4개 모두 본문 존재 (placeholder 아님)" | 5 repo × 4 lang × {README, BRANDING, family} = 60 파일 모두 존재 + LOC ≥ 50 |
| G5 | `docs/i18n/README.md` 신규 — family 단위 i18n 정책 + SOP + glossary 사용법 | "i18n 담당자 onboarding 문서 SSOT" | `test -f docs/i18n/README.md` + LOC ≥ 100 |
| G6 | lefthook 에 `readme-i18n-sync` hook 통합 (5 repo) | "drift 발생 시 commit 자체 차단" | 5 repo `lefthook.yml` 모두 `readme-i18n-sync` 명령 존재 |
| G7 | 우선순위 매트릭스 (P0/P1/P2) 정의 + 각 tier 별 cut-off | "어느 문서를 먼저 번역하는가 우선순위 명확" | 본 spec §4.2 의 P0/P1/P2 표 |
| G8 | 5 sub-spec (S4-A~S4-E) 생성 트리거 정의 — 본 spec 머지가 sub-spec 의 진입 게이트 | "본 마스터 spec 의 Phase 7 완료 = sub-spec 5건 진입 가능 상태" | 본 spec Accepted + portfolio supercycle G3 partial 진척 (4 repo × 1 lang = 4/12 → 5 repo × 3 lang = 15/15) |

### 2.2 Non-Goals (본 마스터 spec 의 OOS)

- ❌ 각 repo 의 *실제 번역 본문 작성* (S4-A~S4-E sub-spec 의 범위)
- ❌ 번역 품질 검수의 *수행* — placeholder `[~]` → `[x]` 승격은 native reviewer 가 별 PR 로 수행 (본 spec 은 SOP 만 정의)
- ❌ 영어 외 4번째 언어 추가 (스페인어 / 프랑스어 등) — 본 spec 은 EN/KO/JA/ZH 4-lang 고정
- ❌ CRD `description` 필드의 i18n — portfolio supercycle G3 의 별 항목. 본 spec 은 *문서* 만.
- ❌ 외부 i18n 호스팅 (Crowdin / Lokalise 등) — 비용 + GHA 의존 위험. *로컬 4계층* 정합 위해 self-hosted script 만.
- ❌ Right-to-Left (Arabic / Hebrew) 지원 — 본 spec OOS, 향후 별 RFC

## 3. 아키텍처 (Architecture)

### 3.1 SSOT 분기 모델

```
operator-commons (i18n SSOT)
  ├─ docs/i18n/glossary-{en,ko,ja,zh}.md     ← 4-lang 용어 사전
  ├─ docs/i18n/README.md                     ← 정책 + SOP + 사용법
  ├─ scripts/check-readme-sync.sh            ← 4-lang drift check (5 repo 동기 배포)
  ├─ scripts/i18n-translate.sh               ← 자동 번역 파이프라인 (glossary forced injection)
  └─ scripts/sync-from-commons.sh            ← 5 repo 의 commons-SSOT 동기화 (RFC-0029 §6.5)
                ↓ sync
  ├─ postgres-operator: scripts/check-readme-sync.sh (synced) + scripts/i18n-translate.sh (synced)
  ├─ mongodb-operator: 동일
  ├─ valkey-operator: 동일 + README.ja.md / README.zh.md 4-lang 골격 신규 작성
  └─ forgewise: 동일 (Python repo — script 호환성 확인 필요)
```

### 3.2 7-Phase 모델

```
Phase 0 사전확인 (5 repo i18n 현황 측정)
  ↓
Phase 1 commons SSOT 정비 (glossary 4-lang 완성 + sync hook 확장 + i18n-translate.sh + docs/i18n/README.md)
  ↓
Phase 2 lefthook readme-i18n-sync hook 추가 (commons)
  ↓
Phase 3 SSOT → 4 repo 배포 (sync-from-commons.sh 패턴)
  ↓
Phase 4 우선순위 P0 번역 — README + BRANDING + family (5 repo × 4-lang = 60 파일)
  ↓
Phase 5 우선순위 P1 번역 — 사용자 운영 문서 (INSTALL / QUICKSTART / Operations / Advanced)
  ↓
Phase 6 우선순위 P2 번역 — 나머지 docs/ (선택적, 비용 결정 의존)
  ↓
Phase 7 native reviewer 승격 SOP 실행 + drift 모니터링 정착
```

각 Phase 별 atomic commit. Phase 4-6 은 sub-spec (S4-A~S4-E) 으로 분리되어 *각 repo 별 PR* 로 머지.

### 3.3 자동 번역 파이프라인 데이터 흐름

```
[ 입력 ]
  EN canonical 문서 (e.g. README.md)
  glossary-{ko,ja,zh}.md (SSOT, forced inject)
  사용자 결정 D1 의 엔진 (DeepL | OpenAI | Claude | Google)
       ↓
[ 처리 ]
  scripts/i18n-translate.sh:
    1. EN 본문 parse (section 단위 split)
    2. glossary 의 *EN → target lang* 매핑 dict 로 강제 치환 (코드 식별자 보호)
    3. 엔진 호출 (section 단위로 prompt 분할, context window 절약)
    4. 결과 merge → target.{ko,ja,zh}.md
    5. drift check (check-readme-sync.sh 의 4-lang 매트릭스)
       ↓
[ 출력 ]
  README.ko.md (`[~]` marker — native reviewer 후속 검수 대상)
  README.ja.md (`[~]`)
  README.zh.md (`[~]`)
  검증 로그 → docs/i18n/translate-log/<date>.md
```

## 4. Phase 상세

### 4.1 Phase 0 — 사전확인 (read-only)

**목적**: 본 spec 작성 시점의 5 repo i18n 현황이 implementation 시점에도 유효한지 검증.

**Run** (각 repo cwd 진입 후):
```bash
# operator-commons
find docs -maxdepth 5 -type f -name '*.md' | grep -v -E '/kb/(adr|rfc|deps)/' | grep -v -E '/(internal|superpowers|specs|plans)/' | wc -l  # 6 기대
ls README*.md | wc -l  # 4 기대

# 4 repo 동일 측정 (postgres / mongodb / valkey / forgewise) — §1.1 표 값과 일치 확인

# glossary 완성도
wc -l docs/i18n/glossary-{ko,ja,zh}.md  # ko=203 / ja=102 / zh=102
grep -c '\[~\]' docs/i18n/glossary-{ja,zh}.md  # ja=N / zh=N (placeholder count)
```

**Gate**: 측정 수치가 §1.1 표와 ±10% 이상 차이 시 본 spec 의 표 갱신 후 재진입.

### 4.2 Phase 1 — commons SSOT 정비

**4.2.1 glossary 4-lang 완성**

- `docs/i18n/glossary-ko.md` (203 LOC) — 본문 완성 상태 유지 (변경 없음)
- `docs/i18n/glossary-ja.md` (102 LOC, placeholder) → 본문 작성 (목표 ≥ 200 LOC, ko 와 동일 구조)
- `docs/i18n/glossary-zh.md` (102 LOC, placeholder) → 본문 작성 (목표 ≥ 200 LOC)
- placeholder marker `[~]` 유지 (RFC-0025 §1.2) — native reviewer 검수 후 별 PR 로 `[x]` 승격

**근거**: glossary 가 본문 부재면 자동 번역 파이프라인의 *forced injection* 이 무의미. ja/zh native reviewer 부재 상황에서도 *기계 번역 초안* 으로 본문 채우고 `[~]` 마킹.

**4.2.2 sync hook 4-lang 확장**

기존 `scripts/check-readme-sync.sh` (EN ↔ KO 한정) → 4-lang 매트릭스로 확장.

매트릭스:
- EN ↔ KO: section header 1:1, line diff ≤ 15%, 양방향 cross-link (기존 유지)
- EN ↔ JA: section header 1:1, line diff ≤ **25%** (한자 압축 + 한국어 대비 짧음), 양방향 cross-link
- EN ↔ ZH: section header 1:1, line diff ≤ **30%** (한자 더욱 짧음), 양방향 cross-link
- KO ↔ JA / KO ↔ ZH / JA ↔ ZH: 검사 안 함 (모두 EN 기준)

D4 결정 필요: 위 임계값 (KO=15%, JA=25%, ZH=30%) 이 *적절한* 값인지 — §7 의 D4 미해결.

**4.2.3 자동 번역 파이프라인 신규**

`scripts/i18n-translate.sh` — 의무 기능:
1. CLI 인터페이스: `./scripts/i18n-translate.sh <source.md> [--lang ko|ja|zh|all] [--engine deepl|openai|claude|google]`
2. glossary forced injection (코드 식별자 보호 + 표준 용어 일관성)
3. 엔진 호출 (D1 결정 사항 의존 — §7 D1 미해결)
4. 결과 검증: line count 임계값 / section header 매핑 / cross-link 자동 삽입
5. 출력: `<source>.<lang>.md` 파일 + `[~]` marker + translate-log 기록
6. Makefile target: `make i18n-translate-readme` / `make i18n-translate-all`

**4.2.4 `docs/i18n/README.md` 신규**

family 단위 i18n 정책 + SOP + glossary 사용법. 골격:

```markdown
# i18n — keiailab operator family 다국어 정책 SSOT

## 1. 정책
- 영어 canonical + 한국어/일본어/중국어 추가
- 번역 대상: 사용자/외부 가시 문서 (README, BRANDING, docs/{getting-started,operations,advanced,comparison,developers}, family.md, 운영 governance)
- 번역 비대상: ADR, RFC, deps log, internal/, superpowers/, specs/, plans/, legal (LICENSE/NOTICE)

## 2. 용어 사전 (Glossary)
- SSOT: operator-commons/docs/i18n/glossary-{ko,ja,zh}.md
- 우선순위: glossary 정의 > 일반 사전. 코드 식별자는 영문 그대로.

## 3. 자동 번역 SOP
- 명령: make i18n-translate-all
- 엔진: <D1 결정>
- 결과: `[~]` placeholder marker — native reviewer 검수 전 까지

## 4. Native Review SOP
- 검수자: <D3 결정>
- 승격 기준: 의역/직역 vs 자연스러움 vs 용어 일관성 3 axis 검증
- `[~]` → `[x]` 승격 PR 양식: `docs(i18n): <lang> <doc> native reviewer 승격`

## 5. Drift Control
- lefthook readme-i18n-sync hook (commit 차단)
- check-readme-sync.sh 의 4-lang 매트릭스
- 임계값: EN↔KO ≤ 15%, EN↔JA ≤ 25%, EN↔ZH ≤ 30%

## 6. 5 repo 동기
- SSOT 변경 시 scripts/sync-from-commons.sh 실행 (RFC-0029 §6.5)
- drift seal 정합

## 7. 비용 추산
- <D1 엔진 결정 의존> — 약 167 × 3 × 200 LOC = 100,000 LOC
- 예상 비용: DeepL ($X), Claude ($Y), OpenAI ($Z), Google ($W)
```

### 4.3 Phase 2 — lefthook readme-i18n-sync hook 통합

operator-commons 의 `lefthook.yml` 에 pre-push 또는 pre-commit hook 추가:

```yaml
    readme-i18n-sync:
      glob: "README*.md"
      run: |
        bash scripts/check-readme-sync.sh
```

위치: pre-commit (변경 즉시 차단) 또는 pre-push (push 직전 검증). 본 spec 권장: **pre-push** — pre-commit 은 작업 흐름을 너무 자주 끊음.

5 repo 모두 동일 hook 추가 (Phase 3 sync 단계에서 일괄 배포).

### 4.4 Phase 3 — SSOT → 4 repo 배포

`scripts/sync-from-commons.sh` 패턴 (RFC-0029 §6.5 sub-repo drift seal):

```bash
# operator-commons → 4 repo 동기 (각 repo 의 scripts/ 에 cp)
for repo in postgres-operator mongodb-operator valkey-operator forgewise; do
  target="/Users/phil/Workspace/keiailab/$repo"
  cp scripts/check-readme-sync.sh "$target/scripts/"
  cp scripts/i18n-translate.sh "$target/scripts/"
  # docs/i18n/README.md 도 mirror (link 또는 copy)
  mkdir -p "$target/docs/i18n"
  cp docs/i18n/README.md "$target/docs/i18n/"
  cp docs/i18n/glossary-{ko,ja,zh}.md "$target/docs/i18n/"
done
```

각 repo 별 별 PR (5건). atomic commit message: `chore(i18n): SSOT sync from operator-commons (glossary + scripts + policy)`.

**valkey-operator 특별 처리**: §1.1 의 README 4-lang 골격 부재 → Phase 3 의 valkey PR 에서 `README.ja.md` + `README.zh.md` placeholder 동시 추가 (헤더 lang switcher + 본문 placeholder).

### 4.5 Phase 4 — P0 번역 (README + BRANDING + family)

우선순위 P0 (최우선):
- README.md → README.{ko,ja,zh}.md (5 repo × 4 lang = **20 파일** — 일부 기존)
- BRANDING.md → BRANDING.{ko,ja,zh}.md (5 repo × 4 = 20)
- docs/family.md → docs/family.{ko,ja,zh}.md (5 repo × 4 = 20)

총 **60 파일**. 모두 기계 번역 + glossary injection + `[~]` marker. 각 repo 별 sub-spec (S4-A~S4-E) 으로 분리.

검증 (per repo):
```bash
for lang in ko ja zh; do
  test -f "README.${lang}.md"
  test -f "BRANDING.${lang}.md"
  test -f "docs/family.${lang}.md"
done
```

### 4.6 Phase 5 — P1 번역 (사용자 운영 문서)

우선순위 P1:
- `docs/getting-started.md` (mongodb, postgres) → 3 lang
- `docs/UPGRADING.md` (mongodb) → 3 lang
- `docs/index.md` (postgres) → 3 lang
- `docs/operations/*.md` (valkey 11건) → 3 lang per file = 33 신규 (기존 ko 12건 활용)
- `docs/advanced/*.md` (mongodb 5건) → 3 lang per = 15 신규
- `docs/comparison/*.md` (mongodb 3건) → 3 lang per = 9 신규
- `docs/developers/*.md` (mongodb 4건) → 3 lang per = 12 신규 (개발자 onboarding 도 사용자 가시로 분류)

추산: 약 30개 × 3 lang = **90 파일**.

### 4.7 Phase 6 — P2 번역 (나머지 docs/)

남은 모든 번역 대상. 본 Phase 는 *선택적* — D2 결정 (P0+P1 만 vs 전체) 에 의존.

추산: 167 - 60(P0) - 30(P1) = **77 파일 × 3 lang = 231 파일**

### 4.8 Phase 7 — native reviewer 승격 SOP 실행 + drift 모니터링

- `[~]` → `[x]` 승격 PR (검수자 별로 분리)
- drift 발생 시 alert (lefthook hook 의 exit 1 + Makefile target `make i18n-check-all`)
- 분기별 audit: `find . -name '*.md' -newer <last-audit-date>` → 신규 EN canonical 추가 시 자동 번역 트리거

## 5. 리스크 + 완화 (Risks)

| ID | 리스크 | 영향 | 완화 |
|---|---|---|---|
| R1 | 자동 번역 품질 — glossary 미적용 시 용어 부정확 (예: `Reconciler` → 잘못된 번역) | 신뢰도 손실 | glossary forced injection (코드 식별자 + 표준 K8s 용어 dict 강제) + `[~]` marker 영구 유지 (native review 전 까지) |
| R2 | 번역 비용 — 167 × 3 × ~200 LOC = ~100,000 LOC. API 비용 (DeepL ~$200, Claude ~$50-100, OpenAI ~$30-80, Google ~$50) | 예산 초과 | D1 결정 + Phase 단계 분할 (P0 만 우선 → 비용 모니터링 → P1/P2 진행 여부 결정) |
| R3 | Drift — EN canonical 갱신 후 ko/ja/zh sync 지연 | 잘못된 정보 제공 | (a) lefthook hook 로 commit 자체 차단 (b) auto-PR 생성 hook (post-merge): 신규 EN 변경 감지 시 자동 번역 PR 생성 |
| R4 | Native reviewer 부재 — ja/zh 모국어 검토자 미확보 | `[~]` 영구 잔존 | D3 결정 — 외주 / 커뮤니티 모집 / 기다림 / placeholder 영구 유지 4안 중 선택 |
| R5 | 5 repo sync 정합 — SSOT 변경 시 동기화 누락 | drift 누적 | sync-from-commons.sh + post-merge lefthook hook (commons 의 변경이 4 repo 로 자동 PR) + RFC-0029 §6.5 drift seal 정합 |
| R6 | valkey 의 기존 12개 `docs/operations/*.ko.md` 의 source language 충돌 — ko 원본인가 EN canonical 부재인가 | 번역 base 모호 | D5 결정 필요. *권장*: EN canonical 강제 (현 ko 파일들 → EN 역번역 후 4-lang 재생성) |
| R7 | 자동 번역 엔진 (D1) 미결정 — 본 spec 의 Phase 1.3 진입 차단 | implementation 시작 불가 | 사용자 결정 필수. 본 spec Accepted 의 *전제 조건* 으로 명시. |
| R8 | forgewise 는 Python repo + 디렉토리 구조 다름 → script 호환성 | sync 실패 | Phase 3 의 forgewise PR 에서 script path adjust (bash script 의 hard-coded 경로 회피 + repo root 자동 감지) |
| R9 | 자동 번역 결과의 PR 폭주 (5 repo × 167 파일 = 835 PR 가능성) | review 부담 | atomic 정책 완화 — i18n 한정으로 *batch PR* (per Phase × per repo) 허용. 단 본 sub-spec 결정 필요 |

## 6. 성공 조건 (Success Criteria)

본 마스터 spec 의 G1~G8 달성 시 다음 8건 모두 PASS:

```bash
# SC1: glossary 4-lang 본문 완성
wc -l /Users/phil/Workspace/keiailab/operator-commons/docs/i18n/glossary-ja.md | awk '$1 >= 150'
wc -l /Users/phil/Workspace/keiailab/operator-commons/docs/i18n/glossary-zh.md | awk '$1 >= 150'

# SC2: sync hook 4-lang 매트릭스
grep -q "README.ja.md" /Users/phil/Workspace/keiailab/operator-commons/scripts/check-readme-sync.sh
grep -q "README.zh.md" /Users/phil/Workspace/keiailab/operator-commons/scripts/check-readme-sync.sh

# SC3: 자동 번역 파이프라인 신규
test -x /Users/phil/Workspace/keiailab/operator-commons/scripts/i18n-translate.sh
make -C /Users/phil/Workspace/keiailab/operator-commons -n i18n-translate-readme

# SC4: i18n 정책 SSOT 문서
test -f /Users/phil/Workspace/keiailab/operator-commons/docs/i18n/README.md
test "$(wc -l < /Users/phil/Workspace/keiailab/operator-commons/docs/i18n/README.md)" -ge 100

# SC5: lefthook hook 통합
grep -q "readme-i18n-sync" /Users/phil/Workspace/keiailab/operator-commons/lefthook.yml

# SC6: 5 repo 동기 — sync script 존재 + 실행 가능
test -x /Users/phil/Workspace/keiailab/operator-commons/scripts/sync-from-commons.sh

# SC7: P0 번역 완료 (5 repo × 4 lang × 3 doc = 60 파일)
total=0
for r in operator-commons postgres-operator mongodb-operator valkey-operator forgewise; do
  for lang in md ko.md ja.md zh.md; do
    for doc in README BRANDING; do
      test -f "/Users/phil/Workspace/keiailab/$r/$doc.$lang" && total=$((total+1))
    done
    test -f "/Users/phil/Workspace/keiailab/$r/docs/family.$lang" && total=$((total+1))
  done
done
test "$total" -ge 60

# SC8: valkey README 4-lang 골격 완성
test -f /Users/phil/Workspace/keiailab/valkey-operator/README.ja.md
test -f /Users/phil/Workspace/keiailab/valkey-operator/README.zh.md
```

## 7. 사용자 결정 필요 (Open Decisions — 미해결)

본 마스터 spec 의 *implementation 진입 전* 사용자 결정 필수:

### D1 자동 번역 엔진 선택

| 옵션 | 장점 | 단점 | 추정 비용 (100k LOC) |
|---|---|---|---|
| (a) DeepL Pro API | 번역 품질 최고 (특히 EN→JA/DE), glossary 지원 | EN→ZH 품질 평이, 코드 영역 처리 약함 | ~$200 |
| (b) OpenAI GPT-4 / 4o API | 코드 인식 양호, prompt 자유도 | EN→KO/JA 품질 평이, 비용 변동 | ~$30-80 |
| (c) Claude API (Sonnet 4.5) | 한국어 품질 최고, 긴 컨텍스트 | EN→JA/ZH 품질 검증 부족 | ~$50-100 |
| (d) Google Cloud Translation | EN→ZH 최고, 빠름 | 격식체/평어 일관성 약함, glossary 제한적 | ~$50 |
| (e) 다중 엔진 (lang 별 best) | 품질 극대화 | 복잡성 증가, sync 어려움 | ~$80-150 |

**권장 (디폴트, 미결정 시)**: (c) Claude — 본 거버넌스 자체가 Claude Code 기반. 한국어 품질 + 컨텍스트 길이 + prompt caching 으로 비용 절감. ja/zh 는 후속 native reviewer 검수 가정.

### D2 번역 우선순위 cut-off

- (a) P0 만 (README + BRANDING + family = 60 파일) — 최소 노출
- (b) P0 + P1 (운영 문서 추가 = 150 파일) — *권장*
- (c) P0 + P1 + P2 (전체 = 501 파일) — 완전 i18n. 비용 + 검수 부담

### D3 native reviewer 확보

- (a) 외주 (전문 번역가 ja/zh 각 1인, ~$1000~5000)
- (b) 커뮤니티 모집 (CONTRIBUTING.md 에 i18n reviewer welcome 명시)
- (c) 기다림 (사용자 증가 후 자연 발생 유도)
- (d) placeholder `[~]` 영구 유지 (기계 번역 = 베스트 에포트)

### D4 drift 임계값 — 한자 압축 고려

- 현 EN↔KO: line diff ≤ 15%
- 본 spec 제안: EN↔JA ≤ 25%, EN↔ZH ≤ 30%
- 사용자 결정 필요: 위 값 (25/30) 적절 or 측정 후 조정 or 다른 metric (글자 수 / 토큰 수) 사용?

### D5 valkey 의 기존 13개 `*.ko.md` 처리

- (a) ko 본 source 강제 — 13 파일을 ko canonical 로 인정 → EN/ja/zh 역번역. ADR 추가 정당화 필요.
- (b) EN canonical 강제 — *권장* — 13 ko 파일을 source 로 EN 역번역 후 EN 을 canonical 로 승격, ja/zh 는 EN 에서 신규 번역. ko 파일은 유지 (단 EN ↔ KO drift check 진입).
- (c) 13 파일 deprecated 표기 후 EN-only 로 시작 — 기존 ko 가치 손실. 권장 안 함.

## 8. 범위 외 (Out-of-Scope)

본 S4 *마스터* spec 은 *정책 + SSOT 자산 + 파이프라인 설계* 만 다룬다. 각 저장소의 *실제 번역 실행* 은 별 sub-spec:

- **S4-A** `postgres-operator` 의 30 docs + 13 top-level × 3 lang 번역 실행
- **S4-B** `mongodb-operator` 의 32 docs + 17 top-level × 3 lang 번역 실행
- **S4-C** `valkey-operator` 의 33 docs + 18 top-level × 3 lang 번역 실행 + README 4-lang 골격
- **S4-D** `operator-commons` 의 6 docs + 12 top-level × 3 lang 번역 실행
- **S4-E** `forgewise` 의 5 docs + 1 top-level × 3 lang 번역 실행

각 sub-spec 은 본 마스터 spec 의 §3 SSOT + §4 Phase 4-6 정책을 따른다.

다음은 본 spec 자체의 OOS (별 cycle):

- ❌ CRD `description` 필드 i18n (Kubernetes API 레벨 — 별 RFC)
- ❌ 외부 호스팅 (docusaurus / docsify / readthedocs i18n plugin) — 본 spec 은 markdown native 만
- ❌ Right-to-Left 언어 (Arabic / Hebrew) — 별 RFC
- ❌ 4번째 추가 언어 (es / fr / de / pt) — 본 spec EN/KO/JA/ZH 4-lang 고정
- ❌ Helm chart values 의 i18n (e.g. validation error message) — `pkg/events` 의 reason constants 와 동일하게 영문 강제
- ❌ Go 코드 에러 메시지 i18n — *영문 유지* (operator 운영 환경 standard)

## 9. 변경 이력

| 날짜 | 변경 | 상태 |
|---|---|---|
| 2026-05-21 | 초안 작성 (portfolio supercycle S4 마스터 spec 분리) | Proposed — 5건 사용자 결정 미해결 |
