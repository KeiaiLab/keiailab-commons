# ROADMAP — operator-commons

> [English](ROADMAP.md) | [한국어](ROADMAP.ko.md) | **日本語** | [中文](ROADMAP.zh.md)

> ⚠️ This translation is AI-generated and pending native review.

本 ROADMAP はライブラリの進化を 3 つの軸 — *API 安定性 tier*、*v1.0.0
昇格基準*、*パッケージ別 follow-up 項目* — で追跡します。本プロジェクトは
時間ベースの締切を維持しません。ライブラリは downstream consumer の
ニーズに応じて進化します。

## チェックボックスの意味

| マーカー | 意味 |
|---|---|
| `[x]` | コード + テスト両方存在。downstream import 動作。 |
| `[~]` | 部分実装 (helper 存在、検証は未完)。 |
| `[ ]` | 未着手。 |

## API 安定性 tier (現在の v0.9.x 候補)

| パッケージ | Tier | 昇格基準 |
|---|---|---|
| `pkg/finalizer` | **Stable** | v1 entry (追加作業なし)。 |
| `pkg/labels` | **Stable** | v1 entry (追加作業なし)。 |
| `pkg/status` | **Stable** | v1 entry (追加作業なし)。 |
| `pkg/storageclass` | **Stable** | 自明な検証 surface (regex + nil チェック)。 |
| `pkg/version` (`Matrix` 含む) | Beta | Generic `Matrix[E]` の cross-repo 検証。 |
| `pkg/monitoring` | Beta | `ServiceMonitor` の downstream 横断同値性 e2e。 |
| `pkg/networkpolicy` | Beta | 4 方向 (ingress / egress × TCP / UDP) 検証。 |
| `pkg/security` | Beta | Restricted PSA ガードの downstream 横断。 |
| `pkg/events` | Beta | Downstream live 採用 + reconciliation regression 0。 |
| `pkg/pvc` | Beta | Downstream PVC expansion live 採用。 |
| `pkg/topology` | Beta | Downstream topology spread live 採用。 |
| `pkg/webhook` | **Experimental** | Multi-downstream 採用 + 安定化。 |
| `pkg/probes` | **Experimental** | 2+ downstream 採用 → Beta。 |

**Tier セマンティクス**:

- **Stable** — patch / minor release で BREAKING CHANGE なし。
  deprecation を使用: 印を付けて 2 minor 維持後に削除。
- **Beta** — minor release で BREAKING CHANGE 許可 (CHANGELOG に記載
  必須)。API 形状はほぼ固まっている。
- **Experimental** — 任意の release で BREAKING CHANGE 可能。呼び出し側が
  リスクを負う。

## v1.0.0 昇格基準 (チェックリスト)

- [ ] 全パッケージが **Stable** tier に到達。
- [ ] 6 つ以上の連続 minor release で BREAKING CHANGE ゼロ。
- [ ] godoc カバレッジ ≥ 80 % (`go doc -all ./...` 基準)。
- [ ] CHANGELOG.md の v0.x 進化履歴 + v1.0.0 release notes。
- [ ] CITATION.cff + DOI (academic citation)。
- [ ] Downstream consumer が v1.0.0 commons を regression 0 で import。
- [x] `go vet ./... && go test ./...` clean (カバレッジ 96.3 % > 85 %
  閾値)。
- [x] API 安定性の約束文書 — [STABILITY.md](STABILITY.ja.md)。
- 検証: downstream consumer CI が `operator-commons v1.0.0` に対する
  全 e2e テストを pass。

## パッケージ別 follow-up

### pkg/finalizer (Stable)

- [x] `Add`、`Remove`、`Contains` helper — `pkg/finalizer/finalizer.go`。
- [x] controller-runtime を回避 (stdlib `slices` のみ)。
- [x] Unit テスト — `pkg/finalizer/finalizer_test.go`。
- [x] 複数 finalizer 順序 helper — `pkg/finalizer/order.go`
  `EnsureOrder`。
- 検証: downstream finalizer regression 0。

### pkg/labels (Stable)

- [x] Recommended Kubernetes labels (`app.kubernetes.io/*`) —
  `pkg/labels/labels.go`。
- [x] Component / instance / part-of マッピング。
- [x] Unit テスト — `pkg/labels/labels_test.go`。
- [x] Recommended labels v2 (K8s 1.30+) — `pkg/labels/v2.go` `AllV2`。
- 検証: downstream `metadata.labels` 一貫性。

### pkg/status (Stable)

- [x] Condition カタログ helper — `pkg/status/conditions.go`。
- [x] `SetAvailable` sugar。
- [x] Unit テスト。
- [x] Condition reason カタログ文書 — `pkg/status/REASONS.md`。
- 検証: downstream `kubectl get <kind> -o yaml`
  `.status.conditions` parity。

### pkg/version (Beta)

- [x] `Matrix[E]` generic — `pkg/version/matrix.go`。
- [x] `SetAvailable` sugar。
- [x] semver-aware version 比較 — `pkg/version/version.go`。
- [x] Cross-version 互換性テスト —
  `pkg/version/api_stability_test.go`。
- [x] `Matrix[E]` serializer (JSON / YAML) —
  `pkg/version/serializer.go`。
- [ ] **Tier 昇格** → Stable。
- 検証: downstream version validation parity。

### pkg/monitoring (Beta)

- [x] Prometheus ServiceMonitor builder — `pkg/monitoring/monitoring.go`。
- [x] Unit テスト。
- [x] PrometheusRule builder (alert / recording 共有) —
  `pkg/monitoring/rule.go`。
- [x] OpenTelemetry exporter helper — `pkg/monitoring/otel.go`。
- [ ] Downstream 同値性 e2e — 同入力 → 同 manifest 出力。
- [ ] **Tier 昇格** → Stable。
- 検証: `monitoring_test.go` の golden file diff = 0。

### pkg/networkpolicy (Beta)

- [x] NetworkPolicy builder — `pkg/networkpolicy/networkpolicy.go`。
- [x] Default-deny + explicit rule helper。
- [x] Unit テスト。
- [x] 4 方向検証 — `pkg/networkpolicy/four_dir_test.go`。
- [x] CIDR + namespace + pod selector combo —
  `pkg/networkpolicy/combo.go`。
- [ ] **Tier 昇格** → Stable。
- 検証: kind cluster で NetworkPolicy を apply、観測された deny / allow
  パスが期待と一致。

### pkg/security (Beta)

- [x] SecurityContext helper (restricted PSA-compliant) —
  `pkg/security/security.go`。
- [x] RBAC helper。
- [x] Unit テスト。
- [x] Restricted PSA regression ガード —
  `pkg/security/psa_guard_test.go`。
- [x] Pod / Container SecurityContext 分割 —
  `pkg/security/split.go`。
- [x] seccompProfile デフォルト helper — `pkg/security/seccomp.go`。
- [ ] **Tier 昇格** → Stable。
- 検証: `kubectl label ns <ns>
  pod-security.kubernetes.io/enforce=restricted` 後、downstream pod が
  Ready に到達。

### pkg/webhook (Experimental)

- [x] Webhook ユーティリティ base — `pkg/webhook/webhook.go`。
- [x] Unit テスト。
- [x] Conversion webhook helper — `pkg/webhook/conversion.go`。
- [x] Validation webhook パターン —
  `pkg/webhook/validation_patterns.go`。
- [ ] Multi-downstream live 採用 → 安定化。
- [ ] **Tier 昇格** → Beta → Stable。
- 検証: 2 つ以上の downstream consumer が同 helper を regression 0 で
  使用。

### pkg/storageclass (Stable)

- [x] DNS-1123 subdomain 検証 — `pkg/storageclass/validator.go`。
- [x] Normalize / MustNormalize — empty → nil (cluster default) + trim
  + pointer return。
- [x] 12 unit テスト (100 % カバレッジ) —
  `pkg/storageclass/validator_test.go`。
- [ ] Downstream live 採用 + regression 0。

### pkg/events (Beta)

- [x] Recorder interface — `client-go` `record.EventRecorder` を
  import せず互換。
- [x] 9 個の Reason 定数。
- [x] Emit / Emitf / EmitWarning / EmitWarningf / WrappedError — 全
  nil-safe。
- [x] Unit テスト (100 % カバレッジ) — `pkg/events/events_test.go`。
- [ ] Downstream live 採用 — Event reason を reconcile path 全体で統一。
- [ ] **Tier 昇格** → Stable。
- 検証: downstream Reconcile path が commons の reason 定数を
  regression 0 で使用。

### pkg/probes (Experimental)

- [x] Fluent builder — HTTP / HTTPS / TCP / Exec handler。
- [x] kubelet デフォルト (Period = 10 s / Timeout = 1 s /
  SuccessThreshold = 1 / FailureThreshold = 3)。
- [x] InitialDelay / Period / Timeout 負値を 0 に clamp。
- [x] handler 未設定時 `Build()` panic (fail-fast 契約)。
- [x] Unit テスト (100 % カバレッジ) — `pkg/probes/builder_test.go`。
- [ ] 2+ downstream live 採用 (Beta 基準)。
- [ ] **Tier 昇格** → Beta → Stable。

### pkg/pvc (Beta)

- [x] PVC expansion helper — `pkg/pvc/expansion.go`。
- [x] Unit テスト — `pkg/pvc/expansion_test.go`。
- [ ] Downstream live 採用 + PVC resize regression 0。
- [ ] **Tier 昇格** → Stable。

### pkg/topology (Beta)

- [x] PVC topology spread helper — `pkg/topology/spread.go`。
- [x] Unit テスト — `pkg/topology/spread_test.go`。
- [ ] Downstream live 採用 + spread constraint 検証。
- [ ] **Tier 昇格** → Stable。

## 依存性ポリシー

- **Kubernetes API のみ** — `k8s.io/api`、`k8s.io/apimachinery`、
  `k8s.io/utils`。controller-runtime 依存性は leaf package に *追加禁止*。
- **Apache-2.0 互換ライセンスのみ** — 依存性追加には ADR 必須。
- **完全な godoc** — すべての新規 public API に godoc 必須。

## ガバナンス / 追跡

- **CHANGELOG.md** — `git-cliff` による自動生成。厳格な semantic
  versioning。
- **CITATION.cff** — 学術引用。v1.0.0 で DOI 発行。
- **ADR** — `docs/kb/adr/` がすべての設計決定を追跡。
- **AGENTS.md** — AI コラボレーションのランブック。

## Non-Goals (意図的にスコープ外)

- ❌ **controller-runtime 依存性** — leaf package は controller-runtime
  free でなければならない。
- ❌ **downstream 固有ロジック** — operator 固有のコードは呼び出し元
  リポジトリに置く。ライブラリは *共有* helper のみ提供。
- ❌ **時間ベースの roadmap** — 機能チェックリストと完了率を使用。
- ❌ **GitHub Actions release ゲート** — ローカル 4 階層に委任。
- ❌ **Plugin / extension SDK ポジション** — これはライブラリであり
  フレームワークではない。
- ❌ **早すぎる v1.0.0** — 昇格基準を満たすまで v0.x に留まる。

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
