package scenes

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Draw(screen *ebiten.Image)
	Update()
}
