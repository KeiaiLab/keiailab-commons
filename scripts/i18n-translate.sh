#!/usr/bin/env bash
# scripts/i18n-translate.sh
# Claude API 추상화 wrapper — 자동 번역 파이프라인 (skeleton, S4 Phase 1).
#
# S4 spec (docs/specs/2026-05-21-i18n-4lang-master-design.md §4.2.3) 명세:
#   1. CLI: ./scripts/i18n-translate.sh <source.md> [--lang ko|ja|zh|all] [--engine deepl|openai|claude|google]
#   2. glossary forced injection (코드 식별자 보호 + 표준 용어 일관성)
#   3. 엔진 호출 (D1 결정: Claude direct — Sonnet 4.5)
#   4. 결과 검증 (line count / section header / cross-link)
#   5. 출력: <source>.<lang>.md + `[~]` marker + translate-log
#
# 본 cycle (S4 Phase 1) 에서는 *문서 + 자동화 골격* 만 작성.
# 실제 API 호출 구현은 별 sub-cycle (S4 Phase 4-6 의 batch 번역 시).
#
# 사용자 결정 D1: Claude direct (현재 세션에서 subagent 가 직접 번역) — API 자동 호출 *옵션*.
# 본 cycle 에서는 subagent 가 Claude 모델 그 자체로서 translate 를 수행하므로,
# 본 script 는 *향후 자동화* 시점의 골격으로만 작성.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 기본값
SOURCE_FILE=""
TARGET_LANGS="all"
ENGINE="claude"
DRY_RUN=0

# 사용법
usage() {
  cat <<EOF
Usage: $0 <source.md> [options]

Arguments:
  <source.md>          번역할 EN canonical 원본 markdown 파일

Options:
  --lang <ko|ja|zh|all>   번역 대상 언어 (default: all)
  --engine <name>         번역 엔진 (default: claude)
                          지원 예정: claude (D1 결정), deepl, openai, google
  --dry-run               실행 안 함, 계획만 출력
  -h, --help              도움말

Example:
  $0 README.md --lang all --engine claude
  $0 docs/getting-started.md --lang ja --dry-run

비고:
  - 본 script 는 S4 Phase 1 시점의 *skeleton*. 실제 API 호출은 별 sub-cycle 구현.
  - 현 cycle 시점에는 subagent (Claude) 가 직접 번역 (D1 결정 사항).
  - 생성된 파일에는 \`[검토 필요]\` marker + warning 배너 강제 삽입.
EOF
}

# 인자 파싱
while [ $# -gt 0 ]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --lang)
      TARGET_LANGS="$2"
      shift 2
      ;;
    --engine)
      ENGINE="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    *)
      if [ -z "$SOURCE_FILE" ]; then
        SOURCE_FILE="$1"
      else
        echo "ERROR: 알 수 없는 인자: $1" >&2
        usage
        exit 1
      fi
      shift
      ;;
  esac
done

if [ -z "$SOURCE_FILE" ]; then
  echo "ERROR: source.md 인자 필수" >&2
  usage
  exit 1
fi

if [ ! -f "$SOURCE_FILE" ]; then
  echo "ERROR: $SOURCE_FILE 파일 없음" >&2
  exit 1
fi

# 언어 list
case "$TARGET_LANGS" in
  all) LANGS="ko ja zh" ;;
  ko|ja|zh) LANGS="$TARGET_LANGS" ;;
  *)
    echo "ERROR: 알 수 없는 lang: $TARGET_LANGS (지원: ko, ja, zh, all)" >&2
    exit 1
    ;;
esac

# Warning 배너 (모든 자동 번역 파일에 강제 삽입)
WARNING_BANNER='> ⚠️ This translation is AI-generated and pending native review.'

echo "==================================================="
echo "i18n-translate.sh — S4 Phase 1 skeleton"
echo "==================================================="
echo "Source:    $SOURCE_FILE"
echo "Languages: $LANGS"
echo "Engine:    $ENGINE"
echo "Dry-run:   $DRY_RUN"
echo "---"

for lang in $LANGS; do
  # 출력 파일 명: foo.md → foo.<lang>.md
  base="${SOURCE_FILE%.md}"
  out="${base}.${lang}.md"

  echo "[$lang] $SOURCE_FILE → $out"

  if [ "$DRY_RUN" = "1" ]; then
    echo "  (dry-run) — 실행 안 함"
    continue
  fi

  # S4 Phase 1 시점: 실제 API 호출 미구현.
  # 향후 (별 sub-cycle): glossary forced injection + API 호출 + 결과 검증.
  cat <<EOF >&2

  [TODO] $lang 번역 실행 미구현 — S4 Phase 1 skeleton.

  현 시점 권장 절차 (사용자 결정 D1: Claude direct):
    1. Claude (subagent) 가 $SOURCE_FILE 본문을 읽고 직접 $lang 로 번역
    2. docs/i18n/glossary-${lang}.md 의 용어 강제 적용
    3. 출력 $out 파일에 다음 헤더 강제 삽입:

       ${WARNING_BANNER}

    4. 변경 후 ./scripts/check-readme-sync.sh 로 drift 검증
    5. commit + PR

  본 script 의 자동 API 호출은 향후 sub-cycle 에서 구현.

EOF
done

echo "---"
echo "본 script 는 향후 자동화 entry point. 현 시점은 manual translation (D1) 권장."
exit 0
