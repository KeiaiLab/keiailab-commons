# ADR-0012: RFC-0002 GitHub Actions Block — lefthook pre-commit hook 자동 강제

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | RFC-0002 (GitHub Actions Permanent Ban), ADR-0011 (lefthook consolidation) |

## Context

`~/.codex/CLAUDE.md` §2 RFC-0002 (2026-04-29): "GitHub Actions 영구 금지 — `.github/workflows/` 디렉토리 추가 / 보존 **금지**".

5 repo audit (2026-05-21) 측정 결과 P2-2 "GHA block hook" 가 모든 repo 에서 ❌:

| repo | P2-2 |
|---|---|
| postgres-operator | ❌ |
| mongodb-operator | ❌ |
| valkey-operator | ❌ |
| operator-commons | ❌ |
| forgewise | ❌ (사전 hook 부재) |

즉 *정책* (RFC-0002) 은 있지만 *자동 강제* (lefthook hook) 이 부재. 사람이 의도하지 않은 `.github/workflows/` 추가 시 차단 불가.

본 ADR 은 operator-commons (다중 언어 SSOT) 에 *최초* gha-block hook 신설 + 후속 4 repo 에 sync 패턴 확립.

### v2.0 정합 고려 (operator family)

2026-05-21 사용자 결정: operator 3 (postgres / mongodb / valkey) 는 v2.0 = GHA *유지* + 통합 ADR (ADR-0048 sister: ADR-0019 / ADR-0033) + 로컬 4계층 dual-track.

따라서 본 hook 은:
- 신규 파일 추가 (`--diff-filter=A`) 만 차단
- 기존 파일 변경 (dependabot 의 actions 버전 bump 등) 은 허용
- 우회: `PLAN_BYPASS=1 git commit` (PR 본문 사유 + ADR 인용 의무)

## Decision

`lefthook.yml` 의 `pre-commit.commands.gha-block` 신설:

```yaml
gha-block:
  run: |
    added=$(git diff --cached --name-only --diff-filter=A | grep "^\.github/workflows/" || true)
    if [ -n "$added" ]; then
      echo "❌ RFC-0002 §2 위반: .github/workflows/ 신규 파일 추가 금지"
      echo "   파일: $added"
      echo "   대체: 로컬 4계층 (lefthook + Makefile + 리뷰어 증거)"
      echo "   예외 (helm-publish/release/scorecard 등) ADR 작성 후 사용자 명시 승인 필요"
      exit 1
    fi
```

## Consequences

- ✅ 신규 GHA workflow 의도치 않은 도입 자동 차단
- ✅ dependabot/renovate 의 기존 workflow 갱신 정상 (modified 만 — 차단 안 함)
- ✅ commons 의 PLAN_BYPASS 패턴 일관 적용 — RFC-0045 §2.5 plan SSOT 차단과 동일 우회 메커니즘
- ⚠️ 4 repo (postgres/mongodb/valkey/forgewise) 에 sync 필요 — `scripts/sync-from-commons.sh` 활용 또는 별 PR 각각
- ⚠️ 본 hook 만으로는 `.github/workflows/` 의 *기존 파일 삭제* 차단 불가 — 의도된 (ADR-0018/0032 의 RFC-0002 strict 노선 일 때 필요한 동작) 이므로 차단 안 함이 정합

## Verification

```bash
# 정상 동작 확인
echo "name: test" > .github/workflows/test.yml
git add .github/workflows/test.yml
lefthook run pre-commit
# 출력에 "RFC-0002 §2 위반" + exit 1

# 기존 파일 변경 허용 확인 (해당 시)
# (commons 는 GHA 0 이므로 본 단계 N/A)

# 정리
git restore --staged .github/workflows/test.yml
rm .github/workflows/test.yml
rmdir .github/workflows 2>/dev/null
```

## Migration

본 ADR 채택 후 후속 sub-cycle:
- S9 sub-cycle: 4 repo (postgres / mongodb / valkey / forgewise) 의 lefthook 에 동일 hook 추가
- valkey 는 ralph-loop 관리 영역 (본 thread 만지지 않음)
- forgewise 는 S4-D 완료 후 (lefthook 변경 충돌 회피)
- postgres + mongodb 는 S7 revert 완료 후
