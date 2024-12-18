package geojson

// featuresJSONInput represents the input structure for a GeoJSON object,
// used to deserialize both single features and feature collections.
type featuresJSONInput struct {
	Type       ObjectType      `json:"type"`       // Specifies the type of GeoJSON object (e.g., "Feature" or "FeatureCollection").
	Geometry   *GeometryObject `json:"geometry"`   // Contains the geometry of the GeoJSON feature (if applicable).
	Properties Properties      `json:"properties"` // Describes additional properties of the GeoJSON feature.
	ID         *ID             `json:"id"`         // Optional identifier for the GeoJSON feature.
	Features   []Feature       `json:"features"`   // An array of features (used if part of a feature collection).
}

// featureCollectionJSONOutput represents the output structure of a GeoJSON FeatureCollection.
// It contains a collection of features and, optionally, a bounding box.
type featureCollectionJSONOutput struct {
	Type     ObjectType  `json:"type"`           // Specifies the type of GeoJSON object (e.g., "FeatureCollection").
	Features []Feature   `json:"features"`       // An array of features within the collection.
	BBox     BoundingBox `json:"bbox,omitempty"` // Optional bounding box that encloses all features in the collection.
}

// featureJSONOutput represents the output structure for a single GeoJSON feature.
// It includes geometry, properties, an optional ID, and an optional bounding box.
type featureJSONOutput struct {
	Type       ObjectType  `json:"type"`                 // Specifies the type of GeoJSON object (e.g., "Feature").
	Geometry   Geometry    `json:"geometry"`             // Contains the geometry of the GeoJSON feature.
	Properties Properties  `json:"properties,omitempty"` // Describes additional properties of the GeoJSON feature.
	ID         *ID         `json:"id,omitempty"`         // Optional identifier for the GeoJSON feature.
	BBox       BoundingBox `json:"bbox,omitempty"`       // Optional bounding box that encloses the feature.
}

// geometryJSONInput represents the input structure for a GeoJSON geometry.
// It captures the type, coordinates, optional bounding box, and sub-geometries when
// handling collections.
type geometryJSONInput struct {
	Type        GeometryType     `json:"type"`        // Specifies the type of geometry (e.g., "Point", "Polygon").
	Coordinates interface{}      `json:"coordinates"` // Contains the coordinates for the geometry.
	Geometries  []GeometryObject `json:"geometries"`  // Contains sub-geometries if part of a geometry collection.
	BBox        BoundingBox      `json:"bbox"`        // Optional bounding box that encloses the geometry.
}

// geometryJSONOutput represents the output structure for a GeoJSON geometry.
// It includes the type, coordinates, and an optional bounding box.
type geometryJSONOutput struct {
	Type        GeometryType `json:"type"`           // Specifies the type of geometry (e.g., "Point", "Polygon").
	Coordinates interface{}  `json:"coordinates"`    // Contains the coordinates for the geometry.
	BBox        BoundingBox  `json:"bbox,omitempty"` // Optional bounding box that encloses the geometry.
}

// geometryCollectionJSONOutput represents the output structure for a GeoJSON geometry collection.
// It specifies the type and contains an array of geometries.
type geometryCollectionJSONOutput struct {
	Type       GeometryType `json:"type"`       // Specifies the type of geometry collection (e.g., "GeometryCollection").
	Geometries []Geometry   `json:"geometries"` // An array of geometries contained in the collection.
}
