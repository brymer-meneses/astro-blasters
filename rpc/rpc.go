package rpc

import (
	"context"
	"log"
	"reflect"
	"sync"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type BaseMessage struct {
	MessageType string
	Payload     msgpack.RawMessage
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		// Create a new buffer if one isn't available in the pool.
		// This example uses a buffer with a size of 4096 bytes (4KB).
		return make([]byte, 4096)
	},
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
	buffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(buffer)

	_, reader, err := conn.Reader(ctx)
	if err != nil {
		return err
	}

	n, err := reader.Read(buffer)
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(buffer[:n], message)
}

func ReceiveExpectedMessage[ExpectedMessage any](ctx context.Context, conn *websocket.Conn, out *ExpectedMessage) error {
	var baseMessage BaseMessage
	if err := ReceiveMessage(ctx, conn, &baseMessage); err != nil {
		return err
	}
	return msgpack.Unmarshal(baseMessage.Payload, out)
}

func DecodeExpectedMessage[ExpectedMessage any](message BaseMessage, out *ExpectedMessage) error {
	if err := msgpack.Unmarshal(message.Payload, out); err != nil {
		return err
	}
	return nil
}
