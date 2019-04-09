package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundsExtend(t *testing.T) {
	for _, tc := range []struct {
		b    *Bounds
		g    T
		want *Bounds
	}{
		{
			b:    NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
			g:    NewPoint(XY).MustSetCoords(Coord{10, -10}),
			want: NewBounds(XY).SetCoords(Coord{0, -10}, Coord{10, 0}),
		},
		{
			b:    NewBounds(XY).SetCoords(Coord{-100, -100}, Coord{100, 100}),
			g:    NewPoint(XY).MustSetCoords(Coord{-10, 10}),
			want: NewBounds(XY).SetCoords(Coord{-100, -100}, Coord{100, 100}),
		},
		{
			b:    NewBounds(XYZ).SetCoords(Coord{0, 0, -1}, Coord{10, 10, 1}),
			g:    NewPoint(XY).MustSetCoords(Coord{5, -10}),
			want: NewBounds(XYZ).SetCoords(Coord{0, -10, -1}, Coord{10, 10, 1}),
		},
		{
			b:    NewBounds(XYZ).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:    NewPoint(XYZ).MustSetCoords(Coord{5, -10, 3}),
			want: NewBounds(XYZ).SetCoords(Coord{0, -10, 0}, Coord{10, 10, 10}),
		},
		{
			b:    NewBounds(XYZ).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:    NewMultiPoint(XYM).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			want: NewBounds(XYZM).SetCoords(Coord{-1, -1, 0, -1}, Coord{11, 11, 10, 11}),
		},
		{
			b:    NewBounds(XY).SetCoords(Coord{0, 0}, Coord{10, 10}),
			g:    NewMultiPoint(XYM).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			want: NewBounds(XYM).SetCoords(Coord{-1, -1, -1}, Coord{11, 11, 11}),
		},
		{
			b:    NewBounds(XY).SetCoords(Coord{0, 0}, Coord{10, 10}),
			g:    NewMultiPoint(XYZ).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			want: NewBounds(XYZ).SetCoords(Coord{-1, -1, -1}, Coord{11, 11, 11}),
		},
		{
			b:    NewBounds(XYM).SetCoords(Coord{0, 0, 0}, Coord{10, 10, 10}),
			g:    NewMultiPoint(XYZ).MustSetCoords([]Coord{{-1, -1, -1}, {11, 11, 11}}),
			want: NewBounds(XYZM).SetCoords(Coord{-1, -1, -1, 0}, Coord{11, 11, 11, 10}),
		},
	} {
		assert.Equal(t, tc.want, tc.b.Clone().Extend(tc.g))
	}
}

func TestBoundsIsEmpty(t *testing.T) {
	for _, testData := range []struct {
		bounds  Bounds
		isEmpty bool
	}{
		{
			bounds:  Bounds{layout: XY, min: Coord{0, 0}, max: Coord{-1, -1}},
			isEmpty: true,
		}, {
			bounds:  Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			isEmpty: false,
		},
		{
			bounds:  Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			isEmpty: false,
		},
	} {
		copy := Bounds{layout: testData.bounds.layout, min: testData.bounds.min, max: testData.bounds.max}
		for j := 0; j < 10; j++ {
			// do multiple checks to verify no obvious side effects are caused
			assert.Equal(t, testData.isEmpty, copy.IsEmpty())
			assert.Equal(t, testData.bounds, copy)
		}
	}
}

func TestBoundsOverlaps(t *testing.T) {
	for _, testData := range []struct {
		bounds, other Bounds
		overlaps      bool
	}{
		{
			bounds:   Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			other:    Bounds{layout: XY, min: Coord{-10, 0}, max: Coord{-5, 10}},
			overlaps: false,
		},
		{
			bounds:   Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			other:    Bounds{layout: XY, min: Coord{-10, 0}, max: Coord{-5, 10}},
			overlaps: true,
		},
		{
			bounds:   Bounds{layout: XY, min: Coord{1, 1}, max: Coord{5, 5}},
			other:    Bounds{layout: XY, min: Coord{-5, -5}, max: Coord{-1, -1}},
			overlaps: false,
		},
		{
			bounds:   Bounds{layout: XYZ, min: Coord{-100, -100, -100}, max: Coord{100, 100, 100}},
			other:    Bounds{layout: XYZ, min: Coord{-10, 0, 0}, max: Coord{-5, 10, 10}},
			overlaps: true,
		},
		{
			bounds:   Bounds{layout: XYZ, min: Coord{0, 0, 0}, max: Coord{100, 100, 100}},
			other:    Bounds{layout: XYZ, min: Coord{5, 5, -10}, max: Coord{10, 10, -5}},
			overlaps: false,
		},
		{
			bounds:   Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			other:    Bounds{layout: XY, min: Coord{-10, -10}, max: Coord{-0.000000000000000000000000000001, 0}},
			overlaps: false,
		},
	} {
		copy := Bounds{layout: testData.bounds.layout, min: testData.bounds.min, max: testData.bounds.max}
		for j := 0; j < 10; j++ {
			// do multiple checks to verify no obvious side effects are caused
			assert.Equal(t, testData.overlaps, copy.Overlaps(testData.bounds.layout, &testData.other))
			assert.Equal(t, testData.bounds, copy)
		}
	}
}

func TestBoundsOverlapsPoint(t *testing.T) {
	for _, testData := range []struct {
		bounds   Bounds
		point    Coord
		overlaps bool
	}{
		{
			bounds:   Bounds{layout: XY, min: Coord{0, 0}, max: Coord{0, 0}},
			point:    Coord{-10, 0},
			overlaps: false,
		},
		{
			bounds:   Bounds{layout: XY, min: Coord{-100, -100}, max: Coord{100, 100}},
			point:    Coord{-10, 0},
			overlaps: true,
		},
		{
			bounds:   Bounds{layout: XYZ, min: Coord{-100, -100, -100}, max: Coord{100, 100, 100}},
			point:    Coord{-5, 10, 10},
			overlaps: true,
		},
		{
			bounds:   Bounds{layout: XYZ, min: Coord{0, 0, 0}, max: Coord{100, 100, 100}},
			point:    Coord{5, 5, -10},
			overlaps: false,
		},
		{
			bounds:   Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}},
			point:    Coord{-0.000000000000000000000000000001, 0},
			overlaps: false,
		},
	} {
		copy := Bounds{layout: testData.bounds.layout, min: testData.bounds.min, max: testData.bounds.max}
		for j := 0; j < 10; j++ {
			// do multiple checks to verify no obvious side effects are caused
			assert.Equal(t, testData.overlaps, copy.OverlapsPoint(testData.bounds.layout, testData.point))
			assert.Equal(t, testData.bounds, copy)
		}
	}
}

func TestBoundsPolygon(t *testing.T) {
	for _, tc := range []struct {
		b    *Bounds
		want *Polygon
	}{
		{
			b:    NewBounds(NoLayout),
			want: NewPolygon(XY),
		},
		{
			b:    NewBounds(XY).Set(0, 0, 1, 1),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:    NewBounds(XYZ).Set(0, 0, 0, 1, 1, 1),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:    NewBounds(XYM).Set(0, 0, 0, 1, 1, 1),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:    NewBounds(XYZM).Set(0, 0, 0, 0, 1, 1, 1, 1),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}}),
		},
		{
			b:    NewBounds(XY).Set(1, 2, 3, 4),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{1, 2}, {1, 4}, {3, 4}, {3, 2}, {1, 2}}}),
		},
		{
			b:    NewBounds(XYZ).Set(1, 2, 3, 4, 5, 6),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{1, 2}, {1, 5}, {4, 5}, {4, 2}, {1, 2}}}),
		},
	} {
		assert.Equal(t, tc.want, tc.b.Polygon())
	}
}

func TestBoundsSet(t *testing.T) {
	bounds := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}}
	bounds.Set(0, 0, 20, 20)
	expected := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{20, 20}}
	assert.Equal(t, expected, bounds)
	assert.Panics(t, func() {
		bounds.Set(2, 2, 2, 2, 2)
	})
}

func TestBoundsSetCoords(t *testing.T) {
	bounds := &Bounds{layout: XY, min: Coord{0, 0}, max: Coord{10, 10}}
	bounds.SetCoords(Coord{0, 0}, Coord{20, 20})
	expected := Bounds{layout: XY, min: Coord{0, 0}, max: Coord{20, 20}}
	assert.Equal(t, expected, *bounds)

	bounds = NewBounds(XY)
	bounds.SetCoords(Coord{0, 0}, Coord{20, 20})
	assert.Equal(t, expected, *bounds)

	bounds = NewBounds(XY)
	bounds.SetCoords(Coord{20, 0}, Coord{0, 20}) // set coords should ensure valid min / max
	assert.Equal(t, expected, *bounds)
}
