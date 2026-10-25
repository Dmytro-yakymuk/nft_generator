package imageproc

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCropAndComposeOver(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	const size = 48
	var full, cropped []image.Image
	for k := 0; k < 5; k++ {
		img := image.NewNRGBA(image.Rect(0, 0, size, size))
		// a background, then semi-transparent rectangles somewhere on the canvas.
		x0, y0 := rnd.Intn(size/2), rnd.Intn(size/2)
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if k == 0 || (x >= x0 && x < x0+12 && y >= y0 && y < y0+9) {
					img.Set(x, y, color.NRGBA{R: uint8(rnd.Intn(256)), G: uint8(rnd.Intn(256)), B: uint8(rnd.Intn(256)), A: uint8(rnd.Intn(256))})
				}
			}
		}
		full = append(full, img)
		crop := CropToVisible(img)
		if k > 0 {
			require.Equal(t, image.Rect(x0, y0, x0+12, y0+9), crop.Bounds(), "layer %d", k)
		}
		cropped = append(cropped, crop)
	}

	// drawing the whole layers over each other gives the same image as drawing the visible parts.
	reference := image.NewRGBA(image.Rect(0, 0, size, size))
	for _, layer := range full {
		draw.Draw(reference, reference.Bounds(), layer, image.Point{}, draw.Over)
	}
	require.Equal(t, reference.Pix, ComposeOver(size, size, cropped).Pix)

	// a fully transparent layer is empty and changes nothing.
	require.True(t, CropToVisible(image.NewNRGBA(image.Rect(0, 0, size, size))).Bounds().Empty())
}
