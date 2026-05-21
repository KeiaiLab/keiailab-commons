# i18n — operator-commons 다국어 정책

본 문서는 `operator-commons` 다국어 문서 운영 정책을 정의합니다. English
canonical + 한국어 / 日本語 / 中文 의 4-lang 골격을 유지합니다.

## §1 정책

### 1.1 기본 원칙

- **canonical 언어**: 영어 (모든 source 의 ground-truth)
- **번역 대상 언어 (3개 추가)**: 한국어 (`ko`), 日本語 (`ja`), 中文 (`zh`, 简体)
- **4-lang 골격**: `README.md` (EN) + `README.ko.md` + `README.ja.md` + `README.zh.md`
- **상위 4 언어 외 추가 불가** — 별 ADR 후에만 (es / fr / de 등 OOS)

### 1.2 번역 대상 vs 비대상

**번역 대상** (사용자 / 외부 가시 문서):

- `README.{md,ko,ja,zh}.md`
- 거버넌스 / 브랜딩 캐논: `docs/BRANDING.md`, `docs/STABILITY.md`,
  `docs/coverage-report.md` 등의 4-lang 매트릭스
- `docs/i18n/glossary-{ko,ja,zh}.md` (본 SSOT)

**번역 비대상** (내부 / 거버넌스 / 법적):

- `docs/kb/adr/*.md` (결정 추적 — 영문 canonical only)
- `docs/kb/deps/*.md` (의존성 audit — 자동 생성)
- `LICENSE`, `NOTICE`, `CITATION.cff`, `.gitignore` (legal / config)
- `MAINTAINERS.md`, `GOVERNANCE.md`, `CHANGELOG.md` (한국어 canonical 정합)

## §2 용어 사전 (Glossary)

### 2.1 위치

SSOT: [`glossary-ko.md`](glossary-ko.md) / [`glossary-ja.md`](glossary-ja.md)
/ [`glossary-zh.md`](glossary-zh.md).

### 2.2 우선순위

1. glossary 의 정의 (최우선).
2. Kubernetes 공식 한국어 / 일본어 / 중문 가이드.
3. 일반 사전.
4. **코드 식별자는 영문 그대로** (절대 번역 금지) — `pkg/probes`,
   `Reconciler`, `kubectl` 등.

### 2.3 일관성 규칙 (4-lang 공통)

- 격식체 / 평어 혼용 금지 (한 문서 내 일관).
- 외부 사용자 가시 문서 = 격식체.
- 내부 문서 (AGENTS) = 평어 또는 자유.
- 첫 등장 시 영문 원어 + 괄호 번역, 이후 번역 단독 가능.

## §3 자동 번역 SOP

### 3.1 엔진

**선택**: Claude direct (AI subagent 가 source 를 읽고 직접 번역).

근거: 한국어 품질, prompt caching 으로 비용 절감. `ja` / `zh` 는 후속 native
reviewer 검수를 가정합니다.

### 3.2 명령

```bash
# 1 파일 번역 (수동)
# subagent 가 source 읽고 직접 번역 → 출력 파일 + warning 배너 강제 삽입

# 자동화 (스크립트)
./scripts/i18n-translate.sh <source.md> --lang all --engine claude

# dry-run
./scripts/i18n-translate.sh README.md --dry-run
```

### 3.3 결과 marker

모든 자동 번역 파일은 다음 warning 배너 삽입:

```markdown
> ⚠️ This translation is AI-generated and pending native review.
```

추가 마킹: `[검토 필요]` (검수 필요), `[검토 완료]` (native reviewer 검수 후
PR 로 승격).

## §4 Native Review SOP

### 4.1 검수자 확보

기본 정책: Claude 자동 번역 + `[검토 필요]` 마킹 + warning 배너. native
reviewer 가 확보될 때까지 placeholder 유지.

### 4.2 승격 기준

`[검토 필요]` → `[검토 완료]` 승격 시 native reviewer 가 확인해야 할 사항:

1. **의역 vs 직역 균형** — 자연스러움 우선, 단 기술 용어 정확성 보존.
2. **용어 일관성** — glossary 적용 확인.
3. **격식체 / 평어 일관** — 한 문서 내 통일.
4. **링크 정합** — cross-link 의 lang 별 정확.
5. **코드 식별자 보존** — 절대 번역 안 됨 확인.

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

### 5.1 lefthook hook

`lefthook.yml` 의 pre-push 에 `readme-i18n-sync` hook 이 강제됩니다:

```yaml
readme-i18n-sync:
  glob: "README*.md"
  run: bash scripts/check-readme-sync.sh
```

`pre-push` 위치 — `pre-commit` 보다 작업 흐름 마찰이 적습니다.

### 5.2 임계값

| lang | line diff 임계값 | 근거 |
|---|---|---|
| ko | ≤ 15 % | 한국어 LOC 는 EN 와 유사 |
| ja | ≤ 25 % | 일본어는 가나 + 한자 혼용 — EN 대비 ~15-20 % 짧음 |
| zh | ≤ 30 % | 중문은 한자 100 % — EN 대비 ~25 % 짧음 |

상위 값은 권장값이며, 실제 측정 후 조정 가능합니다.

### 5.3 4-lang 매트릭스

`scripts/check-readme-sync.sh` 가 4-lang 매트릭스 (EN ↔ {ko, ja, zh}) 자동
검사:

- target lang file 부재 시 skip (4-lang 골격 미완료 시 대응).
- per-lang 우회 가능 (`SKIP_CHECK_README_SYNC_JA=1` 등).

## §6 참조

- glossary 3종: [`glossary-ko.md`](glossary-ko.md) /
  [`glossary-ja.md`](glossary-ja.md) / [`glossary-zh.md`](glossary-zh.md).
- check script: [`../../scripts/check-readme-sync.sh`](../../scripts/check-readme-sync.sh).
- translate script: [`../../scripts/i18n-translate.sh`](../../scripts/i18n-translate.sh).
