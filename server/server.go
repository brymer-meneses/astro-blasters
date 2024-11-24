package server

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	s.serveMux.Handle("/", http.FileServer(http.Dir("server/static/")))

	s.simulation = game.NewGameSimulation()
	s.simulation.OnBulletCollide = s.onBulletCollide
	return s
}

func (self *Server) onBulletCollide(player *donburi.Entry, bullet *donburi.Entry) {
	playerData := component.Player.Get(player)
	playerData.Health -= game.PlayerDamagePerHit

	if playerData.Health > 0 {
		self.broadcastMessage(rpc.NewBaseMessage(messages.EventUpdateHealth{
			PlayerId: playerData.Id,
			Health:   playerData.Health,
		}))
	} else if playerData.Health == 0 {
		bulletData := component.Bullet.Get(bullet)
		scorer := self.simulation.FindCorrespondingPlayer(bulletData.FiredBy)

		scorerData := component.Player.Get(scorer)

		self.broadcastMessage(rpc.NewBaseMessage(messages.EventPlayerDied{
			PlayerId: playerData.Id,
			KilledBy: scorerData.Id,
		}))

		self.simulation.RegisterPlayerDeath(player, scorer)

		go func() {
			time.Sleep(5 * time.Second)
			position := game.GenerateRandomPlayerPosition()
			self.simulation.RespawnPlayer(player, position)

			self.broadcastMessage(rpc.NewBaseMessage(messages.EventPlayerRespawned{
				PlayerId: playerData.Id,
				Position: position,
			}))
		}()
	}
}

func (self *Server) Start(port int) error {
	fmt.Printf("Server started at %s:%d\n", getLocalIP(), port)

	go self.updateState()

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
		case "RegisterPlayerMove":
			var registerPlayerMove messages.RegisterPlayerMove
			if err := rpc.DecodeExpectedMessage(message, &registerPlayerMove); err != nil {
				continue
			}
			player := self.simulation.FindCorrespondingPlayer(playerId)
			expectedPosition := component.Position.Get(player)

			if !isPositionWithinTolerance(*expectedPosition, registerPlayerMove.Position, 3.0) {
				self.broadcastMessage(rpc.NewBaseMessage(messages.UpdatePosition{
					Position: *expectedPosition,
					PlayerId: playerId,
				}))
			}

			self.simulation.RegisterPlayerMove(playerId, registerPlayerMove.Move)
			self.broadcastMessage(rpc.NewBaseMessage(messages.EventPlayerMove{
				Move:     registerPlayerMove.Move,
				PlayerId: playerId,
			}))
		}
	}
	return nil
}

func isPositionWithinTolerance(expected component.PositionData, got component.PositionData, tolerance float64) bool {
	return math.Pow(expected.X-got.X, 2)+math.Pow(expected.Y-got.Y, 2)+math.Pow(expected.Angle-got.Angle, 2) <= math.Pow(tolerance, 2)
}

func (self *Server) updateState() {
	ticker := time.NewTicker(time.Millisecond * 16) // ~60 FPS
	defer ticker.Stop()

	for range ticker.C {
		self.simulation.Update()
	}
}

func (self *Server) sendMessage(playerId types.PlayerId, playerConn *playerConnection, message rpc.BaseMessage) {
	// Lock the mutex to ensure only one goroutine writes at a time
	playerConn.mutex.Lock()
	defer playerConn.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if !playerConn.isConnected {
		return
	}

	// Write the message to the connection
	err := rpc.WriteMessage(ctx, playerConn.conn, message)
	if err != nil {
		log.Printf("Failed to send message to player %d: %v", playerId, err)
	}
}

func (self *Server) broadcastMessage(message rpc.BaseMessage) {
	// For each playerid that does not match the sender, send the message.
	for playerId, playerConn := range self.players {
		go self.sendMessage(playerId, playerConn, message)
	}
}

func (self *Server) broadcastMessageExcept(except types.PlayerId, message rpc.BaseMessage) {
	// For each playerid that does not match the sender, send the message.
	for playerId, playerConn := range self.players {
		if except == playerId {
			continue
		}
		go self.sendMessage(playerId, playerConn, message)
	}
}

func (self *Server) getAvailablePlayerId() (types.PlayerId, error) {
	connectedPlayers := len(self.players)
	if connectedPlayers == 5 {
		return 0, errors.New("Cannot have more than 5 players")
	}
	return types.PlayerId(connectedPlayers), nil
}

func (self *Server) establishConnection(ctx context.Context, connection *websocket.Conn) (types.PlayerId, error) {
	var connectionHandshake messages.ConnectionHandshake
	if err := rpc.ReceiveExpectedMessage(ctx, connection, &connectionHandshake); err != nil {
		return types.InvalidPlayerId, nil
	}

	playerId, err := self.getAvailablePlayerId()
	if err != nil {
		response := rpc.NewBaseMessage(messages.ConnectionHandshakeResponse{IsRoomFull: true})
		return types.InvalidPlayerId, rpc.WriteMessage(ctx, connection, response)
	}

	position := game.GenerateRandomPlayerPosition()

	self.players[playerId] = &playerConnection{
		conn:        connection,
		isConnected: true,
	}

	self.simulation.SpawnPlayer(playerId, &position, connectionHandshake.PlayerName)

	playerData := self.getPlayerData()
	err = rpc.WriteMessage(
		ctx,
		connection,
		rpc.NewBaseMessage(messages.ConnectionHandshakeResponse{
			PlayerId:   playerId,
			PlayerData: playerData,
			IsRoomFull: false,
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	// Tell the other players that this player has joined.
	self.broadcastMessageExcept(playerId, rpc.NewBaseMessage(messages.EventPlayerConnected{
		PlayerId:   playerId,
		PlayerName: connectionHandshake.PlayerName,
		Position:   position,
	}))

	return playerId, nil
}

func (self *Server) getPlayerData() []messages.PlayerData {
	// Get the position data of each player
	enemyData := make([]messages.PlayerData, len(self.players))
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position))
	i := 0

	for player := range query.Iter(self.simulation.ECS.World) {
		data := component.Player.Get(player)
		enemyData[i] = messages.PlayerData{
			PlayerId:   data.Id,
			PlayerName: data.Name,
			Position:   *component.Position.Get(player),
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
