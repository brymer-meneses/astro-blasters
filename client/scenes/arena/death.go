package arena

import (
	"image/color"
	"space-shooter/assets"
	"space-shooter/client/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	delta = 0.01
)

type DeathScene struct {
	fadeInAlpha float64
	config      *config.ClientConfig
}

func NewDeathScene(config *config.ClientConfig) *DeathScene {
	return &DeathScene{
		fadeInAlpha: 0,
		config:      config,
	}
}

func (self *DeathScene) Draw(screen *ebiten.Image) {
	self.fadeInAlpha += delta
	if self.fadeInAlpha > 1.0 {
		self.fadeInAlpha = 1.0 // Cap alpha at 1.0
	}

	{
		overlay := ebiten.NewImage(self.config.ScreenWidth, self.config.ScreenHeight)
		overlay.Fill(color.Black)

		opts := &ebiten.DrawImageOptions{}
		opts.ColorScale.ScaleAlpha(float32(self.fadeInAlpha))
		opts.GeoM.Translate(0, 0)
		screen.DrawImage(overlay, opts)
	}

	{
		font := text.GoTextFace{Source: assets.Munro, Size: 100}
		message := "You Deer"
		width, height := text.Measure(message, &font, 12)

		opts := &text.DrawOptions{}
		opts.GeoM.Translate(-width/2, -height/2)
		opts.GeoM.Translate(float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)/2)

		text.Draw(screen, message, &font, opts)
	}

	{
		font := text.GoTextFace{Source: assets.Munro, Size: 50}
		message := "you will be respawned"
		width, height := text.Measure(message, &font, 12)

		opts := &text.DrawOptions{}
		opts.GeoM.Translate(-width/2, -height/2+70)
		opts.GeoM.Translate(float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)/2)

		text.Draw(screen, message, &font, opts)
	}
}

func (self *DeathScene) Reset() {
	self.fadeInAlpha = 0
}
