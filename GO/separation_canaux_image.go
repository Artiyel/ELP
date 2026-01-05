package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"runtime"
	"sync"
)

func SplitChannelsParallel(src image.Image, workers int) (*image.RGBA, *image.RGBA, *image.RGBA) {
	bounds := src.Bounds() //prends les dimensions de l'image

	rImg := image.NewRGBA(bounds) //crée une image par canal de couleur
	gImg := image.NewRGBA(bounds)
	bImg := image.NewRGBA(bounds)

	height := bounds.Dy()
	chunk := height / workers //defini la taille d'une ligne en fonction du nombre de worker

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ { //pour chaque worker
		yStart := bounds.Min.Y + i*chunk
		yEnd := yStart + chunk //travaille sur une plage définie

		if i == workers-1 {
			yEnd = bounds.Max.Y //si c'est le dernier worker, travaille jusqu'a la fin
		}

		wg.Add(1)
		go splitChannelsChunk( //lance splitchannelschunk avec en parametre le chunk qui a été décidé plus haut
			src, rImg, gImg, bImg,
			yStart, yEnd, &wg,
		)
	}

	wg.Wait()
	return rImg, gImg, bImg
}

func splitChannelsChunk(src image.Image, rImg, gImg, bImg *image.RGBA, yStart, yEnd int, wg *sync.WaitGroup) {
	defer wg.Done()

	bounds := src.Bounds() //reprend les bounds qui ont été mise en parametre

	for y := yStart; y < yEnd; y++ { //pour chaque lignes
		for x := bounds.Min.X; x < bounds.Max.X; x++ { //pour chaque pixel dans une ligne
			r, g, b, _ := src.At(x, y).RGBA() //prend la couleur

			// Conversion 16 bits -> 8 bits
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			rImg.Set(x, y, color.RGBA{r8, 0, 0, 255}) //met la couleur du pixel sur une image a la bonne couleur
			gImg.Set(x, y, color.RGBA{0, g8, 0, 255})
			bImg.Set(x, y, color.RGBA{0, 0, b8, 255})
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
	file, err := os.Open("images/heic1509a.jpg") //ouvre l'image
	if err != nil {
		panic(err) //error handling
	}
	defer file.Close()

	img, format, err := image.Decode(file) //decode l'image
	if err != nil {
		panic(err) //error handling
	}

	println("Format détecté :", format)

	workers := runtime.NumCPU() //workers = nb de coeurs de ton cpu

	rImg, gImg, bImg := SplitChannelsParallel(img, workers)

	// Sauvegarde
	savePNG("images/red.png", rImg)
	savePNG("images/green.png", gImg)
	savePNG("images/blue.png", bImg)
}
