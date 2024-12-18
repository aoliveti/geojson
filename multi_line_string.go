package geojson

import (
	"encoding/json"
	"fmt"
)

var (
	// ErrMultiLineStringTooShort is returned when a MultiLineString has fewer than the minimum required segments.
	ErrMultiLineStringTooShort = fmt.Errorf("line string must have at least one segment")
)

// MultiLineString represents a GeoJSON MultiLineString geometry.
type MultiLineString struct {
	segments      Segments // Segments that define the MultiLineString.
	SerializeBBox bool     // Indicates whether the bounding box should be included during JSON serialization.
}

// BoundingBox calculates and returns the bounding box of the MultiLineString.
func (m *MultiLineString) BoundingBox() BoundingBox {
	return bbox(m.Vertices())
}

// Vertices gathers and returns all vertices from the segments of the MultiLineString.
func (m *MultiLineString) Vertices() Vertices {
	var v Vertices
	for _, s := range m.segments {
		v = append(v, s...)
	}
	return v
}

// Type returns the GeoJSON geometry type of the MultiLineString: "MultiLineString".
func (m *MultiLineString) Type() GeometryType {
	return TypeMultiLineString
}

// Segments returns the collection of segments that define the MultiLineString.
func (m *MultiLineString) Segments() Segments {
	return m.segments
}

// buildCoordinates processes raw GeoJSON coordinates and constructs the segments of the MultiLineString.
func (m *MultiLineString) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	if len(rawSlice) == 0 {
		return ErrMultiLineStringTooShort
	}

	segments := make(Segments, len(rawSlice))
	for i, s := range rawSlice {
		l := LineString{}
		if err := l.buildCoordinates(s); err != nil {
			return err
		}
		segments[i] = l.vertices
	}

	m.segments = segments
	return nil
}

// MarshalJSON converts the MultiLineString into its GeoJSON representation.
func (m *MultiLineString) MarshalJSON() ([]byte, error) {
	out := geometryJSONOutput{
		Type:        m.Type(),
		Coordinates: m.segments,
	}

	if m.SerializeBBox {
		out.BBox = m.BoundingBox()
	}

	return json.Marshal(&out)
}

// UnmarshalJSON parses a GeoJSON representation into a MultiLineString.
func (m *MultiLineString) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", m.Type(), err)
	}

	g, ok := gw.geometry.(*MultiLineString)
	if !ok {
		return ErrInvalidTypeField
	}

	m.segments = g.segments

	return nil
}

// NewMultiLineString creates a new MultiLineString with the provided segments.
func NewMultiLineString(segments Segments) (*MultiLineString, error) {
	if len(segments) == 0 {
		return nil, ErrMultiLineStringTooShort
	}

	for _, s := range segments {
		if _, err := NewLineString(s); err != nil {
			return nil, err
		}
	}

	return &MultiLineString{
		segments: segments,
	}, nil
}

// MustMultiLineString constructs a new MultiLineString and panics if there is an error.
func MustMultiLineString(segments Segments) *MultiLineString {
	mls, err := NewMultiLineString(segments)
	if err != nil {
		panic(err)
	}

	return mls
}
