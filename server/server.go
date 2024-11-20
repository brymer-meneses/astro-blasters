package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"log"
	"net/http"
	"space-shooter/game"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"space-shooter/rpc"
	"space-shooter/server/messages"

	"github.com/coder/websocket"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Server struct {
	serveMux   http.ServeMux
	simulation *game.GameSimulation

	players map[types.PlayerId]*playerConnection
}

type playerConnection struct {
	mutex       sync.Mutex
	conn        *websocket.Conn
	isConnected bool
}

func NewServer() *Server {
	s := &Server{}
	s.players = make(map[types.PlayerId]*playerConnection)

	s.serveMux.HandleFunc("/play/ws", s.ws)
	s.serveMux.HandleFunc("/", http.FileServer(http.Dir("server/static/")).ServeHTTP)

	s.simulation = game.NewGameSimulation()
	return s
}

func (self *Server) Start(port int) error {
	fmt.Printf("Server started at %s:%d\n", getLocalIP(), port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), &self.serveMux)
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
	ctx := context.Background()
	// Register the connected player.
	playerId, err := self.establishConnection(ctx, connection)

	defer func() {
		connection.CloseNow()
		self.players[playerId].isConnected = false
		log.Println("Connection ended")
	}()

	if err != nil {
		return err
	}

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
			if err := rpc.DecodeExpectedMessage(message, &updatePosition); err != nil {
				continue
			}
			self.handleUpdatePosition(&updatePosition)
			self.broadcastMessage(playerId, message)
		case "FireBullet":
			var fireBullet messages.FireBullet
			if err := rpc.DecodeExpectedMessage(message, &fireBullet); err != nil {
				continue
			}
			self.broadcastMessage(playerId, message)
		}
	}

	return nil
}

func (self *Server) broadcastMessage(from types.PlayerId, message rpc.BaseMessage) {
	sendMessage := func(playerId types.PlayerId, playerConn *playerConnection) {
		// Lock the mutex to ensure only one goroutine writes at a time
		playerConn.mutex.Lock()
		defer playerConn.mutex.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Write the message to the connection
		err := rpc.WriteMessage(ctx, playerConn.conn, message)
		if err != nil {
			log.Printf("Failed to send message to player %d: %v", playerId, err)
		}
	}

	// For each playerid that does not match the sender, send the message.
	for playerId, playerConn := range self.players {
		if !playerConn.isConnected {
			continue
		}

		if from == playerId {
			continue
		}
		go sendMessage(playerId, playerConn)
	}
}

func (self *Server) getAvailablePlayerId() (types.PlayerId, error) {
	connectedPlayers := len(self.players)
	if connectedPlayers == 5 {
		return 0, errors.New("Cannot have more than 5 players")
	}
	return types.PlayerId(connectedPlayers), nil
}

func (self *Server) handleUpdatePosition(updatePosition *messages.UpdatePosition) {
	player := self.simulation.FindCorrespondingPlayer(updatePosition.PlayerId)
	if player == nil {
		log.Printf("Cannot find player %d", updatePosition.PlayerId)
		return
	}
	donburi.SetValue(player, component.Position, updatePosition.Position)
}

func (self *Server) establishConnection(ctx context.Context, connection *websocket.Conn) (types.PlayerId, error) {
	playerId, err := self.getAvailablePlayerId()
	if err != nil {
		return types.InvalidPlayerId, rpc.WriteMessage(ctx, connection, rpc.NewBaseMessage(messages.EstablishConnection{IsRoomFull: true}))
	}

	position := component.PositionData{
		X:     0,
		Y:     0,
		Angle: 0,
	}

	self.players[playerId] = &playerConnection{
		conn:        connection,
		isConnected: true,
	}

	self.simulation.SpawnPlayer(playerId, &position)

	playerData := self.getPlayerData()
	err = rpc.WriteMessage(
		ctx,
		connection,
		rpc.NewBaseMessage(messages.EstablishConnection{
			PlayerId:   playerId,
			PlayerData: playerData,
			IsRoomFull: false,
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	self.broadcastMessage(playerId, rpc.NewBaseMessage(messages.PlayerConnected{
		PlayerId: playerId,
		Position: position,
	}))

	return playerId, nil
}

func (self *Server) getPlayerData() []messages.PlayerData {
	// Get the position data of each player
	enemyData := make([]messages.PlayerData, len(self.players))
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position))
	i := 0

	for player := range query.Iter(self.simulation.ECS.World) {
		enemyData[i] = messages.PlayerData{
			PlayerId: component.Player.Get(player).Id,
			Position: *component.Position.Get(player),
		}
		i++
	}
	return enemyData
}

// From: https://stackoverflow.com/a/31551220
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
