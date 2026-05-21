# Maintainers

> **English** | [한국어](MAINTAINERS.ko.md) | [日本語](MAINTAINERS.ja.md)

This document records the maintainers with decision-making authority
for `keiailab/operator-commons`.

## Current maintainers

| Name / team | GitHub | Role | Scope |
|---|---|---|---|
| keiailab maintainers | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | Lead | All |

The GitHub team `@keiailab/maintainers` holds the merge and release-tag
authority for every area of the library.

## Maintainer qualifications

Either a maintainer of a downstream consumer operator, *or* a
contributor who has met the following criteria for at least six months:

- ≥ 10 merged PRs (library PR cadence is lower than a typical operator,
  so the bar is roughly half of that).
- ≥ 20 reviewed PRs (downstream consumer PRs may count).
- Deep familiarity with at least one `pkg/` area (security, labels,
  webhook, monitoring, networkpolicy, version, status, finalizer,
  storageclass, events, probes, pvc, topology).

## Addition procedure

1. An existing maintainer (or the candidate themselves) opens an issue
   or ADR.
2. The `@keiailab/maintainers` team applies lazy consensus (7-day
   comment window).
3. With no dissent, the candidate is added to the GitHub team and this
   file is updated via PR.

## Inactive maintainers

A maintainer who has not been active for six consecutive months is
moved to emeritus (rights revoked, name retained on the honorary roll).

## Cross-repo agreement

A *public-API breaking change* requires LGTM from a downstream consumer
maintainer at the ADR stage — see [GOVERNANCE.md](GOVERNANCE.md).

## i18n document owners

| Language | Owner | Files | Responsibility |
|---|---|---|---|
| English (canonical) | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | `README.md` and canonical docs | Source of truth |
| Korean | TaeHwan Park ([@eightynine01](https://github.com/eightynine01)) | `README.ko.md` and `*.ko.md` | EN canonical sync + translation review |
| Japanese | (recruiting — volunteer via an issue) | `*.ja.md` | AI translation + native review |
| Chinese | (recruiting — volunteer via an issue) | `*.zh.md` | AI translation + native review |

**Drift verification**: `bash scripts/check-readme-sync.sh` — checks
that the file exists, the section header count matches, line counts
differ by less than the per-language threshold, and cross-links are
bidirectional. The lefthook `pre-push` hook `readme-i18n-sync` enforces
this automatically.

## Emeritus

(none yet)

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
