# keiailab-commons (Helm library chart)

operator-commons 의 *Helm library chart*. 4-repo (mongodb-operator /
postgres-operator / valkey-operator) 가 공유하는 partial helper 를 제공.

Plan §2 D14 (Sprint B PR-B2) — RFC-0019 §3.1 implementation.

## 사용

consumer chart 의 `Chart.yaml` 에 의존성 추가:

```yaml
dependencies:
  - name: keiailab-commons
    version: ~0.1.0
    repository: oci://ghcr.io/keiailab/charts
```

`helm dependency update` 후 partial 사용:

```yaml
metadata:
  labels:
    {{- include "keiailab.commonLabels" . | nindent 4 }}
```

## Provided Partials

| Helper | 용도 | RFC-0019 | 추가 시점 |
|---|---|---|---|
| `keiailab.commonLabels` | Helm 표준 공통 label set (Kubernetes Recommended Labels) | §3.1 | v0.1.0 |
| `keiailab.observability.serviceMonitor` | Prometheus Operator ServiceMonitor 공통 spec | §3.1 | v0.1.0 |
| `keiailab.networkpolicy.dataplane` | managed dataplane workload 보호 (default-deny + allow-internal-instance) | §3.2 | v0.2.0 |
| `keiailab.networkpolicy.controlplane` | operator manager pod 자체 보호 (metrics/webhook ingress + API/DNS egress) | §3.2 | v0.2.0 |
| `keiailab.rbac.serviceAccount` | ServiceAccount + ImagePullSecrets + automountToken | §3.5 | v0.3.0 |
| `keiailab.rbac.controllerBase` | controller-runtime base (leader-election + events + service) | §3.5 | v0.3.0 |
| `keiailab.rbac.workloadBase` | managed workload (StatefulSet/Deployment + Service/ConfigMap/Secret/PVC/Pod) | §3.5 | v0.3.0 |
| `keiailab.security.podSecurityContext` | PSS Restricted Pod SecurityContext (runAsNonRoot + seccompProfile) | §3.4 | v0.4.0 |
| `keiailab.security.containerSecurityContext` | PSS Restricted Container SecurityContext (capabilities.drop=ALL + readOnlyRootFilesystem) | §3.4 | v0.4.0 |

### RFC-0019 implementation 완결 (v0.4.0 시점)

§3.1 / §3.2 / §3.4 / §3.5 모두 implementation 완료. consumer chart (mongodb / postgres / valkey operator) 의 partial include 작업 후속.

## License

Apache-2.0

## Refs

- ADR-0005 (commons): RFC-0019 채택 — library chart 패키징.
- RFC-0019: helm partials library 표준.
- Plan §2 D14 (Sprint B PR-B2): operator-commons library chart 신설.
