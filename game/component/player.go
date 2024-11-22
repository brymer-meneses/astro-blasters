package component

import (
	"space-shooter/game/types"

	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Name   string
	Health float64
	Id     types.PlayerId
}

var Player = donburi.NewComponentType[PlayerData]()
