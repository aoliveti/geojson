package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiLineString_BoundingBox(t *testing.T) {
	tests := []struct {
		name     string
		segments Segments
		expected BoundingBox
	}{
		{
			name:     "empty rings",
			segments: Segments{},
			expected: BoundingBox{},
		},
		{
			name: "single line segment",
			segments: Segments{
				{{1, 2}, {3, 4}},
			},
			expected: BoundingBox{1, 2, 3, 4},
		},
		{
			name: "multiple line rings",
			segments: Segments{
				{{1, 2}, {3, 4}},
				{{-1, -2}, {6, 7}},
			},
			expected: BoundingBox{-1, -2, 6, 7},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &MultiLineString{
				segments: tc.segments,
			}
			bbox := m.BoundingBox()
			assert.Equal(t, tc.expected, bbox)
		})
	}
}

func TestMultiLineString_Vertices(t *testing.T) {
	tests := []struct {
		name     string
		segments Segments
		expected Vertices
	}{
		{
			name:     "empty rings",
			segments: Segments{},
			expected: nil,
		},
		{
			name: "single line segment",
			segments: Segments{
				{{1, 2}, {3, 4}},
			},
			expected: Vertices{{1, 2}, {3, 4}},
		},
		{
			name: "multiple line rings",
			segments: Segments{
				{{1, 2}, {3, 4}},
				{{5, 6}, {7, 8}},
			},
			expected: Vertices{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &MultiLineString{
				segments: tc.segments,
			}
			vertices := m.Vertices()
			assert.Equal(t, tc.expected, vertices)
		})
	}
}

func TestMultiLineString_Type(t *testing.T) {
	m := &MultiLineString{}
	expected := TypeMultiLineString
	assert.Equal(t, expected, m.Type())
}

func TestMultiLineString_buildCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expectErr error
	}{
		{
			name: "valid coordinates",
			input: []interface{}{
				[]interface{}{
					[]interface{}{1.0, 2.0},
					[]interface{}{3.0, 4.0},
				},
				[]interface{}{
					[]interface{}{5.0, 6.0},
					[]interface{}{7.0, 8.0},
				},
			},
			expectErr: nil,
		},
		{
			name:      "invalid type",
			input:     "invalid",
			expectErr: ErrInvalidCoordinates,
		},
		{
			name: "inner invalid segment",
			input: []interface{}{
				[]interface{}{
					[]interface{}{1.0, 2.0},
				},
				"invalid",
			},
			expectErr: ErrLineStringTooShort,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &MultiLineString{}
			err := m.buildCoordinates(tc.input)
			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}

func TestMultiLineString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		m        MultiLineString
		expected string
	}{
		{
			name: "multiple rings without bbox",
			m: MultiLineString{
				segments: Segments{
					{{1, 2}, {3, 4}},
					{{5, 6}, {7, 8}},
				},
			},
			expected: `{"type":"MultiLineString","coordinates":[[[1,2],[3,4]],[[5,6],[7,8]]]}`,
		},
		{
			name: "multiple rings with bbox",
			m: MultiLineString{
				segments:      Segments{{{1, 2}, {3, 4}}},
				SerializeBBox: true,
			},
			expected: `{"type":"MultiLineString","coordinates":[[[1,2],[3,4]]],"bbox":[1,2,3,4]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.m.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(data))
		})
	}
}

func TestMultiLineString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expectErr error
	}{
		{
			name:      "valid data",
			data:      `{"type":"MultiLineString","coordinates":[[[1,2],[3,4]],[[5,6],[7,8]]]}`,
			expectErr: nil,
		},
		{
			name:      "multiline too short",
			data:      `{"type":"MultiLineString","coordinates":[]}`,
			expectErr: ErrMultiLineStringTooShort,
		},
		{
			name:      "invalid type",
			data:      `{"type":"Point","coordinates":[1,2]}`,
			expectErr: ErrInvalidTypeField,
		},
		{
			name:      "invalid coordinates",
			data:      `{"type":"MultiLineString","coordinates":"invalid"}`,
			expectErr: ErrInvalidCoordinates,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &MultiLineString{}
			err := m.UnmarshalJSON([]byte(tc.data))
			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}

func TestNewMultiLineString(t *testing.T) {
	tests := []struct {
		name      string
		segments  Segments
		expectErr error
	}{
		{
			name:      "valid single line segment",
			segments:  Segments{{{1, 2}, {3, 4}}},
			expectErr: nil,
		},
		{
			name:      "valid multiple line segments",
			segments:  Segments{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			expectErr: nil,
		},
		{
			name:      "empty segments",
			segments:  Segments{},
			expectErr: ErrMultiLineStringTooShort,
		},
		{
			name:      "line segment too short",
			segments:  Segments{{{1, 2}}},
			expectErr: ErrLineStringTooShort,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewMultiLineString(tc.segments)
			if tc.expectErr != nil {
				require.ErrorIs(t, err, tc.expectErr)
				require.Nil(t, m)
			} else {
				require.NoError(t, err)
				require.NotNil(t, m)
				assert.Equal(t, tc.segments, m.segments)
			}
		})
	}
}

func TestMustMultiLineString(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		segments := Segments{{{1, 2}, {3, 4}}}
		m := MustMultiLineString(segments)
		require.NotNil(t, m)
		assert.Equal(t, segments, m.segments)
	})

	t.Run("invalid input - panic", func(t *testing.T) {
		segments := Segments{}
		assert.PanicsWithError(t, ErrMultiLineStringTooShort.Error(), func() {
			MustMultiLineString(segments)
		})
	})
}

func TestMultiLineString_Segments(t *testing.T) {
	type fields struct {
		segments Segments
	}
	tests := []struct {
		name   string
		fields fields
		want   Segments
	}{
		{
			name: "empty segments",
			fields: fields{
				segments: Segments{},
			},
			want: Segments{},
		},
		{
			name: "single segment",
			fields: fields{
				segments: Segments{{{1, 2}, {3, 4}}},
			},
			want: Segments{{{1, 2}, {3, 4}}},
		},
		{
			name: "multiple segments",
			fields: fields{
				segments: Segments{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
			},
			want: Segments{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiLineString{
				segments: tt.fields.segments,
			}
			assert.Equalf(t, tt.want, m.Segments(), "Segments()")
		})
	}
}
