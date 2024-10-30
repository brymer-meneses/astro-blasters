package assets

import (
	"bufio"
	"log"
	"os"
	"space-shooter/config"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type AssetManager struct {
	Background Sprite
	FontSource *text.GoTextFaceSource

	Ships []Sprite
}

func loadFont(filename string) (*text.GoTextFaceSource, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	source, err := text.NewGoTextFaceSource(reader)

	return source, nil
}

func NewAssetManager(config *config.AppConfig) *AssetManager {

	fontSource, err := loadFont("../assets/MunroFont/munro-narrow.ttf")
	if err != nil {
		log.Fatal(err)
	}

	background := NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/BackGrounds.png").
		CreateTiles(CreateTilesInput{X_start: 0, Y_start: 0, Width: 128, Height: 256, X_count: 3, Y_count: 2}).
		FilterTiles(
			Tile{X: 128, Y: 0, Width: 128, Height: 256},
			Tile{X: 256, Y: 0, Width: 128, Height: 256},
			Tile{X: 0, Y: 256, Width: 128, Height: 256},
		).
		BuildAsBackgroundSprite(config.ScreenWidth, config.ScreenHeight, 128, 256)

	ships := make([]Sprite, 5)
	for i := 0; i < 5; i += 1 {
		ships[i] = NewSpriteBuilder().
			FromFile("../assets/SpaceShooterAssetPack/Ships.png").
			BuildAsSprite(Tile{X: 8, Y: 8 * i, Width: 8, Height: 8})

	}

	return &AssetManager{background, fontSource, ships}
}
