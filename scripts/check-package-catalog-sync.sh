#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

readmes=(README.md README.ko.md README.ja.md README.zh.md)
roadmaps=(docs/ROADMAP.md docs/ROADMAP.ko.md docs/ROADMAP.ja.md docs/ROADMAP.zh.md)
failed=0

report_error() {
  echo "ERROR: $*" >&2
  failed=1
}

for readme in "${readmes[@]}"; do
  if [ ! -f "$readme" ]; then
    report_error "missing file: $readme"
    continue
  fi

  for package_dir in pkg/*/; do
    [ -d "$package_dir" ] || continue
    package="${package_dir#pkg/}"
    package="${package%/}"
    row="| \`pkg/$package\` |"
    row_count="$(grep -Fc "$row" "$readme" || true)"
    if [ "$row_count" -eq 0 ]; then
      report_error "$readme missing package row: \`pkg/$package\`"
    elif [ "$row_count" -ne 1 ]; then
      report_error "$readme has duplicate package row ($row_count): \`pkg/$package\`"
    fi
  done
done

mapfile -t latest_tags < <(git tag --list 'v[0-9]*' --sort=-version:refname | head -n 3)
latest_release=""
if [ "${#latest_tags[@]}" -ne 3 ]; then
  report_error "expected at least 3 release tags, found ${#latest_tags[@]}"
else
  latest_release="${latest_tags[0]}"
  for tag in "${latest_tags[@]}"; do
    version="${tag#v}"
    if ! grep -Eq "^## \\[$version\\]( |$)" CHANGELOG.md; then
      report_error "CHANGELOG.md missing release: $tag"
    fi
  done
fi

if [ ! -f docs/AGENTS.md ]; then
  report_error "missing file: docs/AGENTS.md"
else
  agents_pkg_tree="$(sed -n '/^├── pkg\/$/,/^├── charts\//p' docs/AGENTS.md)"
  for package_dir in pkg/*/; do
    [ -d "$package_dir" ] || continue
    package="${package_dir#pkg/}"
    package="${package%/}"
    if ! grep -Fq "$package/" <<<"$agents_pkg_tree"; then
      report_error "docs/AGENTS.md pkg tree missing: pkg/$package"
    fi
  done

  root_files_claim="$(sed -n '/^Root holds /,/^$/p' docs/AGENTS.md)"
  if [ ! -e NOTICE ] && grep -Fq 'NOTICE' <<<"$root_files_claim"; then
    report_error "docs/AGENTS.md claims missing root file: NOTICE"
  fi
fi

for roadmap in "${roadmaps[@]}"; do
  if [ ! -f "$roadmap" ]; then
    report_error "missing file: $roadmap"
    continue
  fi

  current_tier_heading="$(grep -m1 '^## .*v[0-9]' "$roadmap" || true)"
  if [ -n "$latest_release" ] && [[ "$current_tier_heading" != *"$latest_release"* ]]; then
    report_error "$roadmap current release heading is not $latest_release: ${current_tier_heading:-<missing>}"
  fi
done

if [ "$failed" -ne 0 ]; then
  echo "FAIL: package and release catalogs are out of sync" >&2
  exit 1
fi

echo "PASS: package and release catalogs are in sync"
