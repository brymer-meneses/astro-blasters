package component

import "github.com/yohamta/donburi"
import "math"

type PositionData struct {
	X     float64
	Y     float64
	Angle float64
}

var Position = donburi.NewComponentType[PositionData]()

func (self *PositionData) Forward(magnitude float64) {
	self.Y -= magnitude * math.Cos(self.Angle)
	self.X += magnitude * math.Sin(self.Angle)
}

func (self *PositionData) RotateCounterClockwise(magnitude float64) {
	self.Angle += magnitude * math.Pi / 180
}

func (self *PositionData) RotateClockwise(magnitude float64) {
	self.Angle -= magnitude * math.Pi / 180
}
