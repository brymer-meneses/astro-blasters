package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"space-shooter/game/component"
)

func PlayerMovement(ecs *ecs.ECS) {
	donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))
	query := donburi.NewQuery(filter.Contains(component.Player, component.Position, component.Sprite))

	entry, _ := component.Settings.First(ecs.World)
	settings := component.Settings.Get(entry)

	for player := range query.Iter(ecs.World) {
		if settings.PlayerId != component.Player.GetValue(player).Id {
			continue
		}

		positionData := component.Position.Get(player)
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			positionData.Forward()
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			positionData.RotateClockwise()
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			positionData.RotateCounterClockwise()
		}
	}
}
