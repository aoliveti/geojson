package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeometryCollection_BoundingBox(t *testing.T) {
	tests := []struct {
		name       string
		geometries []Geometry
		expected   BoundingBox
	}{
		{
			"empty geometry collection",
			nil,
			BoundingBox{},
		},
		{
			"single point geometry",
			[]Geometry{MustPoint([]float64{0, 0})},
			BoundingBox{0, 0, 0, 0},
		},
		{
			"multiple geometries",
			[]Geometry{
				MustPoint([]float64{0, 0}),
				MustPoint([]float64{5, 5}),
				MustPoint([]float64{-1, -1}),
			},
			BoundingBox{-1, -1, 5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := NewGeometryCollectionFromSlice(tt.geometries)
			assert.Equal(t, tt.expected, gc.BoundingBox())
		})
	}
}

func TestGeometryCollection_Vertices(t *testing.T) {
	tests := []struct {
		name       string
		geometries []Geometry
		expected   Vertices
	}{
		{
			"empty geometry collection",
			nil,
			nil,
		},
		{
			"single point geometry",
			[]Geometry{MustPoint([]float64{1, 1})},
			Vertices{[]float64{1, 1}},
		},
		{
			"multiple geometries",
			[]Geometry{
				MustPoint([]float64{1, 1}),
				MustPoint([]float64{2, 2}),
			},
			Vertices{
				[]float64{1, 1},
				[]float64{2, 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := NewGeometryCollectionFromSlice(tt.geometries)
			assert.Equal(t, tt.expected, gc.Vertices())
		})
	}
}

func TestGeometryCollection_Type(t *testing.T) {
	gc := NewGeometryCollection()
	assert.Equal(t, TypeGeometryCollection, gc.Type())
}

func TestGeometryCollection_BuildCoordinates(t *testing.T) {
	gc := NewGeometryCollection()
	err := gc.buildCoordinates(nil)
	assert.ErrorIs(t, err, ErrGeometryCollectionBuildCoordinates)
}

func TestGeometryCollection_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		geometries []Geometry
		expected   string
	}{
		{
			"empty geometry collection",
			nil,
			`{"type":"GeometryCollection","geometries":[]}`,
		},
		{
			"single geometry",
			[]Geometry{MustPoint([]float64{1, 1})},
			`{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,1]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := NewGeometryCollectionFromSlice(tt.geometries)
			data, err := gc.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestGeometryCollection_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		expectErr  error
		checkEmpty bool
	}{
		{
			"valid empty geometry collection",
			`{"type":"GeometryCollection","geometries":[]}`,
			nil,
			true,
		},
		{
			"valid single geometry",
			`{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,1]}]}`,
			nil,
			false,
		},
		{
			"different geometry type",
			`{"type":"Point","coordinates":[1, 2]}`,
			ErrInvalidTypeField,
			false,
		},
		{
			"invalid geometry data",
			`{"type":"GeometryCollection","geometries":[{"type":"InvalidType"}]}`,
			ErrInvalidTypeField,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc := NewGeometryCollection()
			err := gc.UnmarshalJSON([]byte(tt.jsonInput))
			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
				if tt.checkEmpty {
					assert.Empty(t, gc.geometries)
				} else {
					assert.NotEmpty(t, gc.geometries)
				}
			}
		})
	}
}

func TestNewGeometryCollection(t *testing.T) {
	gc := NewGeometryCollection()
	assert.Empty(t, gc.geometries)
	assert.IsType(t, &GeometryCollection{}, gc)
}

func TestNewGeometryCollectionFromSlice(t *testing.T) {
	geometries := []Geometry{MustPoint([]float64{1, 1})}
	gc := NewGeometryCollectionFromSlice(geometries)
	assert.Equal(t, len(geometries), len(gc.geometries))
	assert.IsType(t, &GeometryCollection{}, gc)
}

func TestGeometryCollection_Geometries(t *testing.T) {
	type fields struct {
		geometries []Geometry
	}
	tests := []struct {
		name   string
		fields fields
		want   []Geometry
	}{
		{
			"empty geometry collection",
			fields{geometries: nil},
			nil,
		},
		{
			"single geometry",
			fields{geometries: []Geometry{MustPoint([]float64{0, 0})}},
			[]Geometry{MustPoint([]float64{0, 0})},
		},
		{
			"multiple geometries",
			fields{geometries: []Geometry{
				MustPoint([]float64{1, 1}),
				MustPoint([]float64{2, 2}),
			}},
			[]Geometry{
				MustPoint([]float64{1, 1}),
				MustPoint([]float64{2, 2}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeometryCollection{
				geometries: tt.fields.geometries,
			}
			assert.Equalf(t, tt.want, g.Geometries(), "Geometries()")
		})
	}
}
