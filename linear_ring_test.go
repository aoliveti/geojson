package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearRing_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		lr       LinearRing
		expected bool
	}{
		{"valid ring", LinearRing{{0, 0}, {1, 1}, {2, 2}, {0, 0}}, true},
		{"insufficient vertices", LinearRing{{0, 0}, {1, 1}, {0, 0}}, false},
		{"open ring", LinearRing{{0, 0}, {1, 1}, {2, 2}, {3, 3}}, false},
		{"empty ring", LinearRing{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.lr.IsValid(), "IsValid() result mismatch")
		})
	}
}

func TestLinearRing_HasValidSize(t *testing.T) {
	tests := []struct {
		name     string
		lr       LinearRing
		expected bool
	}{
		{"valid size", LinearRing{{0, 0}, {1, 1}, {2, 2}, {0, 0}}, true},
		{"insufficient size", LinearRing{{0, 0}, {1, 1}, {0, 0}}, false},
		{"empty ring", LinearRing{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.lr.HasValidSize(), "HasValidSize() result mismatch")
		})
	}
}

func TestLinearRing_IsClosed(t *testing.T) {
	tests := []struct {
		name     string
		lr       LinearRing
		expected bool
	}{
		{"properly closed", LinearRing{{0, 0}, {1, 1}, {2, 2}, {0, 0}}, true},
		{"not closed", LinearRing{{0, 0}, {1, 1}, {2, 2}}, false},
		{"zero vertices", LinearRing{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.lr.IsClosed(), "IsClosed() result mismatch")
		})
	}
}

func TestNewLinearRing(t *testing.T) {
	tests := []struct {
		name       string
		vertices   Vertices
		expectErr  error
		shouldPass bool
	}{
		{"valid ring", Vertices{{0, 0}, {1, 1}, {2, 2}, {0, 0}}, nil, true},
		{"invalid size", Vertices{{0, 0}, {1, 1}, {0, 0}}, ErrLinearRingSize, false},
		{"not closed", Vertices{{0, 0}, {1, 1}, {2, 2}, {3, 3}}, ErrLinearRingClosed, false},
		{"empty vertices", Vertices{}, ErrLinearRingSize, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lr, err := NewLinearRing(test.vertices)
			if test.shouldPass {
				require.NoError(t, err)
				require.NotNil(t, lr)
			} else {
				require.Error(t, err)
				assert.Equal(t, test.expectErr, err, "Error mismatch")
			}
		})
	}
}

func TestMustLinearRing(t *testing.T) {
	tests := []struct {
		name        string
		vertices    Vertices
		shouldPanic bool
	}{
		{"valid ring", Vertices{{0, 0}, {1, 1}, {2, 2}, {0, 0}}, false},
		{"invalid size", Vertices{{0, 0}, {1, 1}, {0, 0}}, true},
		{"not closed", Vertices{{0, 0}, {1, 1}, {2, 2}}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() {
					MustLinearRing(test.vertices)
				}, "MustLinearRing() did not panic as expected")
			} else {
				assert.NotPanics(t, func() {
					ring := MustLinearRing(test.vertices)
					require.NotNil(t, ring)
				}, "MustLinearRing() panicked unexpectedly")
			}
		})
	}
}

func TestLinearRing_EnsureOrientation(t *testing.T) {
	type args struct {
		shouldBeCounterClockwise bool
	}
	tests := []struct {
		name              string
		lr                LinearRing
		args              args
		expectedIsCorrect bool // Expected correctness of ring orientation after enforcement
	}{
		{"already counterclockwise", LinearRing{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}, args{true}, true},
		{"already clockwise", LinearRing{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}, args{false}, true},
		{"needs reversal for counterclockwise", LinearRing{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}, args{true}, true},
		{"needs reversal for clockwise", LinearRing{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}, args{false}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.lr.EnsureOrientation(tt.args.shouldBeCounterClockwise)
			if tt.args.shouldBeCounterClockwise {
				assert.True(t, tt.lr.IsCounterClockwise(), "Ring is not counterclockwise as expected")
			} else {
				assert.True(t, tt.lr.IsClockwise(), "Ring is not clockwise as expected")
			}
		})
	}
}

func TestLinearRing_IsCounterClockwise(t *testing.T) {
	tests := []struct {
		name string
		lr   LinearRing
		want bool
	}{
		{"properly counterclockwise", LinearRing{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}, true},
		{"properly clockwise", LinearRing{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.lr.IsCounterClockwise(), "IsCounterClockwise()")
		})
	}
}

func TestLinearRing_IsClockwise(t *testing.T) {
	tests := []struct {
		name string
		lr   LinearRing
		want bool
	}{
		{"properly clockwise", LinearRing{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}, true},
		{"properly counterclockwise", LinearRing{{0, 0}, {2, 0}, {2, 2}, {0, 2}, {0, 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.lr.IsClockwise(), "IsClockwise()")
		})
	}
}

func TestLinearRing_Area(t *testing.T) {
	tests := []struct {
		name string
		lr   *LinearRing
		want float64
	}{
		{"simple square", MustLinearRing(Vertices{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}), 4.0},
		{"rectangle", MustLinearRing(Vertices{{0, 0}, {0, 3}, {2, 3}, {2, 0}, {0, 0}}), 6.0},
		{"triangle", MustLinearRing(Vertices{{0, 0}, {2, 0}, {1, 2}, {0, 0}}), 2.0},
		{"complex polygon", MustLinearRing(Vertices{{1, 1}, {4, 1}, {4, 5}, {3, 3}, {2, 4}, {1, 5}, {1, 1}}), 9.0},
		{"zero area (line)", MustLinearRing(Vertices{{0, 0}, {1, 1}, {2, 2}, {0, 0}}), 0.0},
		{"smallest valid ring", MustLinearRing(Vertices{{0, 0}, {1, 1}, {1, 0}, {0, 0}}), 0.5},
		{"area with negative coordinates", MustLinearRing(Vertices{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}, {-1, -1}}), 4.0},
		{"self-intersecting polygon (bowtie)", MustLinearRing(Vertices{{0, 0}, {4, 4}, {4, 0}, {0, 4}, {0, 0}}), 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.lr.Area(), "Area()")
		})
	}
}
