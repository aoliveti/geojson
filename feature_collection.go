package geojson

import (
	"encoding/json"
	"fmt"
)

// FeatureCollection represents a GeoJSON object containing a collection of Features.
type FeatureCollection struct {
	Features      []Feature // Features contains the list of features in the collection.
	SerializeBBox bool      // SerializeBBox determines whether to include the bounding box in the serialized JSON.
}

// BoundingBox calculates and returns the bounding box for all features in the collection.
func (f *FeatureCollection) BoundingBox() BoundingBox {
	return bbox(f.Vertices())
}

// Vertices extracts and returns all vertices from all features in the collection.
func (f *FeatureCollection) Vertices() Vertices {
	var v Vertices
	for _, f := range f.Features {
		v = append(v, f.Vertices()...)
	}
	return v
}

// MarshalJSON serializes the FeatureCollection into GeoJSON format.
// If SerializeBBox is true, it includes the bounding box in the serialized JSON.
func (f *FeatureCollection) MarshalJSON() ([]byte, error) {
	features := f.Features
	if features == nil {
		features = make([]Feature, 0)
	}

	fjc := featureCollectionJSONOutput{
		Type:     TypeFeatureCollection,
		Features: features,
	}

	if f.SerializeBBox {
		fjc.BBox = f.BoundingBox()
	}

	return json.Marshal(&fjc)
}

// UnmarshalJSON deserializes GeoJSON data into a FeatureCollection object.
// Returns an error if the input data cannot be unmarshaled.
func (f *FeatureCollection) UnmarshalJSON(bytes []byte) error {
	few := &Object{}
	err := json.Unmarshal(bytes, few)
	if err != nil {
		return fmt.Errorf("failed to unmarshal feature collection: %w", err)
	}

	*f = *few.features

	return nil
}

// NewFeatureCollection creates and returns a new, empty FeatureCollection.
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{}
}

// NewFeatureCollectionFromFeatures creates and returns a new FeatureCollection
// initialized with the provided slice of features.
func NewFeatureCollectionFromFeatures(features []Feature) *FeatureCollection {
	return &FeatureCollection{
		Features: features,
	}
}
