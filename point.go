package geojson

import (
	"encoding/json"
	"fmt"
)

// Point represents a GeoJSON Point object with coordinates and optional serialization for a bounding box.
type Point struct {
	coords        Coordinates
	SerializeBBox bool
}

// BoundingBox computes the bounding box of the Point.
// If SerializeBBox is true, the bounding box will be included during serialization.
func (p *Point) BoundingBox() BoundingBox {
	return bbox(p.Vertices())
}

// Vertices returns the coordinates of the Point as a slice of Vertices.
func (p *Point) Vertices() Vertices {
	var v Vertices
	v = append(v, p.coords)
	return v
}

// Longitude returns the longitude of the Point.
func (p *Point) Longitude() float64 {
	return p.coords.Longitude()
}

// Latitude returns the latitude of the Point.
func (p *Point) Latitude() float64 {
	return p.coords.Latitude()
}

// Coordinates returns the coordinates of the Point.
func (p *Point) Coordinates() Coordinates {
	return p.coords
}

// HasAltitude checks if the Point includes an altitude value.
func (p *Point) HasAltitude() bool {
	return p.coords.HasAltitude()
}

// Altitude returns the altitude of the Point.
// This should only be called if HasAltitude() returns true.
func (p *Point) Altitude() float64 {
	return p.coords.Altitude()
}

// Type returns the GeoJSON type of the Point as a GeometryType.
func (p *Point) Type() GeometryType {
	return TypePoint
}

// buildCoordinates creates the coordinates for the Point from a raw slice of interface{}.
func (p *Point) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	coords, err := buildCoordinates(rawSlice)
	if err != nil {
		return err
	}

	p.coords = *coords

	return nil
}

// MarshalJSON implements the json.Marshaler interface to serialize the Point into GeoJSON format.
func (p *Point) MarshalJSON() ([]byte, error) {
	out := geometryJSONOutput{
		Type:        p.Type(),
		Coordinates: p.coords,
	}

	if p.SerializeBBox {
		out.BBox = p.BoundingBox()
	}

	return json.Marshal(&out)
}

// UnmarshalJSON implements the json.Unmarshaler interface to parse GeoJSON data into a Point.
func (p *Point) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", p.Type(), err)
	}

	g, ok := gw.geometry.(*Point)
	if !ok {
		return ErrInvalidTypeField
	}

	p.coords = g.coords

	return nil
}

// NewPoint creates a new Point from a slice of float64 coordinates.
// Returns an error if the coordinates are invalid.
func NewPoint(v []float64) (*Point, error) {
	coords, err := NewCoordinates(v)
	if err != nil {
		return nil, err
	}

	return &Point{coords: *coords}, nil
}

// MustPoint creates a new Point and panics if the coordinates are invalid.
// This function should only be used when the coordinates are guaranteed to be valid.
func MustPoint(v []float64) *Point {
	point, err := NewPoint(v)
	if err != nil {
		panic(err)
	}

	return point
}
