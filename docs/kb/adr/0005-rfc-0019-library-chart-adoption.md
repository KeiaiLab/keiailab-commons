# ADR-0005: Helm library chart 정책 §3.1 채택 — keiailab-commons Helm library chart 신설

- Date: 2026-05-09
- Status: Accepted (PR-B2 — §3.1 commonLabels + ServiceMonitor)
- Authors: @eightynine01
- Refs: Helm library chart 정책 (`docs/kb/rfc/0019-helm-partials-library.md`)

## Context

4-repo (mongodb / postgres / valkey operator) Helm chart 의 *공통 partial*
중복 (~230 LOC, Helm library chart 정책 §Motivation). Helm library chart 패키징 표준으로
*single source* 로 통합.

본 ADR 은 Helm library chart 정책 §3.1 (commonLabels + ServiceMonitor) implementation
을 보존한다. §3.2 (NetworkPolicy) / §3.4 (PSS) / §3.5 (RBAC) 는 후속 PR.

## Decision

1. **`charts/keiailab-commons/` 신규 디렉토리** — Helm library chart
   (type: library).

2. **`Chart.yaml`**:
   - `apiVersion: v2`, `type: library`, `version: 0.1.0`, `appVersion: "0.1.0"`.
   - `kubeVersion: ">=1.25.0-0"` (operator support range 동등).
   - `keywords: [keiailab, library, operator-helpers, servicemonitor, rbac]`.

3. **`templates/_helpers.tpl`** — partial 정의:
   - `keiailab.commonLabels` (Kubernetes Recommended Labels).
   - `keiailab.observability.serviceMonitor` (Prometheus Operator
     ServiceMonitor 공통 spec).

4. **dict 인자 패턴**: caller 가 `fullname` / `labels` / `selectorLabels`
   를 dict 로 전달. partial 이 chart 별 helper 이름 무관하게 동작.

5. **OCI publish 분리**: 본 PR 은 *chart 신설 만*. `helm package` +
   `helm push oci://ghcr.io/keiailab/charts/keiailab-commons:0.1.0` 는
   *별 PR (PR-B2.2)* — release 자동화 분리.

6. **§3.2 / §3.4 / §3.5 후속**: NetworkPolicy / PSS / RBAC partial 은
   별 PR (PR-B4 / B5 / C1).

## Consequences

### Positive

- 4-repo cross-cut partial 중복 제거 path 구축. PR-B4/B5 에서 mongodb /
  valkey 의 servicemonitor.yaml 이 partial include 로 전환 가능.
- *single source* 로 정정 — kubectl 검색 정합성 강제.
- Helm 표준 library chart — ArgoCD GitOps 흐름과 정합.

### Negative

- consumer chart 가 *dependency* 추가 — `helm dependency update` 의무.
  CI 흐름 + ArgoCD 와의 정합 검증 필요 (PR-B4 의 검증 단계).
- OCI publish (PR-B2.2) 까지 *consumer 채택 불가* — 본 PR 만으로는 가치
  부재. 단 chart skeleton 가치 보존.

### Trade-offs

- *library chart* (본 ADR) vs *go template 함수* / *Kustomize component* —
  Helm 표준 + consumer chart 모두 Helm 기반 — library chart 가 정합.

## Alternatives Considered

1. **subchart (type: application)** — 거부.
   - subchart 는 *리소스 생성* 용 — 본 chart 는 *helper only*. library
     type 이 정확.

2. **각 chart 의 _helpers.tpl 에 cp 유지** — 거부.
   - 4-repo drift 위험. tooling unification 정책 §3.3 lint 위반.

3. **public Helm repo (Bitnami chart-base 등) 사용** — 거부.
   - keiailab 도메인 표준 (commonLabels 의 prefix 등) 외부 chart 와
     불일치.

## Refs

- Helm library chart 정책 §3.1: commonLabels + ServiceMonitor partial 표준.
- Plan §2 D14 (Sprint B PR-B2).
- 후속 OCI publish: PR-B2.2 (별 PR).
- Helm library chart: <https://helm.sh/docs/topics/library_charts/>
