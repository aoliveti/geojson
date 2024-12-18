package geojson

// GeometryType defines the type of geometry in GeoJSON.
type GeometryType string

// Predefined GeometryType constants representing various geometry types in GeoJSON.
const (
	TypeEmptyGeometry      GeometryType = "Object"
	TypePoint              GeometryType = "Point"
	TypeMultiPoint         GeometryType = "MultiPoint"
	TypeLineString         GeometryType = "LineString"
	TypeMultiLineString    GeometryType = "MultiLineString"
	TypePolygon            GeometryType = "Polygon"
	TypeMultiPolygon       GeometryType = "MultiPolygon"
	TypeGeometryCollection GeometryType = "GeometryCollection"
)

// GeometryIdentifier is an interface for objects that can report their geometry type.
type GeometryIdentifier interface {
	// Type returns the GeometryType of the object.
	Type() GeometryType
}

// geometryBuilder is an interface for building coordinates for geometries.
type geometryBuilder interface {
	// buildCoordinates initializes a geometry's coordinates from the given input.
	buildCoordinates(interface{}) error
}

// Geometry is a composite interface that combines GeometryIdentifier, BoundingBoxer,
// and geometryBuilder, representing a GeoJSON geometry object.
type Geometry interface {
	GeometryIdentifier
	BoundingBoxer
	geometryBuilder
}
