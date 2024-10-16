package scenes

import (
	"context"
	"log"
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/game/component"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameScene struct {
	assetManager *assets.AssetManager
	ecs          *ecs.ECS
	connection   *websocket.Conn
	playerId     component.PlayerId
}

func NewGameScene(config *config.AppConfig, assetManager *assets.AssetManager, playerId component.PlayerId) *GameScene {
	connection, _, err := websocket.Dial(context.Background(), config.ServerWebsocketURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	scene := &GameScene{assetManager: assetManager, playerId: 0, connection: connection}
	scene.ecs =
		ecs.NewECS(donburi.NewWorld()).
			AddRenderer(0, scene.drawEnvironment).
			AddSystem(scene.movePlayer).
			AddSystem(scene.receiveServerUpdates)

	scene.createPlayer(component.PlayerId(0))
	return scene
}

func (self *GameScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	self.assetManager.Background.Render(screen)

	self.ecs.DrawLayer(0, screen)
	self.ecs.Draw(screen)
}

func (self *GameScene) Update() {
	self.ecs.Update()
}

func (self *GameScene) createPlayer(playerId component.PlayerId) {

	world := self.ecs.World
	entity := world.Create(component.Player, component.Position, component.Sprite)
	player := world.Entry(entity)

	donburi.SetValue(
		player,
		component.Player,
		component.PlayerData{
			Name: "Player One",
			Id:   playerId,
		},
	)

	donburi.SetValue(
		player,
		component.Position,
		component.PositionData{
			X:     0,
			Y:     0,
			Angle: 0,
		},
	)

	donburi.SetValue(
		player,
		component.Sprite,
		component.SpriteData{Image: self.assetManager.Ships[playerId].Image},
	)

}

// Draws the game environment.
func (self *GameScene) drawEnvironment(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))

	for player := range query.Iter(self.ecs.World) {
		sprite := component.Sprite.Get(player)
		position := component.Position.Get(player)

		op := &ebiten.DrawImageOptions{}

		x_0 := float64(sprite.Image.Bounds().Dx()) / 2
		y_0 := float64(sprite.Image.Bounds().Dy()) / 2

		op.GeoM.Translate(-x_0, -y_0)

		op.GeoM.Rotate(position.Angle)
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(position.X, position.Y)

		screen.DrawImage(sprite.Image, op)
	}
}

// Handles the movement of the player and sends it to the server.
func (self *GameScene) movePlayer(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))

	for player := range query.Iter(ecs.World) {
		if component.PlayerId(self.playerId) != component.Player.GetValue(player).Id {
			continue
		}

		positionData := component.Position.Get(player)
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			self.connection.Write(context.Background(), 1, []byte("Forward!"))
			positionData.Forward()
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			self.connection.Write(context.Background(), 1, []byte("Clockwise!"))
			positionData.RotateClockwise()
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			self.connection.Write(context.Background(), 1, []byte("Counterclockwise!"))
			positionData.RotateCounterClockwise()
		}
	}
}

// Receives information from the server and updates the game state accordingly.
func (self *GameScene) receiveServerUpdates(_ *ecs.ECS) {
	for {
		_, bytes, err := self.connection.Read(context.Background())
		if err != nil {
			break
		}

		log.Printf("%s\n", string(bytes))
	}
}
