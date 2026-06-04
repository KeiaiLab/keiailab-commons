// SPDX-License-Identifier: MIT

// Package status 는 downstream consumer keiailab operator 공통 Condition Type / Reason
// 카탈로그 + apimeta.SetStatusCondition 헬퍼를 제공한다.
//
// 본 패키지는 pkg/status 표준 의 spec 을 구현한다. downstream operator 의
// internal/controller/status.go 패턴 + Kubernetes deployment controller
// 관용구 (KEP-1623) 를 표준으로 채택.
//
// 외부 의존성: k8s.io/apimachinery 만. controller-runtime 미의존.
package status

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 표준 Condition Type — kubectl describe 출력 정합성을 위해 downstream consumer 동일.
//
// 운영자가 4 operator 모두에서 같은 키워드 (Ready / Progressing / Degraded /
// Available) 로 검색 가능.
const (
	// TypeReady — 운영 가능 상태. False 인 경우 Reason 으로 차단점 식별.
	TypeReady = "Ready"
	// TypeProgressing — reconcile 진행 중. 리소스 생성, shard scaling 등
	// *작업 중* 상태.
	TypeProgressing = "Progressing"
	// TypeDegraded — 부분 장애. 일부 기능은 작동, 일부는 장애.
	TypeDegraded = "Degraded"
	// TypeAvailable — 외부 호출 가능 (Service Endpoint 보유 + readiness probe
	// 통과). Ready 보다 약한 신호.
	TypeAvailable = "Available"
)

// 표준 Reason 카탈로그 — downstream operator conditions 패턴 + Kubernetes
// deployment controller 관용구. Reason 은 *기계 가독* 식별자이므로
// PascalCase + 의미 안정성 보장.
const (
	// ReasonReconciling — Reconcile 진입 직후 + 진행 중인 모든 작업.
	ReasonReconciling = "Reconciling"
	// ReasonAvailable — Ready=True 의 표준 Reason.
	ReasonAvailable = "Available"
	// ReasonNotApplicable — 본 condition 이 현 spec 에서 의미 없음 (예:
	// MonitoringSpec 미설정 시 Monitoring=NotApplicable).
	ReasonNotApplicable = "NotApplicable"
	// ReasonReconcileError — reconcile 중 에러 발생. Message 에 핵심 원인
	// 1줄 요약.
	ReasonReconcileError = "ReconcileError"
	// ReasonExternalDepBlocked — 외부 의존 (cert-manager, secret 등) 미준비
	// 로 차단.
	ReasonExternalDepBlocked = "ExternalDependencyBlocked"
	// ReasonValidationFailed — webhook 또는 controller 측 validation 실패.
	ReasonValidationFailed = "ValidationFailed"
)

// SetReady 는 Ready condition 을 갱신한다. status.LastTransitionTime 은
// apimeta 가 자동 처리. observedGeneration 은 호출자의 obj.Generation 을
// 그대로 전달 (controller-gen 가 status.observedGeneration 에 기록).
//
// 사용 예:
//
//	status.SetReady(&cluster.Status.Conditions, metav1.ConditionTrue,
//	    status.ReasonAvailable, "all members healthy", cluster.Generation)
func SetReady(
	conditions *[]metav1.Condition,
	s metav1.ConditionStatus,
	reason, message string,
	observedGeneration int64,
) {
	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:               TypeReady,
		Status:             s,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: observedGeneration,
	})
}

// SetProgressing 는 Progressing=True 를 기록한다. Reconcile 진입 시 호출.
// 작업 완료 시 SetReady 가 자동으로 Progressing 을 별도 처리하지 *않으므로*,
// 호출자가 SetProgressing(False) 또는 RemoveProgressing 으로 명시 종료.
func SetProgressing(conditions *[]metav1.Condition, reason, message string, observedGeneration int64) {
	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:               TypeProgressing,
		Status:             metav1.ConditionTrue,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: observedGeneration,
	})
}

// SetDegraded 기록. 일부 기능 장애 시.
func SetDegraded(conditions *[]metav1.Condition, reason, message string, observedGeneration int64) {
	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:               TypeDegraded,
		Status:             metav1.ConditionTrue,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: observedGeneration,
	})
}

// SetAvailable — Available condition 갱신. Endpoint 보유 + readiness probe
// 통과 후 호출. Ready 보다 약한 신호 (operator 의 reconcile 완료 vs 외부
// 호출 가능).
//
// downstream consumer 정합 패턴: SetReady=True 직후 SetAvailable=True 호출. reconcile
// 진입 시점에는 SetReady=False/Progressing 만 호출하고 Available 은 변경
// 안 함 (외부 client 의 기존 연결 유지).
func SetAvailable(
	conditions *[]metav1.Condition,
	s metav1.ConditionStatus,
	reason, message string,
	observedGeneration int64,
) {
	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:               TypeAvailable,
		Status:             s,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: observedGeneration,
	})
}

// SetReadyFalse — SetReady 의 False 케이스 슈가.
//
// reconcile 실패 가장 빈번한 호출 패턴 (`SetReady(_, ConditionFalse, …)`)
// 의 가독성 향상용 슈가. 동작은 SetReady(ConditionFalse) 와 동일.
func SetReadyFalse(
	conditions *[]metav1.Condition,
	reason, message string,
	observedGeneration int64,
) {
	SetReady(conditions, metav1.ConditionFalse, reason, message, observedGeneration)
}

// IsReady 는 Ready=True 인지 확인. False / 부재 시 false 반환.
func IsReady(conditions []metav1.Condition) bool {
	return meta.IsStatusConditionTrue(conditions, TypeReady)
}

// FindCondition 은 type 으로 Condition 을 조회. 부재 시 nil.
func FindCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	return meta.FindStatusCondition(conditions, conditionType)
}

// RemoveCondition 은 type 의 condition 을 제거. 없으면 no-op.
func RemoveCondition(conditions *[]metav1.Condition, conditionType string) {
	meta.RemoveStatusCondition(conditions, conditionType)
}
