<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# keiailab 오퍼레이터 패밀리

> ⚠️ This translation is AI-generated and pending native review. — 본 번역은 Claude 기계 번역 결과입니다. 모국어 검토자의 검수 전까지 `[검토 필요]` 상태입니다.

> 공통 기반 위에 구축된 4개의 자매 Kubernetes 오퍼레이터 — `operator-commons` (Go 라이브러리) + Helm partial + Apache-2.0 스택.

본 페이지는 **`operator-commons`** 저장소에서 읽고 계십니다. 본 페이지는 전체 패밀리의 canonical cross-link 입니다.

## 패밀리 개요

| 프로젝트 | 데이터베이스 | 상태 | 저장소 |
|---|---|---|---|
| **`postgres-operator`** | PostgreSQL 18+ | active | https://github.com/keiailab/postgres-operator |
| **`mongodb-operator`** | MongoDB 7.0+ | active | https://github.com/keiailab/mongodb-operator |
| **`valkey-operator`** | Valkey 8.0+ (Redis fork, BSD-3) | active | https://github.com/keiailab/valkey-operator |
| **`operator-commons`** | 공유 Go 라이브러리 | **v0.8.0** (현재 페이지) | https://github.com/keiailab/operator-commons |

## 공유하는 것

4 프로젝트 모두 동일한 운영 primitive 로 수렴합니다:

- **Apache-2.0** end-to-end — SSPL 없음, SaaS 표면 copyleft 없음
- **`operator-commons`** 공유 Go 라이브러리 (v0.8.0+) — 파이널라이저, 라벨, 상태 슈가, security context 빌더, NetworkPolicy / ServiceMonitor partial
- **Helm chart skeleton** — RFC-0027 `default` falsy-toggle 차단, RFC-0026 컴포넌트 키 values, cycle 26 hardening 6 marker (priorityClassName / lifecycle / SA / minReadySeconds / automount / revisionHistoryLimit)
- **OLM bundle parity** — scorecard v1alpha3 6-test matrix
- **i18n** — README + canonical docs 영문 / 한국어 / 日本語 / 中文 (cleanup supercycle 2026-05-21 의 Wave 4)

## 패밀리 안에서 `operator-commons` 의 역할

본 저장소는 **공유 Go 라이브러리** 입니다 — 컨트롤러가 *아닙니다*. 다음을 제공합니다:

| 패키지 | 목적 | Tier |
|---|---|---|
| `pkg/finalizer` | std `slices` 만 사용하는 `Add` / `Remove` / `Has` 파이널라이저 헬퍼 | **Stable** |
| `pkg/labels` | 권장 K8s 라벨 빌더 — `Set`, `All()`, `Selector()` | **Stable** |
| `pkg/status` | 4 표준 Condition Type + 6 Reason catalog + 헬퍼 | **Stable** |
| `pkg/version` | DB 버전 allowlist convention + generic `Matrix[E MatrixEntry]` | Beta |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 빌더 (unstructured) | Beta |
| `pkg/networkpolicy` | Deny-by-default NetworkPolicy 빌더 + functional option | Beta |
| `pkg/security` | PodSecurity *restricted* SecurityContext 빌더 | Beta |
| `pkg/webhook` | Admission validation 헬퍼 | Experimental |

설계 invariant: **leaf 패키지는 stdlib + k8s API 타입만**. controller-runtime 없음, logr 없음, operator-sdk leak 없음.

자세한 패키지 surface 는 [ARCHITECTURE.md](ARCHITECTURE.md), tier 격상 기준은 [ROADMAP.md](ROADMAP.md) 참조.

## 우리가 하지 *않는* 것

- ❌ **upstream 오퍼레이터 임베드 또는 wrap** (PGO, CloudNativePG, MongoDB Community Operator, Sentinel) — license-clean, copyleft 의무 없음
- ❌ **GitHub Actions 를 release gate 로 사용** — 로컬 4계층 + GitLab CI L5 (RFC-0002, RFC-0043 참조)
- ❌ **시간 기반 로드맵 마감일** — 기능 체크리스트 + 완료 비율 (`standards/roadmap.md §1.1` 참조)
- ❌ **Bitnami chart / image** — registry deprecation 위험, Broadcom 인수 (ADR-0136 / ADR-0057 참조)
- ❌ **본 저장소의 CRD / Reconciler** — consumer operator 가 해당 책임 소유

## 시작하는 곳

| 작업 | 진입점 |
|---|---|
| 오퍼레이터에서 `operator-commons` 가져오기 | [README.md](../README.md) Usage 섹션 |
| 아키텍처 읽기 | [ARCHITECTURE.md](ARCHITECTURE.md) |
| 이슈 또는 기능 요청 등록 | https://github.com/keiailab/operator-commons/issues |
| 설계 또는 로드맵 논의 | https://github.com/keiailab/operator-commons/discussions |
| 코드 기여 | [CONTRIBUTING.md](../CONTRIBUTING.md) |
| 보안 이슈 보고 | [SECURITY.md](../SECURITY.md) |
| 브랜드 / 보이스 학습 | [BRANDING.md](BRANDING.md) |
| 도입자 추적 / 사용자 확인 | [ADOPTERS.md](ADOPTERS.md) |
| 메인테이너 찾기 | [MAINTAINERS.md](MAINTAINERS.md) |
| 거버넌스 모델 검토 | [GOVERNANCE.md](GOVERNANCE.md) |
| 진행 예정 작업 확인 | [ROADMAP.md](ROADMAP.md) |
| API 안정성 약속 검토 | [docs/STABILITY.md](STABILITY.md) |

## 패밀리 간 호환성

3 데이터베이스 오퍼레이터 모두 `github.com/keiailab/operator-commons` 를 매칭 버전 (현재 `v0.8.0+`) 으로 가져옵니다:

```go
import (
    "github.com/keiailab/operator-commons/pkg/version"
    "github.com/keiailab/operator-commons/pkg/security"
    "github.com/keiailab/operator-commons/pkg/labels"
    "github.com/keiailab/operator-commons/pkg/monitoring"
    "github.com/keiailab/operator-commons/pkg/finalizer"
    "github.com/keiailab/operator-commons/pkg/status"
)
```

`operator-commons` 의 breaking change 는 3 데이터베이스 오퍼레이터 모두에서 동기 bump 필요 — supercycle Wave 5 의 `make cross-validation` target 으로 검증.

라이브 consumer matrix (3 operator × 8 package × 도입 %) 는 [ADOPTERS.md](ADOPTERS.md) 참조.

## i18n

본 페이지 (및 모든 canonical 프로젝트 문서) 는 4개 언어로 제공됩니다:

- **English** (canonical, 원본 파일)
- [한국어](family.ko.md) (현재 파일)
- [日本語](family.ja.md)
- [中文](family.zh.md)

의심스러울 때는 기술 콘텐츠는 영문 버전이 권위 있으며, 현지화 버전은 동일한 결정을 모국어 표현으로 반영합니다.

---

<p align="center">
  <b>keiailab 오퍼레이터 패밀리</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
