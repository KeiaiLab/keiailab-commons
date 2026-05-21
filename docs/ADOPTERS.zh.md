# Adopters of operator-commons

> ⚠️ This translation is AI-generated and pending native review. — 本翻译为 Claude 机器翻译结果。

本库不是*外部用户直接 import*, 而是 keiailab 的 3 个 operator 作为公共依赖导入的*内部共享库*. 外部用户通过 *consumer operator* 间接使用本库的代码.

## Direct consumers (in-org)

| Operator | 使用包 | 起始版本 | 当前版本 | 最新 commit | 更新日期 |
|---|---|---|---|---|---|
| `keiailab/mongodb-operator` | labels, security, webhook, version, finalizer, networkpolicy | v0.1.0 | **v0.7.0** | `97140db` | 2026-05-20 |
| `keiailab/postgres-operator` | labels, security, webhook, status, version | v0.1.0 | **v0.7.0** | `8c9db39` | 2026-05-20 |
| `keiailab/valkey-operator` | labels, security, webhook, monitoring, finalizer, networkpolicy | v0.1.0 | **v0.6.0** ⚠️ (1 minor lag, I09 upgrade 计划) | `e878420` | 2026-05-20 |

## v0.8.0 candidate — 新 3 包导入预定矩阵

| Operator | `pkg/probes` (Experimental) | `pkg/storageclass` (Stable) | `pkg/events` (Beta) | 导入 PR target |
|---|---|---|---|---|
| `keiailab/postgres-operator` | builders.go L986-998 (2 HTTP probe sites) | builders.go `storageClassPtr()` | RFC-0023 Phase 2 sister (commit `1494ff6` 后续) | 别 PR |
| `keiailab/mongodb-operator` | resources/builder.go L613-626 (2 Exec probe sites, mongosh) | builder.go sister 模式 | candidate (RFC-0023 Phase 2 后续) | 别 PR |
| `keiailab/valkey-operator` | resources/statefulset.go L126-139 (2 Exec probe sites, valkey-cli) | statefulset.go sister 模式 | candidate (RFC-0023 Phase 2 后续) | 别 PR |

> **AST audit evidence (2026-05-21)**: probes builder 9 sites (postgres 2 + mongo 2 + valkey 2 + 3 cross-cutting) / 50-55 LOC reduction estimate. storageclass / events 是 trivial helper.

> **实时 evidence (2026-05-20)**: 本表基于每个 operator 的 `go.mod` 实时 `require github.com/keiailab/operator-commons <ver>` + `grep -rn "github.com/keiailab/operator-commons" --include="*.go"` import 结果.

## External adopters

本库以 *Go module* 公开,任何人都可以使用 `go get github.com/keiailab/operator-commons`. 但在 v0.x 阶段公开 API breaking 可能自由发生, 因此*推荐外部用户使用 v1.0 stable 之后*.

如希望登记外部使用案例,请通过 PR 添加 row:

```markdown
| **<组织 / 项目>** ([profile](<URL>)) | <使用包> | <使用起始版本> | <当前版本> | <登记日期 YYYY-MM-DD> |
```

## CNCF / 许可证

- 许可证: Apache-2.0 only (AGPL/BUSL transitive 依赖 0 件目标)
- 本 ADOPTERS 也作为 CNCF graduation criteria 的 "≥1 public adopter" 同等的 *cross-repo dependency declaration* 使用.

---

<p align="center">
  <b>keiailab 操作器家族</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
