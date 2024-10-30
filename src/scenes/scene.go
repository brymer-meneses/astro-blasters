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
	Channel chan Scene
	once    sync.Once
}

func NewSceneDispatcher() *SceneDispatcher {
	return &SceneDispatcher{
		Channel: make(chan Scene),
	}
}

func (self *SceneDispatcher) DispatchScene(scene Scene) {
	self.once.Do(func() {
		self.Channel <- scene
	})
}

func (self *SceneDispatcher) Reset() {
	self.once = sync.Once{}
}
