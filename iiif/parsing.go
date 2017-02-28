package iiif

import (
	"fmt"
	"image"
	"regexp"
	"strconv"

	"github.com/t11e/picaxe/imageops"
)

func parsePercentageCoords(s string) (*imageops.RelativeRegion, error) {
	parts := percentageCoordsRegexp.FindStringSubmatch(s)
	if len(parts) != 5 {
		return nil, InvalidSpec{
			Message: fmt.Sprintf("Not a valid set of coordinates: %s", s),
		}
	}

	x, err := parsePercentage(parts[1])
	if err != nil {
		return nil, err
	}

	y, err := parsePercentage(parts[2])
	if err != nil {
		return nil, err
	}

	w, err := parsePercentage(parts[3])
	if err != nil {
		return nil, err
	}

	h, err := parsePercentage(parts[4])
	if err != nil {
		return nil, err
	}

	return imageops.NewRelativeRegion(x, y, w, h)
}

func parseRectangle(s string) (*image.Rectangle, error) {
	parts := pixelRectangleRegexp.FindStringSubmatch(s)
	if len(parts) != 5 {
		return nil, InvalidSpec{
			Message: fmt.Sprintf("Not a valid set of coordinates: %s", s),
		}
	}

	x, err := parsePixelComponent(parts[1])
	if err != nil {
		return nil, err
	}

	y, err := parsePixelComponent(parts[2])
	if err != nil {
		return nil, err
	}

	w, err := parsePixelComponent(parts[3])
	if err != nil {
		return nil, err
	}

	h, err := parsePixelComponent(parts[4])
	if err != nil {
		return nil, err
	}

	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}

	rect := image.Rect(x, y, x+w, y+h).Canon()
	return &rect, nil
}

func parseWidthHeight(s string) (width *int, height *int, bestFit bool, err error) {
	parts := pixelWHRegexp.FindStringSubmatch(s)
	if len(parts) != 4 {
		err = InvalidSpec{
			Message: fmt.Sprintf("Not a valid width/height: %s", s),
		}
		return
	}

	if parts[2] != "" {
		var w int
		w, err = parsePixelComponent(parts[2])
		if err != nil {
			return
		}
		width = &w
	}

	if parts[3] != "" {
		var h int
		h, err = parsePixelComponent(parts[3])
		if err != nil {
			return
		}
		height = &h
	}

	bestFit = width != nil && height != nil && parts[1] == "!"
	return
}

var (
	percentageCoordsRegexp = regexp.MustCompile(
		`^(-?[\d]+(?:\.[\d]+)?),(-?[\d]+(?:\.[\d]+)?),(-?[\d]+(?:\.[\d]+)?),(-?[\d]+(?:\.[\d]+)?)$`)
	pixelRectangleRegexp = regexp.MustCompile(
		`^(-?[\d]+),(-?[\d]+),(-?[\d]+),(-?[\d]+)$`)
	pixelWHRegexp = regexp.MustCompile(
		`^(!)?([\d]+)?,([\d]+)?$`)
)

func parsePixelComponent(s string) (int, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func parsePercentage(s string) (float64, error) {
	pcnt, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if pcnt < 0 || pcnt > 100 {
		return 0, InvalidSpec{
			Message: fmt.Sprintf("Invalid percentage number: %s", s),
		}
	}
	return pcnt / 100, nil
}
