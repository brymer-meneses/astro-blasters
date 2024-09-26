package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FrameLeft = iota
	FrameDefault
	FrameRight
)

type Player struct {
	sheet        SpriteSheet
	currentFrame int

	// The horizontal offset of the player.
	X float64
	// The vertical offset of the player.
	Y float64

	// The rotation of the player in radians
	Angle float64
}

func NewPlayer(spriteSheet SpriteSheet, x float64, y float64) Player {
	return Player{sheet: spriteSheet, X: x, Y: y, currentFrame: FrameDefault}
}

func (p *Player) MoveUp() {
	p.Y -= 5 * math.Cos(p.Angle)
	p.X += 5 * math.Sin(p.Angle)
	p.currentFrame = FrameDefault
}

func (p *Player) RotateCounterClockwise() {
	p.Angle += 5 * math.Pi / 180
	p.currentFrame = FrameDefault
}

func (p *Player) RotateClockwise() {
	p.Angle -= 5 * math.Pi / 180
	p.currentFrame = FrameDefault
}

func (p *Player) Render(screen *ebiten.Image) {
	// Define draw options
	op := &ebiten.DrawImageOptions{}

	// The offset (x_0, y_0) for rotation (e.g., the center of the sprite)
	x_0 := float64(p.sheet.sprites[p.currentFrame].image.Bounds().Dx()) / 2
	y_0 := float64(p.sheet.sprites[p.currentFrame].image.Bounds().Dy()) / 2

	// Step 1: Translate to offset (negative offset)
	op.GeoM.Translate(-x_0, -y_0)

	op.GeoM.Rotate(p.Angle)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(p.X, p.Y)

	img := p.sheet.sprites[p.currentFrame].image

	// Draw the image onto the screen
	screen.DrawImage(img, op)
}
