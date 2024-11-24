package component

import (
	"space-shooter/game/types"

	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Name   string
	Health float64
	Score  int
	Id     types.PlayerId

	IsRotatingClockwise        bool
	IsRotatingCounterClockwise bool
	IsMovingForward            bool
	IsFiringBullet             bool
}

var Player = donburi.NewComponentType[PlayerData]()
