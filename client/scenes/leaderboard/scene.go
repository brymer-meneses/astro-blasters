package leaderboard

import (
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LeaderboardScene struct {
	config     *config.ClientConfig
	background *common.Background

	scrollPosition int
}

func NewLeaderboardScene(config *config.ClientConfig) *LeaderboardScene {
	return &LeaderboardScene{config: config, background: common.NewBackground(config.ScreenWidth, config.ScreenHeight)}
}

func (self *LeaderboardScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.DrawImage(self.background.Image, nil)

	// Draw the Title box and Title
	opts1 := &ebiten.DrawImageOptions{}
	imageWidth := assets.Borders.Image.Bounds().Dx()
	opts1.GeoM.Scale(25, 7)
	opts1.GeoM.Translate((float64(self.config.ScreenWidth-imageWidth)/3)+55, 30)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 0}), opts1)

	fontface := text.GoTextFace{Source: assets.MunroNarrow}
	lineSpacing := 10

	self.drawText(screen, "Leaderboard", fontface, 50, 550, 85, lineSpacing)

	// Draw the leaderboard box for the rankings
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 3}), 50, 32, 0, 150, 165, [4]float32{1, 1, 1, 1})
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 50, 32, 0, 150, 165, [4]float32{0.25, 0.25, 0.25, 1})

	self.drawText(screen, "The rankings are determined by who lasted the longest in the battle, ", fontface, 26, 555, 230, lineSpacing)
	self.drawText(screen, "with the bravest holding out until the bitter end. Well done, cadets!", fontface, 26, 555, 260, lineSpacing)

	// Rank 1
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, 290, [4]float32{0.25, 0.25, 0.25, 1})
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 7, Y: 0}), 7, 8, 0, 255, 290, [4]float32{0.8, 0.8, 0.8, 1})
	self.drawText(screen, "1", fontface, 35, 283, 322, lineSpacing)
	self.drawText(screen, "Username", fontface, 50, 440, 322, lineSpacing)
	self.drawText(screen, "00:00", fontface, 50, 760, 322, lineSpacing)

	// Rank 2
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, 360, [4]float32{0.25, 0.25, 0.25, 1})
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 8, Y: 0}), 7, 8, 0, 255, 360, [4]float32{0.9, 0.9, 0.9, 1})
	self.drawText(screen, "2", fontface, 35, 285, 392, lineSpacing)
	self.drawText(screen, "Username", fontface, 50, 440, 392, lineSpacing)
	self.drawText(screen, "00:00", fontface, 50, 760, 392, lineSpacing)

	// Rank 3
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, 430, [4]float32{0.25, 0.25, 0.25, 1})
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 7, Y: 0}), 7, 8, 0, 255, 430, [4]float32{0.8, 0.8, 0.8, 1})
	self.drawText(screen, "3", fontface, 35, 285, 462, lineSpacing)
	self.drawText(screen, "Username", fontface, 50, 440, 462, lineSpacing)
	self.drawText(screen, "00:00", fontface, 50, 760, 462, lineSpacing)

	// Rank 4
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, 500, [4]float32{0.25, 0.25, 0.25, 1})
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 8, Y: 0}), 7, 8, 0, 255, 500, [4]float32{0.9, 0.9, 0.9, 1})
	self.drawText(screen, "4", fontface, 35, 286, 532, lineSpacing)
	self.drawText(screen, "Username", fontface, 50, 440, 532, lineSpacing)
	self.drawText(screen, "00:00", fontface, 50, 760, 532, lineSpacing)

	// Rank 5
	self.drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, 570, [4]float32{0.25, 0.25, 0.25, 1})
	self.drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 7, Y: 0}), 7, 8, 0, 255, 570, [4]float32{0.8, 0.8, 0.8, 1})
	self.drawText(screen, "5", fontface, 35, 285, 602, lineSpacing)
	self.drawText(screen, "Username", fontface, 50, 440, 602, lineSpacing)
	self.drawText(screen, "00:00", fontface, 50, 760, 602, lineSpacing)
}

// Helper function to draw an image with transformations
func (self *LeaderboardScene) drawTransformedImage(screen *ebiten.Image, image *ebiten.Image, scaleX, scaleY, rotate, translateX, translateY float64, colorScale [4]float32) {
	opts := &ebiten.DrawImageOptions{}

	// Apply geometric transformations
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Rotate(rotate) // Rotation in radians
	opts.GeoM.Translate(translateX, translateY)

	// Apply color transformations using ColorScale
	if len(colorScale) == 4 { // Ensure proper length (R, G, B, A)
		opts.ColorScale.Scale(colorScale[0], colorScale[1], colorScale[2], colorScale[3])
	}

	screen.DrawImage(image, opts)
}

// Helper function to draw centered text with specified font size
func (self *LeaderboardScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

func (self *LeaderboardScene) Update(dispatcher *scenes.Dispatcher) {
	// if ebiten.IsKeyPressed(ebiten.KeyShift) {
	// 	os.Exit(0) // Exit the program with a status code of 0 (normal exit)
	// }

	// Scrolling
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		self.scrollPosition--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		self.scrollPosition++
	}
}
