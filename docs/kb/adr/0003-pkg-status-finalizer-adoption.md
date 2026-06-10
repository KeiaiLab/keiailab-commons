# ADR-0003: pkg/status + pkg/finalizer 표준 채택 — pkg/status 슈가 추가 + pkg/finalizer 변경 없음

- Date: 2026-05-09
- Status: Accepted
- Authors: @eightynine01
- Refs: pkg/status + pkg/finalizer 표준

## Context

keiailab-commons v0.5.0 시점에 `pkg/finalizer` 와 `pkg/status` 채택률이
0% 임이 측정됨. downstream consumers 가
각자 finalizer slice 직접 조작 + ConditionType/Reason 분기 구현.

본 ADR 은 pkg/status + pkg/finalizer 표준 의 commons-side implementation 결정을 보존한다.
Consumer migration 자체는 각 repo 별 신규 ADR 에서 추적 (downstream ADR 갱신).

## Decision

1. **`pkg/status` 에 두 함수 추가** (변경 가능):
   - `SetAvailable(conds, status, reason, message, observedGeneration)` —
     Endpoint 보유 + readiness probe 통과 후 호출. `Ready` 보다 약한 신호.
   - `SetReadyFalse(conds, reason, message, observedGeneration)` —
     `SetReady(_, ConditionFalse, _, _, _)` 의 가독성 슈가.

2. **`pkg/finalizer` 는 *변경 없음*** —
   기존 `Add/Remove/Has` + `Prefix` 만 pkg/status + pkg/finalizer 표준 §3.2 의 표준 API 로 *명시*.
   Consumer migration 은 controller-runtime 의 `controllerutil.AddFinalizer`
   에서 본 패키지 `Add` 로 import path 만 교체하는 형태.

3. **`EnsureRemoval` 헬퍼 신설 보류** —
   controller-runtime `client.Client` 의존성 도입은
   `pkg/finalizer/finalizer.go:7-10` docstring 의 *미의존 원칙* 위반.
   별도 sub-package `pkg/finalizer/runtime` 으로 분리 가능하나, 가치
   대비 비용 (테스트 + 의존성 표면 확장) 이 큼. 호출자는
   `if Remove(obj, name) { c.Update(ctx, obj) }` 두 줄 패턴 사용.

4. **Release**: commons v0.6.0 으로 bump (semver minor — API 추가).

## Consequences

### Positive

- downstream operator 가 v0.6.0 import 후 `SetAvailable` 사용 가능 — Phase 보고
  표준화 (Endpoint 가용성 신호 분리).
- `SetReadyFalse` 가 가장 빈번한 reconcile 실패 호출 패턴의 가독성
  향상 — 5 줄 → 2 줄.
- `pkg/finalizer` API 표면 변동 없음 — go.mod consumer 의 v0.5.0 →
  v0.6.0 무영향.

### Negative

- API 표면 +2 함수 — downstream consumer 통일 가치가 슈가 비용 정당화.
- `EnsureRemoval` 부재로 인해 호출자 측 boilerplate 2 줄 유지 — 단,
  controller-runtime 미의존 원칙 보존이 더 큰 가치.

### Trade-offs

- *commons API 추가 최소화* (본 ADR) vs *호출자 편의 극대화* (`EnsureRemoval`
  추가, 거부됨) — `pkg/finalizer` 의 zero-dep 원칙이 downstream consumer cross-cut
  무게에서 우위.
- *generic 4종 ConditionType 만 표준화* (pkg/status + pkg/finalizer 표준 §3.3) vs *전체
  도메인 ConditionType 표준화* — 후자는 무관 type 의 strawman 노출.

## Alternatives Considered

1. **`EnsureRemoval(ctx, c, obj, name) (bool, error)` 추가** — 거부.
   - controller-runtime `client.Client` 의존 → 미의존 원칙 위반.
   - sub-package 분리 비용 > 두 줄 호출자 boilerplate 비용.
2. **`SetReady` 시그니처 변경 (variadic options)** — 거부.
   - 기존 호출자 모두 깨짐 (downstream operator 가 v0.6.0 bump 시 일괄 수정 필요).
   - 슈가 2종 추가가 동일 가치를 호환 보존하며 제공.
3. **`pkg/condition` 신규 패키지로 분리** — 거부.
   - `pkg/status` 가 이미 적합한 이름. 분리 시 import 표면 증가.

## Refs

- pkg/status + pkg/finalizer 표준
- 사례: `pkg/status/conditions.go` (변경 후), `pkg/finalizer/finalizer.go` (변경 없음).
