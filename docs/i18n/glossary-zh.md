# 中文 术语表 (Glossary)

> [English](../../README.md) | [한국어](glossary-ko.md) | [日本語](glossary-ja.md) | **中文**
>
> 本术语表是 keiailab operator family 4 个仓库 (operator-commons + postgres-operator + mongodb-operator + valkey-operator) 的中文翻译时*必须参考*的标准术语表。简体中文 (中国大陆) 表记。
>
> ⚠️ **This translation is AI-generated and pending native review.** — 由 Claude 机器翻译生成。母语审阅者 (native reviewer) 质量验证后升级到 `[x]` 完成状态。本文件所有条目按 `[검토 필요]` (待校验) 状态处理。

## §1 一致性规则

1. **代码标识符保持英文** — 例: `ValkeyCluster`, `kubectl`, `Helm`, `pkg/probes`. 禁止中文翻译。
2. **标准 K8s 术语 = 英文优先 + 中文附注** — 例: `Pod (容器组)`, `Deployment (部署)`. 首次出现时英文 + 括号中文,后续可单独使用中文。
3. **operator-commons API 名称保持英文** — 例: `Reconciler`, `Finalizer`, `EventRecorder`. 可附中文注释 (`Finalizer (终结器)`)。
4. **外部用户可见文档 (README/CONTRIBUTING/SECURITY 等)** = 书面语体. 内部文档 (HANDOFF/AGENTS 等) = 口语体可。
5. **简体中文 GB 标准** — 技术用语优先选用 GB 标准译法,1 个文档内保持一致性。

## §2 Kubernetes 标准术语

| English (canonical) | 中文 |
|---|---|
| CustomResourceDefinition (CRD) | 自定义资源定义 (CRD) |
| Custom Resource (CR) | 自定义资源 (CR) |
| Reconciler | 协调器 (或英文 reconciler) |
| Reconcile Loop | 协调循环 (reconcile 循环) |
| Controller | 控制器 |
| Operator | 操作器 (或英文 operator) |
| Operator Pattern | 操作器模式 |
| Finalizer | 终结器 (或英文 finalizer) |
| Webhook | Webhook (英文为主) |
| Validating Admission Webhook | 验证准入 Webhook |
| Mutating Admission Webhook | 变更准入 Webhook |
| Conversion Webhook | 转换 Webhook |
| Pod | 容器组 (Pod) |
| Deployment | 部署 (Deployment) |
| StatefulSet | 有状态副本集 (StatefulSet) |
| DaemonSet | 守护进程集 (DaemonSet) |
| Job | 任务 (Job) |
| CronJob | 定时任务 (CronJob) |
| Service | 服务 (Service) |
| Ingress | 入口 (Ingress) |
| ConfigMap | 配置映射 (ConfigMap) |
| Secret | 密钥 (Secret) |
| PersistentVolume (PV) | 持久卷 (PV) |
| PersistentVolumeClaim (PVC) | 持久卷声明 (PVC) |
| StorageClass | 存储类 (StorageClass) |
| Namespace | 命名空间 |
| Node | 节点 |
| Cluster | 集群 |
| Control Plane | 控制平面 |
| Worker Node | 工作节点 |
| API Server | API 服务器 |
| Scheduler | 调度器 |
| kubelet | kubelet (英文为主) |
| kube-proxy | kube-proxy (英文为主) |
| etcd | etcd (英文为主) |
| Helm Chart | Helm 图表 (Helm Chart) |
| Helm Release | Helm 发布 |
| Kubernetes RBAC | Kubernetes RBAC |
| ServiceAccount | 服务账号 (ServiceAccount) |
| Role / ClusterRole | 角色 / 集群角色 |
| RoleBinding / ClusterRoleBinding | 角色绑定 / 集群角色绑定 |
| NetworkPolicy | 网络策略 (NetworkPolicy) |
| PodSecurityAdmission (PSA) | Pod 安全准入 (PSA) |
| restricted profile | restricted 配置 (英文为主) |
| baseline profile | baseline 配置 (英文为主) |
| privileged profile | privileged 配置 (英文为主) |
| Probe (liveness/readiness/startup) | 探针 (liveness/readiness/startup) |
| HTTPGetAction | HTTP Get 动作 (HTTPGetAction) |
| TCPSocketAction | TCP 套接字动作 (TCPSocketAction) |
| ExecAction | 执行动作 (ExecAction) |
| InitContainer | 初始化容器 (InitContainer) |
| Sidecar Container | 边车容器 |

## §3 operator-commons 库术语

| English (canonical) | 中文 |
|---|---|
| `pkg/finalizer` | 终结器包 (`pkg/finalizer`) |
| `pkg/labels` | 标签包 (`pkg/labels`) |
| `pkg/status` | 状态条件包 (`pkg/status`) |
| `pkg/version` | 版本兼容性包 (`pkg/version`) |
| `pkg/monitoring` | 监控包 (`pkg/monitoring`) |
| `pkg/networkpolicy` | 网络策略包 (`pkg/networkpolicy`) |
| `pkg/security` | 安全包 (`pkg/security`) |
| `pkg/webhook` | Webhook 包 (`pkg/webhook`) |
| `pkg/probes` (v0.8.0 新增) | 探针构建器包 (`pkg/probes`) |
| `pkg/storageclass` (v0.8.0 新增) | 存储类包 (`pkg/storageclass`) |
| `pkg/events` (v0.8.0 新增) | 事件记录器包 (`pkg/events`) |
| Recorder interface | Recorder 接口 (英文为主) |
| EventType (Normal / Warning) | 事件类型 (Normal / Warning) |
| Reason (event reason) | 事件原因 (Reason) |
| Builder pattern | 构建器模式 |
| Fluent API | 流式 API (Fluent API) |
| API Stability Tier | API 稳定性等级 |
| Stable / Beta / Experimental | 稳定 / Beta / 实验性 (英文推荐 — Stable/Beta/Experimental) |
| Tier 升级 (promotion) | 等级升级 (Tier promotion) |
| Breaking change | 破坏性变更 (Breaking change) |
| Semver (Semantic Versioning) | 语义化版本 (Semver) |
| Deprecated | 已弃用 (Deprecated) |

## §4 reconciler 模式术语

| English (canonical) | 中文 |
|---|---|
| Desired state | 期望状态 |
| Actual state | 实际状态 |
| Reconcile | 协调 (reconcile) |
| Drift | 漂移 (英文为主) |
| Idempotency | 幂等性 |
| Eventual consistency | 最终一致性 |
| Owner reference | 所有者引用 (OwnerReference) |
| Garbage collection | 垃圾回收 (GC) |
| Status condition | 状态条件 (status condition) |
| Available / Ready / Progressing / Degraded | Available / Ready / Progressing / Degraded (英文为主) |
| Failover | 故障切换 (failover) |
| Provisioning | 预置 |
| Rolling update | 滚动更新 |
| Blue-green deployment | 蓝绿部署 |
| Canary deployment | 金丝雀部署 |
| Backup / Restore | 备份 / 恢复 |
| Point-in-time recovery (PITR) | 时间点恢复 (PITR) |
| Horizontal scaling | 水平扩展 |
| Vertical scaling | 垂直扩展 |
| Sharding | 分片 |
| Replication | 复制 (replication) |
| Primary / Replica / Secondary | Primary / Replica / Secondary (英文为主) |

## §5 安全 + 认证术语

| English (canonical) | 中文 |
|---|---|
| Authentication (AuthN) | 认证 (Authentication) |
| Authorization (AuthZ) | 授权 (Authorization) |
| OIDC | OIDC (英文为主) |
| LDAP | LDAP (英文为主) |
| TLS / mTLS | TLS / mTLS (英文为主) |
| Certificate | 证书 (certificate) |
| Service mesh | 服务网格 |
| Tenant / Tenancy | 租户 / 租户机制 |
| Multi-tenancy | 多租户 |
| Encryption at rest | 静态加密 |
| Encryption in transit | 传输加密 |

## §6 运维 + 可观测性术语

| English (canonical) | 中文 |
|---|---|
| Observability | 可观测性 (observability) |
| Metric / Monitoring | 指标 / 监控 |
| Logging | 日志 |
| Tracing | 追踪 |
| Prometheus | Prometheus (英文为主) |
| Grafana | Grafana (英文为主) |
| Alert / Alerting | 告警 / 告警发出 |
| Service Level Objective (SLO) | 服务等级目标 (SLO) |
| Service Level Indicator (SLI) | 服务等级指标 (SLI) |
| Postmortem | 事后分析 (postmortem) |
| Incident | 事件 (incident) |
| Severity (SEV-1 / SEV-2 / SEV-3) | 严重程度 (SEV-1 / SEV-2 / SEV-3、英文为主) |

## §7 治理 + 协作术语

| English (canonical) | 中文 |
|---|---|
| ADR (Architecture Decision Record) | 架构决策记录 (ADR) |
| RFC (Request For Comments) | RFC (英文为主) |
| Roadmap | 路线图 |
| Adopter | 采用者 (Adopter) |
| Maintainer | 维护者 |
| Contributor | 贡献者 |
| Pull Request (PR) | 拉取请求 (PR) |
| Merge Request (MR) | 合并请求 (MR、GitLab) |
| Issue | 议题 |
| Code review | 代码审查 |
| Squash merge | 压缩合并 |
| Rebase | 变基 |
| Cherry-pick | 樱桃挑选 (cherry-pick) |
| CI/CD | CI/CD (英文为主) |
| Pipeline | 流水线 |
| Lint / Linter | 检查 / 检查器 |
| Coverage (test coverage) | 覆盖率 (测试覆盖率) |

## §8 keiailab 运维上下文 (内部术语)

| English (canonical) | 中文 |
|---|---|
| keiailab operator family | keiailab 操作器家族 |
| operator-commons | operator-commons (英文为主 — Go 模块名) |
| supercycle | 超级周期 (supercycle) — 用户命名原本 |
| Wave (Wave 0 ~ Wave 5) | Wave (英文为主) |
| Phase | 阶段 (Phase) |
| Cadence checkpoint | 节奏检查点 |

## §9 参考

- 用户 supercycle plan: `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4 "各语言 *glossary* 1 文件"
- deep-petting-cookie plan: `~/.claude/plans/deep-petting-cookie.md` §Phase 0.8 (i18n-glossary 新增)
- Kubernetes 中文翻译指南: <https://kubernetes.io/zh-cn/docs/contribute/localization_zh/>
- 本库 ROADMAP §API Stability Tier: `../../ROADMAP.md`
- Korean SSOT: [`glossary-ko.md`](glossary-ko.md) (10 sections × ~120 terms 完整版)

## §10 变更历史

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 新增 — 中文标准术语表 placeholder v0.1 (4 sections, ~45 terms, 简体) | deep-petting-cookie §Phase 0.8 + supercycle Wave 4 §4.4 sister |
| 2026-05-21 | S4 Phase 1 — Claude 完整翻译全 10 sections × ~120 terms 扩充 (机器翻译,待母语审阅) | docs/specs/2026-05-21-i18n-4lang-master-design.md §4.2.1 |
