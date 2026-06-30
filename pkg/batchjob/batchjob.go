// SPDX-License-Identifier: MIT

package batchjob

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Params 는 Build 의 입력이다. 컨테이너/볼륨/SecurityContext 는 호출자가 도메인
// 규칙대로 조립해 전달한다.
type Params struct {
	// Name / Namespace / Labels — Job ObjectMeta. Labels 는 Pod template 에도 전파.
	Name      string
	Namespace string
	Labels    map[string]string

	// BackoffLimit / TTLSecondsAfterFinished — nil 이면 해당 필드 미설정.
	BackoffLimit            *int32
	TTLSecondsAfterFinished *int32

	// RestartPolicy — "" → OnFailure (일회성 Job 표준).
	RestartPolicy corev1.RestartPolicy

	// Containers / Volumes / PodSecurityContext — Pod 명세 (도메인 책임).
	Containers         []corev1.Container
	Volumes            []corev1.Volume
	PodSecurityContext *corev1.PodSecurityContext
}

// Build 는 Params 로 일회성 batch/v1 Job 을 조립한다. Labels 는 Job 과 Pod
// template 양쪽에 적용되고, RestartPolicy 기본은 OnFailure 다.
func Build(p Params) *batchv1.Job {
	restart := p.RestartPolicy
	if restart == "" {
		restart = corev1.RestartPolicyOnFailure
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
			Labels:    p.Labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            p.BackoffLimit,
			TTLSecondsAfterFinished: p.TTLSecondsAfterFinished,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: p.Labels},
				Spec: corev1.PodSpec{
					RestartPolicy:   restart,
					Containers:      p.Containers,
					Volumes:         p.Volumes,
					SecurityContext: p.PodSecurityContext,
				},
			},
		},
	}
}
