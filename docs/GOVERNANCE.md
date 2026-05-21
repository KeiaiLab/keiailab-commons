# Governance

> **English** | [한국어](GOVERNANCE.ko.md)

This document defines the decision-making process for the
`keiailab/operator-commons` library. The library is imported in common
by downstream consumer operators, so changes to the *public API* affect
downstream compatibility.

## Principles

1. **Openness** — every decision is made in a public channel (GitHub
   Issue / PR / ADR).
2. **Lazy Consensus** — routine changes proceed unless there is
   explicit dissent.
3. **Explicit Consensus** — public-API breaking changes, new package
   introductions, and license changes require **2/3 supermajority** of
   maintainers *plus* at least one LGTM from a downstream consumer
   maintainer.
4. **Shared responsibility** — maintainers share responsibility for
   library stability, downstream operational impact, and security.

## Decision categories

### Routine changes (Lazy Consensus)

- Bug fixes, documentation improvements, additional tests, minor /
  patch dependency upgrades, internal refactors that leave the public
  API unchanged.
- Process: PR → at least one maintainer LGTM → merge.
- Window: none — when the local gates pass the change can merge. The
  project does not use GitHub Actions; all quality gates are enforced by
  the local four layers (`lefthook.yml`, `Makefile`, reviewer evidence,
  ADR coverage).

### Intermediate changes (Explicit Consensus)

- Adding a new public-API function or type, major dependency upgrades,
  introducing a new `pkg/<sub>` package.
- Process: issue or ADR proposal → 7-day comment window → majority
  maintainer LGTM → merge.
- One or more dissenting opinions trigger a maintainer discussion.

### Public-API breaking changes (ADR required)

- Function signature change, type removal, module-path change, license
  change.
- Process:
  1. Submit `docs/kb/adr/NNNN-<slug>.md`.
  2. 14-day comment window.
  3. 2/3 maintainer supermajority plus at least one downstream consumer
     LGTM.
  4. Move the ADR `Status: Draft → Accepted` and merge the
     implementation PR.

## Security decisions

CVE reports follow [SECURITY.md](../SECURITY.md). Reports are handled
privately. An embargo is kept until downstream consumers can release a
fix.

## Release decisions

- **v0.x**: a single maintainer may tag minor / patch releases under
  Lazy Consensus.
- **v1.0+ (stable)**: strict SemVer — a major bump requires an ADR plus
  a 2/3 supermajority.

## Change history

| Date | Change |
|---|---|
| 2026-05-09 | Document introduced — governance baseline. |

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
