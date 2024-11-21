package menu

import (
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/client/scenes/submenu"

	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	border     *common.Border
	once       sync.Once
	visible    bool
	ticker     *time.Ticker
}

func NewMenuScene(config *config.ClientConfig) *MenuScene {
	return &MenuScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight), border: common.NewBorder(16, 16), visible: true, ticker: time.NewTicker(500 * time.Millisecond)}
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	imageWidth := self.border.Image.Bounds().Dx()
	imageHeight := self.border.Image.Bounds().Dy()
	centerX := (float64(self.config.ScreenWidth-imageWidth) / 4) + 30
	centerY := float64(self.config.ScreenHeight-imageHeight) / 4

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(30, 15)
	opts.GeoM.Translate(centerX, centerY)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 0}), opts)

	fontface := text.GoTextFace{Source: assets.FontNarrow}
	lineSpacing := 10

	// Draw the title
	self.drawText(screen, "Astro", fontface, 80, float64(self.config.ScreenWidth)/2, 260, lineSpacing)
	self.drawText(screen, "Blasters", fontface, 80, float64(self.config.ScreenWidth)/2, 330, lineSpacing)

	// Draw subtext
	if self.visible {
		self.drawText(screen, "Press S To Start the Game", fontface, 40, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-100, lineSpacing)
	}
}

// Helper function to draw centered text with specified font size
func (self *MenuScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x float64, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

// func (self *MenuScene) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return 320, 240
// }

func (self *MenuScene) Update(dispatcher *scenes.Dispatcher) {
	// Toggle visibility every tick
	select {
	case <-self.ticker.C:
		self.visible = !self.visible
	default:
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(submenu.NewSubMenuScene(self.config))
			})
	}
}
