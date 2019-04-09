package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeometryCollectionBounds(t *testing.T) {
	for _, tc := range []struct {
		geoms []T
		want  *Bounds
	}{
		{
			want: NewBounds(NoLayout),
		},
		{
			geoms: []T{
				NewPoint(XY),
			},
			want: NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
		},
		{
			geoms: []T{
				NewPoint(XY),
			},
			want: NewBounds(XY).SetCoords(Coord{0, 0}, Coord{0, 0}),
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XY).MustSetCoords(Coord{3, 4}),
			},
			want: NewBounds(XY).SetCoords(Coord{1, 2}, Coord{3, 4}),
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XYZ).MustSetCoords(Coord{3, 4, 5}),
			},
			want: NewBounds(XYZ).SetCoords(Coord{1, 2, 5}, Coord{3, 4, 5}),
		},
		{
			geoms: []T{
				NewPoint(XY).MustSetCoords(Coord{1, 2}),
				NewPoint(XYM).MustSetCoords(Coord{3, 4, 5}),
			},
			want: NewBounds(XYM).SetCoords(Coord{1, 2, 5}, Coord{3, 4, 5}),
		},
		{
			geoms: []T{
				NewPoint(XYZ).MustSetCoords(Coord{1, 2, 3}),
				NewPoint(XYM).MustSetCoords(Coord{4, 5, 6}),
			},
			want: NewBounds(XYZM).SetCoords(Coord{1, 2, 3, 6}, Coord{4, 5, 3, 6}),
		},
	} {
		assert.Equal(t, tc.want, NewGeometryCollection().MustPush(tc.geoms...).Bounds())
	}
}

func TestGeometryCollectionLayout(t *testing.T) {
	for _, tc := range []struct {
		geoms []T
		want  Layout
	}{
		{
			want: NoLayout,
		},
		{
			geoms: []T{
				NewPoint(XY),
			},
			want: XY,
		},
		{
			geoms: []T{
				NewPoint(XYZ),
			},
			want: XYZ,
		},
		{
			geoms: []T{
				NewPoint(XYM),
			},
			want: XYM,
		},
		{
			geoms: []T{
				NewPoint(XYZM),
			},
			want: XYZM,
		},
		{
			geoms: []T{
				NewPoint(XY),
				NewPoint(XYZ),
			},
			want: XYZ,
		},
		{
			geoms: []T{
				NewPoint(XY),
				NewPoint(XYM),
			},
			want: XYM,
		},
		{
			geoms: []T{
				NewPoint(XY),
				NewPoint(XYZ),
				NewPoint(XYM),
			},
			want: XYZM,
		},
	} {
		assert.Equal(t, tc.want, NewGeometryCollection().MustPush(tc.geoms...).Layout())
	}
}

func TestGeometryCollectionPush(t *testing.T) {
	for _, tc := range []struct {
		srid    int
		geoms   []T
		g       T
		wantErr error
	}{
		{
			g: NewPoint(XY),
		},
		{
			g: NewPoint(XY).SetSRID(4326),
		},
		{
			srid: 4326,
			g:    NewPoint(XY).SetSRID(4326),
		},
		{
			geoms: []T{
				NewPoint(XY).SetSRID(4326),
			},
			g: NewPoint(XY).SetSRID(3857),
		},
		{
			srid: 4326,
			g:    NewPoint(XY).SetSRID(3857),
		},
	} {
		// FIXME this test is obviously wrong as wantErr is nil in every case
		assert.Equal(t, tc.wantErr, NewGeometryCollection().SetSRID(tc.srid).MustPush(tc.geoms...).Push(tc.g))
	}
}
