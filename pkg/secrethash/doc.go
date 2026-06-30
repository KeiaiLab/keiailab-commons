// SPDX-License-Identifier: MIT

// Package secrethash 는 Secret data 의 *결정적* digest 를 계산하는 헬퍼를
// 제공한다. operator 가 cert / password 등 Secret 변경을 감지해 워크로드
// rolling update 를 트리거할 때 사용하는 annotation 값 산출이 주 용도다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator 의 평행 재구현
// (valkey-operator `controller/tls_cert_hash.go` + mongodb-operator
// `controller/password_rotation.go:hashSecretData`) 을 추출하여 신규
// 도입되었다. downstream consumer 동시 적용 회귀 통과 후 Stable Tier
// 격상 예정 (격상 조건은 docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
// valkey 와 mongo 의 secret-hash 구현이 달랐다:
//
//   - valkey `hashTLSSecret`: 고정 키(tls.crt/tls.key/ca.crt)를 *명시 순서* 로
//     누적, full hex. (결정적)
//   - mongo `hashSecretData`: `for k, v := range secret.Data` 로 map 을 *비결정적*
//     순회 후 hex 의 앞 16자만 truncate. Go map 순회 순서는 무작위이므로 동일한
//     Secret 이라도 reconcile 마다 다른 digest 가 나올 수 있었다 (→ 불필요한
//     rolling update 유발 위험 = 잠재 버그).
//
// canonical 채택: 키를 *정렬*(전체) 또는 *명시 순서*(부분)로만 누적 + full hex.
// 이로써 mongo 의 비결정성 버그를 구조적으로 제거하고 valkey 의 결정적 거동과
// 통일한다. 단 mongo 채택 시 digest *형식* 이 16자→64자(full hex)로 바뀌므로
// adopt 시점에 1회 rolling update 가 발생할 수 있다 (annotation 값 변경 = 예상된
// 1회성, 이후 안정).
//
// # 의존성 정책
//
// 본 패키지는 stdlib (crypto/sha256, encoding/hex, sort) 만 의존하는 *순수
// 데이터 변환* 패키지다 (labels / security 와 동일군). K8s client 의존 없음.
//
// # 설계 원칙
//
//   - 결정성: 키 미지정 시 정렬 순서, 지정 시 인자 순서로만 누적.
//   - 키 + 값 양쪽 누적 — 키 자체를 digest 에 포함해 빈 값 / 키 누락을 구분
//     (예: "ca.crt" 만 추가/삭제되어도 digest 가 변한다).
//   - 원본 secret 미노출: SHA256 의 hex digest 만 반환 — annotation 으로
//     K8s API 에 노출돼도 원본 cert / password 는 복원 불가.
package secrethash
