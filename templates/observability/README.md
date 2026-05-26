# operator-commons / templates / observability

Helm library chart partial assets. Consumer charts use the named template
to define Prometheus Operator `ServiceMonitor` resources in a consistent
format.

## Assets

| File | Purpose |
|---|---|
| `_servicemonitor.tpl` | Prometheus Operator `ServiceMonitor` CR partial (named template `keiailab.observability.serviceMonitor`). |

PrometheusRule is deliberately excluded because downstream-specific
metrics dominate — each downstream operator authors its own
PrometheusRule. Common *metric naming conventions* and *recommended alert
patterns* are specified below.

## Standard metric naming convention

Each downstream operator exposes the following via the controller-runtime
`metrics.Registry` (in addition to controller-runtime defaults):

| Metric | Type | Labels | Meaning |
|---|---|---|---|
| `<op>_reconcile_total` | Counter | `result=success\|error\|requeue`, `kind` | Reconcile call count |
| `<op>_reconcile_duration_seconds` | Histogram | `kind` | Reconcile latency |
| `<op>_external_dep_call_total` | Counter | `target`, `result=success\|error` | External dependency calls |
| `<op>_finalizer_pending_total` | Gauge | `kind` | Objects pending deletion |

`<op>` is the operator name in underscore form.

## Recommended alert patterns (downstream PrometheusRule)

The library does not export a PrometheusRule, but recommends the
following three alerts:

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

Operator-specific alerts belong in each downstream operator's
PrometheusRule.

## Usage (consumer chart)

`charts/<name>/templates/servicemonitor.yaml`:

```yaml
{{- if and .Values.metrics.enabled .Values.metrics.serviceMonitor.enabled -}}
{{- include "keiailab.observability.serviceMonitor" . -}}
{{- end }}
```

`charts/<name>/values.yaml` required keys:

```yaml
metrics:
  enabled: true
  secure: false
  serviceMonitor:
    enabled: false
    interval: 30s
    labels: {}
    metricRelabelings: []
    relabelings: []
```

To include the partial in a consumer chart:

1. Sub-chart: copy to `charts/<name>/charts/keiailab-observability/`.
2. Git submodule: reference `operator-commons/templates/observability/`.
3. Build-time copy: `cp operator-commons/templates/observability/*.tpl
   charts/<name>/templates/_partials/` (Makefile target).

Record the adoption method in the consumer's ADR.

## Change policy

- Spec changes to this partial = public API breaking — require an ADR +
  downstream consumer LGTM ([GOVERNANCE.md](../../docs/GOVERNANCE.md)).
- Metric naming convention changes = consensus + downstream batch PR.
