# AGENTS.md — operator-commons

본 문서는 AI 에이전트 / contributor 가 `keiailab/operator-commons` 를 작업할 때 따라야 할 규약을 요약합니다. 글로벌 규약은 `~/Documents/ai-dev/CLAUDE.md` (Tier-1) 와 `~/Documents/ai-dev/standards/*.md` (Tier-2) 를 우선합니다.

## 본 repo 의 특성

- *라이브러리* (Go module). consumer 는 3 keiailab operator (mongodb / postgres / valkey).
- 공개 API breaking 은 항상 RFC + 3 consumer 동시 LGTM 필요 (GOVERNANCE.md).
- *cmd/main.go 없음* — 실행 가능한 바이너리 미생성.
- 테스트는 unit only (envtest / e2e 는 consumer operator 측에서).

## 디렉토리 규약

```
operator-commons/
├── pkg/
│   ├── labels/         # 리소스 라벨 빌더 (4-repo 공통)
│   ├── security/       # PodSecurity restricted SecurityContext
│   ├── webhook/        # 버전 validation 헬퍼
│   ├── monitoring/     # ServiceMonitor / PrometheusRule 빌더 (예약/부분)
│   ├── networkpolicy/  # deny-by-default NetworkPolicy 빌더
│   └── version/        # SemVer 비교 + 지원 버전 매트릭스
├── docs/
│   └── kb/adr/         # Architecture Decision Records (3 operator 와 동일 layout)
├── go.mod
└── *.md                # 거버넌스 산출물 (CONTRIBUTING / GOVERNANCE / SECURITY / ...)
```

## 작업 시 반드시 확인할 것

1. **공개 API 변경 여부**: 함수/타입 시그니처가 바뀌었으면 → consumer 영향 분석 + RFC.
2. **`go mod tidy` drift**: lefthook pre-push 가 강제. 위반 시 push 차단.
3. **테스트 커버리지**: `pkg/<sub>` 변경 시 단위 테스트 의무 (testing.md §1 위반).
4. **Sonatype 검증**: 새 의존성 추가 전 사용자 home `~/.claude` 의 Sonatype guide skill 호출.
5. **Apache-2.0 license-only**: AGPL/BUSL transitive 0건 (ADR-0001 charter).

## 게이트 체계 (RFC-0002 / standards/ci.md)

- L1 pre-commit: gofmt / govet / golangci-lint (`.lefthook.yml`)
- L2 pre-push: go test / govulncheck / gitleaks / go mod tidy drift 검사
- L3 Makefile (도입 예정): `make lint test validate audit`
- L4 PR review: PR 본문에 §2 게이트 PASS 증거 인용

## ADR 위치

`docs/kb/adr/` (3 operator repo 와 동일 — 2026-05-09 `docs/adr/` 에서 이전).

## 빈번한 작업 패턴

- **새 공통 패키지 추가** (예: `pkg/status`, `pkg/finalizer`):
  1. issue + 7일 윈도우
  2. consumer operator 1곳에 cherry-pick PoC
  3. 3 consumer 모두 import 가능 검증
  4. ADR 작성 → release tag bump

- **버그 픽스 (단일 패키지)**:
  1. 단위 테스트 재현
  2. PR + 1 LGTM
  3. patch release tag

- **의존성 업그레이드**:
  1. Renovate / Dependabot 자동 PR 또는 수동
  2. consumer operator 측 go.mod 검증 (replace directive 활용)
  3. 머지 + tag

## 안티패턴

- *consumer-specific* 코드를 commons 에 넣기 — 본 repo 는 *4-repo 공통* 만.
- 새 패키지를 RFC 없이 추가 — 3 consumer 사용 의도 검증 필수.
- 공개 API 시그니처를 silently 변경 — 항상 RFC.

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
