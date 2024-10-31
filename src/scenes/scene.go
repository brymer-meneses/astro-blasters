package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dispatcher *SceneDispatcher)
}

type SceneDispatcher struct {
	channel chan Scene
}

func NewSceneDispatcher() *SceneDispatcher {
	return &SceneDispatcher{
		channel: make(chan Scene, 1),
	}
}

type app interface {
	ChangeScene(scene Scene)
}

func (self *SceneDispatcher) CheckDispatch(app app) {
	for {
		select {
		case scene := <-self.channel:
			app.ChangeScene(scene)
		}
	}
}

func (self *SceneDispatcher) DispatchScene(scene Scene) {
	self.channel <- scene
}
