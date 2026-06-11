// SPDX-License-Identifier: MIT

package reconcile

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const (
	testNS        = "ns"
	testObjName   = "obj-1"
	testComponent = "storage"
)

var testGV = schema.GroupVersion{Group: "test.keiailab.com", Version: "v1"}

// testObject — Statusable 구현 테스트 더블 (CR 모사).
type testObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            testObjectStatus `json:"status,omitempty"`
}

type testObjectStatus struct {
	Phase      string             `json:"phase,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (o *testObject) GetConditions() *[]metav1.Condition { return &o.Status.Conditions }
func (o *testObject) SetPhase(phase string)              { o.Status.Phase = phase }

func (o *testObject) DeepCopyObject() runtime.Object {
	out := &testObject{TypeMeta: o.TypeMeta}
	om := o.ObjectMeta // 구체 타입 로컬 경유 — embedded selector 의 승격 모호 회피 (QF1008)
	om.DeepCopyInto(&out.ObjectMeta)
	out.Status.Phase = o.Status.Phase
	if o.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(o.Status.Conditions))
		for i := range o.Status.Conditions {
			o.Status.Conditions[i].DeepCopyInto(&out.Status.Conditions[i])
		}
	}
	return out
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	s.AddKnownTypes(testGV, &testObject{})
	metav1.AddToGroupVersion(s, testGV)
	if err := corev1.AddToScheme(s); err != nil {
		t.Fatalf("corev1 scheme: %v", err)
	}
	return s
}

func newTestObject(conds ...metav1.Condition) *testObject {
	return &testObject{
		TypeMeta:   metav1.TypeMeta{APIVersion: testGV.String(), Kind: "testObject"},
		ObjectMeta: metav1.ObjectMeta{Name: testObjName, Namespace: testNS},
		Status:     testObjectStatus{Conditions: conds},
	}
}

func newStatusClient(t *testing.T, obj *testObject, funcs ...interceptor.Funcs) client.Client {
	t.Helper()
	b := fake.NewClientBuilder().
		WithScheme(testScheme(t)).
		WithObjects(obj).
		WithStatusSubresource(&testObject{})
	if len(funcs) > 0 {
		b = b.WithInterceptorFuncs(funcs[0])
	}
	return b.Build()
}

// recordedEvent / fakeRecorder — EventRecorder 관측 더블.
type recordedEvent struct {
	eventtype string
	reason    string
	action    string
	note      string
}

type fakeRecorder struct {
	events []recordedEvent
}

func (f *fakeRecorder) Eventf(_ runtime.Object, _ runtime.Object,
	eventtype, reason, action, note string, args ...any) {
	f.events = append(f.events, recordedEvent{
		eventtype: eventtype,
		reason:    reason,
		action:    action,
		note:      fmt.Sprintf(note, args...),
	})
}

func getTestObject(t *testing.T, c client.Client) *testObject {
	t.Helper()
	got := &testObject{}
	if err := c.Get(context.Background(),
		client.ObjectKey{Namespace: testNS, Name: testObjName}, got); err != nil {
		t.Fatalf("re-get: %v", err)
	}
	return got
}

func TestApplyErrorCondition(t *testing.T) {
	reconcileErr := errors.New("disk full")
	tests := []struct {
		name        string
		opts        []Option
		nilRecorder bool
		wantRequeue time.Duration
		wantType    string
		wantReason  string
		wantEvents  int
	}{
		{
			name:        "기본값 — ReconcileError condition + Failed phase + 30s requeue",
			wantRequeue: DefaultRequeueAfter,
			wantType:    ConditionTypeReconcileError,
			wantReason:  ReasonReconcileFailed,
			wantEvents:  1,
		},
		{
			name: "옵션 override — condition type / reason / requeue",
			opts: []Option{
				WithConditionType("Degraded"),
				WithReason("StorageFailed"),
				WithRequeueAfter(5 * time.Second),
			},
			wantRequeue: 5 * time.Second,
			wantType:    "Degraded",
			wantReason:  "StorageFailed",
			wantEvents:  1,
		},
		{
			name:        "recorder nil — 이벤트 발행 없이 정상 동작",
			nilRecorder: true,
			wantRequeue: DefaultRequeueAfter,
			wantType:    ConditionTypeReconcileError,
			wantReason:  ReasonReconcileFailed,
			wantEvents:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			obj := newTestObject()
			c := newStatusClient(t, obj)
			rec := &fakeRecorder{}
			var recorder EventRecorder = rec
			if tt.nilRecorder {
				recorder = nil
			}

			res, err := ApplyErrorCondition(ctx, c, obj, testComponent, reconcileErr, recorder, tt.opts...)

			if !errors.Is(err, reconcileErr) {
				t.Errorf("원본 reconcile 에러 반환 기대, got %v", err)
			}
			if res.RequeueAfter != tt.wantRequeue {
				t.Errorf("RequeueAfter = %s, want %s", res.RequeueAfter, tt.wantRequeue)
			}
			if len(rec.events) != tt.wantEvents {
				t.Errorf("이벤트 %d건, want %d", len(rec.events), tt.wantEvents)
			}
			if tt.wantEvents > 0 {
				ev := rec.events[0]
				if ev.eventtype != corev1.EventTypeWarning ||
					ev.reason != tt.wantType || ev.action != tt.wantType {
					t.Errorf("Warning/%s/%s 이벤트 기대, got %+v", tt.wantType, tt.wantType, ev)
				}
			}

			// 서버 영속 상태 검증 (in-memory 가 아닌 fake API server 진본).
			got := getTestObject(t, c)
			if got.Status.Phase != PhaseFailed {
				t.Errorf("phase = %q, want %q", got.Status.Phase, PhaseFailed)
			}
			cond := meta.FindStatusCondition(got.Status.Conditions, tt.wantType)
			if cond == nil {
				t.Fatalf("condition %q 부재", tt.wantType)
			}
			if cond.Status != metav1.ConditionTrue || cond.Reason != tt.wantReason {
				t.Errorf("condition = %+v, want True/%s", cond, tt.wantReason)
			}
		})
	}
}

func TestApplyErrorCondition_metric_hook(t *testing.T) {
	ctx := context.Background()
	obj := newTestObject()
	c := newStatusClient(t, obj)

	hookCalls := 0
	var gotNS, gotName, gotComponent string
	_, _ = ApplyErrorCondition(ctx, c, obj, testComponent, errors.New("x"), nil,
		WithMetricHook(func(namespace, name, component string) {
			hookCalls++
			gotNS, gotName, gotComponent = namespace, name, component
		}))

	if hookCalls != 1 {
		t.Fatalf("metric hook 1회 호출 기대, got %d", hookCalls)
	}
	if gotNS != testNS || gotName != testObjName || gotComponent != testComponent {
		t.Errorf("hook 라벨 = (%s,%s,%s), want (%s,%s,%s)",
			gotNS, gotName, gotComponent, testNS, testObjName, testComponent)
	}
}

func TestApplyErrorCondition_LastTransitionTime_preserved(t *testing.T) {
	// ADR-0013 canonical 거동 — Status 미변경 (True→True) 시 LastTransitionTime 보존.
	// valkey 구판 (filter+append+Now) 은 매 호출 false transition 을 만들던 버그 —
	// 본 테스트가 그 회귀 가드.
	ctx := context.Background()
	oldTime := metav1.NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	obj := newTestObject(metav1.Condition{
		Type:               ConditionTypeReconcileError,
		Status:             metav1.ConditionTrue,
		Reason:             ReasonReconcileFailed,
		Message:            "previous failure",
		LastTransitionTime: oldTime,
	})
	c := newStatusClient(t, obj)

	_, _ = ApplyErrorCondition(ctx, c, obj, testComponent, errors.New("again"), nil)

	cond := meta.FindStatusCondition(getTestObject(t, c).Status.Conditions, ConditionTypeReconcileError)
	if cond == nil {
		t.Fatal("condition 부재")
	}
	if !cond.LastTransitionTime.Equal(&oldTime) {
		t.Errorf("LastTransitionTime 갱신됨 (%s) — Status 미변경 시 보존 기대 (want %s)",
			cond.LastTransitionTime, oldTime)
	}
	if cond.Message == "previous failure" {
		t.Error("Message 는 최신 에러로 갱신 기대")
	}
}

func TestApplyErrorCondition_conflict_mutate_reapply(t *testing.T) {
	// status update conflict 1회 → UpdateWithRetry 가 refetch 후 mutate 클로저
	// 재적용 — 에러 condition 이 유실 없이 영속되는지 검증.
	ctx := context.Background()
	obj := newTestObject()
	calls := 0
	c := newStatusClient(t, obj, interceptor.Funcs{
		SubResourceUpdate: func(ictx context.Context, inner client.Client, sub string,
			iobj client.Object, opts ...client.SubResourceUpdateOption) error {
			calls++
			if calls == 1 {
				return apierrors.NewConflict(
					schema.GroupResource{Group: testGV.Group, Resource: "testobjects"},
					iobj.GetName(), nil)
			}
			return inner.SubResource(sub).Update(ictx, iobj, opts...)
		},
	})

	_, _ = ApplyErrorCondition(ctx, c, obj, testComponent, errors.New("boom"), nil)

	if calls != 2 {
		t.Fatalf("conflict 1회 후 재시도 — status update 2회 기대, got %d", calls)
	}
	got := getTestObject(t, c)
	if got.Status.Phase != PhaseFailed {
		t.Errorf("conflict 후에도 phase 영속 기대, got %q", got.Status.Phase)
	}
	if meta.FindStatusCondition(got.Status.Conditions, ConditionTypeReconcileError) == nil {
		t.Error("conflict 후에도 condition 영속 기대 (mutate 재적용)")
	}
}
