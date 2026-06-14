// SPDX-License-Identifier: MIT

package reconcile

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/keiailab/keiailab-commons/pkg/events"
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

// EventRecorder 는 pkg/events.Recorder 의 별칭이다 — commons 라이브러리의 단일
// Event 기록 interface (k8s.io/client-go/tools/events.EventRecorder 와 구조 호환).
// 호출자가 controller-runtime Manager 에서 받은 신식 events.EventRecorder 구현체를
// 그대로 전달하면 컴파일 타임에 자동 만족한다. pkg/events 가 canonical interface +
// Reason 상수를 소유(단일 SSOT)하고, 기존 reconcile.EventRecorder 참조는 별칭으로
// 그대로 컴파일된다 (ApplyErrorCondition 시그니처 불변 = downstream operator 무영향).
type EventRecorder = events.Recorder
