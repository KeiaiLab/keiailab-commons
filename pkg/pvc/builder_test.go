// SPDX-License-Identifier: MIT

package pvc_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keiailab/keiailab-commons/pkg/pvc"
)

func TestBuildDataPVC_Defaults(t *testing.T) {
	t.Parallel()
	got := pvc.BuildDataPVC(pvc.DataPVCParams{Size: resource.MustParse("8Gi")})
	if got.Name != "data" {
		t.Fatalf("기본 이름=%q, want data", got.Name)
	}
	if len(got.Spec.AccessModes) != 1 || got.Spec.AccessModes[0] != corev1.ReadWriteOnce {
		t.Fatalf("기본 accessMode=%v, want [RWO]", got.Spec.AccessModes)
	}
	if got.Spec.StorageClassName != nil {
		t.Fatalf("빈 StorageClass 는 nil 이어야 함, got %v", got.Spec.StorageClassName)
	}
	q := got.Spec.Resources.Requests[corev1.ResourceStorage]
	if q.String() != "8Gi" {
		t.Fatalf("size=%s, want 8Gi", q.String())
	}
}

func TestBuildDataPVC_StorageClassNormalized(t *testing.T) {
	t.Parallel()
	got := pvc.BuildDataPVC(pvc.DataPVCParams{StorageClass: "fast-ssd", Size: resource.MustParse("1Gi")})
	if got.Spec.StorageClassName == nil || *got.Spec.StorageClassName != "fast-ssd" {
		t.Fatalf("StorageClassName=%v, want fast-ssd", got.Spec.StorageClassName)
	}
}

func TestBuildDataPVC_LabelsAnnotationsSelector(t *testing.T) {
	t.Parallel()
	sel := &metav1.LabelSelector{MatchLabels: map[string]string{"zone": "a"}}
	got := pvc.BuildDataPVC(pvc.DataPVCParams{
		Name: "data", Labels: map[string]string{"k": "v"}, Annotations: map[string]string{"a": "b"},
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
		Size:        resource.MustParse("2Gi"), Selector: sel,
	})
	if got.Labels["k"] != "v" || got.Annotations["a"] != "b" {
		t.Fatalf("labels/annotations 전달 실패: %+v / %+v", got.Labels, got.Annotations)
	}
	if got.Spec.AccessModes[0] != corev1.ReadWriteMany {
		t.Fatalf("accessMode override 무시: %v", got.Spec.AccessModes)
	}
	if got.Spec.Selector == nil || got.Spec.Selector.MatchLabels["zone"] != "a" {
		t.Fatalf("selector 전달 실패: %+v", got.Spec.Selector)
	}
}
