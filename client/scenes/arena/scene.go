package arena

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"sort"
	"space-shooter/assets"
	"space-shooter/client/config"
	"space-shooter/client/scenes"
	"space-shooter/client/scenes/common"
	"space-shooter/game"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"space-shooter/rpc"
	"space-shooter/server/messages"
	"sync"
	"time"

	dmath "github.com/yohamta/donburi/features/math"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type LeaderboardEntry struct {
	Username string
	Score    int
}

type ArenaScene struct {
	background1 *common.Background
	background2 *common.Background
	config      *config.ClientConfig //testing only
	once        sync.Once            //testing only

	simulation     *game.GameSimulation
	shakeDuration  int
	shakeIntensity float64
	camera         *Camera

	lastFireTime time.Time

	connection *websocket.Conn
	player     *donburi.Entry
	playerName string
	playerId   types.PlayerId

	deathScene *DeathScene

	isAlive bool

	leaderboardData []LeaderboardEntry
	scrollOffset    int
}

func NewArenaScene(config *config.ClientConfig, playerName string) *ArenaScene {
	return &ArenaScene{
		background1:     common.NewBackground(game.MapWidth, game.MapHeight),
		background2:     common.NewBackground(config.ScreenWidth, config.ScreenHeight),
		playerName:      playerName,
		camera:          NewCamera(0, 0, game.MapHeight, game.MapWidth, config),
		deathScene:      NewDeathScene(config),
		isAlive:         true,
		config:          config,
		leaderboardData: []LeaderboardEntry{},
	}
}

func (self *ArenaScene) Configure(controller *scenes.AppController) error {
	controller.ChangeMusic(assets.BattleMusic)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	connection, _, err := websocket.Dial(ctx, self.config.ServerWebsocketURL, nil)

	if err != nil {
		return fmt.Errorf("Failed to connect to the server at %s", self.config.ServerWebsocketURL)
	}

	connectionHandshake := rpc.NewBaseMessage(messages.ConnectionHandshake{PlayerName: self.playerName})
	if err := rpc.WriteMessage(ctx, connection, connectionHandshake); err != nil {
		return fmt.Errorf("Failed to send handshake to the server at %s", self.config.ServerWebsocketURL)
	}

	var response messages.ConnectionHandshakeResponse
	if err := rpc.ReceiveExpectedMessage(ctx, connection, &response); err != nil {
		return fmt.Errorf("Error receiving handshake response: " + err.Error())
	}

	if response.IsRoomFull {
		return fmt.Errorf("Room is full")
	}

	self.connection = connection
	self.simulation = game.NewGameSimulation()

	self.simulation.OnBulletCollide = func(player, bullet *donburi.Entry) {
		if component.Player.Get(player).Id == self.playerId {
			self.startShake(10, 10)
		}
	}

	for _, player := range response.PlayerData {
		if player.PlayerId == response.PlayerId {
			// Focus the camera on the player.
			self.player = self.simulation.SpawnPlayer(player.PlayerId, &player.Position, player.PlayerName)
			self.playerId = player.PlayerId
			self.camera.FocusTarget(player.Position)
			continue
		}
		self.simulation.SpawnPlayer(player.PlayerId, &player.Position, player.PlayerName)
	}

	go self.receiveServerUpdates()
	return nil
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

	if !self.isAlive {
		self.deathScene.Draw(screen)
	}

	if ebiten.IsKeyPressed(ebiten.KeyL) {
		self.showLeaderboard(screen)
	}
}

func (self *ArenaScene) Update(controller *scenes.AppController) {
	if self.isAlive {
		self.handleInput()
	}

	self.simulation.Update()

	position := component.Position.Get(self.player)
	self.camera.FocusTarget(*position)
	self.camera.Constrain()
}

func (self *ArenaScene) handleInput() {
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
	drawSprite := func(position *component.PositionData, scale float64, angleOffset float64, offset dmath.Vec2, sprite *ebiten.Image) {
		// Center the texture.
		x0 := float64(sprite.Bounds().Dx()) / 2
		y0 := float64(sprite.Bounds().Dy()) / 2

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-x0, -y0)
		opts.GeoM.Translate(offset.X, offset.Y)

		opts.GeoM.Rotate(position.Angle)
		opts.GeoM.Rotate(angleOffset)

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

			if !player.IsAlive {
				continue
			}

			font := text.GoTextFace{Source: assets.Munro, Size: 20}
			width, _ := text.Measure(player.Name, &font, 12)

			x := (position.X - width/2) + 6
			y := position.Y - 55

			x += self.camera.X
			y += self.camera.Y

			// Set up the text drawing options
			opts := &text.DrawOptions{}
			opts.GeoM.Translate(x, y)

			text.Draw(screen, player.Name, &font, opts)
			self.drawHealthBar(screen, position, player.Health, 100)

			// Draw the player ship
			drawSprite(position, 4.0, 0, dmath.NewVec2(0, 0), component.Sprite.GetValue(entity))

			if player.Id != self.playerId {
				enemyPosition := component.Position.Get(entity)
				self.drawPointingArrow(screen, enemyPosition)
			} else {
				opts := &text.DrawOptions{}
				opts.GeoM.Translate(5, 5)
				text.Draw(screen, fmt.Sprintf("Score %d", player.Score), &text.GoTextFace{Source: assets.Munro, Size: 20}, opts)
			}

			if player.IsMovingForward {
				exhaust := component.Animation.Get(entity).Frame()
				drawSprite(position, 4.0, 0, dmath.NewVec2(0, 8), exhaust)
			}

		} else if entity.HasComponent(component.Explosion) {
			sprite := component.Animation.Get(entity).Frame()
			position := component.Position.GetValue(entity)
			explosion := component.Explosion.Get(entity)

			for i := 0; i < explosion.Count; i++ {
				position.X += 25 * rand.Float64()
				position.Y += 25 * rand.Float64()
				drawSprite(&position, 4.0, 0, dmath.NewVec2(0, 0), sprite)
			}
		} else if entity.HasComponent(component.Bullet) {
			drawSprite(position, 4.0, -math.Pi/4, dmath.NewVec2(0, 0), component.Sprite.GetValue(entity))
		}
	}
}

func (self *ArenaScene) drawPointingArrow(screen *ebiten.Image, enemyPosition *component.PositionData) {
	ourPosition := component.Position.Get(self.player)
	arrow := assets.Arrows.GetTile(assets.TileIndex{X: 9, Y: 12})

	vec := dmath.NewVec2(enemyPosition.X-ourPosition.X, enemyPosition.Y-ourPosition.Y)
	// do not draw the arrows if they are within the vicinity of the player
	if vec.Magnitude() < 500 {
		return
	}

	normalizedVec := vec.Normalized().MulScalar(100)
	angle := vec.Angle(dmath.NewVec2(1, 0))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(arrow.Bounds().Dx()), -float64(arrow.Bounds().Dy()))
	// The image is rotated by 90 degrees, undo that rotation.
	op.GeoM.Rotate(-math.Pi / 2)
	op.GeoM.Rotate(angle)
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(ourPosition.X, ourPosition.Y)
	op.GeoM.Translate(normalizedVec.X, normalizedVec.Y)
	op.GeoM.Translate(self.camera.X, self.camera.Y)

	screen.DrawImage(arrow, op)
}

func (self *ArenaScene) drawHealthBar(screen *ebiten.Image, position *component.PositionData, health float64, maxHealth float64) {
	if health <= 0 {
		return
	}
	healthBarWidth := 50.0
	healthBarHeight := 3.8

	// Calculate health percentage
	healthPercentage := health / maxHealth

	// Calculate world position for the health bar
	barX := (position.X + self.camera.X - healthBarWidth/2) + 5 // Center horizontally
	barY := position.Y + self.camera.Y - 26                     // Position above ship sprite with 26 offset

	tile := assets.Healthbar.GetTile(assets.TileIndex{X: 0, Y: 10})
	tileWidth := float64(assets.Healthbar.GetTile(assets.TileIndex{X: 0, Y: 10}).Bounds().Dx())

	x := position.X - healthBarWidth/2 + healthBarWidth/2 - tileWidth/2 - 18
	y := position.Y - 30

	x += self.camera.X
	y += self.camera.Y

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(4, 1.4)
	opts.GeoM.Translate(x, y)

	screen.DrawImage(tile, opts)

	// Draw the health bar background
	healthBarBackground := ebiten.NewImage(int(healthBarWidth), int(healthBarHeight))
	healthBarBackground.Fill(color.RGBA{128, 128, 128, 255}) // Light Gray for the background

	opts = &ebiten.DrawImageOptions{}
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

// Receives information from the server and updates the game state accordingly.
func (self *ArenaScene) receiveServerUpdates() {
	for {
		var message rpc.BaseMessage
		if err := rpc.ReceiveMessage(context.Background(), self.connection, &message); err != nil {
			continue
		}

		switch message.MessageType {
		case "UpdatePosition":
			var updatePosition messages.UpdatePosition
			if err := rpc.DecodeExpectedMessage(message, &updatePosition); err != nil {
				continue
			}
			if player := self.simulation.FindCorrespondingPlayer(updatePosition.PlayerId); player != nil {
				component.Position.SetValue(player, updatePosition.Position)
			}
		case "EventPlayerConnected":
			var eventPlayerConnected messages.EventPlayerConnected
			if err := rpc.DecodeExpectedMessage(message, &eventPlayerConnected); err != nil {
				continue
			}
			self.simulation.SpawnPlayer(eventPlayerConnected.PlayerId, &eventPlayerConnected.Position, eventPlayerConnected.PlayerName)
		case "EventPlayerMove":
			var eventPlayerMove messages.EventPlayerMove
			if err := rpc.DecodeExpectedMessage(message, &eventPlayerMove); err != nil {
				continue
			}
			self.simulation.RegisterPlayerMove(eventPlayerMove.PlayerId, eventPlayerMove.Move)
		case "EventUpdateHealth":
			var updateHealth messages.EventUpdateHealth
			if err := rpc.DecodeExpectedMessage(message, &updateHealth); err != nil {
				continue
			}
			self.simulation.UpdatePlayerHealth(updateHealth.PlayerId, updateHealth.Health)
		case "EventPlayerDied":
			var playerDied messages.EventPlayerDied
			if err := rpc.DecodeExpectedMessage(message, &playerDied); err != nil {
				continue
			}

			killed := self.simulation.FindCorrespondingPlayer(playerDied.PlayerId)
			killer := self.simulation.FindCorrespondingPlayer(playerDied.KilledBy)

			self.simulation.RegisterPlayerDeath(killed, killer)
			if playerDied.PlayerId == self.playerId {
				self.deathScene = NewDeathScene(self.config)
				self.isAlive = false
			}
		case "EventPlayerFireBullet":
			var fireBullet messages.EventPlayerFireBullet
			if err := rpc.DecodeExpectedMessage(message, &fireBullet); err != nil {
				continue
			}
			self.simulation.RegisterPlayerFire(self.simulation.FindCorrespondingPlayer(fireBullet.PlayerId))
		case "EventPlayerRespawned":
			var playerRespawn messages.EventPlayerRespawned
			if err := rpc.DecodeExpectedMessage(message, &playerRespawn); err != nil {
				continue
			}

			self.simulation.RespawnPlayer(self.simulation.FindCorrespondingPlayer(playerRespawn.PlayerId), playerRespawn.Position)
			self.isAlive = true
		default:
		}
	}
}

var leaderboardScrollOffset float64 = 0 // Global or struct-level variable

func (self *ArenaScene) scores(screen *ebiten.Image, playerData []component.PlayerData) {
	leaderboardEntries := make([]LeaderboardEntry, len(playerData))
	query := donburi.NewQuery(filter.Contains(component.Player))
	for entity := range query.Iter(self.simulation.ECS.World) {

		if entity.HasComponent(component.Player) {
			player := component.Player.Get(entity)

			for i, player := range playerData {
				leaderboardEntries[i] = LeaderboardEntry{
					Username: player.Name,
					Score:    player.Score,
				}
			}

			log.Print(leaderboardEntries)
			log.Print(player.Name, " ", player.Score)
		}
	}
}

func (self *ArenaScene) showLeaderboard(screen *ebiten.Image) {
	self.scores(screen, []component.PlayerData{})

	// Input handling for mouse scroll
	_, wheelY := ebiten.Wheel()            // Get the mouse scroll input
	leaderboardScrollOffset -= wheelY * 10 // Adjust offset based on scroll direction and speed
	visibleRows := 5
	startY := 290

	// Clamp the scroll offset to prevent out-of-bounds scrolling
	maxOffset := float64(len(self.playerName)*70 - self.config.ScreenHeight/2)
	if leaderboardScrollOffset < 0 {
		leaderboardScrollOffset = 0
	}
	if leaderboardScrollOffset > maxOffset {
		leaderboardScrollOffset = maxOffset
	}

	// Sort LeaderboardData
	sort.SliceStable(self.leaderboardData, func(i, j int) bool {
		return self.leaderboardData[i].Score > self.leaderboardData[j].Score
	})

	// Draw the Title box and Title
	opts1 := &ebiten.DrawImageOptions{}
	imageWidth := assets.Borders.Image.Bounds().Dx()
	opts1.GeoM.Scale(25, 7)
	opts1.GeoM.Translate((float64(self.config.ScreenWidth-imageWidth)/3)+55, 30)
	screen.DrawImage(assets.Borders.GetTile(assets.TileIndex{X: 1, Y: 0}), opts1)

	fontface := text.GoTextFace{Source: assets.MunroNarrow}
	lineSpacing := 10

	drawText(screen, "Leaderboard", fontface, 50, 550, 85, lineSpacing)

	// Draw the leaderboard box for the rankings
	drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 3}), 50, 32, 0, 150, 165, [4]float32{1, 1, 1, 1})
	drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 50, 32, 0, 150, 165, [4]float32{0.25, 0.25, 0.25, 1})

	// // Iterate and render leaderboard entries
	for i, entry := range self.leaderboardData {
		y := float64(startY) + float64(i*70) - leaderboardScrollOffset

		if y < 150 || y > float64(self.config.ScreenHeight)-50 {
			continue
		}

		drawTransformedImage(screen, assets.Borders.GetTile(assets.TileIndex{X: 0, Y: 1}), 38, 4, 0, 255, float64(y), [4]float32{0.25, 0.25, 0.25, 1})
		drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: (i % 2) + 7, Y: 0}), 7, 8, 0, 255, float64(y), [4]float32{0.8, 0.8, 0.8, 1})
		drawText(screen, fmt.Sprintf("%d", i+1), fontface, 35, 283, float64(y+32), lineSpacing)
		drawText(screen, entry.Username, fontface, 50, 440, float64(y+32), lineSpacing)
		drawText(screen, fmt.Sprintf("%d", entry.Score), fontface, 50, 740, float64(y+32), lineSpacing)
	}

	// Scroll indicators
	if self.scrollOffset > 0 {
		drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 8, Y: 0}), 7, 8, 0, 255, 250, [4]float32{1, 1, 1, 1})
	}
	if self.scrollOffset+visibleRows < len(self.leaderboardData) {
		drawTransformedImage(screen, assets.Arrows.GetTile(assets.TileIndex{X: 7, Y: 0}), 7, 8, 0, 255, 650, [4]float32{1, 1, 1, 1})
	}
}

// Helper function to draw centered text with specified font size
func drawText(screen *ebiten.Image, msg string, fontface text.GoTextFace, fontSize float64, x, y float64, lineSpacing int) {
	fontface.Size = fontSize
	width, height := text.Measure(msg, &fontface, 10)

	opts := &text.DrawOptions{}
	opts.LineSpacing = float64(lineSpacing)
	opts.GeoM.Translate(-width/2, -height/2)
	opts.GeoM.Translate(x, y)
	text.Draw(screen, msg, &fontface, opts)
}

func drawTransformedImage(screen *ebiten.Image, image *ebiten.Image, scaleX, scaleY, rotate, translateX, translateY float64, colorScale [4]float32) {
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
