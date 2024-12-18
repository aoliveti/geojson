package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeometryObject_Type(t *testing.T) {
	tests := []struct {
		name     string
		geometry Geometry
		expected GeometryType
	}{
		{"NilGeometry", nil, TypeEmptyGeometry},
		{"PointType", &Point{}, TypePoint},
		{"LineStringType", &LineString{}, TypeLineString},
		{"PolygonType", &Polygon{}, TypePolygon},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := GeometryObject{geometry: test.geometry}
			assert.Equal(t, test.expected, g.Type())
		})
	}
}

func TestGeometryObject_MarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		geometryObject GeometryObject
		expectError    error
	}{
		{"NilGeometry", GeometryObject{}, ErrGeometryNotDefined},
		{"ValidPoint", GeometryObject{geometry: &Point{coords: Coordinates{1.0, 2.0}}}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &test.geometryObject
			data, err := g.MarshalJSON()

			if test.expectError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectError.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}

func TestGeometryObject_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError error
	}{
		{"InvalidJSON", `{"type": "Invalid"}`, ErrInvalidTypeField},
		{"ValidPoint", `{"type": "Point", "coordinates": [1.0, 2.0]}`, nil},
		{"ValidPolygon", `{"type": "Polygon", "coordinates": [[[1.0, 1.0], [2.0, 2.0], [3.0, 3.0], [1.0, 1.0]]]}`, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var g GeometryObject
			err := g.UnmarshalJSON([]byte(test.input))

			if test.expectError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGeometryObject_IsTypeChecks(t *testing.T) {
	tests := []struct {
		name     string
		geometry Geometry
		isValid  func(g *GeometryObject) bool
		expected bool
	}{
		{"IsGeometryCollectionTrue", &GeometryCollection{}, (*GeometryObject).IsGeometryCollection, true},
		{"IsGeometryCollectionFalse", &Point{}, (*GeometryObject).IsGeometryCollection, false},
		{"IsLineStringTrue", &LineString{}, (*GeometryObject).IsLineString, true},
		{"IsLineStringFalse", &Point{}, (*GeometryObject).IsLineString, false},
		{"IsMultiLineStringTrue", &MultiLineString{}, (*GeometryObject).IsMultiLineString, true},
		{"IsMultiLineStringFalse", &Point{}, (*GeometryObject).IsMultiLineString, false},
		{"IsMultiPointTrue", &MultiPoint{}, (*GeometryObject).IsMultiPoint, true},
		{"IsMultiPointFalse", &Polygon{}, (*GeometryObject).IsMultiPoint, false},
		{"IsMultiPolygonTrue", &MultiPolygon{}, (*GeometryObject).IsMultiPolygon, true},
		{"IsMultiPolygonFalse", &Point{}, (*GeometryObject).IsMultiPolygon, false},
		{"IsPointTrue", &Point{}, (*GeometryObject).IsPoint, true},
		{"IsPointFalse", &Polygon{}, (*GeometryObject).IsPoint, false},
		{"IsPolygonTrue", &Polygon{}, (*GeometryObject).IsPolygon, true},
		{"IsPolygonFalse", &Point{}, (*GeometryObject).IsPolygon, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := GeometryObject{geometry: test.geometry}
			assert.Equal(t, test.expected, test.isValid(&g))
		})
	}
}

func TestGeometryObject_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		geometry Geometry
		isEmpty  bool
	}{
		{"EmptyGeometry", nil, true},
		{"NonEmptyGeometry", &Point{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := GeometryObject{geometry: test.geometry}
			assert.Equal(t, test.isEmpty, g.IsEmpty())
		})
	}
}

func TestGeometryObject_ToTypeConversions(t *testing.T) {
	tests := []struct {
		name        string
		input       GeometryObject
		converter   func(g *GeometryObject) (interface{}, error)
		expectError bool
	}{
		{"ToPointSuccess", GeometryObject{geometry: &Point{}}, func(g *GeometryObject) (interface{}, error) { return g.ToPoint() }, false},
		{"ToPointError", GeometryObject{geometry: &Polygon{}}, func(g *GeometryObject) (interface{}, error) { return g.ToPoint() }, true},
		{"ToPointWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToPoint() }, true},

		{"ToLineStringSuccess", GeometryObject{geometry: &LineString{}}, func(g *GeometryObject) (interface{}, error) { return g.ToLineString() }, false},
		{"ToLineStringError", GeometryObject{geometry: &MultiPoint{}}, func(g *GeometryObject) (interface{}, error) { return g.ToLineString() }, true},
		{"ToLineStringWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToLineString() }, true},

		{"ToMultiPointSuccess", GeometryObject{geometry: &MultiPoint{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPoint() }, false},
		{"ToMultiPointError", GeometryObject{geometry: &LineString{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPoint() }, true},
		{"ToMultiPointWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPoint() }, true},

		{"ToMultiLineStringSuccess", GeometryObject{geometry: &MultiLineString{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiLineString() }, false},
		{"ToMultiLineStringError", GeometryObject{geometry: &Point{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiLineString() }, true},
		{"ToMultiLineStringWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiLineString() }, true},

		{"ToPolygonSuccess", GeometryObject{geometry: &Polygon{}}, func(g *GeometryObject) (interface{}, error) { return g.ToPolygon() }, false},
		{"ToPolygonError", GeometryObject{geometry: &MultiPoint{}}, func(g *GeometryObject) (interface{}, error) { return g.ToPolygon() }, true},
		{"ToPolygonWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToPolygon() }, true},

		{"ToMultiPolygonSuccess", GeometryObject{geometry: &MultiPolygon{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPolygon() }, false},
		{"ToMultiPolygonError", GeometryObject{geometry: &GeometryCollection{}}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPolygon() }, true},
		{"ToMultiPolygonWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToMultiPolygon() }, true},

		{"ToGeometryCollectionSuccess", GeometryObject{geometry: &GeometryCollection{}}, func(g *GeometryObject) (interface{}, error) { return g.ToGeometryCollection() }, false},
		{"ToGeometryCollectionError", GeometryObject{geometry: &Point{}}, func(g *GeometryObject) (interface{}, error) { return g.ToGeometryCollection() }, true},
		{"ToGeometryCollectionWithEmpty", GeometryObject{}, func(g *GeometryObject) (interface{}, error) { return g.ToGeometryCollection() }, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.converter(&test.input)

			if test.expectError {
				require.Error(t, err)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestFromGeometry(t *testing.T) {
	type args struct {
		g Geometry
	}
	tests := []struct {
		name string
		args args
		want GeometryObject
	}{
		{"NilGeometry", args{g: nil}, GeometryObject{geometry: nil}},
		{"PointGeometry", args{g: MustPoint([]float64{1, 2})}, GeometryObject{geometry: MustPoint([]float64{1, 2})}},
		{"LineStringGeometry", args{g: MustLineString(Vertices{{1, 2}, {1, 3}})}, GeometryObject{geometry: MustLineString(Vertices{{1, 2}, {1, 3}})}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FromGeometry(tt.args.g), "FromGeometry(%v)", tt.args.g)
		})
	}
}
