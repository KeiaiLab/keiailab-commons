// SPDX-License-Identifier: MIT

package secrethash_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/keiailab/keiailab-commons/pkg/secrethash"
)

// manual 은 테스트 기대값을 독립적으로 재현한다 (구현과 동일 알고리즘을 손으로).
func manual(pairs ...[2]string) string {
	h := sha256.New()
	for _, p := range pairs {
		_, _ = h.Write([]byte(p[0]))
		_, _ = h.Write([]byte(p[1]))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func TestHash_OrderedKeys_MatchesManual(t *testing.T) {
	t.Parallel()
	data := map[string][]byte{
		"tls.crt": []byte("CERT"),
		"tls.key": []byte("KEY"),
		"ca.crt":  []byte("CA"),
	}
	got := secrethash.Hash(data, "tls.crt", "tls.key", "ca.crt")
	want := manual([2]string{"tls.crt", "CERT"}, [2]string{"tls.key", "KEY"}, [2]string{"ca.crt", "CA"})
	if got != want {
		t.Fatalf("ordered-keys digest mismatch:\n got=%s\nwant=%s", got, want)
	}
}

func TestHash_AllKeys_SortedDeterministic(t *testing.T) {
	t.Parallel()
	data := map[string][]byte{"b": []byte("2"), "a": []byte("1"), "c": []byte("3")}
	// keys 미지정 → 정렬(a,b,c) 누적. 명시 정렬 키와 동일해야 한다.
	got := secrethash.Hash(data)
	want := secrethash.Hash(data, "a", "b", "c")
	if got != want {
		t.Fatalf("all-keys digest must equal sorted explicit keys:\n got=%s\nwant=%s", got, want)
	}
}

// TestHash_Deterministic 은 mongodb-operator 의 기존 map-순회 비결정성 버그가
// 재현되지 않음을 보증한다 — 동일 입력은 항상 동일 digest.
func TestHash_Deterministic(t *testing.T) {
	t.Parallel()
	data := map[string][]byte{
		"k0": []byte("v0"), "k1": []byte("v1"), "k2": []byte("v2"),
		"k3": []byte("v3"), "k4": []byte("v4"), "k5": []byte("v5"),
	}
	first := secrethash.Hash(data)
	for i := range 200 {
		if got := secrethash.Hash(data); got != first {
			t.Fatalf("non-deterministic digest on iter %d: %s != %s", i, got, first)
		}
	}
}

func TestHash_ValueChangeChangesDigest(t *testing.T) {
	t.Parallel()
	a := secrethash.Hash(map[string][]byte{"tls.crt": []byte("OLD")}, "tls.crt")
	b := secrethash.Hash(map[string][]byte{"tls.crt": []byte("NEW")}, "tls.crt")
	if a == b {
		t.Fatal("값 변경이 digest 에 반영되지 않음")
	}
}

// TestHash_KeyPresenceChangesDigest — 키 이름을 누적하므로, 빈 값이라도 키가
// 추가/삭제되면 digest 가 변해야 한다 (예: ca.crt 만 추가).
func TestHash_KeyPresenceChangesDigest(t *testing.T) {
	t.Parallel()
	withCA := secrethash.Hash(map[string][]byte{"tls.crt": []byte("C")}, "tls.crt", "ca.crt")
	without := secrethash.Hash(map[string][]byte{"tls.crt": []byte("C")}, "tls.crt")
	if withCA == without {
		t.Fatal("키 누락/추가가 digest 에 반영되지 않음")
	}
}

func TestHash_OrderMattersForExplicitKeys(t *testing.T) {
	t.Parallel()
	data := map[string][]byte{"a": []byte("1"), "b": []byte("2")}
	ab := secrethash.Hash(data, "a", "b")
	ba := secrethash.Hash(data, "b", "a")
	if ab == ba {
		t.Fatal("명시 키 순서가 digest 에 영향을 주지 않음 (순서 의존성 상실)")
	}
}

func TestHash_EmptyInput(t *testing.T) {
	t.Parallel()
	// 빈 입력 → SHA256("") 고정값.
	const sha256Empty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got := secrethash.Hash(nil); got != sha256Empty {
		t.Fatalf("empty digest = %s, want %s", got, sha256Empty)
	}
	if got := secrethash.Hash(map[string][]byte{}); got != sha256Empty {
		t.Fatalf("empty map digest = %s, want %s", got, sha256Empty)
	}
}

func TestHash_MissingKeyTreatedAsEmptyValue(t *testing.T) {
	t.Parallel()
	// 명시 키가 map 에 없으면 빈 값으로 취급 — 키 이름만 누적.
	got := secrethash.Hash(map[string][]byte{}, "absent")
	want := manual([2]string{"absent", ""})
	if got != want {
		t.Fatalf("missing-key digest mismatch:\n got=%s\nwant=%s", got, want)
	}
}
