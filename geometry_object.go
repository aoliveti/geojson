package geojson

import (
	"encoding/json"
	"errors"
)

var (
	// ErrGeometryNotDefined indicates that the geometry is not defined.
	ErrGeometryNotDefined = errors.New("geometry is not defined")

	// ErrGeometryTypeMismatch indicates that the geometry type does not match the expected type.
	ErrGeometryTypeMismatch = errors.New("geometry type mismatch")

	// ErrInvalidTypeField indicates that the type field is invalid or missing in the JSON data.
	ErrInvalidTypeField = errors.New("invalid or missing type field")

	// ErrInvalidCoordinates indicates that the coordinates field is invalid or missing in the JSON data.
	ErrInvalidCoordinates = errors.New("invalid or missing coordinates")
)

// GeometryObject represents a GeoJSON Geometry Object.
type GeometryObject struct {
	geometry Geometry
}

// Type returns the geometry type of the GeometryObject.
func (g *GeometryObject) Type() GeometryType {
	if g.geometry == nil {
		return TypeEmptyGeometry
	}

	return g.geometry.Type()
}

// MarshalJSON marshals the GeometryObject into its JSON representation.
func (g *GeometryObject) MarshalJSON() ([]byte, error) {
	if g.geometry == nil || g.geometry.Type() == TypeEmptyGeometry {
		return nil, ErrGeometryNotDefined
	}

	return json.Marshal(g.geometry)
}

// UnmarshalJSON unmarshals JSON data into the GeometryObject.
func (g *GeometryObject) UnmarshalJSON(bytes []byte) error {
	geometry := geometryJSONInput{}
	if err := json.Unmarshal(bytes, &geometry); err != nil {
		return err
	}

	var v Geometry
	switch geometry.Type {
	case TypePoint:
		v = &Point{}
	case TypeLineString:
		v = &LineString{}
	case TypeMultiPoint:
		v = &MultiPoint{}
	case TypeMultiLineString:
		v = &MultiLineString{}
	case TypePolygon:
		v = &Polygon{}
	case TypeMultiPolygon:
		v = &MultiPolygon{}
	case TypeGeometryCollection:
		gc := &GeometryCollection{}
		for _, gm := range geometry.Geometries {
			gc.geometries = append(gc.geometries, gm.geometry)
		}
		g.geometry = gc
		return nil
	default:
		return ErrInvalidTypeField
	}

	if err := v.buildCoordinates(geometry.Coordinates); err != nil {
		return err
	}

	g.geometry = v

	return nil
}

// IsPoint checks if the GeometryObject is of type Point.
func (g *GeometryObject) IsPoint() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypePoint
}

// IsLineString checks if the GeometryObject is of type LineString.
func (g *GeometryObject) IsLineString() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypeLineString
}

// IsMultiPoint checks if the GeometryObject is of type MultiPoint.
func (g *GeometryObject) IsMultiPoint() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypeMultiPoint
}

// IsMultiLineString checks if the GeometryObject is of type MultiLineString.
func (g *GeometryObject) IsMultiLineString() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypeMultiLineString
}

// IsPolygon checks if the GeometryObject is of type Polygon.
func (g *GeometryObject) IsPolygon() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypePolygon
}

// IsMultiPolygon checks if the GeometryObject is of type MultiPolygon.
func (g *GeometryObject) IsMultiPolygon() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypeMultiPolygon
}

// IsGeometryCollection checks if the GeometryObject is of type GeometryCollection.
func (g *GeometryObject) IsGeometryCollection() bool {
	return !g.IsEmpty() && g.geometry.Type() == TypeGeometryCollection
}

// IsEmpty checks if the GeometryObject is empty or not defined.
func (g *GeometryObject) IsEmpty() bool {
	return g.geometry == nil
}

// ToPoint converts the GeometryObject into a Point, returning an error if the type does not match.
func (g *GeometryObject) ToPoint() (*Point, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*Point)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToLineString converts the GeometryObject into a LineString, returning an error if the type does not match.
func (g *GeometryObject) ToLineString() (*LineString, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*LineString)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToMultiPoint converts the GeometryObject into a MultiPoint, returning an error if the type does not match.
func (g *GeometryObject) ToMultiPoint() (*MultiPoint, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*MultiPoint)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToMultiLineString converts the GeometryObject into a MultiLineString, returning an error if the type does not match.
func (g *GeometryObject) ToMultiLineString() (*MultiLineString, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*MultiLineString)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToPolygon converts the GeometryObject into a Polygon, returning an error if the type does not match.
func (g *GeometryObject) ToPolygon() (*Polygon, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*Polygon)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToMultiPolygon converts the GeometryObject into a MultiPolygon, returning an error if the type does not match.
func (g *GeometryObject) ToMultiPolygon() (*MultiPolygon, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*MultiPolygon)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// ToGeometryCollection converts the GeometryObject into a GeometryCollection, returning an error if the type does not match.
func (g *GeometryObject) ToGeometryCollection() (*GeometryCollection, error) {
	if g.IsEmpty() {
		return nil, ErrGeometryNotDefined
	}

	v, ok := g.geometry.(*GeometryCollection)
	if !ok {
		return nil, ErrGeometryTypeMismatch
	}

	return v, nil
}

// FromGeometry creates and returns a new GeometryObject given a Geometry.
// The input Geometry is assigned to the geometry field of the GeometryObject.
func FromGeometry(g Geometry) GeometryObject {
	return GeometryObject{
		geometry: g,
	}
}
