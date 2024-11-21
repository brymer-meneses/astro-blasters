package game

import (
	"math"
	"math/rand"
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
	for expirable := range donburi.NewQuery(filter.Contains(component.Expirable)).Iter(self.ECS.World) {
		expirableData := component.Expirable.GetValue(expirable)
		if time.Now().After(expirableData.ExpiresWhen) {
			self.ECS.World.Remove(expirable.Entity())
		}
	}

	for bullet := range donburi.NewQuery(filter.Contains(component.Bullet)).Iter(self.ECS.World) {
		futureBulletPosition := component.Position.GetValue(bullet)
		futureBulletPosition.Forward(-10)
		didCollide := false

		for player := range donburi.NewQuery(filter.Contains(component.Player)).Iter(self.ECS.World) {
			if component.Position.Get(player).IntersectsWith(&futureBulletPosition, 10) {
				didCollide = true
			}
		}

		if didCollide {
			self.spawnExplosion(&futureBulletPosition)
			self.ECS.World.Remove(bullet.Entity())
		} else {
			component.Position.SetValue(bullet, futureBulletPosition)
		}
	}
}

func (self *GameSimulation) FireBullet(playerId types.PlayerId) *donburi.Entry {
	player := self.FindCorrespondingPlayer(playerId)
	playerPosition := component.Position.Get(player)
	playerPosition.Forward(-3)

	bulletPosition := *playerPosition
	bulletPosition.Angle += math.Pi
	bulletPosition.Forward(-40)

	entity := self.ECS.World.Create(component.Bullet, component.Animation, component.Position, component.Expirable)
	bullet := self.ECS.World.Entry(entity)

	component.Bullet.SetValue(
		bullet,
		component.BulletData{
			FiredBy: playerId,
		},
	)
	component.Position.SetValue(
		bullet,
		bulletPosition,
	)
	component.Expirable.SetValue(
		bullet,
		component.NewExpirable(time.Second),
	)
	animationIndex := rand.Intn(len(assets.OrangeBulletAnimation))
	component.Animation.SetValue(
		bullet,
		component.NewAnimationData(assets.OrangeBulletAnimation[animationIndex], 5),
	)

	return bullet
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

func (self *GameSimulation) spawnExplosion(position *component.PositionData) {
	world := self.ECS.World
	entity := world.Create(component.Position, component.Explosion, component.Animation, component.Expirable)
	explosion := world.Entry(entity)

	component.Explosion.SetValue(
		explosion,
		component.ExplosionData{
			Count: rand.Intn(3),
		},
	)
	component.Position.SetValue(
		explosion,
		*position,
	)
	component.Animation.SetValue(
		explosion,
		component.NewAnimationData(assets.BlueExplosion, 2),
	)
	component.Expirable.SetValue(
		explosion,
		component.NewExpirable(2*time.Second),
	)
}

func getShipSprite(playerId types.PlayerId) *ebiten.Image {
	i := int(playerId)
	return assets.Ships.GetTile(assets.TileIndex{X: 1, Y: i})
}
