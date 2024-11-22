package common

import (
	"space-shooter/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	Image *ebiten.Image
}

func NewBackground(width int, height int) *Background {
	background := ebiten.NewImage(int(width), int(height))
	tileWidth := 512
	tileHeight := 512

	for x := 0; x < width; x += tileWidth {
		for y := 0; y < height; y += tileHeight {
			opts := ebiten.DrawImageOptions{}
			opts.ColorScale.ScaleAlpha(0.4)
			opts.GeoM.Translate(float64(x), float64(y))

			background.DrawImage(assets.Background.GetTile(assets.TileIndex{X: 0, Y: 0}), &opts)
		}
	}

	return &Background{background}
}
