package assets

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
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

var Miscellaneous Sprite

var Munro *text.GoTextFaceSource
var MunroNarrow *text.GoTextFaceSource

var OrangeExhaustAnimation [4]SpriteSheet
var GreenExhaustAnimation [4]SpriteSheet

var Bullet *ebiten.Image

var OrangeExplosion SpriteSheet

//go:embed sfx/explosion.wav
var Explosion []byte

//go:embed sfx/laser.wav
var LaserAudio []byte

//go:embed sfx/BattleMusic.mp3
var BattleMusic []byte

//go:embed sfx/IntroMusic.mp3
var IntroMusic []byte

//go:embed sfx/hit.wav
var Hit []byte

func init() {
	projectileImage := mustLoadImageFromBytes(projectile)
	projectile := NewSprite(projectileImage, 8, 8)

	iu := mustLoadImageFromBytes(iu)
	Background = NewSprite(mustLoadImageFromBytes(background), 512, 512)
	Ships = NewSprite(mustLoadImageFromBytes(ships), 8, 8)

	Borders = NewSprite(iu, 16, 16)
	Arrows = NewSprite(iu, 8, 8)
	Spacebar = NewSprite(iu, 8, 4)
	Healthbar = NewSprite(iu, 16, 8)
	Messagebar = NewSprite(projectileImage, 24, 8)

	MunroNarrow = mustLoadFontFromBytes(munroNarrow)
	Munro = mustLoadFontFromBytes(munro)

	Miscellaneous := NewSprite(mustLoadImageFromBytes(miscellaneous), 8, 8)

	Bullet = projectile.GetTile(TileIndex{X: 3, Y: 6})

	for i := range 4 {
		OrangeExhaustAnimation[i] = NewSpriteSheet(
			Miscellaneous,
			TileIndex{5 + i, 0},
			TileIndex{5 + i, 1},
			TileIndex{5 + i, 2},
			TileIndex{5 + i, 3},
		)
		GreenExhaustAnimation[i] = NewSpriteSheet(
			Miscellaneous,
			TileIndex{9 + i, 0},
			TileIndex{9 + i, 1},
			TileIndex{9 + i, 2},
			TileIndex{9 + i, 3},
		)
	}

	OrangeExplosion = NewSpriteSheet(
		Miscellaneous,
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

func mustLoadImageFromBytes(data []byte) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return image
}

func mustLoadFontFromBytes(data []byte) *text.GoTextFaceSource {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return fontSource
}
