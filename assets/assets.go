package assets

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Background Sprite
var Ships Sprite
var Borders Sprite
var Arrows Sprite
var Spacebar Sprite
var Healthbar Sprite
var Messagebar Sprite

var Munro *text.GoTextFaceSource
var MunroNarrow *text.GoTextFaceSource

var OrangeBulletAnimation [4]SpriteSheet
var GreenBulletAnimation [4]SpriteSheet

var BlueExplosion SpriteSheet

func init() {
	Background = mustLoadSpriteFromBytes(background, 512, 512)
	Ships = mustLoadSpriteFromBytes(ships, 8, 8)
	Borders = mustLoadSpriteFromBytes(iu, 16, 16)
	Arrows = mustLoadSpriteFromBytes(iu, 8, 8)
	Spacebar = mustLoadSpriteFromBytes(iu, 8, 4)
	Healthbar = mustLoadSpriteFromBytes(iu, 16, 8)
	Messagebar = mustLoadSpriteFromBytes(projectile, 24, 8)

	MunroNarrow = mustLoadFontFromBytes(munroNarrow)
	Munro = mustLoadFontFromBytes(munro)

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

	BlueExplosion = NewSpriteSheet(
		miscSprite,
		TileIndex{12, 6},
		TileIndex{11, 6},
		TileIndex{10, 6},
		TileIndex{9, 6},
	)
}

//go:embed SpaceShooterAssetPack/Miscellaneous.png
var miscellaneous []byte

//go:embed MunroFont/munro-narrow.ttf
var munroNarrow []byte

//go:embed MunroFont/munro.ttf
var munro []byte

//go:embed background.png
var background []byte

//go:embed SpaceShooterAssetPack/Ships.png
var ships []byte

//go:embed SpaceShooterAssetPack/IU.png
var iu []byte

//go:embed SpaceShooterAssetPack/Projectiles.png
var projectile []byte

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
