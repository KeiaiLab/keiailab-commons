# ADR-0018: GitOps overlay + ArtifactHub 검증 파이프라인 표준화

- Date: 2026-06-02
- Status: Accepted
- Authors: @phil

## Context

keiailab operator 4종(mongodb-operator / postgres-operator / valkey-operator + operator-commons
라이브러리)의 cross-repo 표준이 비일치 상태였다:

- **GitOps overlay 경로 drift**: operator 3종 간에도 경로 표준 불일치.
- **ArtifactHub 검증 자동화 부재**: ArtifactHub Signed badge 전제 조건인 PGP
  signingKey(`89A409476828CB992338C378651E51AF520BCB78`) 메타데이터 검증이 수동에만
  의존.
- **operator-commons의 특수성**: `type: library`로 Helm chart를 ArtifactHub에 publish는
  하지만 operator처럼 클러스터에 직접 배포하지 않는다.

operator-commons는 operator 3종이 `go.mod`로 import하는 공유 라이브러리다. 따라서
GitOps overlay(kustomize `deploy/overlays/prod/`)는 적용 대상이 아니며, ArtifactHub
publish(Layer 1)만 해당된다.

ArtifactHub repositoryID는 ArtifactHub에 신규 `operator-commons` repository 등록 후
발급되는 값으로, 등록 완료 시점에 `charts/artifacthub-repo.yml`에 확정 기입한다.

## Decision

**2-레이어 분리**를 전 4종에 적용하되 operator-commons는 **Layer 1만** 해당:

- **Layer 1 — ArtifactHub publish** (4종 모두 해당): helm chart(`charts/operator-commons/`)
  → gh-pages → ArtifactHub Signed badge. 공통 PGP signingKey fingerprint
  `89A409476828CB992338C378651E51AF520BCB78`를 `charts/artifacthub-repo.yml`에 등록.
- **Layer 2 — GitOps 배포 overlay** (operator 3종만): operator-commons는 `type: library`로
  배포 대상이 아니므로 **Layer 2에서 명시적으로 제외**.

**ArtifactHub 검증 파이프라인**:
- `.github/workflows/artifacthub-verify.yml`: `ah lint`(메타데이터 린트) + smoke 테스트
  (gh-pages 인덱싱 확인 + ArtifactHub REST 등록 확인 + `.tgz.prov` 도달성 검증).

**서명 구분**:
- `charts/artifacthub-repo.yml` PGP signingKey → ArtifactHub `Signed` badge.
- cosign(`release.yml`) → GitHub Release `Verified` 레이블.
- 두 서명은 **완전히 별개**다 — 혼동 금지.

**operator-commons 특이사항**:
- `charts/artifacthub-repo.yml`에 `repositoryID`는 ArtifactHub에 신규
  `operator-commons` repository 등록 완료 후 확정. 등록 전까지는 placeholder 상태.
- `charts/operator-commons/Chart.yaml`에 `type: library` 명시로 배포 대상이 아님을
  ArtifactHub에도 전달.
- `artifacthub-verify.yml`이 library chart에 대해서도 `ah lint` + smoke를 수행하여
  메타데이터 품질 보증.

**전파 방식**: Approach A(self-contained) — valkey reference를 각 repo에 복사+적응.
org-level reusable workflow(`uses:`) 방식은 배제. 이유: OSS fork 가능성 +
`keiailab/.github` org repo 2026-05-27 제거됨.

**GH Actions 사용 정당화**: RFC-0002(GitHub Actions 영구 금지)는 GitLab/인프라
closed-source org billing SPOF(2026-04-28 트리거) 컨텍스트의 결정이다. 본 대상은
**GitHub OSS public repo** + **사용자 명시 지시**("GHActions 통해서 artifacthub.io
파이프라인 검증"). 거버넌스 우선순위(사용자 명시 > Tier-1 글로벌)상 OSS public repo의
GH Actions 사용은 정책 위반이 아니다. ADR-0012(GHA block hook)는 GitLab repo 대상이며
본 GitHub OSS repo에는 미적용.

## Consequences

**긍정적**:
- operator-commons의 ArtifactHub publish 검증 자동화 확보.
- library chart임을 명시적으로 표준에 기록 — 다른 operator가 잘못된 GitOps overlay
  적용을 시도하는 혼동 방지.
- `ah lint` + smoke로 `type: library` chart도 메타데이터 품질 보증.
- repositoryID 확정 전 placeholder 정책 명시로 향후 등록 시 일관된 절차 제공.

**부정적 / 트레이드오프**:
- ArtifactHub repositoryID가 등록 완료 전까지 미확정 — REST smoke에서 repositoryID
  기반 확인이 등록 후에야 가능.
- `.tgz.prov` 생성은 현재 로컬 `scripts/release.sh --sign`에서만 동작 — CI 자동화는
  GPG private key secret 결정 후 후속 적용.
- 4종 각각 `artifacthub-verify.yml` 유지 필요(self-contained overhead).

## Alternatives Considered

**org-level reusable workflow(`uses:` 호출)**: 배제. `keiailab/.github` org repo가
2026-05-27 제거됨. OSS repo는 self-contained를 선호(fork 시 의존성 없음). valkey
ADR-0024가 이미 self-contained manual pattern 확립.

**Layer 2(GitOps overlay) operator-commons 포함**: 배제. `type: library`는 클러스터에
직접 배포하는 workload가 아니다. operator 3종이 `go.mod`로 import하는 라이브러리에
kustomize overlay를 두는 것은 의미 없으며 혼동을 초래한다.

**GH Actions 완전 배제(로컬 4계층만)**: ArtifactHub smoke는 gh-pages publish 후
원격 상태를 확인해야 하므로 로컬에서 실행 불가. ADR-0012가 GitLab repo 대상임을 명확히
하고 GitHub OSS repo에서는 publish 파이프라인 GH Actions를 유지한다.

## Refs

- ADR-0012: GitHub Actions block hook (GitLab repo 대상)
- ADR-0014: `scripts/release.sh` — manual library release pipeline
- valkey-operator ADR-0024: Helm chart manual pattern + ArtifactHub (reference 구현)
- valkey-operator ADR-0044: ArtifactHub Signed + Official trust badges
- RFC-0002: GitHub Actions 영구 금지(GitLab/인프라 closed-source 한정)
