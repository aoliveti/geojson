package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		coords   Coordinates
		expected BoundingBox
	}{
		{
			name:     "valid bounding box",
			coords:   Coordinates{1.0, 2.0},
			expected: BoundingBox{1.0, 2.0, 1.0, 2.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			result := point.BoundingBox()
			assert.True(t, result.IsValid())
			assert.Equal(t, len(tt.expected), len(result))
		})
	}
}

func TestVertices(t *testing.T) {
	tests := []struct {
		name     string
		coords   Coordinates
		expected Vertices
	}{
		{
			name:     "valid vertices",
			coords:   Coordinates{1.0, 2.0},
			expected: Vertices{{1.0, 2.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			result := point.Vertices()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLongitude(t *testing.T) {
	tests := []struct {
		name     string
		coords   Coordinates
		expected float64
	}{
		{
			name:     "get longitude",
			coords:   Coordinates{1.0, 2.0},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			result := point.Longitude()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLatitude(t *testing.T) {
	tests := []struct {
		name     string
		coords   Coordinates
		expected float64
	}{
		{
			name:     "get latitude",
			coords:   Coordinates{1.0, 2.0},
			expected: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			result := point.Latitude()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasAltitude(t *testing.T) {
	tests := []struct {
		name     string
		coords   Coordinates
		expected bool
	}{
		{
			name:     "has altitude",
			coords:   Coordinates{1.0, 2.0, 3.0},
			expected: true,
		},
		{
			name:     "no altitude",
			coords:   Coordinates{1.0, 2.0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			result := point.HasAltitude()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAltitude(t *testing.T) {
	tests := []struct {
		name        string
		coords      Coordinates
		expected    float64
		expectPanic bool
	}{
		{
			name:        "get altitude",
			coords:      Coordinates{1.0, 2.0, 3.0},
			expected:    3.0,
			expectPanic: false,
		},
		{
			name:        "no altitude",
			coords:      Coordinates{1.0, 2.0},
			expected:    0.0,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{coords: tt.coords}
			if tt.expectPanic {
				assert.Panics(t, func() { point.Altitude() })
			} else {
				assert.NotPanics(t, func() {
					result := point.Altitude()
					assert.Equal(t, tt.expected, result)
				})
			}
		})
	}
}

func TestType(t *testing.T) {
	point := &Point{}
	expected := TypePoint
	result := point.Type()
	assert.Equal(t, expected, result)
}

func TestPoint_BuildCoordinates(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		hasError bool
	}{
		{
			name:     "valid coordinates",
			input:    []interface{}{1.0, 2.0},
			hasError: false,
		},
		{
			name:     "valid coordinates size",
			input:    []interface{}{1.0},
			hasError: true,
		},
		{
			name:     "invalid coordinates type",
			input:    []interface{}{"1.0", 2.0},
			hasError: true,
		},
		{
			name:     "non-slice input",
			input:    1.0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{}
			err := point.buildCoordinates(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		point    *Point
		expected string
	}{
		{
			name:     "without bbox",
			point:    &Point{coords: Coordinates{1.0, 2.0}, SerializeBBox: false},
			expected: `{"type":"Point","coordinates":[1.0,2.0]}`,
		},
		{
			name:     "with bbox",
			point:    &Point{coords: Coordinates{1.0, 2.0}, SerializeBBox: true},
			expected: `{"type":"Point","bbox":[1.0,2.0,1.0,2.0],"coordinates":[1.0,2.0]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.point.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Coordinates
		hasError bool
	}{
		{
			name:     "valid input",
			input:    `{"type":"Point","coordinates":[1.0,2.0]}`,
			expected: Coordinates{1.0, 2.0},
			hasError: false,
		},
		{
			name:     "invalid type",
			input:    `{"type":"LineString","coordinates":[[1.0,2.0],[3.0,4.0]]}`,
			hasError: true,
		},
		{
			name:     "invalid coordinates",
			input:    `{"type":"Point","coordinates":["1.0",2.0]}`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &Point{}
			err := point.UnmarshalJSON([]byte(tt.input))
			if tt.hasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, point.coords)
			}
		})
	}
}

func TestNewPoint(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		hasError bool
	}{
		{
			name:     "valid coordinates",
			input:    []float64{1.0, 2.0},
			hasError: false,
		},
		{
			name:     "invalid coordinates",
			input:    []float64{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point, err := NewPoint(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, point)
			}
		})
	}
}

func TestMustPoint(t *testing.T) {
	tests := []struct {
		name        string
		input       []float64
		shouldPanic bool
	}{
		{
			name:        "valid coordinates",
			input:       []float64{1.0, 2.0},
			shouldPanic: false,
		},
		{
			name:        "invalid coordinates",
			input:       []float64{},
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { MustPoint(tt.input) })
			} else {
				assert.NotPanics(t, func() { MustPoint(tt.input) })
			}
		})
	}
}

func TestPoint_Coordinates(t *testing.T) {
	type fields struct {
		coords        Coordinates
		SerializeBBox bool
	}
	tests := []struct {
		name   string
		fields fields
		want   Coordinates
	}{
		{
			name: "2D coordinates",
			fields: fields{
				coords: Coordinates{1.0, 2.0},
			},
			want: Coordinates{1.0, 2.0},
		},
		{
			name: "3D coordinates",
			fields: fields{
				coords: Coordinates{1.0, 2.0, 3.0},
			},
			want: Coordinates{1.0, 2.0, 3.0},
		},
		{
			name: "empty coordinates",
			fields: fields{
				coords: Coordinates{},
			},
			want: Coordinates{},
		},
		{
			name:   "nil",
			fields: fields{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{
				coords:        tt.fields.coords,
				SerializeBBox: tt.fields.SerializeBBox,
			}
			assert.Equalf(t, tt.want, p.Coordinates(), "Coordinates()")
		})
	}
}
