package messages

import "space-shooter/game/component"

type EstablishConnection struct {
	PlayerId int
	Position component.PositionData
}

type UpdatePosition struct {
	PlayerId int
	Position component.PositionData
}

type ErrorRoomFull struct {
	PlayerId int
}
