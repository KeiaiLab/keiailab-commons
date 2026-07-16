#!/usr/bin/env bash

set -euo pipefail

ROOT="${DOC_SYNC_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
cd "$ROOT"

readmes=(README.md README.ko.md README.ja.md README.zh.md)
roadmaps=(docs/ROADMAP.md docs/ROADMAP.ko.md docs/ROADMAP.ja.md docs/ROADMAP.zh.md)
failed=0
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

report_error() {
  echo "ERROR: $*" >&2
  failed=1
}

find pkg -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | LC_ALL=C sort >"$tmpdir/actual-packages"
if [ ! -s "$tmpdir/actual-packages" ]; then
  report_error "no public package directories found under pkg/"
fi

compare_package_sets() {
  local label="$1"
  local documented="$2"
  local missing extra duplicates

  missing="$(comm -23 "$tmpdir/actual-packages" "$documented")"
  extra="$(comm -13 "$tmpdir/actual-packages" "$documented")"
  duplicates="$(sort "$documented" | uniq -d)"

  [ -z "$missing" ] || report_error "$label missing packages: $(tr '\n' ' ' <<<"$missing" | sed 's/ $//')"
  [ -z "$extra" ] || report_error "$label contains unknown packages: $(tr '\n' ' ' <<<"$extra" | sed 's/ $//')"
  [ -z "$duplicates" ] || report_error "$label contains duplicate packages: $(tr '\n' ' ' <<<"$duplicates" | sed 's/ $//')"
}

for readme in "${readmes[@]}"; do
  if [ ! -f "$readme" ]; then
    report_error "missing file: $readme"
    continue
  fi

  table="$tmpdir/${readme//\//_}.table"
  rows="$tmpdir/${readme//\//_}.rows"
  awk '
    /^## (Packages|패키지|パッケージ|包)$/ { in_packages=1; next }
    /^## / && in_packages { exit }
    in_packages { print }
  ' "$readme" >"$table"

  if [ ! -s "$table" ]; then
    report_error "$readme missing package table section"
    continue
  fi

  sed -nE 's/^\| `pkg\/([^`]+)` \| (Stable|Beta|Experimental) \| .+ \|$/\1/p' "$table" \
    | LC_ALL=C sort >"$rows"
  package_row_count="$(grep -Ec '^\| `pkg/[^`]+` \|' "$table" || true)"
  parsed_row_count="$(wc -l <"$rows" | tr -d ' ')"
  if [ "$package_row_count" -ne "$parsed_row_count" ]; then
    report_error "$readme has malformed package rows: found=$package_row_count parsed=$parsed_row_count"
  fi
  compare_package_sets "$readme package table" "$rows"
done

if [ ! -f docs/AGENTS.md ]; then
  report_error "missing file: docs/AGENTS.md"
else
  agents_rows="$tmpdir/agents-packages"
  sed -n '/^├── pkg\/$/,/^├── charts\//p' docs/AGENTS.md \
    | sed -nE 's/^│   (├──|└──) ([^/]+)\/.*/\2/p' \
    | LC_ALL=C sort >"$agents_rows"
  compare_package_sets "docs/AGENTS.md pkg tree" "$agents_rows"

  root_files_claim="$(sed -n '/^Root holds /,/^$/p' docs/AGENTS.md)"
  if [ ! -e NOTICE ] && grep -Fq 'NOTICE' <<<"$root_files_claim"; then
    report_error "docs/AGENTS.md claims missing root file: NOTICE"
  fi
fi

mapfile -t latest_tags < <(git tag --list 'v[0-9]*' --sort=-version:refname | head -n 3)
latest_release=""
if [ "${#latest_tags[@]}" -ne 3 ]; then
  report_error "expected at least 3 release tags, found ${#latest_tags[@]}"
else
  latest_release="${latest_tags[0]}"
  for tag in "${latest_tags[@]}"; do
    version="${tag#v}"
    release_date="$(git log -1 --format=%ad --date=short "$tag")"
    expected_heading="## [$version] — $release_date"
    if ! grep -Fxq "$expected_heading" CHANGELOG.md; then
      report_error "CHANGELOG.md missing exact release heading: $expected_heading"
    fi
  done
fi

for roadmap in "${roadmaps[@]}"; do
  if [ ! -f "$roadmap" ]; then
    report_error "missing file: $roadmap"
    continue
  fi

  tier_heading="$(grep -m1 '^## API .*v[0-9]' "$roadmap" || true)"
  if [ -z "$tier_heading" ]; then
    report_error "$roadmap missing API stability tier heading"
  elif [ -n "$latest_release" ] && [[ "$tier_heading" != *"$latest_release"* ]]; then
    report_error "$roadmap API stability tier is not $latest_release: $tier_heading"
  fi
done

if [ "$failed" -ne 0 ]; then
  echo "FAIL: package tables, package tree, release headings, or roadmap release are out of sync" >&2
  exit 1
fi

echo "PASS: package tables/tree match pkg/, latest release headings match tag dates, and roadmaps name the latest release"
