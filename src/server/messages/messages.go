package messages

import (
	"space-shooter/game/component"

	"github.com/vmihailenco/msgpack/v5"
)

type BaseMessage struct {
	MessageType string
	Payload     msgpack.RawMessage
}

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
