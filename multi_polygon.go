package geojson

import (
	"encoding/json"
	"fmt"
)

// MultiPolygon represents a GeoJSON MultiPolygon geometry.
type MultiPolygon struct {
	rings         []LinearRings
	SerializeBBox bool
}

// Type returns the geometry type for MultiPolygon.
func (m *MultiPolygon) Type() GeometryType {
	return TypeMultiPolygon
}

// Vertices collects and returns all vertices contained in the MultiPolygon.
func (m *MultiPolygon) Vertices() Vertices {
	var v Vertices
	for _, segment := range m.rings {
		for _, vertex := range segment {
			v = append(v, vertex...)
		}
	}
	return v
}

// BoundingBox computes and returns the bounding box enclosing the MultiPolygon.
func (m *MultiPolygon) BoundingBox() BoundingBox {
	return bbox(m.Vertices())
}

// LinearRingsSlice returns the MultiPolygon's internal rings as a slice of LinearRings.
func (m *MultiPolygon) LinearRingsSlice() []LinearRings {
	return m.rings
}

// MarshalJSON serializes the MultiPolygon to its GeoJSON representation.
func (m *MultiPolygon) MarshalJSON() ([]byte, error) {
	rings := m.rings
	if len(rings) == 0 {
		rings = make([]LinearRings, 0)
	}

	out := geometryJSONOutput{
		Type:        m.Type(),
		Coordinates: rings,
	}

	if m.SerializeBBox {
		out.BBox = m.BoundingBox()
	}

	return json.Marshal(&out)
}

// UnmarshalJSON deserializes the GeoJSON representation into a MultiPolygon.
func (m *MultiPolygon) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", m.Type(), err)
	}

	g, ok := gw.geometry.(*MultiPolygon)
	if !ok {
		return ErrInvalidTypeField
	}

	m.rings = g.rings

	return nil
}

// NewMultiPolygon creates and returns a new empty MultiPolygon instance.
func NewMultiPolygon() *MultiPolygon {
	return &MultiPolygon{}
}

// NewMultiPolygonFromRingSlice validates the provided slice of LinearRings and creates
// a new MultiPolygon. It ensures each LinearRing within the slice has a valid size
// and is closed. If validation fails, an error is returned.
func NewMultiPolygonFromRingSlice(slice []LinearRings) (*MultiPolygon, error) {
	for _, rings := range slice {
		for _, ring := range rings {
			if !ring.HasValidSize() {
				return nil, ErrLinearRingSize
			}
			if !ring.IsClosed() {
				return nil, ErrLinearRingClosed
			}
		}

		ensureOrientation(rings)
	}

	return &MultiPolygon{
		rings: slice,
	}, nil
}

// MustMultiPolygonFromRingSlice creates a new MultiPolygon from the provided slice of LinearRings.
// It panics if the creation fails due to invalid LinearRings.
func MustMultiPolygonFromRingSlice(slice []LinearRings) *MultiPolygon {
	mp, err := NewMultiPolygonFromRingSlice(slice)
	if err != nil {
		panic(err)
	}

	return mp
}

// buildCoordinates initializes the MultiPolygon rings based on the provided raw coordinate data.
func (m *MultiPolygon) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	lrSlice := make([]LinearRings, len(rawSlice))
	for i, s := range rawSlice {
		p := Polygon{}

		if err := p.buildCoordinates(s); err != nil {
			return err
		}

		lrSlice[i] = p.rings
	}

	m.rings = lrSlice

	return nil
}
