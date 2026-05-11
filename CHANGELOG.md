# Changelog

본 라이브러리의 모든 주요 변경은 본 파일에 기록된다.
형식: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
버저닝: [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

자동 생성: `git-cliff` (RFC-0017 + standards/enforcement.md §2.3) — release tag 시점에 PR 자동 갱신 후속.

## [Unreleased]

## [0.7.0] - 2026-05-09

### Added

- `pkg/version`: generic `Matrix[E MatrixEntry]` 추가 — postgres operator 의 풍부한 Combo struct 등 caller-supplied entry type 위임 지원 (PR-B1, #2, c2fd68f).
- `docs/kb/adr/0004-*` 신규 — Matrix generic 도입 결정 근거.

## [0.6.0] - 2026-05-09

### Added

- `pkg/status`: `SetAvailable` + `SetReadyFalse` 슈가 헬퍼 추가 (#1, 9891bf3).
- `docs/kb/adr/0003-*` 신규 — status 슈가 헬퍼 결정 근거.
- RFC-0018 본문 추가.
- `.codecov.yml` 신규 — 4-repo target 70% 절대 floor 통일 (fb7c8c6).

## [0.5.0] - 2026-05-09

### Added

- audit (4-repo cross-cut, 2026-05-09) 산출물:
  - 거버넌스 doc 7건 신설 (AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY / MAINTAINERS / CODE_OF_CONDUCT / ADOPTERS) — checklist.md §5 위반 해소 (2745d72).
  - `pkg/status/`: 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼 (RFC-0018 §3.1 implementation). 외부 의존: `k8s.io/apimachinery` 만 (eaa3dd1).
  - `pkg/finalizer/`: Add/Remove/Has 헬퍼 (RFC-0018 §3.2 implementation). controller-runtime 의존 회피, std `slices` 만 사용 (eaa3dd1).
  - `templates/observability/_servicemonitor.tpl`: helm chart partial named template "keiailab.observability.serviceMonitor" (RFC-0019 §3.1 implementation, eaa3dd1).
  - `templates/observability/README.md`: 메트릭 명명 규약 + 공통 alert 권장 + consumer chart 사용법.
  - `Makefile` 신규 (lint/test/audit/cover/tidy/tag) — RFC-0017 §3.3 (f796b6a).
  - `.golangci.yml` + `.custom-gcl.yml` 신규 (postgres 패턴 라이브러리 단순화) — RFC-0017 §3.2 (f796b6a).
  - `CHANGELOG.md` 신설 + `docs/kb/deps/2026-05.md` 의존성 audit log 신규 (19d0c8d).
  - `docs/kb/adr/0002-rfc-0017-tooling-unification-adoption.md` + `docs/kb/adr/INDEX.md` 신규.
  - `NOTICE` 신설 (Apache-2.0 §4(d) 정합, dc10f21).
  - `CODEOWNERS` 신설 (4-repo 통일, 97bac9e).
  - README badges 신설 — License/Go/pkg.go.dev/OpenSSF Scorecard (77f6b11) + Discussions badge + Community 섹션 (f0b72a5).
  - `renovate.json` 신설 (4-repo 정합, 2d97d17).
  - `lefthook.yml` 신설 (라이브러리 minimal, 3d1a31e).
  - DCO Signed-off-by warn-only commit-msg gate (21a18d5).
  - 3 operator adoption matrix + commit references (d3e843b).

### Changed

- ADR 디렉토리 이전: `docs/adr/` → `docs/kb/adr/` (3 operator 와 layout 정합, 2745d72).
- `go` directive `1.25.0` → `1.25.7` (4-repo 정합, 0fef7a6).
- `pkg/finalizer` lint fix: 수동 for+== → `slices.Contains`/`slices.Index`/`slices.Delete` (modernize linter, 201816f / cdf524c).
- `pkg/status/conditions.go` SetReady 함수 시그니처 multi-line (lll 통과, cdf524c).

## [v0.4.0] - 2026-05-07

(이전 버전 history 는 git tag log 또는 release notes 참조 — 본 CHANGELOG.md 는 audit 시점에 신설)

## Refs

- 자동 생성 표준: standards/enforcement.md §2.3 (`git-cliff`)
- 본 audit plan: ai-dev plan `mongodb-operator-operator-commons-postgr-tranquil-horizon`
