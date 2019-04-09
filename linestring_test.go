package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLineString struct {
	layout     Layout
	stride     int
	coords     []Coord
	flatCoords []float64
	bounds     *Bounds
}

func assertLineStringEquals(t *testing.T, expected *testLineString, actual *LineString) {
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

func TestLineString(t *testing.T) {
	for _, c := range []struct {
		ls  *LineString
		tls *testLineString
	}{
		{
			ls: NewLineString(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}}),
			tls: &testLineString{
				layout:     XY,
				stride:     2,
				coords:     []Coord{{1, 2}, {3, 4}, {5, 6}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6},
				bounds:     NewBounds(XY).Set(1, 2, 5, 6),
			},
		},
		{
			ls: NewLineString(XYZ).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			tls: &testLineString{
				layout:     XYZ,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				bounds:     NewBounds(XYZ).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			ls: NewLineString(XYM).MustSetCoords([]Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}),
			tls: &testLineString{
				layout:     XYM,
				stride:     3,
				coords:     []Coord{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				bounds:     NewBounds(XYM).Set(1, 2, 3, 7, 8, 9),
			},
		},
		{
			ls: NewLineString(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			tls: &testLineString{
				layout:     XYZM,
				stride:     4,
				coords:     []Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
				flatCoords: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				bounds:     NewBounds(XYZM).Set(1, 2, 3, 4, 9, 10, 11, 12),
			},
		},
	} {
		assertLineStringEquals(t, c.tls, c.ls)
	}
}

func TestLineStringClone(t *testing.T) {
	p1 := NewLineString(XY).MustSetCoords([]Coord{{1, 2}, {3, 4}, {5, 6}})
	assert.False(t, aliases(p1.FlatCoords(), p1.Clone().FlatCoords()))
}

func TestLineStringInterpolate(t *testing.T) {
	ls := NewLineString(XYM).MustSetCoords([]Coord{{1, 2, 0}, {2, 4, 1}, {3, 8, 2}})
	for _, c := range []struct {
		val float64
		dim int
		i   int
		f   float64
	}{
		{val: -0.5, dim: 2, i: 0, f: 0.0},
		{val: 0.0, dim: 2, i: 0, f: 0.0},
		{val: 0.5, dim: 2, i: 0, f: 0.5},
		{val: 1.0, dim: 2, i: 1, f: 0.0},
		{val: 1.5, dim: 2, i: 1, f: 0.5},
		{val: 2.0, dim: 2, i: 2, f: 0.0},
		{val: 2.5, dim: 2, i: 2, f: 0.0},
	} {
		i, f := ls.Interpolate(c.val, c.dim)
		assert.Equal(t, c.i, i)
		assert.Equal(t, c.f, f)
	}
}

func TestLineStringReserve(t *testing.T) {
	ls := NewLineString(XYZM)
	assert.Equal(t, 0, cap(ls.flatCoords))
	ls.Reserve(2)
	assert.Equal(t, 8, cap(ls.flatCoords))
}

func TestLineStringStrideMismatch(t *testing.T) {
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
		_, err := NewLineString(c.layout).SetCoords(c.coords)
		assert.Equal(t, c.err, err)
	}
}
