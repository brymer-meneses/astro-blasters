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
	"time"

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

	camera *Camera
}

func NewGameScene(config *config.AppConfig, assetManager *assets.AssetManager) *GameScene {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
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
		camera:       NewCamera(0, 0, config),
	}

	scene.ecs =
		ecs.NewECS(donburi.NewWorld()).
			AddRenderer(ecs.LayerDefault, scene.drawEnvironment).
			AddSystem(scene.movePlayer)

	scene.spawnPlayer(message.PlayerId, &message.Position)

	// Follow the player.
	scene.camera.FocusTarget(message.Position)

	for _, enemyData := range message.EnemyData {
		scene.spawnPlayer(messages.PlayerId(enemyData.PlayerId), &enemyData.Position)
	}

	go scene.receiveServerUpdates()
	return scene
}

func (self *GameScene) Draw(screen *ebiten.Image) {
	screen.Clear()

	self.ecs.DrawLayer(ecs.LayerDefault, screen)
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
		self.assetManager.Ships[playerId],
	)
}

// Draws the game environment.
func (self *GameScene) drawEnvironment(ecs *ecs.ECS, screen *ebiten.Image) {

	// Draw the background.
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-assets.MapWidth/2, -assets.MapHeight/2)
	opts.GeoM.Translate(self.camera.X, self.camera.Y)
	self.assetManager.Background.RenderWithOptions(screen, opts)

	// Loop through each player and draw each of them.
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))
	for player := range query.Iter(self.ecs.World) {
		sprite := component.Sprite.Get(player)
		position := component.Position.Get(player)

		// Center the texture
		x_0 := float64(sprite.Image.Bounds().Dx()) / 2
		y_0 := float64(sprite.Image.Bounds().Dy()) / 2

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-x_0, -y_0)

		opts.GeoM.Rotate(position.Angle)
		opts.GeoM.Scale(4, 4)
		opts.GeoM.Translate(position.X, position.Y)
		opts.GeoM.Translate(self.camera.X+x_0, self.camera.Y+y_0)

		// Render at this position
		screen.DrawImage(sprite.Image, opts)
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
		self.camera.FocusTarget(*positionData)
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

				player := findCorrespondingPlayer(self.ecs, updatePosition.PlayerId)
				if player != nil {
					component.Position.SetValue(player, updatePosition.Position)
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
		default:
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
