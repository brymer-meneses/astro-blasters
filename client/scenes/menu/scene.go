package menu

import (
	"bytes"
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/client/scenes/submenu"
	"sync"

	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MenuScene struct {
	config     *config.ClientConfig
	background *common.Background
	once       sync.Once
	inputText  string
	isFocused  bool
}

func NewMenuScene(config *config.ClientConfig) *MenuScene {
	// Initialize the audio context and player
	audioContext := audio.NewContext(1000)
	audioData, err := os.ReadFile("assets/sfx/menu-bg.mp3")

	if err != nil {
		log.Fatal("Failed to load audio file:", err)
	}

	// Load and start the audio player
	stream := audio.NewInfiniteLoop(bytes.NewReader(audioData), int64(len(audioData)))
	audioPlayer, err := audioContext.NewPlayer(stream)

	if err != nil {
		log.Fatal("Failed to create audio player:", err)
	}

	audioPlayer.Play()

	return &MenuScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight)}
}

func (self *MenuScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	fontface := text.GoTextFace{Source: assets.FontNarrow}
	lineSpacing := 10

	// Draw the title
	self.drawText(screen, "Cosmic Clash", fontface, 100, float64(self.config.ScreenWidth)/2, 100, lineSpacing)

	// Draw the subtitle
	self.drawText(screen, "Welcome Cadet!", fontface, 50, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-165, lineSpacing)

	boxWidth, boxHeight := 400.0, 50.0
	boxColor := color.RGBA{189, 195, 199, 1}
	boxX := (float64(self.config.ScreenWidth) - boxWidth) / 2
	boxY := float64(self.config.ScreenHeight) - 135

	// Draw the input box
	vector.DrawFilledRect(screen, float32(boxX), float32(boxY), float32(boxWidth), float32(boxHeight), boxColor, false)

	textX, textY := int(boxX)+10, int(boxY)+15
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Press Spacebar to Type your Name: %s", self.inputText), textX, textY)

	// Draw final text
	self.drawText(screen, "Press Enter to Proceed", fontface, 25, float64(self.config.ScreenWidth)/2, float64(self.config.ScreenHeight)-55, lineSpacing)
}

// Helper function to draw centered text with specified font size
func (self *MenuScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

func (self *MenuScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func (self *MenuScene) Update(dispatcher *scenes.Dispatcher) {
	// Toggle focus when the spacebar is pressed (or customize to any other key)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		self.isFocused = !self.isFocused
	}

	// Only handle input if the input box is focused
	if self.isFocused {
		chars := make([]rune, 0)
		for _, r := range ebiten.AppendInputChars(chars) {
			self.inputText += string(r)
		}

		// Handle backspace to remove last character
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(self.inputText) > 0 {
			self.inputText = self.inputText[:len(self.inputText)-1]
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(submenu.NewSubMenuScene(self.config))
			})
	}
}
