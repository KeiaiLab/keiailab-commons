# ADR-0008: Helm library chart 정책 §3.4 채택 — keiailab.security.{pod,container}SecurityContext partials

- Date: 2026-05-09
- Status: Accepted (PR-B5 — chart v0.4.0)
- Authors: @eightynine01
- Refs: Helm library chart 정책 §3.4, ADR-0005/0006/0007 (§3.1/§3.2/§3.5)

## Context

ADR-0005/0006/0007 (PR-B2/B6/C1) 가 commons library chart 의 §3.1
(commonLabels + ServiceMonitor) + §3.2 (NetworkPolicy) + §3.5 (RBAC)
구축. 본 ADR 은 *마지막* §3.4 (PodSecurity Restricted) implementation —
Helm library chart 정책 의 implementation 완결.

K8s Pod Security Standards (PSS) "restricted" profile 표준
(https://kubernetes.io/docs/concepts/security/pod-security-standards/):
- Pod-level: runAsNonRoot + seccompProfile.type=RuntimeDefault.
- Container-level: capabilities.drop=ALL + readOnlyRootFilesystem +
  allowPrivilegeEscalation=false.

downstream consumer 의 manager Deployment / DB workload StatefulSet 이 동일 표준
적용 — 본 partial 로 single source.

## Decision

1. **`charts/keiailab-commons/templates/_security.tpl` 신규**:
   - `keiailab.security.podSecurityContext` — Pod SecurityContext
     (runAsNonRoot + runAsUser/Group + fsGroup + seccompProfile).
   - `keiailab.security.containerSecurityContext` — Container
     SecurityContext (runAsNonRoot + capabilities.drop=ALL +
     readOnlyRootFilesystem + seccompProfile).

2. **`override` 인자 패턴**: caller 가 `override` dict 키로 사용자 정의
   SecurityContext 전달 시 partial 이 *그대로 출력* — v1alpha2
   PodSecurityRestricted=false 시 사용자 override 시나리오 지원
   (downstream operator ADR-0036 정합).

3. **runAsUser/Group default 65532**: distroless nonroot 표준 UID/GID.
   각 operator 의 Dockerfile 와 정합 (`USER 65532:65532`).

4. **chart v0.3.0 → v0.4.0** bump. Helm library chart 정책 implementation 완결.

5. **PR 분할**:
   - **PR-B5** (본 PR): commons partials + chart v0.4.0.
   - 후속 (별 PR): 4 repo Helm chart 의 SecurityContext 가 본 partial
     include 로 교체 — `kubectl describe deploy/...-controller-manager`
     출력 정합 검증.

## Consequences

### Positive

- downstream consumer manager pod 의 PSS Restricted 표준화 — `kubectl describe`
  + Pod Security Admission label 검증 정합.
- override 인자로 v1alpha2 의 PodSecurityRestricted=false 시나리오
  (valkey ADR-0036) 와 정합.
- Helm library chart 정책 §3 implementation 완결 — commons library chart 가 §3.1/
  §3.2/§3.4/§3.5 모두 보유.

### Negative

- chart 표면 +2 partial. consumer chart 의 *override 인자* 학습 의무.
  단 README 표 + 사용 예 명시.
- runAsUser/Group default 65532 가 *암묵적 contract* — Dockerfile 의
  `USER 65532:65532` 와 정합 강제. 다른 operator 가 다른 UID 사용 시
  override 인자로 명시.

### Trade-offs

- *두 partial 분리* (pod + container) vs *single 통합 partial* — 후자는
  Pod/Container SecurityContext 의 *상위 vs 하위* 표면 다름. 분리가
  명시적.
- *PSS Restricted 강제* (default) vs *override 우선* — `override` 인자가
  설정되면 default 무시. valkey ADR-0036 의 secure-by-default + opt-out
  path 와 정합.

## Alternatives Considered

1. **runAsUser/Group default 0 (root)** — 거부.
   - PSS Restricted 위반.

2. **single `keiailab.security.restricted` partial (Pod + Container 통합)** — 거부.
   - Pod 와 Container SecurityContext 가 *서로 다른 location* — single
     partial 로 두 곳 출력 어려움.

3. **values.yaml schema 확장 (commons 가 schema 정의)** — 거부.
   - library chart 는 *helper only* — values schema 는 consumer chart 책임.

## Refs

- Helm library chart 정책 §3.4 (이제 implementation 완료).
- ADR-0005/0006/0007: §3.1/§3.2/§3.5 implementation 선행.
- downstream operator ADR-0036 (PodSecurityRestricted Optional Toggle): 본
  partial 의 override 시나리오 부모.
- 구현 결정.
- K8s PSS: <https://kubernetes.io/docs/concepts/security/pod-security-standards/>
