# ADR-0016: Sprint 1 — pkg/pvc + pkg/topology 신규 추출 (3-operator 중복 ~495 LOC 해소)

- Date: 2026-05-21
- Status: Accepted
- Authors: @eightynine01 (Codex Major #7 — Sprint 1)
- Refs: AGENTS.md "빈번한 작업 패턴 — 새 공통 패키지 추가" (issue + 7일 윈도우 + 3 consumer 검증)

## Context

3 operator (postgres / mongodb / valkey) 의 `internal/controller/pvc_resize.go`
가 ~120 LOC 씩 거의 동일 (~360 LOC 총 중복) — 본문 차이는 `dataPVCNamePrefix`
주석 한 줄 + valkey 의 호출자 시그니처가 multi-STS 미지원이라는 점 뿐.

동시에 `defaultedTopologySpread` (postgres `internal/controller/`, mongodb
`internal/resources/`) + 인라인 valkey 동등 코드 ~135 LOC 총 중복. 임계값만
차이 (postgres `replicas >= 1`, mongodb `members >= 2`).

총 ~495 LOC 의 cross-repo 잡음. Cycle 25 spec-driven 진단에서 *최우선 추출*
대상으로 분류 — 추출 후 operator 측 LOC 가 commons import 2 줄 + 호출 1
줄로 압축됨.

## Decision

1. **`pkg/pvc` 신규** (Stability: Beta).
   - `ExpandDataPVCs(ctx, client.Client, namespace, stsNames, desired, ...Option)`.
   - 함수형 옵션: `WithVCTName(string)` — 기본 `"data"`.
   - 헬퍼: `PVCNamePrefix(vctName, stsName) string`.
   - 의존성 정책: commons 의 *순수 데이터 변환 패키지* (finalizer / status /
     security) 와 달리 본 패키지는 List + Patch + Get 워크플로우를 *imperative*
     하게 수행해야 함 → 단 한 곳에서만 `sigs.k8s.io/controller-runtime/pkg/client`
     의존. 3 consumer 가 이미 controller-runtime 보유 → 신규 의존 부담 0.
   - 로깅은 `controller-runtime/pkg/log.FromContext` — ctx 주입.

2. **`pkg/topology` 신규** (Stability: Beta).
   - `Defaulted(user, replicas, selector, ...Option) []corev1.TopologySpreadConstraint`.
   - 함수형 옵션:
     - `WithMinReplicas(int32)` — 기본 2 (mongodb/valkey 의미), postgres 는 1.
     - `WithTopologyKeys(...string)` — 기본 `[zone, hostname]`.
     - `WithMaxSkew(int32)` — 기본 1.
     - `WithWhenUnsatisfiable(action)` — 기본 ScheduleAnyway.
   - 의존성 정책: 순수 데이터 변환 → controller-runtime 미의존. apimachinery
     + api/core/v1 만.

3. **go.mod 변경 (commons)**.
   - `sigs.k8s.io/controller-runtime v0.24.1` 신규.
   - 전이 의존으로 `k8s.io/api v0.35.0 → v0.36.1`, `k8s.io/apimachinery 동일
     bump`. 3 consumer 가 이미 v0.36.1 사용 중이므로 commons 정합 회복
     (이전 의도적 strict 1.25.7 / 0.35 핀 해제).
   - `go 1.25.7 → 1.26.0` (k8s.io 0.36.0+ 요구). 3 consumer 모두 1.26.0
     이므로 정합.

4. **Release**: commons v0.9.0 (semver minor — API 추가만, breaking 없음).

5. **Consumer migration**: 별도 ADR 으로 각 operator 측에서 추적 (Sprint 1
   Phase 2). 본 ADR 은 commons-side 결정만 보존.

## Consequences

### Positive

- 3 operator 측 ~495 LOC 일괄 삭제 가능 — Phase 2 머지 후.
- 단일 SSOT: PVC expansion / TSC default 로직이 commons 1곳에서만 갱신됨.
- 함수형 옵션 패턴으로 미래 확장 (예: VCT 이름 변경, 추가 topology axis)
  의 API breaking 위험 0.
- 13 + 15 = 28 단위 테스트 (fake K8s client) — 회귀 안전망.

### Negative

- commons go.mod 표면 +1 (controller-runtime). 단, 3 consumer 가 이미
  보유 → 사실상 zero-cost.
- `pkg/pvc` 가 controller-runtime 의존 → commons 의 "controller-runtime
  미의존" 원칙 *부분 완화*. 단 본 ADR 은 그 완화를 *예외적 단일 패키지*
  에 한정하고 명시적으로 사유 (imperative API 호출 필요) 를 기록.
- Beta tier 로 1 cycle 정착 필요 — 3 operator 회귀 통과 후 Stable 격상.

### Trade-offs

- *controller-runtime 의존 도입* (본 ADR) vs *자체 Client 인터페이스 정의*
  (초기 시도, 거부됨) — Patch 인터페이스의 `Data(client.Object) ([]byte,
  error)` 시그니처가 client.Object 자체에 의존하기 때문에 진짜 zero-dep
  구현은 불가능하거나 (모든 타입을 `any` 로 노출하면) 호출자 어댑터 비용이
  너무 큼. controller-runtime 직접 의존이 *명료성* 우위.
- *3 operator 일괄 적용* vs *순차 적용* — 본 ADR 은 후자 채택 (Phase 2
  per-operator PR), 각 operator 의 e2e 영향이 큰 controller 변경이라
  *위험 작업* 분류 (사용자 머지 검토 필요).

## Alternatives Considered

1. **자체 Client 구조적 인터페이스 정의** — 거부.
   - Patch / Object / ObjectList 의 메서드 체인이 controller-runtime 의
     client.Object 에 의존 → 진짜 zero-dep 불가.
   - 어댑터 비용 > controller-runtime 직접 의존 비용.

2. **postgres-operator 측 cherry-pick PoC 1곳** — 거부.
   - AGENTS.md "빈번한 작업 패턴" 의 "consumer operator 1곳에 cherry-pick"
     단계는 *consumer-specific* 코드를 *common* 으로 일반화할 때 의미. 본
     케이스는 *이미 3 consumer 가 동일 코드* 라 단일 PoC 무의미.
   - Sprint 1 Phase 1 (commons) + Phase 2 (3 operator 동시 PR) 의 atomic
     migration 패턴이 더 적합.

3. **`pkg/sts` 로 통합 (PVC + topology 동시 처리)** — 거부.
   - 두 기능의 *호출 주기* 가 다름: TSC default 는 STS 빌드 시 매번,
     PVC expansion 은 reconcile 시 size diff 발견 시만.
   - 패키지 분리가 *호출자 mental model* 에도 더 명확 (resources/ vs
     controller/ 각각 분리되어 있는 현 컨벤션과 일치).

## Refs

- 사례:
  - `pkg/pvc/expansion.go`, `pkg/pvc/expansion_test.go`, `pkg/pvc/doc.go`
  - `pkg/topology/spread.go`, `pkg/topology/spread_test.go`, `pkg/topology/doc.go`
- 글로벌 표준: `standards/adr.md §3` (Nygard 5섹션).
- 추출 원본:
  - `postgres-operator/internal/controller/pvc_resize.go` (~120 LOC).
  - `mongodb-operator/internal/controller/pvc_resize.go` (~120 LOC).
  - `valkey-operator/internal/controller/pvc_resize.go` (~120 LOC).
  - `postgres-operator/internal/controller/topology_spread.go` (~48 LOC).
  - `mongodb-operator/internal/resources/topology_spread.go` (~46 LOC).
  - `valkey-operator/internal/resources/statefulset.go` (~40 LOC inline).
