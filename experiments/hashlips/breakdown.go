//go:build ignore

// Розкладання часу одного зображення для E12 (п. 3.2.11 роботи): скільки платформа витрачає на накладання шарів
// і скільки на кодування PNG - тими самими засобами стандартної бібліотеки Go, що й сервер (image/draw, image/png);
// варіанти шарів декодуються один раз і обрізаються до непрозорих пікселів, як у кеші шарів генерації.
// Для порівняння те саме зображення кодується й з найшвидшим рівнем стиснення (png.BestSpeed).
// Варіанти обираються за вагами з назв файлів у форматі HashLips ("Назва#вага.png"; None - відсутність шару)
// з робочих каталогів experiments/hashlips/run.js. Запускати тією самою версією Go, якою зібрано сервер:
//
//	docker run --rm -v "$PWD":/src -w /src -e GO111MODULE=off golang:1.17-alpine \
//	    go run experiments/hashlips/breakdown.go 300 data/hashlips-runs/E1/layers/L0 ... L4
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const size = 500

type variant struct {
	image  image.Image // nil - шар відсутній
	weight int
}

func main() {
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	var decode time.Duration
	decoded := 0
	layers := make([][]variant, 0, len(os.Args)-2)
	for _, dir := range os.Args[2:] {
		files, err := filepath.Glob(filepath.Join(dir, "*.png"))
		if err != nil {
			panic(err)
		}
		sort.Strings(files)
		variants := make([]variant, 0, len(files))
		for _, file := range files {
			name := strings.TrimSuffix(filepath.Base(file), ".png")
			weight, err := strconv.Atoi(name[strings.LastIndex(name, "#")+1:])
			if err != nil {
				panic(err)
			}
			if strings.HasPrefix(name, "None#") {
				variants = append(variants, variant{weight: weight})
				continue
			}
			start := time.Now()
			img := decodePNG(file)
			decode += time.Since(start)
			decoded++
			variants = append(variants, variant{image: cropToVisible(img), weight: weight})
		}
		layers = append(layers, variants)
	}

	rnd := rand.New(rand.NewSource(1))
	var compose, encode, encodeFast time.Duration
	total, totalFast := 0, 0
	fast := png.Encoder{CompressionLevel: png.BestSpeed}
	for i := 0; i < count; i++ {
		parts := make([]image.Image, 0, len(layers))
		for _, variants := range layers {
			if v := pick(rnd, variants); v.image != nil {
				parts = append(parts, v.image)
			}
		}

		start := time.Now()
		canvas := image.NewRGBA(image.Rect(0, 0, size, size))
		for _, part := range parts {
			draw.Draw(canvas, part.Bounds(), part, part.Bounds().Min, draw.Over)
		}
		compose += time.Since(start)

		start = time.Now()
		var buf bytes.Buffer
		if err := png.Encode(&buf, canvas); err != nil {
			panic(err)
		}
		encode += time.Since(start)
		total += buf.Len()

		start = time.Now()
		buf.Reset()
		if err := fast.Encode(&buf, canvas); err != nil {
			panic(err)
		}
		encodeFast += time.Since(start)
		totalFast += buf.Len()
	}

	ms := func(d time.Duration, n int) float64 { return float64(d) / float64(time.Millisecond) / float64(n) }
	kib := func(bytes int) float64 { return float64(bytes) / float64(count) / 1024 }
	fmt.Printf("Go %s: декодування варіанта %.2f мс (%d варіантів, один раз); на зображення: накладання %.2f мс, "+
		"кодування PNG %.2f мс, %.1f КіБ; з png.BestSpeed - %.2f мс, %.1f КіБ\n",
		strings.TrimPrefix(runtime.Version(), "go"), ms(decode, decoded), decoded, ms(compose, count),
		ms(encode, count), kib(total), ms(encodeFast, count), kib(totalFast))
}

func decodePNG(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

// pick обирає варіант з імовірністю, пропорційною вазі.
func pick(rnd *rand.Rand, variants []variant) variant {
	sum := 0
	for _, v := range variants {
		sum += v.weight
	}
	r := rnd.Intn(sum)
	for _, v := range variants {
		if r -= v.weight; r < 0 {
			return v
		}
	}
	return variants[len(variants)-1]
}

// cropToVisible - копія imageproc.CropToVisible (пакет сервера тут не імпортується, щоб скрипт не потребував модулів).
func cropToVisible(img image.Image) *image.NRGBA {
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
