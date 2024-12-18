package geojson

import (
	"errors"
	"math"
	"slices"
)

const (
	// LinearRingMinimumSize defines the minimum number of coordinates required for a valid LinearRing.
	LinearRingMinimumSize = 4
)

var (
	// ErrLinearRingSize is returned when a LinearRing does not have enough coordinates.
	ErrLinearRingSize = errors.New("linear ring must have at least 4 coordinates")

	// ErrLinearRingClosed is returned when a LinearRing is not closed.
	ErrLinearRingClosed = errors.New("linear ring must be closed")
)

// LinearRing represents a closed linear ring, built from a slice of Vertices.
type LinearRing Vertices

// LinearRings represents a collection of LinearRing objects.
type LinearRings []LinearRing

// IsValid checks if the LinearRing is valid by ensuring it has a valid size
// and is closed.
func (lr *LinearRing) IsValid() bool {
	return lr.HasValidSize() && lr.IsClosed()
}

// HasValidSize verifies if the LinearRing has the minimum required number
// of coordinates.
func (lr *LinearRing) HasValidSize() bool {
	return len(*lr) >= LinearRingMinimumSize
}

// IsClosed checks if the LinearRing is closed by ensuring the first and last
// coordinates are equal.
func (lr *LinearRing) IsClosed() bool {
	if len(*lr) == 0 {
		return false
	}

	first, last := (*lr)[0], (*lr)[len(*lr)-1]
	return first.IsEqual(last)
}

// IsCounterClockwise determines if the LinearRing vertices are ordered in a counterclockwise direction.
// The calculation is based on the signed area of the LinearRing.
// If the result is positive, the vertices are ordered counterclockwise.
func (lr *LinearRing) IsCounterClockwise() bool {
	return signedArea(*lr) > 0
}

// IsClockwise determines if the vertices of the LinearRing
// are ordered in a clockwise direction by checking the signed area.
func (lr *LinearRing) IsClockwise() bool {
	return !lr.IsCounterClockwise()
}

// Area computes the absolute area of a LinearRing.
// It calculates the area using the signed area function, ensuring the result is always positive.
func (lr *LinearRing) Area() float64 {
	return math.Abs(signedArea(*lr))
}

// EnsureOrientation ensures the LinearRing vertices are ordered in the desired direction.
// If the current order is different from the expected order, it reverses the vertices.
// The parameter shouldBeCounterClockwise determines the desired orientation:
// true for counterclockwise, false for clockwise.
func (lr *LinearRing) EnsureOrientation(shouldBeCounterClockwise bool) {
	isCounterClockwise := lr.IsCounterClockwise()
	if shouldBeCounterClockwise == isCounterClockwise {
		return
	}

	slices.Reverse(*lr)
}

// NewLinearRing creates a new LinearRing from the provided vertices.
// It returns an error if the LinearRing has an invalid size or is not closed.
func NewLinearRing(vertices Vertices) (*LinearRing, error) {
	lr := LinearRing(vertices)

	if !lr.HasValidSize() {
		return nil, ErrLinearRingSize
	}

	if !lr.IsClosed() {
		return nil, ErrLinearRingClosed
	}

	return &lr, nil
}

// MustLinearRing creates a new LinearRing from the provided vertices.
// It panics if the LinearRing is invalid.
func MustLinearRing(vertices Vertices) *LinearRing {
	lr, err := NewLinearRing(vertices)
	if err != nil {
		panic(err)
	}

	return lr
}

// signedArea calculates the signed area of a LinearRing using the shoelace formula.
// The formula is: Area = 0.5 * Σ (x(i) * y(i+1) - x(i+1) * y(i))
// A positive result indicates that the vertices are ordered counterclockwise,
// while a negative result indicates that they are ordered clockwise.
func signedArea(ring LinearRing) float64 {
	var v float64
	for i := 0; i < len(ring)-1; i++ {
		x1, y1 := ring[i][idxCoordsLng], ring[i][idxCoordsLat]
		x2, y2 := ring[i+1][idxCoordsLng], ring[i+1][idxCoordsLat]
		v += x1*y2 - x2*y1
	}

	return v * 0.5
}
