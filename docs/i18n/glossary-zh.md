# 中文 术语表 (Glossary)

> [English](../../README.md) | [한국어](glossary-ko.md) | [日本語](glossary-ja.md) (placeholder) | **中文**
>
> 本术语表是 keiailab operator family 4 个仓库 (operator-commons + postgres-operator + mongodb-operator + valkey-operator) 的中文翻译时*必须参考*的标准术语表。简体中文 (大陆) 表记。
>
> **状态**: `[~]` 部分实现 (placeholder) — RFC-0025 §1.2 复选框含义. native reviewer 质量验证后升级到 `[x]` 完成状态.

## §1 一致性规则

1. **代码标识符保持英文** — 例: `ValkeyCluster`, `kubectl`, `Helm`, `pkg/probes`. 禁止中文翻译.
2. **标准 K8s 术语 = 英文优先 + 中文附注** — 例: `Pod (容器组)`, `Deployment (部署)`. 首次出现时英文 + 括号中文,后续可单独使用中文.
3. **operator-commons API 名称保持英文** — 例: `Reconciler`, `Finalizer`, `EventRecorder`.
4. **外部用户可见文档 (README/CONTRIBUTING/SECURITY 等)** = 书面语体. 内部文档 = 口语体可.
5. **简体中文 GB 标准** — 技术用语优先选用 GB 标准译法.

## §2 Kubernetes 标准术语 (placeholder — native reviewer 후 확장)

| English (canonical) | 中文 |
|---|---|
| CustomResourceDefinition (CRD) | 自定义资源定义 (CRD) |
| Custom Resource (CR) | 自定义资源 (CR) |
| Reconciler | 协调器 (或英文 reconciler) |
| Reconcile Loop | 协调循环 |
| Controller | 控制器 |
| Operator | 操作器 (或英文 operator) |
| Finalizer | 终结器 |
| Webhook | Webhook (英文为主) |
| Pod | 容器组 (Pod) |
| Deployment | 部署 (Deployment) |
| StatefulSet | 有状态副本集 (StatefulSet) |
| Service | 服务 (Service) |
| ConfigMap | 配置映射 (ConfigMap) |
| Secret | 密钥 (Secret) |
| PersistentVolume / PVC | 持久卷 / PVC |
| StorageClass | 存储类 (StorageClass) |
| Namespace | 命名空间 |
| Cluster | 集群 |
| Helm Chart | Helm 图表 |
| RBAC | RBAC (英文) |
| ServiceAccount | 服务账号 (ServiceAccount) |
| NetworkPolicy | 网络策略 (NetworkPolicy) |
| PodSecurityAdmission (PSA) | Pod 安全准入 (PSA) |
| Probe (liveness/readiness) | 探针 (liveness/readiness) |

## §3 operator-commons 库术语 (placeholder)

| English (canonical) | 中文 |
|---|---|
| `pkg/probes` (v0.8.0 新增) | 探针构建器包 (`pkg/probes`) |
| `pkg/storageclass` (v0.8.0 新增) | 存储类包 (`pkg/storageclass`) |
| `pkg/events` (v0.8.0 新增) | 事件记录器包 (`pkg/events`) |
| API Stability Tier | API 稳定性等级 |
| Stable / Beta / Experimental | 稳定 / Beta / 实验性 (英文推荐) |
| Builder pattern | 构建器模式 |
| Fluent API | 流式 API |

## §4 reconciler 模式术语 (placeholder)

| English (canonical) | 中文 |
|---|---|
| Desired state | 期望状态 |
| Actual state | 实际状态 |
| Reconcile | 协调 (reconcile) |
| Drift | 漂移 |
| Idempotency | 幂等性 |
| Eventual consistency | 最终一致性 |
| Failover | 故障切换 |
| Rolling update | 滚动更新 |

## §5 Placeholder 状态 (RFC-0025 §1.2 [~])

本文件处于*placeholder* 状态. 达到 Korean glossary (`glossary-ko.md`) 的完整结构 (10 sections × ~120 terms) 需要:

- [x] §1 一致性规则 (基本) — 本 PR
- [x] §2 K8s 标准术语 (~25 terms) — 本 PR (placeholder)
- [x] §3 operator-commons 术语 (~8 terms) — 本 PR (placeholder)
- [x] §4 reconciler 模式 (~8 terms) — 本 PR (placeholder)
- [ ] §5 安全 + 认证术语
- [ ] §6 运维 + 可观测性术语
- [ ] §7 治理 + 协作术语
- [ ] §8 keiailab 内部上下文术语
- [ ] §9 参考
- [ ] §10 变更历史

**native reviewer 必要事项**:
- 简体中文 (大陆) GB 标准译法整合
- 技术术语 (CRD / Reconciler 等) 的行业标准译名一致性
- 中文专有名词 vs 英文保留判断的一致性

## §6 参考

- 用户 supercycle plan: `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4
- Korean SSOT: [`glossary-ko.md`](glossary-ko.md) (~120 terms 完整版, 2026-05-21 commit `4a174b2`)
- Kubernetes 中文翻译指南: <https://kubernetes.io/zh-cn/docs/contribute/localization_zh/>
- 本库 ROADMAP §API Stability Tier: `../../ROADMAP.md`

## §7 变更历史

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 新增 — 中文标准术语表 placeholder v0.1 (4 sections, ~45 terms, 简体) | deep-petting-cookie §Phase 0.8 + supercycle Wave 4 §4.4 sister |
