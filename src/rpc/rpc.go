package rpc

import (
	"context"
	"errors"
	"log"
	"reflect"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type BaseMessage struct {
	MessageType string
	Payload     msgpack.RawMessage
}

func Cast[MessageType any](message *BaseMessage, out *MessageType) error {
	to := reflect.TypeFor[MessageType]().Name()
	expected := message.MessageType

	if to != expected {
		log.Fatalf("Tried to cast message of type %s to %s", expected, to)
	}

	if err := msgpack.Unmarshal(message.Payload, out); err != nil {
		log.Fatalf("Failed to unmarshal: %s", err)
	}

	return nil
}

func ReceiveMessage(ctx context.Context, conn *websocket.Conn, out *BaseMessage) error {
	_, bytes, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	if err := msgpack.Unmarshal(bytes, out); err != nil {
		return errors.New("Invalid message")
	}

	return nil
}

func SendMessage[MessageType any](ctx context.Context, conn *websocket.Conn, message MessageType) error {
	payload, err := msgpack.Marshal(message)
	if err != nil {
		return err
	}

	encodedMessage := BaseMessage{
		MessageType: reflect.TypeFor[MessageType]().Name(),
		Payload:     msgpack.RawMessage(payload),
	}

	bytes, err := msgpack.Marshal(encodedMessage)
	if err != nil {
		return err
	}

	if err := conn.Write(ctx, websocket.MessageBinary, bytes); err != nil {
		return err
	}

	return nil
}
