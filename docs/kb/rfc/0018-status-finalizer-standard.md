# RFC-0018: pkg/status + pkg/finalizer Adoption Standard

- Date: 2026-05-09
- Status: Proposed
- Authors: @eightynine01
- Affected Repos: `operator-commons`, `valkey-operator`, `mongodb-operator`, `postgres-operator`
- Implementation ADR: ADR-0003 (commons), repo 별 신규 ADR (mongodb/valkey/postgres)

## Motivation

operator-commons v0.5.0 시점 측정 결과:

| 패키지 | 채택률 | 비고 |
|---|---|---|
| `pkg/security` | 100% | 3 operator 모두 import |
| `pkg/labels` | 100% | 동일 |
| `pkg/version` | 67% | postgres 가 자체 `internal/version/matrix.go` 사용 (RFC-0019 §3.3 별도 처리) |
| `pkg/webhook` | 67% | mongodb 채택 보류 |
| `pkg/networkpolicy` | 67% | postgres 가 chart 측 미보호 |
| `pkg/monitoring` | 33% | mongodb/postgres 가 helm partial 미사용 |
| **`pkg/finalizer`** | **0%** | 3 operator 가 `controllerutil.AddFinalizer/RemoveFinalizer` 또는 자체 finalizers.go 사용 |
| **`pkg/status`** | **0%** | 3 operator 가 *각자* `setCondition` / `meta.SetStatusCondition` 직접 호출 |

핵심 문제:
- finalizer 키 prefix 가 repo 별로 분기 (`mongodb.keiailab.com/`, `cache.keiailab.io/`, `postgresql.keiailab.com/`).
- ConditionType / Reason 카탈로그가 repo 별로 분기 — `kubectl describe` 출력의 4-repo 정합 부재 (운영자가 도메인 별로 다른 키워드 학습 필요).
- 동일 reconcile 실패가 repo 별로 다른 reason (`ReconcileFailed`, `ReconcileError`, `Failed`) 으로 표현.

본 RFC 는 *commons API 추가 최소화 + consumer migration 우선* 전략을 채택한다.

## Design

### §3.1 pkg/status — Adoption Contract

기존 API (변경 없음):

| 함수 | 사용 시점 |
|---|---|
| `SetReady(conds, status, reason, message, observedGeneration)` | reconcile 종료 시점 — Ready=True/False |
| `SetProgressing(conds, reason, message, gen)` | reconcile 진입 시점 — 작업 중 신호 |
| `SetDegraded(conds, reason, message, gen)` | 부분 장애 — 일부 기능 작동 |
| `IsReady(conds)`, `FindCondition`, `RemoveCondition` | read-side helpers |

신규 슈가 (본 RFC 로 추가):

| 함수 | 사용 시점 |
|---|---|
| `SetAvailable(conds, status, reason, message, gen)` | Endpoint 보유 + readiness probe 통과 후. Ready 보다 약한 신호. reconcile 진입 시 Ready=False 로 전환되어도 Available 유지 가능. |
| `SetReadyFalse(conds, reason, message, gen)` | 빈번한 reconcile 실패 호출의 가독성 슈가 (`SetReady(_, ConditionFalse, _, _, _)` 와 동일). |

표준 Reason 카탈로그 (이미 정의됨, *호출자가 사용 의무*):
- `ReasonReconciling`, `ReasonAvailable`, `ReasonNotApplicable`
- `ReasonReconcileError`, `ReasonExternalDepBlocked`, `ReasonValidationFailed`

### §3.2 pkg/finalizer — Adoption Contract

기존 API (변경 없음, 본 RFC 로 *명시 채택*):

| 함수 | 의미 |
|---|---|
| `Add(obj, name) bool` | finalizer slice 에 name 추가, idempotent |
| `Remove(obj, name) bool` | name 제거, idempotent |
| `Has(obj, name) bool` | 조회 |
| `Prefix` | `keiailab.com/` 표준 prefix |

설계 원칙 (본 RFC 로 *고정*):
- `pkg/finalizer` 는 *controller-runtime 미의존* — `client.Update` 은 호출자 책임.
- 호출 패턴: (1) `Add` → `client.Update` (2) cleanup → `Remove` → `client.Update`.
- 호출자는 `controllerutil.AddFinalizer` 대신 본 패키지 `Add` 를 사용 — 동작 동등이나 *commons 경유 의존성 추적* 강제.

거부된 대안: `EnsureRemoval(ctx, c, obj, name)` — controller-runtime 의존성 도입. ADR-0003 §Alternatives 에 근거 보존.

### §3.3 도메인 특이 ConditionType / Reason 보존

각 repo 의 도메인 특이 ConditionType 은 *유지*:

| Repo | 도메인 특이 ConditionType (예) |
|---|---|
| valkey-operator | `ShardReady`, `ResardingProgress` |
| mongodb-operator | `PrimaryUnreachable`, `ScalePolicyDeliberateFalse` |
| postgres-operator | `ShardsReady`, `RouterReady`, `BackupHealthy`, `AutoSplitEligible` |

commons 는 *generic 4종* (`Ready/Progressing/Degraded/Available`) + 6 Reason 만 강제.

## Migration Plan

### 단계 1 — commons v0.6.0 (PR-A1)

- `pkg/status` 에 `SetAvailable` + `SetReadyFalse` 추가 (본 RFC 의 Implementation ADR).
- 본 RFC 본문 작성.
- 후속 v0.7.0 까지 *기존 호출자 무영향*.

### 단계 2 — 3 operator migration (PR-A5/A6/A7)

| Repo | Migration 변경 | ADR |
|---|---|---|
| valkey-operator (PR-A6) | `controllerutil.AddFinalizer/RemoveFinalizer` → `finalizer.Add/Remove` (4 controller). `setCondition` 호출 시 `status.ReasonReconcileError` 등 commons reason 사용. | 신규 |
| mongodb-operator (PR-A5) | 동일 패턴 (3 controller). | 신규 |
| postgres-operator (PR-A7) | `pkg/status` 만 채택. `pkg/finalizer` 는 ADR-0008 (cascade-delete-by-OwnerReference) 비대칭 보존. | ADR-0008 갱신 |

각 PR 의 *호환성 보장*:
- `controllerutil.AddFinalizer` 와 본 패키지 `Add` 모두 idempotent — 동시 사용 단계에서도 finalizer slice 정합.
- `meta.SetStatusCondition` 와 commons `SetReady` 모두 동일 apimeta 호출 — wire-level 차이 없음.

### 단계 3 — Reason 키 외부 contract 변경 release note

기존 외부 contract:
- valkey 의 `ReasonReconcileFailed` → commons `ReasonReconcileError` 변경.
- mongodb / postgres 도 정합화.

영향: `kubectl jsonpath` 또는 alert rule 에서 `Reason="ReconcileFailed"` 로 매칭하는 사용자.
완화: 1 release window (commons v0.6.0 → v0.7.0) 에서 양쪽 모두 동작 가능 (apimeta dedup).
release note 명시: "Conditions[].reason 카탈로그가 RFC-0018 표준으로 통일됨 — alert rule 갱신 필요".

## Rollout

| 시점 | 작업 |
|---|---|
| commons v0.6.0 | 본 RFC + ADR-0003 + 슈가 2종 추가 |
| valkey/mongodb v1.x | go.mod commons v0.6.0 bump + finalizer + status migration |
| postgres v0.x | status 만 migration (finalizer 비대칭) |
| commons v0.7.0 | 만약 도메인 reason 통합 필요 시 추가 reason enum |

전체 timeline: Sprint A (PR-A1~A7) 내 완료 — 1-2 개발 cycle.

## Alternatives Considered

1. **`EnsureRemoval(ctx, c, obj, name) (bool, error)` 신설** — 거부.
   - controller-runtime `client.Client` 의존성 도입 → 기존 미의존 원칙 (`pkg/finalizer/finalizer.go:7-10` docstring) 위반.
   - 별도 sub-package `pkg/finalizer/runtime` 분리 가능하나, 가치 < 비용.
   - 호출자 코드는 `Remove(obj, name); c.Update(ctx, obj)` 두 줄로 충분.

2. **commons 가 ConditionType / Reason 을 *전부* 표준화** — 거부.
   - 도메인 특이 type (`ShardsReady` 등) 흡수 시 반대편 operator 에 무의미한 type 노출.
   - 표준화 단위는 *generic 4종* 까지.

3. **postgres 도 finalizer 채택 강제** — 거부.
   - postgres ADR-0008 의 cascade-delete-by-OwnerReference 가 *의도된 비대칭*.
   - BackupCleanupJob CRD 가 외부 자원 cleanup 분리 처리.

## Status

- 2026-05-09: Proposed (본 RFC 본문 작성, commons v0.6.0 PR-A1 제출).
- 다음 단계: 3 operator migration PR (PR-A5/A6/A7) 가 본 RFC 를 *Refs* 로 인용. 모두 머지 시 Status: Implemented 로 갱신.

## Refs

- ADR-0003 (commons): pkg/status 슈가 추가 결정 — 본 RFC 의 Implementation.
- ADR-0008 (postgres): cascade-delete-by-OwnerReference 비대칭 결정 — 본 RFC 의 §3.2 Migration 단계 2 가 인용.
- Plan: `~/.claude/plans/1-https-artifacthub-io-packages-helm-clo-synthetic-gem.md` §2 D10 / D11.
- 글로벌 표준: `standards/adr.md §6` (RFC vs ADR 구분).
