package app

import (
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/scenes"
	"space-shooter/scenes/menu"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	config       *config.AppConfig
	assetManager *assets.AssetManager

	sceneDispatcher *scenes.SceneDispatcher
	scene           scenes.Scene
}

func NewApp(config *config.AppConfig) *App {
	assetManager := assets.NewAssetManager(config)
	scene := menu.NewMenuScene(config, assetManager)
	return &App{
		sceneDispatcher: scenes.NewSceneDispatcher(),
		config:          config,
		scene:           scene,
		assetManager:    assetManager,
	}
}

func (self *App) Run() error {
	width, height := 1280, 720
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Handle scene dispatch.
	go self.sceneDispatcher.CheckDispatch(self)

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
	return ebiten.WindowSize()
}

func (self *App) ChangeScene(scene scenes.Scene) {
	self.scene = scene
}
