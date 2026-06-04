// SPDX-License-Identifier: MIT

package monitoring

import (
	"strings"
	"testing"
)

func TestOTelExporterSpec_ContainerEnvs(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if len((OTelExporterSpec{}).ContainerEnvs()) != 0 {
			t.Errorf("expected 0")
		}
	})
	t.Run("full", func(t *testing.T) {
		s := OTelExporterSpec{Endpoint: "c:4317", Protocol: "grpc", Insecure: true, ServiceName: "v-op", Headers: map[string]string{"a": "b"}}
		if len(s.ContainerEnvs()) != 5 {
			t.Errorf("expected 5")
		}
	})
}

func TestOTelExporterSpec_Validate(t *testing.T) {
	if err := (OTelExporterSpec{}).Validate(); err != nil {
		t.Errorf("empty err=%v", err)
	}
	if err := (OTelExporterSpec{Endpoint: "no-port"}).Validate(); err == nil || !strings.Contains(err.Error(), "port") {
		t.Errorf("expected port error")
	}
	if err := (OTelExporterSpec{Endpoint: "c:4317", Protocol: "weird"}).Validate(); err == nil {
		t.Errorf("expected protocol error")
	}
}
