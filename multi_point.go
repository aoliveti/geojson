package geojson

import (
	"encoding/json"
	"fmt"
)

// MultiPoint represents a GeoJSON MultiPoint geometry.
type MultiPoint struct {
	vertices      Vertices // The vertices of the MultiPoint geometry.
	SerializeBBox bool     // Indicates whether to serialize the bounding box.
}

// BoundingBox calculates and returns the bounding box of the MultiPoint geometry.
func (m *MultiPoint) BoundingBox() BoundingBox {
	return bbox(m.Vertices())
}

// Vertices returns the vertices of the MultiPoint geometry.
func (m *MultiPoint) Vertices() Vertices {
	return m.vertices
}

// Type returns the GeoJSON type of the geometry, which is MultiPoint.
func (m *MultiPoint) Type() GeometryType {
	return TypeMultiPoint
}

// buildCoordinates populates the MultiPoint with vertices from the provided raw data.
// It returns an error if the input is invalid.
func (m *MultiPoint) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	vertices := make(Vertices, len(rawSlice))
	for i, s := range rawSlice {
		p := Point{}
		if err := p.buildCoordinates(s); err != nil {
			return err
		}

		vertices[i] = p.coords
	}

	m.vertices = vertices
	return nil
}

// MarshalJSON serializes the MultiPoint to its GeoJSON representation.
// If SerializeBBox is true, the bounding box is included in the output.
func (m *MultiPoint) MarshalJSON() ([]byte, error) {
	vertices := m.vertices
	if len(vertices) == 0 {
		vertices = make(Vertices, 0)
	}

	out := geometryJSONOutput{
		Type:        m.Type(),
		Coordinates: vertices,
	}

	if m.SerializeBBox {
		out.BBox = m.BoundingBox()
	}

	return json.Marshal(&out)
}

// UnmarshalJSON deserializes the GeoJSON representation of a MultiPoint.
// It returns an error if the input data is not valid or doesn't match a MultiPoint.
func (m *MultiPoint) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", m.Type(), err)
	}

	g, ok := gw.geometry.(*MultiPoint)
	if !ok {
		return ErrInvalidTypeField
	}

	m.vertices = g.vertices
	return nil
}

// NewMultiPointFromVertices creates and returns a new MultiPoint from the given vertices.
func NewMultiPointFromVertices(vertices Vertices) *MultiPoint {
	return &MultiPoint{
		vertices: vertices,
	}
}
