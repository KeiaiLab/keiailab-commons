// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/keiailab/keiailab-commons/pkg/status"
)

// 표준 기본값 — downstream operator 양 판 실측 동일 값 (intervals.go requeueSteady 30s
// + condition type/reason 리터럴). WithRequeueAfter / WithConditionType / WithReason
// 옵션으로 호출자가 override 가능.
const (
	// DefaultRequeueAfter — ApplyErrorCondition 의 기본 재큐 간격 (양 repo
	// requeueSteady = 30s 동일 실측).
	DefaultRequeueAfter = 30 * time.Second

	// ConditionTypeReconcileError — 에러 condition 의 기본 Type. 이벤트 reason /
	// action 으로도 동일 값 사용 (원본 양 판 정합).
	ConditionTypeReconcileError = "ReconcileError"

	// ReasonReconcileFailed — 에러 condition 의 기본 Reason (양 repo 동일 값 —
	// mongodb api 상수 / valkey 리터럴).
	ReasonReconcileFailed = "ReconcileFailed"

	// PhaseFailed — ApplyErrorCondition 이 설정하는 status.phase 값.
	PhaseFailed = "Failed"
)

// options — ApplyErrorCondition 의 functional options 누산자.
type options struct {
	requeueAfter  time.Duration
	conditionType string
	reason        string
	metricHook    func(namespace, name, component string)
}

// Option — ApplyErrorCondition 의 functional option.
type Option func(*options)

// WithMetricHook 는 에러 카운터 증분 등 metric 부수효과를 주입한다 — prometheus
// 의존을 commons 에 전파하지 않는 일반화 (valkey
// MetricReconcileErrors.WithLabelValues(ns, name, component).Inc() 인라인 대체).
//
// 사용 예 (downstream operator init):
//
//	reconcile.WithMetricHook(func(ns, name, component string) {
//	    MetricReconcileErrors.WithLabelValues(ns, name, component).Inc()
//	})
func WithMetricHook(hook func(namespace, name, component string)) Option {
	return func(o *options) { o.metricHook = hook }
}

// WithRequeueAfter 는 반환 ctrl.Result 의 RequeueAfter 를 override 한다
// (기본 DefaultRequeueAfter 30s).
func WithRequeueAfter(d time.Duration) Option {
	return func(o *options) { o.requeueAfter = d }
}

// WithReason 는 condition Reason 을 override 한다 (기본 ReasonReconcileFailed).
// pkg/status Reason 카탈로그 스타일 상수 권장.
func WithReason(reason string) Option {
	return func(o *options) { o.reason = reason }
}

// WithConditionType 는 condition Type (+ 이벤트 reason/action) 을 override 한다
// (기본 ConditionTypeReconcileError).
func WithConditionType(conditionType string) Option {
	return func(o *options) { o.conditionType = conditionType }
}

// ApplyErrorCondition 는 reconcile 에러 처리의 표준 패턴 —
//
//  1. logger.Error 출력
//  2. EventRecorder Warning 이벤트 발행 (rec 이 nil 이면 skip)
//  3. metric hook 호출 (WithMetricHook 주입 시)
//  4. Status.Phase = "Failed" + 에러 condition 적용 (meta.SetStatusCondition —
//     LastTransitionTime 은 Status 변경 시만 갱신, K8s convention / ADR-0013)
//  5. status.UpdateWithRetry 로 영속화 — status mutation 을 클로저로 전달해
//     conflict 재시도 시 refetch 후 재적용 (silent-loss 차단)
//  6. RequeueAfter + 원본 reconcileErr 반환
//
// status 갱신 실패는 log 만 남기고 원본 reconcileErr 를 반환한다 (원본 에러가
// requeue 를 이미 보장 — 원본 양 판 동일 거동).
func ApplyErrorCondition(
	ctx context.Context,
	c client.Client,
	obj Statusable,
	component string,
	reconcileErr error,
	rec EventRecorder,
	opts ...Option,
) (ctrl.Result, error) {
	o := options{
		requeueAfter:  DefaultRequeueAfter,
		conditionType: ConditionTypeReconcileError,
		reason:        ReasonReconcileFailed,
	}
	for _, opt := range opts {
		opt(&o)
	}

	logger := log.FromContext(ctx)
	logger.Error(reconcileErr, "Failed to reconcile component", "component", component)
	if o.metricHook != nil {
		o.metricHook(obj.GetNamespace(), obj.GetName(), component)
	}
	if rec != nil {
		rec.Eventf(obj, nil, corev1.EventTypeWarning, o.conditionType, o.conditionType,
			"Failed to reconcile %s: %v", component, reconcileErr)
	}

	// status mutation 클로저 — conflict 재시도 시 refetch 후 재적용되도록
	// UpdateWithRetry 에 동일 클로저를 전달한다 (mongodb PR-2 canonical).
	applyStatus := func() {
		obj.SetPhase(PhaseFailed)
		meta.SetStatusCondition(obj.GetConditions(), metav1.Condition{
			Type:    o.conditionType,
			Status:  metav1.ConditionTrue,
			Reason:  o.reason,
			Message: fmt.Sprintf("Failed to reconcile %s: %v", component, reconcileErr),
		})
	}
	applyStatus()

	if statusErr := status.UpdateWithRetry(ctx, c, obj, applyStatus); statusErr != nil {
		logger.Error(statusErr, "Failed to update status")
	}
	return ctrl.Result{RequeueAfter: o.requeueAfter}, reconcileErr
}
