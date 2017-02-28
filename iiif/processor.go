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

type Processor interface {
	Process(
		spec string,
		resolver resources.Resolver,
		w io.Writer) error
}

type processor struct{}

// Process implements Processor.
func (processor) Process(
	spec string,
	resolver resources.Resolver,
	w io.Writer) error {
	req, err := ParseSpec(spec)
	if err != nil {
		return err
	}

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

	dims, err := req.Size.CalculateDimensions(img.Bounds(), image.Pt(6000, 6000))
	if err != nil {
		return err
	}
	img = imageops.Scale(img, dims)

	switch req.Format {
	case FormatDefault, FormatPNG:
		return png.Encode(w, img)
	case FormatJPEG:
		return jpeg.Encode(w, img, &jpeg.Options{
			Quality: 98,
		})
	case FormatGIF:
		return gif.Encode(w, img, &gif.Options{
			NumColors: 256,
			Quantizer: nil,
			Drawer:    nil,
		})
	}

	return fmt.Errorf("Unexpected format %q", req.Format)
}

var DefaultProcessor = processor{}
