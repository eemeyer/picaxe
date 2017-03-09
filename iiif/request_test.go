package iiif_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/t11e/picaxe/iiif"
	"github.com/t11e/picaxe/imageops"
)

func TestRequest_String(t *testing.T) {
	baseRequest := iiif.Request{
		Identifier: "http://i.imgur.com/J1XaOIa.jpg",
		Region: iiif.Region{
			Kind: iiif.RegionKindFull,
		},
		Size: iiif.Size{
			Kind: iiif.SizeKindMax,
		},
		Format: iiif.FormatPNG,
	}

	t.Run("identifier", func(t *testing.T) {
		req := baseRequest
		req.Identifier = "http://i.imgur.com/J1XaOIa.jpg"
		assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
	})

	t.Run("region", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			req := baseRequest
			req.Region.Kind = iiif.RegionKindFull
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
		t.Run("square", func(t *testing.T) {
			req := baseRequest
			req.Region.Kind = iiif.RegionKindSquare
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/square/max/default.png", req.String())
		})
		t.Run("relative", func(t *testing.T) {
			req := baseRequest
			req.Region.Kind = iiif.RegionKindRelative
			req.Region.Relative = &imageops.RelativeRegion{
				X: 0.255, Y: 0.255, W: 0.5, H: 0.5,
			}
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/pct:25.5,25.5,50,50/max/default.png", req.String())
		})
		t.Run("absolute", func(t *testing.T) {
			req := baseRequest
			req.Region.Kind = iiif.RegionKindAbsolute
			r := image.Rect(10, 10, 442, 244)
			req.Region.Absolute = &r
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/10,10,432,234/max/default.png", req.String())
		})
	})

	t.Run("size", func(t *testing.T) {
		t.Run("full", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindFull
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/full/default.png", req.String())
		})
		t.Run("max", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindMax
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
		t.Run("relative", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindRelative
			req.Size.Relative = newFloat64(0.505)
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/pct:50.5/default.png", req.String())
		})
		t.Run("absolute width and height", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindAbsolute
			req.Size.AbsWidth = newInt(100)
			req.Size.AbsHeight = newInt(200)
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/100,200/default.png", req.String())
		})
		t.Run("absolute width and height, best fit", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindAbsolute
			req.Size.AbsWidth = newInt(100)
			req.Size.AbsHeight = newInt(200)
			req.Size.AbsBestFit = true
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/!100,200/default.png", req.String())
		})
		t.Run("absolute width", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindAbsolute
			req.Size.AbsWidth = newInt(100)
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/100,/default.png", req.String())
		})
		t.Run("absolute height", func(t *testing.T) {
			req := baseRequest
			req.Size.Kind = iiif.SizeKindAbsolute
			req.Size.AbsHeight = newInt(200)
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/,200/default.png", req.String())
		})
	})

	t.Run("format", func(t *testing.T) {
		t.Run("png", func(t *testing.T) {
			req := baseRequest
			req.Format = iiif.FormatPNG
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
		t.Run("jpg", func(t *testing.T) {
			req := baseRequest
			req.Format = iiif.FormatJPEG
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.jpg", req.String())
		})
		t.Run("gif", func(t *testing.T) {
			req := baseRequest
			req.Format = iiif.FormatGIF
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.gif", req.String())
		})
	})

	t.Run("autoOrient", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			req := baseRequest
			req.AutoOrient = true
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png?autoOrient=true", req.String())
		})
		t.Run("false", func(t *testing.T) {
			req := baseRequest
			req.AutoOrient = false
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
	})

	t.Run("trimBorder", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			req := baseRequest
			req.TrimBorder = true
			req.TrimBorderFuzziness = 0.33
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png?trimBorder=0.33", req.String())
		})
		t.Run("false", func(t *testing.T) {
			req := baseRequest
			req.TrimBorder = false
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
	})

	t.Run("scale=down", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			req := baseRequest
			req.Size.AbsDoNotEnlarge = true
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png?scale=down", req.String())
		})
		t.Run("false", func(t *testing.T) {
			req := baseRequest
			req.TrimBorder = false
			assert.Equal(t, "http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/max/default.png", req.String())
		})
	})
}

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

func TestParseSpec(t *testing.T) {
	for _, test := range []struct {
		in            string
		expected      *iiif.Request
		expectedError string
	}{
		{
			in:            "invalid",
			expectedError: `not a valid spec: "invalid"`,
		},
		{
			in: "identifier/full/max/0/default.png",
			expected: &iiif.Request{
				Identifier: "identifier",
				Region:     iiif.Region{Kind: iiif.RegionKindFull},
				Size:       iiif.Size{Kind: iiif.SizeKindMax},
				Format:     iiif.FormatPNG,
			},
		},
		{
			in: "identifier/full/max/0/default.png?trimBorder=0.5",
			expected: &iiif.Request{
				Identifier:          "identifier",
				Region:              iiif.Region{Kind: iiif.RegionKindFull},
				Size:                iiif.Size{Kind: iiif.SizeKindMax},
				Format:              iiif.FormatPNG,
				TrimBorder:          true,
				TrimBorderFuzziness: 0.5,
			},
		},
		{
			in: "identifier/full/max/0/default.png?autoOrient=true",
			expected: &iiif.Request{
				Identifier: "identifier",
				Region:     iiif.Region{Kind: iiif.RegionKindFull},
				Size:       iiif.Size{Kind: iiif.SizeKindMax},
				Format:     iiif.FormatPNG,
				AutoOrient: true,
			},
		},
		{
			in: "identifier/full/max/0/default.png?scale=down",
			expected: &iiif.Request{
				Identifier: "identifier",
				Region:     iiif.Region{Kind: iiif.RegionKindFull},
				Size:       iiif.Size{Kind: iiif.SizeKindMax, AbsDoNotEnlarge: true},
				Format:     iiif.FormatPNG,
			},
		},
		{
			in:            "identifier/full/max/0/default.png?scale=invalid",
			expectedError: `not a valid scale: "invalid"`,
		},
		{
			in: "identifier/full/max/0/default.png?trimBorder=0.5&autoOrient=true&scale=down",
			expected: &iiif.Request{
				Identifier:          "identifier",
				Region:              iiif.Region{Kind: iiif.RegionKindFull},
				Size:                iiif.Size{Kind: iiif.SizeKindMax, AbsDoNotEnlarge: true},
				Format:              iiif.FormatPNG,
				TrimBorder:          true,
				TrimBorderFuzziness: 0.5,
				AutoOrient:          true,
			},
		},
	} {
		t.Run(test.in, func(t *testing.T) {
			actual, err := iiif.ParseSpec(test.in)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSize_CalculateDimensions(t *testing.T) {
	type scenario struct {
		description   string
		in            image.Point
		maxSize       image.Point
		expected      image.Point
		expectedError string
	}
	for _, test := range []struct {
		size      iiif.Size
		scenarios []scenario
	}{
		{
			size: iiif.Size{
				Kind: iiif.SizeKindFull,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(500, 600) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind: iiif.SizeKindMax,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description: "larger than max",
					in:          image.Point{500, 600},
					maxSize:     image.Point{100, 200},
					expected:    image.Point{100, 120},
				},
			},
		},
		{
			size: iiif.Size{
				Kind:      iiif.SizeKindAbsolute,
				AbsWidth:  newInt(300),
				AbsHeight: newInt(400),
			},
			scenarios: []scenario{
				{
					description:   "control",
					in:            image.Point{},
					maxSize:       image.Point{},
					expectedError: "(300, 400) exceeds maximum allowed dimensions (0, 0)",
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{300, 400},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 400},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 400},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:     iiif.SizeKindAbsolute,
				AbsWidth: newInt(300),
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{0, 0},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:      iiif.SizeKindAbsolute,
				AbsHeight: newInt(400),
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description:   "larger than max",
					in:            image.Point{50, 50},
					maxSize:       image.Point{100, 200},
					expectedError: "(400, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:       iiif.SizeKindAbsolute,
				AbsWidth:   newInt(300),
				AbsHeight:  newInt(400),
				AbsBestFit: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:       iiif.SizeKindAbsolute,
				AbsWidth:   newInt(300),
				AbsBestFit: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{0, 0},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:       iiif.SizeKindAbsolute,
				AbsHeight:  newInt(400),
				AbsBestFit: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "scales up, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description:   "larger than max",
					in:            image.Point{50, 50},
					maxSize:       image.Point{100, 200},
					expectedError: "(400, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsWidth:        newInt(300),
				AbsHeight:       newInt(400),
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 400},
				},
				{
					description: "does not enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsWidth:        newInt(300),
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description: "does not enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsHeight:       newInt(400),
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "does not enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{400, 400},
					maxSize:       image.Point{100, 200},
					expectedError: "(400, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsWidth:        newInt(300),
				AbsHeight:       newInt(400),
				AbsBestFit:      true,
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "does not enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsWidth:        newInt(300),
				AbsBestFit:      true,
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{300, 600},
				},
				{
					description: "does not enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(300, 360) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:            iiif.SizeKindAbsolute,
				AbsHeight:       newInt(400),
				AbsBestFit:      true,
				AbsDoNotEnlarge: true,
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "empty image",
					in:          image.Point{},
					maxSize:     image.Point{500, 500},
					expected:    image.Point{},
				},
				{
					description: "scales down",
					in:          image.Point{1000, 2000},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{200, 400},
				},
				{
					description: "does enlarge, smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{100, 200},
				},
				{
					description:   "larger than max",
					in:            image.Point{400, 400},
					maxSize:       image.Point{100, 200},
					expectedError: "(400, 400) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
		{
			size: iiif.Size{
				Kind:     iiif.SizeKindRelative,
				Relative: newFloat64(0.5),
			},
			scenarios: []scenario{
				{
					description: "control",
					in:          image.Point{},
					maxSize:     image.Point{},
					expected:    image.Point{},
				},
				{
					description: "smaller than max",
					in:          image.Point{100, 200},
					maxSize:     image.Point{500, 600},
					expected:    image.Point{50, 100},
				},
				{
					description:   "larger than max",
					in:            image.Point{500, 600},
					maxSize:       image.Point{100, 200},
					expectedError: "(250, 300) exceeds maximum allowed dimensions (100, 200)",
				},
			},
		},
	} {
		description := test.size.String()
		if test.size.AbsDoNotEnlarge {
			description += "?scale=down"
		}
		t.Run(description, func(t *testing.T) {
			for _, scenario := range test.scenarios {
				t.Run(scenario.description, func(t *testing.T) {
					actual, err := test.size.CalculateDimensions(scenario.in, scenario.maxSize)
					if scenario.expectedError != "" {
						assert.EqualError(t, err, scenario.expectedError)
					} else {
						assert.NoError(t, err)
					}
					assert.Equal(t, scenario.expected, actual)
				})
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
