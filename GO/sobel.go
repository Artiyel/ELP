package main

import (
	"image"
	"image/color"
)

func Sobel(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	out := image.NewRGBA(bounds)

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

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

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

			r := clamp(abs(sumRx)+abs(sumRy), 0, 255)
			g := clamp(abs(sumGx)+abs(sumGy), 0, 255)
			b := clamp(abs(sumBx)+abs(sumBy), 0, 255)

			out.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return out
}

// Fonction pour bloquer les valeurs rgb entre 0 et 255 si elles débordent
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
