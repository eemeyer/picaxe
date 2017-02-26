package imageops_test

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertImagesEqual(t *testing.T, expect, actual image.Image) bool {
	if !assert.Equal(t, expect.Bounds().Size(), actual.Bounds().Size(), "sizes are different") {
		return false
	}
	expectBytes := encodeImage(expect)
	actualBytes := encodeImage(actual)
	if !assert.Equal(t, len(expectBytes), len(actualBytes),
		"images are different lengths") {
		return false
	}
	if !reflect.DeepEqual(expectBytes, actualBytes) {
		t.Fatal("Image contents is different")
		return false
	}
	return true
}

func encodeImage(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatalf("Could not encode image: %s", err)
	}
	return buf.Bytes()
}

func loadImage(fileName string) image.Image {
	fileName = "../testdata/" + fileName
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("Unable to load %s: %s", fileName, err)
	}
	return img
}
