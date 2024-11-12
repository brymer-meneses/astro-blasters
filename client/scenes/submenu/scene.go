package submenu

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

type SubMenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	once       sync.Once
}

func NewSubMenuScene(config *config.ClientConfig) *SubMenuScene {
	return &SubMenuScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight)}
}

func (self *SubMenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	fontface := text.GoTextFace{Source: assets.FontNarrow}
	lineSpacing := 10

	// Draw the title
	self.drawText(screen, "How to Play", fontface, 100, float64(self.config.ScreenWidth)/2, 100, lineSpacing)

	self.drawText(screen, "> Press 'W' to move forward <", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-450, lineSpacing)
	self.drawText(screen, "> Press 'A' to rotate clockwise <", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-400, lineSpacing)
	self.drawText(screen, "> Press 'D' to rotate counterclockwise <", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-350, lineSpacing)
	self.drawText(screen, "> Press 'Spacebar' to shoot bullets <", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-300, lineSpacing)

	self.drawText(screen, "You are now ready to play.", fontface, 25, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-200, lineSpacing)

	self.drawText(screen, "Good luck and may the force be with you!", fontface, 25, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-175, lineSpacing)
	// Draw the subtitle
	self.drawText(screen, "Press 'S' to Start the Game", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-75, lineSpacing)
}

// Helper function to draw centered text with specified font size
func (self *SubMenuScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

func (self *SubMenuScene) Update(dispatcher *scenes.Dispatcher) {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(arena.NewArenaScene(self.config))
			})
	}
}
