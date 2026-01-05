package main

import (
	"image"
	"sync"
)

func SobParallel(src image.Image, workers int) *image.RGBA {
	bounds := src.Bounds()
	height := bounds.Dy()

	out := image.NewRGBA(bounds)

	chunk := height / workers
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		yStart := bounds.Min.Y + i*chunk
		yEnd := yStart + chunk
		if i == workers-1 {
			yEnd = bounds.Max.Y
		}

		wg.Add(1)
		go Sobel(src, out, yStart, yEnd, &wg)
	}

	wg.Wait()
	return out
}

func Sobel(src image.Image,
	yStart int, yEnd int,
	wg *sync.WaitGroup){
	
	defer wg.Done()

	bounds := src.Bounds()
	
	GX := [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	GY := [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1},
	}


	for y := yStart; y < yEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumRx, sumGx, sumBx int
			var sumRy, sumGy, sumBy int

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {

					r, g, b, _ := img.At(x+i, y+j).RGBA()
					r8 := int(r >> 8)
					g8 := int(g >> 8)
					b8 := int(b >> 8)

					sumRx += r8 * GX[i+1][j+1]
					sumGx += g8 * GX[i+1][j+1]
					sumBx += b8 * GX[i+1][j+1]

					sumRy += r8 * GY[i+1][j+1]
					sumGy += g8 * GY[i+1][j+1]
					sumBy += b8 * GY[i+1][j+1]

				}
			}
			out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	return out
}

func main() {
	// Ouvrir l'image
	file, err := os.Open(".jpg")
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

	out := SobParallel(img, workers)

	// Sauvegarde
	savePNG("out.png", out)
	
