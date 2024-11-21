package common

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Arrow struct {
	Image *ebiten.Image
}

func NewArrow(width int, height int) *Arrow {
	arrow := ebiten.NewImage(int(width), int(height))

	// Create a new Arrow instance
	a := &Arrow{
		Image: arrow,
	}

	return a
}
