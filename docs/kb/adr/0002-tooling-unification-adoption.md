# ADR-0002: Tooling unification — `.golangci.yml` + `Makefile` 도입

- Date: 2026-05-09
- Status: Accepted
- Authors: @eightynine01
- Tags: tooling, lint, makefile

## Context

본 repo 의 도구 표면을 정비합니다. keiailab-commons 는 *라이브러리* 이므로
일부 항목 (`validate` 타겟, Dockerfile HEALTHCHECK) 은 해당 사항이 없습니다.

본 repo 의 도입 시점 상태 (2026-05-09):

- Hook: `.lefthook.yml` ✓ (이미 채택).
- `.golangci.yml`: **부재**.
- `Makefile`: **부재**.
- 거버넌스 산출물 (CONTRIBUTING / GOVERNANCE / SECURITY / MAINTAINERS /
  CODE_OF_CONDUCT 등): 보강 완료.
- ADR layout: `docs/kb/adr/`.

## Decision

본 repo 에서 다음 도구를 도입합니다:

1. `.lefthook.yml` 변경 없음 (이미 적용).
2. `.golangci.yml` 신규 — 라이브러리 특성에 맞게 조정 (`ginkgolinter` 제외,
   `lll` 임계 완화).
3. `.custom-gcl.yml` 신규 — `logcheck` plugin.
4. `Makefile` 신규 — `lint`, `test`, `audit` 3 타겟 (`validate` 제외, SBOM 은
   라이브러리 특성상 선택).

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

- 라이브러리 측 게이트가 downstream consumer 와 동일 패턴 — cross-repo 작업
  시 학습 비용 감소.
- `govulncheck` 가 라이브러리 단위에서도 강제 — downstream consumer 가
  발견하기 전에 차단.
- ADR 디렉토리 정합 완료 (`docs/kb/adr/`).

### 부정 / 트레이드오프

- 신규 `Makefile` 도입으로 contributor 가 `make lint test` 명령에 익숙해질
  학습 비용 (낮음).

## Alternatives Considered

| 대안 | 거절 사유 |
|------|----------|
| `Makefile` 도입 보류 (라이브러리는 `go test` 만으로 충분) | `govulncheck` / coverage report 자동화 부재, downstream 과의 일관성 결손. |
| `.golangci.yml` 라이브러리 최소화 (5 linter) | 표준 18 linter 위반, downstream 패턴과의 drift 재발. |

## References

- 관련 ADR: [ADR-0001](0001-charter.md) (keiailab-commons charter).
