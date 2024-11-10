package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dispatcher *Dispatcher)
}

type app interface {
	ChangeScene(scenes Scene)
}

type Dispatcher struct {
	app app
}

func NewDispatcher(app app) *Dispatcher {
	return &Dispatcher{app}
}

func (self *Dispatcher) Dispatch(scene Scene) {
	self.app.ChangeScene(scene)
}
