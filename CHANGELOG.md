# Changelog

본 라이브러리의 모든 주요 변경은 본 파일에 기록된다.
형식: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
버저닝: [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

자동 생성: `git-cliff` (RFC-0017 + standards/enforcement.md §2.3) — release tag 시점에 PR 자동 갱신 후속.

## [Unreleased]

### Added

- audit (4-repo cross-cut, 2026-05-09) 산출물:
  - 거버넌스 doc 7건 신설 (AGENTS / GOVERNANCE / CONTRIBUTING / SECURITY / MAINTAINERS / CODE_OF_CONDUCT / ADOPTERS) — checklist.md §5 위반 해소.
  - `pkg/status/`: 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼 (RFC-0018 §3.1 implementation). 외부 의존: `k8s.io/apimachinery` 만.
  - `pkg/finalizer/`: Add/Remove/Has 헬퍼 (RFC-0018 §3.2 implementation). controller-runtime 의존 회피, std `slices` 만 사용.
  - `templates/observability/_servicemonitor.tpl`: helm chart partial named template "keiailab.observability.serviceMonitor" (RFC-0019 §3.1 implementation).
  - `templates/observability/README.md`: 메트릭 명명 규약 + 공통 alert 권장 + consumer chart 사용법.
  - `Makefile` 신규 (lint/test/audit/cover/tidy/tag) — RFC-0017 §3.3.
  - `.golangci.yml` + `.custom-gcl.yml` 신규 (postgres 패턴 라이브러리 단순화) — RFC-0017 §3.2.
  - `docs/kb/adr/0002-rfc-0017-tooling-unification-adoption.md` 신규.
  - `docs/kb/adr/INDEX.md` 신규.
  - `docs/kb/deps/2026-05.md` 신규.

### Changed

- ADR 디렉토리 이전: `docs/adr/` → `docs/kb/adr/` (3 operator 와 layout 정합).
- `pkg/finalizer` lint fix: 수동 for+== → `slices.Contains`/`slices.Index`/`slices.Delete` (modernize linter).
- `pkg/status/conditions.go` SetReady 함수 시그니처 multi-line (lll 통과).

## [v0.4.0] - 2026-05-07

(이전 버전 history 는 git tag log 또는 release notes 참조 — 본 CHANGELOG.md 는 audit 시점에 신설)

## Refs

- 자동 생성 표준: standards/enforcement.md §2.3 (`git-cliff`)
- 본 audit plan: ai-dev plan `mongodb-operator-operator-commons-postgr-tranquil-horizon`
