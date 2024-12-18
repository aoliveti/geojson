package geojson

import "fmt"

var (
	// ErrVerticesEmpty is returned when an attempt is made to add empty vertices to SegmentsBuilder.
	ErrVerticesEmpty = fmt.Errorf("vertices cannot be empty")
)

// Segments represents a collection of Vertices used to define geometric shapes.
type Segments []Vertices

// SegmentsBuilder is a builder for constructing Segments objects incrementally.
type SegmentsBuilder struct {
	segments Segments // The segments being built.
	err      error    // Error encountered during the building process.
}

// NewSegmentsBuilder creates and returns a new instance of SegmentsBuilder.
func NewSegmentsBuilder() *SegmentsBuilder {
	return &SegmentsBuilder{}
}

// Add appends a set of Vertices to the SegmentsBuilder.
// If an error has already occurred or the vertices are empty,
// it updates the error field and returns the builder.
func (sb *SegmentsBuilder) Add(vertices Vertices) *SegmentsBuilder {
	if sb.err != nil {
		return sb
	}

	if len(vertices) == 0 {
		sb.err = ErrVerticesEmpty
		return sb
	}

	sb.segments = append(sb.segments, vertices)
	return sb
}

// Build finalizes the Segments and returns them along with any encountered error.
// If an error occurred during the building process, it returns nil and the error.
func (sb *SegmentsBuilder) Build() (Segments, error) {
	if sb.err != nil {
		return nil, sb.err
	}

	return sb.segments, nil
}
