package main

import (
	"flag"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"
)

// Job correspond à un bloc de lignes à traiter
type Job struct {
	yStart, yEnd int
}

// Worker pool pour séparer les canaux de couleur
func worker(rgba, rImg, gImg, bImg *image.RGBA, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	bounds := rgba.Bounds()
	stride := rgba.Stride

	for job := range jobs {
		for y := job.yStart; y < job.yEnd; y++ {
			rowOffset := (y - bounds.Min.Y) * stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := rowOffset + (x-bounds.Min.X)*4
				r := rgba.Pix[i]
				g := rgba.Pix[i+1]
				b := rgba.Pix[i+2]

				// Rouge
				rImg.Pix[i] = r
				rImg.Pix[i+1] = 0
				rImg.Pix[i+2] = 0
				rImg.Pix[i+3] = 255

				// Vert
				gImg.Pix[i] = 0
				gImg.Pix[i+1] = g
				gImg.Pix[i+2] = 0
				gImg.Pix[i+3] = 255

				// Bleu
				bImg.Pix[i] = 0
				bImg.Pix[i+1] = 0
				bImg.Pix[i+2] = b
				bImg.Pix[i+3] = 255
			}
		}
		wg.Done()
	}
}

// SplitChannelsWorkerPool utilise un worker pool
func SplitChannelsWorkerPool(rgba *image.RGBA, numWorkers, blockSize int) (*image.RGBA, *image.RGBA, *image.RGBA) {
	bounds := rgba.Bounds()

	rImg := image.NewRGBA(bounds)
	gImg := image.NewRGBA(bounds)
	bImg := image.NewRGBA(bounds)

	jobs := make(chan Job, numWorkers*2) // buffer pour éviter blocage

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		go worker(rgba, rImg, gImg, bImg, jobs, &wg)
	}

	// Créer les jobs par bloc de `blockSize` lignes
	for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
		yEnd := y + blockSize
		if yEnd > bounds.Max.Y {
			yEnd = bounds.Max.Y
		}
		wg.Add(1)
		jobs <- Job{yStart: y, yEnd: yEnd}
	}

	close(jobs)
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

	numWorkers := flag.Int("n", runtime.NumCPU(), "Nombre de workers")
	sauvegarde := flag.Bool("s", true, "Active ou non la sauvegarde du fichier résultant")
	path := flag.String("f", "images/Marmite_UMARY_WM-94.jpg", "Précise l'emplacement de l'image a traiter (Présente dans un dossier 'images'")
	flag.Parse()

	println("Nombre de workers:", *numWorkers)

	// Ouvrir l'image
	file, err := os.Open(*path)
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

	// Traitement en worker pool avec blocs de 50 lignes
	tjob := time.Now()
	rImg, gImg, bImg := SplitChannelsWorkerPool(rgba, *numWorkers, 100)
	println("temps de calcul     : ", time.Now().Sub(tjob)/1000000, "ms")

	if *sauvegarde {
		// Sauvegarde en parallèle des différents canaux
		t1 := time.Now()
		var wg sync.WaitGroup
		wg.Add(3)
		go func() { savePNG("images/out_red.png", rImg); wg.Done() }()
		go func() { savePNG("images/out_green.png", gImg); wg.Done() }()
		go func() { savePNG("images/out_blue.png", bImg); wg.Done() }()
		wg.Wait()
		println("temps de sauvegarde : ", time.Now().Sub(t1)/1000000, "ms")
	}
	println("Terminé !")
}
