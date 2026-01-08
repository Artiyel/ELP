package main

import (
	_ "image/jpeg"
)

// Job correspond à un bloc de lignes à traiter
type Job struct {
	yStart, yEnd int
}
