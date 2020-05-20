package geom

import "math"

// A Point represents a single point.
type Point struct {
	geom0
}

// emptyPointFloat64 signifies whether the Point actually POINT EMPTY.
// This complies with Requirement 152 of http://www.geopackage.org/spec/.
var emptyPointFloat64 = math.Float64frombits(0x7FF8000000000000)

// NewPoint allocates a new Point with layout l and all values zero.
func NewPoint(l Layout) *Point {
	return NewPointFlat(l, make([]float64, l.Stride()))
}

// NewPointEmpty allocates a new Point representing an empty Point.
// This follows Requirement 152 of http://www.geopackage.org/spec/:
//   If the geometry is a Point, it SHALL be encoded with each coordinate
//   value set to an IEEE-754 quiet NaN value. GeoPackages SHALL use big
//   endian 0x7ff8000000000000 or little endian 0x000000000000f87f as the
//   binary encoding of the NaN values. (This is because Well-Known Binary
//   as defined in OGC 06-103r4 [9] does not provide a standardized encoding
//   for an empty point set, i.e., 'Point Empty' in Well-Known Text.)
func NewPointEmpty(l Layout) *Point {
	coords := make([]float64, l.Stride())
	for i := 0; i < l.Stride(); i++ {
		coords[i] = emptyPointFloat64
	}
	return NewPointFlat(l, coords)
}

// NewPointFlat allocates a new Point with layout l and flat coordinates flatCoords.
func NewPointFlat(l Layout, flatCoords []float64) *Point {
	g := new(Point)
	g.layout = l
	g.stride = l.Stride()
	g.flatCoords = flatCoords
	return g
}

// Area returns g's area, i.e. zero.
func (g *Point) Area() float64 {
	return 0
}

// Empty returns whether g is empty.
func (g *Point) Empty() bool {
	for i := 0; i < g.Stride(); i++ {
		// We cannot compare NaN vs NaN directly.
		if !math.IsNaN(g.flatCoords[i]) {
			return false
		}
	}
	return true
}

// Clone returns a copy of g that does not alias g.
func (g *Point) Clone() *Point {
	return deriveClonePoint(g)
}

// Length returns the length of g, i.e. zero.
func (g *Point) Length() float64 {
	return 0
}

// MustSetCoords is like SetCoords but panics on any error.
func (g *Point) MustSetCoords(coords Coord) *Point {
	Must(g.SetCoords(coords))
	return g
}

// SetCoords sets the coordinates of g.
func (g *Point) SetCoords(coords Coord) (*Point, error) {
	if err := g.setCoords(coords); err != nil {
		return nil, err
	}
	return g, nil
}

// SetSRID sets the SRID of g.
func (g *Point) SetSRID(srid int) *Point {
	g.srid = srid
	return g
}

// Swap swaps the values of g and g2.
func (g *Point) Swap(g2 *Point) {
	*g, *g2 = *g2, *g
}

// X returns g's X-coordinate.
func (g *Point) X() float64 {
	return g.flatCoords[0]
}

// Y returns g's Y-coordinate.
func (g *Point) Y() float64 {
	return g.flatCoords[1]
}

// Z returns g's Z-coordinate, or zero if g has no Z-coordinate.
func (g *Point) Z() float64 {
	zIndex := g.layout.ZIndex()
	if zIndex == -1 {
		return 0
	}
	return g.flatCoords[zIndex]
}

// M returns g's M-coordinate, or zero if g has no M-coordinate.
func (g *Point) M() float64 {
	mIndex := g.layout.MIndex()
	if mIndex == -1 {
		return 0
	}
	return g.flatCoords[mIndex]
}
