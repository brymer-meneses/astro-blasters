package starter

import (
	"fmt"
	"image/color"
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/arena"
	"space-shooter/client/scenes/common"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type StarterScene struct {
	config        *config.ClientConfig
	background    *common.Background
	once          sync.Once
	inputText     string
	isFocused     bool
	visible       bool
	ticker        *time.Ticker
	cursorVisible bool
	cursorTimer   time.Duration
}

func NewStarterScene(config *config.ClientConfig) *StarterScene {
	return &StarterScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight), visible: true, ticker: time.NewTicker(500 * time.Millisecond)}
}

func (self *StarterScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	fontface := text.GoTextFace{Source: assets.FontNarrow}
	lineSpacing := 10

	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 3}), 60, 15, 0, 50, 185)
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 60, 15, 0, 50, 185)

	// Text blocks
	self.drawText(screen, "Before we take off, cadet, what should we call the brave soul leading this mission?", fontface, 27, 530, 245, lineSpacing)
	self.drawText(screen, "Press 'Enter' to type in your username.", fontface, 27, 530, 280, lineSpacing)

	// Draw Input box for Username
	// opts0 := &ebiten.DrawImageOptions{}
	// imageWidth := self.border.Image.Bounds().Dx()
	// imageHeight := self.border.Image.Bounds().Dx()
	// opts0.GeoM.Scale(30, 7)
	// opts0.GeoM.Translate(float64(self.config.ScreenWidth-imageWidth)/3, float64(self.config.ScreenHeight-imageHeight)/3)
	// screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 3}), opts0)

	// opts1 := &ebiten.DrawImageOptions{}
	// opts1.GeoM.Scale(30, 7)
	// opts1.GeoM.Translate(float64(self.config.ScreenWidth-imageWidth)/3, float64(self.config.ScreenHeight-imageHeight)/3)
	// screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), opts1)

	self.drawText(screen, fmt.Sprintf("> %s", self.inputText), fontface, 30, 530, 330, lineSpacing)

	self.RenderCursor(screen)

	// Draw subtext
	if self.visible {
		self.drawText(screen, "Press Esc To Play the Game", fontface, 40, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-250, lineSpacing)
	}
}

func (self *StarterScene) drawTransformedImage(screen *ebiten.Image, image *ebiten.Image, scaleX, scaleY, rotate, translateX, translateY float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Rotate(rotate) // Rotation in radians; use 0 for no rotation
	opts.GeoM.Translate(translateX, translateY)
	screen.DrawImage(image, opts)
}

// Helper function to draw centered text with specified font size
func (self *StarterScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

// Render the input text and cursor
func (self *StarterScene) RenderCursor(screen *ebiten.Image) {

	// Render blinking cursor (if visible)
	if self.cursorVisible {
		// Cursor placement: calculate X position based on input text length
		cursorX := 535 + len(self.inputText)*5 // Adjust based on font width
		cursorY := 350                         // Cursor vertical position
		cursorImage := ebiten.NewImage(10, 2)  // Create a 2x20 rectangle for the cursor
		cursorImage.Fill(color.White)          // Set cursor color to white

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cursorX), float64(cursorY))
		screen.DrawImage(cursorImage, op) // Draw the cursor on the screen
	}
}

func (self *StarterScene) Update(dispatcher *scenes.Dispatcher) {
	// Toggle focus when the enter is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		self.isFocused = !self.isFocused
		self.cursorVisible = !self.cursorVisible
	}

	// Only handle input if the input box is focused
	if self.isFocused {
		// Update cursor blink timer
		self.cursorTimer += time.Second / time.Duration(ebiten.ActualTPS()) // Use ActualTPS
		if self.cursorTimer > time.Second/2 {
			self.cursorVisible = !self.cursorVisible
			self.cursorTimer = 0
		}

		chars := make([]rune, 0)
		for _, r := range ebiten.AppendInputChars(chars) {
			self.inputText += string(r)
		}

		// Handle backspace to remove last character
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(self.inputText) > 0 {
			self.inputText = self.inputText[:len(self.inputText)-1]
		}
	}

	// Toggle visibility every tick
	select {
	case <-self.ticker.C:
		self.visible = !self.visible
	default:
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(arena.NewArenaScene(self.config))
			})
	}
}
