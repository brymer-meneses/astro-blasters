package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/game/scenes"
)

type Game struct {
	assetManager *assets.AssetManager
	config       *config.AppConfig
	scene        scenes.Scene
}

func NewGame(config *config.AppConfig) Game {
	assetManager := assets.NewAssetManager(config)
	scene := scenes.NewGameScene(config, &assetManager, 0)

	return Game{
		&assetManager,
		config,
		scene,
	}
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
