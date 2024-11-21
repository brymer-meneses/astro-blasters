package common

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Bar struct {
	Image *ebiten.Image
}

func NewBar(width int, height int) *Bar {
	bar := ebiten.NewImage(int(width), int(height))

	// Create a new Bar instance
	r := &Bar{
		Image: bar,
	}

	return r
}
