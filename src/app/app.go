package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"space-shooter/app/state"
)

type App struct {
	screen_width  int
	screen_height int
	stateManager  state.StateManager
}

func NewApp(screen_width, screen_height int) App {
	stateManager := state.NewStateManager(screen_width, screen_height)

	return App{
		screen_width,
		screen_height,
		stateManager,
	}
}

func (g *App) Update() error {
	g.stateManager.HandleUpdate()

	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	g.stateManager.Render(screen)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
