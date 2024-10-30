package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log"
	"net/http"
	"space-shooter/game/component"
	"space-shooter/rpc"
	"space-shooter/server/messages"

	"github.com/coder/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type Server struct {
	serveMux         http.ServeMux
	connectedPlayers int
	ecs              *ecs.ECS
	channel          chan Event
}

type Event struct {
	PlayerId messages.PlayerId
	Message  rpc.BaseMessage
}

func NewServer() *Server {
	cs := &Server{}
	cs.serveMux.HandleFunc("/events/ws", cs.ws)
	cs.serveMux.HandleFunc("/", cs.root)
	cs.ecs = ecs.NewECS(donburi.NewWorld())
	cs.channel = make(chan Event)
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

func (self *Server) handleConnection(connection *websocket.Conn) error {
	defer connection.CloseNow()
	ctx := context.Background()

	// Register the connected player.
	playerId, err := self.establishConnection(ctx, connection)
	if err != nil {
		return rpc.WriteMessage(ctx, connection, rpc.NewBaseMessage(messages.ErrorRoomFull{}))
	}

	// Handle server -> client updates
	go func() {
		for {
			event := <-self.channel
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if event.PlayerId == playerId {
				continue
			}

			if err := rpc.WriteMessage(ctx, connection, event.Message); err != nil {
				log.Println(err)
			}
		}
	}()

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
		switch message.MessageType {
		case "UpdatePosition":
			var updatePosition messages.UpdatePosition
			if err := msgpack.Unmarshal(message.Payload, &updatePosition); err != nil {
				continue
			}
			self.handleUpdatePosition(&updatePosition)
		}
		self.channel <- Event{
			Message:  message,
			PlayerId: playerId,
		}
	}

	log.Println("Connection ended")
	return nil
}

func (self *Server) getAvailablePlayerId() (messages.PlayerId, error) {
	if self.connectedPlayers == 5 {
		return 0, errors.New("Cannot have more than 5 players")
	}

	self.connectedPlayers += 1
	return messages.PlayerId(self.connectedPlayers - 1), nil
}

func (self *Server) handleUpdatePosition(updatePosition *messages.UpdatePosition) {
	player := self.findCorrespondPlayer(updatePosition.PlayerId)
	if player == nil {
		log.Printf("Cannot find player %d", updatePosition.PlayerId)
		return
	}
	donburi.SetValue(player, component.Position, updatePosition.Position)
}

func (self *Server) establishConnection(ctx context.Context, connection *websocket.Conn) (messages.PlayerId, error) {

	// Find a valid playerId for the current connection
	if self.connectedPlayers == 5 {
		return -1, nil
	}
	self.connectedPlayers += 1
	playerId := messages.PlayerId(self.connectedPlayers - 1)

	// Set the position of the current connection
	world := self.ecs.World
	entity := world.Create(component.Player, component.Position)
	player := world.Entry(entity)
	{
		donburi.SetValue(
			player,
			component.Player,
			component.PlayerData{
				Name: "Player One",
				Id:   int(playerId),
			},
		)
		donburi.SetValue(
			player,
			component.Position,
			component.PositionData{
				X:     500,
				Y:     10,
				Angle: 0,
			},
		)
	}

	position := component.Position.Get(player)
	enemyData := self.getEnemyData(playerId)

	establishConnection := rpc.NewBaseMessage(messages.EstablishConnection{
		PlayerId:  messages.PlayerId(playerId),
		Position:  *position,
		EnemyData: enemyData,
	})

	err := rpc.WriteMessage(
		ctx,
		connection,
		establishConnection,
	)

	if err != nil {
		log.Fatal(err)
	}
	if self.connectedPlayers > 1 {
		// Tell the other players that a new player is connected
		self.channel <- Event{
			PlayerId: playerId,
			Message: rpc.NewBaseMessage(messages.PlayerConnected{
				Position: *position,
				PlayerId: playerId,
			}),
		}
	}

	return playerId, nil
}

func (self *Server) findCorrespondPlayer(playerId messages.PlayerId) *donburi.Entry {
	query := donburi.NewQuery(filter.Contains(component.Player))
	for player := range query.Iter(self.ecs.World) {
		if int(playerId) == component.Player.GetValue(player).Id {
			return player
		}
	}
	return nil
}

func (self *Server) getEnemyData(playerId messages.PlayerId) []messages.EnemyData {
	// Get the position data of each player
	enemyData := make([]messages.EnemyData, self.connectedPlayers-1)
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position))
	i := 0

	for player := range query.Iter(self.ecs.World) {
		if playerId == messages.PlayerId(component.Player.Get(player).Id) {
			continue
		}

		enemyData[i] = messages.EnemyData{
			PlayerId: messages.PlayerId(component.Player.Get(player).Id),
			Position: *component.Position.Get(player),
		}
		i++
	}

	return enemyData
}
