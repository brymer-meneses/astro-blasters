package menu_scene

import (
	"space-shooter/assets"
	"space-shooter/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	config       *config.AppConfig
	assetManager *assets.AssetManager
}

func NewMenuScene(config *config.AppConfig, manager *assets.AssetManager) *MenuScene {
	return &MenuScene{config, manager}
}

type FontFace struct {
	text.GoTextFace
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
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

func (self *MenuScene) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyP) {

	}
}
