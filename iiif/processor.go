package iiif

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/t11e/picaxe/imageops"
	"github.com/t11e/picaxe/resources"
)

//go:generate sh -c "mockery -name='Processor' -case=underscore"

// maxScaleSize is the largest image we will scale to.
var maxScaleSize = image.Pt(6000, 6000)

type Result struct {
	ContentType string
}

type Processor interface {
	Process(
		req Request,
		resolver resources.Resolver,
		w io.Writer,
		result *Result) error
}

type processor struct{}

// Process implements Processor.
func (processor) Process(
	req Request,
	resolver resources.Resolver,
	w io.Writer,
	result *Result) error {
	r, err := resolver.GetResource(req.Identifier)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	if req.AutoOrient {
		r.Seek(0, 0)
		metadata := imageops.NewMetadataFromReader(r)
		if metadata.Exif != nil {
			if tag, e := metadata.Exif.Get("Orientation"); e == nil {
				img = imageops.NormalizeOrientation(img, tag.String())
			}
		}
	}

	if req.TrimBorder {
		img = imageops.Trim(img, req.TrimBorderFuzziness)
	}

	switch req.Region.Kind {
	case RegionKindAbsolute:
		img = imageops.CropRect(img, *req.Region.Absolute)
	case RegionKindRelative:
		img = imageops.CropRelative(img, *req.Region.Relative)
	case RegionKindSquare:
		img = imageops.CropSquare(img)
	}

	dims, err := req.Size.CalculateDimensions(img.Bounds().Size(), maxScaleSize)
	if err != nil {
		return err
	}
	img = imageops.Scale(img, dims)

	switch req.Format {
	case FormatPNG:
		if result != nil {
			result.ContentType = "image/png"
		}
		return png.Encode(w, img)
	case FormatJPEG:
		if result != nil {
			result.ContentType = "image/jpeg"
		}
		return jpeg.Encode(w, img, &jpeg.Options{
			Quality: 98,
		})
	case FormatGIF:
		if result != nil {
			result.ContentType = "image/gif"
		}
		return gif.Encode(w, img, &gif.Options{
			NumColors: 256,
			Quantizer: nil,
			Drawer:    nil,
		})
	}

	return fmt.Errorf("Unexpected format %q", req.Format)
}

var DefaultProcessor = processor{}
