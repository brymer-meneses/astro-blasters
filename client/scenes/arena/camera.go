package arena

import (
	"math"
	"space-shooter/client/config"
	"space-shooter/game/component"
)

type Camera struct {
	X           float64
	Y           float64
	SceneWidth  float64
	SceneHeight float64
	config      *config.ClientConfig
}

func NewCamera(x, y, sceneWidth, sceneHeight float64, config *config.ClientConfig) *Camera {
	return &Camera{
		X:           x,
		Y:           y,
		SceneWidth:  sceneWidth,
		SceneHeight: sceneHeight,
		config:      config,
	}
}

func (self *Camera) FocusTarget(target component.PositionData) {
	self.X = -target.X + float64(self.config.ScreenWidth)/2.0
	self.Y = -target.Y + float64(self.config.ScreenHeight)/2.0
}

func (self *Camera) Constrain() {
	self.X = math.Min(self.X, 0)
	self.Y = math.Min(self.Y, 0)

	self.X = math.Max(self.X, -float64(self.SceneWidth)+float64(self.config.ScreenWidth))
	self.Y = math.Max(self.Y, -float64(self.SceneHeight)+float64(self.config.ScreenHeight))
}
