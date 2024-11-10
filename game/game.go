package game

import (
	"space-shooter/assets"
	"space-shooter/game/component"
	"space-shooter/game/types"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameSimulation struct {
	ECS *ecs.ECS
}

func NewGameSimulation() *GameSimulation {
	return &GameSimulation{
		ECS: ecs.NewECS(donburi.NewWorld()),
	}
}

func (self *GameSimulation) Update() {}

func (self *GameSimulation) SpawnPlayer(playerId types.PlayerId, position *component.PositionData) {
	world := self.ECS.World
	entity := world.Create(component.Player, component.Position, component.Sprite)
	player := world.Entry(entity)

	component.Player.SetValue(
		player,
		component.PlayerData{
			Name: "Player One",
			Id:   playerId,
		},
	)
	component.Position.SetValue(
		player,
		*position,
	)
	component.Sprite.SetValue(
		player,
		getShipSprite(playerId),
	)
}

// Returns the ecs entry given the playerId.
func (self *GameSimulation) FindCorrespondingPlayer(playerId types.PlayerId) *donburi.Entry {
	query := donburi.NewQuery(filter.Contains(component.Player))
	for player := range query.Iter(self.ECS.World) {
		if playerId == component.Player.GetValue(player).Id {
			return player
		}
	}
	return nil
}

func getShipSprite(playerId types.PlayerId) *ebiten.Image {
	i := int(playerId)
	return assets.Ships.GetTile(assets.SpriteTile{X0: 1, Y0: i, X1: 2, Y1: i + 1})
}
