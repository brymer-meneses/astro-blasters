package messages

import (
	"space-shooter/game/component"
	"space-shooter/game/types"
)

type PlayerData struct {
	PlayerId types.PlayerId
	Position component.PositionData
}

type EstablishConnection struct {
	IsRoomFull bool
	PlayerId   types.PlayerId
	PlayerData []PlayerData
}

type PlayerConnected struct {
	PlayerId types.PlayerId
	Position component.PositionData
}

type PlayerDisconnected struct {
	PlayerId types.PlayerId
}

type UpdatePosition struct {
	PlayerId types.PlayerId
	Position component.PositionData
}
