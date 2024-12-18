package geojson

import (
	"encoding/json"
	"fmt"
)

const (
	// LineStringMinimumSize defines the minimum number of vertices required for a LineString.
	LineStringMinimumSize = 2
)

var (
	// ErrLineStringTooShort indicates that a LineString must have at least 2 vertices.
	ErrLineStringTooShort = fmt.Errorf("line string must have at least 2 vertices")
)

// LineString represents a GeoJSON LineString geometry, defined by a series of vertices.
type LineString struct {
	vertices      Vertices // Vertices that define the LineString.
	SerializeBBox bool     // Whether to include a bounding box in the JSON serialization.
}

// Type returns the type of the geometry, which is always TypeLineString for LineString.
func (l *LineString) Type() GeometryType {
	return TypeLineString
}

// Vertices returns the Vertices of the LineString.
func (l *LineString) Vertices() Vertices {
	return l.vertices
}

// BoundingBox calculates the bounding box for the LineString.
func (l *LineString) BoundingBox() BoundingBox {
	return bbox(l.Vertices())
}

// buildCoordinates constructs the LineString's vertices from the provided raw data.
// Returns an error if the input is invalid or the number of coordinates is less than the minimum required.
func (l *LineString) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	if len(rawSlice) < LineStringMinimumSize {
		return ErrLineStringTooShort
	}

	vertices := make(Vertices, len(rawSlice))
	for i, s := range rawSlice {
		p := Point{}
		if err := p.buildCoordinates(s); err != nil {
			return err
		}

		vertices[i] = p.coords
	}

	l.vertices = vertices

	return nil
}

// MarshalJSON serializes the LineString as GeoJSON.
// It includes the bounding box (if SerializeBBox is true) and the vertices.
func (l *LineString) MarshalJSON() ([]byte, error) {
	vertices := l.vertices
	if len(vertices) == 0 {
		vertices = make(Vertices, 0)
	}

	out := geometryJSONOutput{
		Type:        l.Type(),
		Coordinates: vertices,
	}

	if l.SerializeBBox {
		out.BBox = l.BoundingBox()
	}

	return json.Marshal(&out)
}

// UnmarshalJSON deserializes the GeoJSON data into a LineString.
// Returns an error if the data is invalid or the type does not match LineString.
func (l *LineString) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", l.Type(), err)
	}

	g, ok := gw.geometry.(*LineString)
	if !ok {
		return ErrInvalidTypeField
	}

	l.vertices = g.vertices

	return nil
}

// NewLineString creates a new LineString from the provided vertices.
// Returns an error if the number of vertices is less than 2.
func NewLineString(v Vertices) (*LineString, error) {
	if len(v) < LineStringMinimumSize {
		return nil, ErrLineStringTooShort
	}

	return &LineString{
		vertices: v,
	}, nil
}

// MustLineString creates a new LineString from the provided vertices.
// It panics if the provided vertices are invalid, such as having fewer than 2 vertices.
func MustLineString(v Vertices) *LineString {
	l, err := NewLineString(v)
	if err != nil {
		panic(err)
	}

	return l
}
