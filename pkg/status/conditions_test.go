package status_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keiailab/operator-commons/pkg/status"
)

// TestSetReady 는 SetReady 가 빈 conditions 에 새 Condition 을 추가하는지,
// 그리고 동일 type 호출 시 *덮어쓰기* (replace) 되는지 검증한다.
//
// AAA 형식:
//
//	Arrange — 빈 conditions 슬라이스
//	Act — SetReady(True, Available) → SetReady(False, ReconcileError)
//	Assert — 마지막 호출만 남고, status/reason 은 두 번째 값
func TestSetReady_AddsAndReplaces(t *testing.T) {
	var conds []metav1.Condition

	status.SetReady(&conds, metav1.ConditionTrue, status.ReasonAvailable, "ok", 5)
	if len(conds) != 1 {
		t.Fatalf("expected 1 condition after first SetReady, got %d", len(conds))
	}
	if conds[0].Type != status.TypeReady {
		t.Errorf("expected Type=%q, got %q", status.TypeReady, conds[0].Type)
	}
	if conds[0].Status != metav1.ConditionTrue {
		t.Errorf("expected Status=True, got %v", conds[0].Status)
	}

	// 두 번째 호출 — 동일 type 이므로 replace.
	status.SetReady(&conds, metav1.ConditionFalse, status.ReasonReconcileError, "boom", 6)
	if len(conds) != 1 {
		t.Fatalf("expected still 1 condition (replace), got %d", len(conds))
	}
	if conds[0].Status != metav1.ConditionFalse {
		t.Errorf("expected Status=False after replace, got %v", conds[0].Status)
	}
	if conds[0].Reason != status.ReasonReconcileError {
		t.Errorf("expected Reason=ReconcileError after replace, got %q", conds[0].Reason)
	}
	if conds[0].ObservedGeneration != 6 {
		t.Errorf("expected ObservedGeneration=6, got %d", conds[0].ObservedGeneration)
	}
}

// TestIsReady — Ready=True 일 때만 true 반환.
func TestIsReady(t *testing.T) {
	tests := []struct {
		name  string
		conds []metav1.Condition
		want  bool
	}{
		{"empty", nil, false},
		{"ready_true", []metav1.Condition{{Type: status.TypeReady, Status: metav1.ConditionTrue}}, true},
		{"ready_false", []metav1.Condition{{Type: status.TypeReady, Status: metav1.ConditionFalse}}, false},
		{"only_progressing", []metav1.Condition{{Type: status.TypeProgressing, Status: metav1.ConditionTrue}}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := status.IsReady(tc.conds)
			if got != tc.want {
				t.Errorf("IsReady(%v) = %v, want %v", tc.conds, got, tc.want)
			}
		})
	}
}

// TestSetProgressing_Coexists — Progressing 과 Ready 가 *공존* 가능 (서로 다른
// type 이므로). meta.SetStatusCondition 은 type 단위로만 replace.
func TestSetProgressing_Coexists(t *testing.T) {
	var conds []metav1.Condition

	status.SetReady(&conds, metav1.ConditionFalse, status.ReasonReconciling, "starting", 1)
	status.SetProgressing(&conds, status.ReasonReconciling, "creating sts", 1)

	if len(conds) != 2 {
		t.Fatalf("expected 2 conditions (Ready + Progressing), got %d", len(conds))
	}
	if status.FindCondition(conds, status.TypeReady) == nil {
		t.Error("expected Ready condition to exist")
	}
	if status.FindCondition(conds, status.TypeProgressing) == nil {
		t.Error("expected Progressing condition to exist")
	}
}

// TestRemoveCondition — 명시 제거.
func TestRemoveCondition(t *testing.T) {
	conds := []metav1.Condition{
		{Type: status.TypeReady, Status: metav1.ConditionTrue},
		{Type: status.TypeDegraded, Status: metav1.ConditionTrue},
	}
	status.RemoveCondition(&conds, status.TypeDegraded)
	if len(conds) != 1 || conds[0].Type != status.TypeReady {
		t.Errorf("expected only Ready remaining, got %v", conds)
	}
	// idempotent: 재제거 시 no-op.
	status.RemoveCondition(&conds, status.TypeDegraded)
	if len(conds) != 1 {
		t.Errorf("expected idempotent remove, got len=%d", len(conds))
	}
}
