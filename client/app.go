package client

import (
	"astro-blasters/client/config"
	"astro-blasters/client/scenes"
	"astro-blasters/client/scenes/common/failure"
	"astro-blasters/client/scenes/menu"
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type App struct {
	config *config.ClientConfig

	controller *scenes.AppController
	scene      scenes.Scene

	player *audio.Player

	audioContext *audio.Context
}

func NewApp(config *config.ClientConfig) *App {
	app := &App{
		config:       config,
		audioContext: audio.NewContext(44100),
	}

	app.controller = scenes.NewAppController(app)
	app.controller.ChangeScene(menu.NewMenuScene(config))
	return app
}

func (self *App) Run() error {
	ebiten.SetWindowSize(self.config.ScreenWidth, self.config.ScreenHeight)
	ebiten.SetWindowTitle("Astro Blasters")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return ebiten.RunGame(self)
}

func (self *App) Update() error {
	self.scene.Update(self.controller)
	return nil
}

func (self *App) Draw(screen *ebiten.Image) {
	self.scene.Draw(screen)
}

func (self *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return self.config.ScreenWidth, self.config.ScreenHeight
}

func (self *App) ChangeScene(scene scenes.Scene) {
	if err := scene.Configure(self.controller); err != nil {
		self.scene = failure.NewFailureScene(self.config, err)
		return
	}

	self.scene = scene
}

func (self *App) ChangeMusic(data []byte) {
	if self.player != nil && self.player.IsPlaying() {
		self.player.Close()
	}

	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	// Turns the byte stream into a reader that will loop
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	self.player, err = self.audioContext.NewPlayer(loop)
	if err != nil {
		panic(err)
	}

	self.player.Play()
}

func (self *App) PlaySfx(data []byte) {
	stream, err := wav.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	player, err := self.audioContext.NewPlayer(stream)
	player.Play()
}
