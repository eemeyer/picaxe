package main

import (
	"image"

	"github.com/disintegration/imaging"
)

func NormalizeOrientation(img image.Image, exifOrientation string) image.Image {
	switch exifOrientation {
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.FlipH(img)
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img) // 90 degrees clockwise
	case "7":
		return imaging.Rotate270(imaging.FlipH(img))
	case "8":
		return imaging.Rotate90(img) // 270 degrees clockwise
	default:
		return img
	}
}
