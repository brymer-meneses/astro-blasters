package game

import (
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/game/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IN_GAME = iota
	IN_MENU
)

type SceneManager struct {
	menu  scenes.MenuScene
	game  scenes.GameScene
	state int
}

func NewStateManager(config *config.AppConfig, assetManager *assets.AssetManager) SceneManager {
	game := scenes.NewGameScene(config, assetManager)
	menu := scenes.NewMenuScene(config, assetManager)

	state := IN_MENU

	return SceneManager{menu, game, state}
}

func (self *SceneManager) HandleUpdate() {

	switch self.state {
	case IN_GAME:
		self.game.HandleUpdate()
		break
	case IN_MENU:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			self.state = IN_GAME
		} else {
			self.menu.HandleUpdate()
		}

		break
	}
}

func (sm *SceneManager) Render(screen *ebiten.Image) {
	switch sm.state {
	case IN_GAME:
		sm.game.Render(screen)
		break
	case IN_MENU:
		sm.menu.Render(screen)
		break
	}

}
