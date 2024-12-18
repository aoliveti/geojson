package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiPoint_BoundingBox(t *testing.T) {
	tests := []struct {
		name        string
		vertices    Vertices
		expectedBox BoundingBox
	}{
		{"empty vertices", Vertices{}, BoundingBox{}},
		{"single point", Vertices{{1, 2}}, BoundingBox{1, 2, 1, 2}},
		{"multiple points 2D", Vertices{{1, 2}, {3, 4}, {0, 5}}, BoundingBox{0, 2, 3, 5}},
		{"multiple points 3D", Vertices{{1, 2, 3}, {4, 5, 6}, {0, 7, 1}}, BoundingBox{0, 2, 1, 4, 7, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPoint{vertices: tt.vertices}
			assert.Equal(t, tt.expectedBox, m.BoundingBox())
		})
	}
}

func TestMultiPoint_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		vertices Vertices
	}{
		{"empty vertices", Vertices{}},
		{"single point", Vertices{{1, 2}}},
		{"multiple points", Vertices{{1, 2}, {3, 4}, {5, 6}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPoint{vertices: tt.vertices}
			assert.Equal(t, tt.vertices, m.Vertices())
		})
	}
}

func TestMultiPoint_Type(t *testing.T) {
	m := &MultiPoint{}
	assert.Equal(t, TypeMultiPoint, m.Type())
}

func TestMultiPoint_buildCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		wantErr     error
		expectedRes Vertices
	}{
		{"invalid input type", "invalid", ErrInvalidCoordinates, nil},
		{"non-slice input", 123, ErrInvalidCoordinates, nil},
		{"valid single point", []interface{}{[]interface{}{1.0, 2.0}}, nil, Vertices{{1.0, 2.0}}},
		{"valid multiple points", []interface{}{
			[]interface{}{1.0, 2.0}, []interface{}{3.0, 4.0},
		}, nil, Vertices{{1.0, 2.0}, {3.0, 4.0}}},
		{"invalid point in slice", []interface{}{[]interface{}{1.0, 2.0}, "invalid"},
			ErrInvalidCoordinates, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPoint{}
			err := m.buildCoordinates(tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRes, m.vertices)
			}
		})
	}
}

func TestMultiPoint_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		multiPoint   *MultiPoint
		expectedJSON string
	}{
		{"empty vertices", &MultiPoint{
			vertices: Vertices{},
		}, `{"type":"MultiPoint","coordinates":[]}`},
		{"single point", &MultiPoint{
			vertices: Vertices{{1, 2}},
		}, `{"type":"MultiPoint","coordinates":[[1,2]]}`},
		{"with BBox", &MultiPoint{
			vertices:      Vertices{{1, 2}, {3, 4}},
			SerializeBBox: true,
		}, `{"type":"MultiPoint","coordinates":[[1,2],[3,4]],"bbox":[1,2,3,4]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.multiPoint.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(got))
		})
	}
}

func TestMultiPoint_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		hasError    bool
		expectedRes Vertices
	}{
		{"invalid coordinates", `{"type":"MultiPoint","coordinates":[[1,2],["3",4]]}`, true, nil},
		{"valid input", `{"type":"MultiPoint","coordinates":[[1,2],[3,4]]}`, false, Vertices{{1, 2}, {3, 4}}},
		{"invalid type", `{"type":"Point","coordinates":[1,2]}`, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiPoint{}
			err := json.Unmarshal([]byte(tt.jsonData), m)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRes, m.vertices)
			}
		})
	}
}

func TestNewMultiPointFromVertices(t *testing.T) {
	tests := []struct {
		name     string
		vertices Vertices
	}{
		{"empty vertices", Vertices{}},
		{"single vertex", Vertices{{1, 2}}},
		{"multiple vertices", Vertices{{1, 2}, {3, 4}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMultiPointFromVertices(tt.vertices)
			assert.Equal(t, tt.vertices, got.vertices)
		})
	}
}
