package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"sync"
	"time"
)

// Job représente une ligne de l'image à traiter
type Job struct {
	y int
}

func main() {

	// 1. Charger l'image
	file, err := os.Open("images/Marmite_UMARY_WM-94.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// 2. Préparer les images (Accès direct Pix)
	bounds := img.Bounds()
	srcRGBA := image.NewRGBA(bounds)
	draw.Draw(srcRGBA, bounds, img, bounds.Min, draw.Src)
	outRGBA := image.NewRGBA(bounds)

	t0 := time.Now()
	// 3. Lancer le traitement parallèle avec Worker Pool
	numWorkers := 8
	SobelWorkerPool(srcRGBA, outRGBA, numWorkers)
	t1 := time.Now()

	// 4. Sauvegarde
	savePNG("images/out.png", outRGBA)
	t2 := time.Now()
	println("Traitement terminé : out.png")
	println("Temps de calcul    : ", t1.Sub(t0))
	println("Temps de sauvegarde :", t2.Sub(t1))
}

func SobelWorkerPool(src *image.RGBA, out *image.RGBA, numWorkers int) {
	bounds := src.Bounds()
	jobs := make(chan Job, bounds.Dy())
	var wg sync.WaitGroup

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		go sobelWorker(src, out, jobs, &wg)
	}

	// Envoyer chaque ligne comme un job (on évite les bords y=0 et y=max)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		wg.Add(1)
		jobs <- Job{y: y}
	}

	close(jobs)
	wg.Wait()
}

func sobelWorker(src *image.RGBA, out *image.RGBA, jobs <-chan Job, wg *sync.WaitGroup) {
	bounds := src.Bounds()
	stride := src.Stride

	// Matrices Sobel
	gxMask := [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gyMask := [3][3]int{{1, 2, 1}, {0, 0, 0}, {-1, -2, -1}}

	for job := range jobs {
		y := job.y
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var Gx, Gy float64

			// Convolution 3x3
			for j := -1; j <= 1; j++ {
				for i := -1; i <= 1; i++ {
					pixelIdx := (y+j-bounds.Min.Y)*stride + (x+i-bounds.Min.X)*4
					r := float64(src.Pix[pixelIdx])
					g := float64(src.Pix[pixelIdx+1])
					b := float64(src.Pix[pixelIdx+2])
					lum := 0.299*r + 0.587*g + 0.114*b

					Gx += lum * float64(gxMask[j+1][i+1])
					Gy += lum * float64(gyMask[j+1][i+1])
				}
			}

			// Calcul manuel de la magnitude (Approximation rapide : |Gx| + |Gy|)
			mag := math.Sqrt(Gx*Gx + Gy*Gy)
			if mag > 255 {
				mag = 255
			}
			val := uint8(mag)

			outIdx := (y-bounds.Min.Y)*stride + (x-bounds.Min.X)*4
			out.Pix[outIdx] = val
			out.Pix[outIdx+1] = val
			out.Pix[outIdx+2] = val
			out.Pix[outIdx+3] = 255
		}
		wg.Done()
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
