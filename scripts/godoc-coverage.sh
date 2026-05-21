#!/usr/bin/env bash
# godoc-coverage.sh — Measures godoc coverage across all pkg/ subdirectories.
#
# Counts: exported identifier with godoc / total exported identifier.
# Output: `coverage=NN%` single line + per-package breakdown on stderr.
#
# Used by: v1.0.0 graduation criterion (≥80%) per docs/STABILITY.md.
#
# Refs: ROADMAP.md 'v1.0.0 졸업 조건' #3

set -euo pipefail

total=0
documented=0

while IFS= read -r pkg; do
  pkg_total=0
  pkg_documented=0
  for file in "$pkg"/*.go; do
    [ -f "$file" ] || continue
    [[ "$file" == *_test.go ]] && continue
    while IFS= read -r ln; do
      line_no="${ln%%:*}"
      pkg_total=$((pkg_total + 1))
      prev=$((line_no - 1))
      prev_line=$(sed -n "${prev}p" "$file" 2>/dev/null || true)
      if echo "$prev_line" | grep -qE '^//'; then
        pkg_documented=$((pkg_documented + 1))
      fi
    done < <(grep -nE '^(func|type|var|const)\s+\(?[A-Z]' "$file" 2>/dev/null || true)
  done
  total=$((total + pkg_total))
  documented=$((documented + pkg_documented))
  if [ "$pkg_total" -gt 0 ]; then
    pct=$((pkg_documented * 100 / pkg_total))
    echo "  $pkg: $pkg_documented/$pkg_total ($pct%)" >&2
  fi
done < <(find pkg -mindepth 1 -maxdepth 1 -type d 2>/dev/null | sort)

if [ "$total" -gt 0 ]; then
  pct=$((documented * 100 / total))
  echo "coverage=${pct}%"
  echo "(${documented}/${total} exported identifiers documented)" >&2
  if [ "$pct" -lt 80 ]; then
    echo "(warning: below v1.0 80% threshold)" >&2
    exit 1
  fi
else
  echo "coverage=0% (no exported identifiers found)"
  exit 1
fi
