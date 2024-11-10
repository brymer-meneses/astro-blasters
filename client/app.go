package client

import (
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/menu"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	config *config.ClientConfig

	sceneDispatcher *scenes.Dispatcher
	scene           scenes.Scene
}

func NewApp(config *config.ClientConfig) *App {
	scene := menu.NewMenuScene(config)
	app := &App{
		config: config,
		scene:  scene,
	}
	app.sceneDispatcher = scenes.NewDispatcher(app)
	return app
}

func (self *App) Run() error {
	ebiten.SetWindowSize(self.config.ScreenWidth, self.config.ScreenHeight)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(self)
}

func (self *App) Update() error {
	self.scene.Update(self.sceneDispatcher)
	return nil
}

func (self *App) Draw(screen *ebiten.Image) {
	self.scene.Draw(screen)
}

func (self *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return self.config.ScreenWidth, self.config.ScreenHeight
}

func (self *App) ChangeScene(scene scenes.Scene) {
	self.scene = scene
}
