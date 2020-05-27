// Package wkt implements Well Known Text encoding and decoding.
package wkt

import (
	"errors"

	"github.com/twpayne/go-geom"
)

const (
	tPoint              = "POINT "
	tMultiPoint         = "MULTIPOINT "
	tLineString         = "LINESTRING "
	tMultiLineString    = "MULTILINESTRING "
	tPolygon            = "POLYGON "
	tMultiPolygon       = "MULTIPOLYGON "
	tGeometryCollection = "GEOMETRYCOLLECTION "
	tZ                  = "Z "
	tM                  = "M "
	tZm                 = "ZM "
	tEmpty              = "EMPTY"
)

// ErrBraceMismatch is returned when braces do not match.
var ErrBraceMismatch = errors.New("wkt: brace mismatch")

// Marshal translates a geometry to the corresponding WKT.
func Marshal(g geom.T) (string, error) {
	return encode(g, defaultMaxDecimalDigits)
}

// Marshal translates a geometry to the corresponding WKT.
// This variant only prints up to the maximum decimal digits for each coord.
func MarshalWithMaxDecimalDigits(g geom.T, maxDecimalDigits int) (string, error) {
	return encode(g, maxDecimalDigits)
}

// Unmarshal translates a WKT to the corresponding geometry.
func Unmarshal(wkt string) (geom.T, error) {
	return decode(wkt)
}
