#!/usr/bin/env bash
# scripts/check-readme-sync.sh
# README.md (EN canonical) ↔ README.ko.md (KR) drift verify.
#
# PR2 (`docs/readme-i18n-ko`, 2026-05-20) — RFC-0045 §2.5 Codex challenge #3
# drift control codify. lefthook pre-commit `readme-i18n-sync` hook sister.
#
# 검사 항목:
#   1. 양 file 존재
#   2. `## ` section header 개수 동일 (1:1 매핑 정합)
#   3. line count diff ≤ 15% (의미 손실 없는 완역 정합)
#   4. 양방향 cross-link 존재 (EN ↔ KR 헤더 link)
#
# exit 0: 정합 ✅
# exit 1: drift ❌ (fix 후 재실행 의무)
#
# 우회: SKIP_CHECK_README_SYNC=1 env (의도적 임시 우회).

set -euo pipefail

if [ "${SKIP_CHECK_README_SYNC:-0}" = "1" ]; then
  echo "[skip] SKIP_CHECK_README_SYNC=1 환경변수로 검사 우회"
  exit 0
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

EN="README.md"
KR="README.ko.md"
THRESHOLD_PCT=15

# 1. 양 file 존재
if [ ! -f "$EN" ]; then
  echo "ERROR: $EN not found in $ROOT"
  exit 1
fi
if [ ! -f "$KR" ]; then
  echo "ERROR: $KR not found in $ROOT"
  exit 1
fi

# 2. section header 개수 정합 (`## ` prefix)
en_sections=$(grep -c '^## ' "$EN" 2>/dev/null || echo 0)
kr_sections=$(grep -c '^## ' "$KR" 2>/dev/null || echo 0)
if [ "$en_sections" != "$kr_sections" ]; then
  echo "ERROR: section header 수 불일치 — $EN=$en_sections vs $KR=$kr_sections"
  echo "  EN 섹션 list:"
  grep '^## ' "$EN" | head -10
  echo "  KR 섹션 list:"
  grep '^## ' "$KR" | head -10
  exit 1
fi

# 3. line count diff ≤ 15%
en_lines=$(wc -l < "$EN" | tr -d ' ')
kr_lines=$(wc -l < "$KR" | tr -d ' ')
if [ "$kr_lines" -gt "$en_lines" ]; then
  diff=$((kr_lines - en_lines))
else
  diff=$((en_lines - kr_lines))
fi
threshold=$((en_lines * THRESHOLD_PCT / 100))
if [ "$diff" -gt "$threshold" ]; then
  echo "ERROR: line count diff > ${THRESHOLD_PCT}% — $EN=$en_lines vs $KR=$kr_lines"
  echo "  diff=$diff, threshold=$threshold (over by $((diff - threshold)))"
  echo "  의미 손실 없는 완역 정합 의무 위반"
  exit 1
fi

# 4. cross-link 양방향 (헤더 또는 본문 어디든)
if ! grep -q 'README\.ko\.md' "$EN"; then
  echo "ERROR: $EN missing link to $KR"
  echo "  fix: $EN line 1 또는 2 에 다음 줄 추가:"
  echo "       > 한국어 (Korean) translation available: [README.ko.md](README.ko.md)"
  exit 1
fi
if ! grep -q 'README\.md' "$KR"; then
  echo "ERROR: $KR missing link to $EN"
  echo "  fix: $KR line 1 에 다음 줄 추가:"
  echo "       > English README: [README.md](README.md) — canonical / 정본"
  exit 1
fi

# 5. 정합 보고
echo "OK: $EN ↔ $KR sync verified"
echo "  sections: $en_sections (양쪽 동일)"
echo "  lines:    EN=$en_lines / KR=$kr_lines (diff=$diff ≤ threshold=$threshold)"
echo "  links:    양방향 cross-link 존재"
exit 0
