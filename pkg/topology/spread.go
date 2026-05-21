package topology

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 기본 임계값 / topology key — 3 operator 의 통합 default.
const (
	// DefaultMinReplicas 는 default TSC 를 자동 주입하기 위한 *최소 replica 수*
	// 의 기본값이다. 본 값 미만 시 nil 반환 (single pod 환경 spread 무의미).
	//
	// mongodb / valkey 가 이 값을 사용. postgres 는 WithMinReplicas(1) override.
	DefaultMinReplicas int32 = 2

	// TopologyKeyZone — 표준 K8s zone label key.
	TopologyKeyZone = "topology.kubernetes.io/zone"

	// TopologyKeyHostname — 표준 K8s hostname label key.
	TopologyKeyHostname = "kubernetes.io/hostname"

	// DefaultMaxSkew — 동일 zone/node 에 +1 이상 몰리는 것 방지.
	DefaultMaxSkew int32 = 1
)

// DefaultTopologyKeys 는 default TSC 의 표준 topology key 순서다.
//
// 순서는 의도적 — zone 먼저, hostname 두번째 (broader → narrower).
func DefaultTopologyKeys() []string {
	return []string{TopologyKeyZone, TopologyKeyHostname}
}

// Option 은 Defaulted 동작을 변경하는 함수형 옵션.
type Option func(*config)

type config struct {
	minReplicas       int32
	topologyKeys      []string
	maxSkew           int32
	whenUnsatisfiable corev1.UnsatisfiableConstraintAction
}

func newConfig(opts ...Option) *config {
	c := &config{
		minReplicas:       DefaultMinReplicas,
		topologyKeys:      DefaultTopologyKeys(),
		maxSkew:           DefaultMaxSkew,
		whenUnsatisfiable: corev1.ScheduleAnyway,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithMinReplicas 는 default TSC 를 주입하는 최소 replica 임계값을 변경한다.
// 기본값은 DefaultMinReplicas (2). downstream operator 처럼 replicas 가 "추가
// 복제본 수" 의미인 경우 WithMinReplicas(1) 사용.
//
// replicas < min → 주입 skip (nil 반환).
func WithMinReplicas(min int32) Option {
	return func(c *config) { c.minReplicas = min }
}

// WithTopologyKeys 는 default TSC 의 topology key 슬라이스를 교체한다.
// 기본값은 DefaultTopologyKeys() (zone + hostname). 1개 또는 3개+ 도 가능.
// nil 또는 빈 슬라이스 전달 시 default 유지 (no-op).
func WithTopologyKeys(keys ...string) Option {
	return func(c *config) {
		if len(keys) == 0 {
			return
		}
		c.topologyKeys = keys
	}
}

// WithMaxSkew 는 default TSC 의 MaxSkew 를 변경한다. 기본값 1.
// 0 이하 입력은 무시 (K8s validation 에서 reject 됨).
func WithMaxSkew(skew int32) Option {
	return func(c *config) {
		if skew <= 0 {
			return
		}
		c.maxSkew = skew
	}
}

// WithWhenUnsatisfiable 는 default TSC 의 WhenUnsatisfiable 액션을 변경한다.
// 기본값 ScheduleAnyway (single-zone cluster 환경 호환).
// DoNotSchedule 로 변경 시 강제 spread — single-zone 환경에서는 pending 위험.
func WithWhenUnsatisfiable(a corev1.UnsatisfiableConstraintAction) Option {
	return func(c *config) {
		if a == "" {
			return
		}
		c.whenUnsatisfiable = a
	}
}

// Defaulted 는 user-provided TSC 가 비어있고 replicas >= minReplicas 일 때
// HA out-of-box default TSC (zone + hostname spread) 를 자동 주입한다.
//
// 동작 우선순위:
//  1. len(user) > 0 → user 그대로 반환 (사용자 명시 override 보장).
//  2. replicas < minReplicas → nil 반환 (단일 pod 환경 → spread 무의미).
//  3. 그 외 → topologyKeys 의 각 키에 대해 MaxSkew=1 + ScheduleAnyway TSC 주입.
//
// 사용 예 (downstream operator):
//
//	tsc := topology.Defaulted(
//	    cluster.Spec.Shards.TopologySpreadConstraints,
//	    cluster.Spec.Shards.Replicas,
//	    labels,
//	    topology.WithMinReplicas(1), // postgres 는 "additional copies" 의미.
//	)
//
// 사용 예 (downstream operator):
//
//	tsc := topology.Defaulted(
//	    nil, // mongodb 는 user TSC 사용 안 함.
//	    mdb.Spec.Members,
//	    labels,
//	)
//
// 사용 예 (downstream operator):
//
//	tsc := topology.Defaulted(
//	    pod.TopologySpreadConstraints,
//	    p.Replicas,
//	    selector,
//	)
//
// selector 가 nil 이거나 빈 map 인 경우 LabelSelector 는 *모든 pod* 매칭이
// 되므로 비정상 동작 — 호출자가 항상 non-nil 의미있는 selector 전달.
func Defaulted(
	user []corev1.TopologySpreadConstraint,
	replicas int32,
	selector map[string]string,
	opts ...Option,
) []corev1.TopologySpreadConstraint {
	if len(user) > 0 {
		return user
	}
	cfg := newConfig(opts...)
	if replicas < cfg.minReplicas {
		return nil
	}
	labelSelector := &metav1.LabelSelector{MatchLabels: selector}
	out := make([]corev1.TopologySpreadConstraint, 0, len(cfg.topologyKeys))
	for _, k := range cfg.topologyKeys {
		out = append(out, corev1.TopologySpreadConstraint{
			MaxSkew:           cfg.maxSkew,
			TopologyKey:       k,
			WhenUnsatisfiable: cfg.whenUnsatisfiable,
			LabelSelector:     labelSelector,
		})
	}
	return out
}
