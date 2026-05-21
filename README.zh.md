# operator-commons

> [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) (placeholder) | **中文** (placeholder)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg)](https://pkg.go.dev/github.com/keiailab/operator-commons)

> **注意 (Notice)**: 此中文 README 为占位文档。正式版本请参考 [README.md](README.md) (English)。完整中文翻译需要母语审校,为后续待办项 (RFC-0025 `[~]` partial marker)。

## 概述 (Overview)

**keiailab** Kubernetes operator (`mongodb-operator`, `valkey-operator`, `postgresql-operator`) 共享的 Go 库。

> 状态: **v0.x — API 可能发生破坏性变更**。v1.0 之后将采用 SemVer (语义化版本) 保持稳定。

完整英文正文请参考 [README.md](README.md)。

## 软件包列表 (Packages, v0.8.0)

| Package | 用途 |
|---|---|
| `pkg/version` | 支持的 DB 版本 allowlist 规约 |
| `pkg/security` | PodSecurity *restricted* SecurityContext 构建器 |
| `pkg/labels` | 推荐 Kubernetes 标签 (`app.kubernetes.io/*`) 构建器 |
| `pkg/monitoring` | Prometheus Operator `ServiceMonitor` 构建器 |
| `pkg/networkpolicy` | NetworkPolicy 构建器 — deny-by-default + functional options |
| `pkg/webhook` | Admission 验证助手 |
| `pkg/finalizer` | Finalizer 助手 (不依赖 controller-runtime) |
| `pkg/status` | 4 个标准 Condition Type + 6 个 Reason 目录 + 助手函数 |

详细的 API 签名和使用示例请参考 [README.md](README.md) 的 `Packages`、`Usage` 章节。

## 许可证 (License)

Apache-2.0 — 详见 [LICENSE](./LICENSE)。每个 minor 版本审计零 AGPL/BUSL 传递依赖为目标。

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a> ·
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/forgewise">forgewise</a>
</p>

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
