package imageops_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t11e/picaxe/imageops"
)

func TestTrim(t *testing.T) {
	for _, test := range []struct {
		fileName     string
		fuzz         float64
		expectedSize image.Point
	}{
		{fileName: "border-1.jpeg", fuzz: 0, expectedSize: image.Pt(812, 622)},
		{fileName: "border-2.jpeg", fuzz: 0, expectedSize: image.Pt(980, 650)},
		{fileName: "border-3.jpeg", fuzz: 0, expectedSize: image.Pt(531, 650)},
		{fileName: "border-4.jpeg", fuzz: 0, expectedSize: image.Pt(808, 618)},

		{fileName: "border-1.jpeg", fuzz: 0.1, expectedSize: image.Pt(784, 590)},
		{fileName: "border-2.jpeg", fuzz: 0.1, expectedSize: image.Pt(980, 650)},
		{fileName: "border-3.jpeg", fuzz: 0.1, expectedSize: image.Pt(493, 650)},
		{fileName: "border-4.jpeg", fuzz: 0.1, expectedSize: image.Pt(784, 590)},

		{fileName: "border-1.jpeg", fuzz: 1.0, expectedSize: image.Pt(0, 0)},
		{fileName: "border-2.jpeg", fuzz: 1.0, expectedSize: image.Pt(0, 0)},
		{fileName: "border-3.jpeg", fuzz: 1.0, expectedSize: image.Pt(0, 0)},
		{fileName: "border-4.jpeg", fuzz: 1.0, expectedSize: image.Pt(0, 0)},
	} {
		t.Run(fmt.Sprintf("%s/fuzz=%f", test.fileName, test.fuzz), func(t *testing.T) {
			img := loadImage(test.fileName)
			trimmed := imageops.Trim(img, test.fuzz)
			assert.Equal(t, test.expectedSize, trimmed.Bounds().Size())
		})
	}
}
