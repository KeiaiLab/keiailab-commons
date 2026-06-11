// SPDX-License-Identifier: MIT

package apply

import (
	"context"
	"errors"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const (
	stsName = "web"
	depName = "web-deploy"
	cName   = "main"
	imgV1   = "app:v1"
	imgV2   = "app:v2"
)

func podTemplate(image string) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{labelKey: "w"}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: cName, Image: image}},
		},
	}
}

func newSTS(image string, replicas int32, mut func(*appsv1.StatefulSet)) *appsv1.StatefulSet {
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{Name: stsName, Namespace: testNS},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    new(replicas),
			ServiceName: "svc-new",
			Selector:    &metav1.LabelSelector{MatchLabels: map[string]string{labelKey: "w"}},
			Template:    podTemplate(image),
		},
	}
	if mut != nil {
		mut(sts)
	}
	return sts
}

func seededSTS(replicas int32, mut func(*appsv1.StatefulSet)) *appsv1.StatefulSet {
	sts := newSTS(imgV1, replicas, mut)
	sts.CreationTimestamp = seedTime()
	return sts
}

func TestStatefulSet(t *testing.T) {
	cases := []struct {
		name             string
		existing         *appsv1.StatefulSet
		desired          *appsv1.StatefulSet
		preserveReplicas bool
		want             func(t *testing.T, got *appsv1.StatefulSet)
	}{
		{
			name:    "신규 생성 시 전체 spec 적용 + ownerRef",
			desired: newSTS(imgV1, 3, nil),
			want: func(t *testing.T, got *appsv1.StatefulSet) {
				if got.Spec.ServiceName != "svc-new" || *got.Spec.Replicas != 3 {
					t.Errorf("spec = %s/%d, 기대 svc-new/3", got.Spec.ServiceName, *got.Spec.Replicas)
				}
				assertOwnerRef(t, got)
			},
		},
		{
			name: "갱신 시 immutable 필드 보존 + mutable 필드 동기화",
			existing: seededSTS(3, func(sts *appsv1.StatefulSet) {
				sts.Spec.ServiceName = "svc-old"
				sts.Spec.Selector = &metav1.LabelSelector{MatchLabels: map[string]string{labelKey: "w-old"}}
				sts.Spec.PodManagementPolicy = appsv1.ParallelPodManagement
				sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Name: "data"}},
				}
			}),
			desired: newSTS(imgV2, 3, func(sts *appsv1.StatefulSet) {
				sts.Spec.Selector = &metav1.LabelSelector{MatchLabels: map[string]string{labelKey: "w-new"}}
				sts.Spec.PodManagementPolicy = appsv1.OrderedReadyPodManagement
				sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Name: "data2"}},
				}
				sts.Spec.MinReadySeconds = 5
			}),
			want: func(t *testing.T, got *appsv1.StatefulSet) {
				if got.Spec.ServiceName != "svc-old" {
					t.Errorf("ServiceName = %q, 기존 svc-old 보존 기대 (immutable)", got.Spec.ServiceName)
				}
				if got.Spec.Selector.MatchLabels[labelKey] != "w-old" {
					t.Errorf("Selector = %v, 기존 보존 기대 (immutable)", got.Spec.Selector.MatchLabels)
				}
				if got.Spec.PodManagementPolicy != appsv1.ParallelPodManagement {
					t.Errorf("PodManagementPolicy = %q, 기존 Parallel 보존 기대", got.Spec.PodManagementPolicy)
				}
				if len(got.Spec.VolumeClaimTemplates) != 1 || got.Spec.VolumeClaimTemplates[0].Name != "data" {
					t.Errorf("VCT = %v, 기존 보존 기대 (immutable)", got.Spec.VolumeClaimTemplates)
				}
				if img := got.Spec.Template.Spec.Containers[0].Image; img != imgV2 {
					t.Errorf("Template image = %q, %q 동기화 기대", img, imgV2)
				}
				if got.Spec.MinReadySeconds != 5 {
					t.Errorf("MinReadySeconds = %d, 5 동기화 기대", got.Spec.MinReadySeconds)
				}
			},
		},
		{
			name:             "갱신 + preserveReplicas=true 시 운영 replicas 보존",
			existing:         seededSTS(5, nil),
			desired:          newSTS(imgV1, 2, nil),
			preserveReplicas: true,
			want: func(t *testing.T, got *appsv1.StatefulSet) {
				if *got.Spec.Replicas != 5 {
					t.Errorf("Replicas = %d, 기존 5 보존 기대 (HPA/scale 가드)", *got.Spec.Replicas)
				}
			},
		},
		{
			name:     "갱신 + preserveReplicas=false 시 desired replicas 동기화",
			existing: seededSTS(5, nil),
			desired:  newSTS(imgV1, 2, nil),
			want: func(t *testing.T, got *appsv1.StatefulSet) {
				if *got.Spec.Replicas != 2 {
					t.Errorf("Replicas = %d, 2 동기화 기대", *got.Spec.Replicas)
				}
			},
		},
		{
			name:             "신규 + preserveReplicas=true 도 첫 deploy 는 desired replicas 사용",
			desired:          newSTS(imgV1, 2, nil),
			preserveReplicas: true,
			want: func(t *testing.T, got *appsv1.StatefulSet) {
				if *got.Spec.Replicas != 2 {
					t.Errorf("Replicas = %d, 첫 Create 는 desired 2 기대", *got.Spec.Replicas)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := testScheme(t)
			b := fake.NewClientBuilder().WithScheme(s)
			if tc.existing != nil {
				b = b.WithObjects(tc.existing)
			}
			c := b.Build()

			if err := StatefulSet(t.Context(), c, s, newOwner(), tc.desired, tc.preserveReplicas); err != nil {
				t.Fatalf("StatefulSet apply 실패: %v", err)
			}

			got := &appsv1.StatefulSet{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: stsName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			tc.want(t, got)
		})
	}
}

// TestStatefulSetConflictRetry 는 update conflict 1회 발생 시 RetryOnConflict 가
// 재시도해 최종 성공하는지 검증한다 (valkey 채택 안전 방향의 회귀 가드).
func TestStatefulSetConflictRetry(t *testing.T) {
	s := testScheme(t)
	conflicted := false
	c := fake.NewClientBuilder().WithScheme(s).WithObjects(seededSTS(3, nil)).
		WithInterceptorFuncs(interceptor.Funcs{
			Update: func(ctx context.Context, cl client.WithWatch, obj client.Object,
				opts ...client.UpdateOption) error {
				if !conflicted {
					conflicted = true
					return apierrors.NewConflict(
						schema.GroupResource{Group: "apps", Resource: "statefulsets"},
						obj.GetName(), errors.New("동시 수정"))
				}
				return cl.Update(ctx, obj, opts...)
			},
		}).Build()

	if err := StatefulSet(t.Context(), c, s, newOwner(), newSTS(imgV2, 3, nil), false); err != nil {
		t.Fatalf("conflict 재시도 후 성공 기대, got %v", err)
	}
	if !conflicted {
		t.Fatal("conflict 경로 미진입 — 테스트 전제 무효")
	}
	got := &appsv1.StatefulSet{}
	if err := c.Get(t.Context(), client.ObjectKey{Name: stsName, Namespace: testNS}, got); err != nil {
		t.Fatalf("get 실패: %v", err)
	}
	if img := got.Spec.Template.Spec.Containers[0].Image; img != imgV2 {
		t.Errorf("재시도 후 image = %q, %q 기대", img, imgV2)
	}
}

func newDeploy(image string, replicas int32, mut func(*appsv1.Deployment)) *appsv1.Deployment {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: depName, Namespace: testNS},
		Spec: appsv1.DeploymentSpec{
			Replicas: new(replicas),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{labelKey: "w"}},
			Template: podTemplate(image),
		},
	}
	if mut != nil {
		mut(dep)
	}
	return dep
}

func seededDeploy(replicas int32, mut func(*appsv1.Deployment)) *appsv1.Deployment {
	dep := newDeploy(imgV1, replicas, mut)
	dep.CreationTimestamp = seedTime()
	return dep
}

// assertTemplateServerDefaults 는 pod template server-default 보존 케이스의 검증
// 본문 — TestDeployment 의 cyclomatic complexity 분리 (gocyclo 한도 준수).
func assertTemplateServerDefaults(t *testing.T, got *appsv1.Deployment) {
	t.Helper()
	spec := got.Spec.Template.Spec
	main := spec.Containers[0]
	if main.Image != imgV2 {
		t.Errorf("image = %q, %q 동기화 기대", main.Image, imgV2)
	}
	if main.ImagePullPolicy != corev1.PullIfNotPresent {
		t.Errorf("ImagePullPolicy = %q, server-default 보존 기대", main.ImagePullPolicy)
	}
	if main.TerminationMessagePath != "/dev/termination-log" {
		t.Errorf("TerminationMessagePath = %q, 보존 기대", main.TerminationMessagePath)
	}
	if main.Ports[0].Protocol != corev1.ProtocolTCP {
		t.Errorf("port Protocol = %q, server-default TCP 보존 기대", main.Ports[0].Protocol)
	}
	probe := main.LivenessProbe
	if probe.TimeoutSeconds != 1 || probe.PeriodSeconds != 10 ||
		probe.SuccessThreshold != 1 || probe.FailureThreshold != 3 {
		t.Errorf("probe server-default 미보존: %+v", probe)
	}
	if spec.RestartPolicy != corev1.RestartPolicyAlways || spec.DNSPolicy != corev1.DNSClusterFirst ||
		spec.SchedulerName != "default-scheduler" {
		t.Errorf("podSpec server-default 미보존: %s/%s/%s",
			spec.RestartPolicy, spec.DNSPolicy, spec.SchedulerName)
	}
	if spec.TerminationGracePeriodSeconds == nil || *spec.TerminationGracePeriodSeconds != 30 {
		t.Errorf("TGPS = %v, 30 보존 기대", spec.TerminationGracePeriodSeconds)
	}
	side := spec.Containers[1]
	if side.Name != "side" || side.ImagePullPolicy != "" {
		t.Errorf("신규 container 는 desired 그대로 기대, got %+v", side)
	}
}

func TestDeployment(t *testing.T) {
	cases := []struct {
		name             string
		existing         *appsv1.Deployment
		desired          *appsv1.Deployment
		preserveReplicas bool
		want             func(t *testing.T, got *appsv1.Deployment)
	}{
		{
			name:    "신규 생성 시 전체 spec 적용 + ownerRef",
			desired: newDeploy(imgV1, 2, nil),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if *got.Spec.Replicas != 2 || got.Spec.Template.Spec.Containers[0].Image != imgV1 {
					t.Errorf("spec 부정합: %+v", got.Spec)
				}
				assertOwnerRef(t, got)
			},
		},
		{
			name: "갱신 시 RevisionHistoryLimit/ProgressDeadlineSeconds desired nil 이면 server-default 보존",
			existing: seededDeploy(2, func(d *appsv1.Deployment) {
				d.Spec.RevisionHistoryLimit = new(int32(10))
				d.Spec.ProgressDeadlineSeconds = new(int32(600))
			}),
			desired: newDeploy(imgV2, 2, nil),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if got.Spec.RevisionHistoryLimit == nil || *got.Spec.RevisionHistoryLimit != 10 {
					t.Errorf("RevisionHistoryLimit = %v, server-default 10 보존 기대 (ping-pong 차단)",
						got.Spec.RevisionHistoryLimit)
				}
				if got.Spec.ProgressDeadlineSeconds == nil || *got.Spec.ProgressDeadlineSeconds != 600 {
					t.Errorf("ProgressDeadlineSeconds = %v, server-default 600 보존 기대",
						got.Spec.ProgressDeadlineSeconds)
				}
				if got.Spec.Template.Spec.Containers[0].Image != imgV2 {
					t.Errorf("image = %q, %q 동기화 기대", got.Spec.Template.Spec.Containers[0].Image, imgV2)
				}
			},
		},
		{
			name: "갱신 시 RevisionHistoryLimit/ProgressDeadlineSeconds desired 명시면 동기화",
			existing: seededDeploy(2, func(d *appsv1.Deployment) {
				d.Spec.RevisionHistoryLimit = new(int32(10))
				d.Spec.ProgressDeadlineSeconds = new(int32(600))
			}),
			desired: newDeploy(imgV1, 2, func(d *appsv1.Deployment) {
				d.Spec.RevisionHistoryLimit = new(int32(3))
				d.Spec.ProgressDeadlineSeconds = new(int32(300))
			}),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if got.Spec.RevisionHistoryLimit == nil || *got.Spec.RevisionHistoryLimit != 3 {
					t.Errorf("RevisionHistoryLimit = %v, 3 동기화 기대", got.Spec.RevisionHistoryLimit)
				}
				if got.Spec.ProgressDeadlineSeconds == nil || *got.Spec.ProgressDeadlineSeconds != 300 {
					t.Errorf("ProgressDeadlineSeconds = %v, 300 동기화 기대", got.Spec.ProgressDeadlineSeconds)
				}
			},
		},
		{
			name: "갱신 시 K8s controller 의 revision annotation 보존",
			existing: seededDeploy(2, func(d *appsv1.Deployment) {
				d.Annotations = map[string]string{deploymentRevisionAnnotation: "7"}
			}),
			desired: newDeploy(imgV2, 2, func(d *appsv1.Deployment) {
				d.Annotations = map[string]string{"x": "y"}
			}),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if got.Annotations[deploymentRevisionAnnotation] != "7" {
					t.Errorf("revision annotation = %q, 기존 7 보존 기대", got.Annotations[deploymentRevisionAnnotation])
				}
				if got.Annotations["x"] != "y" {
					t.Errorf("desired annotation 미반영: %v", got.Annotations)
				}
			},
		},
		{
			name: "갱신 시 desired 가 revision annotation 명시면 desired 우선",
			existing: seededDeploy(2, func(d *appsv1.Deployment) {
				d.Annotations = map[string]string{deploymentRevisionAnnotation: "7"}
			}),
			desired: newDeploy(imgV2, 2, func(d *appsv1.Deployment) {
				d.Annotations = map[string]string{deploymentRevisionAnnotation: "9"}
			}),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if got.Annotations[deploymentRevisionAnnotation] != "9" {
					t.Errorf("revision annotation = %q, desired 9 우선 기대", got.Annotations[deploymentRevisionAnnotation])
				}
			},
		},
		{
			name:             "갱신 + preserveReplicas=true 시 운영 replicas 보존 (HPA 충돌 방지)",
			existing:         seededDeploy(5, nil),
			desired:          newDeploy(imgV1, 2, nil),
			preserveReplicas: true,
			want: func(t *testing.T, got *appsv1.Deployment) {
				if *got.Spec.Replicas != 5 {
					t.Errorf("Replicas = %d, 기존 5 보존 기대", *got.Spec.Replicas)
				}
			},
		},
		{
			name:     "갱신 + preserveReplicas=false 시 desired replicas 동기화",
			existing: seededDeploy(5, nil),
			desired:  newDeploy(imgV1, 2, nil),
			want: func(t *testing.T, got *appsv1.Deployment) {
				if *got.Spec.Replicas != 2 {
					t.Errorf("Replicas = %d, 2 동기화 기대", *got.Spec.Replicas)
				}
			},
		},
		{
			name: "갱신 시 pod template 의 server-default 보존 + desired 의도 필드 동기화",
			existing: seededDeploy(2, func(d *appsv1.Deployment) {
				spec := &d.Spec.Template.Spec
				spec.RestartPolicy = corev1.RestartPolicyAlways
				spec.DNSPolicy = corev1.DNSClusterFirst
				spec.SchedulerName = "default-scheduler"
				spec.TerminationGracePeriodSeconds = new(int64(30))
				main := &spec.Containers[0]
				main.ImagePullPolicy = corev1.PullIfNotPresent
				main.TerminationMessagePath = "/dev/termination-log"
				main.TerminationMessagePolicy = corev1.TerminationMessageReadFile
				main.Ports = []corev1.ContainerPort{{Name: "http", ContainerPort: 8080, Protocol: corev1.ProtocolTCP}}
				main.LivenessProbe = &corev1.Probe{
					ProbeHandler:     corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"true"}}},
					TimeoutSeconds:   1,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 3,
				}
			}),
			desired: newDeploy(imgV2, 2, func(d *appsv1.Deployment) {
				main := &d.Spec.Template.Spec.Containers[0]
				main.Ports = []corev1.ContainerPort{{Name: "http", ContainerPort: 8080}}
				main.LivenessProbe = &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"true"}}},
				}
				d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers,
					corev1.Container{Name: "side", Image: "side:v1"})
			}),
			want: assertTemplateServerDefaults,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := testScheme(t)
			b := fake.NewClientBuilder().WithScheme(s)
			if tc.existing != nil {
				b = b.WithObjects(tc.existing)
			}
			c := b.Build()

			if err := Deployment(t.Context(), c, s, newOwner(), tc.desired, tc.preserveReplicas); err != nil {
				t.Fatalf("Deployment apply 실패: %v", err)
			}

			got := &appsv1.Deployment{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: depName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			tc.want(t, got)
		})
	}
}
