package geojson

import (
	"encoding/json"
	"fmt"
)

var (
	// ErrInvalidFeature is returned when an invalid feature type or format is encountered.
	ErrInvalidFeature = fmt.Errorf("invalid feature type or format")
)

// Object represents a GeoJSON object, which can be either a Feature or a FeatureCollection.
type Object struct {
	featureType ObjectType         // The type of GeoJSON object (e.g., "Feature" or "FeatureCollection").
	feature     *Feature           // The single Feature represented by the object, if applicable.
	features    *FeatureCollection // The FeatureCollection represented by the object, if applicable.
}

// Type returns the type of the GeoJSON object.
func (o *Object) Type() ObjectType {
	if o.featureType == "" {
		return TypeEmptyObject
	}

	return o.featureType
}

// IsFeature checks if the Object is a single Feature.
func (o *Object) IsFeature() bool {
	return o.Type() == TypeFeature
}

// IsFeatureCollection checks if the Object is a FeatureCollection.
func (o *Object) IsFeatureCollection() bool {
	return o.Type() == TypeFeatureCollection
}

// Feature retrieves the single Feature from the Object.
// Returns an error if the Object is not a single Feature.
func (o *Object) Feature() (*Feature, error) {
	if o.IsFeature() {
		return o.feature, nil
	}

	return nil, ErrInvalidFeature
}

// FeatureCollection retrieves the FeatureCollection from the Object.
// Returns an error if the Object is not a FeatureCollection.
func (o *Object) FeatureCollection() (*FeatureCollection, error) {
	if o.IsFeatureCollection() {
		return o.features, nil
	}

	return nil, ErrInvalidFeature
}

// MarshalJSON encodes the Object into JSON.
// This function chooses the appropriate representation based on the Object type.
func (o *Object) MarshalJSON() ([]byte, error) {
	switch o.Type() {
	case TypeFeature:
		return json.Marshal(o.feature)
	case TypeFeatureCollection:
		return json.Marshal(o.features)
	default:
		return nil, ErrInvalidFeature
	}
}

// UnmarshalJSON decodes JSON data into the Object.
// Identifies if the Object is a single Feature or a FeatureCollection, and unmarshals accordingly.
func (o *Object) UnmarshalJSON(bytes []byte) error {
	var feature featuresJSONInput
	if err := json.Unmarshal(bytes, &feature); err != nil {
		return fmt.Errorf("failed to unmarshal features: %w", err)
	}

	if feature.Geometry == nil {
		feature.Geometry = &GeometryObject{}
	}

	switch feature.Type {
	case TypeFeature:
		o.feature = &Feature{
			Geometry:   feature.Geometry.geometry,
			Properties: feature.Properties,
			ID:         feature.ID,
		}
	case TypeFeatureCollection:
		v := NewFeatureCollectionFromFeatures(feature.Features)
		o.features = v
	default:
		return ErrInvalidFeature
	}

	o.featureType = feature.Type

	return nil
}
