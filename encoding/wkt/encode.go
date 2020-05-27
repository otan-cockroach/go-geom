package wkt

import (
	"strconv"
	"strings"

	"github.com/twpayne/go-geom"
)

// defaultMaxDecimalDigits ensure printing as many decimal digits as possible.
const defaultMaxDecimalDigits = -1

// encode translates a geometry to the corresponding WKT.
func encode(g geom.T, maxDecimalDigits int) (string, error) {
	sb := &strings.Builder{}
	if err := write(sb, g, maxDecimalDigits); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func write(sb *strings.Builder, g geom.T, maxDecimalDigits int) error {
	typeString := ""
	switch g := g.(type) {
	case *geom.Point:
		typeString = tPoint
	case *geom.LineString, *geom.LinearRing:
		typeString = tLineString
	case *geom.Polygon:
		typeString = tPolygon
	case *geom.MultiPoint:
		typeString = tMultiPoint
	case *geom.MultiLineString:
		typeString = tMultiLineString
	case *geom.MultiPolygon:
		typeString = tMultiPolygon
	case *geom.GeometryCollection:
		typeString = tGeometryCollection
	default:
		return geom.ErrUnsupportedType{Value: g}
	}
	layout := g.Layout()
	switch layout {
	case geom.NoLayout:
		// Special case for empty GeometryCollections
		if g, ok := g.(*geom.GeometryCollection); !ok || !g.Empty() {
			return geom.ErrUnsupportedLayout(layout)
		}
	case geom.XY:
	case geom.XYZ:
		typeString += tZ
	case geom.XYM:
		typeString += tM
	case geom.XYZM:
		typeString += tZm
	default:
		return geom.ErrUnsupportedLayout(layout)
	}
	if _, err := sb.WriteString(typeString); err != nil {
		return err
	}
	switch g := g.(type) {
	case *geom.Point:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords0(sb, g.FlatCoords(), layout.Stride(), maxDecimalDigits)
	case *geom.LineString:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords1(sb, g.FlatCoords(), layout.Stride(), maxDecimalDigits)
	case *geom.LinearRing:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords1(sb, g.FlatCoords(), layout.Stride(), maxDecimalDigits)
	case *geom.Polygon:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords2(sb, g.FlatCoords(), 0, g.Ends(), layout.Stride(), maxDecimalDigits)
	case *geom.MultiPoint:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords1(sb, g.FlatCoords(), layout.Stride(), maxDecimalDigits)
	case *geom.MultiLineString:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords2(sb, g.FlatCoords(), 0, g.Ends(), layout.Stride(), maxDecimalDigits)
	case *geom.MultiPolygon:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		return writeFlatCoords3(sb, g.FlatCoords(), g.Endss(), layout.Stride(), maxDecimalDigits)
	case *geom.GeometryCollection:
		if g.Empty() {
			return writeEMPTY(sb)
		}
		if _, err := sb.WriteRune('('); err != nil {
			return err
		}
		for i, g := range g.Geoms() {
			if i != 0 {
				if _, err := sb.WriteString(", "); err != nil {
					return err
				}
			}
			if err := write(sb, g, maxDecimalDigits); err != nil {
				return err
			}
		}
		_, err := sb.WriteRune(')')
		return err
	}
	return nil
}

func writeCoord(sb *strings.Builder, coord []float64, maxDecimalDigits int) error {
	for i, x := range coord {
		if i != 0 {
			if _, err := sb.WriteRune(' '); err != nil {
				return err
			}
		}
		coordStr := strconv.FormatFloat(x, 'f', maxDecimalDigits, 64)
		if maxDecimalDigits != -1 {
			coordStr = strings.TrimRight(strings.TrimRight(coordStr, "0"), ".")
		}
		if _, err := sb.WriteString(coordStr); err != nil {
			return err
		}
	}
	return nil
}

//nolint:interfacer
func writeEMPTY(sb *strings.Builder) error {
	_, err := sb.WriteString(tEmpty)
	return err
}

func writeFlatCoords0(
	sb *strings.Builder, flatCoords []float64, stride, maxDecimalDigits int,
) error {
	if _, err := sb.WriteRune('('); err != nil {
		return err
	}
	if err := writeCoord(sb, flatCoords[:stride], maxDecimalDigits); err != nil {
		return err
	}
	_, err := sb.WriteRune(')')
	return err
}

func writeFlatCoords1(
	sb *strings.Builder, flatCoords []float64, stride, maxDecimalDigits int,
) error {
	if _, err := sb.WriteRune('('); err != nil {
		return err
	}
	for i, n := 0, len(flatCoords); i < n; i += stride {
		if i != 0 {
			if _, err := sb.WriteString(", "); err != nil {
				return err
			}
		}
		if err := writeCoord(sb, flatCoords[i:i+stride], maxDecimalDigits); err != nil {
			return err
		}
	}
	_, err := sb.WriteRune(')')
	return err
}

func writeFlatCoords2(
	sb *strings.Builder, flatCoords []float64, start int, ends []int, stride, maxDecimalDigits int,
) error {
	if _, err := sb.WriteRune('('); err != nil {
		return err
	}
	for i, end := range ends {
		if i != 0 {
			if _, err := sb.WriteString(", "); err != nil {
				return err
			}
		}
		if err := writeFlatCoords1(sb, flatCoords[start:end], stride, maxDecimalDigits); err != nil {
			return err
		}
		start = end
	}
	_, err := sb.WriteRune(')')
	return err
}

func writeFlatCoords3(
	sb *strings.Builder, flatCoords []float64, endss [][]int, stride, maxDecimalDigits int,
) error {
	if _, err := sb.WriteRune('('); err != nil {
		return err
	}
	start := 0
	for i, ends := range endss {
		if i != 0 {
			if _, err := sb.WriteString(", "); err != nil {
				return err
			}
		}
		if err := writeFlatCoords2(sb, flatCoords, start, ends, stride, maxDecimalDigits); err != nil {
			return err
		}
		start = ends[len(ends)-1]
	}
	_, err := sb.WriteRune(')')
	return err
}
