package component

import "github.com/yohamta/donburi"
import "math"

type PositionData struct {
	X     float64
	Y     float64
	Angle float64
}

var Position = donburi.NewComponentType[PositionData]()

func (self *PositionData) Forward() {
	self.Y -= 5 * math.Cos(self.Angle)
	self.X += 5 * math.Sin(self.Angle)
}

func (self *PositionData) RotateCounterClockwise() {
	self.Angle += 5 * math.Pi / 180
}

func (self *PositionData) RotateClockwise() {
	self.Angle -= 5 * math.Pi / 180
}
