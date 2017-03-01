package imageops

import (
	"image"

	"github.com/nfnt/resize"
)

func Scale(img image.Image, size image.Point) image.Image {
	w, h, rect := size.X, size.Y, img.Bounds()

	if w == rect.Dx() && h == rect.Dy() {
		return img
	}
	if w <= 0 || h <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	return resize.Resize(uint(w), uint(h), img, resize.Lanczos2)
}
