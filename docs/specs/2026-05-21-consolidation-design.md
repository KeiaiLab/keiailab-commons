# S5 Design Draft: operator-commons 공통화

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Draft (사용자 검토 대기) |
| 위치 | 추후 `operator-commons/docs/specs/2026-05-21-consolidation-design.md` 로 push |
| 의존 | S1+ + S2 + S4 (commons SSOT 정비) + S7 모두 완료 후 진입 |

## 1. 배경

S5 사전 분석 결과 (subagent a44cf6fce8d8a5dad):

### 추출 후보 (commons 부재 + 3 operator 중복) — 10건

| # | 함수/패키지 | 적용 위치 | LOC 감소 추정 |
|---|---|---|---|
| 1 | `pkg/reconcile.HandleFinalizerCleanup` | mongodb `helpers.go:74`, valkey `helpers.go:66` | ~100 |
| 2 | `pkg/reconcile.Statusable` + `ApplyErrorCondition` | mongodb `helpers.go:100,114`, valkey `helpers.go:89,96` | ~160 |
| 3 | `pkg/reconcile.UpdateStatusWithRetry` | mongodb `status_update.go:36`, valkey `status_update.go:27`, postgres `controller.upsert` | ~90 |
| 4 | `pkg/resources.Apply{ConfigMap,Service,STS,NetworkPolicy,PDB,HPA,ServiceMonitor}` | mongodb `resources_apply.go` (334), valkey `resources_apply.go` (176), postgres `upsert` (627) | **~400** + drift 위험 제거 |
| 5 | `pkg/storageclass.ExpandDataPVCs` + `expandSinglePVC` (확장) | postgres `pvc_resize.go:33`, mongodb `pvc_resize.go:33`, valkey `pvc_resize.go:52` | ~250 |
| 6 | `pkg/monitoring.RegisterReconcileSLO` | postgres `metrics.go:37`, mongodb `metrics.go:22`, valkey `metrics.go:78` | ~150 |
| 7 | `pkg/reconcile.RequeueIntervals` (Steady/Provisioning/WaitExternal/Progress const) | mongodb `intervals.go`, valkey `intervals.go` | ~50 |
| 8 | `pkg/apis/common` 의 13개 shared CRD type | mongodb `common_types.go` (13), valkey `common_types.go` (13), postgres 부분 | **~600** + CRD 일관성 |
| 9 | `pkg/reconcile.IsPaused` + `PausedAnnotationKey` | valkey `helpers.go:176`, postgres pooler 변형 | ~30 |
| 10 | `scripts/release-smoke-test.sh` (공유) | 3 operator + commons | ~450 (drift 제거) |

**총 추출 효과**: ~2,280 LOC 감소 + drift 위험 대량 제거.

### 미사용 helper (commons 에 있음 + 3 operator 활용 부족) — 9건

| commons helper | 미사용 operator | 마이그레이션 비용 |
|---|---|---|
| `pkg/status` (SetReady/Progressing/Degraded/Available) | postgres/mongodb/valkey **전부** | low (drop-in) |
| `pkg/finalizer` | postgres | low |
| `pkg/networkpolicy` (combo) | postgres | medium |
| `pkg/monitoring` (otel/rule) | postgres/mongodb | medium |
| `pkg/labels` v2 | (부분) | low |
| `pkg/probes` | valkey | medium |
| Helm `keiailab.security.podSecurityContext/containerSecurityContext` | 전부 | medium |
| Helm `keiailab.rbac.*` / `keiailab.networkpolicy.*` | 전부 | medium |
| `_servicemonitor.tpl` (keiailab-commons chart) | 전부 | low |

## 2. Goals + Non-Goals

### Goals
| ID | 목표 | 검증 |
|---|---|---|
| G1 | pkg/reconcile 신규 패키지 (6 helper) | `go test ./pkg/reconcile/...` PASS + 3 operator 채택 |
| G2 | pkg/resources 신규 패키지 (7 apply 함수) | 동일 |
| G3 | pkg/storageclass + pkg/monitoring 확장 (기존 패키지에 add) | 동일 |
| G4 | pkg/apis/common 신규 (13 shared CRD type) | conversion webhook 통과 |
| G5 | 3 operator 의 중복 코드 제거 + commons import | `wc -l` 측정 ≥ 1,500 LOC 감소 |
| G6 | 9 미사용 helper 의 3 operator 적용 | grep import |
| G7 | commons release `v0.8.0` → `v0.9.0` (minor bump) | tag + CHANGELOG |
| G8 | 3 operator 의 commons 의존 bump | go.mod update + integration test |

### Non-Goals
- 신규 기능 추가 (오직 *추출 + 재사용*)
- breaking API 변경 (v1alpha1 → v1alpha2 conversion 은 후속 cycle)
- 다국어 (S4 별 cycle)

## 3. 아키텍처

```
[Wave 1] commons 패키지 추출
   ↓
   ├─ pkg/reconcile (helper 6) — P0
   ├─ pkg/resources (apply 7) — P0
   ├─ pkg/storageclass.ExpandDataPVCs (확장) — P1
   ├─ pkg/monitoring.RegisterReconcileSLO (확장) — P1
   └─ pkg/apis/common (13 type, deepcopy 포함) — P2 (가장 위험)
   ↓
[Wave 2] commons release v0.9.0 (tag + CHANGELOG)
   ↓
[Wave 3] 3 operator 채택 (sequential — 위험 분산)
   ↓
   ├─ valkey-operator 먼저 (가장 작음)
   ├─ mongodb-operator
   └─ postgres-operator 마지막 (가장 큼)
```

## 4. 단계별 상세

### Phase 0 — pre-flight
- S1+/S2/S4/S7 완료 확인
- commons main fetch
- subagent 분석 결과 (s5-analysis) 인용

### Phase 1 — pkg/reconcile 패키지 신설 (Wave 1, P0)
- branch: `feat/pkg-reconcile-helpers-2026-05-21`
- 신규: `pkg/reconcile/finalizer.go` (HandleFinalizerCleanup)
- 신규: `pkg/reconcile/status.go` (Statusable interface + ApplyErrorCondition + UpdateStatusWithRetry)
- 신규: `pkg/reconcile/intervals.go` (RequeueIntervals const)
- 신규: `pkg/reconcile/pause.go` (IsPaused + PausedAnnotationKey)
- 신규: `pkg/reconcile/doc.go` (API Stability Tier marker)
- 신규: `pkg/reconcile/finalizer_test.go` + status_test.go (envtest)
- ADR `0036-pkg-reconcile-extraction.md`

### Phase 2 — pkg/resources 패키지 신설 (Wave 1, P0)
- branch: `feat/pkg-resources-apply-2026-05-21`
- 신규: `pkg/resources/configmap.go`, `service.go`, `statefulset.go`, `networkpolicy.go`, `pdb.go`, `hpa.go`, `servicemonitor.go`
- 각 `Apply<Resource>(ctx, c, scheme, owner, desired)` 함수
- 사용 패턴: postgres `controller.upsert` (627줄) generic + mongodb `apply*` 도메인 → **둘 다 노출**:
  - `pkg/resources.Upsert` (generic)
  - `pkg/resources.Apply<Resource>` (도메인별 type-safe)
- 신규: `pkg/resources/*_test.go` (envtest + fake client)
- ADR `0037-pkg-resources-extraction.md`

### Phase 3 — pkg/storageclass + pkg/monitoring 확장 (Wave 1, P1)
- branch: `feat/pkg-storageclass-pvc-expand-2026-05-21`
- `pkg/storageclass/pvc_resize.go` 신규 (ExpandDataPVCs + expandSinglePVC)
- branch: `feat/pkg-monitoring-slo-2026-05-21`
- `pkg/monitoring/slo.go` 신규 (RegisterReconcileSLO)
- 각각 test + ADR

### Phase 4 — pkg/apis/common 신규 (Wave 1, P2 — 가장 위험)
- branch: `feat/pkg-apis-common-shared-types-2026-05-21`
- 신규: `pkg/apis/common/v1alpha1/types.go` (13 shared CRD type)
- 신규: `pkg/apis/common/v1alpha1/zz_generated.deepcopy.go` (controller-gen)
- 신규: `pkg/apis/common/v1alpha1/doc.go` (Tier marker, BREAKING 가능성 경고)
- ADR `0038-pkg-apis-common-shared-types.md`
- 채택 전략: v1alpha1 → 각 operator 의 type field 가 `common.StorageSpec` *typealias* — backward-compat. 후속 v1alpha2 에서 native type.

### Wave 2 — commons v0.9.0 release
- branch: `release/v0.9.0`
- CHANGELOG.md 갱신 (Keep a Changelog)
- `git tag -a v0.9.0` + push
- `scripts/release.sh` 실행
- ADR Status: Proposed → Accepted

### Wave 3 — 3 operator 채택 (sequential)

#### Wave 3.1 — valkey-operator
- branch: `feat/adopt-commons-v0.9.0-2026-05-21`
- go.mod: `github.com/keiailab/operator-commons v0.9.0`
- 중복 코드 제거 + import 적용
- 변경 영역: helpers.go (~150 LOC 제거), status_update.go (~50), resources_apply.go (~150), pvc_resize.go (~80), intervals.go (~25), metrics.go (~50) = **~505 LOC**
- e2e: kind cluster + helm install + CR reconcile + delete
- PR + merge

#### Wave 3.2 — mongodb-operator (valkey 패턴 검증 후)
- 동일 패턴, 변경 영역 ~700 LOC

#### Wave 3.3 — postgres-operator (마지막, 가장 큼)
- generic upsert 패턴 → pkg/resources 채택
- 변경 영역 ~600 LOC

### Wave 4 — 미사용 helper 적용 (9 helper)
- Wave 3 와 병행 가능
- 각 operator 별로 pkg/status, pkg/finalizer, pkg/networkpolicy, pkg/monitoring, pkg/probes adoption
- 별 PR 각각

### Wave 5 — 검증 + audit
- audit-production-grade.sh 실행 → P1/P2 정합 확인
- 3 operator e2e 모두 PASS
- ADR 모두 Accepted

## 5. 리스크

| 리스크 | 영향 | 완화 |
|---|---|---|
| pkg/apis/common 의 breaking change | CRD 사용자 영향 (v1alpha1 → v1alpha2 마이그레이션) | type alias 로 backward-compat 유지 + 별 cycle 에서 v1alpha2 |
| postgres `controller.upsert` (627줄) → `pkg/resources.Apply*` 전환 | 큰 변경, 회귀 위험 | Wave 3.3 마지막 (검증된 패턴 후) + integration test 강제 |
| commons 의존성 bump 충돌 (3 operator 동시) | go.mod resolution 충돌 | sequential Wave 3.1 → 3.2 → 3.3 |
| 추출 helper 의 API design 결함 | 후속 변경 부담 | Phase 1-4 의 test 강도 ↑ (≥90% coverage) + API stability Tier marker |
| 사용자 결정 5건 부재 | spec 진행 차단 | default 채택 (본 draft) + 사용자 검토 시 정정 |

## 6. 사용자 결정 5건 (본 draft 의 default)

| # | 결정 | Default 채택 | 변경 가능 |
|---|---|---|---|
| 1 | 패키지 boundary | `pkg/reconcile` (helper) + `pkg/resources` (apply) + `pkg/apis/common` (type) 명확 분리 | ✅ |
| 2 | breaking change 정책 | type alias 로 backward-compat + 별 cycle 에서 v1alpha2 conversion | ✅ |
| 3 | semver | 신규 패키지 = minor bump (v0.8 → v0.9), 기존 API breaking = major bump | ✅ |
| 4 | 추출 우선순위 | P0 (reconcile helper) → P0 (resources) → P1 (storageclass + monitoring) → P2 (apis common, 가장 위험) | ✅ |
| 5 | postgres upsert 통합 | 둘 다 노출 (`pkg/resources.Upsert` generic + `Apply<Resource>` domain) | ✅ |

## 7. 성공 기준

```bash
# 1. commons 신규 패키지 4개 (reconcile + resources + apis/common + 확장 2)
test $(ls operator-commons/pkg/reconcile operator-commons/pkg/resources operator-commons/pkg/apis/common 2>/dev/null | wc -l) -ge 3

# 2. commons v0.9.0 tag
git -C operator-commons tag -l v0.9.0 | grep -q v0.9.0

# 3. 3 operator 의 commons v0.9.0 채택
for repo in postgres-operator mongodb-operator valkey-operator; do
  grep -q "operator-commons v0.9" $repo/go.mod || exit 1
done

# 4. 중복 코드 제거 (helpers.go LOC 50% 이상 감소)
# (각 operator 별 비교)

# 5. e2e 3 operator 모두 PASS
for repo in ...; do make -C $repo integration-test; done

# 6. audit-production-grade.sh — P1/P2 ≥ 90% ✅
./audit-production-grade.sh
```

## 8. 본 draft 의 향후

- 5 cycle (S1+, S2, S4, S6, S7) 모두 완료 후 commons 에 push
- 사용자 검토 + 5건 결정 변경 가능
- 승인 후 writing-plans 호출 → 실 cycle 진입
