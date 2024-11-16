package component

import (
	"space-shooter/game/types"
	"time"

	"github.com/yohamta/donburi"
)

type BulletData struct {
	FiredBy  types.PlayerId
	ShotWhen time.Time
}

var Bullet = donburi.NewComponentType[BulletData]()
