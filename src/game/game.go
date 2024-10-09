package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"space-shooter/game/state"
)

type Game struct {
	screen_width  int
	screen_height int
	stateManager  state.StateManager
}

func NewGame(screen_width, screen_height int) Game {
	stateManager := state.NewStateManager(screen_width, screen_height)

	return Game{
		screen_width,
		screen_height,
		stateManager,
	}
}

func (g *Game) Update() error {
	g.stateManager.HandleUpdate()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.stateManager.Render(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
