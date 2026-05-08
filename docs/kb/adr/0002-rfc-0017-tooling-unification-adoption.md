# ADR-0002: RFC-0017 operator tooling unification 채택

- Date: 2026-05-09
- Status: Proposed
- Authors: @eightynine01
- Tags: tooling, lint, makefile

## Context

ai-dev RFC-0017 가 4 operator repo 도구 통합을 제안한다. operator-commons 는 *라이브러리* 이므로 일부 항목 (validate, Dockerfile HEALTHCHECK) 은 해당 사항 없음.

본 repo 의 현 상태 (2026-05-09 audit):
- Hook: `.lefthook.yml` ✓ (이미 채택)
- `.golangci.yml`: **부재** — RFC-0017 §3.2 위반
- Makefile: **자체 부재** — RFC-0017 §3.3 위반
- 거버넌스 산출물: 2026-05-09 7건 모두 보강 완료 (P0-1 audit action)
- ADR layout: `docs/kb/adr/` (2026-05-09 `docs/adr/` 에서 이전, 3 operator 와 정합)

## Decision

RFC-0017 을 **Accepted** 로 채택하고 본 repo 에서:

1. `.lefthook.yml` 변경 없음 (이미 적용)
2. `.golangci.yml` 신규 — postgres 패턴 단순화 (라이브러리 특성에 맞게 ginkgolinter 제외, lll 임계 완화)
3. `.custom-gcl.yml` 신규 — logcheck plugin
4. **Makefile 신규** — `lint`, `test`, `audit` 3 타겟 (validate 제외, sbom 라이브러리 특성상 선택)

```makefile
.PHONY: lint test audit

lint:
	golangci-lint run --config .golangci.yml ./...

test:
	go test -race -count=1 -coverprofile=cover.out ./...

audit:
	govulncheck ./...
```

## Consequences

### 긍정
- 라이브러리 측 게이트가 consumer operator 와 동일 패턴 — cross-repo 작업 시 학습 비용 감소
- govulncheck 가 라이브러리 단위에서도 강제 — consumer operator 가 발견하기 전에 차단
- ADR 디렉토리 정합 완료 (`docs/kb/adr/`)

### 부정 / 트레이드오프
- 신규 Makefile 도입으로 contributor 가 `make lint test` 명령에 익숙해질 학습 비용 (낮음)

### 후속 작업
- [ ] AI-CM02-1: `.golangci.yml` 적용 후 `pkg/` 18-linter 통과 검증 (Owner: @eightynine01, Due: 2026-05-19)
- [ ] AI-CM02-2: Makefile 신규 + `make lint test audit` PASS 검증 (Owner: @eightynine01, Due: 2026-05-12)
- [ ] AI-CM02-3: RFC-0017 §3.4 follow-up — `pkg/event/reasons.go` 카탈로그 추출 RFC 작성 검토 (Owner: @eightynine01, Due: 2026-05-26)

## Alternatives Considered

| 대안 | 거절 사유 |
|------|----------|
| Makefile 도입 보류 (라이브러리는 go test 만으로 충분) | govulncheck / coverage report 자동화 부재, 4-repo 일관성 결손 |
| `.golangci.yml` 라이브러리 최소화 (5 linter) | RFC-0017 §3.2 18 linter 표준 위반, postgres pattern 과 drift 재발 |

## References

- 글로벌 RFC: `~/Documents/ai-dev/rfcs/0017-operator-tooling-unification.md`
- 관련 ADR: ADR-0001 (operator-commons charter)
- 관련 audit: `~/.claude/plans/mongodb-operator-operator-commons-postgr-tranquil-horizon.md`
