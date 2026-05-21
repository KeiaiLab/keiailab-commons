# i18n — keiailab operator family 다국어 정책 SSOT

> 본 문서는 keiailab operator family (operator-commons + 3 operators + forgewise = 5 repo) 의 i18n 통합 정책 SSOT 입니다. 모든 5 repo 가 본 정책을 *그대로* 따릅니다 (RFC-0029 §6.5 sync drift seal 정합).
>
> **상태**: S4 Phase 1 — 정책 + SSOT 자산 정비 완료. 실제 번역 실행은 별 sub-cycle (S4-A~S4-E + S4 Phase 4-6).

## §1 정책

### 1.1 기본 원칙

- **canonical 언어**: 영어 (모든 source 의 ground-truth)
- **번역 대상 언어 (3개 추가)**: 한국어 (ko) / 日本語 (ja) / 中文 (zh, 简体)
- **4-lang 골격**: README.md (EN) + README.ko.md + README.ja.md + README.zh.md
- **상위 4 언어 외 추가 불가** — 별 RFC 후에만 (es/fr/de 등 OOS)

### 1.2 번역 대상 vs 비대상

**번역 대상** (사용자/외부 가시 문서):
- `README.{md,ko,ja,zh}.md` (5 repo)
- top-level governance/branding: BRANDING / GOVERNANCE / SECURITY / CONTRIBUTING / CODE_OF_CONDUCT / MAINTAINERS / ADOPTERS / ROADMAP / ARCHITECTURE / CHANGELOG / STABILITY
- `docs/family.md` (cross-link footer)
- `docs/getting-started.md`, `docs/UPGRADING.md`, `docs/index.md`
- `docs/advanced/*.md`, `docs/comparison/*.md`, `docs/developers/*.md`
- `docs/operations/*.md`
- `docs/i18n/glossary-{ko,ja,zh}.md` (본 SSOT)

**번역 비대상** (내부/governance/legal):
- `docs/kb/adr/*.md` (ADR — 결정 추적, 영문 canonical only)
- `docs/kb/rfc/*.md` (RFC — 동일)
- `docs/kb/deps/*.md` (deps log — 자동 생성)
- `docs/internal/*.md` (HANDOFF / README / TASKS — 내부 운영자)
- `docs/superpowers/*` (cycle artifacts — 일시적)
- `docs/specs/*.md` (design specs)
- `docs/plans/*.md` (implementation plans)
- `LICENSE`, `NOTICE`, `CITATION.cff`, `.gitignore` 등 (legal/config)

## §2 용어 사전 (Glossary)

### 2.1 위치

SSOT: `operator-commons/docs/i18n/glossary-{ko,ja,zh}.md` (본 디렉토리)

### 2.2 우선순위

1. glossary 의 정의 (최우선)
2. Kubernetes 공식 한국어/일본어/중문 가이드
3. 일반 사전
4. **코드 식별자는 영문 그대로** (절대 번역 금지) — `pkg/probes`, `Reconciler`, `ValkeyCluster`, `kubectl` 등

### 2.3 일관성 규칙 (4-lang 공통)

- 격식체 / 평어 혼용 금지 (한 문서 내 일관)
- 외부 사용자 가시 문서 = 격식체
- 내부 문서 (HANDOFF/AGENTS) = 평어 또는 자유
- 첫 등장 시 영문 원어 + 괄호 번역, 이후 번역 단독 가능

## §3 자동 번역 SOP

### 3.1 엔진 결정 (D1, 2026-05-21)

**선택**: Claude direct (현 cycle 시점에서는 subagent 가 Claude 모델로서 직접 번역).

근거:
- 본 거버넌스 자체가 Claude Code 기반
- 한국어 품질 최고
- 컨텍스트 길이 + prompt caching 으로 비용 절감
- ja/zh 는 후속 native reviewer 검수 가정

### 3.2 명령

```bash
# 1 파일 번역 (수동 — 현 시점 D1)
# subagent 가 source 읽고 직접 번역 → 출력 파일 작성 + warning 배너 강제 삽입

# 자동화 (향후 sub-cycle 구현)
./scripts/i18n-translate.sh <source.md> --lang all --engine claude

# dry-run
./scripts/i18n-translate.sh README.md --dry-run
```

### 3.3 결과 marker

모든 자동 번역 파일은 *반드시* 다음 warning 배너 삽입:

```markdown
> ⚠️ This translation is AI-generated and pending native review.
```

추가 마킹: `[검토 필요]` (검수 필요), `[검토 완료]` (native reviewer 검수 후 PR 로 승격).

## §4 Native Review SOP

### 4.1 검수자 확보 (D3, 2026-05-21)

**선택**: Claude 자동 번역 + `[검토 필요]` 마킹 + warning 배너 (placeholder 영구 유지 — 사용자 / 커뮤니티 검수 유도).

근거:
- 외주 비용 회피
- 커뮤니티 모집은 별 wave (CONTRIBUTING.md i18n reviewer welcome 추가)
- 기다림 동안 자동 번역 = 최소 노출 보장

### 4.2 승격 기준

`[검토 필요]` → `[검토 완료]` 승격 시 native reviewer 가 확인해야 할 사항:

1. **의역 vs 직역 균형** — 자연스러움 우선, 단 기술 용어 정확성 보존
2. **용어 일관성** — glossary 적용 확인
3. **격식체 / 평어 일관** — 한 문서 내 통일
4. **링크 정합** — cross-link 의 lang 별 정확
5. **코드 식별자 보존** — 절대 번역 안 됨 확인

### 4.3 승격 PR 양식

```
title: docs(i18n): <lang> <doc> native reviewer 승격
body:
  ## 검수자
  - <name / handle>

  ## 검수 범위
  - <file path>

  ## 변경 사항
  - [검토 필요] → [검토 완료] marker 변경
  - warning 배너 제거 (또는 partial 유지)
  - <기타 의역 수정>

  ## 검증
  - [x] glossary 일관성
  - [x] 격식체 / 평어 통일
  - [x] cross-link 정합
```

## §5 Drift Control

### 5.1 lefthook hook (S4 Phase 2)

S4 Phase 2 에서 `lefthook.yml` 의 pre-push 에 `readme-i18n-sync` hook 추가:

```yaml
    readme-i18n-sync:
      glob: "README*.md"
      run: bash scripts/check-readme-sync.sh
```

`pre-push` 위치 — `pre-commit` 보다 작업 흐름 마찰 적음.

### 5.2 임계값 (D4, 2026-05-21)

| lang | line diff 임계값 | 근거 |
|---|---|---|
| ko | ≤ 15% | 한국어 LOC 는 EN 와 유사 (기존 EN↔KO 정합 유지) |
| ja | ≤ 25% | 일본어는 가나 + 한자 혼용 — EN 대비 ~15-20% 짧음 (여유 ±5%) |
| zh | ≤ 30% | 중문은 한자 100% — EN 대비 ~25% 짧음 (여유 ±5%) |

상위 값은 spec §4.2.2 의 권장값. 실제 측정 후 조정 가능.

### 5.3 4-lang 매트릭스

`scripts/check-readme-sync.sh` 가 4-lang 매트릭스 (EN ↔ {ko,ja,zh}) 자동 검사:
- target lang file 부재 시 skip (4-lang 골격 미완료 repo 대응)
- per-lang 우회 가능 (`SKIP_CHECK_README_SYNC_JA=1` 등)

## §6 5 repo 동기

### 6.1 SSOT → 4 repo 배포 패턴

`scripts/sync-from-commons.sh` (S4 Phase 3 신규):

```bash
for repo in postgres-operator mongodb-operator valkey-operator forgewise; do
  target="/Users/phil/Workspace/keiailab/$repo"
  cp scripts/check-readme-sync.sh "$target/scripts/"
  cp scripts/i18n-translate.sh "$target/scripts/"
  mkdir -p "$target/docs/i18n"
  cp docs/i18n/README.md "$target/docs/i18n/"
  cp docs/i18n/glossary-{ko,ja,zh}.md "$target/docs/i18n/"
done
```

RFC-0029 §6.5 sub-repo drift seal 정합 — SSOT 본문 그대로 cp.

### 6.2 valkey-operator 특별 처리

valkey-operator 는 README 4-lang 골격 부재 (en+ko 만). S4 Phase 3 의 valkey PR 에서 `README.ja.md` + `README.zh.md` placeholder 동시 추가.

또한 valkey 의 기존 13개 `docs/operations/*.ko.md` 처리 (D5 결정):
- **선택**: (b) EN canonical 강제 — 13 ko 파일을 source 로 EN 역번역 후 EN 을 canonical 로 승격
- 별 cycle (S4-C 의 sub-task) 으로 분리

## §7 우선순위 매트릭스 (D2, 2026-05-21)

| Tier | 범위 | 파일 수 추산 | 결정 |
|---|---|---|---|
| **P0** | README + BRANDING + family (5 repo × 3 doc × 3 lang) | 60 신규 + 일부 기존 갱신 | ✅ 본 cycle 범위 |
| **P1** | 운영 문서 (getting-started, UPGRADING, operations, advanced, comparison, developers) | 약 30 × 3 lang = 90 | ✅ 본 cycle 범위 |
| **P2** | 나머지 docs/ | 약 77 × 3 lang = 231 | ✅ 본 cycle 범위 (D2 결정: 전체) |

D2 결정: **(c) P0 + P1 + P2 (전체)** — 완전 i18n.

## §8 비용 추산

- 총 번역 LOC: 167 × 3 × ~200 = 약 100,000 LOC
- 엔진 (D1): Claude direct
- **현 cycle 비용**: 0 (subagent 가 직접 번역)
- 향후 API 호출 자동화 시 추산: Claude Sonnet 4.5 input $3 / 1M token + output $15 / 1M token → 약 $50-100

## §9 참조

- spec: `docs/specs/2026-05-21-i18n-4lang-master-design.md`
- supercycle spec: `docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md` §4.4 S4
- 본 cycle: `feat/i18n-ssot-baseline-2026-05-21`
- glossary 3종: `glossary-ko.md` / `glossary-ja.md` / `glossary-zh.md`
- check script: `../../scripts/check-readme-sync.sh`
- translate script: `../../scripts/i18n-translate.sh`

## §10 변경 이력

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 신설 — i18n 정책 SSOT v0.1 (S4 Phase 1) | docs/specs/2026-05-21-i18n-4lang-master-design.md §4.2.4 |
