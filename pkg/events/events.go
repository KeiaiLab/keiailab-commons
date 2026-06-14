// SPDX-License-Identifier: MIT

// Package events provides Kubernetes Event recording helpers and standard
// reconciler event reasons shared across keiailab operators.
//
// The Recorder interface matches k8s.io/client-go/tools/events.EventRecorder
// (the modern events API adopted suite-wide in RFC-0023 Phase 2), so callers
// may pass the recorder they obtain from the controller-runtime manager
// without an adapter. The package itself does not depend on client-go.
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
// It matches k8s.io/client-go/tools/events.EventRecorder (the modern events
// API): callers pass the recorder obtained from the controller-runtime manager
// without an adapter. commons does not import client-go tools packages — the
// interface is declared structurally.
//
// The Emit* helpers below map their ergonomic (reason, message) arguments onto
// Eventf with related=nil and action=reason, matching the suite's established
// usage. Call sites needing a distinct action or a related object should call
// Eventf directly.
type Recorder interface {
	Eventf(regarding runtime.Object, related runtime.Object,
		eventtype, reason, action, note string, args ...any)
}

// Emit records a Normal event with the given reason and message.
//
// No-op when rec is nil — convenient for optional event recording paths
// (e.g. unit-tested reconcilers without a real EventRecorder).
func Emit(rec Recorder, obj runtime.Object, reason, message string) {
	if rec == nil {
		return
	}
	rec.Eventf(obj, nil, TypeNormal, reason, reason, "%s", message)
}

// Emitf is Emit with printf-style formatting.
//
// No-op when rec is nil.
func Emitf(rec Recorder, obj runtime.Object, reason, format string, args ...any) {
	if rec == nil {
		return
	}
	rec.Eventf(obj, nil, TypeNormal, reason, reason, format, args...)
}

// EmitWarning records a Warning event using err.Error() as the message.
//
// No-op when rec or err is nil — caller can use this in error-return paths
// without nil-checking err separately.
func EmitWarning(rec Recorder, obj runtime.Object, reason string, err error) {
	if rec == nil || err == nil {
		return
	}
	rec.Eventf(obj, nil, TypeWarning, reason, reason, "%s", err.Error())
}

// EmitWarningf is a Warning event with printf-style formatting.
//
// No-op when rec is nil.
func EmitWarningf(rec Recorder, obj runtime.Object, reason, format string, args ...any) {
	if rec == nil {
		return
	}
	rec.Eventf(obj, nil, TypeWarning, reason, reason, format, args...)
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
