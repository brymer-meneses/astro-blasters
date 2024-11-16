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

	rects := []assets.TileIndex{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
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
