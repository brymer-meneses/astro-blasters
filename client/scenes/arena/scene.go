package arena

import (
	"context"
	"image/color"
	"log"
	"math/rand/v2"
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/client/scenes/leaderboard"
	"space-shooter/game"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"space-shooter/rpc"
	"space-shooter/server/messages"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

const (
	MapWidth      = 4096
	MapHeight     = 4096
	MinimapWidth  = 150
	MinimapHeight = 150
)

type ArenaScene struct {
	connection  *websocket.Conn
	simulation  *game.GameSimulation
	background1 *common.Background
	background2 *common.Background

	config *config.ClientConfig //testing only
	once   sync.Once            //testing only

	lastFireTime time.Time
	camera       *Camera

	shakeDuration  int
	shakeIntensity float64

	player   *donburi.Entry
	playerId types.PlayerId
}

func NewArenaScene(config *config.ClientConfig, playerName string) *ArenaScene {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	connection, _, err := websocket.Dial(ctx, config.ServerWebsocketURL, nil)

	if err != nil {
		log.Fatalf("Failed to connect to the game server at %s\n", config.ServerWebsocketURL)
	}

	connectionHandshake := rpc.NewBaseMessage(messages.ConnectionHandshake{PlayerName: playerName})
	if err := rpc.WriteMessage(ctx, connection, connectionHandshake); err != nil {
		log.Fatalf("Failed to connect to the game server at %s\n", config.ServerWebsocketURL)
	}

	var response messages.ConnectionHandshakeResponse
	if err := rpc.ReceiveExpectedMessage(ctx, connection, &response); err != nil {
		log.Fatal(err)
	}

	if response.IsRoomFull {
		log.Fatal("Room is full")
	}

	camera := NewCamera(0, 0, MapHeight, MapWidth, config)
	simulation := game.NewGameSimulation()
	var mainPlayer *donburi.Entry

	for _, player := range response.PlayerData {
		if player.PlayerId == response.PlayerId {
			// Focus the camera on the player.
			mainPlayer = simulation.SpawnPlayer(player.PlayerId, &player.Position, player.PlayerName)
			camera.FocusTarget(player.Position)
			continue
		}

		simulation.SpawnPlayer(player.PlayerId, &player.Position, player.PlayerName)
	}

	scene := &ArenaScene{
		background1: common.NewBackground(MapWidth, MapHeight),
		background2: common.NewBackground(config.ScreenWidth, config.ScreenHeight),
		playerId:    response.PlayerId,
		player:      mainPlayer,
		simulation:  simulation,
		connection:  connection,
		camera:      camera,
		config:      config,
	}

	go scene.receiveServerUpdates()
	return scene
}

func (self *ArenaScene) Draw(screen *ebiten.Image) {
	screen.Clear()

	if self.shakeDuration > 0 {
		self.camera.X += (rand.Float64()*2 - 1) * self.shakeIntensity
		self.camera.Y += (rand.Float64()*2 - 1) * self.shakeIntensity
		self.shakeDuration -= 1
	}

	self.drawBackground(screen)
	self.drawEntities(screen)
	self.drawMinimap(screen)

	// if player is dead (ui only) {
	// Make the entire screen gray with overlay
	// overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	// overlay.Fill(color.RGBA{64, 64, 64, 128})

	// opts := &ebiten.DrawImageOptions{}
	// screen.DrawImage(overlay, opts)

	// imageWidth := assets.Borders.Image.Bounds().Dx()
	// self.drawMessage(screen, assets.Messagebar.GetTile(assets.TileIndex{X: 0, Y: 8}), 35, 35, 0, float64(self.config.ScreenWidth-imageWidth)/3-158, float64(self.config.ScreenHeight)/3, [4]float32{1, 1, 1, 1})

	// fontface := text.GoTextFace{Source: assets.MunroNarrow}
	// lineSpacing := 10
	// self.drawText(screen, "You have been struck down. Respawn and fight again!", fontface, 35, 560, 380, lineSpacing, [4]float32{0, 0, 0, 1})
	// } //add backend part: 5 second countdown before respawn, can't move when dead, etc.
}

func (self *ArenaScene) Update(dispatcher *scenes.Dispatcher) {
	ctx := context.Background()
	position := component.Position.Get(self.player)

	sendMove := func(move types.PlayerMove) {
		message := rpc.NewBaseMessage(messages.RegisterPlayerMove{Move: move, Position: *position})
		rpc.WriteMessage(ctx, self.connection, message)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		sendMove(types.PlayerStartForward)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyW) || inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		sendMove(types.PlayerStopForward)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		sendMove(types.PlayerStartRotateClockwise)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyD) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		sendMove(types.PlayerStopRotateClockwise)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		sendMove(types.PlayerStartRotateCounterClockwise)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		sendMove(types.PlayerStopRotateCounterClockwise)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		sendMove(types.PlayerStartFireBullet)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		sendMove(types.PlayerStopFireBullet)
	}

	// for testing only
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		self.once.Do(
			func() {
				dispatcher.Dispatch(leaderboard.NewLeaderboardScene(self.config))
			})
	}

	self.simulation.Update()

	self.camera.FocusTarget(*position)
	self.camera.Constrain()
}

func (self *ArenaScene) startShake(duration int, intensity float64) {
	self.shakeDuration = duration
	self.shakeIntensity = intensity
}

// Draw the background.
func (self *ArenaScene) drawBackground(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(2.0, 2.0)
	opts.GeoM.Translate(0.2*self.camera.X, 0.2*self.camera.Y)
	opts.ColorScale.Scale(1, 1, 1, 0.2)
	screen.DrawImage(self.background2.Image, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(self.camera.X, self.camera.Y)
	screen.DrawImage(self.background1.Image, opts)
}

func (self *ArenaScene) drawEntities(screen *ebiten.Image) {
	drawSprite := func(position *component.PositionData, scale float64, sprite *ebiten.Image) {
		// Center the texture.
		x0 := float64(sprite.Bounds().Dx()) / 2
		y0 := float64(sprite.Bounds().Dy()) / 2

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-x0, -y0)

		opts.GeoM.Rotate(position.Angle)
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(position.X, position.Y)
		opts.GeoM.Translate(self.camera.X+x0, self.camera.Y+y0)

		screen.DrawImage(sprite, opts)
	}

	query := donburi.NewQuery(filter.Contains(component.Position))
	for entity := range query.Iter(self.simulation.ECS.World) {
		position := component.Position.Get(entity)

		if entity.HasComponent(component.Player) {
			player := component.Player.Get(entity)

			font := text.GoTextFace{Source: assets.Munro, Size: 20}
			width, _ := text.Measure(player.Name, &font, 12)

			// Calculate the position of the text to center it above the health bar
			x := (position.X - width/2) + 6 // Center horizontally
			y := position.Y - 55            // Above the health bar

			// Apply camera translation
			x += self.camera.X
			y += self.camera.Y

			// Set up the text drawing options
			opts := &text.DrawOptions{}
			opts.GeoM.Translate(x, y)

			text.Draw(screen, player.Name, &font, opts)

			self.drawHealthBar(screen, position, player.Health, 100)
			drawSprite(position, 4.0, component.Sprite.GetValue(entity))
		} else if entity.HasComponent(component.Explosion) {
			sprite := component.Animation.Get(entity).Frame()
			position := component.Position.GetValue(entity)
			explosion := component.Explosion.Get(entity)

			for i := 0; i < explosion.Count; i++ {
				position.X += 25 * rand.Float64()
				position.Y += 25 * rand.Float64()
				drawSprite(&position, 4.0, sprite)
			}
		} else if entity.HasComponent(component.Bullet) {
			drawSprite(position, 4.0, component.Animation.Get(entity).Frame())
		}
	}
}

func (self *ArenaScene) drawTransformedImage(screen *ebiten.Image, tile *ebiten.Image, position *component.PositionData, healthBarWidth float64) {
	tileWidth := float64(tile.Bounds().Dx())

	x := position.X - healthBarWidth/2 + healthBarWidth/2 - tileWidth/2 - 18 // Center horizontally
	y := position.Y - 30                                                     // Position above the health bar

	// Apply camera offsets
	x += self.camera.X
	y += self.camera.Y

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(4, 1.4)
	opts.GeoM.Translate(x, y)

	// Draw the tile on the screen
	screen.DrawImage(tile, opts)
}

func (self *ArenaScene) drawHealthBar(screen *ebiten.Image, position *component.PositionData, health float64, maxHealth float64) {
	healthBarWidth := 50.0
	healthBarHeight := 3.8

	// Calculate health percentage
	healthPercentage := health / maxHealth

	// Calculate world position for the health bar
	barX := (position.X + self.camera.X - healthBarWidth/2) + 5 // Center horizontally
	barY := position.Y + self.camera.Y - 26                     // Position above ship sprite with 26 offset

	self.drawTransformedImage(screen, assets.Healthbar.GetTile(assets.TileIndex{X: 0, Y: 10}), position, healthBarWidth)

	// Draw the health bar background
	healthBarBackground := ebiten.NewImage(int(healthBarWidth), int(healthBarHeight))
	healthBarBackground.Fill(color.RGBA{128, 128, 128, 255}) // Light Gray for the background
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(barX, barY)
	screen.DrawImage(healthBarBackground, opts)

	// Draw the health bar foreground (current health)
	currentHealthWidth := healthBarWidth * healthPercentage
	healthBarForeground := ebiten.NewImage(int(currentHealthWidth), int(healthBarHeight))
	healthBarForeground.Fill(color.RGBA{0, 255, 0, 255}) // Green for the current health
	screen.DrawImage(healthBarForeground, opts)
}

func (self *ArenaScene) drawMessage(screen *ebiten.Image, image *ebiten.Image, scaleX, scaleY, rotate, translateX, translateY float64, colorScale [4]float32) {
	opts := &ebiten.DrawImageOptions{}

	// Apply geometric transformations
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Rotate(rotate) // Rotation in radians
	opts.GeoM.Translate(translateX, translateY)

	// Apply color transformations using ColorScale
	if len(colorScale) == 4 { // Ensure proper length (R, G, B, A)
		opts.ColorScale.Scale(colorScale[0], colorScale[1], colorScale[2], colorScale[3])
	}

	screen.DrawImage(image, opts)
}

// Helper function to draw centered text with specified font size
func (self *ArenaScene) drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int, colorScale [4]float32) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)

	// Apply color transformations using ColorScale
	if len(colorScale) == 4 { // Ensure proper length (R, G, B, A)
		opts.ColorScale.Scale(colorScale[0], colorScale[1], colorScale[2], colorScale[3])
	}

	text.Draw(screen, msg, &fontface, opts)
}

func (self *ArenaScene) drawMinimap(screen *ebiten.Image) {
	// Create the minimap image
	minimap := ebiten.NewImage(MinimapWidth, MinimapHeight)
	minimap.Fill(color.Black)

	// Scale factor to map world coordinates to minimap coordinates
	scaleX := float64(MinimapWidth) / float64(MapWidth)
	scaleY := float64(MinimapHeight) / float64(MapHeight)

	// Get the player's position to center the minimap around it
	playerPos := component.Position.Get(self.simulation.FindCorrespondingPlayer(self.playerId))

	// Calculate offsets to center the player in the minimap
	offsetX := playerPos.X*scaleX - float64(MinimapWidth)/2
	offsetY := playerPos.Y*scaleY - float64(MinimapHeight)/2

	spriteScale := 1.0

	// Draw all players on the minimap, with the main player centered
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))
	for player := range query.Iter(self.simulation.ECS.World) {
		position := component.Position.Get(player)
		sprite := component.Sprite.GetValue(player)

		// Calculate the position of each player relative to the centered player
		minimapX := (position.X * scaleX) - offsetX
		minimapY := (position.Y * scaleY) - offsetY

		// Check if the ship sprite is available
		x_0 := (float64(sprite.Bounds().Dx()) / 2)
		y_0 := (float64(sprite.Bounds().Dy()) / 2)

		// Scale down the ship for the minimap
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-x_0, -y_0)
		opts.GeoM.Rotate(position.Angle)

		opts.GeoM.Scale(spriteScale, spriteScale)
		opts.GeoM.Translate(minimapX, minimapY)
		minimap.DrawImage(sprite, opts)
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
	leftBorder := ebiten.NewImage(borderSize, MinimapHeight+2*borderSize)
	leftBorder.Fill(borderColor)
	screen.DrawImage(leftBorder, minimapScreenOpts)

	// Right border
	rightBorder := ebiten.NewImage(borderSize, MinimapHeight+2*borderSize)
	rightBorder.Fill(borderColor)
	rightBorderOpts := *minimapScreenOpts
	rightBorderOpts.GeoM.Translate(float64(MinimapWidth+borderSize), 0)
	screen.DrawImage(rightBorder, &rightBorderOpts)
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
				if err := rpc.DecodeExpectedMessage(message, &updatePosition); err != nil {
					continue
				}
				if player := self.simulation.FindCorrespondingPlayer(updatePosition.PlayerId); player != nil {
					component.Position.SetValue(player, updatePosition.Position)
				}
			}
		case "EventPlayerConnected":
			{
				var eventPlayerConnected messages.EventPlayerConnected
				if err := rpc.DecodeExpectedMessage(message, &eventPlayerConnected); err != nil {
					continue
				}
				self.simulation.SpawnPlayer(eventPlayerConnected.PlayerId, &eventPlayerConnected.Position, eventPlayerConnected.PlayerName)
			}
		case "EventPlayerMove":
			{
				var eventPlayerMove messages.EventPlayerMove
				if err := rpc.DecodeExpectedMessage(message, &eventPlayerMove); err != nil {
					continue
				}
				self.simulation.RegisterPlayerMove(eventPlayerMove.PlayerId, eventPlayerMove.Move)
			}
		default:
		}
	}
}
