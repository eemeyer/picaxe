package picaxe

import (
	"image"

	"github.com/oliamb/cutter"
)

type CropMode int

const (
	CropModeNone   CropMode = 0
	CropModeCenter          = 1
)

func Crop(img image.Image, mode CropMode, width int, height int) image.Image {
	if mode != CropModeNone {
		bounds := img.Bounds()
		var w, h int
		if bounds.Dx() > width {
			w = width
		} else {
			w = bounds.Dx()
		}
		if bounds.Dy() > height {
			h = height
		} else {
			h = bounds.Dy()
		}
		switch mode {
		case CropModeCenter:
			img, _ = cutter.Crop(img, cutter.Config{
				Mode:   cutter.Centered,
				Width:  w,
				Height: h,
			})
		}
	}
	return img
}
