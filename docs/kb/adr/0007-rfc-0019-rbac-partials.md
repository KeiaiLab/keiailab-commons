# ADR-0007: RFC-0019 §3.5 채택 — keiailab.rbac.{serviceAccount,controllerBase,workloadBase} partials

- Date: 2026-05-09
- Status: Accepted (PR-C1 — chart v0.3.0)
- Authors: @eightynine01
- Refs: RFC-0019 §3.5, ADR-0005 (§3.1), ADR-0006 (§3.2), Plan §2 D15

## Context

ADR-0005/0006 (PR-B2/B6) 이 commons library chart 의 §3.1 (ServiceMonitor +
commonLabels) + §3.2 (NetworkPolicy) 구축. 본 ADR 은 §3.5 (RBAC partials).

4-repo (mongodb / postgres / valkey operator) 의 RBAC YAML 분석:
- 공통 verb set (~80 LOC): leader-election (`coordination.k8s.io/leases`),
  events (core + events.k8s.io), 자체 service watch, managed workload
  (apps/statefulsets, deployments) + dependent (core/services, configmaps,
  secrets, pvc, pods).
- delta (CRD-specific): mongodb `mongodbs/finalizers`, postgres
  `postgresclusters/status`, valkey `valkeys/finalizers` 등 — 각 chart
  자체 yaml 보존.

## Decision

1. **`charts/keiailab-commons/templates/_rbac.tpl` 신규**:
   - `keiailab.rbac.serviceAccount` — SA + ImagePullSecrets +
     automountServiceAccountToken (default true). caller 가 name /
     labels / imagePullSecrets / annotations 를 dict 로 전달.
   - `keiailab.rbac.controllerBase` — *PolicyRule list* 출력 (caller
     의 ClusterRole/Role.rules 에 nindent). leader-election + events +
     service watch.
   - `keiailab.rbac.workloadBase` — *PolicyRule list*. apps/statefulsets,
     deployments + core/services, configmaps, secrets, pvc, pods.

2. **PolicyRule list 패턴**: `serviceAccount` 는 *완전한 리소스* 출력,
   `controllerBase` / `workloadBase` 는 *rule list* 만 — caller 가
   ClusterRole/Role/RoleBinding 의 rules 아래에 include.

3. **chart v0.2.0 → v0.3.0** bump (semver minor).

4. **PR 분할**:
   - **PR-C1** (본 PR): commons partials + chart v0.3.0.
   - **PR-C2~C4** (별 PR, 3 operator): RBAC yaml 의 base verb 들이
     `controllerBase` / `workloadBase` partial include 로 교체. delta
     CRD-specific rules 만 자체 yaml 보존.

## Consequences

### Positive

- 4-repo 공통 RBAC ~80 LOC 추출 — single source. duplicate verb 의 사일런
  drift 차단.
- `controllerBase` 가 *명시적 인텐트* — 새 operator 추가 시 RBAC 표준
  채택 가능 (RFC-0017 §3.3 lint 위반 방지).
- chart v0.3.0 backward 호환 (partial 추가만).

### Negative

- consumer chart RBAC yaml 의 *부분 보존* 패턴 — base 위임 + delta 자체.
  `kubectl auth can-i` 검증 의무 (PR-C2~C4 의 검증).
- *PolicyRule list 출력* 패턴이 §3.1 (full resource 출력) 와 다름 — caller
  의 사용 패턴 학습 필요. README 표 + 사용 예 명시.

### Trade-offs

- *3 partial 분리* (serviceAccount / controllerBase / workloadBase) vs
  *single full RBAC partial* — 후자는 operator 별 ClusterRole 구조 차이
  (Role vs ClusterRole, namespaced vs cluster-scoped) 흡수 어려움. 분리가
  유연.

## Alternatives Considered

1. **single `keiailab.rbac.full` partial (전체 ClusterRole 출력)** — 거부.
   - operator 별 ClusterRole 구조 차이 흡수 어려움.
   - delta CRD-specific rules 와의 통합 복잡.

2. **PolicyRule list 대신 dict (verb→resources map)** — 거부.
   - YAML 출력 시 dict 순서 비결정적 — kubectl diff false-positive 위험.

## Refs

- RFC-0019 §3.5.
- ADR-0005 (§3.1), ADR-0006 (§3.2).
- Plan §2 D15 (Sprint C PR-C1).
- 후속 PR-C2 (mongodb) / PR-C3 (postgres) / PR-C4 (valkey): RBAC partial 채택.
