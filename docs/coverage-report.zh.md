# Coverage Report — v1.0.0 graduation B.11.3

> ⚠️ This translation is AI-generated and pending native review. — 本翻译为 Claude 机器翻译结果。

> `go test ./pkg/... -coverprofile=cover.out` 测量结果。v1.0.0 毕业条件 #7 (`覆盖率 ≥ 85%`) 满足。

## 测量结果 (2026-05-14)

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
| **合计 (statements)** | **96.3%** | **✅ PASS (≥85%)** |

所有 8 个 package 通过 v1.0.0 threshold (85%).

## 重现命令

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.3%
```

## 后续维护

- CI 门控: `pre-push` hook 每次 push 测量覆盖率
- 回归告警: 低于 85% 时 PR block
- 新增 package 时使用同样 ≥85% 标准

## Refs

- ROADMAP.md "v1.0.0 毕业条件" #7 (P-B.11.3)
- docs/STABILITY.md graduation criterion #6
