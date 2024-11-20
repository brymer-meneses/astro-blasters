package assets

import "github.com/hajimehoshi/ebiten/v2"

type SpriteSheet struct {
	Sprite Sprite
	Tiles  []TileIndex
}

func (self *SpriteSheet) GetFrame(index int) *ebiten.Image {
	return self.Sprite.GetTile(self.Tiles[index])
}

func (self *SpriteSheet) TotalFrames() int {
	return len(self.Tiles)
}

func NewSpriteSheet(sprite Sprite, tiles ...TileIndex) SpriteSheet {
	return SpriteSheet{
		Tiles:  tiles,
		Sprite: sprite,
	}
}
