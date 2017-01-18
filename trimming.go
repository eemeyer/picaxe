package main

import (
	"github.com/oliamb/cutter"
	"image"
	"image/color"
	"math"
)

type (
	TrimMode int
)

const (
	TrimModeNone  = 0
	TrimModeFuzzy = 1
)

func Trim(img image.Image, mode TrimMode, fuzzFactor float64) image.Image {
	if mode != TrimModeNone {
		bounds := img.Bounds()
		if !bounds.Empty() {
			cornerColor := img.At(0, 0)

			centerX := bounds.Min.X + bounds.Dx()/2
			centerY := bounds.Min.Y + bounds.Dy()/2

			var xdepth, ydepth int
		OuterX:
			for xdepth = bounds.Min.X; xdepth < centerX; xdepth += 1 {
				for y := bounds.Min.Y; y < bounds.Max.Y; y += 1 {
					if !(colorDistance(img.At(xdepth, y), cornerColor) <= fuzzFactor &&
						colorDistance(img.At(bounds.Max.X-xdepth-1, y), cornerColor) <= fuzzFactor) {
						break OuterX
					}
				}
			}
		OuterY:
			for ydepth = bounds.Min.Y; ydepth < centerY; ydepth += 1 {
				for x := bounds.Min.X; x < bounds.Max.X; x += 1 {
					if !(colorDistance(img.At(x, ydepth), cornerColor) <= fuzzFactor &&
						colorDistance(img.At(x, bounds.Max.Y-ydepth-1), cornerColor) <= fuzzFactor) {
						break OuterY
					}
				}
			}

			if xdepth > 0 || ydepth > 0 {
				img, _ = cutter.Crop(img, cutter.Config{
					Width:  bounds.Dx() - (xdepth * 2),
					Height: bounds.Dy() - (ydepth * 2),
					Anchor: image.Point{xdepth, ydepth},
					Mode:   cutter.TopLeft,
				})
			}
		}
	}
	return img
}

func colorDistance(a color.Color, b color.Color) float64 {
	r1, g1, b1, _ := a.RGBA()
	r2, g2, b2, _ := b.RGBA()
	distance := math.Sqrt(math.Pow(math.Abs(float64(r2)-float64(r1)), 2) +
		math.Pow(math.Abs(float64(g2)-float64(g1)), 2) +
		math.Pow(math.Abs(float64(b2)-float64(b1)), 2))
	return distance / math.Sqrt(math.Pow(65535, 2)*3)
}
