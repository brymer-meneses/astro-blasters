package entities

import (
	"math"
	"space-shooter/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	sprite *assets.Sprite

	// The horizontal offset of the player.
	X float64
	// The vertical offset of the player.
	Y float64

	// The rotation of the player in radians
	Angle float64
}

func NewPlayer(sprite *assets.Sprite) Player {
	return Player{sprite: sprite, X: 0, Y: 0}
}

func (p *Player) MoveUp() {
	p.Y -= 5 * math.Cos(p.Angle)
	p.X += 5 * math.Sin(p.Angle)
}

func (p *Player) RotateCounterClockwise() {
	p.Angle += 5 * math.Pi / 180
}

func (p *Player) RotateClockwise() {
	p.Angle -= 5 * math.Pi / 180
}

func (p *Player) Render(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	x_0 := float64(p.sprite.Image.Bounds().Dx()) / 2
	y_0 := float64(p.sprite.Image.Bounds().Dy()) / 2

	op.GeoM.Translate(-x_0, -y_0)

	op.GeoM.Rotate(p.Angle)
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(p.X, p.Y)

	img := p.sprite.Image

	screen.DrawImage(img, op)
}
