// SPDX-License-Identifier: MIT

package monitoring

import (
	"fmt"
	"strings"
)

// OTelExporterSpec — OpenTelemetry exporter config (caller-supplied).
type OTelExporterSpec struct {
	Endpoint    string
	Protocol    string
	Insecure    bool
	Headers     map[string]string
	ServiceName string
}

// EnvVar — single env var.
type EnvVar struct {
	Name  string
	Value string
}

// ContainerEnvs converts spec to OpenTelemetry SDK auto-config env vars.
func (s OTelExporterSpec) ContainerEnvs() []EnvVar {
	if s.Endpoint == "" {
		return nil
	}
	out := []EnvVar{{Name: "OTEL_EXPORTER_OTLP_ENDPOINT", Value: s.Endpoint}}
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

// Validate sanity-checks the spec.
func (s OTelExporterSpec) Validate() error {
	if s.Endpoint == "" {
		return nil
	}
	if s.Protocol != "" && s.Protocol != "grpc" && s.Protocol != "http/protobuf" {
		return fmt.Errorf("OTelExporterSpec.Protocol must be 'grpc' or 'http/protobuf', got %q", s.Protocol)
	}
	if !strings.Contains(s.Endpoint, ":") {
		return fmt.Errorf("OTelExporterSpec.Endpoint must include port, got %q", s.Endpoint)
	}
	return nil
}
