package imageops

import (
	"image"

	"github.com/oliamb/cutter"
)

// CropSquare crops an image from the center to the shortest dimension that
// will result in a square image.
func CropSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()

	var size int
	if dx > dy {
		size = dy
	} else {
		size = dx
	}

	img, err := cutter.Crop(img, cutter.Config{
		Mode:   cutter.Centered,
		Width:  size,
		Height: size,
	})
	if err != nil {
		panic(err)
	}
	return img
}

// CropRelative crops to relative coordinates.
func CropRelative(img image.Image, region RelativeRegion) image.Image {
	return CropRect(img, region.ToRectangle(img.Bounds()))
}

// CropRect crops an image to a rectangular region.
func CropRect(img image.Image, rect image.Rectangle) image.Image {
	rect = rect.Intersect(img.Bounds())
	img, err := cutter.Crop(img, cutter.Config{
		Mode:   cutter.TopLeft,
		Anchor: rect.Min,
		Width:  rect.Dx(),
		Height: rect.Dy(),
	})
	if err != nil {
		panic(err)
	}
	return img
}
