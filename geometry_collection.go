package geojson

import (
	"encoding/json"
	"fmt"
)

var (
	// ErrGeometryCollectionBuildCoordinates is returned when attempting to build coordinates
	// for a GeometryCollection, which does not directly define coordinates.
	ErrGeometryCollectionBuildCoordinates = fmt.Errorf("%s does not have coordinates to build", TypeGeometryCollection)
)

// GeometryCollection represents a GeoJSON GeometryCollection,
// which is a collection of different geometry objects.
type GeometryCollection struct {
	geometries []Geometry // Slice of Geometry objects contained within the collection.
}

// BoundingBox calculates and returns the BoundingBox for the entire GeometryCollection.
// It computes the bounding box by aggregating the bounding boxes of all geometries in the collection.
func (g *GeometryCollection) BoundingBox() BoundingBox {
	return bbox(g.Vertices())
}

// Vertices aggregates and returns all the vertices from all geometries in the collection.
// This is used for operations like calculating the bounding box of the collection.
func (g *GeometryCollection) Vertices() Vertices {
	var v Vertices
	for _, g := range g.geometries {
		v = append(v, g.Vertices()...)
	}
	return v
}

// Type returns the GeoJSON type for the GeometryCollection, which is "GeometryCollection".
func (g *GeometryCollection) Type() GeometryType {
	return TypeGeometryCollection
}

// Geometries returns the slice of Geometry objects contained in the GeometryCollection.
// It provides access to the individual geometries that make up the collection.
func (g *GeometryCollection) Geometries() []Geometry {
	return g.geometries
}

// MarshalJSON serializes the GeometryCollection into GeoJSON format.
// It outputs the type as "GeometryCollection" and includes child geometries, if any.
func (g *GeometryCollection) MarshalJSON() ([]byte, error) {
	geometries := make([]Geometry, 0)
	if len(g.geometries) > 0 {
		geometries = g.geometries
	}

	out := geometryCollectionJSONOutput{
		Type:       g.Type(),
		Geometries: geometries,
	}

	return json.Marshal(&out)
}

// UnmarshalJSON deserializes the given GeoJSON data into a GeometryCollection.
// It first unmarshals the data into a generic GeometryObject, validates its type,
// and assigns the parsed geometries to the collection.
func (g *GeometryCollection) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", g.Type(), err)
	}

	gc, ok := gw.geometry.(*GeometryCollection)
	if !ok {
		return ErrInvalidTypeField
	}

	g.geometries = gc.geometries

	return nil
}

// buildCoordinates returns an error because GeometryCollection does not directly define coordinates.
// This satisfies the Geometry interface but is unsupported for GeometryCollection.
func (g *GeometryCollection) buildCoordinates(_ interface{}) error {
	return ErrGeometryCollectionBuildCoordinates
}

// NewGeometryCollection creates and returns an empty GeometryCollection.
// This is useful for initializing a collection before adding geometries.
func NewGeometryCollection() *GeometryCollection {
	return &GeometryCollection{}
}

// NewGeometryCollectionFromSlice creates and returns a GeometryCollection initialized with the given geometries.
// This allows creating a collection pre-filled with specific Geometry objects.
func NewGeometryCollectionFromSlice(geometries []Geometry) *GeometryCollection {
	return &GeometryCollection{
		geometries: geometries,
	}
}
