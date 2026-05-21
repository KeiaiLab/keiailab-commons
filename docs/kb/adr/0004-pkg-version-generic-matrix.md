# ADR-0004: pkg/version Generic Matrix[E] 추가 (postgres PR-B3 prerequisite)

- Date: 2026-05-09
- Status: Accepted
- Authors: @eightynine01
- Refs: 기존 downstream 의 downstream operator 의 `internal/version/matrix.go`

## Context

operator-commons v0.5.0 의 `pkg/version` 은 단순 `[]string` 화이트리스트
(`List` + `MustList` + `IsSupported`) 만 제공. downstream operator + valkey-
operator 가 채택했으나 downstream operator 는 *rich entry* (Major / Image
/ Channel / FeatureGate) 가 필요하여 자체 `internal/version/matrix.go`
를 별도 구현 — *commons 채택률 67%*.

기반 — commons 에 generic `Matrix[E]` 추가
하여 postgres 의 `Combo` 같은 struct entry 도 commons 위임 가능하게 한다.

## Decision

1. **`pkg/version/matrix.go` 신규**:
   - `MatrixEntry` interface — `PrimaryKey() string` 의무.
   - `Matrix[E MatrixEntry]` generic struct — element 가 MatrixEntry
     구현 필수.
   - `MustMatrix[E](entries...)` constructor — 빈 리스트 / 빈 PrimaryKey
     / duplicate PrimaryKey 시 init-time panic.
   - `Find(key) (E, bool)` / `IsSupported(key) bool` / `Entries() []E` /
     `Keys() []string` / `Len() int`.

2. **`List` 와 공존**:
   - 기존 `List` API 무변경 — 단순 string 화이트리스트 사용자 (mongodb /
     valkey) 영향 없음.
   - `Matrix[E]` 는 *추가 표면* — semver / channel / feature gate 같은
     메타데이터 필요한 operator (postgres 등) 만 사용.

3. **방어 복사 보존**: `Entries()` 와 `Keys()` 가 internal slice 보호 —
   호출자 mutation 이 commons state 영향 없음.

4. **commons v0.7.0 으로 bump 예정**: API 추가 (semver minor). consumer
   는 `go get @v0.7.0` 후 `pkg/version.MustMatrix` 사용 가능.

## Consequences

### Positive

- postgres `internal/version/matrix.go` 가 commons `Matrix[Combo]` 로
  위임 가능 — *downstream consumer pkg/version 채택률 100%* (PR-B3 후).
- generic 으로 다양한 operator 의 *rich entry* 패턴 표준화 — valkey 의
  `Combo`-ish 진화 후속 가능.
- 기존 `List` 호환 보존 — semver minor bump 만으로 도입.

### Negative

- API 표면 +6 (1 interface + 1 struct + MustMatrix + 5 method). 단
  generic + 단순 method 로 학습 비용 미미.
- Generic 사용으로 Go 1.18+ 강제 — commons go.mod 의 `go 1.25.10`
  와 정합.

### Trade-offs

- *generic Matrix[E]* (본 ADR) vs *interface{} + type assertion* —
  후자는 type safety 부재. 본 ADR 의 generic 이 우위.
- *commons 추가* (본 ADR) vs *postgres 자체 유지* — 후자는 downstream consumer
  cross-cut 변경 시 동기화 부담 (tooling unification 정책 §3.3 lint 위반).

## Alternatives Considered

1. **`Matrix[K comparable, V any]` (key + value 분리)** — 거부.
   - postgres `Combo` 가 *struct 자체가 entry* — key 분리 불필요.
   - `MatrixEntry` interface 가 더 자연스러운 abstraction.

2. **`List` 를 generic 으로 확장** — 거부.
   - 기존 List 사용자 모두 generic 추가 의무 — breaking change. semver
     major bump 필요.
   - `Matrix[E]` *별도 type* 이 호환 보존.

3. **commons 외부 lib (`google/btree` 등) 채택** — 거부.
   - 단순 화이트리스트에 BTree overkill.
   - commons zero-dep 원칙 위반 (k8s.io 외 외부 의존).

## Refs

- 기반 (postgres matrix.go → commons generic Matrix[Combo]).
- downstream operator `internal/version/matrix.go` (현 자체 구현).
- 후속 PR-B3 (postgres): matrix.go 가 commons `Matrix[Combo]` 로 위임.
- commons `pkg/version/version.go` (기존 List API).
