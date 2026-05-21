# keiailab-commons (Helm library chart)

`operator-commons` 의 *Helm library chart*. downstream operator chart 들이
공통으로 사용할 partial helper 를 제공합니다.

## 사용

consumer chart 의 `Chart.yaml` 에 의존성 추가:

```yaml
dependencies:
  - name: keiailab-commons
    version: ~0.4.0
    repository: oci://ghcr.io/keiailab/charts
```

`helm dependency update` 후 partial 사용:

```yaml
metadata:
  labels:
    {{- include "keiailab.commonLabels" . | nindent 4 }}
```

## Provided Partials

| Helper | 용도 | 추가 시점 |
|---|---|---|
| `keiailab.commonLabels` | Helm 표준 공통 label set (Kubernetes Recommended Labels) | v0.1.0 |
| `keiailab.observability.serviceMonitor` | Prometheus Operator ServiceMonitor 공통 spec | v0.1.0 |
| `keiailab.networkpolicy.dataplane` | managed dataplane workload 보호 (default-deny + allow-internal-instance) | v0.2.0 |
| `keiailab.networkpolicy.controlplane` | operator manager pod 자체 보호 (metrics / webhook ingress + API / DNS egress) | v0.2.0 |
| `keiailab.rbac.serviceAccount` | ServiceAccount + ImagePullSecrets + automountToken | v0.3.0 |
| `keiailab.rbac.controllerBase` | controller-runtime base (leader-election + events + service) | v0.3.0 |
| `keiailab.rbac.workloadBase` | managed workload (StatefulSet / Deployment + Service / ConfigMap / Secret / PVC / Pod) | v0.3.0 |
| `keiailab.security.podSecurityContext` | PSS Restricted Pod SecurityContext (`runAsNonRoot` + `seccompProfile`) | v0.4.0 |
| `keiailab.security.containerSecurityContext` | PSS Restricted Container SecurityContext (`capabilities.drop=ALL` + `readOnlyRootFilesystem`) | v0.4.0 |

## License

Apache-2.0 — `../../LICENSE` 참조.

## Refs

- ADR-0005 (commons): Helm library chart 패키징 결정.
- ADR-0006 / ADR-0007 / ADR-0008: networkpolicy / rbac / security partial 채택.
