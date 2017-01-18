package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
)

type ImageFormat int

const (
	ImageFormatSame ImageFormat = 0
	ImageFormatJpeg             = 1
	ImageFormatPng              = 2
	ImageFormatGif              = 3
)

type ProcessingSpec struct {
	Format               ImageFormat
	NormalizeOrientation bool
	Quality              float64
	Trim                 TrimMode
	TrimFuzzFactor       float64
	Scale                ScaleMode
	ScaleWidth           int
	ScaleHeight          int
	Crop                 CropMode
	CropWidth            int
	CropHeight           int
}

func ProcessImage(reader io.ReadSeeker, writer io.Writer, spec ProcessingSpec) error {
	img, formatName, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	reader.Seek(0, 0)
	metadata := NewMetadataFromReader(reader)

	if !img.Bounds().Empty() {
		if spec.NormalizeOrientation && metadata.Exif != nil {
			if tag, err := metadata.Exif.Get("Orientation"); err == nil {
				img = NormalizeOrientation(img, tag.String())
			}
		}

		img = Trim(img, spec.Trim, spec.TrimFuzzFactor)
		img = Scale(img, spec.Scale, spec.ScaleWidth, spec.ScaleHeight)
		img = Crop(img, spec.Crop, spec.CropWidth, spec.CropHeight)
	}

	format := spec.Format
	if format == ImageFormatSame {
		switch formatName {
		case "jpeg":
			format = ImageFormatJpeg
		case "png":
			format = ImageFormatPng
		case "gif":
			format = ImageFormatGif
		default:
			panic(fmt.Sprintf("Don't know format %s", formatName))
		}
	}

	switch format {
	case ImageFormatPng:
		png.Encode(writer, img)
	case ImageFormatJpeg:
		jpeg.Encode(writer, img, &jpeg.Options{
			Quality: int(spec.Quality * 100),
		})
	case ImageFormatGif:
		gif.Encode(writer, img, &gif.Options{
			NumColors: 256,
			Quantizer: nil,
			Drawer:    nil,
		})
	}

	return nil
}
