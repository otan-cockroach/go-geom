package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLinearRing struct {
	layout     Layout
	stride     int
	coords     []Coord
	flatCoords []float64
	bounds     *Bounds
}

func assertLinearRingEquals(t *testing.T, expected *testLinearRing, actual *LinearRing) {
	assert.NoError(t, actual.verify())
	assert.Equal(t, expected.layout, actual.Layout())
	assert.Equal(t, expected.stride, actual.Stride())
	assert.Equal(t, expected.coords, actual.Coords())
	assert.Equal(t, expected.bounds, actual.Bounds())
	assert.Equal(t, len(expected.coords), actual.NumCoords())
	for i, c := range expected.coords {
		assert.Equal(t, c, actual.Coord(i))
	}
}

func TestLinearRing(t *testing.T) {
	for _, c := range []struct {
		lr  *LinearRing
		tlr *testLinearRing
	}{
		{
			lr: NewLinearRing(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}}),
			tlr: &testLinearRing{
				layout:     XY,
				stride:     2,
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			lr: NewLinearRing(XYZ).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			tlr: &testLinearRing{
				layout:     XYZ,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			lr: NewLinearRing(XYM).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			tlr: &testLinearRing{
				layout:     XYM,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			lr: NewLinearRing(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			tlr: &testLinearRing{
				layout:     XYZM,
				stride:     4,
				coords:     []Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 9, 10, 11, 12),
			},
		},
	} {
		assertLinearRingEquals(t, c.tlr, c.lr)
	}
}

func TestLinearRingClone(t *testing.T) {
	p1 := NewLinearRing(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}})
	assert.False(t, aliases(p1.FlatCoords(), p1.Clone().FlatCoords()))
}

func TestLinearRingStrideMismatch(t *testing.T) {
	for _, c := range []struct {
		layout Layout
		coords []Coord
		err    error
	}{
		{
			layout: XY,
			coords: nil,
			err:    nil,
		},
		{
			layout: XY,
			coords: []Coord{},
			err:    nil,
		},
		{
			layout: XY,
			coords: []Coord{{1, 2}, {}},
			err:    ErrStrideMismatch{Got: 0, Want: 2},
		},
		{
			layout: XY,
			coords: []Coord{{1, 2}, {1}},
			err:    ErrStrideMismatch{Got: 1, Want: 2},
		},
		{
			layout: XY,
			coords: []Coord{{1, 2}, {3, 4}},
			err:    nil,
		},
		{
			layout: XY,
			coords: []Coord{{1, 2}, {3, 4, 5}},
			err:    ErrStrideMismatch{Got: 3, Want: 2},
		},
	} {
		_, err := NewLinearRing(c.layout).SetCoords(c.coords)
		assert.Equal(t, c.err, err)
	}
}
