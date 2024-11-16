package assets

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Background Sprite
var Ships Sprite

var FontNarrow *text.GoTextFaceSource

var OrangeBulletAnimation [4]SpriteSheet
var GreenBulletAnimation [4]SpriteSheet

func init() {
	Background = mustLoadSpriteFromBytes(backgrounds, 128, 256)
	Ships = mustLoadSpriteFromBytes(ships, 8, 8)
	FontNarrow = mustLoadFontFromBytes(munroNarrow)

	miscSprite := mustLoadSpriteFromBytes(miscellaneous, 8, 8)
	for i := range 4 {
		OrangeBulletAnimation[i] = NewSpriteSheet(
			miscSprite,
			TileIndex{5 + i, 0},
			TileIndex{5 + i, 1},
			TileIndex{5 + i, 2},
			TileIndex{5 + i, 3},
		)
		GreenBulletAnimation[i] = NewSpriteSheet(
			miscSprite,
			TileIndex{9 + i, 0},
			TileIndex{9 + i, 1},
			TileIndex{9 + i, 2},
			TileIndex{9 + i, 3},
		)
	}
}

//go:embed SpaceShooterAssetPack/Miscellaneous.png
var miscellaneous []byte

//go:embed MunroFont/munro-narrow.ttf
var munroNarrow []byte

//go:embed SpaceShooterAssetPack/BackGrounds.png
var backgrounds []byte

//go:embed SpaceShooterAssetPack/Ships.png
var ships []byte

func mustLoadSpriteFromBytes(data []byte, width, height int) Sprite {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return Sprite{image, width, height}
}

func mustLoadFontFromBytes(data []byte) *text.GoTextFaceSource {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return fontSource
}
