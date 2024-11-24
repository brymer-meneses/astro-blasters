package failure

import (
	"os"
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FailureScene struct {
	error   error
	config  *config.ClientConfig
	ticker  *time.Ticker
	visible bool
	once    sync.Once
}

func NewFailureScene(config *config.ClientConfig, error error) *FailureScene {
	return &FailureScene{
		config:  config,
		error:   error,
		visible: true,
		ticker:  time.NewTicker(500 * time.Millisecond),
	}
}

func (self *FailureScene) Draw(screen *ebiten.Image) {
	font := text.GoTextFace{Source: assets.Munro, Size: 20}

	opts1 := &ebiten.DrawImageOptions{}
	opts1.GeoM.Scale(60, 10)
	opts1.GeoM.Translate(60, 200)

	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 3}), opts1)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), opts1)

	drawText(screen, self.error.Error(), font, 30, float64(self.config.ScreenWidth)/2, 275, 10, [4]float32{255, 255, 255, 255})
	if self.visible {
		drawText(screen, "Press C To Close the Game", font, 30, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-300, 10, [4]float32{255, 255, 255, 255})
	}
}

func (self *FailureScene) Update(controller *scenes.AppController) {
	select {
	case <-self.ticker.C:
		self.visible = !self.visible
	default:
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		self.once.Do(
			func() {
				os.Exit(0)
			})
	}
}

func (self *FailureScene) Configure(controller *scenes.AppController) error {
	return nil
}

func drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int, colorScale [4]float32) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)

	// Apply color transformations using ColorScale
	if len(colorScale) == 4 {
		opts.ColorScale.Scale(colorScale[0], colorScale[1], colorScale[2], colorScale[3])
	}

	text.Draw(screen, msg, &fontface, opts)
}
