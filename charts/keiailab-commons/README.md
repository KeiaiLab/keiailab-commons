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

### 후속 PR 의 partial (별 PR)

| Helper | RFC-0019 | PR |
|---|---|---|
| `keiailab.security.podSecurityRestricted` | §3.4 | PR-B5 |
| `keiailab.rbac.serviceAccount` | §3.5 | PR-C1 |
| `keiailab.rbac.controllerBase` | §3.5 | PR-C1 |
| `keiailab.rbac.workloadBase` | §3.5 | PR-C1 |

## License

Apache-2.0

## Refs

- ADR-0005 (commons): RFC-0019 채택 — library chart 패키징.
- RFC-0019: helm partials library 표준.
- Plan §2 D14 (Sprint B PR-B2): operator-commons library chart 신설.
