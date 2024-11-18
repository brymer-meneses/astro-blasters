package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type ExplosionData struct {
	Count       int
	ExpiresWhen time.Time
}

var Explosion = donburi.NewComponentType[ExplosionData]()
