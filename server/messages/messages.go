package messages

import (
	"space-shooter/game/component"
	"space-shooter/game/types"
)

type EnemyData struct {
	PlayerId types.PlayerId
	Position component.PositionData
}

type EstablishConnection struct {
	IsRoomFull bool
	PlayerId   types.PlayerId
	Position   component.PositionData

	EnemyData []EnemyData
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
