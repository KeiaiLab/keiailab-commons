# keiailab Family — Production-Grade Checklist (v0.1)

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 목적 | 5 repo (postgres / mongodb / valkey / commons / forgewise) 의 *상용 제품 수준* 정의 + 측정 + 적용 절차 |
| 출처 | CLAUDE.md §7 (v3.x-stable 정의) + standards/enforcement.md (P0+P1+P2) + 업계 OSS best practice |
| 사용 | 후속 audit cycle 의 기준 + S5 (공통화) 의 helper 추출 우선순위 + 운영 README 의 컨플라이언스 표 |

## 0. 정의 — "상용 제품 수준" (Production-Grade)

OSS Kubernetes operator + MCP server 가 *운영 환경* 에서 *외부 사용자* 에게 *신뢰* 받을 수 있는 수준. 다음 5 축의 모든 기준 통과:

1. **안전 (Safety)**: secrets / SBOM / dependency CVE / DCO / DCO 같은 기본 가드.
2. **품질 (Quality)**: lint / typecheck / test coverage / CI gates / e2e.
3. **거버넌스 (Governance)**: ADR / RFC / SECURITY / CONTRIBUTING / CoC / CHANGELOG / Roadmap.
4. **운영 (Operations)**: release 자동화 / helm chart / multi-arch (opt-in) / monitoring / SBOM / OLM.
5. **커뮤니티 (Community)**: 다국어 / Adopters / Issue templates / PR templates / CODEOWNERS / 응답 SLA.

## 1. P0 (기본 안전) 자동화 — *모든 PR 에서 강제 차단*

| ID | 항목 | 5 repo 측정 명령 | 게이트 |
|---|---|---|---|
| P0-1 | lefthook pre-commit + pre-push 설치 | `ls .lefthook.yml lefthook.yml` | 둘 중 하나 존재 |
| P0-2 | secrets 스캔 통과 없이 push 불가 | lefthook pre-push `gitleaks` step | grep `gitleaks` in lefthook yaml |
| P0-3 | DCO Signed-off-by 강제 | lefthook commit-msg `dco-signoff` step | grep `dco-signoff` |
| P0-4 | Conventional Commits 강제 | lefthook commit-msg `conventional` step | grep `conventional` |
| P0-5 | 한국어 commit msg 본문 허용 (글로벌 §2) | lefthook 정책 또는 ADR | grep `한국어` 또는 글로벌 import |
| P0-6 | `.github/workflows/` 미존재 (RFC-0002) | `test ! -d .github/workflows` | exit 0 |
| P0-7 | DEX (`docker buildx` default builder, amd64 강제) | Makefile `docker-build` target + lefthook `platforms-amd64-guard` | both present |
| P0-8 | LICENSE 파일 존재 | `test -f LICENSE` | exit 0 |
| P0-9 | go.mod / pyproject.toml drift 차단 | lefthook pre-push `go-mod-tidy` 또는 pyproject lock | grep |

## 2. P1 (품질 게이트) — *PR 머지 전 모두 통과*

| ID | 항목 | 측정 | 임계값 |
|---|---|---|---|
| P1-1 | golangci-lint / ruff 통과 | `make lint` 또는 lefthook pre-push | exit 0 |
| P1-2 | typecheck (gopls / mypy strict) | `make typecheck` 또는 lefthook | exit 0 |
| P1-3 | unit test 통과 | `make test` | exit 0 |
| P1-4 | integration / envtest 통과 | `make integration-test` | exit 0 |
| P1-5 | e2e (kind cluster + CR reconcile) | `make e2e` 또는 별 스크립트 | PASS |
| P1-6 | coverage 임계값 | `go test -cover` 또는 `pytest --cov` | ≥ 80% (신규 코드) |
| P1-7 | govulncheck CVE 0 | lefthook pre-push `govulncheck` | exit 0 |
| P1-8 | trivy-fs / trivy-image 통과 | Makefile `audit` target | no HIGH/CRITICAL |
| P1-9 | gosec 통과 | Makefile `audit` (gosec) | no HIGH |
| P1-10 | helm lint + template 통과 | lefthook pre-push `helm-lint` + `helm-template` | exit 0 |
| P1-11 | kube-linter 통과 | lefthook pre-push `kube-linter` (S1+/S7 phase 1 추가) | exit 0 |
| P1-12 | go-licenses 통과 (forbidden 0) | lefthook pre-push `go-licenses` (S1+/S7 phase 1 추가) | exit 0 |
| P1-13 | markdown-link-check 통과 | lefthook pre-push `markdown-link-check` (S1+/S7 phase 1 추가) | exit 0 |
| P1-14 | docs/specs/ + docs/plans/ 사용 (CLAUDE.md §8) | 디렉토리 존재 | exit 0 |

## 3. P2 (거버넌스 정합)

| ID | 항목 | 측정 | 임계값 |
|---|---|---|---|
| P2-1 | standards/* 일탈 시 ADR 첨부 | `docs/kb/adr/` 안 *최근 30일* ADR 존재 | adr 1+ |
| P2-2 | RFC-0002 (GHA 금지) 자동 차단 | pre-commit hook (gha-block) | grep |
| P2-3 | import-graph 정합 (project → global → standards) | AGENTS.md 의 `@import` 디렉티브 | grep `@~/.codex/standards` 또는 `@./standards` |
| P2-4 | 5 repo BRANDING.md 일관성 | 5 repo 모두 BRANDING.md 존재 + family table 동일 | grep |
| P2-5 | 5 repo docs/family.md 일관성 | 동일 family 표 | diff (allow ≤5%) |
| P2-6 | 5 repo README 헤더/푸터 일관성 | 4-lang switcher + family link + © 2026 | grep |
| P2-7 | SECURITY.md / CONTRIBUTING.md / CoC / CHANGELOG | 5 repo 모두 존재 | `ls SECURITY.md CONTRIBUTING.md CODE_OF_CONDUCT.md CHANGELOG.md` |
| P2-8 | CODEOWNERS | `.github/CODEOWNERS` 존재 | exit 0 |
| P2-9 | Issue templates + PR template | `.github/ISSUE_TEMPLATE/` + `.github/PULL_REQUEST_TEMPLATE.md` 존재 | exit 0 |
| P2-10 | AGENTS.md (Tier-3 프로젝트 override) | 존재 | exit 0 |

## 4. 운영 (Operations)

| ID | 항목 | 측정 | 임계값 |
|---|---|---|---|
| OP-1 | release 자동화 (로컬 스크립트, RFC-0002 정합) | `scripts/release.sh` + Makefile `release` | exit 0 |
| OP-2 | helm chart 자동 publish (로컬 스크립트) | `scripts/helm-publish.sh` + Makefile `helm-publish` | exit 0 |
| OP-3 | semver 정책 (tag 형식) | `git tag` 의 `v` prefix + semver | regex |
| OP-4 | CHANGELOG (Keep a Changelog 양식) | `CHANGELOG.md` `## [Unreleased]` + `## [vX.Y.Z]` 섹션 | grep |
| OP-5 | Adopters 목록 (`ADOPTERS.md`) | 존재 + 1+ entry | grep |
| OP-6 | Roadmap (`ROADMAP.md`) | 존재 + 최근 30일 갱신 | mtime |
| OP-7 | OLM bundle (operator 3개 한정) | `bundle/` 디렉토리 + clusterserviceversion.yaml | exit 0 |
| OP-8 | SBOM 생성 (release 시) | `scripts/sbom.sh` 또는 release 안 통합 | exit 0 |
| OP-9 | Multi-arch opt-in (CLAUDE.md §2: 기본 amd64) | Makefile docker-buildx target + opt-in flag | grep |
| OP-10 | upgrade guide (`docs/upgrade.md` 또는 `UPGRADING.md`) | 존재 | exit 0 |
| OP-11 | troubleshooting 가이드 | `docs/troubleshooting.md` 또는 README 안 | grep |

## 5. 커뮤니티

| ID | 항목 | 측정 | 임계값 |
|---|---|---|---|
| C-1 | 4-lang README (en+ko+ja+zh) | `ls README.md README.ko.md README.ja.md README.zh.md` | 4개 |
| C-2 | 4-lang BRANDING (en+ko+ja+zh) | 동일 | 4개 |
| C-3 | 4-lang family (commons SSOT) | 동일 | 4개 |
| C-4 | 4-lang switcher 상단 nav | grep `> English \| 한국어 \| 日本語 \| 中文` | 모든 4-lang 파일 |
| C-5 | Issue 응답 SLA 명시 (CONTRIBUTING) | grep | exit 0 |
| C-6 | discussion / channel 명시 (README) | 링크 | exist |
| C-7 | 다국어 번역 native review 상태 표시 | `[검토 필요]` 또는 `[reviewed]` 배지 | grep |

## 6. 5 repo 측정 — 현재 baseline + 목표

| repo | P0 | P1 | P2 | OP | C | 총합 |
|---|---|---|---|---|---|---|
| postgres-operator | (S7 진행 중) | (보강 중) | TBD | TBD | TBD | TBD |
| mongodb-operator | (S7 진행 중) | (보강 중) | TBD | TBD | TBD | TBD |
| valkey-operator | (S1+ 진행 중) | (보강 중) | TBD | TBD | TBD | TBD |
| operator-commons | (S2+S4 진행 중) | OK | (S4 진행 중) | TBD | (S4 진행 중) | TBD |
| forgewise | (S6 진행 중) | (보강 중) | (S6 진행 중) | (S6 진행 중) | (S6 진행 중) | TBD |

→ 모든 항목 *측정 + ✅* 시점이 **v3.x-stable 선언** 시점 (CLAUDE.md §7).

## 7. 적용 절차 (audit cycle)

```bash
# 각 항목을 5 repo 별로 측정하는 audit 스크립트 (operator-commons 의 scripts/audit-production-grade.sh)
# Phase 1: 측정 (모든 항목 bash 검증)
# Phase 2: 결과 표 자동 생성 (markdown)
# Phase 3: gap 식별 → 후속 sub-cycle (sub-spec)
# Phase 4: gap 모두 해소 → v3.x-stable 선언 commit
```

## 8. 후속 sub-cycle 식별 (4 cycle 완료 후 즉시 진입)

위 checklist 기반 후속 작업:
- **S5** (operator-commons 공통화): P1 / OP / 운영 helper 추출 (3 operator 공통 패턴)
- **S4-A/B/C/D** (4 repo 다국어 적용): C-1~C-7 충족
- **S8** (audit cycle): 모든 P0/P1/P2/OP/C 항목 자동 측정 + 표 생성 + gap 해소
- **S9** (v3.x-stable 선언): 모든 ✅ 시점에 main commit + ADR + release tag

## 9. 본 문서의 향후 위치

본 문서는 임시 `$CLAUDE_JOB_DIR` 작성. 후속 정리:
- (a) operator-commons/docs/quality/production-grade-checklist.md (commons SSOT 정합)
- (b) keiailab 가족 메타 저장소 (없으면 신설)
- (c) 5 repo 각각 link 만 (commons SSOT)

**제안**: (a) — commons SSOT 정합 + 4-lang 번역 가능 + 5 repo 가 link.

---

## 변경 이력

- 2026-05-21 v0.1: 초기 작성 (4 cycle 진행 중 main thread)
