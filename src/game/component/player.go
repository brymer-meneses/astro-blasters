package component

import (
	"github.com/yohamta/donburi"
)

type PlayerId int

type PlayerData struct {
	Name   string
	Id     PlayerId
	Health float64
	// Profile *ebiten.Image
}

var Player = donburi.NewComponentType[PlayerData]()
