// SPDX-License-Identifier: MIT

package reconcile

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Statusable 는 controller 가 CR 의 status conditions / phase 를 추상화하기 위한
// interface. 각 CR type 이 api 패키지에서 구현한다 — downstream operator 양 판
// (mongodb helpers.go L96-100 ↔ valkey helpers.go L79-83) byte-동일 검증 후 승격.
type Statusable interface {
	client.Object

	// GetConditions 는 status.conditions slice 의 포인터를 반환한다 —
	// meta.SetStatusCondition 이 in-place 로 갱신한다.
	GetConditions() *[]metav1.Condition

	// SetPhase 는 status.phase 를 설정한다 (예: "Failed").
	SetPhase(phase string)
}

// EventRecorder 는 k8s.io/client-go/tools/events.EventRecorder 와 구조 호환되는
// 최소 로컬 interface 다. 호출자가 controller-runtime Manager 에서 받은 신식
// recorder (mgr.GetEventRecorderFor 가 아닌 events.EventRecorder 구현체) 를 그대로
// 전달하면 컴파일 타임에 자동 만족 — commons 가 client-go tools 패키지를 직접
// 의존하지 않기 위한 격리 (doc.go §의존성 정책).
type EventRecorder interface {
	Eventf(regarding runtime.Object, related runtime.Object,
		eventtype, reason, action, note string, args ...any)
}
