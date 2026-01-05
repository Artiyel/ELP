package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"runtime"
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

func Sobel(src image.Image, out *image.RGBA, yStart int, yEnd int, wg *sync.WaitGroup) {
	defer wg.Done()

	bounds := src.Bounds()

	// Matrices de Sobel
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

	// On itère en restant à 1 pixel des bords pour éviter les dépassements d'index avec [i+1][j+1]
	for y := yStart; y < yEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			// Ignorer les pixels de bordure pure pour simplifier
			if x == bounds.Min.X || x >= bounds.Max.X-1 || y == bounds.Min.Y || y >= bounds.Max.Y-1 {
				continue
			}

			var valRx, valGx, valBx int
			var valRy, valGy, valBy int

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					// Correction : Utilisation de src.At et non img.At
					r, g, b, _ := src.At(x+i, y+j).RGBA()

					// Conversion en 8 bits (0-255)
					r8, g8, b8 := int(r>>8), int(g>>8), int(b>>8)

					// Application des masques
					weightX := GX[j+1][i+1] // j=ligne, i=colonne
					weightY := GY[j+1][i+1]

					valRx += r8 * weightX
					valGx += g8 * weightX
					valBx += b8 * weightX

					valRy += r8 * weightY
					valGy += g8 * weightY
					valBy += b8 * weightY
				}
			}

			// Calcul de la magnitude finale : sqrt(Gx^2 + Gy^2)
			resR := math.Sqrt(float64(valRx*valRx + valRy*valRy))
			resG := math.Sqrt(float64(valGx*valGx + valGy*valGy))
			resB := math.Sqrt(float64(valBx*valBx + valBy*valBy))

			// On sature à 255
			out.Set(x, y, color.RGBA{
				uint8(math.Min(255, resR)),
				uint8(math.Min(255, resG)),
				uint8(math.Min(255, resB)),
				255,
			})
		}
	}
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.Open("images/heic1501a.jpg")
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

}
