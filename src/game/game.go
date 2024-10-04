package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	background    Sprite
	screen_width  int
	screen_height int
	player        Player
}

func NewGame(screen_width, screen_height int) Game {
	background := NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/BackGrounds.png").
		CreateTiles(CreateTilesInput{x_start: 0, y_start: 0, width: 128, height: 256, x_count: 3, y_count: 2}).
		FilterTiles(
			Tile{128, 0, 128, 256},
			Tile{256, 0, 128, 256},
			Tile{0, 256, 128, 256},
		).
		BuildAsBackgroundSprite(screen_width, screen_height, 128, 256)

	shipSpriteSheet := NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack/Ships.png").
		CreateTiles(CreateTilesInput{x_start: 0, y_start: 0, width: 8, height: 8, x_count: 3, y_count: 1}).
		BuildAsSpriteSheet()

	player := NewPlayer(shipSpriteSheet, 20, 100)

	return Game{
		screen_width:  screen_width,
		screen_height: screen_height,
		background:    background,
		player:        player,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.RotateClockwise()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.RotateCounterClockwise()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: Make sprites implement the image.Draw interface.
	screen.DrawImage(g.background.image, &ebiten.DrawImageOptions{})

	g.player.Render(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
