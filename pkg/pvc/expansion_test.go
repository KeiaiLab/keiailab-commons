// SPDX-License-Identifier: Apache-2.0

package pvc

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testNS   = "ns"
	size10Gi = "10Gi"
	size20Gi = "20Gi"
)

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	if err := corev1.AddToScheme(s); err != nil {
		t.Fatalf("corev1: %v", err)
	}
	if err := storagev1.AddToScheme(s); err != nil {
		t.Fatalf("storagev1: %v", err)
	}
	return s
}

func boundPVC(name, scName string) *corev1.PersistentVolumeClaim {
	q := resource.MustParse(size10Gi)
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNS},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: q},
			},
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound},
	}
}

func pendingPVC(name, scName string) *corev1.PersistentVolumeClaim {
	q := resource.MustParse(size10Gi)
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNS},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: q},
			},
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimPending},
	}
}

func sc(name string, allowExpansion bool) *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta:           metav1.ObjectMeta{Name: name},
		Provisioner:          "test",
		AllowVolumeExpansion: &allowExpansion,
	}
}

func sizeOf(t *testing.T, c client.Client, name string) string {
	t.Helper()
	p := &corev1.PersistentVolumeClaim{}
	if err := c.Get(context.Background(), types.NamespacedName{Namespace: testNS, Name: name}, p); err != nil {
		t.Fatalf("get %s: %v", name, err)
	}
	q := p.Spec.Resources.Requests[corev1.ResourceStorage]
	return q.String()
}

func TestExpandDataPVCs_grows_when_SC_allows(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
		boundPVC("data-pg-shard-0-1", "gp3"),
	).Build()

	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, n := range []string{"data-pg-shard-0-0", "data-pg-shard-0-1"} {
		if got := sizeOf(t, c, n); got != size20Gi {
			t.Errorf("%s: %s want %s", n, got, size20Gi)
		}
	}
}

func TestExpandDataPVCs_skips_non_matching_PVCs(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
		boundPVC("data-other-cluster-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size20Gi {
		t.Errorf("matching PVC: got %s want %s", got, size20Gi)
	}
	if got := sizeOf(t, c, "data-other-cluster-0"); got != size10Gi {
		t.Errorf("other cluster PVC 변경되면 안됨: got %s", got)
	}
}

func TestExpandDataPVCs_no_expansion_when_SC_disallows(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("standard", false),
		boundPVC("data-pg-shard-0-0", "standard"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("disallow expansion → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_zero_desired_noop(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.Quantity{}); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("zero desired → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_empty_stsNames_noop(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		nil, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("empty stsNames → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_multi_shard(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
		boundPVC("data-pg-shard-1-0", "gp3"),
		boundPVC("data-pg-shard-2-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0", "pg-shard-1", "pg-shard-2"},
		resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, n := range []string{"data-pg-shard-0-0", "data-pg-shard-1-0", "data-pg-shard-2-0"} {
		if got := sizeOf(t, c, n); got != size20Gi {
			t.Errorf("%s: %s want %s", n, got, size20Gi)
		}
	}
}

func TestExpandDataPVCs_skips_pending_PVC(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		pendingPVC("data-pg-shard-0-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("Pending PVC → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_skips_when_desired_le_current(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true),
		boundPVC("data-pg-shard-0-0", "gp3"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse("5Gi")); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("감소 시도 → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_skips_when_SC_not_found(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		boundPVC("data-pg-shard-0-0", "missing"),
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"pg-shard-0"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-pg-shard-0-0"); got != size10Gi {
		t.Errorf("SC missing → unchanged, got %s", got)
	}
}

func TestExpandDataPVCs_custom_VCT_name(t *testing.T) {
	q := resource.MustParse(size10Gi)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "store-myapp-0", Namespace: testNS},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: ptrStr("gp3"),
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: q},
			},
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound},
	}
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(
		sc("gp3", true), pvc,
	).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"myapp"}, resource.MustParse(size20Gi),
		WithVCTName("store")); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "store-myapp-0"); got != size20Gi {
		t.Errorf("custom VCT name: got %s want %s", got, size20Gi)
	}
}

func TestExpandDataPVCs_empty_VCT_name_returns_error(t *testing.T) {
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).Build()
	err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"x"}, resource.MustParse(size20Gi),
		WithVCTName(""))
	if err == nil {
		t.Fatal("expected ErrEmptyVCTName, got nil")
	}
	if err != ErrEmptyVCTName {
		t.Errorf("error: %v, want %v", err, ErrEmptyVCTName)
	}
}

func TestExpandDataPVCs_PVC_without_SC_skips_SC_check(t *testing.T) {
	q := resource.MustParse(size10Gi)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "data-x-0", Namespace: testNS},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: q},
			},
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound},
	}
	c := fake.NewClientBuilder().WithScheme(testScheme(t)).WithObjects(pvc).Build()
	if err := ExpandDataPVCs(context.Background(), c, testNS,
		[]string{"x"}, resource.MustParse(size20Gi)); err != nil {
		t.Fatalf("err: %v", err)
	}
	if got := sizeOf(t, c, "data-x-0"); got != size20Gi {
		t.Errorf("PVC without SC: got %s want %s", got, size20Gi)
	}
}

func TestPVCNamePrefix(t *testing.T) {
	cases := []struct {
		vct, sts string
		want     string
	}{
		{"data", "mycluster", "data-mycluster-"},
		{"", "mycluster", "data-mycluster-"},
		{"store", "x", "store-x-"},
	}
	for _, tc := range cases {
		got := PVCNamePrefix(tc.vct, tc.sts)
		if got != tc.want {
			t.Errorf("PVCNamePrefix(%q,%q) = %q, want %q", tc.vct, tc.sts, got, tc.want)
		}
	}
}

func ptrStr(s string) *string { return &s }
