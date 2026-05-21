# Security Policy

> **English** | [한국어](SECURITY.ko.md) | [日本語](SECURITY.ja.md)

`keiailab/operator-commons` is imported by downstream Kubernetes
operators. A vulnerability in this library can directly affect the
operational security of those downstream consumers.

## Reporting a vulnerability

**Do not file a public issue.**

### Channels

Use one of the following private channels:

1. **GitHub Security Advisory** (preferred):
   `https://github.com/keiailab/operator-commons/security/advisories/new`
2. **Email**: `security@keiailab.com` (PGP optional):
   - PGP fingerprint:
     `89A4 0947 6828 CB99 2338  C378 651E 51AF 520B CB78`.

### What to include

- Affected version (release tag or commit SHA).
- Affected package (`pkg/security`, `pkg/webhook`, etc.).
- Reproduction steps (a minimal repro if possible; declare it when the
  reproduction depends on a downstream environment).
- Impact assessment — the scope of downstream consumer impact.
- A self-assessed CVSS score, if available.

## Response SLAs

| Stage | Time |
|---|---|
| Initial response (acknowledgement) | within 72 hours |
| Severity assessment | within 7 days |
| Patch release | severity-dependent (Critical: 14 days, High: 30, Medium: 60) |
| Public disclosure | 14 days after the patch (or the earliest point at which downstream consumers can release a fix) |

## Embargo handling

Vulnerabilities that affect the public API are embargoed until
downstream consumers can release fixes concurrently. Maintainers share a
private advisory with downstream maintainers ahead of disclosure.

## Supported versions

| Version | Supported |
|---------|-----------|
| 0.x (alpha) | ✅ latest minor only |
| 1.0+ (stable) | TBD — updated after the first stable release |

The library is currently in v0.x. Public APIs may break; security
patches are released only against the latest minor.

## Dependency security

When a dependency is added or upgraded, the PR body cites the license
and CVE review. Dependabot / Renovate automatic-update PRs are
prioritised for review.

## License / supply chain

This library is **Apache-2.0 only**, with a charter goal of zero AGPL /
BUSL transitive dependencies (`docs/kb/adr/0001-charter.md`). A license
audit runs at every minor release.

## Best practices for downstream consumers

Operators that import this library should:

1. **Use `pkg/security`** — call the restricted PodSecurity
   SecurityContext builder rather than rolling your own.
2. **Use `pkg/webhook`** — do not re-implement version validation.
3. **Use `pkg/networkpolicy`** — deny-by-default NetworkPolicy builder.
4. Track the latest patch of
   `github.com/keiailab/operator-commons` in `go.mod` (Renovate
   automatic PRs).

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
