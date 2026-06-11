// SPDX-License-Identifier: MIT

package certmanager

import (
	"reflect"
	"testing"
)

// TestBuildSelfSignedIssuer_Golden 은 valkey-operator BuildSelfSignedIssuer
// 원본 출력 (name/labels 결정만 호출자로 이동) 의 raw map 리터럴 골든 검증.
func TestBuildSelfSignedIssuer_Golden(t *testing.T) {
	t.Parallel()
	got := BuildSelfSignedIssuer("vk1-selfsigned", "cache", map[string]string{
		"app.kubernetes.io/instance": "vk1",
	})
	want := map[string]any{
		"apiVersion": "cert-manager.io/v1",
		"kind":       "Issuer",
		"metadata": map[string]any{
			"name":      "vk1-selfsigned",
			"namespace": "cache",
			"labels": map[string]any{
				"app.kubernetes.io/instance": "vk1",
			},
		},
		"spec": map[string]any{
			"selfSigned": map[string]any{},
		},
	}
	if !reflect.DeepEqual(got.Object, want) {
		t.Errorf("Object mismatch:\n got = %#v\nwant = %#v", got.Object, want)
	}
}

// TestBuildSelfSignedIssuer_NilLabels 는 labels nil 시 metadata.labels 미설정
// 을 검증한다.
func TestBuildSelfSignedIssuer_NilLabels(t *testing.T) {
	t.Parallel()
	got := BuildSelfSignedIssuer("x-selfsigned", "ns1", nil)
	meta := got.Object["metadata"].(map[string]any)
	if _, ok := meta["labels"]; ok {
		t.Errorf("metadata.labels 가 설정됨 (nil 은 미설정이어야 함): %v", meta["labels"])
	}
}

// TestBuildSelfSignedIssuer_GVK 는 IssuerGVK 상수 정합 (cert-manager.io/v1
// Issuer) 과 DeepCopy 안전성을 검증한다.
func TestBuildSelfSignedIssuer_GVK(t *testing.T) {
	t.Parallel()
	got := BuildSelfSignedIssuer("a", "b", nil)
	if gvk := got.GroupVersionKind(); gvk != IssuerGVK {
		t.Errorf("GVK = %v, want %v", gvk, IssuerGVK)
	}
	// spec 이 JSON-safe map 으로 저장됐는지 DeepCopy 로 회귀 가드.
	if cp := got.DeepCopy(); cp == nil {
		t.Fatal("DeepCopy = nil")
	}
}
