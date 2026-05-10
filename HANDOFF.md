# HANDOFF — operator-commons (Sprint A 진입)

> 본 HANDOFF.md 는 token-budget §5 표준 5 항목 (현재 상태 / 다음 단계 /
> 차단점 / 근거 링크 / 의사결정 기록) 을 따른다. 다음 세션 진입 시 가장
> 먼저 읽는다.

## 2026-05-10 cross-repo 세션 결과 (operator-commons 영향)

본 세션 (Ralph loop, 2026-05-10) 의 operator-commons 직접 변경:

| 영역 | PR | 결과 |
|---|---|---|
| HANDOFF Sprint A/B/C 누적 progress 표 | #8 | merged |
| OSS hygiene (CITATION.cff + GH PR/Issue templates) | #9 | merged |

3 operator (valkey + mongodb + postgres) 모두 commons v0.6.0 (pkg/finalizer +
pkg/status sugar) 채택 완료 — `mongodb-operator#119/#120 ADR-0021/0022`,
`postgres-operator#21 ADR-0011`, `valkey-operator#18 ADR-0038` 통해 production
배포 도달.

3 operator 운영 latest version 동기:
- `valkey-operator` v1.0.10 (INC-0001 영구 fix + ADR-0039 self-heal active)
- `mongodb-operator` v1.4.20 (ShardDraining 백오프 regression fix + ADR-0022)
- `postgres-operator` v0.3.0-alpha.16 (OperatorHub bundle scaffold + ADR-0013)

OperatorHub.io 등록 신청 (multi-day review): community-operators PRs
[#8091](https://github.com/k8s-operatorhub/community-operators/pull/8091)
(valkey 1.0.10), [#8092](https://github.com/k8s-operatorhub/community-operators/pull/8092)
(keiailab-mongodb 1.4.20), [#8093](https://github.com/k8s-operatorhub/community-operators/pull/8093)
(keiailab-postgres 0.3.0-alpha.16). 3 PR 모두 orange/Deploy k8s CI PASS.

### 후속 (commons 자체 작업)

1. consumer chart partial include — mongodb / postgres / valkey 의 ServiceMonitor
   / NetworkPolicy / RBAC / Security yaml 이 commons library chart partial include
   로 교체. commons OCI publish (PR-B2.2 별 PR) 의존.
2. v0.7.0 → v0.8.0 — 본 세션의 cumulative consumer 패턴 실측 후 *commons API 표면
   안정화* 결정.

---

## 2026-05-09 Sprint A/B/C 누적 진행 (Ralph iter 1~16, 23 PR 머지)

> Plan: `~/.claude/plans/1-https-artifacthub-io-packages-helm-clo-synthetic-gem.md`

### 머지 PR 표

| repo | PR | 내용 | tag/version |
|---|---|---|---|
| operator-commons | #1 | RFC-0018 + ADR-0003 (status sugars) | v0.6.0 |
| operator-commons | #2 | pkg/version Generic Matrix[E] (ADR-0004) | v0.7.0 |
| operator-commons | #3 | library chart 신설 (RFC-0019, ADR-0005) | chart v0.1.0 |
| operator-commons | #4 | NetworkPolicy partials (ADR-0006) | chart v0.2.0 |
| operator-commons | #5 | RBAC partials (ADR-0007) | chart v0.3.0 |
| operator-commons | #6 | PSS Restricted partials (ADR-0008) — RFC-0019 §3 완결 | chart v0.4.0 |
| valkey-operator | #5 | cosign + SLSA L2 (ADR-0033) | — |
| valkey-operator | #6 | v1alpha2 + AuthSpec.Required (ADR-0034) | — |
| valkey-operator | #7 | commons v0.6.0 bump | — |
| valkey-operator | #8 | pkg/finalizer migration (ADR-0038) | — |
| valkey-operator | #9 | NetworkPolicy.AutoCreate (ADR-0035) | — |
| valkey-operator | #10 | PodSecurityRestricted (ADR-0036) | — |
| valkey-operator | #11 | Sentinel migration runbook (PR-C7) | — |
| valkey-operator | #12 | AuthSpec.RotationPolicy enum (ADR-0031) | — |
| valkey-operator | #13 | ADR-0018 정식 (Cluster Auto-Resharding) | — |
| valkey-operator | #14 | ValkeySpec.Modules (ADR-0032) | — |
| mongodb-operator | #116 | Sprint A HANDOFF entry | — |
| mongodb-operator | #117 | commons v0.6.0 bump | — |
| mongodb-operator | #118 | pkg/finalizer migration (ADR-0021) | — |
| postgres-operator | #19 | Sprint A HANDOFF entry | — |
| postgres-operator | #20 | commons v0.6.0 bump | — |
| postgres-operator | #21 | pkg/status partial adoption (ADR-0011) | — |
| postgres-operator | #22 | matrix.go → commons Matrix[Combo] (ADR-0012) | — |

### Plan §1 Phase 1 갭 해소 현황

| Gap | 상태 | PR |
|---|---|---|
| A: Sentinel 미지원 | 거부 보존 + 운영 runbook | valkey #11 |
| B: Password Rotation | type 추가, controller 후속 | valkey #12 (PR-B7.2 후속) |
| C: HPA | ADR-0027 deferred 보존 | PR-C5 후속 |
| D: Custom Modules | type 추가, controller 후속 | valkey #14 (PR-C6.2 후속) |
| E: 공급망 보증 | cosign + SLSA L2 채택 | valkey #5 |
| F: Cluster Auto-Resharding | ADR-0018 정식, controller 후속 | valkey #13 (PR-B8.2 후속) |

### 사용자 결정 1 (보안 3종 옵션화) 완료

- ADR-0034 Auth Optional / ADR-0035 NetworkPolicy.AutoCreate / ADR-0036
  PSS Restricted Optional — 모두 v1alpha2 type module 에 default=true
  유지로 secure-by-default 보존.
- D4 v1alpha2 + conversion webhook: type module 완료, hub 전환 controller
  분기는 PR-A2.2 후속.

### RFC-0019 implementation 완결

§3.1 (commonLabels + ServiceMonitor) + §3.2 (NetworkPolicy partials) +
§3.4 (PSS partials) + §3.5 (RBAC partials) — commons library chart
v0.4.0 implementation. consumer chart 의 partial include 후속 (commons
OCI publish PR-B2.2 의존).

### 잔여 *대형 controller 작업* (다음 세션 진입점)

모두 valkey PR-A2.2 (v1alpha2 hub 전환) 의존:

1. **PR-A2.2** — valkey hub 전환: 5 CRD × 2 conversion 함수 + 4 controller
   import 변경 + cmd/main.go SchemeBuilder + ensureAuthSecret Required
   분기 + controller-gen 재실행. T3 단일 세션 dedicated.
2. **PR-B7.2** — Password Rotation reconcile: Secret resourceVersion
   watch + rotatePassword helper + replication 의 replica → primary
   순서 강제.
3. **PR-B8.2** — Cluster Resharding reconcile: vk.ClusterMigrateSlots
   helper (16384 slot 256-batch + ASKING) + reconcileResharding phase +
   ReshardingProgress status field + e2e 3→5 shard.
4. **PR-B9** — OperatorHub.io: bundle/manifests + ClusterServiceVersion
   + community-operators repo PR. 외부 visibility.
5. **PR-C5** — HPA reconcile: ValkeySpec.Autoscaling + HPA reconcile +
   ScalePolicy.Deliberate webhook validation.
6. **PR-C6.2** — Custom Modules: statefulset.go init container mount +
   emptyDir + valkey container --loadmodule arg + webhook allow-list +
   e2e (valkey-search FT.SEARCH).

### consumer chart partial include (별도 후속)

mongodb / postgres / valkey 의 ServiceMonitor / NetworkPolicy / RBAC /
Security yaml 이 commons library chart partial include 로 교체. commons
OCI publish (PR-B2.2 별 PR) 의존.

---

## 현재 상태

- 마지막 main commit: `fb7c8c6` chore(audit): .codecov.yml 신규 (4-repo target 70% 절대 floor 통일)
- working tree 상태 (Sprint A PR-A1 진행 중, 미커밋):
  - `pkg/status/conditions.go`: `SetAvailable` + `SetReadyFalse` 슈가 2종 추가
  - `pkg/status/conditions_test.go`: 테스트 2종 추가 (TestSetAvailable_CoexistsWithReady / TestSetReadyFalse_EquivalentToSetReady)
  - `docs/kb/rfc/0018-status-finalizer-standard.md`: 신규 (RFC 본문)
  - `docs/kb/adr/0003-rfc-0018-pkg-status-finalizer-adoption.md`: 신규
  - `docs/kb/adr/INDEX.md`: ADR-0003 행 추가
  - `pkg/finalizer/`: *변경 없음* (ADR-0003 §Decision 2 — `EnsureRemoval` 헬퍼 신설 보류)

검증 통과:
- `go test -C /Users/phil/WorkSpace/public/operator-commons ./pkg/status/... -count=1`: PASS (7 case)
- `go vet -C /Users/phil/WorkSpace/public/operator-commons ./pkg/status/...`: exit 0
- 추가 검증 대기: `make lint` (golangci-lint 1.65+) + `make test` 전체.

## 다음 단계

1. **즉시** (본 세션 또는 다음 세션):
   - `make -C /Users/phil/WorkSpace/public/operator-commons lint` 통과 확인.
   - `make -C /Users/phil/WorkSpace/public/operator-commons test` 전체 통과 확인.
   - commit (사용자 명시 후): `feat(status): SetAvailable + SetReadyFalse 슈가 추가 + RFC-0018 본문 + ADR-0003`
   - tag commons v0.6.0 + push.

2. **후속 (PR-A5/A6/A7 의 의존성)**:
   - 3 operator 의 `go.mod` `github.com/keiailab/operator-commons` v0.5.0 → v0.6.0 bump.
   - RFC-0018 §3.2 Migration 단계 2 진입:
     - valkey-operator (PR-A6): `controllerutil.AddFinalizer/RemoveFinalizer` → `finalizer.Add/Remove` (4 controller).
     - mongodb-operator (PR-A5): 동일 (3 controller).
     - postgres-operator (PR-A7): pkg/status 만 채택 — pkg/finalizer 비대칭 보존 (ADR-0008 갱신).

## 차단점

- 본 PR-A1 자체는 *self-contained* — 차단점 없음.
- 다른 Sprint A PR (A2~A7) 은 commons v0.6.0 release tag 머지 후 진입.

## 근거 링크

- 본 세션 기준 plan: `~/.claude/plans/1-https-artifacthub-io-packages-helm-clo-synthetic-gem.md` §2 D10/D11
- RFC-0018: `docs/kb/rfc/0018-status-finalizer-standard.md`
- ADR-0003: `docs/kb/adr/0003-rfc-0018-pkg-status-finalizer-adoption.md`
- 글로벌 표준: `~/Documents/ai-dev/standards/adr.md §6` (RFC vs ADR 구분)
- 공통화 진단 (Phase 1): valkey-operator 의 `pkg/finalizer` / `pkg/status` 채택률 0% 측정.

## 의사결정 기록 (Sprint A 진입)

본 세션에서 plan 단계 결정 (D10/D11) 을 *코드 base 실측* 후 정제했다:

1. **EnsureRemoval 헬퍼 신설 보류** (ADR-0003 §Decision 3):
   - 근거: `pkg/finalizer/finalizer.go:7-10` docstring 의 *controller-runtime
     미의존 원칙* 보존이 호출자 boilerplate 2 줄 회피보다 우위.
   - 대안 보존: 별도 sub-package `pkg/finalizer/runtime` 분리는 v1.0+ 에서
     재고려 가능 — 단 현 시점은 *commons zero-dep 원칙* 보호.

2. **pkg/status 슈가 2종 (SetAvailable + SetReadyFalse) 만 추가**:
   - 근거: API 변경 최소화 + consumer migration 우선 전략. 기존
     `SetReady/SetProgressing/SetDegraded` + read-side helpers 는 유지.
   - 정합성: `meta.SetStatusCondition` apimeta 위임 패턴 보존 — wire-level
     무변경.

3. **postgres pkg/finalizer 비대칭 보존**:
   - 근거: ADR-0008 cascade-delete-by-OwnerReference 결정이 *의도된 비대칭*.
     BackupCleanupJob CRD 가 외부 자원 cleanup 분리 처리.
   - 영향: RFC-0018 §3.2 Migration 단계 2 에서 postgres 만 pkg/status 채택,
     pkg/finalizer 미채택 — ADR-0008 갱신으로 사유 보존.

4. **Reason 카탈로그 통일** (RFC-0018 §3.3):
   - generic 4종 (Ready/Progressing/Degraded/Available) + 6 Reason 만 commons 강제.
   - 도메인 특이 type (`ShardsReady`, `PrimaryUnreachable`, `BackupHealthy` 등)
     은 각 repo 보존 — 무관 type strawman 노출 회피.

5. **외부 contract 변경 release note 의무**:
   - valkey 의 `ReasonReconcileFailed` → commons `ReasonReconcileError` 전환.
   - 1 release window (commons v0.6.0 → v0.7.0) 에서 양쪽 동작 (apimeta dedup).
   - alert rule 갱신 안내 → CHANGELOG.md + release note 명시.
