# ADR-0019: ArtifactHub publish pipeline — GHA 2-workflow narrow exception (signed OCI)

- Date: 2026-06-12
- Status: Accepted
- Authors: @phil

## Context

keiailab OSS 4 repo 중 operator 3종 (mongodb / postgres / valkey) 은 ArtifactHub 에
`oci://ghcr.io/keiailab/charts/<name>` OCI repo 로 등록 완료 (verified publisher,
공통 서명키 `F1A6893583E632A757FF6767F3CC8C6AEC9CEB08` — 2026-06-10 회전본) 상태인데,
**keiailab-commons 만 미등록 + GHA publish 파이프라인 부재** (2026-06-12 ArtifactHub
API 실측: 3 operator 등록 / commons 404). 수동 `scripts/release.sh` (ADR-0014) 의존이
ArtifactHub stale 의 구조적 원인 — mongodb-operator 가 동일 문제로 GHA OCI publish 를
도입한 선례 (mongodb ADR-0037) 가 있다.

원안 브랜치 `feat/artifacthub-publish-pipeline` (2026-06-02, ADR 번호 0018 선점됨) 은
구 서명키 (`89A40947...`) + 구명칭 (operator-commons) + Apache-2.0 + gh-pages HTTP 추적
전제라 그대로 부활 불가 — 본 ADR 이 신 사실 기준의 재작성 결정 기록이다.

조직 거버넌스 RFC-0002 (GitHub Actions 영구 금지) 와의 관계: 본 repo 는 GitHub 가
진본인 public OSS 로, postgres-operator ADR-0022 (GHA narrow exception 3-workflows) 가
이미 동일 예외 선례를 봉인했고, 사용자 (거버넌스 결정권자) 가 2026-06-12 본 파이프라인
부활을 명시 지시했다.

## Decision

GHA **2-workflow narrow exception** 을 채택한다:

1. **`.github/workflows/helm-publish.yml`** — release published / `v*` tag /
   workflow_dispatch 트리거. chart 를 GPG 서명 (`--sign` + `.prov`) 하여
   `oci://ghcr.io/keiailab/charts/keiailab-commons` 로 push. 동일 chart version 존재 시
   metadata 일치 검증 후 skip (멱등).
2. **`.github/workflows/artifacthub-verify.yml`** — `ah lint` (PR 포함) +
   ArtifactHub 등록·인덱싱·signed=true smoke (`hack/artifacthub_smoke.sh`,
   tag/dispatch 한정).

핵심 매개 결정:

- **tag ↔ chart version 분리**: repo tag (`v*`) 는 Go 모듈 semver, chart version 은
  `charts/keiailab-commons/Chart.yaml` 이 SSOT. operator 식 tag↔chart parity check 는
  **부적용** — OCI 멱등 skip 이 중복 publish 를 무해화한다.
- **gh-pages HTTP repo 단계 의도적 생략**: 소비 경로가 OCI 뿐 (operator 3종
  `dependencies.repository: oci://ghcr.io/keiailab/charts`) + ArtifactHub 추적 소스도
  OCI. mongodb 패턴의 chart-releaser 단계는 수요 0 으로 제외 (Simplicity First).
- **helm-lint.yml 제외**: 품질 게이트는 로컬 lefthook (ADR-0011/0012/0015) + GitLab CI
  shadow 가 담당 — publish 파이프라인 목적 외 (postgres ADR-0022 narrow 원칙 정합).
- **서명키**: 회전본 `F1A68935...` (mongodb/valkey/postgres 와 동일). 공개키는
  `charts/keiailab-helm-signing-public.asc` + `Chart.yaml` `artifacthub.io/signKey`
  annotation. 비밀키는 repo secret `HELM_SIGNING_PRIVATE_KEY`.
- **ADR-0014 (release.sh) 관계**: GHA publish 가 표준 경로, `scripts/release.sh` 는
  로컬 수동 fallback 으로 존치 (supersede 아님 — Go 모듈 릴리스 단계는 release.sh 가
  계속 담당).
- **ArtifactHub 등록**: `hack/artifacthub_register.sh` (AH REST
  `POST /repositories/org/keiailab`, repo name `keiailab-commons`, kind 0, OCI URL).
  등록 후 발급 UUID 를 `charts/artifacthub-repo.yml` `repositoryID` 에 기입하고 동
  파일을 OCI metadata (`ghcr.io/keiailab/charts/keiailab-commons:artifacthub.io` 태그)
  로 push 하여 verified publisher 달성.

## Consequences

- (+) commons chart 신규 버전이 릴리스 즉시 서명되어 OCI 로 publish — ArtifactHub
  stale 구조 해소, 4 repo 파이프라인 일관성 완성.
- (+) Signed badge + verified publisher 로 다운스트림 (operator 3종 + 외부 사용자)
  공급망 신뢰 확보.
- (−) GHA 의존 2 workflow 신설 — RFC-0002 의 narrow exception 누적 (본 ADR 이 경계
  명시: publish/verify 2종 한정, 품질 게이트 GHA 신설 금지 유지).
- (−) repo secrets 3종 (`HELM_SIGNING_PRIVATE_KEY`, `AH_API_KEY_ID`,
  `AH_API_KEY_SECRET`) 운영 부담 — mongodb-operator 와 동일 값 공유.
- 후속: ArtifactHub 등록 (AH API key 필요) → repositoryID 기입 → OCI metadata push.

## Alternatives Considered

- **원안 브랜치 그대로 머지**: 구 키/구명칭/Apache-2.0/gh-pages 전제 — 전면 stale,
  기각.
- **release.sh 에 OCI 서명 publish 추가 (GHA 0)**: 로컬 수동 의존이 stale 의 원인
  그 자체 — 기각 (mongodb ADR-0037 동일 판단).
- **mongodb 식 15-workflow 전체 미러**: scorecard/codeql 등은 별도 결정 사안 —
  본 ADR 범위 밖 (narrow 원칙).

## Refs

- 선례: postgres-operator ADR-0022 (GHA narrow exception) / mongodb-operator ADR-0037
  (gitops-artifacthub-standardization)
- 관련: ADR-0005 (library chart) / ADR-0012 (GHA block hook — 본 예외와 공존) /
  ADR-0014 (release.sh SSOT) / ADR-0018 (rename)
- 거버넌스: RFC-0002 (GHA 금지 — 본 ADR 이 narrow exception 정당화)
- 실측: ArtifactHub API 2026-06-12 (3 operator 등록 / commons 404), 로컬 gpg keyring
  `sec rsa4096 2026-06-10 F1A68935...` (expires 2028-06-09)
