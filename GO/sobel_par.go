package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"
)

// Job représente une ligne de l'image à traiter
type Job struct {
	y int
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 1. Charger l'image
	file, err := os.Open("images/heic1509a.jpg")
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
	numWorkers := 1
	SobelWorkerPool(srcRGBA, outRGBA, numWorkers)
	t1 := time.Now()

	// 4. Sauvegarde
	savePNG("images/out.png", outRGBA)
	t2 := time.Now()
	println("Traitement terminé : out.png")
	println("Temps de calcul : ", t1.Sub(t0))
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
			var rGx, gGx, bGx int
			var rGy, gGy, bGy int

			// Convolution 3x3
			for j := -1; j <= 1; j++ {
				for i := -1; i <= 1; i++ {
					pixelIdx := (y+j)*stride + (x+i)*4

					r := int(src.Pix[pixelIdx])
					g := int(src.Pix[pixelIdx+1])
					b := int(src.Pix[pixelIdx+2])

					wx := gxMask[j+1][i+1]
					wy := gyMask[j+1][i+1]

					rGx += r * wx
					gGx += g * wx
					bGx += b * wx

					rGy += r * wy
					gGy += g * wy
					bGy += b * wy
				}
			}

			// Calcul manuel de la magnitude (Approximation rapide : |Gx| + |Gy|)
			// Rouge
			absRGx, absRGy := rGx, rGy
			if absRGx < 0 {
				absRGx = -absRGx
			}
			if absRGy < 0 {
				absRGy = -absRGy
			}
			magR := absRGx + absRGy

			// Vert
			absGGx, absGGy := gGx, gGy
			if absGGx < 0 {
				absGGx = -absGGx
			}
			if absGGy < 0 {
				absGGy = -absGGy
			}
			magG := absGGx + absGGy

			// Bleu
			absBGx, absBGy := bGx, bGy
			if absBGx < 0 {
				absBGx = -absBGx
			}
			if absBGy < 0 {
				absBGy = -absBGy
			}
			magB := absBGx + absBGy

			// Écriture directe et saturation à 255
			outIdx := y*stride + x*4
			if magR > 255 {
				magR = 255
			}
			if magG > 255 {
				magG = 255
			}
			if magB > 255 {
				magB = 255
			}

			out.Pix[outIdx] = uint8(magR)
			out.Pix[outIdx+1] = uint8(magG)
			out.Pix[outIdx+2] = uint8(magB)
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
