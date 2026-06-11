// SPDX-License-Identifier: MIT

// Package reconcile 는 downstream operator controller 의 공통 reconcile
// orchestration 헬퍼를 제공한다 — status 추상화 interface (Statusable), 에러
// condition 표준 적용 (ApplyErrorCondition), finalizer cleanup 흐름
// (HandleFinalizerCleanup), Secret 멱등 생성 (SecretIfNotExists).
//
// # API Stability Tier
//
// Stability: Beta.
//
// 본 패키지는 downstream operator (mongodb / valkey) 의
// internal/controller/helpers.go + status_update.go 중복 코드 (~420 LOC 영역) 를
// 추출하여 신규 도입되었다. downstream consumer 동시 적용 회귀 통과 후 Stable Tier
// 격상 예정 (자세한 격상 조건은 docs/ROADMAP.md §API Stability Tier).
//
// # 의존성 정책
//
// pkg/pvc 와 동일한 예외 패턴 — 본 패키지는 reconcile 루프의 K8s API 호출
// (status Update / finalizer Update / Secret Get+Create) 을 직접 수행하므로
// controller-runtime/pkg/client 를 의존한다.
//
//   - controller-runtime client.Client 는 downstream operator 가 이미 보유한 표준
//     의존이므로 신규 의존 부담은 0.
//   - Event 발행은 k8s.io/client-go tools/events.EventRecorder 와 구조 호환되는
//     로컬 EventRecorder interface 로 수신 — commons 가 client-go tools 패키지를
//     직접 의존하지 않고, 호출자의 기존 recorder 가 그대로 만족한다.
//   - metric 은 WithMetricHook 옵션 주입 — prometheus 의존을 commons 에 전파하지
//     않는다 (valkey MetricReconcileErrors.WithLabelValues(...).Inc 인라인의 일반화).
//   - 패키지명이 controller-runtime 의 pkg/reconcile 와 동일하므로 양쪽을 함께
//     import 하는 호출자는 alias 권장 (예: commonsreconcile).
//
// # 설계 원칙
//
//   - condition 적용은 k8s.io/apimachinery meta.SetStatusCondition 위임 —
//     LastTransitionTime 은 *Status 변경 시만* 갱신 (K8s convention, mongodb
//     ADR-0013 canonical).
//   - status 영속화는 pkg/status.UpdateWithRetry — conflict 시 refetch 후 mutate
//     클로저 재적용으로 호출자 status 변경의 silent-loss 차단.
//   - Secret 생성은 Get NotFound → Create 의 IsAlreadyExists race-tolerant
//     (병렬 reconcile 안전).
//   - finalizer in-memory 조작은 pkg/finalizer 재사용 — 본 패키지는 그 위의
//     orchestration (cleanup 콜백 + client.Update) 만 담당.
//
// # 원본 두 판 (mongodb / valkey) 차이 + canonical 채택 기록
//
// 원본 두 판이 갈리는 지점은 mongodb 판을 canonical 로 채택했다 (각 차이의 근거):
//
//   - condition 적용: mongodb = meta.SetStatusCondition 위임 (ADR-0013) / valkey =
//     filter+append+LastTransitionTime=Now (매 reconcile false transition 버그) →
//     mongodb 채택. valkey 이관 시 condition LastTransitionTime 거동이 변경되므로
//     condition_age 류 메트릭 관측 변화 주의 (별도 MR + e2e 검증 권장).
//   - status 갱신: mongodb = conflict 시 refetch + mutate 재적용 / valkey = refetch
//     없는 naive retry (stale ResourceVersion 재시도 — 사실상 무효 결함) → mongodb
//     채택 (pkg/status.UpdateWithRetry 위임).
//   - Secret 생성: mongodb = Create 의 IsAlreadyExists race guard (iteration 41) /
//     valkey = guard 부재 (race 시 에러 전파) → mongodb 채택.
//   - 에러 condition metric: valkey = prometheus 카운터 인라인 Inc / mongodb = 없음
//     → WithMetricHook(func(namespace, name, component string)) 옵션 주입으로 일반화.
//   - condition Reason: mongodb = api 패키지 상수 (ReasonReconcileFailed) / valkey =
//     리터럴 — 양쪽 동일 값 "ReconcileFailed". 본 패키지 상수 + WithReason /
//     WithConditionType 옵션으로 API 패키지 결합 해소.
//   - 이벤트/condition 메시지 구분자: mongodb = "Failed to reconcile <c>- <err>" /
//     valkey = "Failed to reconcile <c>: <err>" → ": " 로 표준화 (cosmetic only —
//     meta.SetStatusCondition 은 Message 변경만으로 LastTransitionTime 을 갱신하지
//     않으므로 transition 영향 0).
//   - requeue 간격: 양 repo intervals.go requeueSteady = 30s 동일 실측 →
//     DefaultRequeueAfter 30s + WithRequeueAfter 옵션.
//   - scope 제외 (도메인 정책 잔류): valkey PausedAnnotation / isPaused (ADR-0015) +
//     CondType* 도메인 상수, mongodb clearReconcileErrorCondition (단일 repo 사용 —
//     2nd consumer 등장 시 승격 재검토).
package reconcile
