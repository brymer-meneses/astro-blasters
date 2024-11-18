package arena

import (
	"space-shooter/client/config"
	"space-shooter/game/component"
)

type Camera struct {
	X      float64
	Y      float64
	config *config.ClientConfig
}

func NewCamera(x, y float64, config *config.ClientConfig) *Camera {
	return &Camera{
		X:      x,
		Y:      y,
		config: config,
	}
}

func (self *Camera) FocusTarget(target component.PositionData) {
	self.X = -target.X + float64(self.config.ScreenWidth)/2.0
	self.Y = -target.Y + float64(self.config.ScreenHeight)/2.0
}
