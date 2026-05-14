package monitoring

import (
	"fmt"
	"strings"
)

// OTelExporterSpec — OpenTelemetry exporter config (caller-supplied).
//
// 호출 시:
//
//	spec := monitoring.OTelExporterSpec{
//	    Endpoint: "opentelemetry-collector.observability:4317",
//	    Protocol: "grpc",
//	    Insecure: true,
//	    Headers:  map[string]string{"authorization": "Bearer ..."},
//	}
//	envs := spec.ContainerEnvs()
//
// envs 를 container.env 로 주입하면 OTel SDK 가 자동 pickup.
//
// Refs: ROADMAP.md 'OpenTelemetry exporter helper (선택, 호출자 요구 시)'
//       (P-B.7.3)
type OTelExporterSpec struct {
	// Endpoint — OTLP gRPC endpoint (host:port).
	Endpoint string
	// Protocol — "grpc" or "http/protobuf".
	Protocol string
	// Insecure — skip TLS (dev only).
	Insecure bool
	// Headers — auth headers (e.g., bearer token).
	Headers map[string]string
	// ServiceName — overrides default service.name.
	ServiceName string
}

// EnvVar — single env var (without depending on corev1 — caller maps to corev1.EnvVar).
type EnvVar struct {
	Name  string
	Value string
}

// ContainerEnvs — converts OTelExporterSpec to env vars per
// OpenTelemetry SDK auto-config spec.
//
// Spec reference: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
func (s OTelExporterSpec) ContainerEnvs() []EnvVar {
	if s.Endpoint == "" {
		return nil
	}
	out := []EnvVar{
		{Name: "OTEL_EXPORTER_OTLP_ENDPOINT", Value: s.Endpoint},
	}
	if s.Protocol != "" {
		out = append(out, EnvVar{Name: "OTEL_EXPORTER_OTLP_PROTOCOL", Value: s.Protocol})
	}
	if s.Insecure {
		out = append(out, EnvVar{Name: "OTEL_EXPORTER_OTLP_INSECURE", Value: "true"})
	}
	if s.ServiceName != "" {
		out = append(out, EnvVar{Name: "OTEL_SERVICE_NAME", Value: s.ServiceName})
	}
	if len(s.Headers) > 0 {
		parts := make([]string, 0, len(s.Headers))
		for k, v := range s.Headers {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
		out = append(out, EnvVar{Name: "OTEL_EXPORTER_OTLP_HEADERS", Value: strings.Join(parts, ",")})
	}
	return out
}

// Validate — sanity check on OTelExporterSpec.
func (s OTelExporterSpec) Validate() error {
	if s.Endpoint == "" {
		return nil // disabled — not an error
	}
	if s.Protocol != "" && s.Protocol != "grpc" && s.Protocol != "http/protobuf" {
		return fmt.Errorf("OTelExporterSpec.Protocol must be 'grpc' or 'http/protobuf', got %q", s.Protocol)
	}
	if !strings.Contains(s.Endpoint, ":") {
		return fmt.Errorf("OTelExporterSpec.Endpoint must include port (host:port), got %q", s.Endpoint)
	}
	return nil
}
