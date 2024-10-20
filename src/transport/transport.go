package transport

import (
	"context"
	"errors"
	"reflect"
	"space-shooter/server/messages"
	"time"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type Transport struct {
	connection *websocket.Conn
}

func Connect(socketURL string) (*Transport, error) {
	connection, _, err := websocket.Dial(context.Background(), socketURL, nil)
	if err != nil {
		return nil, err
	}

	return &Transport{connection}, nil
}

func FromConnection(connection *websocket.Conn) Transport {
	return Transport{connection}
}

func (self *Transport) SendMessage(message any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	payload, err := msgpack.Marshal(message)
	if err != nil {
		return err
	}

	encodedMessage := messages.BaseMessage{
		MessageType: reflect.TypeOf(message).Name(),
		Payload:     msgpack.RawMessage(payload),
	}

	bytes, err := msgpack.Marshal(encodedMessage)
	if err != nil {
		return err
	}

	if err := self.connection.Write(ctx, websocket.MessageBinary, bytes); err != nil {
		return err
	}

	return nil
}

func (self *Transport) ReceiveMessage(message any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, bytes, err := self.connection.Read(ctx)
	if err != nil {
		return err
	}

	var baseMessage messages.BaseMessage
	if err := msgpack.Unmarshal(bytes, &baseMessage); err != nil {
		return errors.New("Invalid message")
	}

	if err := msgpack.Unmarshal(baseMessage.Payload, message); err != nil {
		return errors.New("Unmarshalled a type that does not correspond with the type of message")
	}

	return nil
}
