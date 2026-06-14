// SPDX-License-Identifier: MIT

package events_test

import (
	"errors"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/keiailab/keiailab-commons/pkg/events"
)

// mockRecorder implements events.Recorder (the modern events.EventRecorder
// shape) for unit testing.
type mockRecorder struct {
	events []capturedEvent
}

type capturedEvent struct {
	eventType string
	reason    string
	action    string
	message   string
}

func (m *mockRecorder) Eventf(_ runtime.Object, _ runtime.Object, eventtype, reason, action, note string, args ...any) {
	m.events = append(m.events, capturedEvent{eventtype, reason, action, fmt.Sprintf(note, args...)})
}

func newFakeObject() runtime.Object {
	return &corev1.ConfigMap{}
}

func TestEmit_normalEvent(t *testing.T) {
	t.Parallel()

	rec := &mockRecorder{}
	events.Emit(rec, newFakeObject(), events.ReasonCreated, "cluster created")

	if len(rec.events) != 1 {
		t.Fatalf("event count: want=1 got=%d", len(rec.events))
	}
	e := rec.events[0]
	if e.eventType != events.TypeNormal {
		t.Errorf("eventType: want=%q got=%q", events.TypeNormal, e.eventType)
	}
	if e.reason != events.ReasonCreated {
		t.Errorf("reason: want=%q got=%q", events.ReasonCreated, e.reason)
	}
	// Emit maps action=reason (suite's established degenerate usage).
	if e.action != events.ReasonCreated {
		t.Errorf("action: want=%q (==reason) got=%q", events.ReasonCreated, e.action)
	}
	if e.message != "cluster created" {
		t.Errorf("message: want=%q got=%q", "cluster created", e.message)
	}
}

func TestEmit_messageWithPercentNotFormatInjected(t *testing.T) {
	t.Parallel()

	// Emit must treat message as data, not a format string — a literal '%'
	// in the message must survive verbatim (mapped via "%s").
	rec := &mockRecorder{}
	events.Emit(rec, newFakeObject(), events.ReasonProvisioning, "disk 95% full")

	if rec.events[0].message != "disk 95% full" {
		t.Errorf("message with percent: got=%q want=%q", rec.events[0].message, "disk 95% full")
	}
}

func TestEmit_nilRecorderNoop(t *testing.T) {
	t.Parallel()

	// Must not panic.
	events.Emit(nil, newFakeObject(), events.ReasonCreated, "ignored")
}

func TestEmitf_formatting(t *testing.T) {
	t.Parallel()

	rec := &mockRecorder{}
	events.Emitf(rec, newFakeObject(), events.ReasonReconciled, "reconciled %d shards in %s", 3, "default")

	if rec.events[0].message != "reconciled 3 shards in default" {
		t.Errorf("formatted message: got=%q", rec.events[0].message)
	}
}

func TestEmitf_nilRecorderNoop(t *testing.T) {
	t.Parallel()

	events.Emitf(nil, newFakeObject(), events.ReasonReconciled, "ignored %d", 1)
}

func TestEmitWarning_includesErrorMessage(t *testing.T) {
	t.Parallel()

	rec := &mockRecorder{}
	err := errors.New("rbd binding refused")
	events.EmitWarning(rec, newFakeObject(), events.ReasonFailed, err)

	if len(rec.events) != 1 {
		t.Fatalf("event count: want=1 got=%d", len(rec.events))
	}
	if rec.events[0].eventType != events.TypeWarning {
		t.Errorf("eventType: want=Warning got=%q", rec.events[0].eventType)
	}
	if rec.events[0].message != "rbd binding refused" {
		t.Errorf("message: want=err.Error() got=%q", rec.events[0].message)
	}
}

func TestEmitWarning_nilErrorNoop(t *testing.T) {
	t.Parallel()

	rec := &mockRecorder{}
	events.EmitWarning(rec, newFakeObject(), events.ReasonFailed, nil)

	if len(rec.events) != 0 {
		t.Errorf("expected no event when err is nil, got=%d events", len(rec.events))
	}
}

func TestEmitWarning_nilRecorderNoop(t *testing.T) {
	t.Parallel()

	// Must not panic.
	events.EmitWarning(nil, newFakeObject(), events.ReasonFailed, errors.New("ignored"))
}

func TestEmitWarningf_formatting(t *testing.T) {
	t.Parallel()

	rec := &mockRecorder{}
	events.EmitWarningf(rec, newFakeObject(), events.ReasonDegraded, "replicas %d below threshold %d", 1, 3)

	if rec.events[0].eventType != events.TypeWarning {
		t.Errorf("eventType: want=Warning got=%q", rec.events[0].eventType)
	}
	if rec.events[0].message != "replicas 1 below threshold 3" {
		t.Errorf("formatted: got=%q", rec.events[0].message)
	}
}

func TestEmitWarningf_nilRecorderNoop(t *testing.T) {
	t.Parallel()

	events.EmitWarningf(nil, newFakeObject(), events.ReasonDegraded, "ignored %d", 1)
}

func TestWrappedError_withError(t *testing.T) {
	t.Parallel()

	got := events.WrappedError(events.ReasonReconcileErr, errors.New("db unreachable"))
	want := "ReconcileError: db unreachable"
	if got != want {
		t.Errorf("WrappedError: got=%q want=%q", got, want)
	}
}

func TestWrappedError_nilErrorReturnsReasonOnly(t *testing.T) {
	t.Parallel()

	got := events.WrappedError(events.ReasonReady, nil)
	if got != events.ReasonReady {
		t.Errorf("WrappedError(nil): got=%q want=%q", got, events.ReasonReady)
	}
}

func TestReasonConstants_uniqueValues(t *testing.T) {
	t.Parallel()

	// Catch accidental duplicate values when adding new Reason constants.
	reasons := map[string]string{
		events.ReasonCreated:      "ReasonCreated",
		events.ReasonUpdated:      "ReasonUpdated",
		events.ReasonDeleted:      "ReasonDeleted",
		events.ReasonReconciled:   "ReasonReconciled",
		events.ReasonReconcileErr: "ReasonReconcileErr",
		events.ReasonProvisioning: "ReasonProvisioning",
		events.ReasonReady:        "ReasonReady",
		events.ReasonDegraded:     "ReasonDegraded",
		events.ReasonFailed:       "ReasonFailed",
	}

	if len(reasons) != 9 {
		t.Errorf("Reason constant uniqueness: 9 distinct values expected, got=%d", len(reasons))
	}
}

func TestRecorderInterface_acceptsAnyImpl(t *testing.T) {
	t.Parallel()

	// Compile-time check: mockRecorder satisfies events.Recorder (modern API).
	var _ events.Recorder = (*mockRecorder)(nil)
}
