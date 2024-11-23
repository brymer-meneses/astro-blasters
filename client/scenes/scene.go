package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(controller *AppController)
	Configure(controller *AppController)
}

type app interface {
	ChangeScene(scenes Scene)
	ChangeBackgroundMusic(data []byte)
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

func (self *AppController) ChangeBackgroundMusic(data []byte) {
	self.app.ChangeBackgroundMusic(data)
}
