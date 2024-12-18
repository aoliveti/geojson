package geojson

import (
	"encoding/json"
	"fmt"
)

var (
	// ErrInvalidID represents an error for invalid ID types or formats.
	ErrInvalidID = fmt.Errorf("invalid ID: unexpected type or format")
)

// ID represents a GeoJSON feature ID which can be either a string or a number.
type ID struct {
	s *string  // String ID value.
	n *float64 // Numeric ID value.
}

// NewStringID creates a new ID instance initialized with a string value.
func NewStringID(s string) *ID {
	return &ID{s: &s}
}

// NewNumericID creates a new ID instance initialized with a numeric value.
func NewNumericID(n float64) *ID {
	return &ID{n: &n}
}

// StringValue retrieves the string value of the ID.
// Returns the value and a boolean indicating if the value is set.
func (id *ID) StringValue() (string, bool) {
	if id.s != nil {
		return *id.s, true
	}
	return "", false
}

// NumberValue retrieves the numeric value of the ID.
// Returns the value and a boolean indicating if the value is set.
func (id *ID) NumberValue() (float64, bool) {
	if id.n != nil {
		return *id.n, true
	}
	return 0, false
}

// MarshalJSON serializes the ID into its JSON representation.
// It supports both string and numeric values.
func (id *ID) MarshalJSON() ([]byte, error) {
	if id.s != nil {
		return json.Marshal(*id.s)
	}
	if id.n != nil {
		return json.Marshal(*id.n)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON deserializes a JSON value into the ID instance.
// It supports both string and numeric types, returning an error for invalid types.
func (id *ID) UnmarshalJSON(bytes []byte) error {
	var v interface{}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("failed to unmarshal ID: %w", err)
	}

	switch value := v.(type) {
	case string:
		*id = *NewStringID(value)
	case float64:
		*id = *NewNumericID(value)
	default:
		return ErrInvalidID
	}

	return nil
}
