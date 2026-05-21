// SPDX-License-Identifier: Apache-2.0

package version

import (
	"encoding/json"
	"sort"
)

// AsMap 는 Matrix[E] 의 entries 를 PrimaryKey → entry 매핑으로 직렬화한다.
// JSON / YAML 모두 호환 (encoding/json 으로 marshal 가능).
//
// 사용 예:
//
//	data, _ := json.Marshal(matrix.AsMap())
//
// Refs: docs/ROADMAP.md '버전 매트릭스 시리얼라이저 (json/yaml)'
func (m Matrix[E]) AsMap() map[string]E {
	out := make(map[string]E, len(m.entries))
	for _, e := range m.entries {
		out[e.PrimaryKey()] = e
	}
	return out
}

// MarshalJSON — Matrix[E] 의 JSON 표현. PrimaryKey 정렬 후 stable map.
//
// 사용 예:
//
//	data, err := json.Marshal(version.Supported)
func (m Matrix[E]) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(m.entries))
	mp := make(map[string]E, len(m.entries))
	for _, e := range m.entries {
		k := e.PrimaryKey()
		keys = append(keys, k)
		mp[k] = e
	}
	sort.Strings(keys)
	// stable ordered map representation
	type kv struct {
		Key   string `json:"key"`
		Entry E      `json:"entry"`
	}
	out := make([]kv, 0, len(keys))
	for _, k := range keys {
		out = append(out, kv{Key: k, Entry: mp[k]})
	}
	return json.Marshal(out)
}
