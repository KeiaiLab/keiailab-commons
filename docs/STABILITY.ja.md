# API 安定性の約束

> ⚠️ This translation is AI-generated and pending native review. — 本翻訳は Claude による機械翻訳結果です。母語話者検証完了まで `[検証必要]` 状態です。

> operator-commons の API 安定性ティアと breaking-change ポリシー。

## ティア (Tier)

operator-commons は 3-tier 安定性を使用します:

| Tier | 互換性 | Breaking change ポリシー |
|---|---|---|
| **Stable** | minor release 間 backwards-compatible. patch fix のみ. | major version bump (`v2.0.0`) + 6 ヶ月以上の deprecation notice 必要 |
| **Beta** | patch release 間互換. minor release は最低 1-release deprecation window 付きの source-compatible API 改善を含む可 | `CHANGELOG.md` の deprecation 項目と共に minor version bump 許可 |
| **Experimental** | 互換性の約束無し. どの release でも breaking 可. | どの release でも — `CHANGELOG.md` "BREAKING" セクションに flag 必須 |

## 現在の tier matrix

`ROADMAP.md` "API Stability Tier" 表に基づく:

| パッケージ | Tier | 昇格基準 |
|---|---|---|
| `pkg/finalizer` | Stable | (v1 entry — 追加作業無し) |
| `pkg/labels` | Stable | (v1 entry) |
| `pkg/status` | Stable | (v1 entry) |
| `pkg/version` | Beta | Generic `Matrix[E]` 3-repo verify |
| `pkg/monitoring` | Beta | ServiceMonitor 3-repo e2e |
| `pkg/networkpolicy` | Beta | 4-direction TCP/UDP verify |
| `pkg/security` | Beta | restricted PSA 3-repo guard |
| `pkg/webhook` | Experimental | Multi-repo adoption + stabilize |

## 昇格プロセス

1. Sub-task PR が昇格提案と共に open (例: `feat(pkg/X): promote to Stable`)
2. 昇格基準 (ROADMAP 基準) を CI で検証:
   - Cross-repo import passes (3 operator)
   - パッケージ godoc coverage ≥80%
   - Unit + integration test coverage ≥85%
   - exported API に `// TODO` / `// FIXME` 無し
3. ROADMAP.md tier 表を同 PR で更新
4. `CHANGELOG.md` "Changed" エントリー

## Breaking-change ポリシー

**breaking change** = 以下のいずれか:
- exported identifier 削除 (function / type / constant / variable)
- exported signature 変更 (parameter, return type)
- パッケージ削除
- caller code 修正を要する動作変更

### 各 tier 別:

- **Stable**: v2.0.0 まで禁止. deprecation 使用: `// Deprecated: ...` godoc + 新代替追加、旧は 6 ヶ月以上保持
- **Beta**: 1-release deprecation 付きで許可. `CHANGELOG.md` "Deprecated" → "Removed" pipeline 出現必須
- **Experimental**: どの release でも許可、`CHANGELOG.md` "BREAKING" セクション出現必須

## Semantic versioning

`vMAJOR.MINOR.PATCH` per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html):
- **MAJOR**: Stable tier の breaking change — v2.0.0 graduation review 必要
- **MINOR**: 新 feature, Beta tier 追加, non-breaking Stable 改善
- **PATCH**: bug fix のみ、API surface 変更無し

## v1.0.0 graduation

以下*すべて*必要:
1. 8/8 パッケージが Stable tier 到達
2. 6+ 連続 minor release (v0.8 → v0.13) で BREAKING CHANGE 0
3. godoc coverage ≥80% (本文書 + per-package — `scripts/godoc-coverage.sh` で検証)
4. CITATION.cff + Zenodo DOI
5. `v1.0.0-rc.N` の 3-repo import e2e 検証
6. `go vet ./... && go test ./...` clean + coverage ≥85%
7. CHANGELOG.md v0.x evolution history + v1.0.0 release notes
8. 本 `docs/STABILITY.md` (現在のファイル)
9. `pkg/finalizer` multi-finalizer order 保証
10. `pkg/labels` K8s 1.30+ v2 mapping
11. `pkg/status` Condition reason catalog 文書

追跡: `~/.claude/plans/2026-05-14-4-operators-100pct/P-B.md` (29 sub-task).

## Caller の責任

Caller (mongodb-operator, valkey-operator, postgres-operator):
- v1.0.0 まで `go.mod` で `vMAJOR.MINOR.PATCH` で pin
- deprecation warning のため `CHANGELOG.md` 購読
- GA 前 `v1.0.0-rc.N` に対してテスト

## 参考

- `ROADMAP.md` — tier 表 + graduation 基準
- `CHANGELOG.md` — version history
- `CITATION.cff` — academic citation
- `ADOPTERS.md` — 3-repo adoption matrix
- [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html)
