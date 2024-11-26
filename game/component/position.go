package component

import (
	"github.com/yohamta/donburi"
	"math"
)

type PositionData struct {
	X     float64
	Y     float64
	Angle float64
}

var Position = donburi.NewComponentType[PositionData]()

func (self *PositionData) IntersectsWith(other *PositionData, radius float64) bool {
	return math.Pow(other.X-self.X, 2)+math.Pow(other.Y-self.Y, 2) <= math.Pow(radius, 2)
}

func (self *PositionData) Forward(magnitude float64) {
	self.Y -= magnitude * math.Cos(self.Angle)
	self.X += magnitude * math.Sin(self.Angle)
}

func (self *PositionData) Rotate(magnitude float64) {
	self.Angle += magnitude * math.Pi / 180
}
