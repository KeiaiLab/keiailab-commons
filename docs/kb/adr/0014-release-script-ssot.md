# ADR-0014: operator-commons release.sh — 라이브러리 수동 release pipeline

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Author | keiailab |
| Supersedes | (none) |
| Related | GitHub Actions 차단 정책 (GHA 영구 금지), ADR-0011 (lefthook 통합), ADR-0013 (audit SSOT) |

## Context

audit P2/OP 의 OP-1 (release.sh script) 가 operator-commons 에서 ❌. 다른 4 repo:
- downstream operator, downstream operator, downstream operator: 본 ADR 시점 ✅ (S7 + valkey ralph-loop 결과로 이미 보유)
- downstream component: ❌ (Python 패키지, 별 release 양식)

operator-commons 는 *Go module 라이브러리* + *Helm library chart* 의 dual deliverable:
- Go module: tag (`vX.Y.Z`) 만으로 release (`go get` 자동)
- Helm library chart: `charts/keiailab-commons/Chart.yaml` 의 version + 별 package + publish

기존 `make tag` target 은 *사용법 안내* 만 — 실 pipeline 부재. 본 ADR 은 *수동 release pipeline* 정합 — GitHub Actions 차단 정책 (GHA 영구 금지) 준수.

## Decision

`scripts/release.sh` 신설 (10 단계):

1. tag 형식 검증 (vMAJOR.MINOR.PATCH)
2. working tree clean 검증
3. branch=main 검증 + origin/main 동기화
4. 로컬 게이트 — `make all` (lint + test) + `make audit` (govulncheck)
5. version 정합 — `go.mod` + `charts/keiailab-commons/Chart.yaml`
6. CHANGELOG.md 갱신 (git-cliff 자동 또는 사람 prompt)
7. helm chart package
8. tag + push origin
9. `gh release create` — chart .tgz 첨부 + cliff body
10. (옵션) `scripts/helm-publish.sh` 호출 (gh-pages 배포)

`Makefile` 의 `release` target: `make release VERSION=v0.8.0` entry point.

## Consequences

- ✅ audit OP-1 ✅ 해소
- ✅ 모든 release pipeline 로컬화 — GitHub Actions 차단 정책 정합
- ✅ Go module + Helm chart dual deliverable 통합
- ⚠️ `scripts/helm-publish.sh` 가 commons 에서 *아직 없음* — 후속 sub-cycle 작성 필요 (단 commons 는 chart "library" 라 publish 가 다른 양식)
- ⚠️ git-cliff 가용 시 자동, 아니면 사람 prompt — 부분 자동화

## Verification

```bash
# 사전 확인 (실제 release 안 함)
bash -n scripts/release.sh  # syntax check

# 실 사용
make release VERSION=v0.8.0
# 또는
bash scripts/release.sh v0.8.0
```

## Migration

본 ADR 이 OP-1 audit 만족시키나, 후속 sub-cycle 필요:
- `scripts/helm-publish.sh` 작성 (library chart 의 ArtifactHub 배포 양식)
- CHANGELOG.md 양식 (Keep a Changelog v1.1.0) 검증 (이미 존재)
- 첫 실 release (v0.8.0) 시 본 스크립트 검증
