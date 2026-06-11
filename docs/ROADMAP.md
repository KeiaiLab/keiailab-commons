# ROADMAP — keiailab-commons

> **English** | [한국어](ROADMAP.ko.md) | [日本語](ROADMAP.ja.md) | [中文](ROADMAP.zh.md)

This ROADMAP tracks the library's evolution along three axes: *API
stability tier*, *v1.0.0 graduation criteria*, and *per-package
follow-up items*. The project does not maintain time-based deadlines —
the library evolves according to the needs of its downstream consumers.

## Checkbox meaning

| Marker | Meaning |
|---|---|
| `[x]` | Code + tests both present. Downstream import works. |
| `[~]` | Partial implementation (helper present, verification still open). |
| `[ ]` | Not started. |

## API stability tier (current v0.11.0 candidate)

| Package | Tier | Promotion criterion |
|---|---|---|
| `pkg/finalizer` | **Stable** | v1 entry (no additional work). |
| `pkg/labels` | **Stable** | v1 entry (no additional work). |
| `pkg/status` | **Stable** | v1 entry (no additional work). `update.go` (`UpdateWithRetry`) is a Beta surface. |
| `pkg/storageclass` | **Stable** | Trivial validation surface (regex + nil check). |
| `pkg/version` (incl. `Matrix`) | Beta | Generic `Matrix[E]` cross-repo verify. |
| `pkg/monitoring` | Beta | `ServiceMonitor` cross-downstream equivalence e2e. |
| `pkg/networkpolicy` | Beta | 4-direction (ingress / egress × TCP / UDP) verify. |
| `pkg/security` | Beta | Restricted PSA guard across downstream. |
| `pkg/events` | Beta | Downstream live adoption + reconciliation regression 0. |
| `pkg/pvc` | Beta | Downstream PVC expansion live adoption. |
| `pkg/topology` | Beta | Downstream topology spread live adoption. |
| `pkg/apply` | Beta | Downstream idempotent apply live adoption. |
| `pkg/reconcile` | Beta | Downstream reconcile scaffolding live adoption. |
| `pkg/certmanager` | Beta | Downstream Certificate / Issuer render live adoption. |
| `pkg/reconcilemetrics` | Beta | Downstream live adoption + Prometheus series-name parity. |
| `pkg/webhook` | **Experimental** | Multi-downstream adoption + stabilization. |
| `pkg/probes` | **Experimental** | 2+ downstream adoption → Beta. |
| `pkg/bundle` | **Experimental** | 2+ downstream adoption → Beta. |
| Helm `keiailab.secrets.externalSecret` | Beta | Downstream Valkey/MongoDB/PostgreSQL chart render equivalence. |

**Tier semantics**:

- **Stable** — no BREAKING CHANGE in patch / minor releases. Use
  deprecation: mark, keep for 2 minor releases, then remove. Recorded
  exception: the v0.10.0 module path change (`operator-commons` →
  `keiailab-commons`) was an import-path BREAKING CHANGE — permitted in
  a 0.x minor release (SemVer major-version-zero rule, see
  [UPGRADING.md](UPGRADING.md)).
- **Beta** — BREAKING CHANGE allowed in minor releases (must appear in
  CHANGELOG). API shape is mostly settled.
- **Experimental** — BREAKING CHANGE possible at any release. Callers
  bear the risk.

## v1.0.0 graduation criteria (checklist)

- [ ] All packages reach **Stable** tier.
- [ ] Zero BREAKING CHANGE across six or more consecutive minor releases.
- [ ] godoc coverage ≥ 80 % (`go doc -all ./...` basis).
- [ ] CHANGELOG.md v0.x evolution history + v1.0.0 release notes.
- [ ] CITATION.cff + DOI (academic citation).
- [ ] Downstream consumers import v1.0.0 commons with regression 0.
- [x] `go vet ./... && go test ./...` clean (coverage 96.3 % > 85 %
  threshold).
- [x] API stability promise document — [STABILITY.md](STABILITY.md).
- Verify: downstream consumer CI passes all e2e tests against
  `keiailab-commons v1.0.0`.

## Per-package follow-up

### pkg/finalizer (Stable)

- [x] `Add`, `Remove`, `Contains` helpers — `pkg/finalizer/finalizer.go`.
- [x] Avoids controller-runtime (stdlib `slices` only).
- [x] Unit tests — `pkg/finalizer/finalizer_test.go`.
- [x] Multi-finalizer ordering helper — `pkg/finalizer/order.go`
  `EnsureOrder`.
- Verify: downstream finalizer regression 0.

### pkg/labels (Stable)

- [x] Recommended Kubernetes labels (`app.kubernetes.io/*`) —
  `pkg/labels/labels.go`.
- [x] Component / instance / part-of mapping.
- [x] Unit tests — `pkg/labels/labels_test.go`.
- [x] Recommended labels v2 (K8s 1.30+) — `pkg/labels/v2.go` `AllV2`.
- Verify: downstream `metadata.labels` consistency.

### pkg/status (Stable)

- [x] Condition catalog helper — `pkg/status/conditions.go`.
- [x] `SetAvailable` sugar.
- [x] Unit tests.
- [x] Condition reason catalog documentation —
  `pkg/status/REASONS.md`.
- Verify: downstream `kubectl get <kind> -o yaml`
  `.status.conditions` parity.

### pkg/version (Beta)

- [x] `Matrix[E]` generic — `pkg/version/matrix.go`.
- [x] `SetAvailable` sugar.
- [x] semver-aware version compare — `pkg/version/version.go`.
- [x] Cross-version compatibility test —
  `pkg/version/api_stability_test.go`.
- [x] `Matrix[E]` serializer (JSON / YAML) —
  `pkg/version/serializer.go`.
- [ ] **Tier promotion** → Stable.
- Verify: downstream version validation parity.

### pkg/monitoring (Beta)

- [x] Prometheus ServiceMonitor builder — `pkg/monitoring/monitoring.go`.
- [x] Unit tests.
- [x] PrometheusRule builder (alert / recording shared) —
  `pkg/monitoring/rule.go`.
- [x] OpenTelemetry exporter helper — `pkg/monitoring/otel.go`.
- [ ] Downstream equivalence e2e — same input → same manifest output.
- [ ] **Tier promotion** → Stable.
- Verify: `monitoring_test.go` golden file diff = 0.

### Helm secrets partials (Beta)

- [x] `keiailab.secrets.externalSecret` raw YAML helper — ESO/Infisical
  materialization without CRD vendoring.
- [ ] Downstream render equivalence across Valkey, MongoDB, and
  PostgreSQL operator charts.
- Verify: `helm template` with `externalSecrets.enabled=true` renders
  `external-secrets.io/v1` only when the consumer explicitly opts in.

### pkg/networkpolicy (Beta)

- [x] NetworkPolicy builder — `pkg/networkpolicy/networkpolicy.go`.
- [x] Default-deny + explicit rule helper.
- [x] Unit tests.
- [x] 4-direction verification —
  `pkg/networkpolicy/four_dir_test.go`.
- [x] CIDR + namespace + pod selector combo —
  `pkg/networkpolicy/combo.go`.
- [ ] **Tier promotion** → Stable.
- Verify: a kind cluster applies the NetworkPolicy and observed deny /
  allow paths match the expectation.

### pkg/security (Beta)

- [x] SecurityContext helper (restricted PSA-compliant) —
  `pkg/security/security.go`.
- [x] RBAC helper.
- [x] Unit tests.
- [x] Restricted PSA regression guard —
  `pkg/security/psa_guard_test.go`.
- [x] Pod / Container SecurityContext split —
  `pkg/security/split.go`.
- [x] seccompProfile default helper — `pkg/security/seccomp.go`.
- [ ] **Tier promotion** → Stable.
- Verify: after `kubectl label ns <ns>
  pod-security.kubernetes.io/enforce=restricted`, downstream pods reach
  Ready.

### pkg/webhook (Experimental)

- [x] Webhook utility base — `pkg/webhook/webhook.go`.
- [x] Unit tests.
- [x] Conversion webhook helper — `pkg/webhook/conversion.go`.
- [x] Validation webhook patterns —
  `pkg/webhook/validation_patterns.go`.
- [ ] Multi-downstream live adoption → stabilization.
- [ ] **Tier promotion** → Beta → Stable.
- Verify: 2 or more downstream consumers use the same helper with
  regression 0.

### pkg/storageclass (Stable)

- [x] DNS-1123 subdomain validation —
  `pkg/storageclass/validator.go`.
- [x] Normalize / MustNormalize — empty → nil (cluster default) + trim
  + pointer return.
- [x] 12 unit tests (100 % coverage) —
  `pkg/storageclass/validator_test.go`.
- [ ] Downstream live adoption + regression 0.

### pkg/events (Beta)

- [x] Recorder interface — compatible with `client-go`
  `record.EventRecorder` without importing it.
- [x] Nine Reason constants.
- [x] Emit / Emitf / EmitWarning / EmitWarningf / WrappedError — all
  nil-safe.
- [x] Unit tests (100 % coverage) — `pkg/events/events_test.go`.
- [ ] Downstream live adoption — Event reasons unified across
  reconcile path.
- [ ] **Tier promotion** → Stable.
- Verify: downstream Reconcile path uses commons reason constants with
  regression 0.

### pkg/probes (Experimental)

- [x] Fluent builder — HTTP / HTTPS / TCP / Exec handlers.
- [x] kubelet defaults (Period = 10 s / Timeout = 1 s /
  SuccessThreshold = 1 / FailureThreshold = 3).
- [x] InitialDelay / Period / Timeout negative-clamp to 0.
- [x] `Build()` panics when no handler is set (fail-fast contract).
- [x] Unit tests (100 % coverage) — `pkg/probes/builder_test.go`.
- [ ] 2+ downstream live adoption (Beta criterion).
- [ ] **Tier promotion** → Beta → Stable.

### pkg/pvc (Beta)

- [x] PVC expansion helper — `pkg/pvc/expansion.go`.
- [x] Unit tests — `pkg/pvc/expansion_test.go`.
- [ ] Downstream live adoption with PVC resize regression 0.
- [ ] **Tier promotion** → Stable.

### pkg/topology (Beta)

- [x] PVC topology spread helper — `pkg/topology/spread.go`.
- [x] Unit tests — `pkg/topology/spread_test.go`.
- [ ] Downstream live adoption with spread constraint verification.
- [ ] **Tier promotion** → Stable.

### pkg/apply (Beta)

- [x] Idempotent apply helpers — ConfigMap / Secret / Service /
  StatefulSet / Deployment / NetworkPolicy / PodDisruptionBudget /
  HorizontalPodAutoscaler — `pkg/apply/apply.go`,
  `pkg/apply/workload.go`.
- [x] Immutable-field guards — Service ClusterIP / IPFamilies
  create-only, StatefulSet immutable fields + RetryOnConflict,
  Deployment server-default + revision-annotation preservation,
  `preserveReplicas` option (HPA conflict avoidance).
  controller-runtime dependent (non-leaf package).
- [ ] Downstream live adoption with apply regression 0.
- [ ] **Tier promotion** → Stable.
- Verify: `go test ./pkg/apply/...`

### pkg/reconcile (Beta)

- [x] `Statusable` interface (`client.Object` + `GetConditions` +
  `SetPhase`) — `pkg/reconcile/statusable.go`.
- [x] `ApplyErrorCondition` + `HandleFinalizerCleanup` +
  `SecretIfNotExists` helpers. controller-runtime dependent
  (non-leaf package).
- [ ] Downstream live adoption with reconcile-loop regression 0.
- [ ] **Tier promotion** → Stable.
- Verify: `go test ./pkg/reconcile/...`

### pkg/certmanager (Beta)

- [x] `CertParams` + `BuildCertificate` + `BuildSelfSignedIssuer` +
  `ServiceSANs` — `pkg/certmanager/certificate.go`,
  `pkg/certmanager/issuer.go`.
- [x] Unstructured-based — zero cert-manager CRD Go dependency.
- [ ] Downstream live adoption with Certificate / Issuer render
  regression 0.
- [ ] **Tier promotion** → Stable.
- Verify: `go test ./pkg/certmanager/...`

### pkg/reconcilemetrics (Beta)

- [x] `ReconcileMetrics` (Total / Latency / Errors) + `New(subsystem)`
  + `MustRegister` — subsystem injection preserves existing operator
  Prometheus time-series names —
  `pkg/reconcilemetrics/reconcilemetrics.go`.
- [x] `IncTotal` / `ObserveReconcile` / `IncError` / `DeleteFor` /
  `ResultFor` helpers.
- [ ] Downstream live adoption with time-series name parity.
- [ ] **Tier promotion** → Stable.
- Verify: `go test ./pkg/reconcilemetrics/...`

### pkg/bundle (Experimental)

- [x] Bundle annotations — six required registry+v1 annotation constants
  plus `NewAnnotations` builder with `Map()` and `DockerLabels()`.
- [x] FBC schema types — Go structs for `olm.package`, `olm.channel`,
  `olm.bundle`, `olm.deprecations` with JSON serialization.
- [x] Bundle directory validation — `ValidateDir(path)` checks
  `manifests/` + `metadata/` + `annotations.yaml`.
- [x] Unit tests (≥ 85 % coverage).
- [ ] 2+ downstream live adoption (Beta criterion).
- [ ] **Tier promotion** → Beta → Stable.
- Verify: downstream operator bundle build uses commons annotations
  with regression 0.

## Dependency policy

- **Kubernetes API only** — `k8s.io/api`, `k8s.io/apimachinery`,
  `k8s.io/utils`. controller-runtime dependency *must not be added at
  leaf packages*.
- **MIT-compatible licenses only** — every dependency addition
  requires an ADR.
- **Complete godoc** — every new public API requires godoc.

## Governance / tracking

- **CHANGELOG.md** — auto-generated by `git-cliff`. Strict semantic
  versioning.
- **CITATION.cff** — academic citation. DOI issued at v1.0.0.
- **ADR** — `docs/kb/adr/` tracks every design decision.
- **AGENTS.md** — AI-collaboration runbook.

## Non-Goals (deliberately out of scope)

- ❌ **controller-runtime dependency** — leaf packages must remain
  controller-runtime free.
- ❌ **downstream-specific logic** — operator-specific code lives in
  the caller's repository. The library ships only *shared* helpers.
- ❌ **Time-based roadmap** — use a feature checklist plus completion
  percentages.
- ❌ **GitHub Actions release gates** — delegate to the local four
  layers.
- ❌ **Plugin / extension SDK positioning** — this is a library, not a
  framework.
- ❌ **Premature v1.0.0** — stay in v0.x until the graduation
  criteria are met.

## Adopters

| Repo | Packages used | Import version |
|---|---|---|
| `mongodb-operator` | finalizer / version / webhook / pvc / topology / security | v0.10.0 (v0.11.0 migration planned) |
| `postgres-operator` | topology / pvc / status / security / version / webhook | v0.10.0 (v0.11.0 migration planned) |
| `valkey-operator` | finalizer / version / security / pvc / networkpolicy / monitoring | v0.10.0 (v0.11.0 migration planned) |

## Change history

| Date | Change | Refs |
|---|---|---|
| 2026-06-11 | v0.11.0 candidate: four new Beta packages (`pkg/apply` / `pkg/reconcile` / `pkg/certmanager` / `pkg/reconcilemetrics`) + `pkg/status` `UpdateWithRetry` Beta surface + Adopters table + v0.10.0 module-path exception note. | v0.11.0 / [UPGRADING.md](UPGRADING.md) |

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
