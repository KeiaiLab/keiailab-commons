# ADR-0016: `pkg/pvc` + `pkg/topology` 신규 추출

- Date: 2026-05-21
- Status: Accepted
- Authors: @eightynine01

## Context

downstream Kubernetes operator 들이 `internal/controller/pvc_resize.go`
형태로 ~120 LOC 의 거의 동일한 PVC expansion 코드를 각자 보유하고
있었습니다. 본문 차이는 호출자 시그니처의 다중 StatefulSet 지원 여부 등
사소한 부분이었으며, ~360 LOC 의 cross-downstream 중복으로 분류됩니다.

동시에 `defaultedTopologySpread` 의 인라인 동등 코드 ~135 LOC 가 별도로
중복되었습니다. 임계값 (`replicas >= 1` vs `members >= 2`) 만 다를 뿐
본질은 동일합니다.

총 ~495 LOC 의 cross-downstream 잡음 — *최우선 추출* 대상으로 분류했으며,
추출 후 downstream 측 LOC 는 commons import 두 줄 + 호출 한 줄로
압축됩니다.

## Decision

1. **`pkg/pvc` 신규** (Stability: Beta).
   - `ExpandDataPVCs(ctx, client.Client, namespace, stsNames, desired, ...Option)`.
   - 함수형 옵션: `WithVCTName(string)` — 기본 `"data"`.
   - 헬퍼: `PVCNamePrefix(vctName, stsName) string`.
   - 의존성 정책: commons 의 *순수 데이터 변환 패키지* (finalizer / status /
     security) 와 달리 본 패키지는 List + Patch + Get 워크플로우를
     *imperative* 하게 수행해야 함 → 단 한 곳에서만
     `sigs.k8s.io/controller-runtime/pkg/client` 의존. downstream 이 이미
     controller-runtime 보유 → 신규 의존 부담 0.
   - 로깅은 `controller-runtime/pkg/log.FromContext` — ctx 주입.

2. **`pkg/topology` 신규** (Stability: Beta).
   - `Defaulted(user, replicas, selector, ...Option) []corev1.TopologySpreadConstraint`.
   - 함수형 옵션:
     - `WithMinReplicas(int32)` — 기본 2.
     - `WithTopologyKeys(...string)` — 기본 `[zone, hostname]`.
     - `WithMaxSkew(int32)` — 기본 1.
     - `WithWhenUnsatisfiable(action)` — 기본 `ScheduleAnyway`.
   - 의존성 정책: 순수 데이터 변환 → controller-runtime 미의존.
     `apimachinery` + `api/core/v1` 만.

3. **go.mod 변경 (commons)**.
   - `sigs.k8s.io/controller-runtime v0.24.1` 신규.
   - 전이 의존으로 `k8s.io/api v0.35.0 → v0.36.1`, `k8s.io/apimachinery`
     동일 bump.
   - `go 1.25.7 → 1.26.0` (k8s.io 0.36.0+ 요구).

4. **Release**: commons v0.9.0 (semver minor — API 추가만, breaking 없음).

5. **Consumer migration**: 별 ADR 으로 downstream 측에서 추적. 본 ADR 은
   commons-side 결정만 보존.

## Consequences

### Positive

- downstream 측 ~495 LOC 일괄 삭제 가능 — migration 머지 후.
- 단일 SSOT: PVC expansion / TSC default 로직이 commons 1 곳에서만 갱신됨.
- 함수형 옵션 패턴으로 미래 확장 (예: VCT 이름 변경, 추가 topology axis)
  의 API breaking 위험 0.
- 28 단위 테스트 (fake K8s client) — 회귀 안전망.

### Negative

- commons `go.mod` 표면 +1 (controller-runtime). 단, downstream 이 이미
  보유 → 사실상 zero-cost.
- `pkg/pvc` 가 controller-runtime 의존 → commons 의 "controller-runtime
  미의존" 원칙 *부분 완화*. 본 ADR 은 그 완화를 *예외적 단일 패키지*
  에 한정하고 명시적으로 사유 (imperative API 호출 필요) 를 기록.
- Beta tier 로 1 cycle 정착 필요 — downstream 회귀 통과 후 Stable 격상.

### Trade-offs

- *controller-runtime 의존 도입* (본 ADR) vs *자체 Client 인터페이스 정의*
  (초기 시도, 거부됨) — Patch 인터페이스의 `Data(client.Object) ([]byte,
  error)` 시그니처가 `client.Object` 자체에 의존하기 때문에 진짜 zero-dep
  구현은 불가능하거나 (모든 타입을 `any` 로 노출하면) 호출자 어댑터 비용이
  너무 큽니다. controller-runtime 직접 의존이 *명료성* 우위.
- *downstream 일괄 적용* vs *순차 적용* — 본 ADR 은 후자 채택. 각 downstream
  의 e2e 영향이 큰 controller 변경이라 *위험 작업* 분류 (사용자 머지
  검토가 필요).

## Alternatives Considered

1. **자체 Client 구조적 인터페이스 정의** — 거부.
   - Patch / Object / ObjectList 의 메서드 체인이 controller-runtime 의
     `client.Object` 에 의존 → 진짜 zero-dep 불가.
   - 어댑터 비용 > controller-runtime 직접 의존 비용.

2. **단일 downstream cherry-pick PoC** — 거부.
   - 본 케이스는 *이미 다중 downstream 이 동일 코드* 라 단일 PoC 무의미.
   - commons 변경 후 downstream 측 동시 PR 의 atomic migration 패턴이 더
     적합.

3. **`pkg/sts` 로 통합 (PVC + topology 동시 처리)** — 거부.
   - 두 기능의 *호출 주기* 가 다름: TSC default 는 STS 빌드 시 매번,
     PVC expansion 은 reconcile 시 size diff 발견 시만.
   - 패키지 분리가 *호출자 mental model* 에도 더 명확.

## Refs

- 사례:
  - `pkg/pvc/expansion.go`, `pkg/pvc/expansion_test.go`, `pkg/pvc/doc.go`.
  - `pkg/topology/spread.go`, `pkg/topology/spread_test.go`,
    `pkg/topology/doc.go`.
- 양식: Nygard 5 섹션.
