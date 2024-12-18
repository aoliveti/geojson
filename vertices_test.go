package geojson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerticesBuilder_Add(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		wantErr bool
	}{
		{
			name:    "valid coordinates",
			input:   []float64{1.0, 2.0},
			wantErr: false,
		},
		{
			name:    "valid 3D coordinates",
			input:   []float64{1.0, 2.0, 3.0},
			wantErr: false,
		},
		{
			name:    "empty slice",
			input:   []float64{},
			wantErr: true,
		},
		{
			name:    "odd number of values",
			input:   []float64{1.0},
			wantErr: true,
		},
		{
			name:    "valid with zeros",
			input:   []float64{0.0, 0.0},
			wantErr: false,
		},
		{
			name:    "invalid type handling",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewVerticesBuilder()
			builder.Add(tt.input)
			_, err := builder.Build()
			if tt.wantErr {
				require.Error(t, err, "expected an error but got none")
			} else {
				require.NoError(t, err, "unexpected error")
			}
		})
	}
}

func TestVerticesBuilder_Build(t *testing.T) {
	tests := []struct {
		name              string
		coordinates       [][]float64
		expectVerticesLen int
		wantErr           bool
	}{
		{
			name: "multiple valid coordinates",
			coordinates: [][]float64{
				{1.0, 2.0},
				{3.0, 4.0},
			},
			expectVerticesLen: 2,
			wantErr:           false,
		},
		{
			name:              "no coordinates added",
			coordinates:       [][]float64{},
			expectVerticesLen: 0,
			wantErr:           false,
		},
		{
			name: "one valid and one invalid coordinate",
			coordinates: [][]float64{
				{1.0, 2.0},
				{},
			},
			expectVerticesLen: 0,
			wantErr:           true,
		},
		{
			name: "invalid coordinates only",
			coordinates: [][]float64{
				{},
				nil,
			},
			expectVerticesLen: 0,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewVerticesBuilder()
			for _, coords := range tt.coordinates {
				builder.Add(coords)
			}
			vertices, err := builder.Build()
			if tt.wantErr {
				require.Error(t, err, "expected an error but got none")
			} else {
				require.NoError(t, err, "unexpected error")
			}
			require.Equal(t, tt.expectVerticesLen, len(vertices), "unexpected number of vertices")
		})
	}
}

func TestNewVerticesBuilder(t *testing.T) {
	t.Run("initialize VerticesBuilder", func(t *testing.T) {
		builder := NewVerticesBuilder()
		require.NotNil(t, builder, "NewVerticesBuilder() returned nil")
		require.NoError(t, builder.err, "NewVerticesBuilder() initialization error")
		require.Nil(t, builder.vertices, "NewVerticesBuilder() vertices is nil, expected empty slice")
	})
}
