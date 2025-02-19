package geojson

import (
	"math"
)

const (
	// bboxSize2D specifies the size of a 2D bounding box.
	bboxSize2D = 4
	// bboxSize3D specifies the size of a 3D bounding box.
	bboxSize3D = 6
)

// BoundingBoxer is an interface that defines methods for calculating the bounding box
// and retrieving the vertices of a geometry.
type BoundingBoxer interface {
	// BoundingBox returns the bounding box for a geometry.
	BoundingBox() BoundingBox
	// Vertices returns the vertices that form the geometry.
	Vertices() Vertices
}

// BoundingBox represents a geographic bounding box, either 2D or 3D, as a slice of float64 values.
type BoundingBox []float64

// Is2D checks if the bounding box is a valid 2D bounding box.
func (b *BoundingBox) Is2D() bool {
	return len(*b) == bboxSize2D
}

// Is3D checks if the bounding box is a valid 3D bounding box.
func (b *BoundingBox) Is3D() bool {
	return len(*b) == bboxSize3D
}

// IsZero checks if the bounding box is empty (contains no values).
func (b *BoundingBox) IsZero() bool {
	return len(*b) == 0
}

// IsValid checks if the bounding box is either empty, a 2D bounding box, or a 3D bounding box.
func (b *BoundingBox) IsValid() bool {
	return b.IsZero() || b.Is2D() || b.Is3D()
}

// updateRange updates the minimum and maximum float64 values based on the provided value.
func updateRange(value float64, minVal, maxVal *float64) {
	if value < *minVal {
		*minVal = value
	}
	if value > *maxVal {
		*maxVal = value
	}
}

// bbox calculates the minimum bounding box for a set of vertices, supporting both 2D and 3D bounding boxes.
// It iterates over the provided vertices to determine the minimum and maximum bounds for longitude,
// latitude, and optionally altitude, constructing a bounding box based on the data available.
func bbox(vertices Vertices) BoundingBox {
	// If no vertices are provided, return an empty bounding box.
	if len(vertices) == 0 {
		return BoundingBox{}
	}

	// Initialize the minimum and maximum values for longitude, latitude, and altitude.
	minLng, minLat, maxLng, maxLat := LongitudeMax, LatitudeMax, LongitudeMin, LatitudeMin
	minAlt, maxAlt := math.MaxFloat64, -math.MaxFloat64

	altitudeCount := 0 // Tracks the number of vertices with altitude information.

	// Iterate over each vertex to calculate bounding box boundaries.
	for _, v := range vertices {
		// Update minimum and maximum longitude and latitude values.
		updateRange(v.Longitude(), &minLng, &maxLng)
		updateRange(v.Latitude(), &minLat, &maxLat)

		// Update minimum and maximum altitude values if altitude is present.
		if v.HasAltitude() {
			altitudeCount++ // Increment 3D vertex counter.

			updateRange(v.Altitude(), &minAlt, &maxAlt)
		}
	}

	// Adjust the altitude bounds for vertices that do not include altitude.
	if altitudeCount != len(vertices) {
		for _, v := range vertices {
			if !v.HasAltitude() {
				// If altitude is missing, ensure it defaults to 0 within the bounds.
				updateRange(0, &minAlt, &maxAlt)
			}
		}
	}

	if altitudeCount > 0 {
		// Return a 3D bounding box with longitude, latitude, and altitude values.
		return BoundingBox{
			minLng, // Minimum longitude.
			minLat, // Minimum latitude.
			minAlt, // Minimum altitude.
			maxLng, // Maximum longitude.
			maxLat, // Maximum latitude.
			maxAlt, // Maximum altitude.
		}
	}

	// Return a 2D bounding box when no altitude information exists.
	return BoundingBox{
		minLng, // Minimum longitude.
		minLat, // Minimum latitude.
		maxLng, // Maximum longitude.
		maxLat, // Maximum latitude.
	}
}
