// SPDX-License-Identifier: MIT

package apply

import (
	"context"
	"maps"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// deploymentRevisionAnnotation 은 K8s Deployment controller 가 관리하는 rollout
// revision annotation 이다 — operator 가 덮어쓰면 rollout 추적이 깨지므로 보존한다.
const deploymentRevisionAnnotation = "deployment.kubernetes.io/revision"

// StatefulSet 은 desired StatefulSet 을 idempotent 하게 적용한다.
// Selector / ServiceName / VolumeClaimTemplates / PodManagementPolicy 등
// immutable 필드는 Create 시점에만 설정하고, 운영 중에는 Replicas / Template /
// UpdateStrategy / MinReadySeconds / PersistentVolumeClaimRetentionPolicy 만
// 동기화한다.
//
// preserveReplicas=true 시 운영 중인 STS 의 spec.Replicas 를 desired 값으로
// 덮어쓰지 않는다. 두 가지 경우에 사용:
//   - HPA 가 활성화돼 HPA controller 가 Replicas 를 자체 patch 하는 경우.
//   - scale 정책 가드로 replicas 변경이 보류된 경우.
//
// 첫 Create 시점에는 desired.Spec.Replicas 를 그대로 사용해 첫 deploy 를 막지
// 않는다. update conflict 는 retry.RetryOnConflict 로 재시도한다 (안전 방향).
func StatefulSet(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *appsv1.StatefulSet,
	preserveReplicas bool,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		target := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
		}
		_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
			target.Labels = desired.Labels
			target.Annotations = desired.Annotations
			if target.CreationTimestamp.IsZero() {
				target.Spec = desired.Spec
			} else {
				// 운영 중인 STS 는 immutable 필드 (Selector / ServiceName /
				// VolumeClaimTemplates / PodManagementPolicy) 를 그대로 두고
				// 변경 가능한 필드만 동기화.
				if !preserveReplicas {
					target.Spec.Replicas = desired.Spec.Replicas
				}
				target.Spec.Template = desired.Spec.Template
				target.Spec.UpdateStrategy = desired.Spec.UpdateStrategy
				target.Spec.MinReadySeconds = desired.Spec.MinReadySeconds
				target.Spec.PersistentVolumeClaimRetentionPolicy = desired.Spec.PersistentVolumeClaimRetentionPolicy
			}
			return controllerutil.SetControllerReference(owner, target, scheme)
		})
		return err
	})
}

// Deployment 는 desired Deployment 를 idempotent 하게 적용한다.
// Selector 는 immutable 이므로 Create 시점에만 설정한다.
//
// preserveReplicas=true 시 운영 중인 Deployment 의 spec.Replicas 를 보존한다 —
// HPA controller 의 patch 와 operator 의 reconcile 충돌 방지.
//
// server-default 보존 (P0 fix 계보): RevisionHistoryLimit /
// ProgressDeadlineSeconds 는 K8s Deployment controller 가 자동으로 서버
// 기본값 (10, 600) 을 채우는 *pointer 필드* 다. 빌더가 nil 을 desired 로
// 넘기면 매 reconcile 마다 nil 로 덮어씌워지고 K8s 가 즉시 기본값을 재주입
// → 무한 ping-pong (generation 116k+ 폭주 재현). desired 가 비어 있으면
// 운영 중인 객체의 server-defaulted 값을 그대로 두어 fight 를 차단한다.
// pod template 내부의 probe / containerPort / podSpec 기본값도 동일하게
// 보존한다 (deploymentTemplateWithServerDefaults).
func Deployment(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *appsv1.Deployment,
	preserveReplicas bool,
) error {
	target := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = deploymentAnnotationsWithControllerRevision(desired.Annotations, target.Annotations)
		if target.CreationTimestamp.IsZero() {
			target.Spec = desired.Spec
		} else {
			currentTemplate := target.Spec.Template.DeepCopy()
			if !preserveReplicas {
				target.Spec.Replicas = desired.Spec.Replicas
			}
			target.Spec.Template = deploymentTemplateWithServerDefaults(desired.Spec.Template, *currentTemplate)
			target.Spec.Strategy = desired.Spec.Strategy
			target.Spec.MinReadySeconds = desired.Spec.MinReadySeconds
			// nil-guard: desired 가 비어 있으면 server-default 보존.
			if desired.Spec.RevisionHistoryLimit != nil {
				target.Spec.RevisionHistoryLimit = desired.Spec.RevisionHistoryLimit
			}
			if desired.Spec.ProgressDeadlineSeconds != nil {
				target.Spec.ProgressDeadlineSeconds = desired.Spec.ProgressDeadlineSeconds
			}
			target.Spec.Paused = desired.Spec.Paused
		}
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// deploymentAnnotationsWithControllerRevision 은 desired annotation 에 K8s
// Deployment controller 가 부여한 revision annotation 을 합성한다 — desired 가
// 명시하지 않은 경우에만 current 값을 보존.
func deploymentAnnotationsWithControllerRevision(desired, current map[string]string) map[string]string {
	if len(desired) == 0 && current[deploymentRevisionAnnotation] == "" {
		return desired
	}
	out := make(map[string]string, len(desired)+1)
	maps.Copy(out, desired)
	if _, ok := out[deploymentRevisionAnnotation]; !ok {
		if revision := current[deploymentRevisionAnnotation]; revision != "" {
			out[deploymentRevisionAnnotation] = revision
		}
	}
	return out
}

// deploymentTemplateWithServerDefaults 는 current template 을 기저로 desired 의
// 의도 필드만 덮어쓴 pod template 을 만든다 — desired 가 비워둔 server-default
// 필드 (restartPolicy / dnsPolicy / probe 기본값 등) 는 current 값이 잔존해
// update fight 를 차단한다.
func deploymentTemplateWithServerDefaults(desired, current corev1.PodTemplateSpec) corev1.PodTemplateSpec {
	desiredCopy := desired.DeepCopy()
	out := *current.DeepCopy()
	out.ObjectMeta = desiredCopy.ObjectMeta
	out.Spec.SecurityContext = desiredCopy.Spec.SecurityContext
	out.Spec.InitContainers = desiredCopy.Spec.InitContainers
	out.Spec.Containers = desiredCopy.Spec.Containers
	out.Spec.Volumes = desiredCopy.Spec.Volumes
	out.Spec.Affinity = desiredCopy.Spec.Affinity
	out.Spec.NodeSelector = desiredCopy.Spec.NodeSelector
	out.Spec.Tolerations = desiredCopy.Spec.Tolerations
	out.Spec.TopologySpreadConstraints = desiredCopy.Spec.TopologySpreadConstraints
	out.Spec.ImagePullSecrets = desiredCopy.Spec.ImagePullSecrets
	out.Spec.ServiceAccountName = desiredCopy.Spec.ServiceAccountName
	out.Spec.PriorityClassName = desiredCopy.Spec.PriorityClassName
	out.Spec.RuntimeClassName = desiredCopy.Spec.RuntimeClassName
	out.Spec.DNSConfig = desiredCopy.Spec.DNSConfig
	out.Spec.HostAliases = desiredCopy.Spec.HostAliases
	preservePodSpecServerDefaults(&out.Spec, current.Spec)
	return out
}

// preservePodSpecServerDefaults 는 desired pod spec 의 zero-value 필드를 current
// 의 server-defaulted 값으로 backfill 한다.
func preservePodSpecServerDefaults(desired *corev1.PodSpec, current corev1.PodSpec) {
	if desired.RestartPolicy == "" {
		desired.RestartPolicy = current.RestartPolicy
	}
	if desired.DNSPolicy == "" {
		desired.DNSPolicy = current.DNSPolicy
	}
	if desired.SchedulerName == "" {
		desired.SchedulerName = current.SchedulerName
	}
	if desired.TerminationGracePeriodSeconds == nil {
		desired.TerminationGracePeriodSeconds = current.TerminationGracePeriodSeconds
	}
	preserveContainerListServerDefaults(desired.Containers, current.Containers)
	preserveContainerListServerDefaults(desired.InitContainers, current.InitContainers)
}

// preserveContainerListServerDefaults 는 이름이 일치하는 container 쌍에 대해
// server-default 보존을 적용한다 — current 에 없는 신규 container 는 desired 그대로.
func preserveContainerListServerDefaults(desired, current []corev1.Container) {
	currentByName := map[string]corev1.Container{}
	for _, container := range current {
		currentByName[container.Name] = container
	}
	for i := range desired {
		if currentContainer, ok := currentByName[desired[i].Name]; ok {
			preserveContainerServerDefaults(&desired[i], currentContainer)
		}
	}
}

// preserveContainerServerDefaults 는 단일 container 의 zero-value 필드
// (imagePullPolicy / terminationMessage* / probe / port protocol) 를 current
// 의 server-defaulted 값으로 backfill 한다.
func preserveContainerServerDefaults(desired *corev1.Container, current corev1.Container) {
	if desired.ImagePullPolicy == "" {
		desired.ImagePullPolicy = current.ImagePullPolicy
	}
	if desired.TerminationMessagePath == "" {
		desired.TerminationMessagePath = current.TerminationMessagePath
	}
	if desired.TerminationMessagePolicy == "" {
		desired.TerminationMessagePolicy = current.TerminationMessagePolicy
	}
	preserveProbeServerDefaults(desired.LivenessProbe, current.LivenessProbe)
	preserveProbeServerDefaults(desired.ReadinessProbe, current.ReadinessProbe)
	preserveProbeServerDefaults(desired.StartupProbe, current.StartupProbe)
	preserveContainerPortServerDefaults(desired.Ports, current.Ports)
}

// preserveProbeServerDefaults 는 probe 의 zero-value 필드를 current 의
// server-defaulted 값 (timeoutSeconds=1, periodSeconds=10 등) 으로 backfill 한다.
func preserveProbeServerDefaults(desired, current *corev1.Probe) {
	if desired == nil || current == nil {
		return
	}
	if desired.TimeoutSeconds == 0 {
		desired.TimeoutSeconds = current.TimeoutSeconds
	}
	if desired.PeriodSeconds == 0 {
		desired.PeriodSeconds = current.PeriodSeconds
	}
	if desired.SuccessThreshold == 0 {
		desired.SuccessThreshold = current.SuccessThreshold
	}
	if desired.FailureThreshold == 0 {
		desired.FailureThreshold = current.FailureThreshold
	}
	if desired.TerminationGracePeriodSeconds == nil {
		desired.TerminationGracePeriodSeconds = current.TerminationGracePeriodSeconds
	}
}

// preserveContainerPortServerDefaults 는 이름이 일치하는 port 의 protocol
// zero-value 를 current 의 server-defaulted 값 (TCP) 으로 backfill 한다.
func preserveContainerPortServerDefaults(desired, current []corev1.ContainerPort) {
	currentByName := map[string]corev1.ContainerPort{}
	for _, port := range current {
		currentByName[port.Name] = port
	}
	for i := range desired {
		if desired[i].Protocol != "" {
			continue
		}
		if currentPort, ok := currentByName[desired[i].Name]; ok {
			desired[i].Protocol = currentPort.Protocol
		}
	}
}
