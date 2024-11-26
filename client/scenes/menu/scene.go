package menu

import (
	"astro-blasters/assets"
	"astro-blasters/client/config"
	"astro-blasters/client/scenes"
	"astro-blasters/client/scenes/common"
	"astro-blasters/client/scenes/submenu"

	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	once       sync.Once
	visible    bool
	ticker     *time.Ticker
}

func NewMenuScene(config *config.ClientConfig) *MenuScene {
	return &MenuScene{
		config:     config,
		background: common.NewBackground(config.ScreenWidth, config.ScreenHeight),
		ticker:     time.NewTicker(500 * time.Millisecond),
		visible:    true,
	}
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	imageWidth := assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 0}).Bounds().Dx()
	imageHeight := assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 0}).Bounds().Dy()

	centerX := (float64(self.config.ScreenWidth-imageWidth) / 4) + 30
	centerY := float64(self.config.ScreenHeight-imageHeight) / 4

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(30, 15)
	opts.GeoM.Translate(centerX, centerY)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 0}), opts)

	fontface := text.GoTextFace{Source: assets.MunroNarrow}
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

func (self *MenuScene) Update(controller *scenes.AppController) {
	// Toggle visibility every tick
	select {
	case <-self.ticker.C:
		self.visible = !self.visible
	default:
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		self.once.Do(
			func() {
				controller.ChangeScene(submenu.NewSubMenuScene(self.config))
			})
	}
}

func (self *MenuScene) Configure(controller *scenes.AppController) error {
	controller.ChangeMusic(assets.IntroMusic)
	return nil
}
