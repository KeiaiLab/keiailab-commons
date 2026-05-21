#!/usr/bin/env bash
# scripts/check-readme-sync.sh
# README.md (EN canonical) ↔ README.{ko,ja,zh}.md drift verify (4-lang matrix).
#
# S4 Phase 1 (2026-05-21, docs/specs/2026-05-21-i18n-4lang-master-design.md §4.2.2):
#   기존 EN↔KO 한정 → EN↔{KO,JA,ZH} 매트릭스 확장.
#   언어별 임계값 (한자 압축 반영): KO=15% / JA=25% / ZH=30%.
#
# 원본: PR2 (`docs/readme-i18n-ko`, 2026-05-20) — RFC-0045 §2.5 Codex challenge #3
#       drift control codify. lefthook pre-commit `readme-i18n-sync` hook sister.
#
# 검사 항목 (per target lang):
#   1. 양 file 존재
#   2. `## ` section header 개수 동일 (1:1 매핑 정합)
#   3. line count diff ≤ 임계값% (의미 손실 없는 완역 정합)
#   4. 양방향 cross-link 존재
#
# exit 0: 4-lang 매트릭스 모두 정합 (또는 target lang file 부재 시 skip)
# exit 1: 1건이라도 drift 발생
#
# 우회: SKIP_CHECK_README_SYNC=1 env (의도적 임시 우회).
# 부분 우회: SKIP_CHECK_README_SYNC_JA=1 / SKIP_CHECK_README_SYNC_ZH=1 (lang 별 우회).

set -euo pipefail

if [ "${SKIP_CHECK_README_SYNC:-0}" = "1" ]; then
  echo "[skip] SKIP_CHECK_README_SYNC=1 환경변수로 전체 검사 우회"
  exit 0
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

EN="README.md"

# lang 별 (target_lang, target_file, threshold_pct, skip_env_var)
# 한자 압축 반영 — JA/ZH 는 EN 대비 짧음. spec §4.2.2 참조.
LANG_MATRIX=(
  "ko:README.ko.md:15:SKIP_CHECK_README_SYNC_KO"
  "ja:README.ja.md:25:SKIP_CHECK_README_SYNC_JA"
  "zh:README.zh.md:30:SKIP_CHECK_README_SYNC_ZH"
)

# 1. EN canonical 존재 확인
if [ ! -f "$EN" ]; then
  echo "ERROR: $EN (EN canonical) not found in $ROOT"
  exit 1
fi

en_sections=$(grep -c '^## ' "$EN" 2>/dev/null || echo 0)
en_lines=$(wc -l < "$EN" | tr -d ' ')

overall_pass=1
checked_count=0

for entry in "${LANG_MATRIX[@]}"; do
  IFS=':' read -r lang target_file threshold_pct skip_env <<< "$entry"

  # lang 별 우회
  skip_val=$(eval "echo \${$skip_env:-0}")
  if [ "$skip_val" = "1" ]; then
    echo "[skip $lang] $skip_env=1 환경변수로 우회"
    continue
  fi

  # 2. target file 존재 (없으면 skip — 4-lang 골격 부재 repo 대응)
  if [ ! -f "$target_file" ]; then
    echo "[skip $lang] $target_file 부재 — 4-lang 골격 미완료 (별 PR 대상)"
    continue
  fi

  checked_count=$((checked_count + 1))

  # 3. section header 개수 정합
  tgt_sections=$(grep -c '^## ' "$target_file" 2>/dev/null || echo 0)
  if [ "$en_sections" != "$tgt_sections" ]; then
    echo "ERROR [$lang]: section header 수 불일치 — $EN=$en_sections vs $target_file=$tgt_sections"
    overall_pass=0
    continue
  fi

  # 4. line count diff ≤ 임계값
  tgt_lines=$(wc -l < "$target_file" | tr -d ' ')
  if [ "$tgt_lines" -gt "$en_lines" ]; then
    diff=$((tgt_lines - en_lines))
  else
    diff=$((en_lines - tgt_lines))
  fi
  threshold=$((en_lines * threshold_pct / 100))
  if [ "$diff" -gt "$threshold" ]; then
    echo "ERROR [$lang]: line count diff > ${threshold_pct}% — $EN=$en_lines vs $target_file=$tgt_lines"
    echo "  diff=$diff, threshold=$threshold (over by $((diff - threshold)))"
    overall_pass=0
    continue
  fi

  # 5. cross-link 양방향
  if ! grep -q "$(basename "$target_file" .md | sed 's/\./\\./g')\\.md" "$EN"; then
    echo "ERROR [$lang]: $EN missing link to $target_file"
    overall_pass=0
    continue
  fi
  if ! grep -q 'README\.md' "$target_file"; then
    echo "ERROR [$lang]: $target_file missing link to $EN"
    overall_pass=0
    continue
  fi

  echo "OK [$lang]: $EN ↔ $target_file sync verified (sections=$en_sections, lines diff=$diff/$threshold)"
done

if [ "$checked_count" = "0" ]; then
  echo "WARNING: no target language file present — skipping all checks"
fi

if [ "$overall_pass" = "1" ]; then
  echo "PASS: 4-lang sync matrix — $checked_count language(s) verified"
  exit 0
else
  echo "FAIL: 4-lang sync matrix — 1 건 이상 drift 발생"
  exit 1
fi
