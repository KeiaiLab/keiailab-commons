<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab operator 가 공유하는 Go 라이브러리 — finalizer / labels / status / version / security / monitoring partials**
>
> [English](README.md) | **한국어** | [日本語](README.ja.md) | [中文](README.zh.md)

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go" alt="Go Version"/></a>
  <a href="https://pkg.go.dev/github.com/keiailab/operator-commons"><img src="https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg" alt="Go Reference"/></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons"><img src="https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge" alt="OpenSSF Scorecard"/></a>
  <a href="https://github.com/keiailab/operator-commons/discussions"><img src="https://img.shields.io/github/discussions/keiailab/operator-commons?label=discussions&logo=github" alt="GitHub Discussions"/></a>
  <a href="https://github.com/keiailab/operator-commons/blob/main/docs/quality/audit-history.md"><img src="https://img.shields.io/badge/keiailab-v3.x--stable-success?style=flat-square" alt="keiailab v3.x-stable"/></a>
  <a href="https://github.com/keiailab/operator-commons/blob/main/scripts/audit-production-grade.sh"><img src="https://img.shields.io/badge/audit-100%25-success?style=flat-square" alt="audit"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <b>한국어</b> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>

---

**keiailab** Kubernetes operator (`mongodb-operator`, `valkey-operator`, `postgresql-operator`) 가 공유하는 Go 라이브러리입니다.

> 상태: **v0.x — API 가 깨질 수 있음**. v1.0 이후부터 SemVer (의미론적 버전 관리) stable.

## Why (왜 만들었는가)

3개 operator 가 *독립적으로* 동일한 scaffolding 코드 (PodSecurity restricted context, 버전 allowlist, NetworkPolicy template, ServiceMonitor builder) 를 구현해 왔습니다. repo 간 드리프트가 이미 일관성 불일치를 만들고 있어 — 본 라이브러리가 단일 진본 (single source of truth) 입니다.

## Packages (v0.8.0)

| Package | Purpose |
|---|---|
| `pkg/version` | 지원 DB 버전 allowlist 컨벤션 (`MustList`, `IsSupported`, `Strings`, `Default`) + 제네릭 `Matrix[E MatrixEntry]`. |
| `pkg/security` | functional options 패턴의 PodSecurity *restricted* SecurityContext 빌더. |
| `pkg/labels` | 권장 Kubernetes 레이블 (`app.kubernetes.io/*`) 빌더 — `Set`, `All()`, `Selector()` (버전 인식 분기). |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 빌더 (unstructured — CRD-soft 방식). |
| `pkg/networkpolicy` | NetworkPolicy 빌더 — deny-by-default + functional options (`WithSelfIngress`, `WithIngressFromPeers`, `WithDenyEgress`, `WithEgressToPeers`). |
| `pkg/webhook` | Admission 검증 헬퍼 — `ValidateAllowedVersion` (완전 일치), `ValidateWithPredicate` (호출자 제공 매처, 예: semver-prefix). |
| `pkg/finalizer` | Finalizer 헬퍼 — `Add`/`Remove`/`Has` (controller-runtime 의존 회피, 표준 `slices` 만 사용). |
| `pkg/status` | 4 표준 Condition Type + 6 Reason 카탈로그 + 헬퍼 (`SetReady`, `SetAvailable`, `SetReadyFalse`). |

`pkg/conditions` 는 *upstream `k8s.io/apimachinery/pkg/api/meta.SetStatusCondition` 활용 권장* (commons 미추가 결정 — boundary 분석 결과, 자세히는 mongodb-operator HANDOFF iteration 32 참조).

## Adoption Matrix (3 operator)

| Operator | sec | ver | lab | mon | np | wh | 채택률 |
|---|---|---|---|---|---|---|---|
| [mongodb-operator](https://github.com/keiailab/mongodb-operator) | ✅ | ✅ | ✅ | ⏳ | ✅ | ⏳ | **4/6 (67%)** |
| [valkey-operator](https://github.com/keiailab/valkey-operator) | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **6/6 (100%)** 🎉 |
| [postgres-operator](https://github.com/keiailab/postgres-operator) | ✅ | ⏳ | ✅ | ⏳ | ⏳ | ✅ | **3/6 (50%)** |

valkey-operator 가 *최초 100% 채택* — 다른 operator 의 carbon-copy reference 역할을 합니다. 적용 사례 commits:
- `pkg/security`: it8 (3 operator cross-cut) — `23fd3da` mongodb / `a0be4cf` valkey / `ac2e647` postgres
- `pkg/version`: mongodb it9 `a8db040`, valkey it8
- `pkg/labels`: mongodb it27 `ebc5803`, postgres it28 `c68b451`, valkey it29 `e8428b1`
- `pkg/monitoring`: valkey it23 `1765b54`
- `pkg/networkpolicy`: valkey it25 `97162b5`, mongodb it26 `ca0ec27`
- `pkg/webhook`: valkey it31 `14be0db`, postgres it34 `1d8fa17`

⏳ 영역은 *기능 추가 동반* (예: mongodb webhook server / ServiceMonitor reconciler) 또는 *별 추상화 적합* (postgres version matrix.go 의 Combo struct 가 commons.MustList 보다 풍부 — 위임 부적합) 으로 *deepening 보류* 상태입니다.

## Usage (사용 방법)

```go
import (
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/version"
)

var SupportedMongoDBVersions = version.MustList("8.0", "8.2", "8.3")

func buildContainerSecurityContext() *corev1.SecurityContext {
    return security.RestrictedContainer(
        security.WithRunAsUser(999),
        security.WithRunAsGroup(999),
    )
}
```

## Versioning + Release (버전 관리 및 릴리스)

- v0.x: API 파괴적 변경 허용. 각 태그 (`v0.N.M`) 는 패키지, public API, 또는 주요 동작 변경 시 발행됩니다.
- 각 호출자 operator 는 `go.mod` `require` 로 버전을 고정합니다 — 본 repo 와 3개 operator 를 함께 로컬 개발할 때는 `replace` 디렉티브 사용도 허용됩니다.
- v1.0 이후: Semantic Versioning. 파괴적 변경은 RFC 절차를 거쳐야 합니다.

## Community (커뮤니티)

- **Discussions**: [GitHub Discussions](https://github.com/keiailab/operator-commons/discussions) — pkg API 질문, integration 사례, 새 helper 제안
- **Issues**: [GitHub Issues](https://github.com/keiailab/operator-commons/issues) — 버그 / API 요청
- **Downstream**: 3 operator (mongodb-operator / postgres-operator / valkey-operator) — `go.mod replace` 또는 직접 `require` 로 사용
- **Stability matrix**: `pkg/labels`, `pkg/security`, `pkg/version`, `pkg/webhook` (stable / 안정, v0.5+) / `pkg/networkpolicy`, `pkg/monitoring` (experimental / 실험적)

## License (라이선스)

Apache-2.0 — [LICENSE](./LICENSE) 참조. minor 릴리스마다 감사하는 AGPL/BUSL 전이 의존성 제로 목표를 유지합니다.

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a> ·
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/forgewise">forgewise</a>
</p>

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
