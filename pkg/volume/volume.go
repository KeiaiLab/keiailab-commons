// SPDX-License-Identifier: MIT

package volume

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// CertDefaultMode — TLS cert Secret 의 DefaultMode (owner read-only). cert/key
// 파일 권한 검사를 통과하고 최소권한을 보장하는 보안 불변식.
const CertDefaultMode int32 = 0o400

// TLSSecretMount — cert-manager 가 발급한 TLS Secret 을 read-only(0o400)로
// 마운트하는 Volume + VolumeMount 쌍을 반환한다.
//
// volName 은 Volume 과 VolumeMount 를 잇는 이름(예: "tls" / "tls-server"),
// secretName 은 Secret 이름, mountPath 는 컨테이너 마운트 경로다. empty-guard
// (secretName=="" 시 미마운트)는 호출자 책임이다.
func TLSSecretMount(volName, secretName, mountPath string) (corev1.Volume, corev1.VolumeMount) {
	vol := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  secretName,
				DefaultMode: ptr.To(CertDefaultMode),
			},
		},
	}
	mount := corev1.VolumeMount{
		Name:      volName,
		MountPath: mountPath,
		ReadOnly:  true,
	}
	return vol, mount
}
