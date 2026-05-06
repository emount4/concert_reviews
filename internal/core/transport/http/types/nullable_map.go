// core_http_types/nullable_map.go
package core_http_types

import (
	"encoding/json"
	"fmt"

	core_types "github.com/emount4/concert_reviews/internal/core/types"
)

type NullableMapStringString struct {
	core_types.Nullable[map[string]string]
}

func (n *NullableMapStringString) UnmarshalJSON(b []byte) error {
	n.Set = true
	if string(b) == "null" {
		n.Value = nil
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("invalid social_links: %w", err)
	}
	n.Value = &m
	return nil
}

func (n NullableMapStringString) MarshalJSON() ([]byte, error) {
	if !n.Set || n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

// Хелперы для удобства
func NewNullableMapStringString(m map[string]string, set bool) NullableMapStringString {
	return NullableMapStringString{core_types.Nullable[map[string]string]{Value: &m, Set: set}}
}

func NullableMapStringStringNull() NullableMapStringString {
	return NullableMapStringString{core_types.Nullable[map[string]string]{Value: nil, Set: true}}
}
