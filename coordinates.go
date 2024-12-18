package geojson

import (
	"encoding/json"
	"fmt"
	"slices"
)

const (
	// LongitudeMin defines the minimum valid longitude.
	LongitudeMin float64 = -180
	// LongitudeMax defines the maximum valid longitude.
	LongitudeMax float64 = 180
	// LatitudeMin defines the minimum valid latitude.
	LatitudeMin float64 = -90
	// LatitudeMax defines the maximum valid latitude.
	LatitudeMax float64 = 90
)

const (
	// Indices for accessing longitude, latitude, and altitude in the coordinates array.
	idxCoordsLng = iota
	idxCoordsLat
	idxCoordsAlt
)

const (
	// coordsMinLen is the minimum number of elements in a valid coordinates array.
	coordsMinLen = 2
	// coordsMaxLen is the maximum number of elements in a valid coordinates array.
	coordsMaxLen = 3
)

var (
	// ErrLongitudeRange is returned when a longitude value is out of range.
	ErrLongitudeRange = fmt.Errorf("longitude must be between -180 and 180")
	// ErrLatitudeRange is returned when a latitude value is out of range.
	ErrLatitudeRange = fmt.Errorf("latitude must be between -90 and 90")
	// ErrCoordinatesSize is returned when the coordinates array does not have 2 or 3 elements.
	ErrCoordinatesSize = fmt.Errorf("coordinates must have 2 or 3 elements")
)

// Coordinates represents a GeoJSON coordinate array.
type Coordinates []float64

// NewCoordinates creates a new Coordinates object from a float64 array.
// Returns an error if the input array is invalid.
func NewCoordinates(v []float64) (*Coordinates, error) {
	if len(v) != coordsMinLen && len(v) != coordsMaxLen {
		return nil, ErrCoordinatesSize
	}

	c := make(Coordinates, len(v))
	c[idxCoordsLng] = v[idxCoordsLng]
	c[idxCoordsLat] = v[idxCoordsLat]

	// Validate longitude and latitude values.
	if err := validateCoordinates(c[idxCoordsLng], c[idxCoordsLat]); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	// Include altitude if present.
	if len(v) == coordsMaxLen {
		c[idxCoordsAlt] = v[idxCoordsAlt]
	}

	return &c, nil
}

// MustCoordinates is a helper that panics on error when creating Coordinates.
// It should only be used when failure is not an acceptable outcome.
func MustCoordinates(v []float64) *Coordinates {
	k, err := NewCoordinates(v)
	if err != nil {
		panic(err)
	}

	return k
}

// Longitude returns the longitude value of the coordinates.
func (c *Coordinates) Longitude() float64 {
	return (*c)[idxCoordsLng]
}

// Latitude returns the latitude value of the coordinates.
func (c *Coordinates) Latitude() float64 {
	return (*c)[idxCoordsLat]
}

// HasAltitude checks if the coordinates include an altitude value.
func (c *Coordinates) HasAltitude() bool {
	return len(*c) == coordsMaxLen
}

// Altitude returns the altitude value of the coordinates.
// This should only be called if HasAltitude() returns true.
func (c *Coordinates) Altitude() float64 {
	return (*c)[idxCoordsAlt]
}

// IsEqual checks if the current Coordinates are equal to the provided Coordinates.
// It returns true if both have the same values in the same order, false otherwise.
func (c *Coordinates) IsEqual(v Coordinates) bool {
	return slices.Compare(*c, v) == 0
}

// String returns a string representation of the coordinates in GeoJSON format.
func (c *Coordinates) String() string {
	if c.HasAltitude() {
		return fmt.Sprintf("[ %g, %g, %g ]", c.Longitude(), c.Latitude(), c.Altitude())
	}

	return fmt.Sprintf("[ %g, %g ]", c.Longitude(), c.Latitude())
}

// UnmarshalJSON implements the json.Unmarshaler interface to parse a GeoJSON coordinates array.
func (c *Coordinates) UnmarshalJSON(data []byte) error {
	var v []float64
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to unmarshal coordinates: %w", err)
	}

	if len(v) != coordsMinLen && len(v) != coordsMaxLen {
		return ErrCoordinatesSize
	}

	// Validate longitude and latitude values.
	if err := validateCoordinates(v[idxCoordsLng], v[idxCoordsLat]); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	*c = v
	return nil
}

// validateCoordinates checks if the provided latitude and longitude are within valid ranges.
func validateCoordinates(longitude, latitude float64) error {
	if longitude < LongitudeMin || longitude > LongitudeMax {
		return ErrLongitudeRange
	}
	if latitude < LatitudeMin || latitude > LatitudeMax {
		return ErrLatitudeRange
	}

	return nil
}

// buildCoordinates constructs a Coordinates object from a generic interface.
// The input must be a slice of interface{} with 2 or 3 float64 elements,
// representing the longitude, latitude, and optionally altitude.
// Returns an error if the input is invalid or contains out-of-range values.
func buildCoordinates(v interface{}) (*Coordinates, error) {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return nil, ErrInvalidCoordinates
	}

	// Ensure the slice contains 2 or 3 elements.
	if len(rawSlice) != coordsMinLen && len(rawSlice) != coordsMaxLen {
		return nil, ErrCoordinatesSize
	}

	slice := make([]float64, len(rawSlice))
	for i, s := range rawSlice {
		switch c := s.(type) {
		case float64:
			slice[i] = c
		case int:
			slice[i] = float64(c)
		default:
			return nil, ErrInvalidCoordinates
		}
	}

	// Validate the longitude and latitude values.
	coords, err := NewCoordinates(slice)
	if err != nil {
		return nil, err
	}

	return coords, nil
}
