# ADR-0017: v3.x-stable baseline 인정 (audit ❌ 0 충족, 본 repo = audit SSOT)

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | commons-ADR/0013 (audit SSOT, self-reference), commons-ADR/0014 (release.sh), commons-ADR/0015 (lefthook P1-12/13), commons-ADR/0016 (Sprint 1 pvc + topology), CLAUDE.md §7 (v3.x-stable 정의) |

## Context

CLAUDE.md §7: "본 규약은 **상용 제품 수준**의 다중 프로젝트 일관성을 목표로 한다 — `standards/enforcement.md`의 P0+P1+P2 자동화 모두 충족 시 *v3.x-stable* 선언."

본 repo (operator-commons) 는 두 역할을 동시에 수행한다:
1. **audit 측정의 SSOT** — `scripts/audit-production-grade.sh` (commons-ADR/0013) 가 5 repo 전체를 측정 (self-reference).
2. **3 operator 의 공유 Go 라이브러리 + Helm partials** — 회귀 차단 핵심 SPOF.

따라서 본 repo 의 v3.x-stable 인정은 *self-audit* 의 성격을 갖는다 — 측정 도구 자체가 자기 자신을 측정. 검증 정합성은 시계열 audit-history (commons docs/quality/audit-history.md) 에 누적된 5 repo 결과 + 본 repo 의 P0~C 항목 통과 시점이 SSOT.

### 1. audit ❌ 0 측정 — 2026-05-21 15:30 (self-reference)

`commons/scripts/audit-production-grade.sh /Users/phil/Workspace/keiailab` 출력의 5 repo 표 행:

- P0 (기본 안전): ✅ pre-commit / pre-push / secrets / 한국어 검사 (lefthook + gitleaks)
- P1 (품질 게이트): ✅ lint (`golangci-lint`) / test / typecheck (`go vet`) / build / audit / import-graph (`go mod graph`) / kube-linter (commons-ADR/0015 P1-11) / go-licenses (commons-ADR/0015 P1-12) / markdown-link-check (commons-ADR/0015 P1-13)
- P2 (거버넌스): ✅ ADR coverage (0001~0016) / RFC-0002 GHA block hook 강제 (commons-ADR/0012) / standards/* 정합
- OP (운영): ✅ release.sh 자동화 (commons-ADR/0014) / Go module publish (pkg.go.dev) / Helm library chart .tgz publish
- C (커뮤니티): ✅ ADOPTERS.md (3 operator 매트릭스) / CONTRIBUTING / CODE_OF_CONDUCT / SECURITY / GOVERNANCE / MAINTAINERS (i18n 4-lang) 정합

audit 시계열 기록: [`docs/quality/audit-history.md`](../../quality/audit-history.md) → "🎉 2026-05-21 15:30 — audit ❌ 0 달성" 섹션.

### 2. 거버넌스 baseline

- **RFC-0002 정합** (GitHub Actions 영구 금지) — 본 repo 의 lefthook pre-commit hook 이 `.github/workflows/` 추가 자동 차단 (commons-ADR/0012).
- **i18n 4-lang** (en/ko/ja/zh) README + 거버넌스 문서 + docs/family.{ko,ja,zh}.md — supercycle 2026-05-21 Wave 4 완료.
- **Sprint 1 추출** — commons-ADR/0016 의 pkg/pvc + pkg/topology 신규로 3 operator ~495 LOC 일괄 삭제 가능. v0.9.0 minor bump 대상.
- **라이브러리 분리 원칙**: leaf 패키지 stdlib + k8s API types only. controller-runtime / logr / operator-sdk 누설 금지 (단 pkg/pvc 의 예외는 commons-ADR/0016 명시).

## Decision

본 repo (`keiailab/operator-commons`) 를 **v3.x-stable** 로 인정한다.

- *3 operator (postgres / mongodb / valkey) 의 공유 라이브러리* 로서 외부 사용자 (3 operator 의 fork / 채택자) 대상 운영 등급.
- 후속 release tag `v0.9.0` 권장 (Sprint 1 결과 + audit SSOT 통합) — 구체 버전 (v0.9.0 vs v1.0.0 GA) 은 별 사용자 결정 + 별 ADR 로 추적.
- 본 ADR 자체는 *baseline 인정* 만 — 실 tag 행위는 사용자가 별도 명시.

## Consequences

### Positive

- **외부 신뢰** — audit 자동 측정 (commons-ADR/0013) + 본 baseline ADR + 거버넌스 4종 + i18n 4-lang + 3 operator 채택 매트릭스의 5 축이 *상용 등급* 라이브러리 신호로 작용.
- **거버넌스 SSOT** — 5 repo 중 본 repo 가 audit 측정 + audit-history + family 페이지 + standards 거버넌스 SSOT 의 단일점.
- **회귀 차단** — 본 repo 의 break 는 3 operator 의 동시 break. CI 게이트 보강 (cross-validation Wave 5 후속) 으로 미연 방지.
- **타 family 의 reference** — 후속 sister project (forgewise + 미래 추가) 가 본 audit 패턴 + ADR 패턴 + i18n 패턴을 reference.

### Negative / 회귀 차단 조건

- **audit ❌ ≥ 1 회귀 시** — 본 repo 가 SSOT 이므로 *자기 평가의 모순*. v3.x-stable 인정 *유지 불가* + audit-history 시계열 기록 + 본 ADR 갱신 필수.
- **API 표면 breaking change** — 3 operator 동기 bump 필수. commons-ADR/0016 의 cross-validation Makefile target 강제.
- **standards/* 일탈 시** — ADR 부재면 §5 실패. 일탈 자체는 ADR 동반 시 허용.
- **i18n drift** — 4-lang 거버넌스 문서 sync 강제 (readme-i18n-sync hook).

### Trade-offs

- *v3.x-stable 본 선언* (본 ADR) vs *RFC-0005 글로벌 선언 대기* — 본 repo 는 baseline 만 인정하고 글로벌 RFC-0005 는 별 사용자 결정 영역으로 분리. 글로벌 선언 부재 시에도 본 repo 의 audit ❌ 0 자체가 *측정 가능한 운영 등급 신호*.
- *현 v0.8.0 alpha 인정* vs *v1.0 GA 격상* — 본 ADR 은 격상 강제 안 함. Sprint 1 결과 v0.9.0 정합 (commons-ADR/0016) 후 v1.0.0 별 결정.

## 후속 (v3.1+)

본 baseline 후 v3.1+ 진화 후보:
- P3 성능 게이트 (commons 패키지 benchmark + budget) — 별 RFC
- P4 DR 게이트 (3 operator 의 backup/restore 공통 helper 추출) — 별 RFC
- P5 커뮤니티 KPI (이슈 응답 SLA + 외부 adopter 성장 + go-doc 게시) — 별 RFC
- audit 자동 측정 cron (월 1회) + audit-history 자동 갱신 — 별 ADR
- 5 repo cross-validation Wave 5 자동화 (Makefile target) — 별 ADR

## 참조

- commons-ADR/0012: RFC-0002 GHA block hook
- commons-ADR/0013: `audit-production-grade.sh` 5 repo SSOT 측정 자동화
- commons-ADR/0014: release.sh — 라이브러리 수동 release pipeline
- commons-ADR/0015: lefthook P1-12 + P1-13 보강 (go-licenses + markdown-link-check)
- commons-ADR/0016: Sprint 1 pkg/pvc + pkg/topology 신규 추출
- commons audit-history (시계열): [docs/quality/audit-history.md](../../quality/audit-history.md)
- CLAUDE.md §7 (v3.x-stable 정의): https://github.com/keiailab/.codex (글로벌 standards, private)
- ADOPTERS.md (3 operator 채택 매트릭스): https://github.com/keiailab/operator-commons/blob/main/docs/ADOPTERS.md
