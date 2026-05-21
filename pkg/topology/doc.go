// Package topology 는 keiailab operator (mongodb / postgres / valkey) 의 공통
// TopologySpreadConstraints 헬퍼를 제공한다. 현재 단일 기능:
// HA out-of-box default 자동 주입 (zone + node 2-축).
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 Sprint 1 (2026-05-21) 에서 3 operator 의 *defaultedTopologySpread*
// 중복 코드 (~135 LOC) 를 추출하여 신규 도입되었다. 3 consumer 동시 적용 회귀
// 통과 후 Stable Tier 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 설계 원칙
//
//   - controller-runtime / K8s client 미의존 — 순수 데이터 변환.
//   - 입력 (user TSC, replicas, label selector) → 출력 ([]TopologySpreadConstraint).
//   - 사용자가 user-provided TSC 명시 시 그대로 우선 (override 보장).
//   - 임계값 (최소 replica 수) 은 함수형 옵션 — operator 별 의미론 차이를 표현.
//
// # operator 별 임계값 매핑 (참고)
//
// downstream operator:
//   - Spec.Shards.Replicas 가 "*추가* 복제본 수" 의미 → 1 이상이면 총 pod 2 이상.
//   - WithMinReplicas(1).
//
// downstream operator:
//   - Spec.Members 가 "*전체* 멤버 수" 의미 → 2 이상.
//   - WithMinReplicas(2) (기본값).
//
// downstream operator:
//   - Spec.Replicas 가 "*전체* 복제본 수" 의미 → 2 이상.
//   - WithMinReplicas(2) (기본값).
package topology
