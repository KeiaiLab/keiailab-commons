# Adopters of operator-commons

본 라이브러리는 *외부 사용자가 직접 import* 하기보다는 keiailab 의 3 operator 가 공통 의존성으로 import 하는 *내부 공유 라이브러리* 입니다. 외부 사용자는 *consumer operator* 를 통해 간접적으로 본 라이브러리의 코드를 사용하게 됩니다.

## Direct consumers (in-org)

| Operator | 사용 패키지 | 시작 버전 | 현재 버전 | 최근 commit | 갱신 일자 |
|---|---|---|---|---|---|
| `keiailab/mongodb-operator` | labels, security, webhook, version, finalizer, networkpolicy | v0.1.0 | **v0.7.0** | `97140db` | 2026-05-20 |
| `keiailab/postgres-operator` | labels, security, webhook, status, version | v0.1.0 | **v0.7.0** | `8c9db39` | 2026-05-20 |
| `keiailab/valkey-operator` | labels, security, webhook, monitoring, finalizer, networkpolicy | v0.1.0 | **v0.6.0** ⚠️ (1 minor lag, I09 upgrade 예정) | `e878420` | 2026-05-20 |

> **라이브 evidence (2026-05-20)**: 본 표는 각 operator 의 `go.mod` 라이브 `require github.com/keiailab/operator-commons <ver>` + `grep -rn "github.com/keiailab/operator-commons" --include="*.go"` import 결과 기반. Codex challenge #1 (Plan Adversarial Review RFC-0045 §2.5) 정합 — valkey 가 v0.6.0 에 머물러 있어 v0.7.0 신규 export 미수용 (I09 별 PR 진행 예정).

## External adopters

본 라이브러리는 *Go module* 로 공개되어 있어 누구나 `go get github.com/keiailab/operator-commons` 으로 사용 가능합니다. 그러나 v0.x 단계에서는 공개 API breaking 이 자유롭게 일어날 수 있으므로 *외부 사용자에게는 v1.0 stable 이후를 권장* 합니다.

외부 사용 사례 등재를 원하시면 PR 로 row 추가:

```markdown
| **<조직 / 프로젝트>** ([profile](<URL>)) | <사용 패키지> | <사용 시작 버전> | <현재 버전> | <등재 일자 YYYY-MM-DD> |
```

## CNCF / 라이선스

- 라이선스: Apache-2.0 only (AGPL/BUSL transitive 의존성 0건 목표)
- 본 ADOPTERS 는 CNCF graduation criteria 의 "≥1 public adopter" 와 동등한 *cross-repo dependency declaration* 으로도 활용됩니다.
