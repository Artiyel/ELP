package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"runtime"
	"sync"
)

// Job correspond à une ligne à traiter
type Job struct {
	y int
}

// Worker pool pour séparer les canaux de couleur
func worker(rgba *image.RGBA, rImg, gImg, bImg *image.RGBA, jobs <-chan Job, wg *sync.WaitGroup) {
	bounds := rgba.Bounds()
	stride := rgba.Stride

	for job := range jobs {
		y := job.y
		rowOffset := (y - bounds.Min.Y) * stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := rowOffset + (x-bounds.Min.X)*4
			r := rgba.Pix[i]
			g := rgba.Pix[i+1]
			b := rgba.Pix[i+2]

			rImg.Pix[i] = r
			rImg.Pix[i+1] = 0
			rImg.Pix[i+2] = 0
			rImg.Pix[i+3] = 255

			gImg.Pix[i] = 0
			gImg.Pix[i+1] = g
			gImg.Pix[i+2] = 0
			gImg.Pix[i+3] = 255

			bImg.Pix[i] = 0
			bImg.Pix[i+1] = 0
			bImg.Pix[i+2] = b
			bImg.Pix[i+3] = 255
		}
		wg.Done()
	}
}

// SplitChannelsWorkerPool utilise un worker pool
func SplitChannelsWorkerPool(rgba *image.RGBA, numWorkers int) (*image.RGBA, *image.RGBA, *image.RGBA) {
	bounds := rgba.Bounds()

	rImg := image.NewRGBA(bounds)
	gImg := image.NewRGBA(bounds)
	bImg := image.NewRGBA(bounds)

	jobs := make(chan Job, bounds.Dy())
	var wg sync.WaitGroup

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		go worker(rgba, rImg, gImg, bImg, jobs, &wg)
	}

	// Envoyer les jobs (une ligne = un job)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		jobs <- Job{y: y}
	}

	// Tous les jobs sont envoyés, fermer le channel
	close(jobs)

	// Attendre que tout soit fini
	wg.Wait()
	return rImg, gImg, bImg
}

// Sauvegarde en PNG
func savePNG(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	// Forcer l'utilisation de tous les cœurs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Ouvrir l'image
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

	// Convertir en RGBA pour un accès direct à la mémoire
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	numWorkers := runtime.NumCPU()
	rImg, gImg, bImg := SplitChannelsWorkerPool(rgba, numWorkers)

	// Sauvegarde
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { savePNG("images/red.png", rImg); wg.Done() }()
	go func() { savePNG("images/green.png", gImg); wg.Done() }()
	go func() { savePNG("images/blue.png", bImg); wg.Done() }()
	wg.Wait()

	println("Terminé !")
}
