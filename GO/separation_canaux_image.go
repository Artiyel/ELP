package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"sync"
)

func splitChannelsChunk(
	src image.Image,
	rImg, gImg, bImg *image.RGBA,
	yStart, yEnd int,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	bounds := src.Bounds()

	for y := yStart; y < yEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()

			// Conversion 16 bits -> 8 bits
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			rImg.Set(x, y, color.RGBA{r8, 0, 0, 255})
			gImg.Set(x, y, color.RGBA{0, g8, 0, 255})
			bImg.Set(x, y, color.RGBA{0, 0, b8, 255})
		}
	}
}

func SplitChannelsParallel(src image.Image, workers int) (
	*image.RGBA, *image.RGBA, *image.RGBA,
) {
	bounds := src.Bounds()

	rImg := image.NewRGBA(bounds)
	gImg := image.NewRGBA(bounds)
	bImg := image.NewRGBA(bounds)

	height := bounds.Dy()
	chunk := height / workers

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		yStart := bounds.Min.Y + i*chunk
		yEnd := yStart + chunk

		if i == workers-1 {
			yEnd = bounds.Max.Y
		}

		wg.Add(1)
		go splitChannelsChunk(
			src, rImg, gImg, bImg,
			yStart, yEnd, &wg,
		)
	}

	wg.Wait()
	return rImg, gImg, bImg
}

func savePNG(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, img)
}

func main() {
	// Ouvrir l'image
	file, err := os.Open("input.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	println("Format détecté :", format)

	workers := runtime.NumCPU()

	rImg, gImg, bImg := SplitChannelsParallel(img, workers)

	// Sauvegarde
	savePNG("red.png", rImg)
	savePNG("green.png", gImg)
	savePNG("blue.png", bImg)
}
