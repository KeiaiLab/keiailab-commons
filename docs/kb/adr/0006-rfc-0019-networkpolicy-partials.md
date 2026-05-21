# ADR-0006: Helm library chart 정책 §3.2 채택 — keiailab.networkpolicy.{dataplane,controlplane} partials

- Date: 2026-05-09
- Status: Accepted (PR-B6 — chart v0.2.0)
- Authors: @eightynine01
- Refs: Helm library chart 정책 §3.2, ADR-0005 (Helm library chart 정책 §3.1 채택)

## Context

ADR-0005 (PR-B2) 가 `keiailab-commons` library chart 신설 + §3.1
(commonLabels + ServiceMonitor) implementation 완료. 본 ADR 은 §3.2
(NetworkPolicy partials) implementation.

downstream operator 와 downstream operator 가 *서로 다른 두 NetworkPolicy
패턴* 을 보유:

| 패턴 | 출처 | 용도 |
|---|---|---|
| dataplane (default-deny + allow-internal-instance) | downstream operator/charts/downstream operator/templates/networkpolicy.yaml | managed workload (DB pod) 자체 보호 |
| controlplane (manager 자체 + metrics/webhook + API/DNS) | downstream operator/charts/downstream operator/templates/networkpolicy.yaml | operator manager pod 자체 보호 |

두 패턴이 *서로 다른 추상* — single partial 통합 대신 *two partials*.

## Decision

1. **`charts/keiailab-commons/templates/_networkpolicy.tpl` 신규**:
   - `keiailab.networkpolicy.dataplane` — postgres 패턴 추출. caller 가
     fullname/labels/managedBy/port 를 dict 인자로 전달.
   - `keiailab.networkpolicy.controlplane` — valkey 패턴 추출. caller 가
     fullname/labels/selectorLabels/metricsPort/webhookEnabled/webhookPort
     + additionalIngress/Egress 를 dict 로 전달.

2. **두 patterns 분리 보존**: dataplane 은 *managed-by label 기반 selector*,
   controlplane 은 *selectorLabels (operator 자체)*. 하나의 partial 로
   통합 시 *추상 모호*.

3. **chart v0.1.0 → v0.2.0** bump (semver minor — partial 추가).

4. **PR 분할**:
   - PR-B6 (본 PR): commons partials + chart v0.2.0.
   - PR-B7 (mongodb): networkpolicy.yaml 신규 + dataplane partial include
     (mongodb 가 chart 에 networkpolicy 미보유 — 본 PR-B6 후 추가).
   - PR-B8 (postgres / valkey): 기존 networkpolicy.yaml → partial include
     로 교체.

## Consequences

### Positive

- mongodb 의 networkpolicy 미존재 차단 — PR-B7 가 commons partial 로
  단순 추가 가능.
- downstream consumer NetworkPolicy 정합성 — 동일 default-deny + allow-internal
  semantics.
- chart v0.2.0 backward 호환 — v0.1.0 사용자 영향 없음 (partial 추가만).

### Negative

- chart 표면 +2 partial. consumer chart 가 *dict 인자* 패턴 학습 의무.
  단 §3.1 (ServiceMonitor) 와 동일 패턴 — 추가 비용 미미.
- *두 patterns 분리* 로 인해 consumer chart 가 *어떤 partial 사용?* 의사
  결정. 단 README 표 + 두 출처 (postgres/valkey) 명시로 가이드.

### Trade-offs

- *두 partials* (본 ADR) vs *single unified partial* — 후자는 *추상 모호*.
  본 ADR 의 분리가 *명시 의미*.

## Alternatives Considered

1. **single partial + scope enum** (`scope: dataplane|controlplane`) — 거부.
   - dict 인자 분기 복잡. 두 partials 가 *명시적*.

2. **partial 미추출, 각 chart 자체 유지** — 거부.
   - downstream consumer cross-cut drift 위험 (tooling unification 정책 §3.3 lint 위반).
   - mongodb chart 의 networkpolicy 부재 차단 path 부재.

## Refs

- Helm library chart 정책 §3.2.
- ADR-0005 (commons): §3.1 implementation.
- downstream operator/charts/downstream operator/templates/networkpolicy.yaml (dataplane reference).
- downstream operator/charts/downstream operator/templates/networkpolicy.yaml (controlplane reference).
- 구현 결정.
- 후속 PR-B7 (mongodb networkpolicy.yaml 신규) / PR-B8 (postgres+valkey 기존 → partial include).
