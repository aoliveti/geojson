package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProperties_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		wantErr error
	}{
		{"valid key and value", "key1", "value1", nil},
		{"overwrite existing key", "key1", "value2", nil},
		{"empty key", "", "value", ErrKeyEmpty},
		{"value is nil", "key2", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Properties
			err := p.Set(tt.key, tt.value)
			assert.ErrorIs(t, err, tt.wantErr, "unexpected error result")

			if tt.key != "" && tt.wantErr == nil {
				got, ok := p[tt.key]
				assert.True(t, ok, "key should exist in properties")
				assert.Equal(t, tt.value, got, "unexpected value set in properties")
			}
		})
	}
}

func TestProperties_Get(t *testing.T) {
	p := Properties{"key1": "value1", "key2": nil}

	tests := []struct {
		name     string
		key      string
		wantVal  interface{}
		wantDoes bool
	}{
		{"existing key", "key1", "value1", true},
		{"non-existing key", "key3", nil, false},
		{"existing key with nil value", "key2", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotDoes := p.Get(tt.key)
			assert.Equal(t, tt.wantDoes, gotDoes)
			assert.Equal(t, tt.wantVal, gotVal)
		})
	}
}

func TestProperties_GetString(t *testing.T) {
	p := Properties{"key1": "value1", "key2": 123, "key3": nil}

	tests := []struct {
		name      string
		key       string
		wantStr   string
		wantError error
	}{
		{"existing string key", "key1", "value1", nil},
		{"non-existing key", "key4", "", ErrPropertyNotFound},
		{"key with non-string value", "key2", "", ErrInvalidString},
		{"key with nil value", "key3", "", ErrInvalidString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, err := p.GetString(tt.key)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantStr, gotStr)
		})
	}
}

func TestProperties_GetInt(t *testing.T) {
	p := Properties{"key1": 123.0, "key2": "value1", "key3": nil}

	tests := []struct {
		name      string
		key       string
		wantInt   int
		wantError error
	}{
		{"existing int key", "key1", 123, nil},
		{"non-existing key", "key4", 0, ErrPropertyNotFound},
		{"key with non-int value", "key2", 0, ErrInvalidInt},
		{"key with nil value", "key3", 0, ErrInvalidInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInt, err := p.GetInt(tt.key)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantInt, gotInt)
		})
	}
}

func TestProperties_GetFloat(t *testing.T) {
	p := Properties{"key1": 123.45, "key2": "value1", "key3": nil}

	tests := []struct {
		name      string
		key       string
		wantFloat float64
		wantError error
	}{
		{"existing float key", "key1", 123.45, nil},
		{"non-existing key", "key4", 0, ErrPropertyNotFound},
		{"key with non-float value", "key2", 0, ErrInvalidFloat},
		{"key with nil value", "key3", 0, ErrInvalidFloat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFloat, err := p.GetFloat(tt.key)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantFloat, gotFloat)
		})
	}
}

func TestProperties_GetBool(t *testing.T) {
	p := Properties{"key1": true, "key2": "value1", "key3": nil}

	tests := []struct {
		name      string
		key       string
		wantBool  bool
		wantError error
	}{
		{"existing bool key", "key1", true, nil},
		{"non-existing key", "key4", false, ErrPropertyNotFound},
		{"key with non-bool value", "key2", false, ErrInvalidBool},
		{"key with nil value", "key3", false, ErrInvalidBool},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, err := p.GetBool(tt.key)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantBool, gotBool)
		})
	}
}

func TestProperties_GetFromNil(t *testing.T) {
	t.Run("Get methods for all types with non-existing key and uninitialized map", func(t *testing.T) {
		var p Properties

		nonExistingKey := "nonExistingKey"

		// Test Get
		val, exists := p.Get(nonExistingKey)
		assert.False(t, exists, "Get should indicate that the key does not exist")
		assert.Nil(t, val, "Get should return nil for a non-existing key")

		// Test GetString
		strVal, err := p.GetString(nonExistingKey)
		assert.ErrorIs(t, err, ErrPropertyNotFound, "GetString should return ErrPropertyNotFound for a non-existing key")
		assert.Empty(t, strVal, "GetString should return an empty string for a non-existing key")

		// Test GetInt
		intVal, err := p.GetInt(nonExistingKey)
		assert.ErrorIs(t, err, ErrPropertyNotFound, "GetInt should return ErrPropertyNotFound for a non-existing key")
		assert.Equal(t, 0, intVal, "GetInt should return 0 for a non-existing key")

		// Test GetFloat
		floatVal, err := p.GetFloat(nonExistingKey)
		assert.ErrorIs(t, err, ErrPropertyNotFound, "GetFloat should return ErrPropertyNotFound for a non-existing key")
		assert.Equal(t, float64(0), floatVal, "GetFloat should return 0 for a non-existing key")

		// Test GetBool
		boolVal, err := p.GetBool(nonExistingKey)
		assert.ErrorIs(t, err, ErrPropertyNotFound, "GetBool should return ErrPropertyNotFound for a non-existing key")
		assert.False(t, boolVal, "GetBool should return false for a non-existing key")
	})
}

func TestProperties_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		props    Properties
		wantJSON string
	}{
		{"empty properties", Properties{}, "null"},
		{"non-empty properties", Properties{"key1": "value1"}, `{"key1":"value1"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON, err := tt.props.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(gotJSON))
		})
	}
}

func TestProperties_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     Properties
		wantErr  error
	}{
		{"null JSON", "null", nil, nil},
		{"valid JSON", `{"key1":"value1"}`, Properties{"key1": "value1"}, nil},
		{"empty JSON", `{}`, Properties{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Properties
			err := json.Unmarshal([]byte(tt.jsonData), &p)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, p)
		})
	}
}
