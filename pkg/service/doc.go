// SPDX-License-Identifier: MIT

// Package service 는 downstream operator 공통의 Kubernetes Service 조립 빌더를
// 제공한다. headless(pod-to-pod stable DNS) / client(ClusterIP) Service 의
// ObjectMeta + Spec 조립을 단일 진실원으로 고정한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 valkey-operator (`resources/service.go` BuildHeadlessService /
// BuildClientService / BuildPrimaryService / BuildMetricsService) 와
// mongodb-operator (`resources/builder.go` BuildHeadlessService /
// BuildClientService) 의 평행 재구현을 추출하여 신규 도입되었다. downstream
// 동시 적용 회귀 통과 후 Stable 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
//   - valkey: ServiceSpec(Labels/Annotations/Type/IPFamilyPolicy/IPFamilies) 옵션
//     병합 + clientPorts 헬퍼 + 4종(headless/client/primary/metrics). 풍부.
//   - mongo: 옵션 없는 단순 2종(headless/client), 단일/이중 포트.
//   - canonical: commons 는 *조립* 만 책임지는 Build(Params) 1개를 제공한다.
//     포트 구성·옵션 병합·서비스 종류 선택 같은 도메인 로직은 호출자에 잔류하고,
//     완성된 labels/annotations/selector/ports/type 을 Params 로 받아 Service 를
//     조립한다. Headless=true 면 ClusterIP="None" + PublishNotReadyAddresses=true.
//
// # 의존성 정책
//
// k8s.io/api/core/v1 + apimachinery 만 의존하는 *순수 데이터 변환* 패키지
// (labels / pdb / hpa 빌더군과 동일). K8s client 의존 없음 (persist 는 pkg/apply.Service).
//
// # 설계 원칙
//
//   - CR 타입 비의존: Ports/Selector/Labels 는 호출자가 조립해 전달.
//   - Headless 와 Type 은 상호배타: Headless=true 면 Type 무시(ClusterIP "None").
//   - 빈 Type → ClusterIP 기본 (non-headless).
package service
