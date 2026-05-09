# HANDOFF — operator-commons (Sprint A 진입)

> 본 HANDOFF.md 는 token-budget §5 표준 5 항목 (현재 상태 / 다음 단계 /
> 차단점 / 근거 링크 / 의사결정 기록) 을 따른다. 다음 세션 진입 시 가장
> 먼저 읽는다.

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
