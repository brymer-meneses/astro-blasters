package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(controller *AppController)
	Configure(controller *AppController) error
}

type app interface {
	ChangeScene(scenes Scene)
	ChangeMusic(data []byte)
	PlaySfx(data []byte)
}

type AppController struct {
	app app
}

func NewAppController(app app) *AppController {
	return &AppController{app}
}

func (self *AppController) ChangeScene(scene Scene) {
	self.app.ChangeScene(scene)
}

func (self *AppController) ChangeMusic(data []byte) {
	self.app.ChangeMusic(data)
}

func (self *AppController) PlaySfx(data []byte) {
	self.app.PlaySfx(data)
}
