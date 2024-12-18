package geojson

import (
	"encoding/json"
	"fmt"
)

var (
	// ErrPolygonLinearRingCount is an error indicating that a polygon must consist of at least one linear ring.
	ErrPolygonLinearRingCount = fmt.Errorf("polygon must have at least one linear ring")
)

// Polygon represents a geometric rings defined by a series of rings.
type Polygon struct {
	rings         LinearRings // The rings that comprise the polygon.
	SerializeBBox bool        // Flag to indicate if the bounding box should be serialized.
}

// BoundingBox calculates and returns the minimum bounding box for the polygon.
func (p *Polygon) BoundingBox() BoundingBox {
	return bbox(p.Vertices())
}

// Vertices retrieves all the vertices that make up the polygon, combining all the vertices from its rings.
func (p *Polygon) Vertices() Vertices {
	var v Vertices
	for _, segment := range p.rings {
		v = append(v, segment...)
	}
	return v
}

// Type returns the geometry type of the polygon, which is TypePolygon.
func (p *Polygon) Type() GeometryType {
	return TypePolygon
}

// LinearRings returns the collection of linear rings that make up the polygon.
// The first ring represents the outer boundary, and subsequent rings represent holes.
func (p *Polygon) LinearRings() LinearRings {
	return p.rings
}

// OuterRing returns the outer ring (boundary) of the Polygon.
// If the Polygon has no rings, it returns nil.
func (p *Polygon) OuterRing() LinearRing {
	if len(p.rings) == 0 {
		return nil
	}

	return p.rings[0]
}

// InnerRings returns the inner rings (holes) of the Polygon.
// If the Polygon has no rings, it returns nil.
func (p *Polygon) InnerRings() LinearRings {
	if len(p.rings) == 0 {
		return nil
	}

	return p.rings[1:]
}

// MarshalJSON converts the polygon into its JSON representation as per the GeoJSON specification.
// If SerializeBBox is enabled, the bounding box will also be included in the output.
func (p *Polygon) MarshalJSON() ([]byte, error) {
	// Prepare the GeoJSON output structure.
	out := geometryJSONOutput{
		Type:        p.Type(),
		Coordinates: p.rings,
	}

	// Include the bounding box if SerializeBBox is enabled.
	if p.SerializeBBox {
		out.BBox = p.BoundingBox()
	}

	// Convert the structure to JSON bytes.
	return json.Marshal(&out)
}

// UnmarshalJSON parses the polygon data from its JSON representation.
// It ensures the parsed data matches the structure of a valid polygon.
func (p *Polygon) UnmarshalJSON(data []byte) error {
	gw := &GeometryObject{}

	// Parse the JSON data into a temporary GeometryObject structure.
	if err := gw.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", p.Type(), err)
	}

	// Check if the geometry type is a Polygon.
	g, ok := gw.geometry.(*Polygon)
	if !ok {
		return ErrInvalidTypeField
	}

	p.rings = g.rings

	return nil
}

// NewPolygon creates a new Polygon instance initialized with the provided linear rings.
// Returns an error if the number of rings is zero.
func NewPolygon(rings LinearRings) (*Polygon, error) {
	// Validate the input to ensure at least one ring is provided.
	if len(rings) == 0 {
		return nil, ErrPolygonLinearRingCount
	}

	for _, ring := range rings {
		if !ring.HasValidSize() {
			return nil, ErrLinearRingSize
		}
		if !ring.IsClosed() {
			return nil, ErrLinearRingClosed
		}
	}

	ensureOrientation(rings)

	return &Polygon{rings: rings}, nil
}

// MustPolygon creates a new Polygon and panics if the provided rings are invalid.
// This is a helper function for scenarios where error handling can be deferred to the caller.
func MustPolygon(rings LinearRings) *Polygon {
	polygon, err := NewPolygon(rings)
	if err != nil {
		panic(err)
	}

	return polygon
}

// buildCoordinates populates the polygon's rings from the provided raw coordinate data.
// It validates and converts the raw data into a series of segments representing the rings of the polygon.
func (p *Polygon) buildCoordinates(v interface{}) error {
	rawSlice, ok := v.([]interface{})
	if !ok {
		return ErrInvalidCoordinates
	}

	// Loop through each raw ring representation and convert it into LinearRings.
	rings := make(LinearRings, len(rawSlice))
	for i, r := range rawSlice {
		rawRing, ok := r.([]interface{})
		if !ok {
			return ErrInvalidCoordinates
		}

		ring := make(Vertices, len(rawRing))
		for j, rv := range rawRing {
			coords, err := buildCoordinates(rv)
			if err != nil {
				return err
			}

			ring[j] = *coords
		}

		// Create a LinearRing from the vertices and validate it.
		lr, err := NewLinearRing(ring)
		if err != nil {
			return err
		}

		rings[i] = *lr
	}

	if len(rings) == 0 {
		return ErrPolygonLinearRingCount
	}

	ensureOrientation(rings)

	p.rings = rings

	return nil
}

// ensureOrientation ensures the rings in a LinearRings collection
// are properly oriented according to their roles in a polygon.
// The first ring (outer ring) is oriented in a counterclockwise direction,
// while all inner rings (holes) are oriented in a clockwise direction.
func ensureOrientation(rings LinearRings) {
	if len(rings) == 0 {
		return
	}

	// Ensure the first ring is oriented in a counterclockwise direction.
	rings[0].EnsureOrientation(true)
	// Ensure all inner rings are oriented in a clockwise direction.
	for i := 1; i < len(rings); i++ {
		rings[i].EnsureOrientation(false)
	}
}
