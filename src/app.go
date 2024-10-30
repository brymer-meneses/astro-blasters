package app

import (
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/scenes"
	"space-shooter/scenes/game"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	config       *config.AppConfig
	assetManager *assets.AssetManager
	scene        scenes.Scene
}

func NewApp(config *config.AppConfig) App {
	assetManager := assets.NewAssetManager(config)
	scene := game.NewGameScene(config, assetManager)
	return App{config, assetManager, scene}
}

func (self *App) Run() error {
	width, height := 1280, 720
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(self); err != nil {
		return err
	}

	return nil
}

func (self *App) Update() error {
	self.scene.Update()
	return nil
}

func (self *App) Draw(screen *ebiten.Image) {
	self.scene.Draw(screen)
}

func (self *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
