package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_StringValue(t *testing.T) {
	tests := []struct {
		name     string
		id       *ID
		expected string
		ok       bool
	}{
		{"string value", NewStringID("test"), "test", true},
		{"no string value", NewNumericID(42), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.id.StringValue()
			assert.Equal(t, tt.expected, val)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestID_NumberValue(t *testing.T) {
	tests := []struct {
		name     string
		id       *ID
		expected float64
		ok       bool
	}{
		{"numeric value", NewNumericID(42), 42, true},
		{"no numeric value", NewStringID("test"), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.id.NumberValue()
			assert.Equal(t, tt.expected, val)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		id       *ID
		expected string
	}{
		{"string ID", NewStringID("test"), `"test"`},
		{"numeric ID", NewNumericID(42), `42`},
		{"nil ID", &ID{}, `null`},
		{"nil ID", nil, `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&tt.id)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *ID
		expectError bool
	}{
		{"valid string", `"test"`, NewStringID("test"), false},
		{"int number", `42`, NewNumericID(42), false},
		{"null value", `null`, nil, false},
		{"invalid value", `true`, nil, true},
		{"empty JSON", `{}`, nil, true},
		{"empty string", ``, nil, true},
		{"float number", `42.0`, NewNumericID(42.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id *ID
			err := json.Unmarshal([]byte(tt.input), &id)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.expected, id, "expected: %v, got: %v", tt.expected, id)
			}
		})
	}
}
