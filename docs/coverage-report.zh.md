# 覆盖率报告

> [English](coverage-report.md) | [한국어](coverage-report.ko.md) | [日本語](coverage-report.ja.md) | **中文**

> `go test ./pkg/... -coverprofile=cover.out` 测量结果。
> v1.0.0 毕业条件: 覆盖率 ≥ 85 %。

## 测量结果 (2026-05-26)

| 包 | 覆盖率 | 状态 |
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
| **合计 (statements)** | **96.2%** | **✅ PASS (≥85%)** |

全部13个包均超过v1.0.0阈值（85%）。

## 重现命令

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.2%
```

## 维护

- 质量门: `pre-push` hook在每次push时测量覆盖率。
- 回归警报: 任何包低于85%时阻止PR。
- 新包必须满足相同的 ≥ 85% 标准。

## 参考

- [ROADMAP.md](ROADMAP.md) — v1.0.0 毕业条件 #7。
- [STABILITY.md](STABILITY.md) — 毕业条件 #6。
