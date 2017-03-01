package imageops

import (
	"errors"
	"image"
	"math"
)

type RelativeRegion struct {
	X, Y, W, H float64
}

func NewRelativeRegion(x, y, w, h float64) (*RelativeRegion, error) {
	if !(x >= 0 && x <= 1 && y >= 0 && y <= 1 && w >= 0 && w <= 1 && h >= 0 && h <= 1) {
		return nil, errors.New("Invalid coordinates")
	}
	return &RelativeRegion{X: x, Y: y, W: w, H: h}, nil
}

func (r RelativeRegion) IsEmpty() bool {
	return r.W <= 0 || r.H <= 0
}

func (r RelativeRegion) ToRectangle(rect image.Rectangle) image.Rectangle {
	w, h := float64(rect.Dx()), float64(rect.Dy())
	return image.Rect(
		round(r.X*w),
		round(r.Y*h),
		round(r.X*w+r.W*w),
		round(r.Y*h+r.H*h)).Canon().Intersect(rect)
}

type RelativePoint struct {
	X, Y float64
}

func NewRelativePoint(x, y float64) (*RelativePoint, error) {
	if x < 0 || x > 1.0 || y < 0 || y > 1.0 {
		return nil, errors.New("Invalid coordinates")
	}
	return &RelativePoint{X: x, Y: y}, nil
}

func (r RelativePoint) ToPoint(rect image.Rectangle) image.Point {
	w, h := float64(rect.Dx()), float64(rect.Dy())
	return image.Pt(round(r.X*w), round(r.Y*h))
}

// FitDimensions scales (down, or up if necessary) a set of dimensions to
// fit with w, h, preserving the aspect ratio of the input.
//
// * If w and h are nil, the input size is returned.
// * If w is nil, then the dimensions are scaled to fit within the height.
// * If h is nil, then the dimensions are scaled to fit within the width.
//
func FitDimensions(size image.Point, w, h *int) image.Point {
	if w == nil && h == nil {
		return size
	}

	if size.Y <= 0 || size.X <= 0 || (w != nil && *w <= 0) || (h != nil && *h <= 0) {
		return image.Pt(0, 0)
	}

	sw, sh := float64(size.X), float64(size.Y)

	if w != nil && h != nil {
		tw, th := float64(*w), float64(*h)
		scale := min(tw/sw, th/sh)
		return image.Pt(round(sw*scale), round(sh*scale))
	}

	aspect := sw / sh
	if aspect <= 0 {
		return image.Pt(0, 0)
	}
	if w != nil {
		return image.Pt(*w, round(float64(*w)/aspect))
	}
	return image.Pt(round(float64(*h)*aspect), *h)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func round(f float64) int {
	return int(math.Floor(f + .5))
}
