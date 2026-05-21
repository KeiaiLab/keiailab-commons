# operator-commons / templates / observability

RFC-0019 §3.1 표준 — 4-repo keiailab operator 공통 *helm chart partial* 자산.

## 보관 자산

| 파일 | 용도 |
|---|---|
| `_servicemonitor.tpl` | prometheus-operator `ServiceMonitor` CR partial (named template `keiailab.observability.serviceMonitor`) |

PrometheusRule 은 *operator 특화 메트릭* (예: valkey 의 `valkey_cluster_state_ok`, mongodb 의 `mongodb_replicaset_primary_count`) 이 많아 4-repo 공통 partial 이 부적합 → **각 repo 자체 작성**. 단 *공통 메트릭 명명 규약* + *공통 alert* 일부는 본 README 에 명시.

## 표준 메트릭 명명 규약 (4-repo 공통)

각 operator 는 `cmd/main.go` 에서 controller-runtime metrics.Registry 를 활용해 다음을 노출 (controller-runtime 기본 메트릭 외):

| 메트릭 | 타입 | 라벨 | 의미 |
|---|---|---|---|
| `<op>_reconcile_total` | Counter | `result=success\|error\|requeue`, `kind` | reconcile 호출 횟수 |
| `<op>_reconcile_duration_seconds` | Histogram | `kind` | reconcile latency |
| `<op>_external_dep_call_total` | Counter | `target`, `result=success\|error` | 외부 의존 호출 (Mongo / Postgres / Valkey instance) |
| `<op>_finalizer_pending_total` | Gauge | `kind` | 삭제 대기 중인 객체 수 |

`<op>` 는 operator 이름의 *underscore* 형 (예: `mongodb_operator`, `postgres_operator`, `valkey_operator`).

## 공통 alert 권장 (각 repo 에 PrometheusRule 로 추가)

본 라이브러리는 PrometheusRule 자체를 export 하지 않지만, 4-repo 가 *동일 양식* 으로 다음 3건은 보유하는 것을 권장:

```yaml
- alert: <Op>OperatorDown
  expr: up{job=~"<op>-operator.*"} == 0
  for: 2m
  labels: { severity: critical }
- alert: <Op>OperatorReconcileErrorsHigh
  expr: rate(<op>_reconcile_total{result="error"}[5m]) > 0.1
  for: 5m
  labels: { severity: warning }
- alert: <Op>OperatorReconcileSlow
  expr: histogram_quantile(0.99, rate(<op>_reconcile_duration_seconds_bucket[5m])) > 30
  for: 10m
  labels: { severity: info }
```

operator-specific alert (예: redis cluster slot 분배, postgres replication lag, mongodb arbiter 상태) 는 각 repo PrometheusRule 에 추가.

## 사용법 (consumer chart)

`charts/<name>/templates/servicemonitor.yaml`:

```yaml
{{- if and .Values.metrics.enabled .Values.metrics.serviceMonitor.enabled -}}
{{- include "keiailab.observability.serviceMonitor" . -}}
{{- end }}
```

`charts/<name>/values.yaml` 에 다음 키 추가:

```yaml
metrics:
  enabled: true
  secure: false
  serviceMonitor:
    enabled: false                # opt-in (prometheus-operator 미설치 환경 보호)
    interval: 30s
    labels: {}
    metricRelabelings: []
    relabelings: []
```

partial 자체를 chart 에 *포함* 하려면:

1. `charts/<name>/charts/keiailab-observability/` (sub-chart) 로 복사
2. 또는 git submodule 로 `operator-commons/templates/observability/` 참조
3. 또는 빌드 시점 `cp operator-commons/templates/observability/*.tpl charts/<name>/templates/_partials/` (Makefile 타겟)

각 repo 의 ADR 에 채택 방식 기록.

## 변경 정책

- 본 partial 의 spec 변경 = *공개 API breaking* — RFC + 3 consumer LGTM (docs/GOVERNANCE.md "공개 API breaking" 절차).
- 메트릭 명명 규약 변경 = *합의 후 4-repo 동시 PR* (operator-specific 로 발산하면 본 표준 의미 상실).
