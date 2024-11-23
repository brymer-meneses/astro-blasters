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

type TileIndex struct {
	X int
	Y int
}

func NewSprite(image *ebiten.Image, tileWidth int, TileHeight int) Sprite {
	return Sprite{Image: image, TileWidth: tileWidth, TileHeight: TileHeight}
}

func (self *Sprite) GetTile(tile TileIndex) *ebiten.Image {
	x0 := tile.X * self.TileWidth
	y0 := tile.Y * self.TileHeight

	x1 := (tile.X + 1) * self.TileWidth
	y1 := (tile.Y + 1) * self.TileHeight

	rect := image.Rect(x0, y0, x1, y1)

	return self.Image.SubImage(rect).(*ebiten.Image)
}
