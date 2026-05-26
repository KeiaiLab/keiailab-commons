# カバレッジレポート

> [English](coverage-report.md) | [한국어](coverage-report.ko.md) | **日本語** | [中文](coverage-report.zh.md)

> `go test ./pkg/... -coverprofile=cover.out` 測定結果。
> v1.0.0 卒業条件: カバレッジ ≥ 85 %。

## 測定結果 (2026-05-26)

| パッケージ | カバレッジ | 状態 |
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
| **合計 (statements)** | **96.2%** | **✅ PASS (≥85%)** |

全13パッケージがv1.0.0閾値（85%）を超過しています。

## 再現コマンド

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.2%
```

## 保守

- 品質ゲート: `pre-push` hookが毎pushでカバレッジを測定します。
- 回帰警報: パッケージが85%未満に低下するとPRをブロックします。
- 新規パッケージも同じ ≥ 85% 基準を満たす必要があります。

## 参照

- [ROADMAP.md](ROADMAP.md) — v1.0.0 卒業条件 #7。
- [STABILITY.md](STABILITY.md) — 卒業条件 #6。
