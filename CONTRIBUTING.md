# Contributing to keiailab-commons

> **English** | [한국어](CONTRIBUTING.ko.md) | [日本語](CONTRIBUTING.ja.md) | [中文](CONTRIBUTING.zh.md)

`keiailab/keiailab-commons` is a Go library imported by downstream
Kubernetes operators. All contributions follow
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) and
[docs/GOVERNANCE.md](docs/GOVERNANCE.md).

## Contribution flow

1. **Share intent via an issue or ADR** (for non-trivial changes).
2. **Fork + feature branch** — use `feat/<slug>`, `fix/<slug>`,
   `docs/<slug>`, or `refactor/<slug>`.
3. **Verify the local gates**:
   - `lefthook install --force`
   - `lefthook run pre-commit --all-files`
   - `make lint test`
4. **Open the PR** — Conventional Commits format; English or Korean body
   is acceptable.
5. **Review SLA**: maintainers reply within 24 hours.

## PR checklist (author)

- [ ] PR title: Conventional Commits format (`feat`, `fix`, `docs`,
  `refactor`, `test`, `chore`).
- [ ] PR body: change summary + verification commands + cited output.
- [ ] Unit tests added or updated (mandatory for any change inside
  `pkg/<sub>`).
- [ ] godoc updated when the public API changes.
- [ ] For **public-API breaking changes**: linked ADR plus
  downstream-consumer impact analysis.
- [ ] `go.mod` / `go.sum` drift = 0 (running `go mod tidy` produces no
  change).
- [ ] New dependencies: cite the license and CVE review in the PR body.

## Local development (cross-cut work with a downstream consumer)

When a change touches both `keiailab-commons` and a downstream operator:

```fish
# 1. add a replace directive in the consumer operator's go.mod
#    (local-only; do not commit it)
# go.mod tail:
#   replace github.com/keiailab/keiailab-commons => ../keiailab-commons

# 2. edit both sides + run `go test ./...` on each

# 3. split the PRs:
#    - keiailab-commons side: merge + tag (e.g. v0.9.0)
#    - consumer side: bump the require directive (remove replace)
```

## Adding a new `pkg/<sub>` package

Follow the "intermediate change" process in
[docs/GOVERNANCE.md](docs/GOVERNANCE.md):

1. Open an issue or ADR explaining *why* it belongs in commons and
   *which* downstream consumer will use it.
2. Wait through the 7-day comment window.
3. Merge after multiple maintainer LGTMs.

## Release

- v0.x SemVer: every minor bump represents a public-API change or a
  meaningful behaviour change.
- Release procedure:
  1. `git tag v0.X.Y` (annotated).
  2. `git-cliff` regenerates the CHANGELOG PR.
  3. `git push origin v0.X.Y`.
  4. Open a follow-up PR in each downstream consumer to bump the
     `require` directive.

## Security vulnerabilities

Use the private disclosure process in [SECURITY.md](SECURITY.md). Public
issues are not the right channel.

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
