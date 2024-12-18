package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject_Type(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected ObjectType
	}{
		{"defaultType", Object{}, TypeEmptyObject},
		{"customTypeFeature", Object{featureType: TypeFeature}, TypeFeature},
		{"customTypeCollection", Object{featureType: TypeFeatureCollection}, TypeFeatureCollection},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.object.Type())
		})
	}
}

func TestObject_IsFeature(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected bool
	}{
		{"isFeatureTrue", Object{featureType: TypeFeature}, true},
		{"isFeatureFalse", Object{featureType: TypeFeatureCollection}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.object.IsFeature())
		})
	}
}

func TestObject_IsFeatureCollection(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected bool
	}{
		{"isCollectionFalse", Object{featureType: TypeEmptyObject}, false},
		{"isCollectionTrue", Object{featureType: TypeFeatureCollection}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.object.IsFeatureCollection())
		})
	}
}

func TestObject_Feature(t *testing.T) {
	tests := []struct {
		name        string
		object      Object
		wantFeature *Feature
		wantErr     error
	}{
		{"validFeature", Object{featureType: TypeFeature, feature: &Feature{}}, &Feature{}, nil},
		{"invalidFeature", Object{featureType: TypeFeatureCollection}, nil, ErrInvalidFeature},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feature, err := tt.object.Feature()
			assert.Equal(t, tt.wantFeature, feature)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObject_FeatureCollection(t *testing.T) {
	tests := []struct {
		name         string
		object       Object
		wantFeatures *FeatureCollection
		wantErr      error
	}{
		{"validCollection", Object{featureType: TypeFeatureCollection, features: &FeatureCollection{}}, &FeatureCollection{}, nil},
		{"invalidCollection", Object{featureType: TypeEmptyObject}, nil, ErrInvalidFeature},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			features, err := tt.object.FeatureCollection()
			assert.Equal(t, tt.wantFeatures, features)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObject_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		object   Object
		expected string
		wantErr  error
	}{
		{"marshalFeature", Object{featureType: TypeFeature, feature: &Feature{}}, `{"type":"Feature","geometry":null}`, nil},
		{"marshalCollection", Object{featureType: TypeFeatureCollection, features: &FeatureCollection{}}, `{"type":"FeatureCollection","features":[]}`, nil},
		{"marshalInvalid", Object{}, "", ErrInvalidFeature},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := tt.object.MarshalJSON()
			assert.Equal(t, tt.expected, string(bytes))
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObject_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		inputJSON   string
		expected    Object
		expectedErr bool
	}{
		{
			"unmarshalFeature",
			`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,1]}}`,
			Object{
				featureType: TypeFeature,
				feature: &Feature{
					Geometry: MustPoint([]float64{1, 1}),
				},
			},
			false,
		},
		{
			"unmarshalCollection",
			`{"type":"FeatureCollection","features":[]}`,
			Object{
				featureType: TypeFeatureCollection,
				features:    &FeatureCollection{Features: []Feature{}},
			},
			false,
		},
		{
			"invalidType",
			`{"type":"InvalidType"}`,
			Object{},
			true,
		},
		{
			"invalidJSON",
			`invalid`,
			Object{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj Object
			err := json.Unmarshal([]byte(tt.inputJSON), &obj)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.featureType, obj.featureType)
				assert.Equal(t, tt.expected.feature, obj.feature)
				assert.Equal(t, tt.expected.features, obj.features)
			}
		})
	}
}
