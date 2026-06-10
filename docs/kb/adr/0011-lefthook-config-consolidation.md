# ADR-0011 — lefthook 설정 통합 (.lefthook.yml → lefthook.yml)

| 메타 | 값 |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Supersedes | — |
| Refs | (none) |

## Context

keiailab-commons 는 2026-05-20 시점에 template SSOT template 의 `lefthook.yml` (407 LOC) 을 도입하면서 기존 GitHub Actions 차단 정책 §1 정합 minimal 버전 `.lefthook.yml` (65 LOC) 을 *중복 보존* 했다.

lefthook CLI 는 `lefthook.yml` 을 우선 로드하고 `.lefthook.yml` 은 *fallback override* 로 해석한다. 양 파일 공존 시:
- 신규 `lefthook.yml` 의 Go 가드 (go-vet / go-test / golangci-lint) 는 활성
- 레거시 `.lefthook.yml` 의 Go pre-commit 가드 (gofmt write-staged) 는 *부분적으로 무시* 가능 — drift 부채
- 레거시의 commit-msg (Conventional Commits + DCO) 는 신규의 commit-msg 와 *중복 또는 conflict* 가능

또한 신규 `lefthook.yml` 의 commit-msg `commitlint` 는 *Node repo 만* 동작 (`package.json` 부재 시 skip) — Go-only repo 에서는 Conventional Commits 강제 부재.

## Decision

1. 레거시 `.lefthook.yml` (65 LOC) 을 `git rm` 으로 삭제.
2. 신규 `lefthook.yml` 에 다음 hook 통합 (레거시 hook 의 의도 그대로 이식):
   - **pre-commit**: `go-fmt` 를 *write-staged* 동작으로 변경 (`gofmt -l -w {staged_files}` + `stage_fixed: true`) — 기존은 check-only.
   - **commit-msg**: `conventional-commits` 신규 hook (Go-only repo 대응, commitlint 미동작 환경).
   - **commit-msg**: `dco-signoff` 신규 hook (A-P0-5, `DCO_STRICT=0` 일시 우회 가능).
   - **pre-push**: `govulncheck` 신규 hook (Go module CVE call-graph 분석).
   - **pre-push**: `go-mod-tidy` drift check 신규 hook (go.mod/go.sum 정합 강제).
3. 단일 SSOT 파일 = `lefthook.yml` 로 통일.

## Consequences

- (+) Go 가드 (gofmt write-staged) 가 pre-commit 시점에 *auto-fix* 적용 (개발자 마찰 0).
- (+) Conventional Commits + DCO 가 Go-only repo 에서도 강제 (GitHub Actions 차단 정책 §1 로컬 4계층 정합).
- (+) Go module CVE 가드 (govulncheck) 활성 — library 다운스트림 (downstream operator) 보호.
- (+) go.mod / go.sum drift 차단 — 배포 신뢰성 게이트.
- (+) 설정 파일 1개 = drift 발생 가능성 0.
- (+) template SSOT (sub-repo drift 정책 §6.5 sync drift seal) 정합.
- (-) 레거시 minimal 정의는 git 이력으로만 추적 (commit hash + 본 ADR cross-link).
- (-) lefthook.yml 의 LOC 증가 (407 → ~470) — 단일 파일 길이 부담, sub-yaml include 미사용.

## Alternatives

- (B) 양 파일 보존 — lefthook fallback 의 anti-pattern. ADR 부재 시 §5 실패.
- (C) 신규 파일 폐기 + 레거시 유지 — multi-stack 가드 자산 손실.

## Cross-link

- 관련 RFC: GitHub Actions 차단 정책 (GHA 영구 금지 → 로컬 4 계층 일원화)
