package version_test

import (
	"testing"

	"github.com/keiailab/operator-commons/pkg/version"
)

// testEntry — 테스트용 MatrixEntry 구현. postgres 의 Combo 패턴 모사.
type testEntry struct {
	Version string
	Channel string
	Image   string
}

func (t testEntry) PrimaryKey() string { return t.Version }

// TestMustMatrix_FindAndIsSupported — 정상 생성 + 검색.
func TestMustMatrix_FindAndIsSupported(t *testing.T) {
	m := version.MustMatrix(
		testEntry{Version: "1.0", Channel: "stable", Image: "img:1.0"},
		testEntry{Version: "2.0", Channel: "beta", Image: "img:2.0"},
	)

	e, ok := m.Find("1.0")
	if !ok {
		t.Fatal("Find(1.0) ok=false")
	}
	if e.Channel != "stable" || e.Image != "img:1.0" {
		t.Errorf("Find(1.0) = %+v, want Channel=stable Image=img:1.0", e)
	}

	if !m.IsSupported("1.0") {
		t.Error("IsSupported(1.0) = false")
	}
	if !m.IsSupported("2.0") {
		t.Error("IsSupported(2.0) = false")
	}
	if m.IsSupported("3.0") {
		t.Error("IsSupported(3.0) = true (unexpected)")
	}
}

// TestMustMatrix_NotFound — 없는 key 는 zero E + false.
func TestMustMatrix_NotFound(t *testing.T) {
	m := version.MustMatrix(testEntry{Version: "1.0"})
	e, ok := m.Find("99.0")
	if ok {
		t.Errorf("Find(99.0) ok=true, want false")
	}
	if e.Version != "" || e.Channel != "" {
		t.Errorf("Find(99.0) zero E = %+v", e)
	}
}

// TestMustMatrix_KeysAndEntriesAndLen — 외부 노출 helper.
func TestMustMatrix_KeysAndEntriesAndLen(t *testing.T) {
	m := version.MustMatrix(
		testEntry{Version: "1.0"},
		testEntry{Version: "2.0"},
		testEntry{Version: "3.0"},
	)

	if got := m.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
	keys := m.Keys()
	if len(keys) != 3 || keys[0] != "1.0" || keys[2] != "3.0" {
		t.Errorf("Keys() = %v, want [1.0 2.0 3.0]", keys)
	}
	entries := m.Entries()
	if len(entries) != 3 {
		t.Errorf("Entries() len = %d, want 3", len(entries))
	}
	// 방어 복사 확인 — 호출자 mutation 이 내부 영향 없음.
	entries[0] = testEntry{Version: "MUTATED"}
	if e, _ := m.Find("1.0"); e.Version != "1.0" {
		t.Errorf("내부 mutation 노출 — Entries 가 방어 복사 미실행: %+v", e)
	}
}

// TestMustMatrix_EmptyPanics — 빈 리스트 panic.
func TestMustMatrix_EmptyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty matrix")
		}
	}()
	version.MustMatrix[testEntry]()
}

// TestMustMatrix_EmptyPrimaryKeyPanics — 빈 PrimaryKey panic.
func TestMustMatrix_EmptyPrimaryKeyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty PrimaryKey")
		}
	}()
	version.MustMatrix(testEntry{Version: "", Channel: "stable"})
}

// TestMustMatrix_DuplicatePrimaryKeyPanics — 중복 PrimaryKey panic.
func TestMustMatrix_DuplicatePrimaryKeyPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate PrimaryKey")
		}
	}()
	version.MustMatrix(
		testEntry{Version: "1.0", Channel: "stable"},
		testEntry{Version: "1.0", Channel: "beta"},
	)
}

// TestMustMatrix_ListCoexistence — 기존 List 와 Matrix[E] 가 동일 패키지
// 에서 공존. 기존 호출자 영향 없음을 검증.
func TestMustMatrix_ListCoexistence(t *testing.T) {
	l := version.MustList("a", "b", "c")
	m := version.MustMatrix(
		testEntry{Version: "1.0"},
		testEntry{Version: "2.0"},
	)
	if !l.IsSupported("a") {
		t.Error("List.IsSupported regression")
	}
	if !m.IsSupported("1.0") {
		t.Error("Matrix.IsSupported failed in coexistence")
	}
}
