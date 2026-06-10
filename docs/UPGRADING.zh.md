# Upgrading keiailab-commons

> [English](UPGRADING.md) | [한국어](UPGRADING.ko.md) | [日本語](UPGRADING.ja.md) | **中文**

> ⚠️ This translation is AI-generated and pending native review.

本文档汇总在 bump `github.com/keiailab/keiailab-commons` Go
module 的 minor 或 major 版本时所需的迁移步骤。它是
下游 consumer 的通用入口点。

## 0. 版本政策（SemVer）

| 变更类型 | SemVer bump | 示例 |
|---|---|---|
| 新包添加 | minor (v0.X → v0.X+1) | `pkg/events`、`pkg/storageclass` 引入 |
| 现有 API signature 变更（breaking） | major (v0.X → v1.0 / v1.X → v2.0) | `pkg/status.SetReady()` signature 变更 |
| 包内行为变更（non-breaking） | patch (v0.X.Y → v0.X.Y+1) | bug 修复 |
| ADR 偏离 | major + Deprecated notice | API stability tier 变更 |

API stability tier（`pkg/<name>/doc.go` marker）：

- **Stable** — minor release 之间向后兼容。
- **Beta** — 可能在下一个 minor 中变更。
- **Experimental** — 任何时候都可能变更。

## 1. v0.7.x → v0.8.x

### Helm library chart consumer

```bash
helm dep update charts/<your-operator>
helm template <your-operator> charts/<your-operator>
```

`keiailab-commons` chart v0.8.0 partial（`_servicemonitor.tpl`、
`_rbac.tpl`、`_networkpolicy.tpl`）无需额外工作。

### Go module consumer

```bash
go get github.com/keiailab/keiailab-commons@v0.8.0
go mod tidy
```

无额外工作 — 向后兼容。

## 2. v0.8.x → v0.9.x

### 新包（minor bump）

| Package | Purpose | Tier |
|---|---|---|
| `pkg/pvc` | PVC expansion helper | Beta |
| `pkg/topology` | PVC topology spread + zone-aware affinity | Beta |

### 迁移

在下游 operator 中添加 import：

```go
import (
    "github.com/keiailab/keiailab-commons/pkg/pvc"
    "github.com/keiailab/keiailab-commons/pkg/topology"
)
```

### 向后兼容性

- 现有包（`pkg/status`、`pkg/finalizer`、`pkg/networkpolicy`、
  `pkg/monitoring`、`pkg/probes`、`pkg/labels`、`pkg/storageclass`、
  `pkg/webhook`、`pkg/events`、`pkg/security`、`pkg/version`）保留
  其 signature。
- `keiailab-commons` chart 的 `_security.tpl` 和
  `_servicemonitor.tpl` partial 是 *opt-in*；保持现有
  inline 定义不变没有任何影响。

### 推荐迁移流程

```bash
# 1. bump 依赖
go get github.com/keiailab/keiailab-commons@v0.9.0
go mod tidy

# 2. 验证
make verify  # lint + test + build

# 3. e2e (kind)
kind create cluster
helm install <operator> charts/<operator>
kubectl apply -f config/samples/
kubectl get <CR> -A  # 观察 reconciliation
```

## 3. v0.9.x → v0.10.x

### Repository / module rename

从 `v0.10.0` 开始，Go module path 使用：

```bash
github.com/keiailab/keiailab-commons
```

下游 operator 需要同时更新 import path 和 dependency pin：

```bash
go get github.com/keiailab/keiailab-commons@v0.10.0
go mod tidy
```

既有 `v0.9.x` tag 声明的是 `github.com/keiailab/operator-commons`
module path，因此不能通过新的 module path 消费。

## 4. v0.9.x → v1.0.0

当满足 v1.0.0 graduation 标准（参见
[STABILITY.md](STABILITY.zh.md) "v1.0.0 graduation"）时推进：

- 所有包达到 Stable tier。
- v0.x → v1.0 是 *naming* 变更 — 语义不变（无
  breaking change）。

## 5. 通用迁移检查清单

升级前：

- [ ] `go mod tidy` 无变化（drift = 0）。
- [ ] `make audit` 通过（govulncheck CVE = 0）。
- [ ] 现有 e2e suite 通过。

升级后：

- [ ] 下游 import path 更新（`go get -u` 或固定版本）。
- [ ] `make verify` 通过。
- [ ] e2e 通过。
- [ ] Helm chart `charts/<operator>` `dependencies:` 更新。

## 6. Breaking-change 通知政策

- **Deprecation**：在新 minor 中添加 `// Deprecated:` 注释；两个
  minor 之后移除。
- **Breaking**：major bump + 本文件中的专用 section + ADR。
- **无静默 breaking change**：每个 breaking change 都至少有
  一个 minor 的事前 deprecation。

## 参考

- ADR index：[`docs/kb/adr/INDEX.md`](kb/adr/INDEX.zh.md)。
- API stability：`pkg/<name>/doc.go` tier marker。
- i18n：[`docs/i18n/README.md`](i18n/README.zh.md)（多语言政策）。

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
