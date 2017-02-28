package iiif

import (
	"fmt"
	"image"
	"math"
	"net/url"
	"regexp"
	"strings"

	"github.com/t11e/picaxe/imageops"
)

type InvalidSpec struct {
	Message string
}

func (e InvalidSpec) Error() string {
	return e.Message
}

const (
	RegionStringFull   = "full"
	RegionStringSquare = "square"
	RegionStringPct    = "pct:"
)

type RegionKind int

const (
	RegionKindFull RegionKind = iota
	RegionKindSquare
	RegionKindAbsolute
	RegionKindRelative
)

type Region struct {
	Kind     RegionKind
	Absolute *image.Rectangle
	Relative *imageops.RelativeRegion
}

const (
	SizeStringFull = "full"
	SizeStringMax  = "max"
	SizeStringPct  = "pct:"
)

type SizeKind int

const (
	SizeKindFull SizeKind = iota
	SizeKindMax
	SizeKindAbsolute
	SizeKindRelative
)

type Size struct {
	Kind       SizeKind
	AbsWidth   *int
	AbsHeight  *int
	AbsBestFit bool
	Relative   *float64
}

func (size Size) CalculateDimensions(
	rect image.Rectangle, limit image.Point) (image.Point, error) {
	s := rect.Size()
	switch size.Kind {
	case SizeKindFull:
		return checkDimensions(limit, s)
	case SizeKindMax:
		// TODO: Do aspect
		if s.X > limit.X && s.Y > limit.Y {
			if s.X > s.Y {
				w := limit.X
				return checkDimensions(limit, computeDimensions(rect, &w, nil))
			}
			h := limit.Y
			return checkDimensions(limit, computeDimensions(rect, nil, &h))
		} else if s.X > limit.X {
			w := limit.X
			return checkDimensions(limit, computeDimensions(rect, &w, nil))
		} else if s.Y > limit.Y {
			h := limit.Y
			return checkDimensions(limit, computeDimensions(rect, nil, &h))
		}
		return checkDimensions(limit, s)
	case SizeKindAbsolute:
		// TODO: AbsBestFit
		return checkDimensions(limit, computeDimensions(rect, size.AbsWidth, size.AbsHeight))
	case SizeKindRelative:
		w := round(float64(s.X) * *size.Relative)
		h := round(float64(s.Y) * *size.Relative)
		return checkDimensions(limit, computeDimensions(rect, &w, &h))
	}
	panic("Invalid size specification")
}

func checkDimensions(limit image.Point, p image.Point) (image.Point, error) {
	if p.X > limit.X || p.Y > limit.Y {
		return image.Point{}, fmt.Errorf("(%d, %d) exceeds maximum allowed dimensions (%d, %d)",
			p.X, p.Y, limit.X, limit.Y)
	}
	return p, nil
}

func computeDimensions(rect image.Rectangle, w, h *int) image.Point {
	if rect.Dy() <= 0 || rect.Dx() <= 0 {
		return image.Pt(0, 0)
	}

	aspect := float64(rect.Dy()) / float64(rect.Dx())
	if aspect <= 0 {
		return image.Pt(0, 0)
	}

	if w != nil {
		if h != nil {
			return image.Pt(*w, *h)
		}
		return image.Pt(*w, int(math.Floor(float64(*w)*aspect)))
	} else if h != nil {
		return image.Pt(int(math.Floor(float64(*h)/aspect)), *h)
	}
	return rect.Size()
}

type Format string

const (
	FormatDefault Format = ""
	FormatJPEG           = "jpg"
	FormatPNG            = "png"
	FormatGIF            = "gif"
)

type Request struct {
	Identifier string
	Region     Region
	Size       Size
	Format     Format
}

var specRegexp = regexp.MustCompile(`([^/]+)/([^/]+)/([^/]+)/([^/]+)/([^/]+)\.(.+)$`)

func ParseSpec(spec string) (*Request, error) {
	parts := specRegexp.FindStringSubmatch(spec)
	if len(parts) != 7 {
		return nil, InvalidSpec{
			Message: fmt.Sprintf("not a valid spec: %q", spec),
		}
	}

	var req Request

	if id, err := url.QueryUnescape(parts[1]); err == nil {
		req.Identifier = id
	} else {
		return nil, err
	}

	if err := parseRegion(parts[2], &req.Region); err != nil {
		return nil, err
	}

	if err := parseSize(parts[3], &req.Size); err != nil {
		return nil, err
	}

	if rotation := parts[4]; rotation != "" {
		switch rotation {
		case "0":
			// OK
		default:
			return nil, InvalidSpec{
				Message: fmt.Sprintf("unsupported rotation %q", rotation),
			}
		}
	}

	if quality := parts[5]; quality != "" {
		switch quality {
		case "color", "default":
			// OK
		default:
			return nil, InvalidSpec{
				Message: fmt.Sprintf("unsupported quality %q", quality),
			}
		}
	}

	if format := parts[6]; format != "" {
		name, ok := formatNameMap[format]
		if !ok {
			return nil, InvalidSpec{
				Message: fmt.Sprintf("unsupported format %q", format),
			}
		}
		req.Format = name
	} else {
		req.Format = FormatDefault
	}

	return &req, nil
}

func parseRegion(regionValue string, region *Region) error {
	switch regionValue {
	case RegionStringFull, "":
		region.Kind = RegionKindFull
		return nil
	case RegionStringSquare:
		region.Kind = RegionKindSquare
		return nil
	}

	if strings.HasPrefix(regionValue, RegionStringPct) {
		var err error
		region.Kind = RegionKindRelative
		region.Relative, err = parsePercentageCoords(regionValue[len(RegionStringPct):])
		if err != nil {
			return err
		}
		return nil
	}

	var err error
	region.Kind = RegionKindAbsolute
	region.Absolute, err = parseRectangle(regionValue)
	return err
}

func parseSize(sizeValue string, size *Size) error {
	switch sizeValue {
	case SizeStringFull, "":
		size.Kind = SizeKindFull
		return nil
	case SizeStringMax:
		size.Kind = SizeKindMax
		return nil
	}

	if strings.HasPrefix(sizeValue, SizeStringPct) {
		pcnt, err := parsePercentage(sizeValue[len(SizeStringPct):])
		if err != nil {
			return err
		}

		size.Kind = SizeKindRelative
		size.Relative = &pcnt
		return nil
	}

	var err error
	size.Kind = SizeKindAbsolute
	size.AbsWidth, size.AbsHeight, size.AbsBestFit, err = parseWidthHeight(sizeValue)
	return err
}

func round(f float64) int {
	return int(math.Floor(f + .5))
}

var formatNameMap map[string]Format

func init() {
	formats := []Format{FormatJPEG, FormatPNG, FormatGIF}
	formatNameMap = make(map[string]Format, len(formats))
	for _, n := range formats {
		formatNameMap[string(n)] = n
	}
}
