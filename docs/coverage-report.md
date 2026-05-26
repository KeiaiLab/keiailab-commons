# Coverage Report

> **English** | [한국어](coverage-report.ko.md) | [日本語](coverage-report.ja.md) | [中文](coverage-report.zh.md)

> `go test ./pkg/... -coverprofile=cover.out` measurement results.
> v1.0.0 graduation criterion: coverage ≥ 85 %.

## Results (2026-05-26)

| Package | Coverage | Status |
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
| **Total (statements)** | **96.2%** | **✅ PASS (≥85%)** |

All 13 packages exceed the v1.0.0 threshold (85%).

## Reproduce

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.2%
```

## Maintenance

- Quality gate: `pre-push` hook measures coverage on every push.
- Regression alarm: block PR when any package drops below 85 %.
- New packages must meet the same ≥ 85 % bar.

## References

- [ROADMAP.md](ROADMAP.md) — v1.0.0 graduation criterion #7.
- [STABILITY.md](STABILITY.md) — graduation criterion #6.
