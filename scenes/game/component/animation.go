package component

import "github.com/yohamta/donburi"

type AnimationData struct {
	sheet []SpriteData
	frame uint16
}

func (self *AnimationData) Frame() *SpriteData {
	return &self.sheet[self.frame]
}

func (self *AnimationData) Update() {
	self.frame = (self.frame + 1) % uint16(len(self.sheet))
}

var Animation = donburi.NewComponentType[AnimationData]()
