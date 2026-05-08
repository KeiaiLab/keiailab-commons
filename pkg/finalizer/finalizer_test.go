package finalizer_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keiailab/operator-commons/pkg/finalizer"
)

// fakeObj — metav1.Object 인터페이스 minimal 구현 (테스트 헬퍼).
//
// 실제 K8s controller 코드는 client.Object (= runtime.Object + metav1.Object)
// 를 받지만 본 패키지는 metav1.Object 만 요구하므로 ObjectMeta 직접 사용.
type fakeObj struct {
	metav1.ObjectMeta
}

// TestAdd_AddsAndIdempotent — Add 가 부재 시 추가하고, 재호출 시 no-op.
func TestAdd_AddsAndIdempotent(t *testing.T) {
	obj := &fakeObj{}
	const name = "test.keiailab.com/finalizer"

	if changed := finalizer.Add(&obj.ObjectMeta, name); !changed {
		t.Fatal("expected first Add to return true")
	}
	if got := obj.GetFinalizers(); len(got) != 1 || got[0] != name {
		t.Errorf("expected finalizers=[%q], got %v", name, got)
	}

	if changed := finalizer.Add(&obj.ObjectMeta, name); changed {
		t.Error("expected idempotent second Add to return false")
	}
	if got := obj.GetFinalizers(); len(got) != 1 {
		t.Errorf("expected len=1 after idempotent add, got %d", len(got))
	}
}

// TestRemove_RemovesAndIdempotent.
func TestRemove_RemovesAndIdempotent(t *testing.T) {
	obj := &fakeObj{}
	obj.SetFinalizers([]string{"a", "b", "c"})

	if changed := finalizer.Remove(&obj.ObjectMeta, "b"); !changed {
		t.Error("expected Remove(b) to return true")
	}
	got := obj.GetFinalizers()
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Errorf("expected [a c], got %v", got)
	}

	if changed := finalizer.Remove(&obj.ObjectMeta, "b"); changed {
		t.Error("expected idempotent second Remove to return false")
	}
}

// TestHas — 등록 / 부재 확인.
func TestHas(t *testing.T) {
	obj := &fakeObj{}
	obj.SetFinalizers([]string{"x", "y"})

	if !finalizer.Has(&obj.ObjectMeta, "x") {
		t.Error("expected Has(x)=true")
	}
	if finalizer.Has(&obj.ObjectMeta, "z") {
		t.Error("expected Has(z)=false")
	}
}

// TestAdd_PreservesExisting — 기존 finalizer 유지하면서 추가.
func TestAdd_PreservesExisting(t *testing.T) {
	obj := &fakeObj{}
	obj.SetFinalizers([]string{"existing"})

	finalizer.Add(&obj.ObjectMeta, "new")
	got := obj.GetFinalizers()
	if len(got) != 2 || got[0] != "existing" || got[1] != "new" {
		t.Errorf("expected [existing new], got %v", got)
	}
}
