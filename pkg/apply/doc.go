// SPDX-License-Identifier: MIT

// Package apply 는 K8s 리소스의 idempotent apply 헬퍼를 제공한다 —
// controllerutil.CreateOrUpdate 를 도메인 타입별로 감싸 immutable 필드의
// create-only 가드와 server-default 보존 (nil-guard) 을 일관 적용한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator 의 *resources_apply.go* near-identical
// 중복 (~530 LOC) 을 추출하여 신규 도입되었다. downstream consumer 동시 적용
// 회귀 통과 후 Stable Tier 격상 예정 (자세한 격상 조건은 docs/ROADMAP.md
// §API Stability Tier).
//
// # 의존성 정책
//
// commons 의 *순수 데이터 변환* 패키지 (finalizer, status, security 등) 와
// 달리 본 패키지는 K8s API 호출을 직접 수행한다 — pkg/pvc 와 동일한
// controller-runtime 의존 예외군에 속한다.
//
//   - controller-runtime client.Client + controllerutil.CreateOrUpdate —
//     downstream operator 가 이미 보유한 표준 의존이므로 신규 부담 0.
//   - k8s.io/client-go/util/retry — StatefulSet 의 update conflict 재시도.
//
// # 설계 원칙
//
//   - mutateFn 패턴: 갱신할 필드를 각 함수가 *명시적으로* 선언한다. desired
//     전체 spec 덮어쓰기는 Create 시점에만 — partial desired 로 인한 기존
//     spec 손실 위험 차단.
//   - immutable 필드는 Create 시점에만 설정: Service 의 ClusterIP /
//     IPFamilies / IPFamilyPolicy, StatefulSet 의 Selector / ServiceName /
//     VolumeClaimTemplates / PodManagementPolicy, Deployment 의 Selector.
//   - server-default pointer 필드 nil-guard: desired 가 nil 이면 운영 중
//     객체의 server-defaulted 값을 보존한다. nil 덮어쓰기 ↔ 서버 기본값
//     재주입의 무한 ping-pong (Deployment generation 116k+ 폭주 사고) 차단.
//   - Secret 은 create-once 만 제공 (SecretIfNotExists): 랜덤 credential 을
//     CreateOrUpdate 로 다루면 매 reconcile 마다 재생성되므로 의도적 제외.
//   - CRD 의존 리소스 (ServiceMonitor / Certificate 등 unstructured) 는
//     failSoftNoMatch=true 로 CRD 미설치 클러스터에서 fail-soft 가능 —
//     CRD 설치 후 다음 reconcile 에서 자동 생성된다.
//   - 모든 함수는 controllerutil.SetControllerReference 로 owner 를 일관
//     설정한다 — owner 삭제 시 GC cascade 보장.
package apply
