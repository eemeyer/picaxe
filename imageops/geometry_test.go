package imageops_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t11e/picaxe/imageops"
)

func TestFitDimensions(t *testing.T) {
	t.Run("landscape", func(t *testing.T) {
		t.Run("w and h", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 100),
					imageops.FitDimensions(image.Pt(100, 50), newInt(200), newInt(200)))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 150),
					imageops.FitDimensions(image.Pt(400, 300), newInt(200), newInt(200)))
			})
		})

		t.Run("w", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 100),
					imageops.FitDimensions(image.Pt(100, 50), newInt(200), nil))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 150),
					imageops.FitDimensions(image.Pt(400, 300), newInt(200), nil))
			})
		})

		t.Run("h", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(400, 200),
					imageops.FitDimensions(image.Pt(100, 50), nil, newInt(200)))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(267, 200),
					imageops.FitDimensions(image.Pt(400, 300), nil, newInt(200)))
			})
		})

		t.Run("no max", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(100, 50),
					imageops.FitDimensions(image.Pt(100, 50), nil, nil))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(400, 300),
					imageops.FitDimensions(image.Pt(400, 300), nil, nil))
			})
		})
	})

	t.Run("portrait", func(t *testing.T) {
		t.Run("w and h", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(100, 200),
					imageops.FitDimensions(image.Pt(50, 100), newInt(200), newInt(200)))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(150, 200),
					imageops.FitDimensions(image.Pt(300, 400), newInt(200), newInt(200)))
			})
		})

		t.Run("w", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 400),
					imageops.FitDimensions(image.Pt(50, 100), newInt(200), nil))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(200, 267),
					imageops.FitDimensions(image.Pt(300, 400), newInt(200), nil))
			})
		})

		t.Run("h", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(100, 200),
					imageops.FitDimensions(image.Pt(50, 100), nil, newInt(200)))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(150, 200),
					imageops.FitDimensions(image.Pt(300, 400), nil, newInt(200)))
			})
		})

		t.Run("no max", func(t *testing.T) {
			t.Run("input smaller", func(t *testing.T) {
				assert.Equal(t, image.Pt(50, 100),
					imageops.FitDimensions(image.Pt(50, 100), nil, nil))
			})

			t.Run("input larger", func(t *testing.T) {
				assert.Equal(t, image.Pt(300, 400),
					imageops.FitDimensions(image.Pt(300, 400), nil, nil))
			})
		})
	})
}

func newInt(v int) *int {
	return &v
}
