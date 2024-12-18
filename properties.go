package geojson

import (
	"encoding/json"
	"errors"
)

// Error definitions for operations on the Properties type.
var (
	ErrKeyEmpty         = errors.New("property key cannot be empty")
	ErrPropertyNotFound = errors.New("property not found")
	ErrInvalidString    = errors.New("property is not a string")
	ErrInvalidInt       = errors.New("property is not an integer")
	ErrInvalidFloat     = errors.New("property is not a float")
	ErrInvalidBool      = errors.New("property is not a boolean")
)

// Properties represents a map of key-value pairs used as metadata for a GeoJSON feature.
// It adheres to the GeoJSON specification (RFC 7946) by supporting arbitrary key-value data.
type Properties map[string]interface{}

// Set assigns a value to a specified key in the Properties map.
// If the key already exists, its value gets updated.
// Returns an error if the key is empty.
func (p *Properties) Set(key string, value interface{}) error {
	if key == "" {
		return ErrKeyEmpty
	}

	if *p == nil {
		*p = make(map[string]interface{})
	}

	(*p)[key] = value

	return nil
}

// Get fetches the value associated with a key in the Properties map.
// Returns the value and a boolean indicating whether the key exists.
func (p *Properties) Get(key string) (interface{}, bool) {
	if p == nil || len(*p) == 0 {
		return nil, false
	}

	value, ok := (*p)[key]
	return value, ok
}

// GetString retrieves the value for the given key as a string.
// Returns an error if the key does not exist or the value is not a string.
func (p *Properties) GetString(key string) (string, error) {
	if p == nil || len(*p) == 0 {
		return "", ErrPropertyNotFound
	}

	value, ok := (*p)[key]
	if !ok {
		return "", ErrPropertyNotFound
	}

	strValue, ok := value.(string)
	if !ok {
		return "", ErrInvalidString
	}

	return strValue, nil
}

// GetInt retrieves the value for the given key as an integer.
// Returns an error if the key does not exist or the value is not an integer.
func (p *Properties) GetInt(key string) (int, error) {
	if p == nil || len(*p) == 0 {
		return 0, ErrPropertyNotFound
	}

	value, ok := (*p)[key]
	if !ok {
		return 0, ErrPropertyNotFound
	}

	intValue, ok := value.(float64)
	if !ok {
		return 0, ErrInvalidInt
	}

	return int(intValue), nil
}

// GetFloat retrieves the value for the given key as a float64.
// Returns an error if the key does not exist or the value is not a float64.
func (p *Properties) GetFloat(key string) (float64, error) {
	if p == nil || len(*p) == 0 {
		return 0, ErrPropertyNotFound
	}

	value, ok := (*p)[key]
	if !ok {
		return 0, ErrPropertyNotFound
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, ErrInvalidFloat
	}

	return floatValue, nil
}

// GetBool retrieves the value for the given key as a boolean.
// Returns an error if the key does not exist or the value is not a boolean.
func (p *Properties) GetBool(key string) (bool, error) {
	if p == nil || len(*p) == 0 {
		return false, ErrPropertyNotFound
	}

	value, ok := (*p)[key]
	if !ok {
		return false, ErrPropertyNotFound
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, ErrInvalidBool
	}

	return boolValue, nil
}

// MarshalJSON converts the Properties map to a JSON-encoded byte slice.
// Serializes to null if the map is nil or empty.
func (p *Properties) MarshalJSON() ([]byte, error) {
	if p == nil || len(*p) == 0 {
		return json.Marshal(nil)
	}

	return json.Marshal(map[string]interface{}(*p))
}

// UnmarshalJSON parses a JSON-encoded byte slice and stores the result into the Properties map.
func (p *Properties) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = raw

	return nil
}
