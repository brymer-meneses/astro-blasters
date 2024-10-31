package game

import (
	"context"
	"log"
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/rpc"
	"space-shooter/scenes"
	"space-shooter/scenes/game/component"
	"space-shooter/server/messages"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameScene struct {
	assetManager *assets.AssetManager
	ecs          *ecs.ECS
	connection   *websocket.Conn
	playerId     messages.PlayerId
}

func NewGameScene(config *config.AppConfig, assetManager *assets.AssetManager) *GameScene {
	ctx := context.Background()
	connection, _, err := websocket.Dial(ctx, config.ServerWebsocketURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to the game server at %s\n", config.ServerWebsocketURL)
	}

	var message messages.EstablishConnection
	if err := rpc.ReceiveExpectedMessage(ctx, connection, &message); err != nil {
		log.Fatal(err)
	}

	if message.IsRoomFull {
		log.Fatal("Room is full")
	}

	scene := &GameScene{
		assetManager: assetManager,
		playerId:     message.PlayerId,
		connection:   connection,
	}

	scene.ecs =
		ecs.NewECS(donburi.NewWorld()).
			AddRenderer(0, scene.drawEnvironment).
			AddSystem(scene.movePlayer)

	scene.spawnPlayer(message.PlayerId, &message.Position)

	for _, enemyData := range message.EnemyData {
		scene.spawnPlayer(messages.PlayerId(enemyData.PlayerId), &enemyData.Position)
	}

	go scene.receiveServerUpdates()
	return scene
}

func (self *GameScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	self.assetManager.Background.Render(screen)

	self.ecs.DrawLayer(0, screen)
	self.ecs.Draw(screen)
}

func (self *GameScene) Update(dispatcher *scenes.SceneDispatcher) {
	self.ecs.Update()
}

// Spawns the player in the game.
func (self *GameScene) spawnPlayer(playerId messages.PlayerId, position *component.PositionData) {
	world := self.ecs.World
	entity := world.Create(component.Player, component.Position, component.Sprite)
	player := world.Entry(entity)

	component.Player.SetValue(
		player,
		component.PlayerData{
			Name: "Player One",
			Id:   int(playerId),
		},
	)

	component.Position.SetValue(
		player,
		*position,
	)

	component.Sprite.SetValue(
		player,
		component.SpriteData{Image: self.assetManager.Ships[playerId].Image},
	)
}

// Draws the game environment.
func (self *GameScene) drawEnvironment(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))

	// Loop each player
	for player := range query.Iter(self.ecs.World) {
		sprite := component.Sprite.Get(player)
		position := component.Position.Get(player)

		op := &ebiten.DrawImageOptions{}

		// Center the texture
		x_0 := float64(sprite.Image.Bounds().Dx()) / 2
		y_0 := float64(sprite.Image.Bounds().Dy()) / 2
		op.GeoM.Translate(-x_0, -y_0)

		op.GeoM.Rotate(position.Angle)
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(position.X, position.Y)

		// Render at this position
		screen.DrawImage(sprite.Image, op)
	}
}

// Handles the movement of the player and sends it to the server.
func (self *GameScene) movePlayer(ecs *ecs.ECS) {

	updatePosition := func(positionData *component.PositionData) {
		message := rpc.NewBaseMessage(
			messages.UpdatePosition{
				PlayerId: self.playerId,
				Position: *positionData,
			})

		rpc.WriteMessage(context.Background(), self.connection, message)
	}

	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))

	for player := range query.Iter(ecs.World) {
		if int(self.playerId) != component.Player.GetValue(player).Id {
			continue
		}

		positionData := component.Position.Get(player)
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			positionData.Forward()
			updatePosition(positionData)
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			positionData.RotateClockwise()
			updatePosition(positionData)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			positionData.RotateCounterClockwise()
			updatePosition(positionData)
		}

	}
}

// Receives information from the server and updates the game state accordingly.
func (self *GameScene) receiveServerUpdates() {
	for {
		var message rpc.BaseMessage
		if err := rpc.ReceiveMessage(context.Background(), self.connection, &message); err != nil {
			continue
		}

		switch message.MessageType {
		case "UpdatePosition":
			{
				var updatePosition messages.UpdatePosition
				if err := msgpack.Unmarshal(message.Payload, &updatePosition); err != nil {
					continue
				}

				entry := findCorrespondingPlayer(self.ecs, updatePosition.PlayerId)
				if entry != nil {
					component.Position.SetValue(entry, updatePosition.Position)
				}
			}
		case "PlayerConnected":
			{
				var playerConnected messages.PlayerConnected
				if err := msgpack.Unmarshal(message.Payload, &playerConnected); err != nil {
					continue
				}

				self.spawnPlayer(playerConnected.PlayerId, &playerConnected.Position)
			}
		}
	}

}

// Returns the ecs entry given the playerId.
func findCorrespondingPlayer(ecs *ecs.ECS, playerId messages.PlayerId) *donburi.Entry {
	query := donburi.NewQuery(filter.Contains(component.Player))
	for player := range query.Iter(ecs.World) {
		if int(playerId) == component.Player.GetValue(player).Id {
			return player
		}
	}
	return nil
}
