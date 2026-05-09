# RFC-0019: Helm Partials Library Chart 표준

- Date: 2026-05-09
- Status: Proposed (§3.1 ServiceMonitor + commonLabels Implemented in PR-B2)
- Authors: @eightynine01
- Affected Repos: `operator-commons`, `valkey-operator`, `mongodb-operator`, `postgres-operator`
- Implementation ADR: ADR-0005 (commons), 후속 repo 별 ADR

## Motivation

4-repo (mongodb / postgres / valkey operator) 의 Helm chart 가 *동일
구조의 partial* 을 각자 정의:

| Partial | mongodb | postgres | valkey | 중복 LOC |
|---|---|---|---|---|
| ServiceMonitor (`metrics.enabled` + `serviceMonitor.enabled` 분기) | 36 | 0 (chart 부재) | 36 | 72 |
| NetworkPolicy (deny-by-default + 같은 instance 허용) | 0 (chart 부재) | 부분 | 부분 | ~120 |
| commonLabels (Kubernetes Recommended Labels) | 7 | 7 | 7 | 21 |
| RBAC base (leader-election + events) | 6 | 6 | 6 | 18 |

총 ~230 LOC 중복. *Helm library chart* 패키징으로 *single source* 로 통합.

## Design

### §3.1 commonLabels + observability.serviceMonitor (PR-B2 implementation)

**Helper**:
- `keiailab.commonLabels` — Kubernetes Recommended Labels (`managed-by`,
  `instance`, `helm.sh/chart`, `version`).
- `keiailab.observability.serviceMonitor` — ServiceMonitor 공통 spec.

**Caller 패턴** (dict 인자):
```yaml
{{- include "keiailab.observability.serviceMonitor" (dict
    "ctx" .
    "fullname" (include "<chart>.fullname" .)
    "labels" (include "<chart>.labels" .)
    "selectorLabels" (include "<chart>.selectorLabels" .)) }}
```

caller 가 *자체 fullname/labels/selectorLabels* helper 를 dict 로 전달
— partial 이 chart 별 helper 이름 무관하게 동작.

### §3.2 NetworkPolicy partials (PR-B4 후속)

- `keiailab.networkpolicy.dataplane` — postgres 패턴 (default-deny +
  allow-internal-instance).
- `keiailab.networkpolicy.controlplane` — valkey 패턴 (manager NP, webhook
  + metrics 분기).

### §3.3 (예약) — pkg/version Generic Matrix (commons-ADR-0004 implementation)

본 RFC 외 — postgres PR-B3 에서 commons `Matrix[E]` 위임 완료.

### §3.4 PodSecurity Restricted partial (PR-B5 후속)

- `keiailab.security.podSecurityRestricted` — PSA "restricted" profile
  SecurityContext 표준.

### §3.5 RBAC partial (PR-C1 후속)

- `keiailab.rbac.serviceAccount` — SA + ImagePullSecrets + 자동 token 비활성.
- `keiailab.rbac.controllerBase` — leader-election + events + service watch.
- `keiailab.rbac.workloadBase` — statefulsets/deployments + services/configmaps/secrets.

## Migration Plan

| 단계 | PR | 산출물 |
|---|---|---|
| commons §3.1 implementation | PR-B2 (본 PR) | `charts/keiailab-commons/` library chart + commonLabels + ServiceMonitor |
| commons §3.1 OCI publish | PR-B2.2 (별 PR) | `helm package` + `helm push oci://ghcr.io/keiailab/charts` 자동화 |
| consumer §3.1 채택 | PR-B4 (mongodb) / PR-B5 (valkey) | servicemonitor.yaml → partial include |
| commons §3.2 (NetworkPolicy) | PR-B6 | dataplane / controlplane partials |
| consumer §3.2 채택 | PR-B7 (mongodb networkpolicy 신규) / PR-B8 (valkey/postgres) | networkpolicy.yaml → partial include |
| commons §3.5 (RBAC) | PR-C1 | controllerBase / workloadBase partials |
| consumer §3.5 채택 | PR-C2~C4 | RBAC partial include |

## Rollout

| 시점 | Library version |
|---|---|
| PR-B2 (본 PR) | keiailab-commons v0.1.0 — commonLabels + ServiceMonitor |
| PR-B6 | v0.2.0 — NetworkPolicy partials |
| PR-C1 | v0.3.0 — RBAC partials |

OCI publish: `oci://ghcr.io/keiailab/charts/keiailab-commons:<version>`.
consumer chart 가 `helm dependency update` 후 사용.

## Alternatives Considered

1. **각 chart 자체 partial 보존** — 거부. 4-repo cross-cut drift 위험
   (RFC-0017 §3.3 lint 위반).
2. **Helm library chart 대신 Kustomize component** — 거부. 4-repo Helm
   기반 — Kustomize 통합 비용 > 가치.
3. **gomplate 외부 도구** — 거부. Helm 표준 외 도구 — operator 표준
   GitOps 흐름 (ArgoCD) 와 정합 부족.

## Status

- 2026-05-09: Proposed (§3.1 implementation in PR-B2 — Accepted via ADR-0005).
- 후속 §3.2 / §3.4 / §3.5 implementation 완료 시 Status: Implemented.

## Refs

- ADR-0005 (commons): RFC-0019 §3.1 채택 implementation.
- commons-ADR-0004: pkg/version Matrix[E] generic — §3.3 implementation (별 흐름).
- Plan §2 D13/D14/D15.
- Helm library chart docs: <https://helm.sh/docs/topics/library_charts/>
