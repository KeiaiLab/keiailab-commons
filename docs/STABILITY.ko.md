# API 안정성 약속

> ⚠️ This translation is AI-generated and pending native review. — 본 번역은 Claude 기계 번역 결과입니다. 모국어 검토자의 검수 전까지 `[검토 필요]` 상태입니다.

> operator-commons 의 API 안정성 등급과 breaking-change 정책.

## 등급 (Tier)

operator-commons 는 3-tier 안정성을 사용합니다:

| Tier | 호환성 | Breaking change 정책 |
|---|---|---|
| **Stable** | minor release 간 backwards-compatible. patch fix 만. | major version bump (`v2.0.0`) + 6개월 이상의 deprecation notice 필요 |
| **Beta** | patch release 간 호환. minor release 는 최소 1-release deprecation window 와 함께 source-compatible API 개선 포함 가능 | `CHANGELOG.md` 의 deprecation 항목과 함께 minor version bump 허용 |
| **Experimental** | 호환성 약속 없음. 어느 release 에서나 breaking 가능. | 어느 release 든 — `CHANGELOG.md` "BREAKING" 섹션에 flag 필수 |

## 현재 tier matrix

`ROADMAP.md` "API Stability Tier" 표 기준:

| 패키지 | Tier | 격상 기준 |
|---|---|---|
| `pkg/finalizer` | Stable | (v1 진입 — 추가 작업 없음) |
| `pkg/labels` | Stable | (v1 진입) |
| `pkg/status` | Stable | (v1 진입) |
| `pkg/version` | Beta | Generic `Matrix[E]` 3-repo verify |
| `pkg/monitoring` | Beta | ServiceMonitor 3-repo e2e |
| `pkg/networkpolicy` | Beta | 4-direction TCP/UDP verify |
| `pkg/security` | Beta | restricted PSA 3-repo guard |
| `pkg/webhook` | Experimental | Multi-repo adoption + stabilize |

## 격상 절차

1. Sub-task PR 가 격상 제안과 함께 열림 (예: `feat(pkg/X): promote to Stable`)
2. 격상 기준 (ROADMAP 기준) 이 CI 로 검증:
   - Cross-repo import 통과 (3 operator)
   - 패키지 godoc coverage ≥80%
   - Unit + integration test coverage ≥85%
   - exported API 에 `// TODO` / `// FIXME` 없음
3. ROADMAP.md tier 표 동일 PR 에서 업데이트
4. `CHANGELOG.md` "Changed" 항목 추가

## Breaking-change 정책

**breaking change** = 다음 중 하나:
- exported identifier 제거 (function / type / constant / variable)
- exported signature 변경 (parameter, return type)
- 패키지 제거
- caller code 수정이 필요한 동작 변경

### 각 tier 별:

- **Stable**: v2.0.0 까지 금지. deprecation 사용: `// Deprecated: ...` godoc + 신규 대안 추가, 기존은 6개월 이상 유지
- **Beta**: 1-release deprecation 과 함께 허용. `CHANGELOG.md` "Deprecated" → "Removed" pipeline 에 등장 필수
- **Experimental**: 어느 release 든 허용, `CHANGELOG.md` "BREAKING" 섹션 등장 필수

## Semantic versioning

`vMAJOR.MINOR.PATCH` per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html):
- **MAJOR**: Stable tier 의 breaking change — v2.0.0 graduation review 필요
- **MINOR**: 신규 feature, Beta tier 추가, non-breaking Stable 개선
- **PATCH**: bug fix 만, API surface 변경 없음

## v1.0.0 graduation

다음 *모두* 필요:
1. 8/8 패키지가 Stable tier 도달
2. 6+ 연속 minor release (v0.8 → v0.13) 동안 BREAKING CHANGE 0
3. godoc coverage ≥80% (본 문서 + per-package — `scripts/godoc-coverage.sh` 로 검증)
4. CITATION.cff + Zenodo DOI
5. `v1.0.0-rc.N` 의 3-repo import e2e 검증
6. `go vet ./... && go test ./...` clean + coverage ≥85%
7. CHANGELOG.md v0.x evolution history + v1.0.0 release notes
8. 본 `docs/STABILITY.md` (현재 파일)
9. `pkg/finalizer` multi-finalizer order 보장
10. `pkg/labels` K8s 1.30+ v2 mapping
11. `pkg/status` Condition reason catalog 문서

추적: `~/.claude/plans/2026-05-14-4-operators-100pct/P-B.md` (29 sub-task).

## Caller 책임

Caller (mongodb-operator, valkey-operator, postgres-operator):
- v1.0.0 까지 `go.mod` 에서 `vMAJOR.MINOR.PATCH` 로 pin
- deprecation warning 을 위해 `CHANGELOG.md` 구독
- GA 전 `v1.0.0-rc.N` 에 대해 테스트

## 참조

- `ROADMAP.md` — tier 표 + graduation 기준
- `CHANGELOG.md` — 버전 history
- `CITATION.cff` — academic citation
- `ADOPTERS.md` — 3-repo adoption matrix
- [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html)
