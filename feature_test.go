package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeature_BoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		feature  Feature
		expected BoundingBox
	}{
		{
			name:     "nil geometry",
			feature:  Feature{},
			expected: BoundingBox{},
		},
		{
			name: "Point geometry",
			feature: Feature{
				Geometry: MustPoint(Coordinates{1.0, 2.0}),
			},
			expected: BoundingBox{1.0, 2.0, 1.0, 2.0},
		},
		{
			name: "LineString geometry",
			feature: Feature{
				Geometry: MustLineString([]Coordinates{
					{1.0, 2.0}, {3.0, 4.0},
				}),
			},
			expected: BoundingBox{1.0, 2.0, 3.0, 4.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.feature.BoundingBox()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFeature_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		feature  Feature
		expected Vertices
	}{
		{
			name:     "nil geometry",
			feature:  Feature{},
			expected: nil,
		},
		{
			name: "Point geometry",
			feature: Feature{
				Geometry: MustPoint(Coordinates{1.0, 2.0}),
			},
			expected: Vertices{{1.0, 2.0}},
		},
		{
			name: "LineString geometry",
			feature: Feature{
				Geometry: MustLineString([]Coordinates{
					{1.0, 2.0}, {3.0, 4.0},
				}),
			},
			expected: Vertices{{1.0, 2.0}, {3.0, 4.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.feature.Vertices()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFeature_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		feature     Feature
		expected    string
		expectError bool
	}{
		{
			name: "basic feature serialization",
			feature: Feature{
				Properties: Properties{
					"name": "test",
				},
				SerializeBBox: false,
			},
			expected:    `{"type":"Feature","geometry":null,"properties":{"name":"test"}}`,
			expectError: false,
		},
		{
			name: "basic feature serialization with ID",
			feature: Feature{
				Properties: Properties{
					"name": "test",
				},
				ID:            NewNumericID(1),
				SerializeBBox: false,
			},
			expected:    `{"type":"Feature","geometry":null,"properties":{"name":"test"}, "id":1}`,
			expectError: false,
		},
		{
			name: "Point geometry with BoundingBox",
			feature: Feature{
				Geometry:      MustPoint(Coordinates{1.0, 2.0}),
				SerializeBBox: true,
			},
			expected:    `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"bbox":[1,2,1,2]}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.feature.MarshalJSON()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.expected, string(result))
			}
		})
	}
}

func TestFeature_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		expectError bool
		validate    func(*Feature)
	}{
		{
			name:        "invalid JSON",
			jsonInput:   `invalid-json`,
			expectError: true,
		},
		{
			name:        "invalid feature type",
			jsonInput:   `{"type":"Invalid","geometry":null}`,
			expectError: true,
		},
		{
			name:        "different feature type",
			jsonInput:   `{"type":"FeatureCollection","geometry":null}`,
			expectError: true,
		},
		{
			name:        "valid Feature with Point geometry",
			jsonInput:   `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{"name":"test"}}`,
			expectError: false,
			validate: func(f *Feature) {
				assert.NotNil(t, f.Geometry)
				assert.Equal(t, MustPoint(Coordinates{1.0, 2.0}), f.Geometry)
				value, ok := f.Properties.Get("name")
				assert.True(t, ok)
				assert.Equal(t, "test", value)
			},
		},
		{
			name:        "valid Feature with Properties and ID",
			jsonInput:   `{"type":"Feature","properties":{"name":"test"},"id":"test"}`,
			expectError: false,
			validate: func(f *Feature) {
				assert.Nil(t, f.Geometry)
				assert.NotNil(t, f.ID)

				idValue, ok := f.ID.StringValue()
				assert.True(t, ok)
				assert.Equal(t, "test", idValue)

				value, ok := f.Properties.Get("name")
				assert.True(t, ok)
				assert.Equal(t, "test", value)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var feature Feature
			err := feature.UnmarshalJSON([]byte(tt.jsonInput))
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.validate != nil {
					tt.validate(&feature)
				}
			}
		})
	}
}

func TestFeatureBuilder(t *testing.T) {
	t.Run("Build empty feature", func(t *testing.T) {
		builder := NewFeatureBuilder()
		feature := builder.Build()
		assert.Nil(t, feature.Geometry)
		assert.Nil(t, feature.Properties)
		assert.Nil(t, feature.ID)
	})

	t.Run("SetGeometry", func(t *testing.T) {
		geometry := MustPoint(Coordinates{1.0, 2.0})
		builder := NewFeatureBuilder()
		feature := builder.SetGeometry(geometry).Build()
		assert.Equal(t, geometry, feature.Geometry)
	})

	t.Run("SetProperties", func(t *testing.T) {
		properties := Properties{"key": "value"}
		builder := NewFeatureBuilder()
		feature := builder.SetProperties(properties).Build()
		assert.Equal(t, properties, feature.Properties)
	})

	t.Run("SetID", func(t *testing.T) {
		id := NewStringID("1234")
		builder := NewFeatureBuilder()
		feature := builder.SetID(*id).Build()
		assert.Equal(t, id, feature.ID)
	})
}

func TestFeature_GeometryObject(t *testing.T) {
	type fields struct {
		Geometry      Geometry
		Properties    Properties
		ID            *ID
		SerializeBBox bool
	}
	tests := []struct {
		name   string
		fields fields
		want   GeometryObject
	}{
		{
			name:   "nil geometry",
			fields: fields{Geometry: nil},
			want:   GeometryObject{geometry: nil},
		},
		{
			name:   "Point geometry",
			fields: fields{Geometry: MustPoint(Coordinates{1.0, 2.0})},
			want:   GeometryObject{geometry: MustPoint(Coordinates{1.0, 2.0})},
		},
		{
			name:   "LineString geometry",
			fields: fields{Geometry: MustLineString([]Coordinates{{1.0, 2.0}, {3.0, 4.0}})},
			want:   GeometryObject{geometry: MustLineString([]Coordinates{{1.0, 2.0}, {3.0, 4.0}})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feature{
				Geometry:      tt.fields.Geometry,
				Properties:    tt.fields.Properties,
				ID:            tt.fields.ID,
				SerializeBBox: tt.fields.SerializeBBox,
			}
			assert.Equalf(t, tt.want, f.GeometryObject(), "GeometryObject()")
		})
	}
}
