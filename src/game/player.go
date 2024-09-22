package game

import (
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
	X int
	// The vertical offset of the player.
	Y int

	// The rotation of the player in radians
	Angle float64
}

func NewPlayer(spriteSheet SpriteSheet, x int, y int) Player {
	return Player{sheet: spriteSheet, X: x, Y: y, currentFrame: FrameDefault}
}

func (p *Player) MoveUp() {
	p.Y -= 5
	p.currentFrame = FrameDefault
}

func (p *Player) MoveDown() {
	p.Y += 5
	p.currentFrame = FrameDefault
}

func (p *Player) MoveLeft() {
	p.X -= 5
	p.currentFrame = FrameLeft
}

func (p *Player) MoveRight() {
	p.X += 5
	p.currentFrame = FrameRight
}
func (p *Player) Render(screen *ebiten.Image) {
	// Define draw options
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	op.GeoM.Translate(float64(p.X), float64(p.Y))

	img := p.sheet.sprites[p.currentFrame].image

	// Draw the image onto the screen
	screen.DrawImage(img, op)
}
