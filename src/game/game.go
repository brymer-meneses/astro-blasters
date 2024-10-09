package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"space-shooter/assets"
	"space-shooter/config"
)

type Game struct {
	stateManager SceneManager
	assetManager *assets.AssetManager
	config       *config.AppConfig
}

func NewGame(screenWidth, screenHeight int) Game {
	config := config.AppConfig{ScreenHeight: screenHeight, ScreenWidth: screenWidth}

	assetManager := assets.NewAssetManager(&config)
	stateManager := NewStateManager(&config, &assetManager)

	return Game{
		stateManager,
		&assetManager,
		&config,
	}
}

func (g *Game) Update() error {
	g.stateManager.HandleUpdate()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.stateManager.Render(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
