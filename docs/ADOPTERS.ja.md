# Adopters of operator-commons

> ⚠️ This translation is AI-generated and pending native review. — 本翻訳は Claude による機械翻訳結果です。

本ライブラリは*外部ユーザーが直接 import* するよりも、keiailab の 3 operator が共通依存性として import する*内部共有ライブラリ*です。外部ユーザーは *consumer operator* を通じて間接的に本ライブラリのコードを使用します。

## Direct consumers (in-org)

| Operator | 使用パッケージ | 開始バージョン | 現在のバージョン | 最新 commit | 更新日 |
|---|---|---|---|---|---|
| `keiailab/mongodb-operator` | labels, security, webhook, version, finalizer, networkpolicy | v0.1.0 | **v0.7.0** | `97140db` | 2026-05-20 |
| `keiailab/postgres-operator` | labels, security, webhook, status, version | v0.1.0 | **v0.7.0** | `8c9db39` | 2026-05-20 |
| `keiailab/valkey-operator` | labels, security, webhook, monitoring, finalizer, networkpolicy | v0.1.0 | **v0.6.0** ⚠️ (1 minor lag, I09 upgrade 予定) | `e878420` | 2026-05-20 |

## v0.8.0 candidate — 新 3 パッケージ導入予定マトリクス

| Operator | `pkg/probes` (Experimental) | `pkg/storageclass` (Stable) | `pkg/events` (Beta) | 導入 PR target |
|---|---|---|---|---|
| `keiailab/postgres-operator` | builders.go L986-998 (2 HTTP probe sites) | builders.go `storageClassPtr()` | RFC-0023 Phase 2 sister (commit `1494ff6` 後続) | 別 PR |
| `keiailab/mongodb-operator` | resources/builder.go L613-626 (2 Exec probe sites, mongosh) | builder.go sister パターン | candidate (RFC-0023 Phase 2 後続) | 別 PR |
| `keiailab/valkey-operator` | resources/statefulset.go L126-139 (2 Exec probe sites, valkey-cli) | statefulset.go sister パターン | candidate (RFC-0023 Phase 2 後続) | 別 PR |

> **AST audit evidence (2026-05-21)**: probes builder 9 sites (postgres 2 + mongo 2 + valkey 2 + 3 cross-cutting) / 50-55 LOC reduction estimate. storageclass / events は trivial helper.

> **ライブ evidence (2026-05-20)**: 本表は各 operator の `go.mod` ライブ `require github.com/keiailab/operator-commons <ver>` + `grep -rn "github.com/keiailab/operator-commons" --include="*.go"` import 結果に基づく.

## External adopters

本ライブラリは *Go module* として公開されており、誰でも `go get github.com/keiailab/operator-commons` で使用可能です。ただし v0.x 段階では公開 API breaking が自由に発生するため、*外部ユーザーには v1.0 stable 以降を推奨*します。

外部使用事例登録を希望される場合は PR で row 追加:

```markdown
| **<組織 / プロジェクト>** ([profile](<URL>)) | <使用パッケージ> | <使用開始バージョン> | <現在のバージョン> | <登録日 YYYY-MM-DD> |
```

## CNCF / ライセンス

- ライセンス: Apache-2.0 only (AGPL/BUSL transitive 依存性 0 件目標)
- 本 ADOPTERS は CNCF graduation criteria の "≥1 public adopter" と同等の *cross-repo dependency declaration* としても活用されます.

---

<p align="center">
  <b>keiailab オペレーターファミリー</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
