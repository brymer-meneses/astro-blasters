package server

import (
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"space-shooter/game/component"
	"space-shooter/rpc"
	"space-shooter/server/messages"

	"github.com/coder/websocket"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type PlayerId int

type Server struct {
	serveMux         http.ServeMux
	connectedPlayers int
	ecs              *ecs.ECS
	channel          chan rpc.BaseMessage
}

func NewServer() *Server {
	cs := &Server{}
	cs.serveMux.HandleFunc("/events/ws", cs.ws)
	cs.serveMux.HandleFunc("/", cs.root)
	cs.ecs = ecs.NewECS(donburi.NewWorld())
	cs.channel = make(chan rpc.BaseMessage)
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
	defer connection.CloseNow()

	ctx := context.Background()

	playerId, err := self.getAvailablePlayerId()
	if err != nil {
		if err := rpc.SendMessage(ctx, connection, messages.ErrorRoomFull{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	position := self.registerPlayer(playerId)
	err = rpc.SendMessage(
		ctx,
		connection,
		messages.EstablishConnection{
			PlayerId: playerId,
			Position: *position,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	// Handle client -> server updates
	for {
		var message rpc.BaseMessage
		err := rpc.ReceiveMessage(ctx, connection, &message)
		status := websocket.CloseStatus(err)

		if status == websocket.StatusGoingAway || status == websocket.StatusAbnormalClosure {
			break
		}

		if err != nil {
			break
		}

		log.Printf("%+v", message.MessageType)

		switch message.MessageType {
		case "UpdatePosition":
			var updatePosition messages.UpdatePosition
			rpc.Cast(&message, &updatePosition)

			self.handleUpdatePosition(&updatePosition)
		}
	}

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
