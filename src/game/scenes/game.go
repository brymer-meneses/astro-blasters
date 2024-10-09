package scenes

import (
	"space-shooter/assets"
	"space-shooter/config"
	"space-shooter/game/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene struct {
	assetManager *assets.AssetManager
	players      []entities.Player
}

func NewGameScene(config *config.AppConfig, assetManager *assets.AssetManager) GameScene {
	players := make([]entities.Player, 5)
	for i := 0; i < 5; i++ {
		players[i] = entities.NewPlayer(&assetManager.Ships[i])
	}

	return GameScene{assetManager, players}
}

func (gs *GameScene) Render(screen *ebiten.Image) {
	gs.assetManager.Background.Render(screen)
	gs.players[0].Render(screen)
}

func (gs *GameScene) HandleUpdate() {

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		gs.players[0].MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		gs.players[0].RotateClockwise()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		gs.players[0].RotateCounterClockwise()
	}
}
