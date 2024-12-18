package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentsBuilder_Add(t *testing.T) {
	tests := []struct {
		name             string
		vertices         Vertices
		initialErr       error
		expectedErr      error
		expectedSegments Segments
	}{
		{"valid vertices", Vertices{{1, 1}, {2, 2}}, nil, nil, Segments{Vertices{{1, 1}, {2, 2}}}},
		{"empty vertices", nil, nil, ErrVerticesEmpty, nil},
		{"builder with error", Vertices{{1, 1}, {2, 2}}, ErrVerticesEmpty, ErrVerticesEmpty, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &SegmentsBuilder{
				err: tt.initialErr,
			}
			resultBuilder := builder.Add(tt.vertices)

			// Assert error
			if tt.expectedErr != nil {
				require.Error(t, resultBuilder.err)
				assert.ErrorIs(t, resultBuilder.err, tt.expectedErr)
			} else {
				assert.NoError(t, resultBuilder.err)
			}

			// Assert rings
			assert.Equal(t, tt.expectedSegments, resultBuilder.segments)
		})
	}
}

func TestSegmentsBuilder_Build(t *testing.T) {
	tests := []struct {
		name             string
		initialErr       error
		initialSegments  Segments
		expectedErr      error
		expectedSegments Segments
	}{
		{"valid build", nil, Segments{Vertices{{1, 1}, {2, 2}}}, nil, Segments{Vertices{{1, 1}, {2, 2}}}},
		{"builder with error", ErrVerticesEmpty, nil, ErrVerticesEmpty, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &SegmentsBuilder{
				segments: tt.initialSegments,
				err:      tt.initialErr,
			}
			segments, err := builder.Build()

			// Assert error
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			// Assert rings
			assert.Equal(t, tt.expectedSegments, segments)
		})
	}
}

func TestNewSegmentsBuilder(t *testing.T) {
	builder := NewSegmentsBuilder()

	// Assert builder is not nil
	require.NotNil(t, builder)

	// Assert initial state
	assert.Empty(t, builder.segments)
	assert.NoError(t, builder.err)
}
