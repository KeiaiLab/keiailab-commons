# API 稳定性承诺

> ⚠️ This translation is AI-generated and pending native review. — 本翻译为 Claude 机器翻译结果。母语审阅者校验前为 `[待校验]` 状态。

> operator-commons 的 API 稳定性等级和 breaking-change 策略。

## 等级 (Tier)

operator-commons 使用 3-tier 稳定性:

| Tier | 兼容性 | Breaking change 策略 |
|---|---|---|
| **Stable** | minor release 间向后兼容. 仅 patch fix. | 需要 major version bump (`v2.0.0`) + 6+ 个月 deprecation notice |
| **Beta** | patch release 间兼容. minor release 可包含至少 1-release deprecation window 的 source-compatible API 改进 | 允许 minor version bump 配合 `CHANGELOG.md` 中的 deprecation 条目 |
| **Experimental** | 无兼容性承诺. 任何 release 都可 break. | 任何 release — 必须在 `CHANGELOG.md` "BREAKING" 段标记 |

## 当前 tier matrix

基于 `ROADMAP.md` "API Stability Tier" 表:

| 包 | Tier | 升级标准 |
|---|---|---|
| `pkg/finalizer` | Stable | (v1 entry — 无需追加工作) |
| `pkg/labels` | Stable | (v1 entry) |
| `pkg/status` | Stable | (v1 entry) |
| `pkg/version` | Beta | Generic `Matrix[E]` 3-repo verify |
| `pkg/monitoring` | Beta | ServiceMonitor 3-repo e2e |
| `pkg/networkpolicy` | Beta | 4-direction TCP/UDP verify |
| `pkg/security` | Beta | restricted PSA 3-repo guard |
| `pkg/webhook` | Experimental | Multi-repo adoption + stabilize |

## 升级流程

1. Sub-task PR 与升级提案一同 open (例如 `feat(pkg/X): promote to Stable`)
2. 升级标准 (ROADMAP 标准) 通过 CI 验证:
   - Cross-repo import 通过 (3 个 operator)
   - 包 godoc coverage ≥80%
   - Unit + integration test coverage ≥85%
   - exported API 中无 `// TODO` / `// FIXME`
3. ROADMAP.md tier 表在同一 PR 更新
4. `CHANGELOG.md` 添加 "Changed" 条目

## Breaking-change 策略

**breaking change** = 以下任一:
- 移除 exported identifier (function / type / constant / variable)
- 变更 exported signature (parameter, return type)
- 移除包
- 需要 caller code 修改的行为变更

### 各 tier:

- **Stable**: 在 v2.0.0 之前禁止. 使用 deprecation: 添加 `// Deprecated: ...` godoc + 新替代,保留旧 6+ 个月
- **Beta**: 允许 1-release deprecation. 必须出现在 `CHANGELOG.md` "Deprecated" → "Removed" pipeline 中
- **Experimental**: 任何 release 都允许,必须出现在 `CHANGELOG.md` "BREAKING" 段

## Semantic versioning

`vMAJOR.MINOR.PATCH` per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html):
- **MAJOR**: Stable tier 的 breaking change — 需要 v2.0.0 graduation review
- **MINOR**: 新 feature, Beta tier 增加, non-breaking Stable 改进
- **PATCH**: 仅 bug fix, 无 API surface 变更

## v1.0.0 graduation

需要*全部*:
1. 8/8 包到达 Stable tier
2. 在 6+ 连续 minor release (v0.8 → v0.13) 中 BREAKING CHANGE 为 0
3. godoc coverage ≥80% (本文档 + per-package — 用 `scripts/godoc-coverage.sh` 验证)
4. CITATION.cff + Zenodo DOI
5. `v1.0.0-rc.N` 的 3-repo import e2e 验证
6. `go vet ./... && go test ./...` clean + coverage ≥85%
7. CHANGELOG.md v0.x evolution history + v1.0.0 release notes
8. 本 `docs/STABILITY.md` (当前文件)
9. `pkg/finalizer` multi-finalizer order 保证
10. `pkg/labels` K8s 1.30+ v2 mapping
11. `pkg/status` Condition reason catalog 文档

跟踪: `~/.claude/plans/2026-05-14-4-operators-100pct/P-B.md` (29 sub-task).

## Caller 责任

Caller (mongodb-operator, valkey-operator, postgres-operator):
- 在 v1.0.0 之前在 `go.mod` 中以 `vMAJOR.MINOR.PATCH` pin
- 订阅 `CHANGELOG.md` 以获取 deprecation warning
- 在 GA 之前对 `v1.0.0-rc.N` 进行测试

## 参考

- `ROADMAP.md` — tier 表 + graduation 标准
- `CHANGELOG.md` — version history
- `CITATION.cff` — academic citation
- `ADOPTERS.md` — 3-repo adoption matrix
- [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html)
