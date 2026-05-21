# 日本語 用語集 (Glossary)

> [English](../../README.md) | [한국어](glossary-ko.md) | **日本語** | [中文](glossary-zh.md) (予定)
>
> 本用語集は keiailab operator family 4 リポジトリ (operator-commons + postgres-operator + mongodb-operator + valkey-operator) の日本語翻訳時に*必ず参照*する標準用語集です。
>
> **状態**: `[~]` 部分実装 (placeholder) — RFC-0025 §1.2 体크박스 의미. native reviewer による品質検証後 `[x]` 完了状態へ昇格.

## §1 一貫性ルール

1. **コード識別子は英文そのまま** — 例: `ValkeyCluster`, `kubectl`, `Helm`, `pkg/probes`. 日本語翻訳禁止.
2. **標準 K8s 用語は英文優先 + 日本語併記** — 例: `Pod (ポッド)`, `Deployment (デプロイメント)`. 本文初出時に英文 + 括弧日本語、以降は日本語単独可.
3. **operator-commons API 名称は英文そのまま** — 例: `Reconciler`, `Finalizer`, `EventRecorder`.
4. **外部ユーザー可視文書 (README/CONTRIBUTING/SECURITY 등)** = 敬体 (`です/ます調`). 内部文書 = 常体可.
5. **敬体と常体の混在禁止** — 1 文書内一貫性.

## §2 Kubernetes 標準用語 (placeholder — native reviewer 後拡張)

| English (canonical) | 日本語 |
|---|---|
| CustomResourceDefinition (CRD) | カスタムリソース定義 (CRD) |
| Custom Resource (CR) | カスタムリソース (CR) |
| Reconciler | リコンサイラー (または英文そのまま) |
| Reconcile Loop | 調整ループ (reconcile ループ) |
| Controller | コントローラー |
| Operator | オペレーター |
| Finalizer | ファイナライザー |
| Webhook | ウェブフック |
| Pod | ポッド (Pod) |
| Deployment | デプロイメント (Deployment) |
| StatefulSet | ステートフルセット (StatefulSet) |
| Service | サービス (Service) |
| ConfigMap | コンフィグマップ (ConfigMap) |
| Secret | シークレット (Secret) |
| PersistentVolume / PVC | 永続ボリューム / PVC |
| StorageClass | ストレージクラス (StorageClass) |
| Namespace | ネームスペース |
| Cluster | クラスター |
| Helm Chart | ヘルムチャート (Helm Chart) |
| RBAC | RBAC (英文そのまま) |
| ServiceAccount | サービスアカウント (ServiceAccount) |
| NetworkPolicy | ネットワークポリシー (NetworkPolicy) |
| PodSecurityAdmission (PSA) | ポッドセキュリティアドミッション (PSA) |
| Probe (liveness/readiness) | プローブ (liveness/readiness) |

## §3 operator-commons ライブラリ用語 (placeholder)

| English (canonical) | 日本語 |
|---|---|
| `pkg/probes` (v0.8.0 新規) | プローブビルダーパッケージ (`pkg/probes`) |
| `pkg/storageclass` (v0.8.0 新規) | ストレージクラスパッケージ (`pkg/storageclass`) |
| `pkg/events` (v0.8.0 新規) | イベントレコーダーパッケージ (`pkg/events`) |
| API Stability Tier | API 安定性ティア |
| Stable / Beta / Experimental | 安定 / ベータ / 実験的 (英文推奨) |
| Builder pattern | ビルダーパターン |
| Fluent API | フルーエント API |

## §4 reconciler パターン用語 (placeholder)

| English (canonical) | 日本語 |
|---|---|
| Desired state | 望ましい状態 |
| Actual state | 実際の状態 |
| Reconcile | 調整 (reconcile) |
| Drift | ドリフト |
| Idempotency | 冪等性 |
| Eventual consistency | 結果整合性 |
| Failover | フェイルオーバー |
| Rolling update | ローリングアップデート |

## §5 Placeholder 状態 (RFC-0025 §1.2 [~])

本ファイルは*placeholder* 状態であり、Korean glossary (`glossary-ko.md`) の本格的な構造 (10 sections × ~120 terms) に到達するためには:

- [x] §1 一貫性ルール (基本) — 本 PR
- [x] §2 K8s 標準用語 (~25 terms) — 本 PR (placeholder)
- [x] §3 operator-commons 用語 (~8 terms) — 本 PR (placeholder)
- [x] §4 reconciler パターン (~8 terms) — 本 PR (placeholder)
- [ ] §5 セキュリティ + 認証用語
- [ ] §6 運用 + 観察用語
- [ ] §7 ガバナンス + 協業用語
- [ ] §8 keiailab 社内コンテキスト用語
- [ ] §9 参照
- [ ] §10 変更履歴

**native reviewer 必要事項**:
- 業務敬語 (`です/ます`) の自然な fluency 検証
- 技術用語 (CRD / Reconciler 等) の業界標準訳に整合
- カタカナ表記 vs 英文そのまま判断の一貫性

## §6 参照

- 사용자 supercycle plan: `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4
- Korean SSOT: [`glossary-ko.md`](glossary-ko.md) (~120 terms 完全版, 2026-05-21 commit `4a174b2`)
- Kubernetes 日本語翻訳ガイド: <https://kubernetes.io/ja/docs/contribute/localization_ja/>
- 本ライブラリ ROADMAP §API Stability Tier: `../../ROADMAP.md`

## §7 変更履歴

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 新規 — 日本語標準用語集 placeholder v0.1 (4 sections, ~45 terms) | deep-petting-cookie §Phase 0.8 + supercycle Wave 4 §4.4 sister |
