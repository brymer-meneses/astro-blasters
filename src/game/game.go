package game

import (
	"fmt"

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
		FromFile("../assets/SpaceShooterAssetPack_BackGrounds.png").
		CreateTiles(CreateTilesInput{x_start: 0, y_start: 0, width: 128, height: 256, x_count: 3, y_count: 2}).
		FilterTiles([]Tile{
			{128, 0, 128, 256},
			{256, 0, 128, 256},
			{0, 256, 128, 256},
		}).
		BuildAsBackgroundSprite(screen_width, screen_height, 128, 256)

	shipSpriteSheet := NewSpriteBuilder().
		FromFile("../assets/SpaceShooterAssetPack_Ships.png").
		CreateTiles(CreateTilesInput{x_start: 0, y_start: 0, width: 8, height: 8, x_count: 3, y_count: 1}).
		BuildAsSpriteSheet()

	fmt.Printf("ship %d\n", len(shipSpriteSheet.sprites))

	player := NewPlayer(shipSpriteSheet, 10, 0)

	return Game{
		screen_width:  screen_width,
		screen_height: screen_height,
		background:    background,
		player:        player,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MoveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MoveRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.MoveDown()
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
