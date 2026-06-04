// SPDX-License-Identifier: MIT

// Package events provides Kubernetes Event recording helpers and standard
// reconciler event reasons shared across keiailab operators.
//
// The Recorder interface is structurally compatible with
// k8s.io/client-go/tools/record.EventRecorder, so callers may pass their
// existing recorder without an adapter. The package itself does not depend
// on client-go.
package events

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

// Event type constants mirror corev1.EventType values.
// Defined locally to avoid pulling k8s.io/api into helper-only call sites.
const (
	TypeNormal  = "Normal"
	TypeWarning = "Warning"
)

// Standard reconciler event reason constants shared across keiailab operators.
//
// Use these in place of repo-specific string literals for cross-operator
// consistency in Event streams. Operator-specific reasons (e.g. database
// state transitions) should remain in operator packages.
const (
	ReasonCreated      = "Created"
	ReasonUpdated      = "Updated"
	ReasonDeleted      = "Deleted"
	ReasonReconciled   = "Reconciled"
	ReasonReconcileErr = "ReconcileError"
	ReasonProvisioning = "Provisioning"
	ReasonReady        = "Ready"
	ReasonDegraded     = "Degraded"
	ReasonFailed       = "Failed"
)

// Recorder is the minimal Event recording interface used by this package.
//
// Structurally compatible with k8s.io/client-go/tools/record.EventRecorder —
// callers passing their existing record.EventRecorder need no adapter.
type Recorder interface {
	Event(object runtime.Object, eventtype, reason, message string)
	Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...any)
}

// Emit records a Normal event with the given reason and message.
//
// No-op when rec is nil — convenient for optional event recording paths
// (e.g. unit-tested reconcilers without a real EventRecorder).
func Emit(rec Recorder, obj runtime.Object, reason, message string) {
	if rec == nil {
		return
	}
	rec.Event(obj, TypeNormal, reason, message)
}

// Emitf is Emit with printf-style formatting.
//
// No-op when rec is nil.
func Emitf(rec Recorder, obj runtime.Object, reason, format string, args ...any) {
	if rec == nil {
		return
	}
	rec.Eventf(obj, TypeNormal, reason, format, args...)
}

// EmitWarning records a Warning event using err.Error() as the message.
//
// No-op when rec or err is nil — caller can use this in error-return paths
// without nil-checking err separately.
func EmitWarning(rec Recorder, obj runtime.Object, reason string, err error) {
	if rec == nil || err == nil {
		return
	}
	rec.Event(obj, TypeWarning, reason, err.Error())
}

// EmitWarningf is a Warning event with printf-style formatting.
//
// No-op when rec is nil.
func EmitWarningf(rec Recorder, obj runtime.Object, reason, format string, args ...any) {
	if rec == nil {
		return
	}
	rec.Eventf(obj, TypeWarning, reason, format, args...)
}

// WrappedError formats a reason and error into a single status-line string
// suitable for Status.Conditions[].Message.
//
// Returns reason unchanged when err is nil.
func WrappedError(reason string, err error) string {
	if err == nil {
		return reason
	}
	return fmt.Sprintf("%s: %s", reason, err.Error())
}
