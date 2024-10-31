package scenes

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update(dispatcher *SceneDispatcher)
}

type SceneDispatcher struct {
	channel chan Scene
	once    sync.Once
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
	self.once.Do(func() {
		self.channel <- scene
	})
}

func (self *SceneDispatcher) Reset() {
	self.once = sync.Once{}
}
