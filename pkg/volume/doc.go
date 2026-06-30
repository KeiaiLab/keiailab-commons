// SPDX-License-Identifier: MIT

// Package volume 은 downstream operator 공통의 Pod Volume / VolumeMount 빌더를
// 제공한다. 현재: cert-manager 발급 TLS Secret 의 read-only 마운트.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 valkey-operator (`resources/statefulset.go` tlsVolumes/tlsVolumeMounts)
// 와 mongodb-operator (`resources/builder.go` buildTLSServerVolume/buildTLSServerMount)
// 의 평행 재구현을 추출하여 신규 도입되었다. downstream 동시 적용 회귀 통과 후
// Stable 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
//   - 공통 불변식: TLS Secret 을 DefaultMode 0o400 + ReadOnly 로 마운트
//     (cert/key 파일 권한 검사 통과 + 최소권한).
//   - 차이: volume 이름(valkey "tls" / mongo "tls-server"), mountPath, 반환 형태
//     (valkey []slice + empty-guard / mongo 단일). → 이름·경로는 인자로, empty-guard
//     는 호출자 잔류. canonical = 0o400 + ReadOnly 불변식을 단일 SSOT 로 고정.
//
// # 의존성 정책
//
// k8s.io/api/core/v1 + k8s.io/utils/ptr 만 의존하는 *순수 데이터 변환* 패키지.
//
// # 설계 원칙
//
//   - 0o400 (owner read-only) + ReadOnly mount 는 항상 적용 (cert 보안 불변식).
//   - volume 이름 / mountPath 는 호출자 도메인 규칙으로 전달.
package volume
