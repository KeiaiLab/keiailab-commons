// SPDX-License-Identifier: MIT

package version_test

import (
	"encoding/json"
	"testing"

	"github.com/keiailab/keiailab-commons/pkg/version"
)

// stableTestEntry — api stability test 전용 MatrixEntry 구현.
// (matrix_test.go 의 testEntry 와 분리 위해 별 이름.)
type stableTestEntry struct {
	Key string
}

func (s stableTestEntry) PrimaryKey() string { return s.Key }

// TestAPIStability — public API surface 가 v0.x 진화에서 유지되는지 확인.
//
// 본 테스트는 downstream operator 의 호출 패턴을 모방한
// *cross-version compatibility 표면 가드*. ROADMAP §"Cross-version
// compatibility test — 3 repo CI 통합" 의 first slice.
//
// 호출자가 의존하는 표면 (List / IsSupported / Strings / Default) +
// Matrix[E] (MustMatrix / IsSupported / Find / Entries / Keys / Len) +
// serializer (AsMap / MarshalJSON) 변경 시 본 테스트 fail.
func TestAPIStability(t *testing.T) {
	// 표면 1: List
	supported := version.MustList("1.0", "2.0")
	if !supported.IsSupported("1.0") {
		t.Errorf("List.IsSupported('1.0') = false, want true")
	}
	if got := supported.Default(); got != "1.0" {
		t.Errorf("List.Default() = %q, want %q", got, "1.0")
	}
	if got := supported.Strings(); len(got) != 2 {
		t.Errorf("List.Strings() len = %d, want 2", len(got))
	}

	// 표면 2: Matrix[E]
	matrix := version.MustMatrix(stableTestEntry{Key: "a"}, stableTestEntry{Key: "b"})
	if !matrix.IsSupported("a") {
		t.Errorf("Matrix.IsSupported('a') = false, want true")
	}
	if _, ok := matrix.Find("b"); !ok {
		t.Errorf("Matrix.Find('b') ok = false, want true")
	}
	if matrix.Len() != 2 {
		t.Errorf("Matrix.Len() = %d, want 2", matrix.Len())
	}
	if got := matrix.Keys(); len(got) != 2 {
		t.Errorf("Matrix.Keys() len = %d, want 2", len(got))
	}
	if got := matrix.Entries(); len(got) != 2 {
		t.Errorf("Matrix.Entries() len = %d, want 2", len(got))
	}

	// 표면 3: serializer (B.6.2)
	asMap := matrix.AsMap()
	if len(asMap) != 2 {
		t.Errorf("AsMap len = %d, want 2", len(asMap))
	}
	if asMap["a"].Key != "a" {
		t.Errorf("AsMap['a'].Key = %q, want %q", asMap["a"].Key, "a")
	}
	data, err := json.Marshal(matrix)
	if err != nil {
		t.Errorf("MarshalJSON err = %v", err)
	}
	if len(data) == 0 {
		t.Errorf("MarshalJSON empty output")
	}
}
