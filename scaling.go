package main

import (
	"image"
	"log"
	"math"

	"github.com/nfnt/resize"
)

type ScaleMode int

const (
	// ScaleModeNone is no scaling; target dimensions will be ignored.
	ScaleModeNone ScaleMode = 0

	// ScaleModeUp scales up to target dimensions if smaller.
	ScaleModeUp = 1

	// ScaleModeDown scales down to fit within the target dimensions; if already
	// smaller, no scaling is performed.
	ScaleModeDown = 2

	// ScaleModeCover covers target dimensions, even if this will exceed the
	// dimensions. If the input image is smaller or bigger than the target
	// dimensions, it will be scaled up or down so that it covers or exceeds the
	// target dimensions. Combined with crop, this ensures a cropped image is
	// always at least the crop size, which is ideal for thumbnails.
	ScaleModeCover = 3
)

func Scale(img image.Image, mode ScaleMode, targetWidth, targetHeight int) image.Image {
	if mode != ScaleModeNone {
		rect := img.Bounds()
		var w, h int = computeDimensions(rect, mode, targetWidth, targetHeight)
		log.Print("scale to: ", w, h)
		if w != rect.Dx() || h != rect.Dy() {
			img = resize.Resize(uint(w), uint(h), img, resize.Lanczos2)
		}
	}
	return img
}

func min(a float64, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func computeDimensions(
	rect image.Rectangle,
	mode ScaleMode,
	targetWidth, targetHeight int) (newWidth, newHeight int) {
	targetW := float64(targetWidth)
	targetH := float64(targetHeight)
	width := float64(rect.Dx())
	height := float64(rect.Dy())
	aspect := float64(height) / float64(width)

	var w, h float64
	switch mode {
	case ScaleModeUp:
		if height < targetW || width < targetH {
			if height > width {
				w = targetW
				h = math.Floor(w / aspect)
			} else {
				w = math.Floor(targetH / aspect)
				h = targetH
			}
		} else {
			return computeDimensions(rect, ScaleModeCover, targetWidth, targetHeight)
		}
	case ScaleModeDown:
		if width > targetH {
			if targetH/aspect > targetW {
				w = math.Floor(min(height, targetW))
				h = math.Floor(w * aspect)
			} else {
				h = targetH
				w = math.Floor(float64(h) / aspect)
			}
		} else if height > targetW {
			if targetW*aspect > targetH {
				h = math.Floor(min(width, targetH))
				w = math.Floor(h / aspect)
			} else {
				w = targetW
				h = math.Floor(w * aspect)
			}
		} else {
			w, h = height, width
		}
	case ScaleModeCover:
		if math.Floor(targetW*aspect) < targetH {
			h = targetH
			w = math.Floor(targetH / aspect)
		} else if math.Floor(targetH/aspect) < targetW {
			w = targetW
			h = math.Floor(targetW * aspect)
		}
	}
	return int(w), int(h)
}
