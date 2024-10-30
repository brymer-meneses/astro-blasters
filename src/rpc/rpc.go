package rpc

import (
	"context"
	"log"
	"reflect"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type BaseMessage struct {
	MessageType string
	Payload     msgpack.RawMessage
}

func NewBaseMessage(message any) BaseMessage {
	payload, err := msgpack.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	encoded := BaseMessage{
		MessageType: reflect.TypeOf(message).Name(),
		Payload:     payload,
	}
	return encoded
}

func WriteMessage(ctx context.Context, conn *websocket.Conn, message BaseMessage) error {
	marshaled, err := msgpack.Marshal(message)
	if err != nil {
		return err
	}
	return conn.Write(ctx, websocket.MessageBinary, marshaled)
}

func ReceiveMessage(ctx context.Context, conn *websocket.Conn, message *BaseMessage) error {
	_, bytes, err := conn.Read(ctx)
	if err != nil {
		return nil
	}
	return msgpack.Unmarshal(bytes, message)
}

func ReceiveExpectedMessage[ExpectedMessage any](ctx context.Context, conn *websocket.Conn, out *ExpectedMessage) error {
	var baseMessage BaseMessage
	if err := ReceiveMessage(ctx, conn, &baseMessage); err != nil {
		return err
	}
	return msgpack.Unmarshal(baseMessage.Payload, out)
}
