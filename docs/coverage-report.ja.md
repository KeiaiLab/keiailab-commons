# Coverage Report — v1.0.0 graduation B.11.3

> ⚠️ This translation is AI-generated and pending native review. — 本翻訳は Claude による機械翻訳結果です。

> `go test ./pkg/... -coverprofile=cover.out` 測定結果。v1.0.0 卒業条件 #7 (`カバレッジ ≥ 85%`) 充足。

## 測定結果 (2026-05-14)

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
| **合計 (statements)** | **96.3%** | **✅ PASS (≥85%)** |

8 package すべてが v1.0.0 threshold (85%) を通過。

## 再現コマンド

```bash
go test -coverprofile=cover.out -covermode=atomic ./pkg/...
go tool cover -func=cover.out | tail -1
# total:                                  (statements)            96.3%
```

## 今後の維持

- CI ゲート: `pre-push` hook で push の度にカバレッジ測定
- Regression alarm: 85% 以下に落ちると PR ブロック
- 新規 package 追加時に同じ ≥85% 基準

## Refs

- ROADMAP.md "v1.0.0 卒業条件" #7
- docs/STABILITY.md graduation criterion #6
