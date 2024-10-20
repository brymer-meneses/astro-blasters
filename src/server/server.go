package server

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"space-shooter/game/component"
	"space-shooter/server/messages"
	"space-shooter/transport"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerId int

type Server struct {
	serveMux         http.ServeMux
	connectedPlayers int
	ecs              *ecs.ECS
	channel          chan messages.BaseMessage
}

func NewServer() *Server {
	cs := &Server{}
	cs.serveMux.HandleFunc("/events/ws", cs.ws)
	cs.serveMux.HandleFunc("/", cs.root)
	cs.ecs = ecs.NewECS(donburi.NewWorld())
	cs.channel = make(chan messages.BaseMessage)
	return cs
}

func (self *Server) Start(port string) error {
	log.Printf("Listening at %s", port)
	return http.ListenAndServe(":"+port, &self.serveMux)
}

func (self *Server) root(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func (self *Server) ws(w http.ResponseWriter, r *http.Request) {
	connection, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "Connection Failed")
		return
	}

	self.handleConnection(connection)
}

func (self *Server) handleConnection(connection *websocket.Conn) {
	transport := transport.FromConnection(connection)
	playerId, err := self.getAvailablePlayerId()
	if err != nil {
		transport.SendMessage(messages.ErrorRoomFull{})
		return
	}

	position := self.registerPlayer(playerId)
	err = transport.SendMessage(
		messages.EstablishConnection{
			PlayerId: playerId,
			Position: *position,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Handle client -> server updates
	go func() {
		defer connection.CloseNow()
		for {
			_, bytes, err := connection.Read(ctx)
			status := websocket.CloseStatus(err)
			if status == websocket.StatusGoingAway || status == websocket.StatusAbnormalClosure {
				break
			}

			var message messages.BaseMessage
			if err := msgpack.Unmarshal(bytes, &message); err != nil {
				log.Printf("Invalid Message Format")
				break
			}

			switch message.MessageType {
			case "UpdatePosition":
				var payload messages.UpdatePosition
				if err := msgpack.Unmarshal(message.Payload, &payload); err != nil {
					log.Fatal(err)
				}
				self.handleUpdatePosition(&payload)
			}

			// Publish this message to others
			self.channel <- message
		}
	}()

	// Handle server -> client updates
	go func() {
		for {
			message := <-self.channel
			bytes, err := msgpack.Marshal(message)

			if err != nil {
				log.Printf("Invalid Message Format")
				break
			}

			connection.Write(context.Background(), websocket.MessageBinary, bytes)
		}
	}()

	log.Println("Connection ended")
}

func (self *Server) getAvailablePlayerId() (int, error) {
	if self.connectedPlayers == 5 {
		return 0, errors.New("Cannot have more than 5 players")
	}

	self.connectedPlayers += 1
	return self.connectedPlayers - 1, nil
}

func (self *Server) handleUpdatePosition(updatePosition *messages.UpdatePosition) {
	log.Printf("%+v", updatePosition)
}

func (self *Server) registerPlayer(playerId int) *component.PositionData {
	world := self.ecs.World
	entity := world.Create(component.Player, component.Position)
	player := world.Entry(entity)

	donburi.SetValue(
		player,
		component.Player,
		component.PlayerData{
			Name: "Player One",
			Id:   playerId,
		},
	)

	// TODO: Make this random for each player, like make them sparse
	donburi.SetValue(
		player,
		component.Position,
		component.PositionData{
			X:     500,
			Y:     10,
			Angle: 0,
		},
	)

	return component.Position.Get(player)
}
