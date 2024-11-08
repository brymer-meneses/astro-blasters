package menu

import (
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/game"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	config       *config.ClientConfig
	assetManager *assets.AssetManager
	once         sync.Once
}

func NewMenuScene(config *config.ClientConfig, manager *assets.AssetManager) *MenuScene {
	return &MenuScene{config: config, assetManager: manager}
}

type FontFace struct {
	text.GoTextFace
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	self.assetManager.Background.Render(screen)

	msg := "Space Shooter"

	fontface := text.GoTextFace{
		Source: self.assetManager.FontSource,
		Size:   100,
	}

	lineSpacing := 10

	width, height := text.Measure(msg, &fontface, 10)

	ops := &text.DrawOptions{}
	ops.LineSpacing = float64(lineSpacing)
	ops.GeoM.Translate(-width/2, -height/2)

	ops.GeoM.Translate(float64(self.config.ScreenWidth)/2, 100)

	text.Draw(screen, msg, &fontface, ops)
}

func (self *MenuScene) Update(dispatcher *scenes.SceneDispatcher) {
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		self.once.Do(
			func() {
				dispatcher.DispatchScene(game.NewGameScene(self.config, self.assetManager))
			})
	}
}
