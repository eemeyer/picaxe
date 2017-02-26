package imageops

import (
	"image"
	"image/color"
	"math"

	"github.com/oliamb/cutter"
)

// Trim trims borders from images. For example:
//
//                  mmmmmm           mmXmmm
//   mmmm           mmmmmm           mmmmmX
//   mAAm -> AA     mmAAmm -> AA     mmAAmm -> (same)
//   mmmm           mmmmmm           mmmmmm
//                  mmmmmm           mmXmmm
//
// The edge of the image is considered a trimmable border iff it is
// contiguous with respect to color distance. A color is contiguious iff
// the distance to the adjacent pixel's color is less than or equal to
// the fuzz factor. (With a fuzz factor of 0.0, all colors are distinct.)
// Furthermore, the border must extend around the entire rectangular edge
// of the image. The algorithm trims the outer edge concentrically until
// a non-consecutive edge is found.
//
func Trim(img image.Image, fuzzFactor float64) image.Image {
	bounds := img.Bounds()
	if bounds.Empty() {
		return img
	}

	cornerColor := img.At(0, 0)

	centerX := bounds.Min.X + bounds.Dx()/2
	centerY := bounds.Min.Y + bounds.Dy()/2

	// Very stupid algorithm to concentrically determine a trimmable border
	var xdepth, ydepth int
OuterX:
	for xdepth = bounds.Min.X; xdepth < centerX; xdepth++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if !(colorDistance(img.At(xdepth, y), cornerColor) <= fuzzFactor &&
				colorDistance(img.At(bounds.Max.X-xdepth-1, y), cornerColor) <= fuzzFactor) {
				break OuterX
			}
		}
	}
OuterY:
	for ydepth = bounds.Min.Y; ydepth < centerY; ydepth++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
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
