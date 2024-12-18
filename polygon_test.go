package geojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygon_Vertices(t *testing.T) {
	tests := []struct {
		name string
		p    *Polygon
		want Vertices
	}{
		{
			name: "empty polygon",
			p:    &Polygon{},
			want: nil,
		},
		{
			name: "polygon with rings",
			p: MustPolygon([]LinearRing{
				*MustLinearRing([]Coordinates{{10, 20}, {30, 40}, {50, 60}, {10, 20}}),
				*MustLinearRing([]Coordinates{{5, 5}, {10, 10}, {15, 5}, {5, 5}}),
			}),
			want: Vertices{
				{10, 20}, {50, 60}, {30, 40}, {10, 20},
				{5, 5}, {10, 10}, {15, 5}, {5, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.Vertices()
			assert.Equal(t, tt.want, got, "Vertices() mismatch")
		})
	}
}

func TestPolygon_Type(t *testing.T) {
	p := &Polygon{}
	assert.Equal(t, TypePolygon, p.Type(), "Type() mismatch")
}

func TestPolygon_buildCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr error
	}{
		{
			name:    "valid coordinates",
			input:   []interface{}{[]interface{}{[]interface{}{0, 0}, []interface{}{10, 0}, []interface{}{0, 10}, []interface{}{0, 0}}},
			wantErr: nil,
		},
		{
			name:    "empty coordinates",
			input:   []interface{}{},
			wantErr: ErrPolygonLinearRingCount,
		},
		{
			name:    "invalid type",
			input:   "invalid",
			wantErr: ErrInvalidCoordinates,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Polygon{}
			err := p.buildCoordinates(tt.input)
			assert.ErrorIs(t, err, tt.wantErr, "buildCoordinates() mismatch")
		})
	}
}

func TestPolygon_MarshalJSON(t *testing.T) {
	p := MustPolygon([]LinearRing{
		*MustLinearRing([]Coordinates{{10, 20}, {30, 40}, {50, 60}, {10, 20}}),
	})
	p.SerializeBBox = true
	data, err := p.MarshalJSON()
	assert.NoError(t, err, "MarshalJSON() error")
	expected := `{"type":"Polygon","coordinates":[[[10,20],[50,60],[30,40],[10,20]]],"bbox":[10,20,50,60]}`
	assert.JSONEq(t, expected, string(data), "MarshalJSON() mismatch")
}

func TestPolygon_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid polygon",
			input:   `{"type":"Polygon","coordinates":[[[10,20],[30,40],[50,60],[10,20]]]}`,
			wantErr: nil,
		},
		{
			name:    "polygon not closed",
			input:   `{"type":"Polygon","coordinates":[[[10,20],[30,40],[50,60],[10,10]]]}`,
			wantErr: ErrLinearRingClosed,
		},
		{
			name:    "empty polygon",
			input:   `{"type":"Polygon","coordinates":[[]]}`,
			wantErr: ErrLinearRingSize,
		},
		{
			name:    "invalid type",
			input:   `{"type":"Invalid","coordinates":[[[10,20],[30,40],[50,60],[10,20]]]}`,
			wantErr: ErrInvalidTypeField,
		},
		{
			name:    "different type",
			input:   `{"type":"Point","coordinates":[10,20]}`,
			wantErr: ErrInvalidTypeField,
		},
		{
			name:    "invalid data type for coordinates",
			input:   `{"type":"Polygon","coordinates":42}`,
			wantErr: ErrInvalidCoordinates,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Polygon{}
			err := json.Unmarshal([]byte(tt.input), p)
			assert.ErrorIs(t, err, tt.wantErr, "UnmarshalJSON() mismatch")
		})
	}
}

func TestNewPolygon(t *testing.T) {
	tests := []struct {
		name    string
		rings   []LinearRing
		wantErr error
	}{
		{
			name:    "valid rings",
			rings:   []LinearRing{*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}})},
			wantErr: nil,
		},
		{
			name:    "no rings",
			rings:   []LinearRing{},
			wantErr: ErrPolygonLinearRingCount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPolygon(tt.rings)
			assert.ErrorIs(t, err, tt.wantErr, "NewPolygon() mismatch")
		})
	}
}

func TestMustPolygon(t *testing.T) {
	tests := []struct {
		name        string
		input       []LinearRing
		shouldPanic bool
	}{
		{
			name: "valid polygon",
			input: []LinearRing{
				*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}),
			},
			shouldPanic: false,
		},
		{
			name:        "no linear rings",
			input:       []LinearRing{},
			shouldPanic: true,
		},
		{
			name: "invalid linear ring",
			input: []LinearRing{
				{}, // Invalid empty ring
			},
			shouldPanic: true,
		},
		{
			name: "non-closed linear ring",
			input: []LinearRing{
				LinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 1}}), // Not closed
			},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Did not expect panic, but got %v", r)
					}
				} else if tt.shouldPanic {
					t.Errorf("Expected panic, but did not get one")
				}
			}()
			_ = MustPolygon(tt.input)
		})
	}
}

func Test_ensureOrientation(t *testing.T) {
	type args struct {
		rings LinearRings
	}
	tests := []struct {
		name     string
		args     args
		expected LinearRings
	}{
		{
			name: "single exterior ring",
			args: args{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}),
				},
			},
			expected: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}), // Should remain unchanged
			},
		},
		{
			name: "exterior and interior ring",
			args: args{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {0, 10}, {10, 0}, {0, 0}}),
					*MustLinearRing([]Coordinates{{1, 1}, {3, 1}, {1, 3}, {1, 1}}),
				},
			},
			expected: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}), // Exterior is CCW
				*MustLinearRing([]Coordinates{{1, 1}, {1, 3}, {3, 1}, {1, 1}}),   // Interior is CW
			},
		},
		{
			name: "already correct orientation",
			args: args{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}), // CCW
					*MustLinearRing([]Coordinates{{1, 1}, {1, 3}, {3, 1}, {1, 1}}),   // CW
				},
			},
			expected: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}), // Should remain unchanged
				*MustLinearRing([]Coordinates{{1, 1}, {1, 3}, {3, 1}, {1, 1}}),   // Should remain unchanged
			},
		},
		{
			name: "inverted orientation",
			args: args{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {0, 10}, {10, 0}, {0, 0}}), // CW
					*MustLinearRing([]Coordinates{{1, 1}, {3, 1}, {1, 3}, {1, 1}}),   // CCW
				},
			},
			expected: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}), // Changed to CCW
				*MustLinearRing([]Coordinates{{1, 1}, {1, 3}, {3, 1}, {1, 1}}),   // Changed to CW
			},
		},
		{
			name: "empty rings",
			args: args{
				rings: LinearRings{},
			},
			expected: LinearRings{}, // Should remain empty
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureOrientation(tt.args.rings)
			assert.Equal(t, tt.expected, tt.args.rings, "ensureOrientation() mismatch for %v", tt.name)
		})
	}
}

func TestPolygon_LinearRings(t *testing.T) {
	type fields struct {
		rings LinearRings
	}
	tests := []struct {
		name   string
		fields fields
		want   LinearRings
	}{
		{
			name: "empty LinearRings",
			fields: fields{
				rings: LinearRings{},
			},
			want: LinearRings{},
		},
		{
			name: "single LinearRing",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {1, 0}, {0, 1}, {0, 0}}),
				},
			},
			want: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {1, 0}, {0, 1}, {0, 0}}),
			},
		},
		{
			name: "multiple LinearRings",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {2, 0}, {0, 2}, {0, 0}}),
					*MustLinearRing([]Coordinates{{1, 1}, {1, 1.5}, {1.5, 1}, {1, 1}}),
				},
			},
			want: LinearRings{
				*MustLinearRing([]Coordinates{{0, 0}, {2, 0}, {0, 2}, {0, 0}}),
				*MustLinearRing([]Coordinates{{1, 1}, {1, 1.5}, {1.5, 1}, {1, 1}}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Polygon{
				rings: tt.fields.rings,
			}
			assert.Equalf(t, tt.want, p.LinearRings(), "LinearRings()")
		})
	}
}

func TestPolygon_OuterRing(t *testing.T) {
	type fields struct {
		rings LinearRings
	}
	tests := []struct {
		name   string
		fields fields
		want   LinearRing
	}{
		{
			name: "polygon with no rings",
			fields: fields{
				rings: LinearRings{},
			},
			want: nil,
		},
		{
			name: "polygon with only outer ring",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}),
				},
			},
			want: *MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}),
		},
		{
			name: "polygon with outer and inner rings",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {2, 0}, {0, 2}, {0, 0}}),
					*MustLinearRing([]Coordinates{{1, 1}, {1, 1.5}, {1.5, 1}, {1, 1}}),
				},
			},
			want: *MustLinearRing([]Coordinates{{0, 0}, {2, 0}, {0, 2}, {0, 0}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Polygon{
				rings: tt.fields.rings,
			}
			assert.Equalf(t, tt.want, p.OuterRing(), "OuterRing()")
		})
	}
}

func TestPolygon_InnerRings(t *testing.T) {
	type fields struct {
		rings LinearRings
	}
	tests := []struct {
		name   string
		fields fields
		want   LinearRings
	}{
		{
			name: "empty ring",
			fields: fields{
				rings: LinearRings{},
			},
		},
		{
			name: "polygon with no inner rings",
			fields: fields{
				rings: LinearRings{*MustLinearRing([]Coordinates{{10, 10}, {20, 20}, {10, 20}, {10, 10}})},
			},
			want: LinearRings{},
		},
		{
			name: "polygon with one inner ring",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {10, 0}, {0, 10}, {0, 0}}),
					*MustLinearRing([]Coordinates{{3, 3}, {5, 5}, {5, 3}, {3, 3}}),
				},
			},
			want: LinearRings{
				*MustLinearRing([]Coordinates{{3, 3}, {5, 5}, {5, 3}, {3, 3}}),
			},
		},
		{
			name: "polygon with multiple inner rings",
			fields: fields{
				rings: LinearRings{
					*MustLinearRing([]Coordinates{{0, 0}, {2, 0}, {0, 2}, {0, 0}}),
					*MustLinearRing([]Coordinates{{1, 1}, {1, 1.5}, {1.5, 1}, {1, 1}}),
					*MustLinearRing([]Coordinates{{0.5, 0.5}, {1, 1}, {0.5, 1}, {0.5, 0.5}}),
				},
			},
			want: LinearRings{
				*MustLinearRing([]Coordinates{{1, 1}, {1, 1.5}, {1.5, 1}, {1, 1}}),
				*MustLinearRing([]Coordinates{{0.5, 0.5}, {1, 1}, {0.5, 1}, {0.5, 0.5}}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Polygon{
				rings: tt.fields.rings,
			}
			assert.Equalf(t, tt.want, p.InnerRings(), "InnerRings()")
		})
	}
}
