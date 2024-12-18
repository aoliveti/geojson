package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureCollection_BoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		features []Feature
		expected BoundingBox
	}{
		{
			name:     "empty features",
			features: []Feature{},
			expected: BoundingBox{},
		},
		{
			name: "single feature",
			features: []Feature{
				{Geometry: MustPoint([]float64{1.0, 2.0})},
			},
			expected: BoundingBox{1.0, 2.0, 1.0, 2.0},
		},
		{
			name: "multiple features",
			features: []Feature{
				{Geometry: MustPoint([]float64{1.0, 2.0})},
				{Geometry: MustPoint([]float64{3.0, 4.0})},
			},
			expected: BoundingBox{1.0, 2.0, 3.0, 4.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewFeatureCollectionFromFeatures(tt.features)
			result := fc.BoundingBox()
			assert.True(t, result.IsValid(), "bounding box should be valid")
			assert.Equal(t, tt.expected, result, "bounding box mismatch")
		})
	}
}

func TestFeatureCollection_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		features []Feature
		expected Vertices
	}{
		{
			name:     "empty features",
			features: []Feature{},
			expected: nil,
		},
		{
			name: "single feature",
			features: []Feature{
				{Geometry: MustPoint([]float64{1.0, 2.0})},
			},
			expected: Vertices{*MustCoordinates([]float64{1.0, 2.0})},
		},
		{
			name: "multiple features",
			features: []Feature{
				{Geometry: MustPoint([]float64{1.0, 2.0})},
				{Geometry: MustPoint([]float64{3.0, 4.0})},
			},
			expected: Vertices{
				*MustCoordinates([]float64{1.0, 2.0}),
				*MustCoordinates([]float64{3.0, 4.0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewFeatureCollectionFromFeatures(tt.features)
			result := fc.Vertices()
			assert.Equal(t, tt.expected, result, "vertices mismatch")
		})
	}
}

func TestFeatureCollection_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		featureColl  *FeatureCollection
		expectedJSON string
		expectError  bool
	}{
		{
			name: "with bounding box",
			featureColl: &FeatureCollection{
				Features: []Feature{
					{Geometry: MustPoint([]float64{1.0, 2.0})},
				},
				SerializeBBox: true,
			},
			expectedJSON: `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}],"bbox":[1,2,1,2]}`,
			expectError:  false,
		},
		{
			name: "without bounding box",
			featureColl: &FeatureCollection{
				Features: []Feature{
					{Geometry: MustPoint([]float64{1.0, 2.0})},
				},
				SerializeBBox: false,
			},
			expectedJSON: `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}]}`,
			expectError:  false,
		},
		{
			name:         "empty FeatureCollection",
			featureColl:  &FeatureCollection{},
			expectedJSON: `{"type":"FeatureCollection","features":[]}`,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := tt.featureColl.MarshalJSON()
			if tt.expectError {
				assert.Error(t, err, "marshal should return an error")
			} else {
				require.NoError(t, err, "marshal should not return an error")
				assert.JSONEq(t, tt.expectedJSON, string(marshaled), "JSON output mismatch")
			}
		})
	}
}

func TestFeatureCollection_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid input",
			input:       `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}]}`,
			expectError: false,
		},
		{
			name:        "invalid input",
			input:       `{"invalid":"data"}`,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fc FeatureCollection
			err := fc.UnmarshalJSON([]byte(tt.input))
			if tt.expectError {
				assert.Error(t, err, "unmarshal should return an error")
			} else {
				require.NoError(t, err, "unmarshal should not return an error")
			}
		})
	}
}

func TestNewFeatureCollection(t *testing.T) {
	fc := NewFeatureCollection()
	assert.Empty(t, fc.Features, "features should be empty")
	assert.False(t, fc.SerializeBBox, "SerializeBBox should be false by default")
}

func TestNewFeatureCollectionFromFeatures(t *testing.T) {
	features := []Feature{
		{Geometry: MustPoint([]float64{1.0, 2.0})},
		{Geometry: MustPoint([]float64{3.0, 4.0})},
	}
	fc := NewFeatureCollectionFromFeatures(features)
	assert.Equal(t, features, fc.Features, "features mismatch")
}
