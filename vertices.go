package geojson

import "fmt"

// Vertices represents a slice of Coordinates, used to define geometric shapes.
type Vertices []Coordinates

// VerticesBuilder is a builder for constructing Vertices objects.
type VerticesBuilder struct {
	vertices Vertices
	err      error
}

// NewVerticesBuilder initializes and returns a new instance of VerticesBuilder.
func NewVerticesBuilder() *VerticesBuilder {
	return &VerticesBuilder{}
}

// Add appends a new set of coordinates to the builder.
// It validates and creates the Coordinates object before adding it to Vertices.
// Returns the same instance of VerticesBuilder to allow method chaining.
func (vb *VerticesBuilder) Add(v []float64) *VerticesBuilder {
	if vb.err != nil {
		return vb
	}

	coords, err := NewCoordinates(v)
	if err != nil {
		vb.err = fmt.Errorf(
			"failed to create coordinates: %w",
			err)
		return vb
	}

	vb.vertices = append(vb.vertices, *coords)
	return vb
}

// Build finalizes and returns the constructed Vertices object and any error encountered during its construction.
func (vb *VerticesBuilder) Build() (Vertices, error) {
	if vb.err != nil {
		return nil, vb.err
	}

	return vb.vertices, nil
}
