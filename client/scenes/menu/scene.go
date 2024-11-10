package menu

import (
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/arena"
	"space-shooter/client/scenes/common"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	once       sync.Once
}

func NewMenuScene(config *config.ClientConfig) *MenuScene {
	return &MenuScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight)}
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	msg := "Space Shooter"

	fontface := text.GoTextFace{
		Source: assets.FontNarrow,
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

func (self *MenuScene) Update(dispatcher *scenes.Dispatcher) {
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(arena.NewArenaScene(self.config))
			})
	}
}
