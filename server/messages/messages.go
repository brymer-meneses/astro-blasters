package messages

import (
	"space-shooter/scenes/game/component"
)

type PlayerId int

type EnemyData struct {
	PlayerId PlayerId
	Position component.PositionData
}

type EstablishConnection struct {
	IsRoomFull bool
	PlayerId   PlayerId
	Position   component.PositionData

	EnemyData []EnemyData
}

type PlayerConnected struct {
	PlayerId PlayerId
	Position component.PositionData
}

type PlayerDisconnected struct {
	PlayerId PlayerId
}

type UpdatePosition struct {
	PlayerId PlayerId
	Position component.PositionData
}
