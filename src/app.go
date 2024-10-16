package app

import (
	"space-shooter/config"
	"space-shooter/game"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	config *config.AppConfig
}

func NewApp(config *config.AppConfig) App {
	return App{config}
}

func (self *App) RunApp() error {
	width, height := 1280, 720
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := game.NewGame(self.config)

	err := ebiten.RunGame(&game)

	return err
}
