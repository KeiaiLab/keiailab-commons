# ADR-0015: lefthook P1-12 + P1-13 보강 — go-licenses + markdown-link-check (audit 정합)

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | ADR-0011 (lefthook 통합), ADR-0013 (audit SSOT), S7 (postgres + mongodb 의 동일 패턴) |

## Context

audit P1 의 다음 항목이 commons 에서 ❌:
- P1-12 go-licenses (forbidden / restricted 라이선스 차단)
- P1-13 markdown-link-check (broken link 차단)

postgres + mongodb 는 S7 의 lefthook 3종 보강 PR (#88 / #196) 머지로 이미 ✅.
operator family 정합성 위해 commons 도 동일 3종 보강 (kube-linter 는 operator
전용이라 N/A — commons 는 Helm *library* chart 만, 실 manifest 없음).

## Decision

`lefthook.yml` 의 pre-push 에 두 hook 추가:

```yaml
go-licenses:
  glob: "*.go"
  run: |
    if command -v go-licenses >/dev/null; then
      go-licenses check ./... --disallowed_types=forbidden,restricted 2>&1 | head -20 || echo "[warn] ..."
    else
      echo "[info] go-licenses 미설치 — skip"
    fi
markdown-link-check:
  glob: "*.md"
  run: |
    if command -v markdown-link-check >/dev/null; then
      for f in {staged_files}; do
        markdown-link-check --quiet "$f" 2>&1 | ...
      done
    else
      echo "[info] markdown-link-check 미설치 — skip"
    fi
```

두 hook 모두 *미설치 시 skip* — 개발자 환경 차이 허용. 강제 install 안 함.

## Consequences

- ✅ audit P1-12 + P1-13 commons ✅ 적용
- ✅ operator family 3 operator + commons 일관성 (forgewise 는 Python — go-licenses N/A)
- ⚠️ go-licenses + markdown-link-check 미설치 시 *skip* — 실 차단 아님 (경고만)
  - 개발자 워크플로 마찰 회피
  - CI (현재 GHA 영구 금지로 부재) 부재 대체로 충분치 않음 — 후속 cycle 에서 install 강제 검토
- ⚠️ markdown-link-check 가 일부 false positive (rate limit, 임시 outage) 가능 — 경고만 유지

## Verification

```bash
# 도구 설치 후
brew install markdown-link-check 2>&1 || npm i -g markdown-link-check
go install github.com/google/go-licenses@latest

# lefthook 검증
lefthook run pre-push
# go-licenses + markdown-link-check 두 hook 통과 또는 skip 출력
```

## Migration

- 4 repo (postgres / mongodb / valkey / forgewise) 의 lefthook 에 동일 패턴
- forgewise 는 go-licenses N/A (Python — pyproject.toml 의 license-file 검사 또는 pip-licenses 별 hook)
- valkey 는 ralph-loop 관리 → 본 ADR 의 검토 후 ralph-loop 가 반영
