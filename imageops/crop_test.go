package imageops_test

import (
	"image"
	_ "image/jpeg"
	"log"
	"testing"

	"github.com/t11e/picaxe/imageops"
)

func TestCropSquare(t *testing.T) {
	assertImagesEqual(t,
		loadImage("hippos-crop-square.png"),
		imageops.CropSquare(loadImage("hippos.png")))
}

func TestCropRect(t *testing.T) {
	assertImagesEqual(t,
		loadImage("hippos-crop-50,50,150,150.png"),
		imageops.CropRect(loadImage("hippos.png"), image.Rect(50, 50, 150, 150)))
}

func TestCropRelative(t *testing.T) {
	region, err := imageops.NewRelativeRegion(0.5, 0.5, 0.25, 0.25)
	if err != nil {
		log.Fatal(err)
	}
	assertImagesEqual(t,
		loadImage("hippos-crop-320,240,480,360.png"),
		imageops.CropRelative(loadImage("hippos.png"), *region))
}
