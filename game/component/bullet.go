package component

import (
	"space-shooter/game/types"

	"github.com/yohamta/donburi"
)

type BulletData struct {
	FiredBy types.PlayerId
}

var Bullet = donburi.NewComponentType[BulletData]()
