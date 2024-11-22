package submenu

import (
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/client/scenes/starter"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type SubMenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	once       sync.Once
	visible    bool
	ticker     *time.Ticker
}

func NewSubMenuScene(config *config.ClientConfig) *SubMenuScene {
	return &SubMenuScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight), visible: true, ticker: time.NewTicker(500 * time.Millisecond)}
}

func (self *SubMenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	// Draw the Title box and Title
	opts1 := &ebiten.DrawImageOptions{}
	imageWidth := assets.Borders.Image.Bounds().Dx()
	opts1.GeoM.Scale(25, 7)
	opts1.GeoM.Translate((float64(self.config.ScreenWidth-imageWidth)/3)+55, 50)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 3}), opts1)

	fontface := text.GoTextFace{Source: assets.FontNarrow}
	lineSpacing := 10

	self.drawText(screen, "Welcome Cadet!", fontface, 50, 550, 105, lineSpacing)

	// Draw the Instructions box for the controls
	// Borders and Arrows
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 3}), 60, 31, 0, 50, 185)
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 60, 31, 0, 50, 185)

	// Text blocks
	self.drawText(screen, "The galaxy needs a hero and YOU are our last hope. Blast enemy ships and protect the", fontface, 26, 530, 255, lineSpacing)
	self.drawText(screen, "fate of the stars! Before you take-off, here's your mission briefing on the controls.", fontface, 26, 530, 285, lineSpacing)
	self.drawText(screen, "Use the arrow keys to navigate your ship—learn them well and may your aim be true!", fontface, 26, 530, 315, lineSpacing)

	// Right arrow instruction
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 3.5, 3.5, 0, 200, 350)
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 6, Y: 10}), 5, 5, 0, 206, 360)
	self.drawText(screen, "Press the right arrow key to rotate clockwise", fontface, 30, 510, 380, lineSpacing)

	// Up arrow instruction
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 3.5, 3.5, 0, 200, 410)
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 6, Y: 10}), 5, 5, -1.5708, 210, 460)
	self.drawText(screen, "Press the up arrow key to move forward", fontface, 30, 485, 440, lineSpacing)

	// Left arrow instruction
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 3.5, 3.5, 0, 200, 470)
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 6, Y: 10}), 5, 5, 3.14159, 248, 516)
	self.drawText(screen, "Press the left arrow key to rotate counterclockwise", fontface, 30, 543, 500, lineSpacing)

	// Spacebar instruction
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 3.5, 3.5, 0, 200, 530)
	self.drawTransformedImage(screen, assets.Bar.GetTile(assets.TileIndex{X: 5, Y: 24}), 5, 5, 0, 208, 555)
	self.drawText(screen, "Press the spacebar to shoot bullets", fontface, 30, 458, 560, lineSpacing)

	// Draw subtext
	if self.visible {
		self.drawText(screen, "Press P To Proceed", fontface, 40, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-110, lineSpacing)
	}
}

// Helper function to draw an image with transformations
func (self *SubMenuScene) drawTransformedImage(screen *ebiten.Image, image *ebiten.Image, scaleX, scaleY, rotate, translateX, translateY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Rotate(rotate) // Rotation in radians; use 0 for no rotation
	opts.GeoM.Translate(translateX, translateY)
	screen.DrawImage(image, opts)
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
	// Toggle visibility every tick
	select {
	case <-self.ticker.C:
		self.visible = !self.visible
	default:
	}

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(starter.NewStarterScene(self.config))
			})
	}
}
