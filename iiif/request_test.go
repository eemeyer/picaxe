package iiif_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/t11e/picaxe/iiif"
	"github.com/t11e/picaxe/imageops"
)

func TestParseSpec_invalid(t *testing.T) {
	var err error

	_, err = iiif.ParseSpec("some-identifier/full/max/0")
	assert.Error(t, err)

	_, err = iiif.ParseSpec("full/max/0/default.png")
	assert.Error(t, err)

	_, err = iiif.ParseSpec("default.png")
	assert.Error(t, err)

	_, err = iiif.ParseSpec("stuff")
	assert.Error(t, err)
}

func TestParseSpec_rotation(t *testing.T) {
	var err error

	_, err = iiif.ParseSpec("some-identifier/full/max/0/default.png")
	assert.NoError(t, err)

	for i := 1; i <= 360; i++ {
		_, err = iiif.ParseSpec(fmt.Sprintf("some-identifier/full/max/%d/default.png", i))
		assert.Error(t, err)
	}
}

func TestParseSpec_quality(t *testing.T) {
	var err error

	_, err = iiif.ParseSpec("some-identifier/full/max/0/default.png")
	assert.NoError(t, err)

	_, err = iiif.ParseSpec("some-identifier/full/max/0/color.png")
	assert.NoError(t, err)

	_, err = iiif.ParseSpec("some-identifier/full/max/0/grayscale.png")
	assert.Error(t, err)
}

func TestParseSpec_format(t *testing.T) {
	for _, test := range []struct {
		format       string
		expectFormat iiif.Format
		expectError  string
	}{
		{format: "png", expectFormat: iiif.FormatPNG},
		{format: "jpg", expectFormat: iiif.FormatJPEG},
		{format: "gif", expectFormat: iiif.FormatGIF},
		{format: "", expectError: `not a valid spec: "some-identifier/full/max/0/default."`},
		{format: "tif", expectError: "unsupported format \"tif\""},
		{format: "jp2", expectError: "unsupported format \"jp2\""},
		{format: "pdf", expectError: "unsupported format \"pdf\""},
		{format: "webp", expectError: "unsupported format \"webp\""},
		{format: "a4", expectError: "unsupported format \"a4\""},
	} {
		t.Run(test.format, func(t *testing.T) {
			req, err := iiif.ParseSpec(
				fmt.Sprintf("some-identifier/full/max/0/default.%s", test.format))
			if test.expectError != "" {
				if !assert.Error(t, err) {
					return
				}
				assert.Equal(t, test.expectError, err.Error())
			} else {
				if !assert.NoError(t, err) {
					return
				}
				if !assert.NotNil(t, req) {
					return
				}
				assert.Equal(t, test.expectFormat, req.Format)
			}
		})
	}
}

func TestParseSpec_region(t *testing.T) {
	for _, test := range []struct {
		region       string
		expectResult iiif.Region
		expectError  string
	}{
		{
			region: "full",
			expectResult: iiif.Region{
				Kind: iiif.RegionKindFull,
			},
		},

		{
			region: "square",
			expectResult: iiif.Region{
				Kind: iiif.RegionKindSquare,
			},
		},

		{
			region: "pct:0,0,100,100",
			expectResult: iiif.Region{
				Kind:     iiif.RegionKindRelative,
				Relative: &imageops.RelativeRegion{X: 0, Y: 0, W: 1, H: 1},
			},
		},

		{
			region: "pct:0.0,0.0,100.0,100.0",
			expectResult: iiif.Region{
				Kind:     iiif.RegionKindRelative,
				Relative: &imageops.RelativeRegion{X: 0, Y: 0, W: 1, H: 1},
			},
		},

		{region: "pct:-1.0,0.0,100,100", expectError: "Invalid percentage number: -1.0"},
		{region: "pct:200.0,0.0,100,100", expectError: "Invalid percentage number: 200.0"},
		{region: "pct:0,-1.0,100,100", expectError: "Invalid percentage number: -1.0"},
		{region: "pct:0,200.0,100,100", expectError: "Invalid percentage number: 200.0"},
		{region: "pct:0,0,-1.0,100", expectError: "Invalid percentage number: -1.0"},
		{region: "pct:0,0,200.0,100", expectError: "Invalid percentage number: 200.0"},
		{region: "pct:0,0,100,-1.0", expectError: "Invalid percentage number: -1.0"},
		{region: "pct:0,0,100,200.0", expectError: "Invalid percentage number: 200.0"},

		{
			region: "0,0,100,100",
			expectResult: iiif.Region{
				Kind: iiif.RegionKindAbsolute,
				Absolute: &image.Rectangle{
					Min: image.Pt(0, 0),
					Max: image.Pt(100, 100),
				},
			},
		},

		{
			region: "-10,-10,100,100",
			expectResult: iiif.Region{
				Kind: iiif.RegionKindAbsolute,
				Absolute: &image.Rectangle{
					Min: image.Pt(-10, -10),
					Max: image.Pt(90, 90),
				},
			},
		},

		{
			region: "0,0,-100,-100",
			expectResult: iiif.Region{
				Kind: iiif.RegionKindAbsolute,
				Absolute: &image.Rectangle{
					Min: image.Pt(0, 0),
					Max: image.Pt(0, 0),
				},
			},
		},

		{
			region:      "not,0,coordinates,1",
			expectError: "Not a valid set of coordinates: not,0,coordinates,1",
		},
	} {
		t.Run(test.region, func(t *testing.T) {
			req, err := iiif.ParseSpec(
				fmt.Sprintf("some-identifier/%s/max/0/default.png", test.region))
			if test.expectError != "" {
				if !assert.Error(t, err) {
					return
				}
				assert.Equal(t, test.expectError, err.Error())
			} else {
				if !assert.NoError(t, err) {
					return
				}
				if !assert.NotNil(t, req) {
					return
				}
				assert.Equal(t, test.expectResult, req.Region)
			}
		})
	}
}

func TestParseSpec_size(t *testing.T) {
	for _, test := range []struct {
		size         string
		expectResult iiif.Size
		expectError  string
	}{
		{
			size: "full",
			expectResult: iiif.Size{
				Kind: iiif.SizeKindFull,
			},
		},

		{
			size: "max",
			expectResult: iiif.Size{
				Kind: iiif.SizeKindMax,
			},
		},

		{
			size: "pct:50",
			expectResult: iiif.Size{
				Kind:     iiif.SizeKindRelative,
				Relative: newFloat64(0.5),
			},
		},

		{
			size: "100,200",
			expectResult: iiif.Size{
				Kind:      iiif.SizeKindAbsolute,
				AbsWidth:  newInt(100),
				AbsHeight: newInt(200),
			},
		},

		{
			size: "100,",
			expectResult: iiif.Size{
				Kind:     iiif.SizeKindAbsolute,
				AbsWidth: newInt(100),
			},
		},

		{
			size: ",200",
			expectResult: iiif.Size{
				Kind:      iiif.SizeKindAbsolute,
				AbsHeight: newInt(200),
			},
		},

		{
			size: "!100,200",
			expectResult: iiif.Size{
				Kind:       iiif.SizeKindAbsolute,
				AbsWidth:   newInt(100),
				AbsHeight:  newInt(200),
				AbsBestFit: true,
			},
		},

		{size: "pct:101", expectError: "Invalid percentage number: 101"},
		{size: "pct:-1", expectError: "Invalid percentage number: -1"},

		{
			size:        "not,anything,valid",
			expectError: "Not a valid width/height: not,anything,valid",
		},
	} {
		t.Run(test.size, func(t *testing.T) {
			req, err := iiif.ParseSpec(
				fmt.Sprintf("some-identifier/full/%s/0/default.png", test.size))
			if test.expectError != "" {
				if !assert.Error(t, err) {
					return
				}

			} else {
				if !assert.NoError(t, err) {
					return
				}
				if !assert.NotNil(t, req) {
					return
				}
				assert.Equal(t, test.expectResult, req.Size)
			}
		})
	}
}

func newFloat64(v float64) *float64 {
	return &v
}

func newInt(v int) *int {
	return &v
}
