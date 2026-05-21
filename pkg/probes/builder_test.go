// SPDX-License-Identifier: Apache-2.0

package probes_test

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/keiailab/operator-commons/pkg/probes"
)

func TestNew_appliesKubeletDefaults(t *testing.T) {
	t.Parallel()

	got := probes.New().TCP(8080).Build()

	if got.PeriodSeconds != 10 {
		t.Errorf("PeriodSeconds default: want=10 got=%d", got.PeriodSeconds)
	}
	if got.TimeoutSeconds != 1 {
		t.Errorf("TimeoutSeconds default: want=1 got=%d", got.TimeoutSeconds)
	}
	if got.SuccessThreshold != 1 {
		t.Errorf("SuccessThreshold default: want=1 got=%d", got.SuccessThreshold)
	}
	if got.FailureThreshold != 3 {
		t.Errorf("FailureThreshold default: want=3 got=%d", got.FailureThreshold)
	}
	if got.InitialDelaySeconds != 0 {
		t.Errorf("InitialDelaySeconds default: want=0 got=%d", got.InitialDelaySeconds)
	}
}

func TestHTTP_setsHTTPGetActionWithHTTPScheme(t *testing.T) {
	t.Parallel()

	got := probes.New().HTTP("/readyz", 8080).Build()

	if got.HTTPGet == nil {
		t.Fatal("HTTPGet handler nil")
	}
	if got.HTTPGet.Path != "/readyz" {
		t.Errorf("Path: want=/readyz got=%q", got.HTTPGet.Path)
	}
	if got.HTTPGet.Port != intstr.FromInt32(8080) {
		t.Errorf("Port: want=8080 got=%v", got.HTTPGet.Port)
	}
	if got.HTTPGet.Scheme != "" {
		t.Errorf("HTTP Scheme: want=empty (defaults to HTTP) got=%q", got.HTTPGet.Scheme)
	}
}

func TestHTTPS_setsHTTPSScheme(t *testing.T) {
	t.Parallel()

	got := probes.New().HTTPS("/healthz", 8443).Build()

	if got.HTTPGet == nil {
		t.Fatal("HTTPGet handler nil")
	}
	if got.HTTPGet.Scheme != corev1.URISchemeHTTPS {
		t.Errorf("Scheme: want=HTTPS got=%q", got.HTTPGet.Scheme)
	}
}

func TestTCP_setsTCPSocketAction(t *testing.T) {
	t.Parallel()

	got := probes.New().TCP(6379).Build()

	if got.TCPSocket == nil {
		t.Fatal("TCPSocket handler nil")
	}
	if got.TCPSocket.Port != intstr.FromInt32(6379) {
		t.Errorf("Port: want=6379 got=%v", got.TCPSocket.Port)
	}
	if got.HTTPGet != nil || got.Exec != nil {
		t.Error("Only TCPSocket should be set")
	}
}

func TestExec_setsExecActionCommand(t *testing.T) {
	t.Parallel()

	cmd := []string{"valkey-cli", "-a", "$(VALKEY_PASSWORD)", "ping"}
	got := probes.New().Exec(cmd...).Build()

	if got.Exec == nil {
		t.Fatal("Exec handler nil")
	}
	if len(got.Exec.Command) != len(cmd) {
		t.Fatalf("Command length: want=%d got=%d", len(cmd), len(got.Exec.Command))
	}
	for i, want := range cmd {
		if got.Exec.Command[i] != want {
			t.Errorf("Command[%d]: want=%q got=%q", i, want, got.Exec.Command[i])
		}
	}
}

func TestFluentChain_overridesAllFields(t *testing.T) {
	t.Parallel()

	got := probes.New().
		HTTP("/livez", 9090).
		InitialDelay(15 * time.Second).
		Period(20 * time.Second).
		Timeout(5 * time.Second).
		SuccessThreshold(1).
		FailureThreshold(6).
		Build()

	if got.InitialDelaySeconds != 15 {
		t.Errorf("InitialDelaySeconds: want=15 got=%d", got.InitialDelaySeconds)
	}
	if got.PeriodSeconds != 20 {
		t.Errorf("PeriodSeconds: want=20 got=%d", got.PeriodSeconds)
	}
	if got.TimeoutSeconds != 5 {
		t.Errorf("TimeoutSeconds: want=5 got=%d", got.TimeoutSeconds)
	}
	if got.FailureThreshold != 6 {
		t.Errorf("FailureThreshold: want=6 got=%d", got.FailureThreshold)
	}
}

func TestInitialDelay_clampsNegativeToZero(t *testing.T) {
	t.Parallel()

	got := probes.New().TCP(1234).InitialDelay(-5 * time.Second).Build()

	if got.InitialDelaySeconds != 0 {
		t.Errorf("Negative duration clamp: want=0 got=%d", got.InitialDelaySeconds)
	}
}

func TestHandlerOverride_lastWins(t *testing.T) {
	t.Parallel()

	got := probes.New().TCP(1234).HTTP("/readyz", 8080).Build()

	if got.HTTPGet == nil {
		t.Error("HTTPGet should be set after override")
	}
	if got.TCPSocket != nil {
		t.Error("TCPSocket should be cleared by HTTP override")
	}
}

func TestBuild_panicsWhenNoHandler(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Build() without handler must panic")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic value type: want=string got=%T", r)
		}
		if msg == "" {
			t.Error("panic message empty")
		}
	}()

	_ = probes.New().Build()
}

func TestPostgresOperatorReadinessPattern(t *testing.T) {
	t.Parallel()

	// Replaces: downstream operator builders source
	got := probes.New().
		HTTP("/readyz", 8080).
		InitialDelay(5 * time.Second).
		Period(10 * time.Second).
		Build()

	if got.HTTPGet.Path != "/readyz" || got.InitialDelaySeconds != 5 {
		t.Errorf("postgres readiness pattern mismatch: %+v", got)
	}
}

func TestMongoDBOperatorLivenessExecPattern(t *testing.T) {
	t.Parallel()

	// Replaces: downstream operator builders source
	got := probes.New().
		Exec("mongosh", "--quiet", "--eval", "db.adminCommand('ping')").
		InitialDelay(30 * time.Second).
		Period(10 * time.Second).
		Build()

	if got.Exec.Command[0] != "mongosh" || got.InitialDelaySeconds != 30 {
		t.Errorf("mongodb liveness pattern mismatch: %+v", got)
	}
}

func TestValkeyOperatorReadinessExecPattern(t *testing.T) {
	t.Parallel()

	// Replaces: downstream operator builders source
	got := probes.New().
		Exec("valkey-cli", "-a", "$(VALKEY_PASSWORD)", "ping").
		InitialDelay(5 * time.Second).
		Period(10 * time.Second).
		FailureThreshold(3).
		Build()

	if got.Exec.Command[0] != "valkey-cli" || got.FailureThreshold != 3 {
		t.Errorf("valkey readiness pattern mismatch: %+v", got)
	}
}
