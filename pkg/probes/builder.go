// Package probes provides a fluent builder for corev1.Probe.
//
// Replaces 9 duplicated probe construction sites across keiailab
// postgres-operator (2 HTTP sites), mongodb-operator (2 Exec sites),
// and valkey-operator (2 Exec sites) plus 3 cross-cutting patterns.
package probes

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// kubelet default values for Probe fields when not explicitly set.
// Reference: k8s.io/api/core/v1.Probe field comments + kubelet defaults.
const (
	defaultPeriodSeconds    int32 = 10
	defaultTimeoutSeconds   int32 = 1
	defaultSuccessThreshold int32 = 1
	defaultFailureThreshold int32 = 3
)

// Builder constructs corev1.Probe with sensible kubelet defaults and a fluent API.
//
// Zero value is not directly usable. Always start with New().
//
// Example (HTTP readiness):
//
//	probe := probes.New().
//	    HTTP("/readyz", 8080).
//	    InitialDelay(5*time.Second).
//	    Period(10*time.Second).
//	    Build()
//
// Example (Exec liveness with auth-aware command):
//
//	probe := probes.New().
//	    Exec("valkey-cli", "-a", "$(VALKEY_PASSWORD)", "ping").
//	    InitialDelay(20*time.Second).
//	    FailureThreshold(5).
//	    Build()
type Builder struct {
	initialDelay     int32
	period           int32
	timeout          int32
	successThreshold int32
	failureThreshold int32
	handler          corev1.ProbeHandler
	handlerSet       bool
}

// New returns a Builder seeded with kubelet default values.
// Caller must set exactly one handler (HTTP / HTTPS / TCP / Exec) before Build.
func New() *Builder {
	return &Builder{
		period:           defaultPeriodSeconds,
		timeout:          defaultTimeoutSeconds,
		successThreshold: defaultSuccessThreshold,
		failureThreshold: defaultFailureThreshold,
	}
}

// HTTP sets HTTPGetAction handler with default Scheme=HTTP.
// path is the URL path; port is the container TCP port.
func (b *Builder) HTTP(path string, port int32) *Builder {
	b.handler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: path,
			Port: intstr.FromInt32(port),
		},
	}
	b.handlerSet = true
	return b
}

// HTTPS sets HTTPGetAction handler with Scheme=HTTPS.
func (b *Builder) HTTPS(path string, port int32) *Builder {
	b.handler = corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{
			Path:   path,
			Port:   intstr.FromInt32(port),
			Scheme: corev1.URISchemeHTTPS,
		},
	}
	b.handlerSet = true
	return b
}

// TCP sets TCPSocketAction handler.
func (b *Builder) TCP(port int32) *Builder {
	b.handler = corev1.ProbeHandler{
		TCPSocket: &corev1.TCPSocketAction{
			Port: intstr.FromInt32(port),
		},
	}
	b.handlerSet = true
	return b
}

// Exec sets ExecAction handler with the given command slice.
// At least one element required (operator-specific health command).
func (b *Builder) Exec(cmd ...string) *Builder {
	b.handler = corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: cmd,
		},
	}
	b.handlerSet = true
	return b
}

// InitialDelay sets InitialDelaySeconds. Negative durations are treated as 0.
func (b *Builder) InitialDelay(d time.Duration) *Builder {
	b.initialDelay = secondsOf(d)
	return b
}

// Period sets PeriodSeconds. Default 10s if not set.
func (b *Builder) Period(d time.Duration) *Builder {
	b.period = secondsOf(d)
	return b
}

// Timeout sets TimeoutSeconds. Default 1s if not set.
func (b *Builder) Timeout(d time.Duration) *Builder {
	b.timeout = secondsOf(d)
	return b
}

// SuccessThreshold sets SuccessThreshold. Default 1.
// For liveness/startup probes, Kubernetes requires this to be 1.
func (b *Builder) SuccessThreshold(n int32) *Builder {
	b.successThreshold = n
	return b
}

// FailureThreshold sets FailureThreshold. Default 3.
func (b *Builder) FailureThreshold(n int32) *Builder {
	b.failureThreshold = n
	return b
}

// Build returns the configured corev1.Probe.
// Panics if no handler was set (HTTP / HTTPS / TCP / Exec must be called once).
func (b *Builder) Build() *corev1.Probe {
	if !b.handlerSet {
		panic("probes.Builder: handler not set (call HTTP/HTTPS/TCP/Exec before Build)")
	}
	return &corev1.Probe{
		ProbeHandler:        b.handler,
		InitialDelaySeconds: b.initialDelay,
		PeriodSeconds:       b.period,
		TimeoutSeconds:      b.timeout,
		SuccessThreshold:    b.successThreshold,
		FailureThreshold:    b.failureThreshold,
	}
}

// secondsOf converts a Duration to int32 seconds, clamping negatives to 0.
func secondsOf(d time.Duration) int32 {
	if d < 0 {
		return 0
	}
	return int32(d.Seconds())
}
