package component

import (
	"space-shooter/assets"

	"github.com/yohamta/donburi"
)

type AnimationData struct {
	sheet []assets.Sprite
	frame uint16
}

func (self *AnimationData) Frame() assets.Sprite {
	return self.sheet[self.frame]
}

func (self *AnimationData) Update() {
	self.frame = (self.frame + 1) % uint16(len(self.sheet))
}

var Animation = donburi.NewComponentType[AnimationData]()
