package geojson

import (
	"encoding/json"
	"fmt"
)

// Feature represents a GeoJSON feature with a geometry, properties, an optional ID, and bounding box toggling.
type Feature struct {
	Geometry      Geometry   // Geometry specifies the spatial information of the feature.
	Properties    Properties // Properties contains supplementary data about the feature.
	ID            *ID        // ID is an optional identifier for the feature.
	SerializeBBox bool       // SerializeBBox determines whether to include the bounding box in the serialized JSON.
}

// BoundingBox calculates and returns the bounding box for the feature's geometry.
func (f *Feature) BoundingBox() BoundingBox {
	return bbox(f.Vertices())
}

// Vertices extracts and returns all vertices present in the feature's geometry.
func (f *Feature) Vertices() Vertices {
	if f.Geometry == nil {
		return nil
	}

	var v Vertices
	v = append(v, f.Geometry.Vertices()...)
	return v
}

// GeometryObject converts the Feature's geometry into a GeometryObject.
func (f *Feature) GeometryObject() GeometryObject {
	return GeometryObject{
		geometry: f.Geometry,
	}
}

// UnmarshalJSON deserializes GeoJSON data into a Feature object.
func (f *Feature) UnmarshalJSON(bytes []byte) error {
	few := &Object{}
	if err := json.Unmarshal(bytes, few); err != nil {
		return fmt.Errorf("failed to unmarshal feature: %w", err)
	}

	if few.feature == nil {
		return ErrInvalidFeature
	}

	f.Geometry = few.feature.Geometry
	f.Properties = few.feature.Properties
	f.ID = few.feature.ID

	return nil
}

// MarshalJSON serializes a Feature object into GeoJSON format.
func (f *Feature) MarshalJSON() ([]byte, error) {
	fj := &featureJSONOutput{
		Type:       TypeFeature,
		Geometry:   f.Geometry,
		Properties: f.Properties,
		ID:         f.ID,
	}

	if f.SerializeBBox {
		fj.BBox = f.BoundingBox()
	}

	return json.Marshal(fj)
}

// FeatureBuilder is a builder for constructing Feature objects.
type FeatureBuilder struct {
	feature Feature // feature holds the Feature object being constructed.
}

// NewFeatureBuilder creates and returns a new instance of FeatureBuilder.
func NewFeatureBuilder() *FeatureBuilder {
	return &FeatureBuilder{}
}

// Build finalizes and returns the constructed Feature object.
func (fb *FeatureBuilder) Build() Feature {
	return fb.feature
}

// SetGeometry sets the geometry of the feature and returns the builder.
func (fb *FeatureBuilder) SetGeometry(geometry Geometry) *FeatureBuilder {
	fb.feature.Geometry = geometry
	return fb
}

// SetProperties sets the properties of the feature and returns the builder.
func (fb *FeatureBuilder) SetProperties(properties Properties) *FeatureBuilder {
	fb.feature.Properties = properties
	return fb
}

// SetID sets the ID of the feature and returns the builder.
func (fb *FeatureBuilder) SetID(id ID) *FeatureBuilder {
	fb.feature.ID = &id
	return fb
}
