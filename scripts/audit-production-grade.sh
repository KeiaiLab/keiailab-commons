#!/usr/bin/env bash
# keiailab Family — Production-Grade Audit (v0.1)
#
# Usage:
#   $ audit-production-grade.sh [/path/to/5-repo-parent]
#   (default: /Users/phil/Workspace/keiailab)
#
# Output:
#   - 표준 출력: 5 repo × 5 축 (P0/P1/P2/OP/C) 점수 표 (markdown)
#   - 상세: /tmp/audit-production-grade-<timestamp>.md
#
# 정의: /Users/phil/.claude/jobs/61af9ffe/production-grade-checklist.md

set -euo pipefail

PARENT="${1:-/Users/phil/Workspace/keiailab}"
REPOS=(postgres-operator mongodb-operator valkey-operator operator-commons forgewise)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
DETAILED="/tmp/audit-production-grade-${TIMESTAMP}.md"

# ━━━ 측정 함수 ━━━

# Returns "✅" / "❌" / "—" (not applicable)
check_P0_1_lefthook() {
  local repo=$1
  [[ -f "$PARENT/$repo/.lefthook.yml" || -f "$PARENT/$repo/lefthook.yml" ]] && echo "✅" || echo "❌"
}

check_P0_2_gitleaks() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -qE "gitleaks|detect-secrets" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P0_3_dco() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -q "dco-signoff\|dco_signoff\|Signed-off-by" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P0_4_conventional() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -qE "conventional|commitlint|^feat|^fix" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P0_6_no_workflows() {
  local repo=$1
  # 케이스 1: workflow 부재 → ✅ (RFC-0002 strict)
  if [[ ! -d "$PARENT/$repo/.github/workflows" || -z "$(ls -A "$PARENT/$repo/.github/workflows" 2>/dev/null)" ]]; then
    echo "✅"
    return
  fi
  # 케이스 2: workflow 존재 + v2.0 정합 ADR 첨부 → ✅ (RFC-0002 §2 일탈 ADR)
  # 패턴: gha-retention, restore-github-actions, gha-retain
  if find "$PARENT/$repo/docs/kb/adr/" -type f \( -name "*gha-retention*" -o -name "*restore-github-actions*" -o -name "*gha-retain*" -o -name "*gha-removal*" \) 2>/dev/null | grep -q .; then
    echo "✅"
    return
  fi
  echo "❌"
}

check_P0_8_license() {
  local repo=$1
  [[ -f "$PARENT/$repo/LICENSE" ]] && echo "✅" || echo "❌"
}

check_P0_9_mod_drift_hook() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -qE "go-mod-tidy|mod tidy|uv lock|uv-lock-drift|uv sync --check" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

# P1
check_P1_1_lint_target() {
  local repo=$1
  [[ -f "$PARENT/$repo/Makefile" ]] && grep -qE "^lint:" "$PARENT/$repo/Makefile" && echo "✅" || echo "❌"
}

check_P1_3_test_target() {
  local repo=$1
  [[ -f "$PARENT/$repo/Makefile" ]] && grep -qE "^test:" "$PARENT/$repo/Makefile" && echo "✅" || echo "❌"
}

check_P1_11_kube_linter() {
  local repo=$1
  # kube-linter 는 operator 만 (forgewise / commons 는 N/A)
  if [[ "$repo" == "forgewise" || "$repo" == "operator-commons" ]]; then
    echo "—"; return
  fi
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -q "kube-linter" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P1_12_go_licenses() {
  local repo=$1
  if [[ "$repo" == "forgewise" ]]; then echo "—"; return; fi
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml" "$PARENT/$repo/Makefile"; do
    [[ -f $f ]] && grep -q "go-licenses\|go_licenses" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P1_13_md_link() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml"; do
    [[ -f $f ]] && grep -q "markdown-link-check\|md-link" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

# P2 (Governance)
check_P2_2_gha_block_hook() {
  local repo=$1
  for f in "$PARENT/$repo/.lefthook.yml" "$PARENT/$repo/lefthook.yml" "$PARENT/$repo/.husky/pre-commit"; do
    [[ -f $f ]] && grep -q "gha-block\|workflows.*forbid\|RFC-0002\|github_actions" "$f" && { echo "✅"; return; }
  done
  echo "❌"
}

check_P2_7_governance_files() {
  local repo=$1
  local count=0
  for f in SECURITY.md CONTRIBUTING.md CODE_OF_CONDUCT.md CHANGELOG.md; do
    [[ -f "$PARENT/$repo/$f" ]] && count=$((count+1))
  done
  echo "${count}/4"
}

check_P2_4_branding() {
  local repo=$1
  [[ -f "$PARENT/$repo/BRANDING.md" ]] && echo "✅" || echo "❌"
}

check_P2_5_family() {
  local repo=$1
  [[ -f "$PARENT/$repo/docs/family.md" || -f "$PARENT/$repo/family.md" ]] && echo "✅" || echo "❌"
}

check_P2_8_codeowners() {
  local repo=$1
  [[ -f "$PARENT/$repo/.github/CODEOWNERS" || -f "$PARENT/$repo/CODEOWNERS" || -f "$PARENT/$repo/docs/CODEOWNERS" ]] && echo "✅" || echo "❌"
}

check_P2_9_pr_template() {
  local repo=$1
  [[ -f "$PARENT/$repo/.github/PULL_REQUEST_TEMPLATE.md" ]] && echo "✅" || echo "❌"
}

check_P2_10_agents_md() {
  local repo=$1
  [[ -f "$PARENT/$repo/AGENTS.md" ]] && echo "✅" || echo "❌"
}

# Operations
check_OP_1_release_script() {
  local repo=$1
  [[ -f "$PARENT/$repo/scripts/release.sh" ]] && echo "✅" || echo "❌"
}

check_OP_2_helm_publish_script() {
  local repo=$1
  if [[ "$repo" == "forgewise" || "$repo" == "operator-commons" ]]; then echo "—"; return; fi
  [[ -f "$PARENT/$repo/scripts/helm-publish.sh" ]] && echo "✅" || echo "❌"
}

check_OP_5_adopters() {
  local repo=$1
  [[ -f "$PARENT/$repo/ADOPTERS.md" ]] && echo "✅" || echo "❌"
}

check_OP_6_roadmap() {
  local repo=$1
  [[ -f "$PARENT/$repo/ROADMAP.md" || -f "$PARENT/$repo/docs/ROADMAP.md" ]] && echo "✅" || echo "❌"
}

check_OP_7_olm_bundle() {
  local repo=$1
  if [[ "$repo" == "forgewise" || "$repo" == "operator-commons" ]]; then echo "—"; return; fi
  [[ -d "$PARENT/$repo/bundle" ]] && echo "✅" || echo "❌"
}

check_OP_10_upgrade_guide() {
  local repo=$1
  [[ -f "$PARENT/$repo/docs/upgrade.md" || -f "$PARENT/$repo/UPGRADING.md" || -f "$PARENT/$repo/docs/UPGRADING.md" ]] && echo "✅" || echo "❌"
}

# Community
check_C_1_readme_4lang() {
  local repo=$1
  local count=0
  for f in README.md README.ko.md README.ja.md README.zh.md; do
    [[ -f "$PARENT/$repo/$f" ]] && count=$((count+1))
  done
  echo "${count}/4"
}

check_C_2_branding_4lang() {
  local repo=$1
  local count=0
  for f in BRANDING.md BRANDING.ko.md BRANDING.ja.md BRANDING.zh.md; do
    [[ -f "$PARENT/$repo/$f" ]] && count=$((count+1))
  done
  echo "${count}/4"
}

# ━━━ 표 생성 ━━━

print_header() {
  printf "| ID | Check | postgres | mongodb | valkey | commons | forgewise |\n"
  printf "|---|---|---|---|---|---|---|\n"
}

print_row() {
  local id=$1; local label=$2; local fn=$3
  printf "| %s | %s |" "$id" "$label"
  for repo in "${REPOS[@]}"; do
    printf " %s |" "$($fn $repo)"
  done
  printf "\n"
}

main() {
  echo "# keiailab Production-Grade Audit Result"
  echo
  echo "**Generated**: $(date -Iseconds)"
  echo "**Parent**: $PARENT"
  echo
  echo "## P0 — 기본 안전"
  print_header
  print_row "P0-1" "lefthook 설치" check_P0_1_lefthook
  print_row "P0-2" "gitleaks (secrets)" check_P0_2_gitleaks
  print_row "P0-3" "DCO sign-off" check_P0_3_dco
  print_row "P0-4" "Conventional Commits" check_P0_4_conventional
  print_row "P0-6" "RFC-0002 (GHA 0)" check_P0_6_no_workflows
  print_row "P0-8" "LICENSE" check_P0_8_license
  print_row "P0-9" "mod drift hook" check_P0_9_mod_drift_hook
  echo
  echo "## P1 — 품질 게이트"
  print_header
  print_row "P1-1" "make lint target" check_P1_1_lint_target
  print_row "P1-3" "make test target" check_P1_3_test_target
  print_row "P1-11" "kube-linter hook (operator only)" check_P1_11_kube_linter
  print_row "P1-12" "go-licenses (go repo only)" check_P1_12_go_licenses
  print_row "P1-13" "markdown-link-check hook" check_P1_13_md_link
  echo
  echo "## P2 — 거버넌스"
  print_header
  print_row "P2-2" "GHA block hook" check_P2_2_gha_block_hook
  print_row "P2-4" "BRANDING.md" check_P2_4_branding
  print_row "P2-5" "family.md" check_P2_5_family
  print_row "P2-7" "거버넌스 4종 (SEC/CONT/CoC/CHG)" check_P2_7_governance_files
  print_row "P2-8" "CODEOWNERS" check_P2_8_codeowners
  print_row "P2-9" "PR template" check_P2_9_pr_template
  print_row "P2-10" "AGENTS.md" check_P2_10_agents_md
  echo
  echo "## OP — 운영"
  print_header
  print_row "OP-1" "release.sh script" check_OP_1_release_script
  print_row "OP-2" "helm-publish.sh (operator only)" check_OP_2_helm_publish_script
  print_row "OP-5" "ADOPTERS.md" check_OP_5_adopters
  print_row "OP-6" "ROADMAP.md" check_OP_6_roadmap
  print_row "OP-7" "OLM bundle (operator only)" check_OP_7_olm_bundle
  print_row "OP-10" "upgrade guide" check_OP_10_upgrade_guide
  echo
  echo "## C — 커뮤니티"
  print_header
  print_row "C-1" "README 4-lang" check_C_1_readme_4lang
  print_row "C-2" "BRANDING 4-lang" check_C_2_branding_4lang
  echo
  echo "## 요약"
  echo
  echo "- 각 항목 ✅ = 통과 / ❌ = 미통과 / — = N/A"
  echo "- 모든 ✅ 시점 = CLAUDE.md §7 v3.x-stable 선언 조건"
  echo
}

main "$@" | tee "$DETAILED"
echo
echo "상세 저장: $DETAILED"
