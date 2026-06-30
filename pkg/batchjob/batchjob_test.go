// SPDX-License-Identifier: MIT

package batchjob_test

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	"github.com/keiailab/keiailab-commons/pkg/batchjob"
)

func TestBuild_EnvelopeAndLabelPropagation(t *testing.T) {
	t.Parallel()
	got := batchjob.Build(batchjob.Params{
		Name: "db-backup", Namespace: "ns",
		Labels:                  map[string]string{"app.kubernetes.io/component": "backup"},
		BackoffLimit:            ptr.To[int32](2),
		TTLSecondsAfterFinished: ptr.To[int32](86400),
		Containers:              []corev1.Container{{Name: "c", Image: "img"}},
	})
	if got.Name != "db-backup" || got.Namespace != "ns" {
		t.Fatalf("ObjectMeta mismatch: %s/%s", got.Namespace, got.Name)
	}
	// Labels 는 Job 과 Pod template 양쪽.
	if got.Labels["app.kubernetes.io/component"] != "backup" {
		t.Fatalf("Job labels 누락: %v", got.Labels)
	}
	if got.Spec.Template.Labels["app.kubernetes.io/component"] != "backup" {
		t.Fatalf("Pod template labels 전파 실패: %v", got.Spec.Template.Labels)
	}
	if got.Spec.BackoffLimit == nil || *got.Spec.BackoffLimit != 2 {
		t.Fatalf("BackoffLimit=%v, want 2", got.Spec.BackoffLimit)
	}
	if got.Spec.TTLSecondsAfterFinished == nil || *got.Spec.TTLSecondsAfterFinished != 86400 {
		t.Fatalf("TTL=%v, want 86400", got.Spec.TTLSecondsAfterFinished)
	}
}

func TestBuild_RestartPolicyDefaultsOnFailure(t *testing.T) {
	t.Parallel()
	got := batchjob.Build(batchjob.Params{Name: "x"})
	if got.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyOnFailure {
		t.Fatalf("기본 RestartPolicy=%q, want OnFailure", got.Spec.Template.Spec.RestartPolicy)
	}
}

func TestBuild_RestartPolicyOverride(t *testing.T) {
	t.Parallel()
	got := batchjob.Build(batchjob.Params{Name: "x", RestartPolicy: corev1.RestartPolicyNever})
	if got.Spec.Template.Spec.RestartPolicy != corev1.RestartPolicyNever {
		t.Fatalf("RestartPolicy override 무시: %q", got.Spec.Template.Spec.RestartPolicy)
	}
}

func TestBuild_NilLimitsUnset(t *testing.T) {
	t.Parallel()
	got := batchjob.Build(batchjob.Params{Name: "x"})
	if got.Spec.BackoffLimit != nil || got.Spec.TTLSecondsAfterFinished != nil {
		t.Fatal("nil limit 은 미설정이어야 함")
	}
}

func TestBuild_ContainersVolumesSC(t *testing.T) {
	t.Parallel()
	sc := &corev1.PodSecurityContext{RunAsNonRoot: ptr.To(true)}
	got := batchjob.Build(batchjob.Params{
		Name:               "x",
		Containers:         []corev1.Container{{Name: "a"}, {Name: "b"}},
		Volumes:            []corev1.Volume{{Name: "v"}},
		PodSecurityContext: sc,
	})
	if len(got.Spec.Template.Spec.Containers) != 2 || len(got.Spec.Template.Spec.Volumes) != 1 {
		t.Fatal("containers/volumes 전달 실패")
	}
	if got.Spec.Template.Spec.SecurityContext == nil || got.Spec.Template.Spec.SecurityContext.RunAsNonRoot == nil {
		t.Fatal("PodSecurityContext 전달 실패")
	}
}
