package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinates_Longitude(t *testing.T) {
	tests := []struct {
		name     string
		input    Coordinates
		expected float64
	}{
		{"valid longitude", Coordinates{12.34, 56.78}, 12.34},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Longitude())
		})
	}
}

func TestCoordinates_Latitude(t *testing.T) {
	tests := []struct {
		name     string
		input    Coordinates
		expected float64
	}{
		{"valid latitude", Coordinates{12.34, 56.78}, 56.78},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Latitude())
		})
	}
}

func TestCoordinates_HasAltitude(t *testing.T) {
	tests := []struct {
		name     string
		input    Coordinates
		expected bool
	}{
		{"no altitude", Coordinates{12.34, 56.78}, false},
		{"has altitude", Coordinates{12.34, 56.78, 100.0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.HasAltitude())
		})
	}
}

func TestCoordinates_Altitude(t *testing.T) {
	tests := []struct {
		name     string
		input    Coordinates
		expected float64
	}{
		{"valid altitude", Coordinates{12.34, 56.78, 100.0}, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.Altitude())
		})
	}
}

func TestCoordinates_IsEqual(t *testing.T) {
	tests := []struct {
		name     string
		c1       Coordinates
		c2       Coordinates
		expected bool
	}{
		{"identical 2D coordinates", Coordinates{12.34, 56.78}, Coordinates{12.34, 56.78}, true},
		{"identical 3D coordinates", Coordinates{12.34, 56.78, 100.0}, Coordinates{12.34, 56.78, 100.0}, true},
		{"different longitudes", Coordinates{12.34, 56.78}, Coordinates{56.78, 56.78}, false},
		{"different latitudes", Coordinates{12.34, 56.78}, Coordinates{12.34, 12.34}, false},
		{"different altitudes", Coordinates{12.34, 56.78, 100.0}, Coordinates{12.34, 56.78, 200.0}, false},
		{"3D vs 2D comparison", Coordinates{12.34, 56.78, 100.0}, Coordinates{12.34, 56.78}, false},
		{"empty vs non-empty", Coordinates{}, Coordinates{12.34, 56.78}, false},
		{"both empty", Coordinates{}, Coordinates{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.c1.IsEqual(tt.c2))
		})
	}
}

func TestCoordinates_String(t *testing.T) {
	tests := []struct {
		name     string
		input    Coordinates
		expected string
	}{
		{"no altitude", Coordinates{12.34, 56.78}, "[ 12.34, 56.78 ]"},
		{"has altitude", Coordinates{12.34, 56.78, 100.0}, "[ 12.34, 56.78, 100 ]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}

func TestCoordinates_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
		expected  Coordinates
	}{
		{"empty JSON", ``, true, nil},
		{"invalid JSON format", `invalid`, true, nil},
		{"empty coordinates", `[]`, true, Coordinates{}},
		{"valid 2D coordinates", `[12.34, 56.78]`, false, Coordinates{12.34, 56.78}},
		{"valid 3D coordinates", `[12.34, 56.78, 100.0]`, false, Coordinates{12.34, 56.78, 100.0}},
		{"invalid coordinate count", `[12.34]`, true, nil},
		{"out of range longitude", `[200.0, 56.78]`, true, nil},
		{"out of range latitude", `[12.34, 100.0]`, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Coordinates
			err := json.Unmarshal([]byte(tt.input), &c)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, c)
			}
		})
	}
}

func TestNewCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		input     []float64
		expectErr bool
		expected  Coordinates
	}{
		{"empty coordinates", []float64{}, true, nil},
		{"valid 2D coordinates", []float64{12.34, 56.78}, false, Coordinates{12.34, 56.78}},
		{"valid 3D coordinates", []float64{12.34, 56.78, 100.0}, false, Coordinates{12.34, 56.78, 100.0}},
		{"invalid coordinate count", []float64{12.34}, true, nil},
		{"out of range longitude", []float64{200.0, 56.78}, true, nil},
		{"out of range latitude", []float64{12.34, 100.0}, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coords, err := NewCoordinates(tt.input)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, *coords)
			}
		})
	}
}

func TestMustCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		input       []float64
		expectPanic bool
	}{
		{"valid 2D coordinates", []float64{12.34, 56.78}, false},
		{"valid 3D coordinates", []float64{12.34, 56.78, 100.0}, false},
		{"invalid coordinate count", []float64{12.34}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					_ = MustCoordinates(tt.input)
				})
			} else {
				assert.NotPanics(t, func() {
					_ = MustCoordinates(tt.input)
				})
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		latitude  float64
		expectErr bool
	}{
		{"valid coordinates", 12.34, 56.78, false},
		{"invalid longitude", 200.0, 56.78, true},
		{"invalid latitude", 12.34, -100.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinates(tt.longitude, tt.latitude)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr bool
		expected  Coordinates
	}{
		{"valid 2D coordinates", []interface{}{12.34, 56.78}, false, Coordinates{12.34, 56.78}},
		{"valid 3D coordinates", []interface{}{12.34, 56.78, 100.0}, false, Coordinates{12.34, 56.78, 100.0}},
		{"invalid coordinate size", []interface{}{12.34}, true, nil},
		{"non-numeric input", []interface{}{"12.34", 56.78}, true, nil},
		{"invalid longitude range", []interface{}{200.0, 56.78}, true, nil},
		{"invalid latitude range", []interface{}{12.34, 100.0}, true, nil},
		{"non-slice input", map[string]interface{}{"lng": 12.34, "lat": 56.78}, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coords, err := buildCoordinates(tt.input)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, &tt.expected, coords)
			}
		})
	}
}
