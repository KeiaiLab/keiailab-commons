// SPDX-License-Identifier: MIT

package pdb_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/keiailab/keiailab-commons/pkg/pdb"
)

func ptr(v intstr.IntOrString) *intstr.IntOrString { return &v }

func TestBuild_ObjectMetaAndSelector(t *testing.T) {
	t.Parallel()
	got := pdb.Build(pdb.Params{
		Name:      "cache-0-pdb",
		Namespace: "ns1",
		Labels:    map[string]string{"app.kubernetes.io/name": "valkey"},
		Selector:  map[string]string{"app.kubernetes.io/instance": "cache-0"},
		Replicas:  3,
	})
	if got.Name != "cache-0-pdb" || got.Namespace != "ns1" {
		t.Fatalf("ObjectMeta mismatch: %s/%s", got.Namespace, got.Name)
	}
	if got.Labels["app.kubernetes.io/name"] != "valkey" {
		t.Fatalf("labels not propagated: %v", got.Labels)
	}
	if got.Spec.Selector == nil || got.Spec.Selector.MatchLabels["app.kubernetes.io/instance"] != "cache-0" {
		t.Fatalf("selector mismatch: %+v", got.Spec.Selector)
	}
}

func TestBuild_MinAvailableTakesPrecedence(t *testing.T) {
	t.Parallel()
	mn := ptr(intstr.FromInt(2))
	mx := ptr(intstr.FromInt(1))
	got := pdb.Build(pdb.Params{Name: "x", Replicas: 5, MinAvailable: mn, MaxUnavailable: mx})
	if got.Spec.MinAvailable == nil || got.Spec.MinAvailable.IntValue() != 2 {
		t.Fatalf("MinAvailable not honored: %+v", got.Spec.MinAvailable)
	}
	if got.Spec.MaxUnavailable != nil {
		t.Fatalf("MaxUnavailable must be nil when MinAvailable set: %+v", got.Spec.MaxUnavailable)
	}
}

func TestBuild_MaxUnavailableWhenNoMin(t *testing.T) {
	t.Parallel()
	mx := ptr(intstr.FromString("25%"))
	got := pdb.Build(pdb.Params{Name: "x", Replicas: 5, MaxUnavailable: mx})
	if got.Spec.MaxUnavailable == nil || got.Spec.MaxUnavailable.StrVal != "25%" {
		t.Fatalf("MaxUnavailable not honored: %+v", got.Spec.MaxUnavailable)
	}
	if got.Spec.MinAvailable != nil {
		t.Fatalf("MinAvailable must be nil: %+v", got.Spec.MinAvailable)
	}
}

// 기본 정책: DefaultFloor 로 operator 별 차이를 흡수한다.
func TestBuild_DefaultMinAvailable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		replicas int32
		floor    int
		want     int
	}{
		{"valkey 3노드 floor1", 3, 1, 2}, // max(2,1)=2
		{"valkey 1노드 floor1", 1, 1, 1}, // max(0,1)=1 — primary 보존
		{"mongo 3멤버 floor0", 3, 0, 2},  // max(2,0)=2
		{"mongo 1멤버 floor0", 1, 0, 0},  // max(0,0)=0
		{"replicas0 floor0", 0, 0, 0},  // max(-1,0)=0
		{"replicas0 floor1", 0, 1, 1},  // max(-1,1)=1
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := pdb.Build(pdb.Params{Name: "x", Replicas: c.replicas, DefaultFloor: c.floor})
			if got.Spec.MinAvailable == nil {
				t.Fatalf("MinAvailable nil")
			}
			if got.Spec.MinAvailable.IntValue() != c.want {
				t.Fatalf("minAvailable=%d, want %d", got.Spec.MinAvailable.IntValue(), c.want)
			}
			if got.Spec.MaxUnavailable != nil {
				t.Fatalf("MaxUnavailable must be nil in default policy")
			}
		})
	}
}

// 회귀 가드: commons Build 가 valkey BuildPDB 와 동일 결과(default floor=1).
func TestBuild_MatchesValkeyDefault(t *testing.T) {
	t.Parallel()
	got := pdb.Build(pdb.Params{
		Name: "v-pdb", Namespace: "ns", Replicas: 3, DefaultFloor: 1,
		Labels:   map[string]string{"k": "v"},
		Selector: map[string]string{"s": "1"},
	})
	want := intstr.FromInt(2)
	if *got.Spec.MinAvailable != want {
		t.Fatalf("valkey 등가성 실패: %v != %v", *got.Spec.MinAvailable, want)
	}
	_ = metav1.ObjectMeta{} // import 사용 보장
}
