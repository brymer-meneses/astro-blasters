package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dispatcher *SceneDispatcher)
}

type app interface {
	ChangeScene(scenes Scene)
}

type SceneDispatcher struct {
	app app
}

func NewSceneDispatcher(app app) *SceneDispatcher {
	return &SceneDispatcher{app}
}

func (self *SceneDispatcher) DispatchScene(scene Scene) {
	self.app.ChangeScene(scene)
}
