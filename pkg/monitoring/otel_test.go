package monitoring

import (
	"strings"
	"testing"
)

func TestOTelExporterSpec_ContainerEnvs(t *testing.T) {
	t.Run("empty endpoint = no envs", func(t *testing.T) {
		s := OTelExporterSpec{}
		if len(s.ContainerEnvs()) != 0 {
			t.Errorf("expected 0 envs, got %d", len(s.ContainerEnvs()))
		}
	})

	t.Run("endpoint only", func(t *testing.T) {
		s := OTelExporterSpec{Endpoint: "collector:4317"}
		envs := s.ContainerEnvs()
		if len(envs) != 1 || envs[0].Name != "OTEL_EXPORTER_OTLP_ENDPOINT" {
			t.Errorf("expected 1 env (endpoint), got %+v", envs)
		}
	})

	t.Run("full spec", func(t *testing.T) {
		s := OTelExporterSpec{
			Endpoint:    "collector:4317",
			Protocol:    "grpc",
			Insecure:    true,
			ServiceName: "valkey-operator",
			Headers:     map[string]string{"authorization": "Bearer xyz"},
		}
		envs := s.ContainerEnvs()
		if len(envs) != 5 {
			t.Errorf("expected 5 envs, got %d: %+v", len(envs), envs)
		}
		found := map[string]bool{}
		for _, e := range envs {
			found[e.Name] = true
		}
		want := []string{"OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_PROTOCOL",
			"OTEL_EXPORTER_OTLP_INSECURE", "OTEL_SERVICE_NAME", "OTEL_EXPORTER_OTLP_HEADERS"}
		for _, name := range want {
			if !found[name] {
				t.Errorf("missing env %q", name)
			}
		}
	})
}

func TestOTelExporterSpec_Validate(t *testing.T) {
	t.Run("empty = ok (disabled)", func(t *testing.T) {
		if err := (OTelExporterSpec{}).Validate(); err != nil {
			t.Errorf("empty Validate = %v, want nil", err)
		}
	})
	t.Run("missing port", func(t *testing.T) {
		err := OTelExporterSpec{Endpoint: "collector-no-port"}.Validate()
		if err == nil || !strings.Contains(err.Error(), "port") {
			t.Errorf("expected port error, got %v", err)
		}
	})
	t.Run("invalid protocol", func(t *testing.T) {
		err := OTelExporterSpec{Endpoint: "c:4317", Protocol: "weird"}.Validate()
		if err == nil {
			t.Errorf("expected protocol error")
		}
	})
	t.Run("valid grpc", func(t *testing.T) {
		err := OTelExporterSpec{Endpoint: "c:4317", Protocol: "grpc"}.Validate()
		if err != nil {
			t.Errorf("valid spec err = %v", err)
		}
	})
}
