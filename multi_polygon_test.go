package geojson

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiPolygon_Type(t *testing.T) {
	mp := NewMultiPolygon()
	if got := mp.Type(); got != TypeMultiPolygon {
		t.Errorf("Expected %v, got %v", TypeMultiPolygon, got)
	}
}

func TestMultiPolygon_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		input    []LinearRings
		expected Vertices
	}{
		{
			name:     "empty",
			input:    []LinearRings{},
			expected: Vertices{},
		},
		{
			name: "single ring slice",
			input: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
			},
			expected: Vertices{{1, 1}, {2, 2}, {3, 3}, {1, 1}},
		},
		{
			name: "multiple ring slices",
			input: []LinearRings{
				{{{1, 1}, {3, 3}, {1, 0}, {1, 1}}},
				{{{5, 5}, {6, 6}, {5, 0}, {5, 5}}},
			},
			expected: Vertices{{1, 1}, {3, 3}, {1, 0}, {1, 1}, {5, 5}, {6, 6}, {5, 0}, {5, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := MustMultiPolygonFromRingSlice(tt.input)
			got := mp.Vertices()
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %v vertices, got %v", len(tt.expected), len(got))
			}
		})
	}
}

func TestMultiPolygon_BoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		input    []LinearRings
		expected bool
	}{
		{
			name:     "empty",
			input:    []LinearRings{},
			expected: true,
		},
		{
			name: "valid bounding box",
			input: []LinearRings{
				{{{0, 0}, {4, 0}, {4, 4}, {0, 4}, {0, 0}}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := MustMultiPolygonFromRingSlice(tt.input)
			bb := mp.BoundingBox()
			got := bb.IsValid()
			if got != tt.expected {
				t.Errorf("Expected validity %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestMultiPolygon_MarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []LinearRings
		serializeBBox bool
		expectedErr   error
	}{
		{
			name:          "empty MultiPolygon",
			input:         []LinearRings{},
			serializeBBox: false,
			expectedErr:   nil,
		},
		{
			name: "valid MultiPolygon",
			input: []LinearRings{
				{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}},
			},
			serializeBBox: true,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := MustMultiPolygonFromRingSlice(tt.input)
			mp.SerializeBBox = tt.serializeBBox
			_, err := json.Marshal(mp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestMultiPolygon_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		inputJSON    string
		expectedErr  error
		expectedType GeometryType
	}{
		{
			name: "valid MultiPolygon JSON",
			inputJSON: `{
				"type": "MultiPolygon",
				"coordinates": [
					[
						[[1, 1], [2, 2], [3, 3], [1, 1]]
					]
				]
			}`,
			expectedErr:  nil,
			expectedType: TypeMultiPolygon,
		},
		{
			name: "different type",
			inputJSON: `{
				"type": "Point",
				"coordinates": [1, 1]
			}`,
			expectedErr:  ErrInvalidTypeField,
			expectedType: "",
		},
		{
			name:         "invalid JSON",
			inputJSON:    `{"type": "Invalid"}`,
			expectedErr:  ErrInvalidTypeField,
			expectedType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := NewMultiPolygon()
			err := json.Unmarshal([]byte(tt.inputJSON), mp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && mp.Type() != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, mp.Type())
			}
		})
	}
}

func TestNewMultiPolygonFromRingSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       []LinearRings
		expectedErr error
	}{
		{
			name: "valid rings",
			input: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
			},
			expectedErr: nil,
		},
		{
			name: "invalid ring size",
			input: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}}},
			},
			expectedErr: ErrLinearRingSize,
		},
		{
			name: "ring not closed",
			input: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}, {4, 4}}},
			},
			expectedErr: ErrLinearRingClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMultiPolygonFromRingSlice(tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestMustMultiPolygonFromRingSlice(t *testing.T) {
	type args struct {
		slice []LinearRings
	}
	tests := []struct {
		name        string
		args        args
		want        *MultiPolygon
		expectPanic bool
	}{
		{
			name: "valid input",
			args: args{
				slice: []LinearRings{
					{{{1, 1}, {2, 1}, {1, 0}, {1, 1}}},
				},
			},
			want: &MultiPolygon{
				rings: []LinearRings{
					{{{1, 1}, {1, 0}, {2, 1}, {1, 1}}},
				},
			},
			expectPanic: false,
		},
		{
			name: "invalid input - ring size",
			args: args{
				slice: []LinearRings{
					{{{1, 1}, {2, 2}, {3, 3}}},
				},
			},
			want:        nil,
			expectPanic: true,
		},
		{
			name: "invalid input - ring not closed",
			args: args{
				slice: []LinearRings{
					{{{1, 1}, {2, 2}, {3, 3}, {4, 4}}},
				},
			},
			want:        nil,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				require.Panics(t, func() {
					MustMultiPolygonFromRingSlice(tt.args.slice)
				}, "MustMultiPolygonFromRingSlice(%v)", tt.args.slice)
			} else {
				assert.Equalf(t, tt.want, MustMultiPolygonFromRingSlice(tt.args.slice), "MustMultiPolygonFromRingSlice(%v)", tt.args.slice)
			}
		})
	}
}

func TestMultiPolygon_buildCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid input coordinates",
			input: []interface{}{
				[]interface{}{
					[]interface{}{
						[]interface{}{1.0, 1.0}, []interface{}{2.0, 2.0}, []interface{}{3.0, 3.0}, []interface{}{1.0, 1.0},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid coordinates type",
			input:   "invalid data type",
			wantErr: true,
		},
		{
			name: "invalid linear ring structure",
			input: []interface{}{
				[]interface{}{
					[]interface{}{
						[]interface{}{1.0, 1.0}, []interface{}{2.0, 2.0}, []interface{}{3.0, 3.0}, []interface{}{3.0, 2.0}, // Not closed
					},
				},
			},
			wantErr: true,
		},
		{
			name:    "empty coordinates",
			input:   []interface{}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPolygon{}
			err := m.buildCoordinates(tt.input)
			if tt.wantErr {
				require.Error(t, err, "Expected an error for input: %v", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %v", tt.input)
			}
		})
	}
}

func TestMultiPolygon_LinearRingsSlice(t *testing.T) {
	type fields struct {
		rings []LinearRings
	}
	tests := []struct {
		name   string
		fields fields
		want   []LinearRings
	}{
		{
			name: "empty rings",
			fields: fields{
				rings: []LinearRings{},
			},
			want: []LinearRings{},
		},
		{
			name: "single ring slice",
			fields: fields{
				rings: []LinearRings{
					{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
				},
			},
			want: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
			},
		},
		{
			name: "multiple rings slices",
			fields: fields{
				rings: []LinearRings{
					{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
					{{{4, 4}, {5, 5}, {6, 6}, {4, 4}}},
				},
			},
			want: []LinearRings{
				{{{1, 1}, {2, 2}, {3, 3}, {1, 1}}},
				{{{4, 4}, {5, 5}, {6, 6}, {4, 4}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPolygon{
				rings: tt.fields.rings,
			}
			assert.Equalf(t, tt.want, m.LinearRingsSlice(), "LinearRingsSlice()")
		})
	}
}
