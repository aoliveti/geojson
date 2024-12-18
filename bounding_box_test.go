package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundingBox_Is2D(t *testing.T) {
	tests := []struct {
		name     string
		bbox     BoundingBox
		expected bool
	}{
		{"empty", BoundingBox{}, false},
		{"2D bbox", BoundingBox{0, 0, 1, 1}, true},
		{"3D bbox", BoundingBox{0, 0, 0, 1, 1, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bbox.Is2D())
		})
	}
}

func TestBoundingBox_Is3D(t *testing.T) {
	tests := []struct {
		name     string
		bbox     BoundingBox
		expected bool
	}{
		{"empty", BoundingBox{}, false},
		{"2D bbox", BoundingBox{0, 0, 1, 1}, false},
		{"3D bbox", BoundingBox{0, 0, 0, 1, 1, 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bbox.Is3D())
		})
	}
}

func TestBoundingBox_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		bbox     BoundingBox
		expected bool
	}{
		{"empty", BoundingBox{}, true},
		{"2D bbox", BoundingBox{0, 0, 1, 1}, false},
		{"3D bbox", BoundingBox{0, 0, 0, 1, 1, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bbox.IsZero())
		})
	}
}

func TestBoundingBox_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		bbox     BoundingBox
		expected bool
	}{
		{"empty", BoundingBox{}, true},
		{"2D bbox", BoundingBox{0, 0, 1, 1}, true},
		{"3D bbox", BoundingBox{0, 0, 0, 1, 1, 1}, true},
		{"invalid with 5 coords", BoundingBox{0, 0, 0, 1, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bbox.IsValid())
		})
	}
}

func TestBbox(t *testing.T) {
	tests := []struct {
		name     string
		vertices Vertices
		expected BoundingBox
	}{
		{
			name:     "empty vertices",
			vertices: Vertices{},
			expected: BoundingBox{},
		},
		{
			name: "2D vertices",
			vertices: Vertices{
				{-10.0, 0.0},
				{10.0, 20.0},
			},
			expected: BoundingBox{-10.0, 0.0, 10.0, 20.0},
		},
		{
			name: "3D vertices",
			vertices: Vertices{
				{-10.0, 0.0, 100.0},
				{10.0, 20.0, 200.0},
			},
			expected: BoundingBox{-10.0, 0.0, 100.0, 10.0, 20.0, 200.0},
		},
		{
			name: "mixed 2D and 3D vertices",
			vertices: Vertices{
				{-10.0, 0.0},
				{10.0, 20.0, 200.0},
			},
			expected: BoundingBox{-10.0, 0.0, 0.0, 10.0, 20.0, 200.0},
		},
		{
			name: "mixed 2D and 3D vertices with one negative altitude",
			vertices: Vertices{
				{-10.0, 0.0},
				{10.0, 20.0, -200.0},
			},
			expected: BoundingBox{-10.0, 0.0, -200.0, 10.0, 20.0, 0.0},
		},
		{
			name: "mixed 2D and 3D vertices with negative altitudes",
			vertices: Vertices{
				{-10.0, 0.0, -100.0},
				{10.0, 20.0, -200.0},
			},
			expected: BoundingBox{-10.0, 0.0, -200.0, 10.0, 20.0, -100.0},
		},
	}

	const epsilon = 0.1
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bbox(tt.vertices)

			require.Equal(t, len(tt.expected), len(got), "bbox length mismatch")
			for i, val := range got {
				// epsilon Avoid division by zero error
				assert.InEpsilon(t, tt.expected[i]+epsilon, val+epsilon, epsilon, "bbox value mismatch at index %d", i)
			}
		})
	}
}
