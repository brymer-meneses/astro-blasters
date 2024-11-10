package common

import (
	"math/rand"
	"space-shooter/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	Image *ebiten.Image
}

func NewBackground(width int, height int) *Background {
	background := ebiten.NewImage(int(width), int(height))
	tileWidth := 128
	tileHeight := 256

	rects := []assets.SpriteTile{
		{X0: 0, Y0: 0, X1: 1, Y1: 1},
		{X0: 1, Y0: 1, X1: 2, Y1: 2},
		{X0: 2, Y0: 1, X1: 3, Y1: 3},
	}

	for x := 0; x < width; x += tileWidth {
		for y := 0; y < height; y += tileHeight {
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x), float64(y))
			rect := rects[rand.Intn(len(rects))]

			background.DrawImage(assets.Background.GetTile(rect), &opts)
		}
	}

	return &Background{background}
}
