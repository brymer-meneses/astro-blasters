package state

import (
	"bufio"
	"log"
	"os"
	"space-shooter/game/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// The state for the game.
type MenuState struct {
	background sprite.Sprite
	fontSource *text.GoTextFaceSource

	screenWidth  int
	screenHeight int
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

func NewMenuState(screen_width, screen_height int) MenuState {
	background := sprite.NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/BackGrounds.png").
		CreateTiles(sprite.CreateTilesInput{X_start: 0, Y_start: 0, Width: 128, Height: 256, X_count: 3, Y_count: 2}).
		FilterTiles(
			sprite.Tile{X: 128, Y: 0, Width: 128, Height: 256},
			sprite.Tile{X: 256, Y: 0, Width: 128, Height: 256},
			sprite.Tile{X: 0, Y: 256, Width: 128, Height: 256},
		).
		BuildAsBackgroundSprite(screen_width, screen_height, 128, 256)

	fontSource, err := loadFont("../assets/MunroFont/munro-narrow.ttf")
	if err != nil {
		log.Fatal(err)
	}

	return MenuState{background, fontSource, screen_width, screen_height}
}

type FontFace struct {
	text.GoTextFace
}

func (ms *MenuState) Render(screen *ebiten.Image) {
	ms.background.Render(screen)

	msg := "Space Shooter"

	fontface := text.GoTextFace{
		Source: ms.fontSource,
		Size:   100,
	}

	lineSpacing := 10

	width, height := text.Measure(msg, &fontface, 10)

	ops := &text.DrawOptions{}
	ops.LineSpacing = float64(lineSpacing)
	ops.GeoM.Translate(-width/2, -height/2)
	ops.GeoM.Translate(float64(ms.screenWidth)/2, 100)

	text.Draw(screen, msg, &fontface, ops)

}

func (m *MenuState) HandleUpdate() {}
