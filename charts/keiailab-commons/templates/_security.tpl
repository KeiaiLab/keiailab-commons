{{/*
keiailab-commons — PodSecurity Restricted partials (library chart partial 표준).

기반: commons-ADR-0008.

K8s Pod Security Standards "restricted" profile (PSS Restricted) 의
표준 SecurityContext / ContainerSecurityContext 출력. downstream consumer (mongodb /
postgres / valkey operator) 의 *manager pod* 보안 표준화.

reference: https://kubernetes.io/docs/concepts/security/pod-security-standards/

두 partial:
  - keiailab.security.podSecurityContext       — Pod-level SecurityContext.
  - keiailab.security.containerSecurityContext — Container-level SecurityContext.
*/}}


{{/*
keiailab.security.podSecurityContext — PSS Restricted Pod SecurityContext.

caller 인자 (dict):
  - runAsUser:    UID (default 65532 — distroless nonroot 표준)
  - runAsGroup:   GID (default 65532)
  - fsGroup:      fsGroup (default 65532)
  - override:     사용자 정의 PodSecurityContext (raw map, optional —
                  PodSecurityRestricted=false 시 사용자 override)

사용 예 (consumer chart 의 templates/deployment.yaml — pod spec):

  spec:
    template:
      spec:
        securityContext:
          {{- include "keiailab.security.podSecurityContext" (dict
              "runAsUser" 65532
              "runAsGroup" 65532
              "fsGroup" 65532) | nindent 10 }}

또는 사용자 override (v1alpha2 PodSecurityRestricted=false 시):

  spec:
    template:
      spec:
        securityContext:
          {{- include "keiailab.security.podSecurityContext" (dict
              "override" .Values.podSecurityContext) | nindent 10 }}
*/}}
{{- define "keiailab.security.podSecurityContext" -}}
{{- if .override -}}
{{- toYaml .override }}
{{- else -}}
runAsNonRoot: true
runAsUser: {{ .runAsUser | default 65532 }}
runAsGroup: {{ .runAsGroup | default 65532 }}
fsGroup: {{ .fsGroup | default 65532 }}
seccompProfile:
  type: RuntimeDefault
{{- end -}}
{{- end }}


{{/*
keiailab.security.containerSecurityContext — PSS Restricted Container
SecurityContext.

caller 인자 (dict):
  - runAsUser:                 UID (default 65532)
  - runAsGroup:                GID (default 65532)
  - readOnlyRootFilesystem:    bool (default true)
  - override:                  사용자 정의 SecurityContext (raw map, optional)

사용 예 (consumer chart 의 templates/deployment.yaml — container spec):

  containers:
    - name: manager
      image: ...
      securityContext:
        {{- include "keiailab.security.containerSecurityContext" (dict
            "runAsUser" 65532
            "runAsGroup" 65532) | nindent 8 }}
*/}}
{{- define "keiailab.security.containerSecurityContext" -}}
{{- if .override -}}
{{- toYaml .override }}
{{- else -}}
runAsNonRoot: true
runAsUser: {{ .runAsUser | default 65532 }}
runAsGroup: {{ .runAsGroup | default 65532 }}
readOnlyRootFilesystem: {{ default true .readOnlyRootFilesystem }}
allowPrivilegeEscalation: false
capabilities:
  drop:
    - ALL
seccompProfile:
  type: RuntimeDefault
{{- end -}}
{{- end }}
