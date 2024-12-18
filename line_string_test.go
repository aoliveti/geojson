package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineString_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		line     *LineString
		expected Vertices
	}{
		{
			name:     "empty vertices",
			line:     &LineString{vertices: nil},
			expected: nil,
		},
		{
			name:     "single vertex",
			line:     &LineString{vertices: Vertices{{1.1, 2.2}}},
			expected: Vertices{{1.1, 2.2}},
		},
		{
			name:     "multiple vertices",
			line:     &LineString{vertices: Vertices{{1.1, 2.2}, {3.3, 4.4}}},
			expected: Vertices{{1.1, 2.2}, {3.3, 4.4}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.line.Vertices())
		})
	}
}

func TestLineString_Type(t *testing.T) {
	line := &LineString{}
	expected := TypeLineString

	assert.Equal(t, expected, line.Type())
}

func TestLineString_BoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		line     *LineString
		expected BoundingBox
	}{
		{
			name:     "empty vertices",
			line:     &LineString{vertices: nil},
			expected: BoundingBox{},
		},
		{
			name:     "valid bounding box",
			line:     &LineString{vertices: Vertices{{1.1, 2.2}, {3.3, 4.4}}},
			expected: bbox(Vertices{{1.1, 2.2}, {3.3, 4.4}}),
		},
		{
			name:     "3D coordinates",
			line:     &LineString{vertices: Vertices{{1.1, 2.2, 5.5}, {3.3, 4.4, 6.6}}},
			expected: bbox(Vertices{{1.1, 2.2, 5.5}, {3.3, 4.4, 6.6}}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.line.BoundingBox())
		})
	}
}

func TestLineString_buildCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr error
	}{
		{
			name:      "invalid input type",
			input:     "invalid",
			expectErr: ErrInvalidCoordinates,
		},
		{
			name:      "lesser then two coordinates",
			input:     []interface{}{[]interface{}{1.1, 2.2}},
			expectErr: ErrLineStringTooShort,
		},
		{
			name:      "valid coordinates",
			input:     []interface{}{[]interface{}{1.1, 2.2}, []interface{}{3.3, 4.4}},
			expectErr: nil,
		},
		{
			name:      "invalid coordinate in slice",
			input:     []interface{}{[]interface{}{1.1, 2.2}, "invalid"},
			expectErr: ErrInvalidCoordinates,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line := &LineString{}
			err := line.buildCoordinates(tc.input)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLineString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		line           *LineString
		expectedOutput string
	}{
		{
			name:           "empty vertices",
			line:           &LineString{vertices: Vertices{}},
			expectedOutput: `{"type":"LineString","coordinates":[]}`,
		},
		{
			name:           "with vertices",
			line:           &LineString{vertices: Vertices{{1.1, 2.2}, {3.3, 4.4}}},
			expectedOutput: `{"type":"LineString","coordinates":[[1.1,2.2],[3.3,4.4]]}`,
		},
		{
			name: "with BBox serialization",
			line: &LineString{
				vertices:      Vertices{{1.1, 2.2}, {3.3, 4.4}},
				SerializeBBox: true,
			},
			expectedOutput: `{"type":"LineString","coordinates":[[1.1,2.2],[3.3,4.4]],"bbox":[1.1,2.2,3.3,4.4]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tc.line.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tc.expectedOutput, string(output))
		})
	}
}

func TestLineString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		inputJSON string
		expectErr error
	}{
		{
			name:      "valid LineString",
			inputJSON: `{"type":"LineString","coordinates":[[1.1,2.2],[3.3,4.4]]}`,
			expectErr: nil,
		},
		{
			name:      "invalid type",
			inputJSON: `{"type":"Point","coordinates":[3.3,4.4]}`,
			expectErr: ErrInvalidTypeField,
		},
		{
			name:      "invalid JSON structure",
			inputJSON: `{"type":"LineString","coordinates":"invalid"}`,
			expectErr: ErrInvalidCoordinates,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line := &LineString{}
			err := line.UnmarshalJSON([]byte(tc.inputJSON))
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewLineStringFromVertices(t *testing.T) {
	tests := []struct {
		name      string
		vertices  Vertices
		expectErr error
	}{
		{
			name:      "valid vertices",
			vertices:  Vertices{{1.1, 2.2}, {3.3, 4.4}},
			expectErr: nil,
		},
		{
			name:      "less than two vertices",
			vertices:  Vertices{{1.1, 2.2}},
			expectErr: ErrLineStringTooShort,
		},
		{
			name:      "no vertices",
			vertices:  Vertices{},
			expectErr: ErrLineStringTooShort,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line, err := NewLineString(tc.vertices)
			if tc.expectErr != nil {
				assert.Nil(t, line)
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NotNil(t, line)
				assert.Equal(t, tc.vertices, line.vertices)
				assert.NoError(t, err)
			}
		})
	}
}

func TestMustLineStringFromVertices(t *testing.T) {
	tests := []struct {
		name         string
		vertices     Vertices
		expectPanics bool
	}{
		{
			name:         "valid vertices",
			vertices:     Vertices{{1.1, 2.2}, {3.3, 4.4}},
			expectPanics: false,
		},
		{
			name:         "empty vertices",
			vertices:     Vertices{},
			expectPanics: true,
		},
		{
			name:         "less than two vertices",
			vertices:     Vertices{{1.1, 2.2}},
			expectPanics: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanics {
				assert.Panics(t, func() { MustLineString(tc.vertices) })
			} else {
				assert.NotPanics(t, func() {
					line := MustLineString(tc.vertices)
					assert.NotNil(t, line)
					assert.Equal(t, tc.vertices, line.vertices)
				})
			}
		})
	}
}
