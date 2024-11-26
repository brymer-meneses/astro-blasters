package component

import (
	"astro-blasters/game/types"

	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Name   string
	Health float64
	Score  int
	Id     types.PlayerId

	IsAlive     bool
	IsConnected bool

	IsRotatingClockwise        bool
	IsRotatingCounterClockwise bool
	IsMovingForward            bool
	IsFiringBullet             bool
}

var Player = donburi.NewComponentType[PlayerData]()
