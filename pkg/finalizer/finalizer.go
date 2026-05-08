// Package finalizer 는 4-repo keiailab operator 공통 Finalizer 헬퍼를
// 제공한다.
//
// 본 패키지는 RFC-0018 §3.2 spec 을 구현한다. mongodb-operator
// api/v1alpha1/finalizers.go 패턴 (상수 + Add/Remove 헬퍼) 을 표준화.
//
// 의존성: k8s.io/apimachinery 만. controller-runtime 의 controllerutil
// 미의존 — finalizer slice 직접 조작 (commons 의 의존성 표면을 최소화).
// client.Update 은 호출자가 직접 수행 (본 패키지는 *in-memory* 변경만).
package finalizer

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Prefix — keiailab 표준 finalizer prefix. repo 별 specific suffix 는 호출자가
// 부여 (예: "mongodb.keiailab.com/finalizer", "valkey.keiailab.com/finalizer").
const Prefix = "keiailab.com/"

// Add 는 obj 의 finalizer slice 에 name 이 없으면 추가한다.
//
// 반환값: 변경 발생 시 true (호출자는 client.Update(ctx, obj) 로 persist).
// 이미 존재하면 false (no-op, Update 불필요).
//
// idempotent — 동일 name 으로 두 번 호출해도 두 번째는 false.
func Add(obj metav1.Object, name string) bool {
	finalizers := obj.GetFinalizers()
	for _, f := range finalizers {
		if f == name {
			return false
		}
	}
	obj.SetFinalizers(append(finalizers, name))
	return true
}

// Remove 는 obj 의 finalizer slice 에서 name 을 제거한다.
//
// 반환값: 변경 발생 시 true (호출자는 client.Update). 없으면 false.
func Remove(obj metav1.Object, name string) bool {
	finalizers := obj.GetFinalizers()
	for i, f := range finalizers {
		if f == name {
			obj.SetFinalizers(append(finalizers[:i], finalizers[i+1:]...))
			return true
		}
	}
	return false
}

// Has 는 obj 에 name finalizer 가 등록되어 있는지 확인.
func Has(obj metav1.Object, name string) bool {
	for _, f := range obj.GetFinalizers() {
		if f == name {
			return true
		}
	}
	return false
}
