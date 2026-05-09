{{/*
keiailab-commons — NetworkPolicy partials (RFC-0019 §3.2).

Plan §2 D13 / commons-ADR-0006.

두 패턴 추출:
  - keiailab.networkpolicy.dataplane    — managed workload 보호 (postgres 패턴).
  - keiailab.networkpolicy.controlplane — operator manager 자체 보호 (valkey 패턴).
*/}}


{{/*
keiailab.networkpolicy.dataplane — managed dataplane workload 보호 NetworkPolicy 묶음.

postgres 패턴: default-deny + allow-internal-instance (같은 managed-by
label 의 pod 간만 dataplane port 허용).

caller 인자 (dict):
  - ctx:        helm context (.)
  - fullname:   caller chart 의 fullname (NetworkPolicy name prefix)
  - labels:     caller chart 의 labels
  - managedBy:  managed-by label 값 (예: "keiailab-postgres-operator")
  - port:       dataplane port (예: 5432 / 27017 / 6379)

사용 예 (consumer chart 의 templates/networkpolicy.yaml):

  {{- if .Values.networkPolicies.enabled -}}
  {{ include "keiailab.networkpolicy.dataplane" (dict
      "ctx" .
      "fullname" (include "postgres-operator.fullname" .)
      "labels" (include "postgres-operator.labels" .)
      "managedBy" "keiailab-postgres-operator"
      "port" 5432) }}
  {{- end }}
*/}}
{{- define "keiailab.networkpolicy.dataplane" -}}
{{- $ctx := .ctx -}}
{{- $managedBy := .managedBy -}}
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ .fullname }}-dataplane-default-deny
  namespace: {{ $ctx.Release.Namespace }}
  labels:
    {{- .labels | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/managed-by: {{ $managedBy }}
  policyTypes: [Ingress, Egress]
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ .fullname }}-dataplane-allow-internal
  namespace: {{ $ctx.Release.Namespace }}
  labels:
    {{- .labels | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/managed-by: {{ $managedBy }}
  policyTypes: [Ingress, Egress]
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app.kubernetes.io/managed-by: {{ $managedBy }}
      ports:
        - port: {{ .port }}
          protocol: TCP
  egress:
    - to:
        - podSelector:
            matchLabels:
              app.kubernetes.io/managed-by: {{ $managedBy }}
      ports:
        - port: {{ .port }}
          protocol: TCP
{{- end }}


{{/*
keiailab.networkpolicy.controlplane — operator manager pod 자체 보호.

valkey 패턴: manager pod 의 ingress (metrics + 선택적 webhook) +
egress (K8s API + DNS + 사용자 추가).

caller 인자 (dict):
  - ctx:               helm context (.)
  - fullname:          caller chart 의 fullname
  - labels:            caller chart 의 labels
  - selectorLabels:    caller chart 의 selectorLabels
  - metricsPort:       metrics port (default 8443)
  - webhookEnabled:    webhook 활성 여부 (bool, optional)
  - webhookPort:       webhook port (optional, default 9443)
  - additionalIngress: 추가 ingress rules (list, optional, raw YAML)
  - additionalEgress:  추가 egress rules (list, optional, raw YAML)

사용 예 (consumer chart 의 templates/networkpolicy.yaml):

  {{- if .Values.networkPolicy.enabled }}
  {{ include "keiailab.networkpolicy.controlplane" (dict
      "ctx" .
      "fullname" (include "valkey-operator.fullname" .)
      "labels" (include "valkey-operator.labels" .)
      "selectorLabels" (include "valkey-operator.selectorLabels" .)
      "metricsPort" .Values.service.metricsPort
      "webhookEnabled" .Values.webhook.enabled
      "webhookPort" .Values.webhook.port
      "additionalIngress" .Values.networkPolicy.ingress
      "additionalEgress" .Values.networkPolicy.egress) }}
  {{- end }}
*/}}
{{- define "keiailab.networkpolicy.controlplane" -}}
{{- $ctx := .ctx -}}
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ .fullname }}
  namespace: {{ $ctx.Release.Namespace }}
  labels:
    {{- .labels | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      {{- .selectorLabels | nindent 6 }}
  policyTypes:
    - Ingress
    - Egress
  ingress:
    # Prometheus metrics scrape.
    - ports:
        - protocol: TCP
          port: {{ .metricsPort | default 8443 }}
    {{- if .webhookEnabled }}
    # Admission webhook — kube-apiserver 호출.
    - ports:
        - protocol: TCP
          port: {{ .webhookPort | default 9443 }}
    {{- end }}
    {{- with .additionalIngress }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
  egress:
    # K8s API server.
    - to:
        - namespaceSelector: {}
      ports:
        - protocol: TCP
          port: 443
    # DNS.
    - to:
        - namespaceSelector: {}
      ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
    {{- with .additionalEgress }}
    {{- toYaml . | nindent 4 }}
    {{- end }}
{{- end }}
