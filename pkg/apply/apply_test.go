// SPDX-License-Identifier: MIT

package apply

import (
	"context"
	"maps"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const (
	testNS    = "ns"
	ownerName = "owner"
	labelKey  = "app"
	svcName   = "svc"
	smName    = "sm"
)

// testGVK 는 CRD 의존 unstructured 리소스 (ServiceMonitor 류) 의 테스트 GVK.
var testGVK = schema.GroupVersionKind{Group: "monitoring.example.com", Version: "v1", Kind: "ServiceMonitor"}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{
		corev1.AddToScheme,
		appsv1.AddToScheme,
		networkingv1.AddToScheme,
		policyv1.AddToScheme,
		autoscalingv2.AddToScheme,
	} {
		if err := add(s); err != nil {
			t.Fatalf("scheme 등록 실패: %v", err)
		}
	}
	s.AddKnownTypeWithName(testGVK, &unstructured.Unstructured{})
	s.AddKnownTypeWithName(testGVK.GroupVersion().WithKind(testGVK.Kind+"List"), &unstructured.UnstructuredList{})
	return s
}

// newOwner 는 owner reference 검증용 더미 owner 객체.
func newOwner() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: ownerName, Namespace: testNS, UID: "owner-uid"},
	}
}

// seedTime 은 update 경로 판별 (CreationTimestamp.IsZero) 을 위해 seed 객체에
// 부여하는 고정 생성 시각.
func seedTime() metav1.Time { return metav1.NewTime(time.Unix(1000, 0).UTC()) }

func assertOwnerRef(t *testing.T, obj metav1.Object) {
	t.Helper()
	refs := obj.GetOwnerReferences()
	if len(refs) != 1 || refs[0].Name != ownerName {
		t.Fatalf("owner reference 1건 (owner=%s) 기대, got %+v", ownerName, refs)
	}
	if refs[0].Controller == nil || !*refs[0].Controller {
		t.Fatalf("controller=true 기대, got %+v", refs[0])
	}
}

func TestConfigMap(t *testing.T) {
	const cmName = "cm"
	newCM := func(data map[string]string) *corev1.ConfigMap {
		return &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: cmName, Namespace: testNS},
			Data:       data,
		}
	}
	cases := []struct {
		name     string
		existing *corev1.ConfigMap
		desired  *corev1.ConfigMap
		wantData map[string]string
	}{
		{
			name:     "신규 생성 시 Data 적용 + ownerRef 설정",
			desired:  newCM(map[string]string{"k": "v"}),
			wantData: map[string]string{"k": "v"},
		},
		{
			name: "기존 객체 갱신 시 Data 동기화",
			existing: func() *corev1.ConfigMap {
				cm := newCM(map[string]string{"k": "stale"})
				cm.CreationTimestamp = seedTime()
				return cm
			}(),
			desired:  newCM(map[string]string{"k": "fresh", "extra": "1"}),
			wantData: map[string]string{"k": "fresh", "extra": "1"},
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

			if err := ConfigMap(t.Context(), c, s, newOwner(), tc.desired); err != nil {
				t.Fatalf("ConfigMap apply 실패: %v", err)
			}

			got := &corev1.ConfigMap{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: cmName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			if !maps.Equal(got.Data, tc.wantData) {
				t.Errorf("Data = %v, 기대 %v", got.Data, tc.wantData)
			}
			assertOwnerRef(t, got)
		})
	}
}

func baseSvc() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: svcName, Namespace: testNS},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: map[string]string{labelKey: "x"},
			Ports:    []corev1.ServicePort{{Name: "p", Port: 80}},
		},
	}
}

func seededSvc(mut func(*corev1.Service)) *corev1.Service {
	s := baseSvc()
	s.CreationTimestamp = seedTime()
	if mut != nil {
		mut(s)
	}
	return s
}

func desiredSvc(mut func(*corev1.Service)) *corev1.Service {
	s := baseSvc()
	if mut != nil {
		mut(s)
	}
	return s
}

func TestService(t *testing.T) {
	lbClass := "service.k8s.aws/nlb"
	cases := []struct {
		name     string
		existing *corev1.Service
		desired  *corev1.Service
		want     func(t *testing.T, got *corev1.Service)
	}{
		{
			name: "신규 생성 시 create-only 필드 desired 적용",
			desired: desiredSvc(func(s *corev1.Service) {
				s.Spec.ClusterIP = corev1.ClusterIPNone
				s.Spec.IPFamilies = []corev1.IPFamily{corev1.IPv4Protocol}
				s.Spec.IPFamilyPolicy = ptr.To(corev1.IPFamilyPolicySingleStack)
			}),
			want: func(t *testing.T, got *corev1.Service) {
				if got.Spec.ClusterIP != corev1.ClusterIPNone {
					t.Errorf("ClusterIP = %q, 기대 None", got.Spec.ClusterIP)
				}
				if len(got.Spec.IPFamilies) != 1 || got.Spec.IPFamilies[0] != corev1.IPv4Protocol {
					t.Errorf("IPFamilies = %v, 기대 [IPv4]", got.Spec.IPFamilies)
				}
				assertOwnerRef(t, got)
			},
		},
		{
			name: "갱신 시 immutable 필드 (ClusterIP/IPFamilies/IPFamilyPolicy) 보존 + Ports 동기화",
			existing: seededSvc(func(s *corev1.Service) {
				s.Spec.ClusterIP = "10.0.0.1"
				s.Spec.IPFamilies = []corev1.IPFamily{corev1.IPv4Protocol}
				s.Spec.IPFamilyPolicy = ptr.To(corev1.IPFamilyPolicySingleStack)
			}),
			desired: desiredSvc(func(s *corev1.Service) {
				s.Spec.ClusterIP = corev1.ClusterIPNone
				s.Spec.IPFamilies = []corev1.IPFamily{corev1.IPv6Protocol}
				s.Spec.IPFamilyPolicy = ptr.To(corev1.IPFamilyPolicyRequireDualStack)
				s.Spec.Ports = []corev1.ServicePort{{Name: "p", Port: 81}}
			}),
			want: func(t *testing.T, got *corev1.Service) {
				if got.Spec.ClusterIP != "10.0.0.1" {
					t.Errorf("ClusterIP = %q, 기존 10.0.0.1 보존 기대", got.Spec.ClusterIP)
				}
				if len(got.Spec.IPFamilies) != 1 || got.Spec.IPFamilies[0] != corev1.IPv4Protocol {
					t.Errorf("IPFamilies = %v, 기존 [IPv4] 보존 기대", got.Spec.IPFamilies)
				}
				if got.Spec.IPFamilyPolicy == nil || *got.Spec.IPFamilyPolicy != corev1.IPFamilyPolicySingleStack {
					t.Errorf("IPFamilyPolicy = %v, 기존 SingleStack 보존 기대", got.Spec.IPFamilyPolicy)
				}
				if got.Spec.Ports[0].Port != 81 {
					t.Errorf("Ports[0].Port = %d, 81 동기화 기대", got.Spec.Ports[0].Port)
				}
			},
		},
		{
			name: "갱신 시 LB 운영 필드 동기화 (LoadBalancerIP 고착 P0 회귀 가드)",
			existing: seededSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
				s.Spec.LoadBalancerIP = "1.1.1.1"
				s.Spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyCluster
			}),
			desired: desiredSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
				s.Spec.LoadBalancerIP = "2.2.2.2"
				s.Spec.ExternalTrafficPolicy = corev1.ServiceExternalTrafficPolicyLocal
				s.Spec.LoadBalancerSourceRanges = []string{"10.0.0.0/8"}
			}),
			want: func(t *testing.T, got *corev1.Service) {
				if got.Spec.LoadBalancerIP != "2.2.2.2" {
					t.Errorf("LoadBalancerIP = %q, 2.2.2.2 동기화 기대 (고착 P0)", got.Spec.LoadBalancerIP)
				}
				if got.Spec.ExternalTrafficPolicy != corev1.ServiceExternalTrafficPolicyLocal {
					t.Errorf("ExternalTrafficPolicy = %q, Local 기대", got.Spec.ExternalTrafficPolicy)
				}
				if len(got.Spec.LoadBalancerSourceRanges) != 1 {
					t.Errorf("LoadBalancerSourceRanges = %v", got.Spec.LoadBalancerSourceRanges)
				}
			},
		},
		{
			name: "갱신 시 pointer 필드 desired nil 이면 server-default 보존 (ping-pong 차단)",
			existing: seededSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
				s.Spec.AllocateLoadBalancerNodePorts = new(true)
				s.Spec.LoadBalancerClass = &lbClass
				s.Spec.InternalTrafficPolicy = ptr.To(corev1.ServiceInternalTrafficPolicyCluster)
			}),
			desired: desiredSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
			}),
			want: func(t *testing.T, got *corev1.Service) {
				if got.Spec.AllocateLoadBalancerNodePorts == nil || !*got.Spec.AllocateLoadBalancerNodePorts {
					t.Errorf("AllocateLoadBalancerNodePorts = %v, 기존 true 보존 기대", got.Spec.AllocateLoadBalancerNodePorts)
				}
				if got.Spec.LoadBalancerClass == nil || *got.Spec.LoadBalancerClass != lbClass {
					t.Errorf("LoadBalancerClass = %v, 기존 보존 기대", got.Spec.LoadBalancerClass)
				}
				if got.Spec.InternalTrafficPolicy == nil ||
					*got.Spec.InternalTrafficPolicy != corev1.ServiceInternalTrafficPolicyCluster {
					t.Errorf("InternalTrafficPolicy = %v, 기존 Cluster 보존 기대", got.Spec.InternalTrafficPolicy)
				}
			},
		},
		{
			name: "갱신 시 pointer 필드 desired 명시면 동기화",
			existing: seededSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
				s.Spec.AllocateLoadBalancerNodePorts = new(true)
				s.Spec.InternalTrafficPolicy = ptr.To(corev1.ServiceInternalTrafficPolicyCluster)
			}),
			desired: desiredSvc(func(s *corev1.Service) {
				s.Spec.Type = corev1.ServiceTypeLoadBalancer
				s.Spec.AllocateLoadBalancerNodePorts = new(false)
				s.Spec.LoadBalancerClass = &lbClass
				s.Spec.InternalTrafficPolicy = ptr.To(corev1.ServiceInternalTrafficPolicyLocal)
			}),
			want: func(t *testing.T, got *corev1.Service) {
				if got.Spec.AllocateLoadBalancerNodePorts == nil || *got.Spec.AllocateLoadBalancerNodePorts {
					t.Errorf("AllocateLoadBalancerNodePorts = %v, false 동기화 기대", got.Spec.AllocateLoadBalancerNodePorts)
				}
				if got.Spec.LoadBalancerClass == nil || *got.Spec.LoadBalancerClass != lbClass {
					t.Errorf("LoadBalancerClass = %v, %q 동기화 기대", got.Spec.LoadBalancerClass, lbClass)
				}
				if got.Spec.InternalTrafficPolicy == nil ||
					*got.Spec.InternalTrafficPolicy != corev1.ServiceInternalTrafficPolicyLocal {
					t.Errorf("InternalTrafficPolicy = %v, Local 동기화 기대", got.Spec.InternalTrafficPolicy)
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

			if err := Service(t.Context(), c, s, newOwner(), tc.desired); err != nil {
				t.Fatalf("Service apply 실패: %v", err)
			}

			got := &corev1.Service{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: svcName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			tc.want(t, got)
		})
	}
}

func TestNetworkPolicy(t *testing.T) {
	const npName = "np"
	newNP := func(selVal string, ingressLen int) *networkingv1.NetworkPolicy {
		np := &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{Name: npName, Namespace: testNS},
			Spec: networkingv1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{MatchLabels: map[string]string{labelKey: selVal}},
				PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
			},
		}
		for range ingressLen {
			np.Spec.Ingress = append(np.Spec.Ingress, networkingv1.NetworkPolicyIngressRule{})
		}
		return np
	}
	cases := []struct {
		name       string
		existing   *networkingv1.NetworkPolicy
		desired    *networkingv1.NetworkPolicy
		wantSelVal string
		wantRules  int
	}{
		{
			name:       "신규 생성",
			desired:    newNP("v1", 1),
			wantSelVal: "v1",
			wantRules:  1,
		},
		{
			name: "갱신 시 PodSelector/Ingress 매번 동기화",
			existing: func() *networkingv1.NetworkPolicy {
				np := newNP("v1", 1)
				np.CreationTimestamp = seedTime()
				return np
			}(),
			desired:    newNP("v2", 2),
			wantSelVal: "v2",
			wantRules:  2,
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

			if err := NetworkPolicy(t.Context(), c, s, newOwner(), tc.desired); err != nil {
				t.Fatalf("NetworkPolicy apply 실패: %v", err)
			}

			got := &networkingv1.NetworkPolicy{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: npName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			if got.Spec.PodSelector.MatchLabels[labelKey] != tc.wantSelVal {
				t.Errorf("PodSelector = %v, 기대 %s", got.Spec.PodSelector.MatchLabels, tc.wantSelVal)
			}
			if len(got.Spec.Ingress) != tc.wantRules {
				t.Errorf("Ingress rule 수 = %d, 기대 %d", len(got.Spec.Ingress), tc.wantRules)
			}
			assertOwnerRef(t, got)
		})
	}
}

func TestPDB(t *testing.T) {
	const pdbName = "pdb"
	newPDB := func(minAvailable int32) *policyv1.PodDisruptionBudget {
		return &policyv1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{Name: pdbName, Namespace: testNS},
			Spec: policyv1.PodDisruptionBudgetSpec{
				Selector:     &metav1.LabelSelector{MatchLabels: map[string]string{labelKey: "x"}},
				MinAvailable: new(intstr.FromInt32(minAvailable)),
			},
		}
	}
	cases := []struct {
		name     string
		existing *policyv1.PodDisruptionBudget
		desired  *policyv1.PodDisruptionBudget
		wantMin  int32
	}{
		{name: "신규 생성", desired: newPDB(1), wantMin: 1},
		{
			name: "갱신 시 MinAvailable 동기화",
			existing: func() *policyv1.PodDisruptionBudget {
				pdb := newPDB(1)
				pdb.CreationTimestamp = seedTime()
				return pdb
			}(),
			desired: newPDB(2),
			wantMin: 2,
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

			if err := PDB(t.Context(), c, s, newOwner(), tc.desired); err != nil {
				t.Fatalf("PDB apply 실패: %v", err)
			}

			got := &policyv1.PodDisruptionBudget{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: pdbName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			if got.Spec.MinAvailable == nil || got.Spec.MinAvailable.IntValue() != int(tc.wantMin) {
				t.Errorf("MinAvailable = %v, 기대 %d", got.Spec.MinAvailable, tc.wantMin)
			}
			assertOwnerRef(t, got)
		})
	}
}

func TestHPA(t *testing.T) {
	const hpaName = "hpa"
	newHPA := func(maxReplicas int32) *autoscalingv2.HorizontalPodAutoscaler {
		return &autoscalingv2.HorizontalPodAutoscaler{
			ObjectMeta: metav1.ObjectMeta{Name: hpaName, Namespace: testNS},
			Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
				ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
					APIVersion: "apps/v1", Kind: "StatefulSet", Name: "web",
				},
				MaxReplicas: maxReplicas,
			},
		}
	}
	cases := []struct {
		name       string
		existing   *autoscalingv2.HorizontalPodAutoscaler
		desired    *autoscalingv2.HorizontalPodAutoscaler
		wantAbsent bool
		wantMax    int32
	}{
		{name: "desired nil + 부재 → no-op", wantAbsent: true},
		{
			name: "desired nil + 존재 → 삭제 (toggle off 회복)",
			existing: func() *autoscalingv2.HorizontalPodAutoscaler {
				h := newHPA(3)
				h.CreationTimestamp = seedTime()
				return h
			}(),
			wantAbsent: true,
		},
		{name: "신규 생성", desired: newHPA(5), wantMax: 5},
		{
			name: "갱신 시 spec 동기화",
			existing: func() *autoscalingv2.HorizontalPodAutoscaler {
				h := newHPA(3)
				h.CreationTimestamp = seedTime()
				return h
			}(),
			desired: newHPA(7),
			wantMax: 7,
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

			if err := HPA(t.Context(), c, s, newOwner(), hpaName, testNS, tc.desired); err != nil {
				t.Fatalf("HPA apply 실패: %v", err)
			}

			got := &autoscalingv2.HorizontalPodAutoscaler{}
			err := c.Get(t.Context(), client.ObjectKey{Name: hpaName, Namespace: testNS}, got)
			if tc.wantAbsent {
				if !apierrors.IsNotFound(err) {
					t.Fatalf("HPA 부재 기대, got err=%v obj=%+v", err, got.Spec)
				}
				return
			}
			if err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			if got.Spec.MaxReplicas != tc.wantMax {
				t.Errorf("MaxReplicas = %d, 기대 %d", got.Spec.MaxReplicas, tc.wantMax)
			}
			assertOwnerRef(t, got)
		})
	}
}

func newSM(specA string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{Object: map[string]any{
		"spec": map[string]any{"a": specA},
	}}
	u.SetGroupVersionKind(testGVK)
	u.SetName(smName)
	u.SetNamespace(testNS)
	return u
}

func TestUnstructured(t *testing.T) {
	noMatchErr := &meta.NoKindMatchError{
		GroupKind:        testGVK.GroupKind(),
		SearchedVersions: []string{testGVK.Version},
	}
	cases := []struct {
		name     string
		existing *unstructured.Unstructured
		desired  *unstructured.Unstructured
		noMatch  bool // Get 이 NoKindMatchError 반환 (CRD 미설치 모사)
		failSoft bool
		wantErr  bool
		wantA    string // "" 이면 객체 검증 skip
	}{
		{name: "desired nil 은 no-op", desired: nil},
		{name: "신규 생성 - spec + ownerRef", desired: newSM("v1"), wantA: "v1"},
		{name: "기존 갱신 - spec 교체", existing: newSM("stale"), desired: newSM("fresh"), wantA: "fresh"},
		{name: "CRD 미설치 + failSoft=true → nil", desired: newSM("v1"), noMatch: true, failSoft: true},
		{name: "CRD 미설치 + failSoft=false → 에러 전파", desired: newSM("v1"), noMatch: true, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := testScheme(t)
			b := fake.NewClientBuilder().WithScheme(s)
			if tc.existing != nil {
				b = b.WithObjects(tc.existing)
			}
			if tc.noMatch {
				b = b.WithInterceptorFuncs(interceptor.Funcs{
					Get: func(_ context.Context, _ client.WithWatch, _ client.ObjectKey,
						_ client.Object, _ ...client.GetOption) error {
						return noMatchErr
					},
				})
			}
			c := b.Build()

			err := Unstructured(t.Context(), c, s, newOwner(), tc.desired, tc.failSoft)
			if tc.wantErr {
				if err == nil || !meta.IsNoMatchError(err) {
					t.Fatalf("NoMatch 에러 전파 기대, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unstructured apply 실패: %v", err)
			}
			if tc.wantA == "" {
				return
			}

			got := &unstructured.Unstructured{}
			got.SetGroupVersionKind(testGVK)
			if err := c.Get(t.Context(), client.ObjectKey{Name: smName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			a, _, _ := unstructured.NestedString(got.Object, "spec", "a")
			if a != tc.wantA {
				t.Errorf("spec.a = %q, 기대 %q", a, tc.wantA)
			}
			assertOwnerRef(t, got)
		})
	}
}

func TestSecretIfNotExists(t *testing.T) {
	const secName = "sec"
	cases := []struct {
		name       string
		existing   *corev1.Secret
		raceCreate bool // Create 가 AlreadyExists 반환 (병렬 reconcile race 모사)
		wantBuilds int
		wantValue  string // "" 이면 객체 내용 검증 skip
	}{
		{name: "부재 시 build 호출 + 생성 + ownerRef", wantBuilds: 1, wantValue: "generated"},
		{
			name: "존재 시 no-op (build 미호출 — credential 재생성 차단)",
			existing: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: secName, Namespace: testNS},
				Data:       map[string][]byte{"k": []byte("kept")},
			},
			wantBuilds: 0,
			wantValue:  "kept",
		},
		{name: "Create AlreadyExists race 흡수 (에러 아님)", raceCreate: true, wantBuilds: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := testScheme(t)
			b := fake.NewClientBuilder().WithScheme(s)
			if tc.existing != nil {
				b = b.WithObjects(tc.existing)
			}
			if tc.raceCreate {
				b = b.WithInterceptorFuncs(interceptor.Funcs{
					Create: func(_ context.Context, _ client.WithWatch, obj client.Object,
						_ ...client.CreateOption) error {
						return apierrors.NewAlreadyExists(
							schema.GroupResource{Resource: "secrets"}, obj.GetName())
					},
				})
			}
			c := b.Build()

			builds := 0
			build := func() *corev1.Secret {
				builds++
				return &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: secName, Namespace: testNS},
					Data:       map[string][]byte{"k": []byte("generated")},
				}
			}

			if err := SecretIfNotExists(t.Context(), c, s, newOwner(), secName, build); err != nil {
				t.Fatalf("SecretIfNotExists 실패: %v", err)
			}
			if builds != tc.wantBuilds {
				t.Errorf("build 호출 수 = %d, 기대 %d", builds, tc.wantBuilds)
			}
			if tc.wantValue == "" {
				return
			}

			got := &corev1.Secret{}
			if err := c.Get(t.Context(), client.ObjectKey{Name: secName, Namespace: testNS}, got); err != nil {
				t.Fatalf("get 실패: %v", err)
			}
			if string(got.Data["k"]) != tc.wantValue {
				t.Errorf("Data[k] = %q, 기대 %q", got.Data["k"], tc.wantValue)
			}
			if tc.wantBuilds > 0 {
				assertOwnerRef(t, got)
			}
		})
	}
}
