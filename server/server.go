package server

import (
	"context"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"astro-blasters/game"
	"astro-blasters/game/component"
	"astro-blasters/game/types"
	"astro-blasters/rpc"
	"astro-blasters/server/messages"
	"log"
	"net/http"

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
	mutex          sync.Mutex
	conn           *websocket.Conn
	isConnected    bool
	lastBulletFire time.Time
}

func NewServer() *Server {
	s := &Server{}
	s.players = make(map[types.PlayerId]*playerConnection)

	s.serveMux.HandleFunc("/play/ws", s.ws)
	s.serveMux.Handle("/", http.FileServer(http.Dir("server/static/")))

	s.simulation = game.NewGameSimulation()

	s.simulation.OnBulletCollide = s.onBulletCollide
	s.simulation.OnBulletFire = s.onBulletFire
	return s
}

func (self *Server) onBulletFire(player *donburi.Entry) {
	playerId := component.Player.Get(player).Id
	connection := self.players[playerId]
	now := time.Now()

	if connection.lastBulletFire.IsZero() || now.Sub(connection.lastBulletFire) >= 300*time.Millisecond {
		connection.lastBulletFire = now
		self.broadcastMessage(rpc.NewBaseMessage(messages.EventPlayerFireBullet{
			PlayerId: playerId,
		}))
		self.simulation.RegisterPlayerFire(player)
	}
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
	if err != nil {
		return err
	}

	defer func() {
		connection.CloseNow()
		player := self.simulation.FindCorrespondingPlayer(playerId)
		self.simulation.RegisterPlayerDisconnection(player)
		self.players[playerId].isConnected = false
		self.broadcastMessageExcept(playerId, rpc.NewBaseMessage(messages.EventPlayerDisconnected{
			PlayerId: playerId,
		}))
	}()

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
	for playerId, playerConn := range self.players {
		go self.sendMessage(playerId, playerConn, message)
	}
}

// For each playerid that does not match the sender, send the message.
func (self *Server) broadcastMessageExcept(except types.PlayerId, message rpc.BaseMessage) {
	for playerId, playerConn := range self.players {
		if except == playerId {
			continue
		}
		go self.sendMessage(playerId, playerConn, message)
	}
}

func (self *Server) getAvailablePlayerId() types.PlayerId {
	return types.PlayerId(len(self.players))
}

func (self *Server) establishConnection(ctx context.Context, connection *websocket.Conn) (types.PlayerId, error) {
	var connectionHandshake messages.ConnectionHandshake
	if err := rpc.ReceiveExpectedMessage(ctx, connection, &connectionHandshake); err != nil {
		return types.InvalidPlayerId, nil
	}

	playerId := self.getAvailablePlayerId()
	position := game.GenerateRandomPlayerPosition()

	self.players[playerId] = &playerConnection{
		conn:        connection,
		isConnected: true,
	}

	self.simulation.CreatePlayer(playerId, &position, connectionHandshake.PlayerName, true)

	playerData := self.getPlayerData()
	err := rpc.WriteMessage(
		ctx,
		connection,
		rpc.NewBaseMessage(messages.ConnectionHandshakeResponse{
			PlayerId:   playerId,
			PlayerData: playerData,
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
	enemyData := []messages.PlayerData{}
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position))

	for player := range query.Iter(self.simulation.ECS.World) {
		data := component.Player.Get(player)

		enemyData = append(enemyData,
			messages.PlayerData{
				PlayerId:    data.Id,
				PlayerName:  data.Name,
				IsConnected: data.IsConnected,
				Position:    *component.Position.Get(player),
			},
		)
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
