package assets

import (
	"bytes"
	"log"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed MunroFont/munro-narrow.ttf
var munroNarrow []byte

//go:embed SpaceShooterAssetPack/BackGrounds.png
var backgroundsAsset []byte

//go:embed SpaceShooterAssetPack/Ships.png
var shipsAsset []byte

const (
	MapWidth  = 4000
	MapHeight = 4000
)

type AssetManager struct {
	Background Sprite
	FontSource *text.GoTextFaceSource

	Ships []Sprite
}

func NewAssetManager() *AssetManager {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(munroNarrow))
	if err != nil {
		log.Fatal(err)
	}

	background := NewSpriteBuilder().
		FromBytes(backgroundsAsset).
		CreateTiles(CreateTilesInput{X_start: 0, Y_start: 0, Width: 128, Height: 256, X_count: 3, Y_count: 2}).
		FilterTiles(
			Tile{X: 128, Y: 0, Width: 128, Height: 256},
			Tile{X: 256, Y: 0, Width: 128, Height: 256},
			Tile{X: 0, Y: 256, Width: 128, Height: 256},
		).
		BuildAsBackgroundSprite(MapWidth, MapHeight, 128, 256)

	ships := make([]Sprite, 5)
	for i := 0; i < 5; i += 1 {
		ships[i] = NewSpriteBuilder().
			FromBytes(shipsAsset).
			BuildAsSprite(Tile{X: 8, Y: 8 * i, Width: 8, Height: 8})

	}

	return &AssetManager{background, fontSource, ships}
}
