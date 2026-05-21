# 日本語 用語集 (Glossary)

> [English](../../README.md) | [한국어](glossary-ko.md) | **日本語** | [中文](glossary-zh.md)
>
> 本用語集は keiailab operator family 4 リポジトリ (operator-commons + postgres-operator + mongodb-operator + valkey-operator) の日本語翻訳時に*必ず参照*する標準用語集です。
>
> ⚠️ **This translation is AI-generated and pending native review.** — Claude による機械翻訳。母語話者 (native reviewer) による品質検証後 `[x]` 完了状態へ昇格します。本ファイルの全項目は `[検証必要]` 状態として扱ってください。

## §1 一貫性ルール

1. **コード識別子は英文そのまま** — 例: `ValkeyCluster`, `kubectl`, `Helm`, `pkg/probes`. 日本語翻訳禁止。
2. **標準 K8s 用語は英文優先 + 日本語併記** — 例: `Pod (ポッド)`, `Deployment (デプロイメント)`. 本文初出時に英文 + 括弧日本語、以降は日本語単独可。
3. **operator-commons API 名称は英文そのまま** — 例: `Reconciler`, `Finalizer`, `EventRecorder`. 日本語併記可 (`Finalizer (ファイナライザー)`)。
4. **外部ユーザー可視文書 (README/CONTRIBUTING/SECURITY 等)** = 敬体 (`です/ます調`). 内部文書 (HANDOFF/AGENTS 等) = 常体または自由。
5. **敬体と常体の混在禁止** — 1 文書内で一貫性を保つこと。

## §2 Kubernetes 標準用語

| English (canonical) | 日本語 |
|---|---|
| CustomResourceDefinition (CRD) | カスタムリソース定義 (CRD) |
| Custom Resource (CR) | カスタムリソース (CR) |
| Reconciler | リコンサイラー (または英文そのまま) |
| Reconcile Loop | 調整ループ (reconcile ループ) |
| Controller | コントローラー |
| Operator | オペレーター (または英文そのまま) |
| Operator Pattern | オペレーターパターン |
| Finalizer | ファイナライザー (または英文そのまま) |
| Webhook | ウェブフック (または英文そのまま) |
| Validating Admission Webhook | 検証用アドミッションウェブフック |
| Mutating Admission Webhook | 変更用アドミッションウェブフック |
| Conversion Webhook | 変換ウェブフック |
| Pod | ポッド (Pod) |
| Deployment | デプロイメント (Deployment) |
| StatefulSet | ステートフルセット (StatefulSet) |
| DaemonSet | デーモンセット (DaemonSet) |
| Job | ジョブ (Job) |
| CronJob | クロンジョブ (CronJob) |
| Service | サービス (Service) |
| Ingress | イングレス (Ingress) |
| ConfigMap | コンフィグマップ (ConfigMap) |
| Secret | シークレット (Secret) |
| PersistentVolume (PV) | 永続ボリューム (PV) |
| PersistentVolumeClaim (PVC) | 永続ボリュームクレーム (PVC) |
| StorageClass | ストレージクラス (StorageClass) |
| Namespace | ネームスペース |
| Node | ノード |
| Cluster | クラスター |
| Control Plane | コントロールプレーン |
| Worker Node | ワーカーノード |
| API Server | API サーバー |
| Scheduler | スケジューラー |
| kubelet | kubelet (英文そのまま) |
| kube-proxy | kube-proxy (英文そのまま) |
| etcd | etcd (英文そのまま) |
| Helm Chart | ヘルムチャート (Helm Chart) |
| Helm Release | ヘルムリリース |
| Kubernetes RBAC | Kubernetes RBAC |
| ServiceAccount | サービスアカウント (ServiceAccount) |
| Role / ClusterRole | ロール / クラスターロール |
| RoleBinding / ClusterRoleBinding | ロールバインディング / クラスターロールバインディング |
| NetworkPolicy | ネットワークポリシー (NetworkPolicy) |
| PodSecurityAdmission (PSA) | ポッドセキュリティアドミッション (PSA) |
| restricted profile | restricted プロファイル (英文そのまま) |
| baseline profile | baseline プロファイル (英文そのまま) |
| privileged profile | privileged プロファイル (英文そのまま) |
| Probe (liveness/readiness/startup) | プローブ (liveness/readiness/startup) |
| HTTPGetAction | HTTP Get アクション (HTTPGetAction) |
| TCPSocketAction | TCP ソケットアクション (TCPSocketAction) |
| ExecAction | 実行アクション (ExecAction) |
| InitContainer | 初期化コンテナ (InitContainer) |
| Sidecar Container | サイドカーコンテナ |

## §3 operator-commons ライブラリ用語

| English (canonical) | 日本語 |
|---|---|
| `pkg/finalizer` | ファイナライザーパッケージ (`pkg/finalizer`) |
| `pkg/labels` | ラベルパッケージ (`pkg/labels`) |
| `pkg/status` | 状態条件パッケージ (`pkg/status`) |
| `pkg/version` | バージョン互換性パッケージ (`pkg/version`) |
| `pkg/monitoring` | モニタリングパッケージ (`pkg/monitoring`) |
| `pkg/networkpolicy` | ネットワークポリシーパッケージ (`pkg/networkpolicy`) |
| `pkg/security` | セキュリティパッケージ (`pkg/security`) |
| `pkg/webhook` | ウェブフックパッケージ (`pkg/webhook`) |
| `pkg/probes` (v0.8.0 新規) | プローブビルダーパッケージ (`pkg/probes`) |
| `pkg/storageclass` (v0.8.0 新規) | ストレージクラスパッケージ (`pkg/storageclass`) |
| `pkg/events` (v0.8.0 新規) | イベントレコーダーパッケージ (`pkg/events`) |
| Recorder interface | Recorder インターフェース (英文そのまま) |
| EventType (Normal / Warning) | イベントタイプ (Normal / Warning) |
| Reason (event reason) | イベント理由 (Reason) |
| Builder pattern | ビルダーパターン |
| Fluent API | フルーエント API (Fluent API) |
| API Stability Tier | API 安定性ティア |
| Stable / Beta / Experimental | 安定 / ベータ / 実験的 (英文推奨 — Stable/Beta/Experimental) |
| Tier 昇格 (promotion) | ティア昇格 (Tier promotion) |
| Breaking change | 破壊的変更 (Breaking change) |
| Semver (Semantic Versioning) | セマンティックバージョニング (Semver) |
| Deprecated | 非推奨 (Deprecated) |

## §4 reconciler パターン用語

| English (canonical) | 日本語 |
|---|---|
| Desired state | 望ましい状態 |
| Actual state | 実際の状態 |
| Reconcile | 調整 (reconcile) |
| Drift | ドリフト (英文そのまま) |
| Idempotency | 冪等性 |
| Eventual consistency | 結果整合性 |
| Owner reference | オーナー参照 (OwnerReference) |
| Garbage collection | ガベージコレクション (GC) |
| Status condition | 状態条件 (status condition) |
| Available / Ready / Progressing / Degraded | Available / Ready / Progressing / Degraded (英文そのまま) |
| Failover | フェイルオーバー (failover) |
| Provisioning | プロビジョニング |
| Rolling update | ローリングアップデート |
| Blue-green deployment | ブルーグリーンデプロイメント |
| Canary deployment | カナリアデプロイメント |
| Backup / Restore | バックアップ / リストア |
| Point-in-time recovery (PITR) | ポイントインタイムリカバリ (PITR) |
| Horizontal scaling | 水平スケーリング |
| Vertical scaling | 垂直スケーリング |
| Sharding | シャーディング |
| Replication | レプリケーション (replication) |
| Primary / Replica / Secondary | Primary / Replica / Secondary (英文そのまま) |

## §5 セキュリティ + 認証用語

| English (canonical) | 日本語 |
|---|---|
| Authentication (AuthN) | 認証 (Authentication) |
| Authorization (AuthZ) | 認可 (Authorization) |
| OIDC | OIDC (英文そのまま) |
| LDAP | LDAP (英文そのまま) |
| TLS / mTLS | TLS / mTLS (英文そのまま) |
| Certificate | 証明書 (certificate) |
| Service mesh | サービスメッシュ |
| Tenant / Tenancy | テナント / テナンシー |
| Multi-tenancy | マルチテナンシー |
| Encryption at rest | 保存時暗号化 |
| Encryption in transit | 通信時暗号化 |

## §6 運用 + 観察用語

| English (canonical) | 日本語 |
|---|---|
| Observability | 可観測性 (observability) |
| Metric / Monitoring | メトリック / モニタリング |
| Logging | ロギング |
| Tracing | トレーシング |
| Prometheus | Prometheus (英文そのまま) |
| Grafana | Grafana (英文そのまま) |
| Alert / Alerting | アラート / アラート発報 |
| Service Level Objective (SLO) | サービスレベル目標 (SLO) |
| Service Level Indicator (SLI) | サービスレベル指標 (SLI) |
| Postmortem | ポストモーテム (postmortem) |
| Incident | インシデント (incident) |
| Severity (SEV-1 / SEV-2 / SEV-3) | 重大度 (SEV-1 / SEV-2 / SEV-3、英文そのまま) |

## §7 ガバナンス + 協業用語

| English (canonical) | 日本語 |
|---|---|
| ADR (Architecture Decision Record) | アーキテクチャ決定記録 (ADR) |
| RFC (Request For Comments) | RFC (英文そのまま) |
| Roadmap | ロードマップ |
| Adopter | アダプター / 採用者 (Adopter) |
| Maintainer | メンテナー |
| Contributor | コントリビューター |
| Pull Request (PR) | プルリクエスト (PR) |
| Merge Request (MR) | マージリクエスト (MR、GitLab) |
| Issue | イシュー |
| Code review | コードレビュー |
| Squash merge | スカッシュマージ |
| Rebase | リベース |
| Cherry-pick | チェリーピック (cherry-pick) |
| CI/CD | CI/CD (英文そのまま) |
| Pipeline | パイプライン |
| Lint / Linter | リント / リンター |
| Coverage (test coverage) | カバレッジ (テストカバレッジ) |

## §8 keiailab 運用コンテキスト (社内用語)

| English (canonical) | 日本語 |
|---|---|
| keiailab operator family | keiailab オペレーターファミリー |
| operator-commons | operator-commons (英文そのまま — Go モジュール名) |
| supercycle | スーパーサイクル (supercycle) — ユーザー命名 originalsiv |
| Wave (Wave 0 ~ Wave 5) | Wave (英文そのまま) |
| Phase | フェーズ (Phase) |
| Cadence checkpoint | ケイデンスチェックポイント |

## §9 参照

- ユーザー supercycle plan: `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4 「各言語別の *glossary* 1 ファイル」
- deep-petting-cookie plan: `~/.claude/plans/deep-petting-cookie.md` §Phase 0.8 (i18n-glossary 新規)
- Kubernetes 日本語翻訳ガイド: <https://kubernetes.io/ja/docs/contribute/localization_ja/>
- 本ライブラリ ROADMAP §API Stability Tier: `../../ROADMAP.md`
- Korean SSOT: [`glossary-ko.md`](glossary-ko.md) (10 sections × ~120 terms 完全版)

## §10 変更履歴

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 新規 — 日本語標準用語集 placeholder v0.1 (4 sections, ~45 terms) | deep-petting-cookie §Phase 0.8 + supercycle Wave 4 §4.4 sister |
| 2026-05-21 | S4 Phase 1 — Claude による全 10 sections × ~120 terms 拡充 (機械翻訳、native review 待ち) | docs/specs/2026-05-21-i18n-4lang-master-design.md §4.2.1 |
