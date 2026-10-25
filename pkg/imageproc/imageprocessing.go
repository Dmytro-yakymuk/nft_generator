package imageproc

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fogleman/gg"
	"github.com/zeebo/errs"
)

// TypeFile defines the list of possible type of files.
type TypeFile string

const (
	// TypeFilePNG indicates that the type file is png.
	TypeFilePNG TypeFile = "png"
	// TypeFileSVG indicates that the type file is svg.
	TypeFileSVG TypeFile = "svg"
	// TypeFileJPEG indicates that the type file is jpeg.
	TypeFileJPEG TypeFile = "jpeg"
	// TypeFileJSON indicates that the type file is json.
	TypeFileJSON TypeFile = "json"
)

// LayerComponentsCount searches count files in the specified path and by name of file.
func LayerComponentsCount(pathToLayerComponents, nameFile string) (int, error) {
	files, err := ioutil.ReadDir(pathToLayerComponents)
	if err != nil {
		return 0, fmt.Errorf(pathToLayerComponents + " - folder does not exist")
	}

	var count int
	for _, file := range files {
		isMatched, err := regexp.Match(fmt.Sprintf(nameFile, `\d`), []byte(file.Name()))
		if err != nil {
			return 0, err
		}
		if isMatched {
			count++
		}
	}

	return count, nil
}

// CreateLayer searches and decodes image to layer.
func CreateLayer(path, name string) (image.Image, error) {
	image, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errs.Combine(err, image.Close())
	}()

	layer, err := png.Decode(image)
	if err != nil {
		return nil, err
	}
	return layer, err
}

// CropToVisible returns a copy of the part of the image within the bounding box of its pixels that are not
// fully transparent, in the same coordinates, so that drawing it touches only those pixels and keeping it
// in memory takes only their size. A fully transparent image gives an empty one.
func CropToVisible(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	nrgba, ok := img.(*image.NRGBA)
	if !ok {
		nrgba = image.NewNRGBA(bounds)
		draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	}

	minX, minY, maxX, maxY := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := nrgba.Pix[(y-bounds.Min.Y)*nrgba.Stride:]
		for x := 0; x < bounds.Dx(); x++ {
			if row[x*4+3] == 0 {
				continue
			}
			if bounds.Min.X+x < minX {
				minX = bounds.Min.X + x
			}
			if bounds.Min.X+x+1 > maxX {
				maxX = bounds.Min.X + x + 1
			}
			if y < minY {
				minY = y
			}
			maxY = y + 1
		}
	}

	visible := image.Rectangle{}
	if minX < maxX {
		visible = image.Rect(minX, minY, maxX, maxY)
	}

	cropped := image.NewNRGBA(visible)
	draw.Draw(cropped, visible, nrgba, visible.Min, draw.Src)
	return cropped
}

// ComposeOver draws the layers in order on a transparent canvas of the given size by the Over operator,
// each layer within its own bounds. On the transparent canvas the first layer is copied as it is.
func ComposeOver(width, height int, layers []image.Image) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, layer := range layers {
		draw.Draw(canvas, layer.Bounds(), layer, layer.Bounds().Min, draw.Over)
	}
	return canvas
}

// SaveImage saves image by path.
func SaveImage(path, fullPath string, baseImage image.Image) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	resultImage, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	if err = png.Encode(resultImage, baseImage); err != nil {
		return err
	}
	defer func() {
		err = errs.Combine(err, resultImage.Close())
	}()

	return nil
}

// Inscription entity describes values required to apply inscription to the image.
type Inscription struct {
	Img         image.Image
	Width       int
	Height      int
	PathToFonts string
	FontSize    float64
	FontColor   string
	Text        string
	X           float64
	Y           float64
	XAlign      float64
	YAlign      float64
	Rotate      float64
}

// ApplyInscription overlays the inscription on the image.
func ApplyInscription(inscription Inscription) (image.Image, error) {
	dc := gg.NewContext(inscription.Width, inscription.Height)
	if err := dc.LoadFontFace(inscription.PathToFonts, inscription.FontSize); err != nil {
		return nil, err
	}

	dc.SetHexColor(inscription.FontColor)
	dc.DrawImage(inscription.Img, 0, 0)
	dc.RotateAbout(inscription.Rotate, inscription.X, inscription.Y)
	dc.DrawStringAnchored(inscription.Text, inscription.X, inscription.Y, inscription.XAlign, inscription.YAlign)
	dc.Clip()
	return dc.Image(), nil
}
