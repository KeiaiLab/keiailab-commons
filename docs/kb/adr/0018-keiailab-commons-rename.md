# ADR-0018: keiailab-commons repository and module rename

| Meta | Value |
|---|---|
| Status | Accepted |
| Date | 2026-06-11 |
| Author | keiailab |
| Supersedes | (none) |
| Related | ADR-0001, ADR-0005, ADR-0014 |

## Context

The shared operator library was originally published as `operator-commons`,
while its Helm library chart already used the clearer product-facing name
`keiailab-commons`. Downstream operators therefore had two names for one
responsibility boundary: Go imports used `operator-commons`, while Helm
consumers used `keiailab-commons`.

This mismatch made GitOps and ArtifactHub guidance harder to audit because the
same shared layer appeared under different names.

## Decision

Rename the GitHub repository and Go module to `keiailab-commons`:

- repository: `github.com/keiailab/keiailab-commons`
- Go module: `github.com/keiailab/keiailab-commons`
- Helm library chart: `keiailab-commons`
- OCI chart package: `oci://ghcr.io/keiailab/charts/keiailab-commons`

The library remains a shared implementation package only. It does not own
operator CRDs, controllers, workload instances, or environment-specific GitOps
values.

## Consequences

- Downstream operators must update imports and `go.mod` to the new module path.
- A new Go module tag is required because old tags declare the old module path.
- GitHub's repository redirect is sufficient for human navigation, but new
  documentation and dependencies must use the new name.
- Helm chart consumers continue to use the existing `keiailab-commons` chart
  name.

## Verification

```bash
go test ./...
helm lint charts/keiailab-commons
helm package charts/keiailab-commons
helm show chart oci://ghcr.io/keiailab/charts/keiailab-commons --version 0.8.0
```

## Refs

- GitHub repository: <https://github.com/keiailab/keiailab-commons>
- ArtifactHub/GHCR OCI namespace: `oci://ghcr.io/keiailab/charts/*`
