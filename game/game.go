package game

import (
	"math"
	"space-shooter/assets"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"time"

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

func (self *GameSimulation) Update() {
	component.Bullet.Each(self.ECS.World, func(bullet *donburi.Entry) {
		position := component.Position.Get(bullet)
		position.Forward(-10)
	})
}

func (self *GameSimulation) FireBullet(playerId types.PlayerId) {
	player := self.FindCorrespondingPlayer(playerId)
	position := component.Position.GetValue(player)
	position.Angle += math.Pi
	position.Forward(-40)

	entity := self.ECS.World.Create(component.Bullet, component.Animation, component.Position)
	bullet := self.ECS.World.Entry(entity)

	component.Bullet.SetValue(
		bullet,
		component.BulletData{
			FiredBy:  playerId,
			ShotWhen: time.Now(),
		},
	)
	component.Position.SetValue(
		bullet,
		position,
	)
	component.Animation.SetValue(
		bullet,
		component.NewAnimationData(assets.OrangeBulletAnimation[0], 5),
	)
}

func (self *GameSimulation) SpawnPlayer(playerId types.PlayerId, position *component.PositionData) *donburi.Entry {
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

	return player
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
	return assets.Ships.GetTile(assets.TileIndex{X: 1, Y: i})
}
