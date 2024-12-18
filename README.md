# GeoJSON
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aoliveti/geojson)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/aoliveti/geojson/go.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/aoliveti/geojson)](https://pkg.go.dev/github.com/aoliveti/geojson)
[![codecov](https://codecov.io/gh/aoliveti/geojson/graph/badge.svg?token=TmNqSDCdJ9)](https://codecov.io/gh/aoliveti/geojson)
[![Go Report Card](https://goreportcard.com/badge/github.com/aoliveti/geojson)](https://goreportcard.com/report/github.com/aoliveti/geojson)
![GitHub License](https://img.shields.io/github/license/aoliveti/geojson)

**geojson** is a Go package implementing the RFC 7946 standard.
It provides full support for GeoJSON objects, including **Features**, **FeatureCollections**, and **Geometries**, offering comprehensive functionality for serializing and deserializing GeoJSON data.

- Strict adherence to the [RFC 7946 GeoJSON standard](https://datatracker.ietf.org/doc/html/rfc7946).
- Support for GeoJSON object types:
  - **Feature**
    - with ID and Properties
  - **FeatureCollection**
  - **Geometries**
    - Supported geometries: `Point`, `LineString`, `Polygon`, `MultiPoint`, `MultiLineString`, `MultiPolygon`, and
      `GeometryCollection`.
- Robust validation of object types and coordinates.
- Bounding box generation.
- Polygon and MultiPolygon **linear rings are automatically oriented to follow the right-hand rule**: outer rings are counterclockwise, holes are clockwise.
- All wrappers, features and geometries support JSON marshalling and unmarshalling

## Install

```bash
go get github.com/aoliveti/geojson
```

---

## Usage
### GeoJSON Object (`geojson.Object`)

The `geojson.Object` is a generic representation of a GeoJSON object, providing the tools to parse, inspect, and work with GeoJSON features and metadata. This structure is useful for handling GeoJSON data programmatically by offering methods to determine the type of GeoJSON object and interact with its contents.

- Methods:
  - `Type()` - Returns the GeoJSON type (`Feature` or `FeatureCollection`).
  - `IsFeature()` - Checks if the object is a **Feature**.
  - `IsFeatureCollection()` - Checks if the object is a **FeatureCollection**.
  - `Feature()` - Retrieves the **Feature** object (error if not a feature).
  - `FeatureCollection()` - Retrieves the **FeatureCollection** (error if not a collection).
  - `MarshalJSON()` / `UnmarshalJSON()` - Converts the object to/from JSON.

The possible types include:

- **`TypeEmptyObject`**: Represents an invalid or uninitialized GeoJSON object.
- **`TypeFeature`**: Indicates the object is a single GeoJSON Feature.
- **`TypeFeatureCollection`**: Indicates the object is a GeoJSON FeatureCollection.

### Unmarshalling GeoJSON
Here is a quick overview of how to work with `geojson.Object` and its methods:
```go
  data := `{
		"type": "Feature",
		"geometry": {
			"type": "Point",
			"coordinates": [102.0, 0.5]
		},
		"properties": {
			"name": "Sample Point"
		}
	}`

  var obj geojson.Object
  if err := obj.UnmarshalJSON([]byte(data)); err != nil {
    ...
  }

  if !obj.IsFeature() {
    ...
  }

  feature, _ := obj.Feature()
  fmt.Println(feature)
}
```
### Feature (`geojson.Feature`)

The `geojson.Feature` represents a single GeoJSON Feature. It consists in a spatial geometry, properties, and optionally an ID. The `Feature` object also supports operations such as calculating bounding boxes and extracting vertices.

- **Attributes**:
  - `Geometry`: Specifies the spatial information of the feature (point, line, polygon, etc.).
  - `Properties`: A map containing additional metadata about the feature.
  - `ID`: An optional identifier for the feature.
  - `SerializeBBox`: A flag determining whether to include a bounding box in serialized output.

- **Feature Builder**:
  Use the `FeatureBuilder` to programmatically build a `Feature` with the following methods:
  - `SetGeometry(geometry Geometry)`
  - `SetProperties(properties Properties)`
  - `SetID(id ID)`
  - `Build()`

#### Example

```go
geometry := geojson.MustPoint([]float64{102.0, 0.5})
properties := geojson.Properties{"name": "Sample Point"}

builder := geojson.NewFeatureBuilder()
feature := builder.SetGeometry(geometry).
    SetProperties(properties).
    Build()

data, err := json.Marshal(&feature)
if err != nil {
    ...
}

fmt.Println(string(data))
```

---

### FeatureCollection (`geojson.FeatureCollection`)
The `geojson.FeatureCollection` represents a GeoJSON object that contains a collection of `Feature` objects.

- **Attributes**:
  - `Features`: A list of `Feature` objects in the collection.
  - `SerializeBBox`: A flag determining whether to include a bounding box in serialized output.

- **Constructor Methods**:
  - **`NewFeatureCollection()`**: Creates an empty `FeatureCollection`.
  - **`NewFeatureCollectionFromFeatures(features []Feature)`**: Initializes a `FeatureCollection` with the given slice of features.

#### Example
```go
features := []geojson.Feature{feature1, feature2, feature3}
collection := geojson.NewFeatureCollectionFromFeatures(features)
```

---

### Working with Properties
The `Properties` type manages metadata as key-value pairs for GeoJSON features.

- **Add/Update**: Use `Set("key", value)` to add or update properties.
- **Retrieve**: Use `Get("key")` or typed methods (`GetString`, `GetInt`, etc.) for type-safe access.

Example:
```go
properties := geojson.Properties{}
properties.Set("name", "Foo")
name, _ := properties.GetString("name")
```

Typed methods return errors for missing keys or type mismatches, such as `ErrPropertyNotFound`.

---

### GeometryObject (`geojson.GeometryObject`)
The `geojson.GeometryObject` serves as a generic wrapper for GeoJSON geometries. Its main purpose is to handle GeoJSON geometry data during the unmarshalling process, allowing for flexible deserialization and easy identification of the actual geometry type. Additionally, it provides methods to safely access specific geometry types and convert the wrapped geometry into the desired structure.
- **Generic Geometry Wrapping**:  
  The `GeometryObject` wraps any GeoJSON-compatible geometry, including:
  - Point
  - LineString
  - MultiPoint
  - MultiLineString
  - Polygon
  - MultiPolygon
  - GeometryCollection

- **Geometry Identification**:  
  The type of the contained geometry can be identified using the **`Type()`** method, which returns a `GeometryType`.

- **Safe Type Conversion**:  
  The `GeometryObject` provides methods to check whether the contained geometry matches a specific type and to convert the contained geometry into a strongly typed object (e.g., `Point`, `LineString`, etc.).

- **Marshalling/Unmarshalling Support**:  
  The `GeometryObject` handles serialization to and deserialization from GeoJSON formatted data, ensuring compatibility with the GeoJSON format.

#### Key Methods
1. **Geometry Type and State Methods**:
  - **`func (g *GeometryObject) Type() GeometryType`**  
    Returns the type of the contained geometry (e.g., `Point`, `Polygon`, etc.).

  - **`func (g *GeometryObject) IsEmpty() bool`**  
    Checks if the `GeometryObject` is empty or not defined.

  - Specific Type Check Methods:
    - `IsPoint()`, `IsLineString()`, `IsMultiPoint()`, `IsMultiLineString()`, `IsPolygon()`, `IsMultiPolygon()`, `IsGeometryCollection()`

2. **Type Conversion Methods**:  
   Converts the contained geometry into a strongly typed object and returns an error if the type does not match:
  - `ToPoint() (*Point, error)`
  - `ToLineString() (*LineString, error)`
  - `ToMultiPoint() (*MultiPoint, error)`
  - `ToMultiLineString() (*MultiLineString, error)`
  - `ToPolygon() (*Polygon, error)`
  - `ToMultiPolygon() (*MultiPolygon, error)`
  - `ToGeometryCollection() (*GeometryCollection, error)`

3. **Marshalling and Unmarshalling**:
  - **`func (g *GeometryObject) MarshalJSON() ([]byte, error)`**  
    Serializes the `GeometryObject` into GeoJSON format. Returns an error if the geometry is not defined.

  - **`func (g *GeometryObject) UnmarshalJSON(data []byte) error`**  
    Deserializes GeoJSON data into the `GeometryObject`. Automatically detects and handles the actual geometry type (e.g., Point, Polygon, etc.).

4. **Wrap existing Geometry**:
  - **`static func FromGeometry(g Geometry) GeometryObject`**  
    Creates a new `GeometryObject` from an existing `Geometry`.

#### Example

```go
// Unmarshal GeoJSON data (e.g., for a Point)
geoJSONData := []byte(`{
        "type": "Point",
        "coordinates": [100.0, 0.0]
    }`)

err := gObj.UnmarshalJSON(geoJSONData)
if err != nil {
    ...
}

// Identify Geometry Type
if gObj.IsPoint() {
    point, err := gObj.ToPoint()
    if err != nil {
        ...
    }

    // Use point
}
```
```go
// Use a switch to check and handle the geometry type
switch gObj.Type() {
    case geojson.TypePoint:
      point, err := gObj.ToPoint()
      if err != nil {
          ...
      }

    case geojson.TypeLineString:
		lineString, err := gObj.ToLineString()

        if err != nil {
            ...
        }
		
    ...
}
```

---

### Coordinates (`geojson.Coordinates`)
The `geojson.Coordinates` type represents a GeoJSON coordinate array (WGS84), containing **longitude**, **latitude**, and optionally **altitude**.

  - Longitude: `-180 <= longitude <= 180`
  - Latitude: `-90 <= latitude <= 90`
  - Coordinates must have either 2 elements (`[longitude, latitude]`) or 3 elements (`[longitude, latitude, altitude]`).
  - **`Longitude()`**: Returns the longitude value of the coordinates.
  - **`Latitude()`**: Returns the latitude value of the coordinates.
  - **`HasAltitude()`**: Checks if the coordinates include an altitude.
  - **`Altitude()`**: Returns the altitude value (if present).
  - **`NewCoordinates([]float64) (*Coordinates, error)`**: Creates a new `Coordinates` object from a float64 array. Returns an error for invalid input.
  - **`MustCoordinates([]float64) *Coordinates`**: Creates a `Coordinates` object and panics on error.
  
#### Example
```go
coordinates, err := geojson.NewCoordinates([]float64{12.4924, 41.8902, 45})
if err != nil {
    ...
}

fmt.Println("longitude:", coordinates.Longitude())
fmt.Println("latitude:", coordinates.Latitude())

if coordinates.HasAltitude() {
    fmt.Println("altitude:", coordinates.Altitude())
}
```

---

### Vertices (`geojson.Vertices`)
The `geojson.Vertices` type is a collection of `Coordinates` used to define geometric shapes (e.g., a line or polygon).
Allows grouping multiple `Coordinates` into a structured, ordered collection.

### Segments (`geojson.Segments`)

The `geojson.Segments` type represents a collection of **line segments**, constructed from groups of `Vertices`. These segments are commonly used to define continuous or discrete line geometries, such as **MultilineStrings**.

### LinearRing (`geojson.LinearRing`)

The `geojson.LinearRing` type represents a **closed linear ring**, a sequence of connected vertices forming a continuous loop. Linear rings are commonly used in **GeoJSON Polygons** and **MultiPolygons**, defining their outer and inner boundaries (holes). This type includes methods for validation, orientation, and basic geometric calculations.

1. **Validation**:
  - Ensures the linear ring has at least the required minimum points (`4` coordinates, including the closing point).
  - Ensures the first and last coordinates are identical (i.e., the ring is closed).

2. **Orientation**:
  - Determines whether the linear ring is ordered **clockwise** or **counterclockwise** based on the **signed area** of its vertices.
  - Provides functionality to adjust the orientation to meet specific requirements (e.g., outer rings counterclockwise and inner rings clockwise).

3. **Geometry Calculations**:
  - Computes the **absolute area** of the linear ring.

---

### Marshalling and Unmarshalling Geometries
All geometries in this package support **JSON marshalling** and **unmarshalling**
Here's an example of how to marshal and unmarshal a single geometry, such as a `Point`:

#### Example: Marshalling a `Point` to JSON

```go
point := geojson.MustPoint([]float64{102.0, 0.5})
pointJSON, err := json.Marshal(point)
if err != nil {
	...
}

fmt.Println(string(pointJSON))
```

#### Example: Unmarshalling a `Point` from JSON

```go
pointJSON := `{"type": "Point", "coordinates": [12.4924, 41.8902]}`

var point geojson.Point
err := json.Unmarshal([]byte(pointJSON), &point)
if err != nil {
    ...
}

fmt.Println(point.Coordinates())
```

---

## Contributing

Contributions are welcome! If you'd like to contribute, please submit a pull request or open an issue to discuss changes.

## License

This library is licensed under the MIT License. See [LICENSE](LICENSE) for details.
