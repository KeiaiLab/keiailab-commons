# Upgrading keiailab-commons

> [English](UPGRADING.md) | **한국어** | [日本語](UPGRADING.ja.md) | [中文](UPGRADING.zh.md)

본 문서는 `github.com/keiailab/keiailab-commons` Go 모듈의 minor / major
버전 업그레이드 시 필요한 마이그레이션 작업을 정리합니다. downstream
consumer 가 본 라이브러리를 import 할 때의 *공통 진입점* 입니다.

## 0. 버전 정책 (semver)

| 변경 유형 | semver bump | 예시 |
|---|---|---|
| 신규 패키지 추가 | minor (v0.X → v0.X+1) | `pkg/events`, `pkg/storageclass` 신설 |
| 기존 API 시그니처 변경 (breaking) | major (v0.X → v1.0 / v1.X → v2.0) | `pkg/status.SetReady()` 시그니처 변경 |
| 패키지 *내부* 동작 변경 (non-breaking) | patch (v0.X.Y → v0.X.Y+1) | bug fix |
| ADR 일탈 결정 | major + Deprecated 안내 | API stability tier 변경 |

API stability tier (`pkg/<name>/doc.go` 의 marker):

- **Stable** — minor 범위 backward-compat 보장.
- **Beta** — 다음 minor 에서 변경 가능.
- **Experimental** — 언제든 변경 (위험 부담은 호출자).

## 1. v0.7.x → v0.8.x

### Helm library chart 사용자

```bash
helm dep update charts/<your-operator>
helm template <your-operator> charts/<your-operator>
```

`keiailab-commons` chart v0.8.0 의 `_servicemonitor.tpl` / `_rbac.tpl` /
`_networkpolicy.tpl` partial 사용 시 추가 작업이 필요하지 않습니다.

### Go 모듈 사용자

```bash
go get github.com/keiailab/keiailab-commons@v0.8.0
go mod tidy
```

추가 작업 없음 — backward-compatible.

## 2. v0.8.x → v0.9.x

### 신규 패키지 추가 (minor bump)

| 패키지 | 목적 | Tier |
|---|---|---|
| `pkg/pvc` | PVC 헬퍼 (확장, in-place patch) | Beta |
| `pkg/topology` | PVC topology spread + zone-aware affinity | Beta |

### Migration

downstream operator 의 import path 추가:

```go
import (
    "github.com/keiailab/keiailab-commons/pkg/pvc"
    "github.com/keiailab/keiailab-commons/pkg/topology"
)
```

### Backward-compat

- 기존 패키지 (`pkg/status`, `pkg/finalizer`, `pkg/networkpolicy`,
  `pkg/monitoring`, `pkg/probes`, `pkg/labels`, `pkg/storageclass`,
  `pkg/webhook`, `pkg/events`, `pkg/security`, `pkg/version`) 시그니처
  변경 없음.
- Helm library chart `keiailab-commons` 의 `_security.tpl` /
  `_servicemonitor.tpl` 추가 사용은 *opt-in* — 기존 인라인 정의를 그대로
  두어도 영향 없음.

### 권장 마이그레이션 절차

```bash
# 1. 의존 bump
go get github.com/keiailab/keiailab-commons@v0.9.0
go mod tidy

# 2. 검증
make verify  # lint + test + build

# 3. e2e (kind)
kind create cluster
helm install <operator> charts/<operator>
kubectl apply -f config/samples/
kubectl get <CR> -A  # reconcile 결과 확인
```

## 3. v0.9.x → v0.10.x

### repository / module rename

`v0.10.0`부터 Go module path는 다음 경로를 사용합니다.

```bash
github.com/keiailab/keiailab-commons
```

downstream operator는 import path와 dependency pin을 함께 갱신합니다.

```bash
go get github.com/keiailab/keiailab-commons@v0.10.0
go mod tidy
```

기존 `v0.9.x` tag는 `github.com/keiailab/operator-commons` module path를
선언하므로 새 module path로 소비할 수 없습니다.

## 4. v0.9.x → v1.0.0

v1.0.0 졸업 조건 충족 시점에 진행합니다 ([STABILITY.md](STABILITY.md) "v1.0.0
graduation" 참조).

- 모든 패키지가 Stable tier 로 격상.
- v0.x → v1.0 은 *명명만* 변경, semantic 동일 (breaking change 없음).

## 5. 일반 마이그레이션 체크리스트

업그레이드 전:

- [ ] `go mod tidy` 후 `go.mod` 변경 없음 (drift 차단).
- [ ] `make audit` 통과 (govulncheck CVE 0).
- [ ] 기존 e2e 스위트 PASS.

업그레이드 후:

- [ ] downstream operator 의 import path 변경 (`go get -u` 또는 명시 버전).
- [ ] `make verify` 통과.
- [ ] e2e PASS.
- [ ] Helm chart `charts/<operator>` 의 `dependencies:` 갱신.

## 6. 비호환 변경 안내 정책

- **Deprecation**: 신규 minor 에서 `// Deprecated:` 주석 + 2 minor 후 제거.
- **Breaking**: major bump + 본 UPGRADING.md 의 별 섹션 + ADR 작성.
- **사후 통보 안 함**: 모든 breaking 변경은 *최소 1 minor* 사전 deprecation
  거침.

## 참고

- ADR 목록: [`docs/kb/adr/INDEX.ko.md`](kb/adr/INDEX.ko.md).
- API stability: `pkg/<name>/doc.go` 의 Tier marker.
- i18n: [`docs/i18n/README.ko.md`](i18n/README.ko.md) (다국어 정책).
