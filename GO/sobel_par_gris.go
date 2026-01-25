package main

import (
	"flag"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
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
	// Nombre de workers (par défaut = nombre de cœurs CPU)
	numWorkers := flag.Int("n", runtime.NumCPU(), "Nombre de workers")
	// Active ou non la sauvegarde du résultat
	sauvegarde := flag.Bool("s", true, "Active ou non la sauvegarde du fichier résultant")
	// Chemin du fichier image à traiter
	path := flag.String("f", "images/Marmite_UMARY_WM-94.jpg", "Précise l'emplacement de l'image a traiter (Présente dans un dossier 'images'")
	flag.Parse() // Lecture des arguments

	// Charger l'image
	file, err := os.Open(*path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Décodage de l’image
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()                                // Récupération des dimensions de l’image
	srcRGBA := image.NewRGBA(bounds)                      // Création d’une image RGBA pour accès direct à la mémoire
	draw.Draw(srcRGBA, bounds, img, bounds.Min, draw.Src) // Copie de l’image décodée dans l’image RGBA
	outRGBA := image.NewRGBA(bounds)                      // Image de sortie

	t0 := time.Now()
	// Traitement Sobel parallèle via worker pool
	SobelWorkerPool(srcRGBA, outRGBA, *numWorkers)
	t1 := time.Now()
	println("Temps de calcul    : ", t1.Sub(t0))
	if *sauvegarde {
		// Sauvegarde
		savePNG("images/out_gris.png", outRGBA)
		t2 := time.Now()
		println("Temps de sauvegarde :", t2.Sub(t1))
	}
	println("Traitement terminé")

}

// SobelWorkerPool lance un pool de workers pour appliquer le filtre Sobel
func SobelWorkerPool(src *image.RGBA, out *image.RGBA, numWorkers int) {
	bounds := src.Bounds()              // Dimensions de l’image
	jobs := make(chan Job, bounds.Dy()) // Canal de jobs (1 job = 1 ligne)
	var wg sync.WaitGroup               // Synchronisation des jobs

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		go sobelWorker(src, out, jobs, &wg)
	}

	// Envoyer chaque ligne comme un job (on évite les bords y=0 et y=max)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		wg.Add(1)
		jobs <- Job{y: y}
	}

	close(jobs) // Plus aucun job à envoyer
	wg.Wait()   // Attente de la fin de tous les jobs
}

// sobelWorker applique le filtre Sobel sur les lignes reçues
func sobelWorker(src *image.RGBA, out *image.RGBA, jobs <-chan Job, wg *sync.WaitGroup) {
	bounds := src.Bounds() // Dimensions de l’image
	stride := src.Stride   // Nombre d’octets par ligne

	// Matrices Sobel horizontal et vertical
	gxMask := [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	gyMask := [3][3]int{{1, 2, 1}, {0, 0, 0}, {-1, -2, -1}}

	// Traitement des jobs reçus
	for job := range jobs {
		y := job.y // Ligne à traiter
		// Parcours des pixels de la ligne (hors bords)
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var Gx, Gy float64 // Gradients horizontal et vertical du pixel courant

			// Convolution 3x3 : le pixel central est comparé à ses voisins
			for j := -1; j <= 1; j++ { // Décalage vertical
				for i := -1; i <= 1; i++ { // Décalage horizontal
					pixelIdx := (y+j-bounds.Min.Y)*stride + (x+i-bounds.Min.X)*4 // Calcul de l’index mémoire du pixel voisin
					// Lecture des composantes RGB du pixel voisin
					r := float64(src.Pix[pixelIdx])
					g := float64(src.Pix[pixelIdx+1])
					b := float64(src.Pix[pixelIdx+2])
					// Conversion du pixel couleur en luminance (niveau de gris)
					lum := 0.299*r + 0.587*g + 0.114*b
					// Accumulation des gradients selon X et Y
					Gx += lum * float64(gxMask[j+1][i+1])
					Gy += lum * float64(gyMask[j+1][i+1])
				}
			}

			// Calcul manuel de la magnitude (Approximation rapide : |Gx| + |Gy|)
			mag := math.Sqrt(Gx*Gx + Gy*Gy)
			// Saturation de la valeur pour rester dans [0, 255]
			if mag > 255 {
				mag = 255
			}
			val := uint8(mag) // Conversion en valeur entière 8 bits

			outIdx := (y-bounds.Min.Y)*stride + (x-bounds.Min.X)*4 // Index mémoire du pixel de sortie
			// Écriture du pixel résultat
			out.Pix[outIdx] = val
			out.Pix[outIdx+1] = val
			out.Pix[outIdx+2] = val
			out.Pix[outIdx+3] = 255
		}
		wg.Done() // Signal de fin du job
	}
}

// savePNG enregistre l'image avec le nom donné
func savePNG(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
