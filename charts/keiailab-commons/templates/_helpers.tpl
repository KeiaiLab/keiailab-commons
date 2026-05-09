{{/*
keiailab-commons — Helm library chart helper partial.

Plan §2 D14 (Sprint B PR-B2) — RFC-0019 §3.1 implementation.

본 chart 는 type: library — consumer chart (mongodb-operator /
postgres-operator / valkey-operator) 가 dependency 로 import 후 include
호출.

Provided helpers:
  - keiailab.commonLabels                     — Helm 표준 공통 label.
  - keiailab.observability.serviceMonitor     — ServiceMonitor 공통 spec.

future helpers (별 PR):
  - keiailab.networkpolicy.dataplane          — RFC-0019 §3.2
  - keiailab.networkpolicy.controlplane       — RFC-0019 §3.2
  - keiailab.security.podSecurityRestricted   — RFC-0019 §3.4
  - keiailab.rbac.serviceAccount              — RFC-0019 §3.5
  - keiailab.rbac.controllerBase              — RFC-0019 §3.5
*/}}


{{/*
keiailab.commonLabels — Helm 표준 공통 label set.

4-repo operator chart 에서 동일 적용 — kubectl 검색 정합성 보장.
컨벤션: Kubernetes Recommended Labels
(https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/).

사용 예 (consumer chart 의 templates/* metadata.labels):

  metadata:
    labels:
      {{- include "keiailab.commonLabels" . | nindent 4 }}
*/}}
{{- define "keiailab.commonLabels" -}}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}


{{/*
keiailab.observability.serviceMonitor — Prometheus Operator ServiceMonitor
공통 spec.

전제 values (caller chart 의 .Values):
  - .Values.metrics.enabled                        (bool, 필수)
  - .Values.metrics.serviceMonitor.enabled         (bool, 필수)
  - .Values.metrics.serviceMonitor.interval        (string, default "30s")
  - .Values.metrics.serviceMonitor.labels          (map, optional)
  - .Values.metrics.secure                         (bool, optional, default false)

전제: caller 가 dict 로 fullname / labels / selectorLabels 를 전달
(consumer chart 별 helper 이름이 다르므로 partial 이 직접 호출 불가).

사용 예 (consumer chart 의 templates/servicemonitor.yaml):

  {{- include "keiailab.observability.serviceMonitor" (dict
      "ctx" .
      "fullname" (include "mongodb-operator.fullname" .)
      "labels" (include "mongodb-operator.labels" .)
      "selectorLabels" (include "mongodb-operator.selectorLabels" .)) }}
*/}}
{{- define "keiailab.observability.serviceMonitor" -}}
{{- $ctx := .ctx -}}
{{- if and $ctx.Values.metrics.enabled $ctx.Values.metrics.serviceMonitor.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ .fullname }}
  namespace: {{ $ctx.Release.Namespace }}
  labels:
    {{- .labels | nindent 4 }}
    {{- with $ctx.Values.metrics.serviceMonitor.labels }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
spec:
  endpoints:
    - port: metrics
      {{- if $ctx.Values.metrics.secure }}
      scheme: https
      tlsConfig:
        insecureSkipVerify: true
      {{- else }}
      scheme: http
      {{- end }}
      interval: {{ $ctx.Values.metrics.serviceMonitor.interval | default "30s" }}
  namespaceSelector:
    matchNames:
      - {{ $ctx.Release.Namespace }}
  selector:
    matchLabels:
      {{- .selectorLabels | nindent 6 }}
{{- end }}
{{- end }}
