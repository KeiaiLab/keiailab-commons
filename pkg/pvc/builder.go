// SPDX-License-Identifier: MIT

package pvc

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keiailab/keiailab-commons/pkg/storageclass"
)

// DataPVCParams 는 BuildDataPVC 의 입력이다.
type DataPVCParams struct {
	// Name — PVC/VCT 이름. "" → DefaultVCTName ("data").
	Name        string
	Labels      map[string]string
	Annotations map[string]string

	// AccessModes — 빈 슬라이스 → [ReadWriteOnce] 기본.
	AccessModes []corev1.PersistentVolumeAccessMode
	// StorageClass — storageclass.Normalize 적용: "" → nil (cluster default).
	StorageClass string
	// Size — spec.resources.requests.storage.
	Size resource.Quantity
	// Selector — PV 바인딩 selector (선택).
	Selector *metav1.LabelSelector
}

// BuildDataPVC 는 StatefulSet volumeClaimTemplate 용 data PVC 를 조립한다.
//
// AccessModes 기본 ReadWriteOnce, StorageClass 빈 값은 commons storageclass.Normalize
// 로 nil(cluster default) 처리한다. ObjectMeta.Name 기본은 DefaultVCTName("data").
func BuildDataPVC(p DataPVCParams) corev1.PersistentVolumeClaim {
	name := p.Name
	if name == "" {
		name = DefaultVCTName
	}
	accessModes := p.AccessModes
	if len(accessModes) == 0 {
		accessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	}
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      p.Labels,
			Annotations: p.Annotations,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: storageclass.Normalize(p.StorageClass),
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: p.Size},
			},
			Selector: p.Selector,
		},
	}
}
