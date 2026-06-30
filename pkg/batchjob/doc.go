// SPDX-License-Identifier: MIT

// Package batchjob 은 downstream operator 공통의 batch/v1 Job 엔벨로프 빌더를
// 제공한다. 백업/복원/업로드 등 일회성 Job 의 BackoffLimit / TTLSecondsAfterFinished /
// RestartPolicy / 라벨 전파 / PodSecurityContext 컨벤션을 단일 진실원으로 고정한다.
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 valkey-operator (`resources/{backup,upload,download}_job.go`) 와
// mongodb-operator (`resources/builder.go` BuildBackupJob/BuildRestoreJob) 의
// Job 엔벨로프 평행 재구현을 추출하여 신규 도입되었다. downstream 동시 적용
// 회귀 통과 후 Stable 격상 예정 (docs/ROADMAP.md §API Stability Tier).
//
// # 원본 두 판 차이 + canonical 채택
//
//   - 공통: BackoffLimit + TTLSecondsAfterFinished(=86400, 24h 자동정리) +
//     RestartPolicyOnFailure + Job/Pod 양쪽 라벨 + restricted PodSecurityContext.
//   - 차이: 컨테이너 명령(valkey-cli --rdb / mongodump|mongorestore), 볼륨(backup
//     PVC / TLS secret), 이미지. → 컨테이너·볼륨·SC 는 호출자가 조립해 전달하고,
//     commons 는 *Job/PodTemplate 엔벨로프* 만 책임진다.
//
// # 의존성 정책
//
// k8s.io/api/batch/v1 + core/v1 + apimachinery 만 의존하는 *순수 데이터 변환* 패키지.
//
// # 설계 원칙
//
//   - RestartPolicy 빈 값 → OnFailure (일회성 Job 표준).
//   - Labels 는 Job ObjectMeta 와 Pod template 양쪽에 전파 (selector/관측 일관).
//   - 컨테이너/볼륨/PodSecurityContext 는 호출자 도메인 책임.
package batchjob
