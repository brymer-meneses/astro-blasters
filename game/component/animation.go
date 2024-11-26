package component

import (
	"astro-blasters/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type AnimationData struct {
	sheet            assets.SpriteSheet
	activeFrameIndex int
	frameTimeCounter int
	framesPerSecond  int
}

func NewAnimationData(sheet assets.SpriteSheet, framesPerSecond int) AnimationData {
	return AnimationData{
		sheet:            sheet,
		framesPerSecond:  framesPerSecond,
		frameTimeCounter: 0,
		activeFrameIndex: 0,
	}
}

func (self *AnimationData) Frame() *ebiten.Image {
	if self.frameTimeCounter == self.framesPerSecond {
		self.frameTimeCounter = 0
		self.activeFrameIndex = (self.activeFrameIndex + 1) % self.sheet.TotalFrames()
	} else {
		self.frameTimeCounter += 1
	}

	return self.sheet.GetFrame(self.activeFrameIndex)
}

var Animation = donburi.NewComponentType[AnimationData]()
