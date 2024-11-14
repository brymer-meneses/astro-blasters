package arena

import (
	"context"
	"log"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/game"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"space-shooter/rpc"
	"space-shooter/server/messages"
	"time"
	"image/color"
	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	MapWidth  = 4000
	MapHeight = 4000
	MinimapWidth  = 150
	MinimapHeight = 150
)

type ArenaScene struct {
	connection *websocket.Conn
	simulation *game.GameSimulation
	background *common.Background
	playerId   types.PlayerId
	camera     *Camera
}

func NewArenaScene(config *config.ClientConfig) *ArenaScene {
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

	scene := &ArenaScene{
		background: common.NewBackground(MapWidth, MapHeight),
		playerId:   message.PlayerId,
		simulation: game.NewGameSimulation(),
		connection: connection,
		camera:     NewCamera(0, 0, config),
	}

	scene.simulation.SpawnPlayer(message.PlayerId, &message.Position)
	for _, enemy := range message.EnemyData {
		scene.simulation.SpawnPlayer(enemy.PlayerId, &enemy.Position)
	}

	// Focus the camera on the player.
	scene.camera.FocusTarget(message.Position)

	go scene.receiveServerUpdates()
	return scene
}

func (self *ArenaScene) Draw(screen *ebiten.Image) {
	screen.Clear()
	// Draw the background.
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-MapWidth/2, -MapHeight/2)
	opts.GeoM.Translate(self.camera.X, self.camera.Y)
	screen.DrawImage(self.background.Image, opts)

	// Loop through each player and draw each of them.
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))
	for player := range query.Iter(self.simulation.ECS.World) {
		sprite := component.Sprite.GetValue(player)
		position := component.Position.Get(player)

		// Center the texture
		x_0 := float64(sprite.Bounds().Dx()) / 2
		y_0 := float64(sprite.Bounds().Dy()) / 2

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-x_0, -y_0)

		opts.GeoM.Rotate(position.Angle)
		opts.GeoM.Scale(4, 4)
		opts.GeoM.Translate(position.X, position.Y)
		opts.GeoM.Translate(self.camera.X+x_0, self.camera.Y+y_0)

		// Render at this position
		screen.DrawImage(sprite, opts)
	}
	self.drawMinimap(screen)
}

func (self *ArenaScene) Update(dispatcher *scenes.Dispatcher) {
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

	for player := range query.Iter(self.simulation.ECS.World) {
		if self.playerId != component.Player.GetValue(player).Id {
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
func (self *ArenaScene) receiveServerUpdates() {
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

				player := self.simulation.FindCorrespondingPlayer(updatePosition.PlayerId)
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

				self.simulation.SpawnPlayer(playerConnected.PlayerId, &playerConnected.Position)
			}
		default:
		}
	}

}

func (self *ArenaScene) drawMinimap(screen *ebiten.Image) {
    // Create the minimap image
    minimap := ebiten.NewImage(MinimapWidth, MinimapHeight)
    minimap.Fill(color.Black)

    // Scale factor to map world coordinates to minimap coordinates
    scaleX := float64(MinimapWidth) / float64(MapWidth)
    scaleY := float64(MinimapHeight) / float64(MapHeight)

    // Draw the scaled-down game background on the minimap
    bgOpts := &ebiten.DrawImageOptions{}
    bgOpts.GeoM.Scale(scaleX, scaleY)
    minimap.DrawImage(self.background.Image, bgOpts)

    // Get the player's position to center the minimap around it
    playerPos := component.Position.Get(self.simulation.FindCorrespondingPlayer(self.playerId))

    // Calculate offsets to center the player in the minimap
    offsetX := playerPos.X*scaleX - float64(MinimapWidth)/2
    offsetY := playerPos.Y*scaleY - float64(MinimapHeight)/2

    // Draw all players on the minimap, with the main player centered
    query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))
    for player := range query.Iter(self.simulation.ECS.World) {
        position := component.Position.Get(player)

        // Calculate the position of each player relative to the centered player
        minimapX := (position.X * scaleX) - offsetX
        minimapY := (position.Y * scaleY) - offsetY

        // Check if the ship sprite is available
        playerSpritePtr := component.Sprite.Get(player)
        if playerSpritePtr != nil {
            playerSprite := *playerSpritePtr // Dereference to get *ebiten.Image

            // Scale down the ship for the minimap
            spriteScale := 0.5 // Adjust this scale factor to change ship size

            opts := &ebiten.DrawImageOptions{}
            opts.GeoM.Scale(spriteScale, spriteScale)
            opts.GeoM.Translate(minimapX, minimapY)
            minimap.DrawImage(playerSprite, opts)
        }
    }

    // Draw the minimap onto the main screen in the upper-left corner
    minimapScreenOpts := &ebiten.DrawImageOptions{}
    minimapScreenOpts.GeoM.Translate(10, 10) // Minimap position on the screen
    screen.DrawImage(minimap, minimapScreenOpts)

    // Draw a white border around the minimap
    borderColor := color.RGBA{255, 255, 255, 255} // White color for the border
    borderSize := 2

    // Top border
    topBorder := ebiten.NewImage(MinimapWidth+2*borderSize, borderSize)
    topBorder.Fill(borderColor)
    screen.DrawImage(topBorder, minimapScreenOpts)

    // Bottom border
    bottomBorder := ebiten.NewImage(MinimapWidth+2*borderSize, borderSize)
    bottomBorder.Fill(borderColor)
    bottomBorderOpts := *minimapScreenOpts
    bottomBorderOpts.GeoM.Translate(0, float64(MinimapHeight+borderSize))
    screen.DrawImage(bottomBorder, &bottomBorderOpts)

    // Left border
    leftBorder := ebiten.NewImage(borderSize, MinimapHeight)
    leftBorder.Fill(borderColor)
    screen.DrawImage(leftBorder, minimapScreenOpts)

    // Right border
    rightBorder := ebiten.NewImage(borderSize, MinimapHeight)
    rightBorder.Fill(borderColor)
    rightBorderOpts := *minimapScreenOpts
    rightBorderOpts.GeoM.Translate(float64(MinimapWidth+borderSize), 0)
    screen.DrawImage(rightBorder, &rightBorderOpts)
}

