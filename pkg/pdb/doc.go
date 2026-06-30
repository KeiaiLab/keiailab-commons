// SPDX-License-Identifier: MIT

// Package pdb 는 downstream operator 공통의 PodDisruptionBudget 빌더를
// 제공한다. drain 시 quorum / primary 동시 evict 를 막는 운영 불변식을 단일
// 진실원으로 고정한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 valkey-operator (`resources/pdb.go` BuildPDB / BuildShardPDB) 와
// mongodb-operator (`resources/builder.go` BuildMongoDBPDB / pdbBaseSpec) 의
// 평행 재구현을 추출하여 신규 도입되었다. downstream 동시 적용 회귀 통과 후
// Stable 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
// valkey 와 mongo 의 PDB 기본 정책이 달랐다:
//
//   - default minAvailable floor: valkey=max(replicas-1, 1)(primary 항상 보존) /
//     mongo=max(replicas-1, 0). → Params.DefaultFloor 로 파라미터화(호출자가 1/0 전달).
//   - min/max 우선순위: valkey 는 MaxUnavailable 먼저 검사, mongo 는 MinAvailable
//     먼저. → commons canonical = MinAvailable 우선 (둘 다 설정은 webhook 이 reject
//     하므로 실사용 동작 동일, 무효 입력에서만 차이).
//
// # 의존성 정책
//
// stdlib + k8s.io/api/policy/v1 + apimachinery 만 의존하는 *순수 데이터 변환*
// 패키지 (labels / security / networkpolicy 빌더군과 동일). K8s client 의존 없음.
//
// # 설계 원칙
//
//   - CR 타입 비의존: commons 는 operator 의 CRD 타입을 import 하지 않는다.
//     호출자가 자신의 spec 에서 min/max(*intstr.IntOrString) 를 추출해 Params 로
//     전달한다 (certmanager.CertParams 패턴 정합).
//   - 기본 정책: MinAvailable / MaxUnavailable 모두 미지정 시
//     minAvailable = max(replicas-1, DefaultFloor). DefaultFloor 로 operator 별
//     차이(valkey=1 primary-보존 / mongo=0)를 흡수한다.
//   - opt-in(enabled) 판정은 *호출자 책임* — 비활성 시 호출자가 nil 처리.
package pdb
