package assets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image      *ebiten.Image
	TileWidth  int
	TileHeight int
}

type SpriteTile struct {
	X0, Y0, X1, Y1 int
}

func (self *Sprite) GetTile(tile SpriteTile) *ebiten.Image {
	rect := image.Rect(
		tile.X0*self.TileWidth, tile.Y0*self.TileHeight,
		tile.X1*self.TileWidth, tile.Y1*self.TileHeight,
	)

	return self.Image.SubImage(rect).(*ebiten.Image)
}
