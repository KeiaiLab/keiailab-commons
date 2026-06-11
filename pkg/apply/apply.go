// SPDX-License-Identifier: MIT

package apply

import (
	"context"
	"fmt"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ConfigMap 은 desired ConfigMap 을 idempotent 하게 적용한다.
// Data / BinaryData / Labels / Annotations 만 갱신하며 owner reference 를 설정한다.
func ConfigMap(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *corev1.ConfigMap,
) error {
	target := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = desired.Annotations
		target.Data = desired.Data
		target.BinaryData = desired.BinaryData
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// Service 는 desired Service 를 idempotent 하게 적용한다.
//
// Spec.ClusterIP / IPFamilies / IPFamilyPolicy 는 immutable 이므로 Create
// 시점에만 설정한다 — Create 시 desired 그대로 (headless 면 "None"), Update
// 시 K8s 가 할당한 값을 그대로 둔다.
//
// LoadBalancer / NodePort 운영 필드 (LoadBalancerIP, externalIPs, PROXY
// protocol 헬스 체크 NodePort 등) 는 첫 Create 후에도 반영되어야 한다 —
// 과거 누락 시 LoadBalancerIP 영구 고착 P0 결함의 회귀 가드.
//
// server-default pointer 필드 (LoadBalancerClass / AllocateLoadBalancerNodePorts /
// InternalTrafficPolicy) 는 desired 가 nil 이면 운영 중 객체의 server-defaulted
// 값을 보존한다 — Deployment 의 RevisionHistoryLimit 가드와 동일한 ping-pong
// 클래스 방어. SessionAffinityConfig 는 affinity=None 전환 시 잔존 값이 API
// validation 에 걸리므로 nil-guard 없이 desired 를 단일 진실로 동기화한다.
func Service(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *corev1.Service,
) error {
	target := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = desired.Annotations
		if target.CreationTimestamp.IsZero() {
			target.Spec.ClusterIP = desired.Spec.ClusterIP
			target.Spec.IPFamilies = desired.Spec.IPFamilies
			target.Spec.IPFamilyPolicy = desired.Spec.IPFamilyPolicy
		}
		target.Spec.Ports = desired.Spec.Ports
		target.Spec.Selector = desired.Spec.Selector
		target.Spec.Type = desired.Spec.Type
		target.Spec.PublishNotReadyAddresses = desired.Spec.PublishNotReadyAddresses
		target.Spec.SessionAffinity = desired.Spec.SessionAffinity
		target.Spec.SessionAffinityConfig = desired.Spec.SessionAffinityConfig
		target.Spec.LoadBalancerIP = desired.Spec.LoadBalancerIP
		target.Spec.LoadBalancerSourceRanges = desired.Spec.LoadBalancerSourceRanges
		target.Spec.ExternalIPs = desired.Spec.ExternalIPs
		target.Spec.ExternalName = desired.Spec.ExternalName
		target.Spec.ExternalTrafficPolicy = desired.Spec.ExternalTrafficPolicy
		target.Spec.HealthCheckNodePort = desired.Spec.HealthCheckNodePort
		// nil-guard: desired 가 비어 있으면 server-default 보존 (ping-pong 차단).
		if desired.Spec.LoadBalancerClass != nil {
			target.Spec.LoadBalancerClass = desired.Spec.LoadBalancerClass
		}
		if desired.Spec.AllocateLoadBalancerNodePorts != nil {
			target.Spec.AllocateLoadBalancerNodePorts = desired.Spec.AllocateLoadBalancerNodePorts
		}
		if desired.Spec.InternalTrafficPolicy != nil {
			target.Spec.InternalTrafficPolicy = desired.Spec.InternalTrafficPolicy
		}
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// NetworkPolicy 는 desired NetworkPolicy 를 idempotent 하게 적용한다.
// PodSelector / PolicyTypes / Ingress / Egress 를 매번 desired 그대로 동기화 —
// 멤버 라벨 / 포트가 spec 변경에 따라 갱신되어도 추적 가능.
func NetworkPolicy(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *networkingv1.NetworkPolicy,
) error {
	target := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = desired.Annotations
		target.Spec.PodSelector = desired.Spec.PodSelector
		target.Spec.PolicyTypes = desired.Spec.PolicyTypes
		target.Spec.Ingress = desired.Spec.Ingress
		target.Spec.Egress = desired.Spec.Egress
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// PDB 는 desired PodDisruptionBudget 을 idempotent 하게 적용한다.
// MinAvailable / MaxUnavailable 은 K8s 제약상 둘 중 하나만 설정 가능하므로
// desired 그대로 mutate 한다 (둘 중 하나만 non-nil 이라는 호출자 책임).
func PDB(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *policyv1.PodDisruptionBudget,
) error {
	target := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = desired.Annotations
		target.Spec.Selector = desired.Spec.Selector
		target.Spec.MinAvailable = desired.Spec.MinAvailable
		target.Spec.MaxUnavailable = desired.Spec.MaxUnavailable
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// HPA 는 desired HorizontalPodAutoscaler 를 idempotent 하게 적용한다.
// desired==nil 이면 *기존 HPA 를 삭제* 한다 — autoscaling toggle off 회복 경로.
// name / namespace 는 삭제 경로의 lookup 키로 desired 와 무관하게 항상 전달한다.
func HPA(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	name, namespace string,
	desired *autoscalingv2.HorizontalPodAutoscaler,
) error {
	if desired == nil {
		// disable: 기존 HPA 삭제 (있으면).
		existing := &autoscalingv2.HorizontalPodAutoscaler{}
		if err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, existing); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		return c.Delete(ctx, existing)
	}
	target := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{Name: desired.Name, Namespace: desired.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.Labels = desired.Labels
		target.Annotations = desired.Annotations
		target.Spec = desired.Spec
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	return err
}

// Unstructured 는 desired unstructured 객체 (ServiceMonitor / Certificate 등
// CRD 의존 리소스) 를 idempotent 하게 적용한다. desired==nil 이면 no-op.
//
// failSoftNoMatch=true 면 대상 CRD 가 클러스터에 미설치인 경우 (NotFound /
// NoMatch 에러) 를 fail-soft 처리한다 (nil 반환) — CRD 설치 후 다음 reconcile
// 에서 자동 생성된다.
func Unstructured(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	desired *unstructured.Unstructured,
	failSoftNoMatch bool,
) error {
	if desired == nil {
		return nil
	}
	target := &unstructured.Unstructured{}
	target.SetGroupVersionKind(desired.GroupVersionKind())
	target.SetName(desired.GetName())
	target.SetNamespace(desired.GetNamespace())
	_, err := controllerutil.CreateOrUpdate(ctx, c, target, func() error {
		target.SetLabels(desired.GetLabels())
		target.SetAnnotations(desired.GetAnnotations())
		target.Object["spec"] = desired.Object["spec"]
		return controllerutil.SetControllerReference(owner, target, scheme)
	})
	if err != nil && failSoftNoMatch && (apierrors.IsNotFound(err) || meta.IsNoMatchError(err)) {
		return nil
	}
	return err
}

// SecretIfNotExists 는 Secret 이 존재하지 않으면 build() 로 생성한다.
// 존재 시 no-op (멱등) 이며 build 는 호출되지 않는다 — 랜덤 credential 생성
// 비용 / 재생성 사고 차단. owner reference + scheme 은 GC 를 위해 필수.
//
// keyfile / password Secret 처럼 *immutable + create-once* 자원에 사용한다.
// CreateOrUpdate 변형은 의도적으로 제공하지 않는다 (doc.go 설계 원칙).
//
// 병렬 reconcile 의 *Get NotFound → Create AlreadyExists* race 는 에러가
// 아니라 정상 경로로 흡수한다 (race-tolerant).
func SecretIfNotExists(
	ctx context.Context,
	c client.Client,
	scheme *runtime.Scheme,
	owner client.Object,
	secretName string,
	build func() *corev1.Secret,
) error {
	existing := &corev1.Secret{}
	err := c.Get(ctx, client.ObjectKey{Name: secretName, Namespace: owner.GetNamespace()}, existing)
	if err == nil {
		return nil // 이미 존재 — no-op
	}
	if !apierrors.IsNotFound(err) {
		return err
	}
	secret := build()
	if err := controllerutil.SetControllerReference(owner, secret, scheme); err != nil {
		return fmt.Errorf("set owner ref: %w", err)
	}
	if err := c.Create(ctx, secret); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}
