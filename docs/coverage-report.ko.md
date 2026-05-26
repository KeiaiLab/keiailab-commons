# 커버리지 보고서

> [English](coverage-report.md) | **한국어** | [日本語](coverage-report.ja.md) | [中文](coverage-report.zh.md)

> `go test ./pkg/... -coverprofile=cover.out` 측정 결과.
> v1.0.0 졸업 조건: 커버리지 ≥ 85 %.

## 측정 결과 (2026-05-26)

| 패키지 | 커버리지 | 상태 |
|---|---:|---|
| `pkg/events` | **100.0%** | ✅ |
| `pkg/finalizer` | **96.3%** | ✅ |
| `pkg/labels` | **100.0%** | ✅ |
| `pkg/monitoring` | **97.6%** | ✅ |
| `pkg/networkpolicy` | **89.2%** | ✅ |
| `pkg/probes` | **100.0%** | ✅ |
| `pkg/pvc` | **90.2%** | ✅ |
| `pkg/security` | **100.0%** | ✅ |
| `pkg/status` | **87.5%** | ✅ |
| `pkg/storageclass` | **100.0%** | ✅ |
| `pkg/topology` | **100.0%** | ✅ |
| `pkg/version` | **100.0%** | ✅ |
| `pkg/webhook` | **97.2%** | ✅ |
| **합계 (statements)** | **96.2%** | **✅ PASS (≥85%)** |

전체 13개 패키지가 v1.0.0 임계값(85%)을 초과합니다.

## 재현 명령

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.2%
```

## 유지 관리

- 품질 게이트: `pre-push` hook이 매 push마다 커버리지를 측정합니다.
- 회귀 경보: 패키지가 85% 미만으로 하락하면 PR을 차단합니다.
- 신규 패키지도 동일한 ≥ 85% 기준을 충족해야 합니다.

## 참조

- [ROADMAP.md](ROADMAP.md) — v1.0.0 졸업 조건 #7.
- [STABILITY.md](STABILITY.md) — 졸업 조건 #6.
