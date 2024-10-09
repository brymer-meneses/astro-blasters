package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"space-shooter/game/sprite"
)

type GameState struct {
	background sprite.Sprite
	player     Player
}

func NewGameState(screen_width, screen_height int) GameState {
	background := sprite.NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/BackGrounds.png").
		CreateTiles(sprite.CreateTilesInput{X_start: 0, Y_start: 0, Width: 128, Height: 256, X_count: 3, Y_count: 2}).
		FilterTiles(
			sprite.Tile{X: 128, Y: 0, Width: 128, Height: 256},
			sprite.Tile{X: 256, Y: 0, Width: 128, Height: 256},
			sprite.Tile{X: 0, Y: 256, Width: 128, Height: 256},
		).
		BuildAsBackgroundSprite(screen_width, screen_height, 128, 256)

	shipSpriteSheet := sprite.NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/Ships.png").
		CreateTiles(sprite.CreateTilesInput{X_start: 0, Y_start: 0, Width: 8, Height: 8, X_count: 3, Y_count: 1}).
		BuildAsSpriteSheet()

	player := NewPlayer(shipSpriteSheet, 20, 100)

	return GameState{background: background, player: player}
}

func (gs *GameState) Render(screen *ebiten.Image) {
	gs.background.Render(screen)
	gs.player.Render(screen)
}

func (gs *GameState) HandleUpdate() {

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		gs.player.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		gs.player.RotateClockwise()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		gs.player.RotateCounterClockwise()
	}
}
