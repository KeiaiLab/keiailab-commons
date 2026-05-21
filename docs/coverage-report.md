# Coverage Report — v1.0.0 graduation B.11.3

> `go test ./pkg/... -coverprofile=cover.out` 측정 결과. v1.0.0 졸업 조건 #7 (`커버리지 ≥ 85%`) 충족.

## 측정 결과 (2026-05-14)

| Package | Coverage | Status |
|---|---:|---|
| `pkg/finalizer` | **96.3%** | ✅ |
| `pkg/labels` | **100.0%** | ✅ |
| `pkg/monitoring` | **97.6%** | ✅ |
| `pkg/networkpolicy` | **89.2%** | ✅ |
| `pkg/security` | **100.0%** | ✅ |
| `pkg/status` | **87.5%** | ✅ |
| `pkg/version` | **100.0%** | ✅ |
| `pkg/webhook` | **97.2%** | ✅ |
| **합계 (statements)** | **96.3%** | **✅ PASS (≥85%)** |

모든 8 package 가 v1.0.0 threshold (85%) 통과. 

## 재현 명령

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.3%
```

## 향후 유지

- CI 게이트: `pre-push` hook 으로 매 push 마다 coverage 측정
- Regression alarm: 85% 이하 떨어지면 PR block
- 신규 package 추가 시 동일 ≥85% 기준

## Refs

- ROADMAP.md "v1.0.0 졸업 조건" #7
- docs/STABILITY.md graduation criterion #6
