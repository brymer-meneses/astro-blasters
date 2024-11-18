package component

import (
	"github.com/yohamta/donburi"
)

type ExplosionData struct {
	Count int
}

var Explosion = donburi.NewComponentType[ExplosionData]()
