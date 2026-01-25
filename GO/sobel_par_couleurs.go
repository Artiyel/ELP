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
	srcRGBA := image.NewRGBA(bounds)                      // Création d’une image RGBA pour accès direct à la mémoire (plus rapide)
	draw.Draw(srcRGBA, bounds, img, bounds.Min, draw.Src) // Copie de l’image décodée dans l’image RGBA
	outRGBA := image.NewRGBA(bounds)                      // Image de sortie

	t0 := time.Now()
	// Traitement Sobel parallèle via worker pool
	SobelWorkerPool(srcRGBA, outRGBA, *numWorkers)
	t1 := time.Now()
	println("Temps de calcul    : ", t1.Sub(t0))
	if *sauvegarde {
		// Sauvegarde
		savePNG("images/out_couleur.png", outRGBA)
		t2 := time.Now()
		println("Temps de sauvegarde :", t2.Sub(t1))
	}
	println("Traitement terminé")

}

// SobelWorkerPool lance un pool de workers pour appliquer le filtre Sobel
func SobelWorkerPool(src *image.RGBA, out *image.RGBA, numWorkers int) {
	bounds := src.Bounds()              // Dimensions de l’image
	jobs := make(chan Job, bounds.Dy()) // Canal de jobs (1 job = 1 ligne). File d'attente.
	var wg sync.WaitGroup               // Synchronisation des jobs

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		go sobelWorker(src, out, jobs, &wg)
	}

	// Chaque ligne de l'image devient un job (on évite les bords y=0 et y=max)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		wg.Add(1)         // Un job en plus à attendre
		jobs <- Job{y: y} // Envoi du job
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
			var rGx, gGx, bGx int
			var rGy, gGy, bGy int

			// Convolution 3x3 : le pixel central est comparé à ses voisins
			for j := -1; j <= 1; j++ {
				for i := -1; i <= 1; i++ {
					pixelIdx := (y+j)*stride + (x+i)*4 // Calcul de l’index mémoire du pixel voisin
					// Lecture des composantes RGB
					r := int(src.Pix[pixelIdx])
					g := int(src.Pix[pixelIdx+1])
					b := int(src.Pix[pixelIdx+2])

					// Poids Sobel
					wx := gxMask[j+1][i+1]
					wy := gyMask[j+1][i+1]
					// Accumulation Gx
					rGx += r * wx
					gGx += g * wx
					bGx += b * wx
					// Accumulation Gy
					rGy += r * wy
					gGy += g * wy
					bGy += b * wy
				}
			}

			// Calcul du gradient/magnitude (Approximation rapide : |Gx| + |Gy|)
			// Rouge
			absRGx, absRGy := rGx, rGy
			if absRGx < 0 {
				absRGx = -absRGx
			}
			if absRGy < 0 {
				absRGy = -absRGy
			}
			magR := absRGx + absRGy // Plus |Gx| + |Gy| est grand, plus il y a un bord

			// Vert
			absGGx, absGGy := gGx, gGy
			if absGGx < 0 {
				absGGx = -absGGx
			}
			if absGGy < 0 {
				absGGy = -absGGy
			}
			magG := absGGx + absGGy // Plus |Gx| + |Gy| est grand, plus il y a un bord

			// Bleu
			absBGx, absBGy := bGx, bGy
			if absBGx < 0 {
				absBGx = -absBGx
			}
			if absBGy < 0 {
				absBGy = -absBGy
			}
			magB := absBGx + absBGy // Plus |Gx| + |Gy| est grand, plus il y a un bord

			// Écriture directe et saturation à 255
			outIdx := y*stride + x*4 // Index du pixel de sortie
			if magR > 255 {
				magR = 255
			}
			if magG > 255 {
				magG = 255
			}
			if magB > 255 {
				magB = 255
			}
			// Écriture du pixel résultat
			out.Pix[outIdx] = uint8(magR)
			out.Pix[outIdx+1] = uint8(magG)
			out.Pix[outIdx+2] = uint8(magB)
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
