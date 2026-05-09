{{/*
keiailab-commons — RBAC partials (RFC-0019 §3.5).

Plan §2 D15 / commons-ADR-0007.

세 partial 추출:
  - keiailab.rbac.serviceAccount    — ServiceAccount + ImagePullSecrets + 자동 token 비활성.
  - keiailab.rbac.controllerBase    — controller-runtime base RBAC (leader-election + events + service watch).
  - keiailab.rbac.workloadBase      — managed workload (StatefulSet/Deployment) + Service/ConfigMap/Secret RBAC.

caller 가 *delta rule* (CRD-specific verb 등) 만 자체 yaml 보존, base RBAC 은 본 partial 위임.
*/}}


{{/*
keiailab.rbac.serviceAccount — 표준 ServiceAccount 정의.

caller 인자 (dict):
  - ctx:                    helm context (.)
  - name:                   SA 이름 (보통 fullname)
  - labels:                 caller chart 의 labels
  - imagePullSecrets:       []corev1.LocalObjectReference (optional)
  - automountToken:         bool (optional, default true — operator 가 SA token 사용)
  - annotations:            map (optional, 예: aws.amazon.com/eks-iam-role)

사용 예:

  {{ include "keiailab.rbac.serviceAccount" (dict
      "ctx" .
      "name" (include "valkey-operator.fullname" .)
      "labels" (include "valkey-operator.labels" .)
      "imagePullSecrets" .Values.imagePullSecrets) }}
*/}}
{{- define "keiailab.rbac.serviceAccount" -}}
{{- $ctx := .ctx -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .name }}
  namespace: {{ $ctx.Release.Namespace }}
  labels:
    {{- .labels | nindent 4 }}
  {{- with .annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- with .imagePullSecrets }}
imagePullSecrets:
  {{- toYaml . | nindent 2 }}
{{- end }}
{{- if hasKey . "automountToken" }}
automountServiceAccountToken: {{ .automountToken }}
{{- else }}
automountServiceAccountToken: true
{{- end }}
{{- end }}


{{/*
keiailab.rbac.controllerBase — controller-runtime 표준 base RBAC rules.

leader-election (coordination.k8s.io/leases) + events (events.k8s.io,
core/events) + 자체 service watch — controller manager 모든 operator
공통 의존.

본 partial 은 *PolicyRule list* 만 출력 — caller 가 ClusterRole /
Role 의 rules: 아래에 nindent 호출.

사용 예:

  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRole
  metadata:
    name: {{ include "valkey-operator.fullname" . }}-manager-role
  rules:
    {{ include "keiailab.rbac.controllerBase" . | nindent 4 }}
    # delta CRD-specific rules:
    - apiGroups: ["cache.keiailab.io"]
      resources: ["valkeys", "valkeyclusters"]
      verbs: ["*"]
*/}}
{{- define "keiailab.rbac.controllerBase" -}}
- apiGroups: ["coordination.k8s.io"]
  resources: ["leases"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["events"]
  verbs: ["create", "patch"]
- apiGroups: ["events.k8s.io"]
  resources: ["events"]
  verbs: ["create", "patch"]
- apiGroups: [""]
  resources: ["services"]
  verbs: ["get", "list", "watch"]
{{- end }}


{{/*
keiailab.rbac.workloadBase — managed workload (StatefulSet / Deployment)
+ 의존 리소스 (Service / ConfigMap / Secret) RBAC.

operator 가 reconcile 하는 workload 생성/갱신/삭제 권한. CRD-specific
rules 와 별개 — 모든 workload-managing operator 공통.

사용 예:

  - apiGroups: ["apps"]    # 본 partial 가 cover
  - apiGroups: [""]         # 본 partial 가 cover (services/configmaps/secrets)
  - apiGroups: ["cache.keiailab.io"]
    resources: ["valkeys"]  # delta — caller 자체 yaml
    verbs: ["*"]
*/}}
{{- define "keiailab.rbac.workloadBase" -}}
- apiGroups: ["apps"]
  resources: ["statefulsets", "deployments"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["services", "configmaps", "secrets", "serviceaccounts"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["persistentvolumeclaims"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
{{- end }}
