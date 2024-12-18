package geojson

// TypeEmptyObject represents an empty GeoJSON object type.
const (
	TypeEmptyObject ObjectType = "Object"

	// TypeFeature represents a single GeoJSON feature type.
	TypeFeature ObjectType = "Feature"

	// TypeFeatureCollection represents a GeoJSON feature collection type.
	TypeFeatureCollection ObjectType = "FeatureCollection"
)

// ObjectType defines the type of a GeoJSON object as a string.
type ObjectType string

// FeatureIdentifier is an interface for objects that can return their GeoJSON type.
type FeatureIdentifier interface {
	// Type returns the type of the GeoJSON object.
	Type() ObjectType
}
