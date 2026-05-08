{{- /*
keiailab 4-repo 공통 ServiceMonitor partial — RFC-0019 §3.1.

본 partial 은 controller-runtime 의 metrics endpoint (port=metrics, path=/metrics)
를 prometheus-operator 의 ServiceMonitor CR 로 노출하는 *generic* 자산이다.
operator-specific 부분 (alert rules, custom metric 명명) 은 각 repo 의 별도
PrometheusRule 에 둔다.

사용법 (consumer chart 측):

  {{- if and .Values.metrics.enabled .Values.metrics.serviceMonitor.enabled -}}
  {{- include "keiailab.observability.serviceMonitor" . }}
  {{- end }}

호출자 chart 가 보유해야 할 helper templates:
  - {{ .Values.operatorName }} — 또는 named template "<chart>.fullname"
  - "<chart>.labels" — 표준 라벨 set
  - "<chart>.selectorLabels" — Service selector

values.yaml 가 보유해야 할 키:
  metrics:
    enabled: true                                  # 마스터 토글
    secure: false                                  # https + tlsConfig (cert-manager 환경)
    serviceMonitor:
      enabled: true                                # ServiceMonitor CR 생성 토글
      interval: 30s
      labels: {}                                   # 추가 라벨 (예: prometheus instance selector)
      metricRelabelings: []                        # 표준 relabel 쿼리
      relabelings: []
*/}}
{{- define "keiailab.observability.serviceMonitor" -}}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ include (printf "%s.fullname" .Chart.Name) . }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include (printf "%s.labels" .Chart.Name) . | nindent 4 }}
    {{- with .Values.metrics.serviceMonitor.labels }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
spec:
  endpoints:
    - port: metrics
      {{- if .Values.metrics.secure }}
      scheme: https
      tlsConfig:
        insecureSkipVerify: true
      {{- else }}
      scheme: http
      {{- end }}
      interval: {{ .Values.metrics.serviceMonitor.interval | default "30s" }}
      {{- with .Values.metrics.serviceMonitor.metricRelabelings }}
      metricRelabelings:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.metrics.serviceMonitor.relabelings }}
      relabelings:
        {{- toYaml . | nindent 8 }}
      {{- end }}
  namespaceSelector:
    matchNames:
      - {{ .Release.Namespace }}
  selector:
    matchLabels:
      {{- include (printf "%s.selectorLabels" .Chart.Name) . | nindent 6 }}
{{- end -}}
