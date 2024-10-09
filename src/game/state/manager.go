package state

import "github.com/hajimehoshi/ebiten/v2"

const (
	IN_GAME = iota
	IN_MENU
)

type StateManager struct {
	menu  MenuState
	game  GameState
	state int
}

func NewStateManager(screen_width, screen_height int) StateManager {
	game := NewGameState(screen_width, screen_height)
	menu := NewMenuState(screen_width, screen_height)

	state := IN_MENU

	return StateManager{menu, game, state}
}

func (sm *StateManager) HandleUpdate() {
	switch sm.state {
	case IN_GAME:
		sm.game.HandleUpdate()
		break
	case IN_MENU:
		sm.menu.HandleUpdate()
		break
	}
}

func (sm *StateManager) Render(screen *ebiten.Image) {
	switch sm.state {
	case IN_GAME:
		sm.game.Render(screen)
		break
	case IN_MENU:
		sm.menu.Render(screen)
		break
	}

}
