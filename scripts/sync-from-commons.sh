#!/usr/bin/env bash
# scripts/sync-from-commons.sh
# keiailab-commons → 4 repo (postgres + mongodb + valkey + forgewise) 의 i18n SSOT 동기.
#
# S4 Phase 3 (2026-05-21, docs/specs/2026-05-21-i18n-4lang-master-design.md §4.4):
#   RFC-0029 §6.5 sub-repo sync drift seal 정합 — SSOT 본문 그대로 cp.
#
# 동기 대상:
#   - scripts/check-readme-sync.sh (4-lang 매트릭스 script)
#   - scripts/i18n-translate.sh (Claude API skeleton)
#   - docs/i18n/README.md (i18n 정책 SSOT)
#   - docs/i18n/glossary-{ko,ja,zh}.md (용어 사전 3종)
#
# 사용법:
#   ./scripts/sync-from-commons.sh [--dry-run] [--repo <name>]
#
# 옵션:
#   --dry-run            실행 안 함, 계획만 출력
#   --repo <name>        특정 repo 만 동기 (postgres-operator | mongodb-operator | valkey-operator | forgewise)
#   -h, --help           도움말

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMMONS_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WORKSPACE_ROOT="$(cd "$COMMONS_ROOT/.." && pwd)"

# 4 sub-repo (commons 제외)
ALL_REPOS=(
  "postgres-operator"
  "mongodb-operator"
  "valkey-operator"
  "forgewise"
)

# 기본값
DRY_RUN=0
TARGET_REPO=""

usage() {
  cat <<EOF
Usage: $0 [options]

Options:
  --dry-run            실행 안 함, 계획만 출력
  --repo <name>        특정 repo 만 동기 (default: 모든 4 repo)
                       지원: postgres-operator | mongodb-operator | valkey-operator | forgewise
  -h, --help           도움말

Example:
  $0 --dry-run
  $0 --repo postgres-operator
  $0  # 모든 4 repo 동기

동기 대상:
  - scripts/check-readme-sync.sh
  - scripts/i18n-translate.sh
  - docs/i18n/README.md
  - docs/i18n/glossary-{ko,ja,zh}.md

RFC-0029 §6.5 정합 — SSOT 본문 그대로 cp.
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --repo)
      TARGET_REPO="$2"
      shift 2
      ;;
    *)
      echo "ERROR: 알 수 없는 인자: $1" >&2
      usage
      exit 1
      ;;
  esac
done

# 동기 대상 source 파일 (commons SSOT)
SOURCE_FILES=(
  "scripts/check-readme-sync.sh"
  "scripts/i18n-translate.sh"
  "docs/i18n/README.md"
  "docs/i18n/glossary-ko.md"
  "docs/i18n/glossary-ja.md"
  "docs/i18n/glossary-zh.md"
)

# repo list 결정
if [ -n "$TARGET_REPO" ]; then
  REPOS=("$TARGET_REPO")
else
  REPOS=("${ALL_REPOS[@]}")
fi

echo "==================================================="
echo "sync-from-commons.sh — S4 Phase 3"
echo "==================================================="
echo "Commons root: $COMMONS_ROOT"
echo "Workspace:    $WORKSPACE_ROOT"
echo "Dry-run:      $DRY_RUN"
echo "Target repos: ${REPOS[*]}"
echo "---"

# source 파일 사전 존재 확인
for src in "${SOURCE_FILES[@]}"; do
  if [ ! -f "$COMMONS_ROOT/$src" ]; then
    echo "ERROR: source 파일 부재 — $COMMONS_ROOT/$src" >&2
    exit 1
  fi
done

# 동기 실행
synced_count=0
for repo in "${REPOS[@]}"; do
  target_root="$WORKSPACE_ROOT/$repo"

  if [ ! -d "$target_root" ]; then
    echo "[skip $repo] $target_root 디렉토리 부재"
    continue
  fi

  echo "[$repo] $target_root"

  for src in "${SOURCE_FILES[@]}"; do
    src_full="$COMMONS_ROOT/$src"
    tgt_full="$target_root/$src"
    tgt_dir="$(dirname "$tgt_full")"

    # target 디렉토리 생성 (필요시)
    if [ "$DRY_RUN" = "0" ]; then
      mkdir -p "$tgt_dir"
    fi

    # 기존 파일과 동일 시 skip
    if [ -f "$tgt_full" ] && cmp -s "$src_full" "$tgt_full"; then
      echo "  [same] $src (변경 없음)"
      continue
    fi

    if [ "$DRY_RUN" = "1" ]; then
      if [ -f "$tgt_full" ]; then
        echo "  [would-update] $src"
      else
        echo "  [would-create] $src"
      fi
    else
      cp "$src_full" "$tgt_full"
      # 실행 권한 보존
      if [ -x "$src_full" ]; then
        chmod +x "$tgt_full"
      fi
      if [ -f "$tgt_full.bak" ]; then
        rm -f "$tgt_full.bak"
      fi
      synced_count=$((synced_count + 1))
      echo "  [sync] $src"
    fi
  done
  echo ""
done

echo "---"
if [ "$DRY_RUN" = "1" ]; then
  echo "dry-run 완료 — 실제 동기 없음. 적용은 --dry-run 제거 후 재실행."
else
  echo "동기 완료 — $synced_count 파일 갱신."
  echo ""
  echo "후속 절차:"
  echo "  1. 각 sub-repo 에서 'git status' + 'git diff' 로 변경 확인"
  echo "  2. 각 repo 별 PR 생성 (commit message: 'chore(i18n): SSOT sync from keiailab-commons')"
  echo "  3. 5 PR 머지 후 4-lang drift check 통합 검증"
fi
exit 0
