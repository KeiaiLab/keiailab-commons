{{/*
keiailab-commons — External Secrets Operator partials.

본 partial 은 ExternalSecret manifest 를 raw YAML 로 출력한다. CRD 타입을
vendoring 하지 않으므로 consumer chart 는 ESO/Infisical 사용 여부를 values 로
선택할 수 있고, CRD 미설치 클러스터에서는 호출 자체를 비활성화하면 된다.
*/}}

{{/*
keiailab.secrets.externalSecret — Infisical/ESO 기반 Secret materialization.

caller 인자 (dict):
  - ctx:              helm context (.)
  - name:             ExternalSecret 이름
  - namespace:        namespace (optional, default .Release.Namespace)
  - labels:           metadata.labels YAML string 또는 map
  - refreshInterval:  refreshInterval (optional, default "1h")
  - secretStoreKind:  ClusterSecretStore | SecretStore (optional, default ClusterSecretStore)
  - secretStoreName:  secret store 이름 (필수)
  - targetName:       생성할 Secret 이름 (optional, default name)
  - creationPolicy:   Owner | Merge | Orphan | None (optional, default Owner)
  - data:             list of {secretKey, remoteKey, property?, decodingStrategy?}
  - targetTemplate:   target.template raw map (optional)

사용 예:

  {{- if .Values.externalSecrets.enabled }}
  {{ include "keiailab.secrets.externalSecret" (dict
      "ctx" .
      "name" "mongodb-admin"
      "labels" (include "mongodb-operator.labels" .)
      "secretStoreName" .Values.externalSecrets.clusterSecretStore
      "targetName" "mongodb-admin"
      "data" (list (dict "secretKey" "password" "remoteKey" "/data/mongodb/admin/password"))) }}
  {{- end }}
*/}}
{{- define "keiailab.secrets.externalSecret" -}}
{{- $ctx := .ctx -}}
{{- $namespace := .namespace | default $ctx.Release.Namespace -}}
{{- $refreshInterval := .refreshInterval | default "1h" -}}
{{- $secretStoreKind := .secretStoreKind | default "ClusterSecretStore" -}}
{{- $secretStoreName := required "keiailab.secrets.externalSecret: secretStoreName is required" .secretStoreName -}}
{{- $targetName := .targetName | default .name -}}
{{- $creationPolicy := .creationPolicy | default "Owner" -}}
apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata:
  name: {{ .name }}
  namespace: {{ $namespace }}
  labels:
    {{- if kindIs "string" .labels }}
    {{- .labels | nindent 4 }}
    {{- else }}
    {{- toYaml .labels | nindent 4 }}
    {{- end }}
spec:
  refreshInterval: {{ $refreshInterval | quote }}
  secretStoreRef:
    kind: {{ $secretStoreKind }}
    name: {{ $secretStoreName }}
  target:
    name: {{ $targetName }}
    creationPolicy: {{ $creationPolicy }}
    {{- with .targetTemplate }}
    template:
      {{- toYaml . | nindent 6 }}
    {{- end }}
  data:
    {{- range .data }}
    - secretKey: {{ required "keiailab.secrets.externalSecret: data[].secretKey is required" .secretKey }}
      remoteRef:
        key: {{ required "keiailab.secrets.externalSecret: data[].remoteKey is required" .remoteKey | quote }}
        {{- with .property }}
        property: {{ . | quote }}
        {{- end }}
        {{- with .decodingStrategy }}
        decodingStrategy: {{ . }}
        {{- end }}
    {{- end }}
{{- end }}
