#!/usr/bin/env bash
#
# operator-commons 수동 release 스크립트.
#
# 사용:
#   bash scripts/release.sh v0.8.0
#
# 동작:
#   1. tag 형식 검증 (vMAJOR.MINOR.PATCH).
#   2. working tree clean 검증.
#   3. branch=main 검증.
#   4. 로컬 게이트 — make all (lint + test) + make audit (govulncheck).
#   5. version 정합 검증 — go.mod + charts/keiailab-commons/Chart.yaml.
#   6. CHANGELOG.md 갱신 (cliff 가용 시 자동, 아니면 사람 확인 prompt).
#   7. helm chart package — charts/keiailab-commons → /tmp/keiailab-commons-vX.Y.Z.tgz.
#   8. tag 생성 + push origin.
#   9. gh release create — chart .tgz 첨부 + cliff body (가용 시).
#  10. (옵션) helm publish — bash scripts/helm-publish.sh 호출.
#
# 사전조건:
#   - git remote 'origin' 설정 + main branch 권한.
#   - gh CLI 인증 (gh auth status).
#   - (선택) git-cliff: brew install git-cliff (CHANGELOG/release body 자동 생성).
#   - (선택) syft: brew install syft (SBOM 생성).
#
# 본 스크립트는 *수동* 실행 (GitHub Actions 미사용).

set -euo pipefail

usage() {
  echo "Usage: $0 <version>"
  echo "  version: vMAJOR.MINOR.PATCH (e.g. v0.8.0)"
  exit 1
}

VERSION="${1:-}"
[[ -z "$VERSION" ]] && usage

# 1. tag 형식 검증
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "❌ Invalid version: $VERSION (expected vMAJOR.MINOR.PATCH, e.g. v0.8.0)" >&2
  exit 1
fi
VERSION_NUMERIC="${VERSION#v}"

# 2. working tree clean
if [[ -n "$(git status --porcelain)" ]]; then
  echo "❌ Working tree not clean. Stash or commit first." >&2
  git status --short >&2
  exit 1
fi

# 3. branch=main
BRANCH=$(git branch --show-current)
if [[ "$BRANCH" != "main" ]]; then
  echo "❌ Not on main branch (current: $BRANCH)" >&2
  exit 1
fi

# 3.5. 동기화
git fetch origin main
if ! git merge-base --is-ancestor origin/main HEAD; then
  echo "❌ Local main is behind origin/main. Pull first." >&2
  exit 1
fi

# 4. 로컬 게이트
echo "🔍 Running make all (lint + test)..."
make all
echo "🔍 Running make audit (govulncheck)..."
make audit || echo "[warn] audit fail — 진행 전 확인"

# 5. version 정합
CHART_VERSION=$(grep '^version:' charts/keiailab-commons/Chart.yaml | awk '{print $2}' | tr -d '"')
if [[ "$CHART_VERSION" != "$VERSION_NUMERIC" ]]; then
  echo "❌ Version mismatch: Chart.yaml=$CHART_VERSION vs requested=$VERSION_NUMERIC" >&2
  echo "   Update charts/keiailab-commons/Chart.yaml version field first." >&2
  exit 1
fi

# 6. CHANGELOG.md
if [[ -f CHANGELOG.md ]]; then
  if command -v git-cliff >/dev/null; then
    git-cliff --tag "$VERSION" -o CHANGELOG.md
    git add CHANGELOG.md
    GIT_AUTHOR_NAME="${GIT_AUTHOR_NAME:-$(git config user.name)}" \
    GIT_AUTHOR_EMAIL="${GIT_AUTHOR_EMAIL:-$(git config user.email)}" \
    git commit -s -m "chore(release): CHANGELOG for $VERSION (git-cliff 자동 생성)"
    git push origin main
  else
    echo "⚠️  git-cliff 미설치. CHANGELOG.md 수동 갱신 후 commit + push 필요."
    read -p "CHANGELOG.md 갱신 + commit 완료했나? (y/N) " ans
    [[ "$ans" != "y" ]] && exit 1
  fi
fi

# 7. helm chart package
echo "📦 Packaging helm chart..."
PACKAGE_DIR=$(mktemp -d)
helm package charts/keiailab-commons -d "$PACKAGE_DIR"
CHART_PKG=$(ls "$PACKAGE_DIR"/*.tgz | head -1)
echo "   → $CHART_PKG"

# 8. tag + push
echo "🏷️  Creating tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION

operator-commons $VERSION

See CHANGELOG.md for details.
"
git push origin "$VERSION"

# 9. GitHub Release
RELEASE_BODY=""
if command -v git-cliff >/dev/null; then
  RELEASE_BODY=$(git-cliff --tag "$VERSION" --strip all)
elif [[ -f CHANGELOG.md ]]; then
  RELEASE_BODY=$(awk "/## \[$VERSION_NUMERIC\]/,/## \[/" CHANGELOG.md | head -n -1)
fi

echo "🚀 Creating GitHub release..."
if [[ -n "$RELEASE_BODY" ]]; then
  gh release create "$VERSION" "$CHART_PKG" --title "$VERSION" --notes "$RELEASE_BODY"
else
  gh release create "$VERSION" "$CHART_PKG" --title "$VERSION" --generate-notes
fi

# 10. (옵션) helm publish
if [[ -f scripts/helm-publish.sh ]]; then
  read -p "📤 helm-publish 실행? (y/N) " ans
  if [[ "$ans" == "y" ]]; then
    bash scripts/helm-publish.sh "$VERSION"
  fi
fi

echo "✅ Release $VERSION complete."
echo ""
echo "다음 단계:"
echo "  1. downstream operator 의 go.mod 에 operator-commons $VERSION 채택"
echo "  2. UPGRADING.md 의 v$VERSION_NUMERIC 섹션 검토"
echo "  3. (Sprint 1 / S5) pkg/* 추출 패키지 변경 시 deprecation 경로 확인"

# cleanup
rm -rf "$PACKAGE_DIR"
