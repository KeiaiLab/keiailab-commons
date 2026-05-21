# Upgrading operator-commons

본 문서는 `github.com/keiailab/operator-commons` Go module 의 마이너/메이저 버전
업그레이드 시 필요한 마이그레이션 작업을 정리한다. 본 라이브러리를 import 하는
3 operator (postgres / mongodb / valkey) 의 *공통 진입점*.

## 0. 버전 정책 (semver)

| 변경 유형 | semver bump | 예시 |
|---|---|---|
| 신규 패키지 추가 | minor (v0.X → v0.X+1) | pkg/reconcile, pkg/resources 신설 (Sprint 1) |
| 기존 API 시그니처 변경 (breaking) | major (v0.X → v1.0 / v1.X → v2.0) | pkg/status.SetReady() 시그니처 변경 |
| 패키지 *내부* 동작 변경 (non-breaking) | patch (v0.X.Y → v0.X.Y+1) | bug fix |
| ADR 일탈 결정 | major + Deprecated 안내 | API stability Tier 변경 |

API stability Tier (각 pkg/<name>/doc.go 의 marker):
- **Stable** — 1 minor 안 backward-compat 보장
- **Beta** — 다음 minor 에서 변경 가능
- **Alpha** — 언제든 변경 (실험 단계)

## 1. v0.7.x → v0.8.x

### Helm library chart 사용자

```bash
helm dep update charts/<your-operator>
helm template <your-operator> charts/<your-operator> --set features.cluster.enabled=true
```

`keiailab-commons` chart v0.8.0 의 `_servicemonitor.tpl` / `_rbac.tpl` /
`_networkpolicy.tpl` partial 사용 시 추가 작업 없음. ADR-0006 + ADR-0007 + ADR-0008
정합.

### Go module 사용자

```bash
go get github.com/keiailab/operator-commons@v0.8.0
go mod tidy
```

추가 작업:
- (없음 — backward-compat)

## 2. v0.8.x → v0.9.x (예정 — Sprint 1 + S5 cycle 결과)

### 신규 패키지 추가 (minor bump)

| 패키지 | 목적 | Tier | 추출 출처 |
|---|---|---|---|
| `pkg/pvc` | PVC 헬퍼 (확장, topology) | Beta | 3 operator 중복 (~495 LOC, PR #52) |
| `pkg/topology` | PVC topology 추출 | Beta | 3 operator 중복 |
| `pkg/reconcile` (계획) | HandleFinalizerCleanup + Statusable + UpdateStatusWithRetry + RequeueIntervals + IsPaused | Beta | S5 cycle |
| `pkg/resources` (계획) | Apply{ConfigMap,Service,STS,NetworkPolicy,PDB,HPA,ServiceMonitor} + Upsert generic | Beta | S5 cycle |

### Migration

3 operator (postgres / mongodb / valkey) 의 import path 추가:

```go
import (
    "github.com/keiailab/operator-commons/pkg/pvc"
    "github.com/keiailab/operator-commons/pkg/topology"
    // 후속: pkg/reconcile, pkg/resources
)
```

### Backward-compat

- 기존 `pkg/status` / `pkg/finalizer` / `pkg/networkpolicy` / `pkg/monitoring` / `pkg/probes` / `pkg/labels` / `pkg/storageclass` / `pkg/webhook` 시그니처 변경 없음
- Helm library chart `keiailab-commons` v0.9.0 의 `_security.tpl` / `_servicemonitor.tpl` 추가 사용은 *opt-in* (기존 인라인 정의 그대로 두면 영향 없음)

### 권장 마이그레이션 절차 (3 operator 공통)

```bash
# 1. 의존 bump
go get github.com/keiailab/operator-commons@v0.9.0
go mod tidy

# 2. 중복 코드 제거 (sub-spec 별 단계)
# 예: helpers.go 의 HandleFinalizerCleanup → pkg/reconcile.HandleFinalizerCleanup
# 예: resources_apply.go 의 ApplyConfigMap → pkg/resources.ApplyConfigMap

# 3. 검증
make verify  # lint + test + build + audit
make integration-test  # envtest

# 4. e2e (kind)
kind create cluster
helm install <operator> charts/<operator>
kubectl apply -f config/samples/
kubectl get <CR> -A  # reconcile 결과 확인
```

## 3. v0.9.x → v1.0.0 (예정 — v3.x-stable 선언 시점)

CLAUDE.md §7 의 *상용 제품 수준* (P0+P1+P2+OP+C 모두 ✅) 도달 시.

- 모든 API stability Tier `Stable` 승격
- breaking change 없음 (v0.x → v1.0 은 *명명만* 변경, semantic 동일)
- 5 repo 의 일관성 보장: `docs/quality/production-grade-checklist.md` 의 모든 항목 ✅

상세: ADR-0013 (audit-production-grade.sh) 참조.

## 4. 일반 마이그레이션 체크리스트

업그레이드 전:
- [ ] `go mod tidy` 후 `go.mod` 변경 없음 (drift 차단)
- [ ] `make audit` 통과 (govulncheck CVE 0)
- [ ] 기존 e2e 스위트 PASS
- [ ] 의존 3 operator 의 `vendor/` 또는 `go.sum` 갱신 검토

업그레이드 후:
- [ ] 3 operator 의 import path 변경 (`go get -u` 또는 명시 버전)
- [ ] 각 operator 의 `make verify` 통과
- [ ] 각 operator 의 e2e PASS
- [ ] Helm chart `charts/<operator>` 의 `dependencies:` 갱신 (keiailab-commons 의 chart 버전)

## 5. 비호환 변경 안내 정책

- **Deprecation**: 신규 minor 에서 `// Deprecated:` 주석 + 2 minor 후 제거
- **Breaking**: major bump + 본 UPGRADING.md 의 별 섹션 + ADR 작성
- **사후 통보 안 함**: 모든 breaking 변경은 *최소 1 minor* 사전 deprecation 거침

## 참고

- ADR 목록: `docs/kb/adr/INDEX.md`
- API stability: `pkg/<name>/doc.go` 의 Tier marker
- audit: `make audit-quality` (5 repo 측정, ADR-0013)
- i18n: `docs/i18n/README.md` (4-lang 정책 SSOT)
