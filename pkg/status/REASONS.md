# pkg/status — Condition Reason 카탈로그

> 4-repo (mongodb / valkey / postgres / commons) 공통 Condition Reason 사용 가이드. 본 문서는 `conditions.go` 의 const 선언을 *기계 가독* + *운영 가독* 양쪽으로 사용 가능하도록 표준화.

## Reason 별 사용 매트릭스

| Reason | Type 조합 | Status | 사용 시점 | Message 예시 |
|---|---|---|---|---|
| `Reconciling` | `Progressing` | True | Reconcile 진입 직후 + 진행 중인 모든 작업 | `"reconciling shard 2/4: scaling replicas 3->5"` |
| `Available` | `Ready` / `Available` | True | 모든 리소스 reconcile 완료 + 외부 호출 가능 | `"all 3 members healthy, primary at index 0"` |
| `NotApplicable` | (any) | False | 본 condition 이 현 spec 에서 의미 없음 (e.g. `MonitoringSpec` 미설정 시 `Monitoring=NotApplicable`) | `"monitoring disabled via spec.monitoring.enabled=false"` |
| `ReconcileError` | `Ready` / `Degraded` | False | reconcile 중 에러 발생 (transient retry 포함) | `"failed to create StatefulSet: <root cause>"` |
| `ExternalDependencyBlocked` | `Ready` | False | 외부 의존 (cert-manager, Secret, ConfigMap, CRD) 미준비로 차단 | `"cert-manager Certificate 'X' not ready: pending issuance"` |
| `ValidationFailed` | `Ready` / `Degraded` | False | webhook 또는 controller 측 validation 실패 (사용자 spec 오류) | `"spec.replicas must be >= 3 for HA topology"` |

## 사용 흐름 (4-repo 표준)

### 정상 Reconcile

```
1. Reconcile 진입 시점
   SetProgressing(conditions, ReasonReconciling, "reconciling cluster", gen)

2. 모든 리소스 reconcile 완료
   SetReady(conditions, ConditionTrue, ReasonAvailable, "all healthy", gen)
   SetAvailable(conditions, ConditionTrue, ReasonAvailable, "endpoints ready", gen)
   RemoveCondition(conditions, TypeProgressing)
```

### 에러 발생

```
1. 에러 catch 시점
   SetReadyFalse(conditions, ReasonReconcileError, err.Error(), gen)

2. 외부 의존 차단 발견 시점
   SetReadyFalse(conditions, ReasonExternalDependencyBlocked,
       "cert-manager Certificate X pending", gen)
   // Available 은 변경 안 함 (기존 연결 유지)

3. 사용자 spec 오류
   SetReadyFalse(conditions, ReasonValidationFailed,
       "spec.replicas must be >= 3", gen)
   SetDegraded(conditions, ReasonValidationFailed, "...", gen)
```

### 부분 장애

```
1. 일부 기능 장애 발견 시점
   SetDegraded(conditions, ReasonReconcileError,
       "backup CronJob failing, runtime healthy", gen)
   // Ready / Available 은 유지 (운영 가능)
```

### 기능 비활성

```
1. Optional spec 미설정 시점 (예: monitoring)
   meta.SetStatusCondition(conditions, metav1.Condition{
       Type:               "Monitoring",
       Status:             metav1.ConditionFalse,
       Reason:             ReasonNotApplicable,
       Message:            "monitoring disabled via spec",
       ObservedGeneration: gen,
   })
```

## 호출자 검증

`kubectl get <cr> -o jsonpath='{.status.conditions}'` 출력으로 매트릭스 검증:

```bash
kubectl get postgrescluster prod -o jsonpath='{range .status.conditions[*]}{.type}={.status}/{.reason}{"\n"}{end}'
# Ready=True/Available
# Progressing=False/Available
# Available=True/Available
```

## 4-repo 정합 게이트

본 카탈로그 추가/변경 시 *3 operator 모두* 반영. cross-repo audit:

```bash
for repo in mongodb-operator valkey-operator postgres-operator; do
  grep -rE "ReasonReconciling|ReasonAvailable|ReasonReconcileError" \
    /Users/phil/WorkSpace/public/$repo/internal/ | wc -l
done
# 0 = commons import 안 함 / >0 = import 사용 중
```

## Reason 추가 정책

신규 Reason 추가는 *cross-repo audit 가 1+* 인 경우만:
1. PR 본문에 *해당 reason 의 4-repo 사용 사례 1+ 인용*
2. 본 REASONS.md 매트릭스에 행 추가
3. CHANGELOG.md `Added` 섹션 entry
4. `pkg/status/conditions_test.go` 에 reason value 테스트

## References

- `conditions.go` — Reason / Type const 선언 + helper 함수
- `conditions_test.go` — unit test
- RFC-0018 §3.1 — pkg/status 스펙
- KEP-1623 — Kubernetes Standard Conditions
- ROADMAP.md "Condition reason 표준 카탈로그 문서화" (P-B.5.1 마감)
