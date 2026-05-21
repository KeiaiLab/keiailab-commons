<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **Kubernetes operator 공통 scaffolding 을 위한 Go 라이브러리 — finalizer / labels / status / version / security / monitoring partials.**
>
> [English](README.md) | **한국어** | [日本語](README.ja.md) | [中文](README.zh.md)

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <b>한국어</b> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>

---

Kubernetes operator 코드베이스의 scaffolding drift 를 제거하는 재사용 가능한
Go 라이브러리입니다 — PodSecurity restricted context, 지원 버전 allowlist,
NetworkPolicy 템플릿, ServiceMonitor 빌더, finalizer / status 헬퍼,
Helm library chart partial 을 작고 안정된 API 표면 뒤에 묶어 제공합니다.

> 상태: **v0.x — API 변경 가능.** v1.0 부터 SemVer stable.

## Why

Operator 작성자는 동일한 scaffolding 을 반복적으로 구현합니다 — restricted
PodSecurity context, 지원 버전 매트릭스, default-deny NetworkPolicy,
ServiceMonitor 빌더, finalizer 헬퍼, status condition 카탈로그. 이를 독립적으로
재구현하면 비슷한 reconciler 간 조용한 불일치가 생기고, minor 리비전을
거치며 점차 분기됩니다. `operator-commons` 는 이 scaffolding 의 단일 소스 입니다 —
헬퍼를 import 하고 canonical 구현을 받아, 모든 저장소에서 다시 발명하지 않습니다.

## 패키지

| 패키지 | Tier | 목적 |
|---|---|---|
| `pkg/finalizer` | Stable | Finalizer 헬퍼 — `Add` / `Remove` / `Has` / `EnsureOrder` (stdlib `slices` 만, controller-runtime 의존 없음). |
| `pkg/labels` | Stable | Kubernetes 권장 라벨 (`app.kubernetes.io/*`) 빌더 — `Set`, `All()`, `Selector()`, v2 매핑 (`AllV2`). |
| `pkg/status` | Stable | 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼 (`SetReady`, `SetAvailable`, `SetReadyFalse`). |
| `pkg/storageclass` | Stable | DNS-1123 storageClass validator + `Normalize` / `MustNormalize` (empty → cluster default 포인터). |
| `pkg/version` | Beta | 버전 allowlist 규약 (`MustList`, `IsSupported`, `Strings`, `Default`) + generic `Matrix[E MatrixEntry]` + serializer. |
| `pkg/monitoring` | Beta | Prometheus Operator `ServiceMonitor` 및 `PrometheusRule` 빌더 (unstructured — CRD-soft). |
| `pkg/networkpolicy` | Beta | Deny-by-default NetworkPolicy 빌더 + functional options (`WithSelfIngress`, `WithIngressFromPeers`, `WithDenyEgress`, `WithEgressToPeers`, `ComboPeer`). |
| `pkg/security` | Beta | PodSecurity *restricted* SecurityContext 빌더 + Pod / Container 분리 + seccomp profile 포인터. |
| `pkg/events` | Beta | 최소 `Recorder` 인터페이스 + 9 표준 `Reason` 상수 + `Emit` / `EmitWarning` / `WrappedError` (nil-safe). |
| `pkg/probes` | Experimental | `corev1.Probe` fluent 빌더 — HTTP / HTTPS / TCP / Exec, kubelet default + clamp. |
| `pkg/webhook` | Experimental | Admission validation 헬퍼 — `ValidateAllowedVersion`, `ValidateWithPredicate`, conversion registry. |

[docs/STABILITY.md](docs/STABILITY.md) 가 tier 약속을 정의합니다.
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) 는 패키지 표면과 설계 불변식을
설명합니다. [docs/ROADMAP.md](docs/ROADMAP.md) 는 tier 격상 기준과 v1.0 졸업
체크리스트를 추적합니다.

## 사용 예

```go
import (
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/version"
    corev1 "k8s.io/api/core/v1"
)

var supportedVersions = version.MustList("1.0", "1.1", "1.2")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

각 패키지의 자세한 예제는 `pkg/<name>/doc.go` 의 패키지 문서에 있습니다
(`go doc github.com/keiailab/operator-commons/pkg/<name>`).

## 버저닝 및 릴리스

- **v0.x**: API breaking 허용. 각 태그 (`v0.N.M`) 는 패키지의 공개 API
  또는 의미 있는 동작 변경을 동반합니다. 소비자는 `go.mod` 로 특정 버전을 핀합니다.
- **v1.0 이후**: Semantic Versioning. Breaking change 는 ADR (`docs/kb/adr/`) 필수.
- 로컬 `replace` directive 는 cross-repo 개발에 허용; 릴리스 태그는 항상
  canonical module path 를 보존합니다.

## 커뮤니티

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — 패키지 API 질문, integration 사례, 새 helper 제안.
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — 버그 및 구체적 기능 요청.
- **Security**: 비공개 공시 절차는 [SECURITY.md](SECURITY.md) 참조.
- **Contributing**: 개발 워크플로는 [CONTRIBUTING.md](CONTRIBUTING.md) 참조.

## 라이선스

Apache-2.0 — [LICENSE](LICENSE) 참조. AGPL / BUSL transitive 의존성 0건 목표
(매 minor 릴리스마다 감사).

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
