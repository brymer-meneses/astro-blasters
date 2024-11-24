package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"log"
	"math"
	"math/rand"
	"space-shooter/assets"
	"space-shooter/game/component"
	"space-shooter/game/types"
	"time"
)

const (
	PlayerDamagePerHit  = 5
	PlayerMovementSpeed = 5
	PlayerRotationSpeed = 5
	MapWidth            = 4096
	MapHeight           = 4096

	ShipWidth  = 32
	ShipHeight = 32
)

type GameSimulation struct {
	ECS       *ecs.ECS
	OnCollide func(player *donburi.Entry)
}

func NewGameSimulation(onCollide func(player *donburi.Entry)) *GameSimulation {
	return &GameSimulation{
		ECS:       ecs.NewECS(donburi.NewWorld()),
		OnCollide: onCollide,
	}
}

func (self *GameSimulation) UpdatePlayerHealth(playerId types.PlayerId, health float64) {
	player := self.FindCorrespondingPlayer(playerId)
	playerData := component.Player.Get(player)
	playerData.Health = health
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
			if component.Position.Get(player).IntersectsWith(&futureBulletPosition, 20) {
				didCollide = true
				self.OnCollide(player)
			}
		}

		if didCollide {
			self.spawnExplosion(&futureBulletPosition)
			self.ECS.World.Remove(bullet.Entity())
		} else {
			component.Position.SetValue(bullet, futureBulletPosition)
		}
	}

	for player := range donburi.NewQuery(filter.Contains(component.Player)).Iter(self.ECS.World) {
		playerData := component.Player.Get(player)

		if playerData.IsFiringBullet {
			self.fireBullet(player)
		}

		futurePosition := component.Position.GetValue(player)
		if playerData.IsMovingForward {
			futurePosition.Forward(PlayerMovementSpeed)
		}

		if playerData.IsRotatingClockwise {
			futurePosition.Rotate(PlayerRotationSpeed)
		}

		if playerData.IsRotatingCounterClockwise {
			futurePosition.Rotate(-PlayerRotationSpeed)
		}

		if futurePosition.X < ShipWidth || futurePosition.X > MapWidth-ShipWidth {
			continue
		}
		if futurePosition.Y < ShipHeight || futurePosition.Y > MapHeight-ShipHeight {
			continue
		}

		component.Position.SetValue(player, futurePosition)
	}
}

func (self *GameSimulation) RegisterPlayerMove(playerId types.PlayerId, move types.PlayerMove) {
	player := self.FindCorrespondingPlayer(playerId)
	if player == nil {
		log.Fatal("Invalid player Id")
	}

	playerData := component.Player.Get(player)
	switch move {
	case types.PlayerStartFireBullet:
		playerData.IsFiringBullet = true
	case types.PlayerStopFireBullet:
		playerData.IsFiringBullet = false

	case types.PlayerStartForward:
		playerData.IsMovingForward = true
	case types.PlayerStopForward:
		playerData.IsMovingForward = false

	case types.PlayerStartRotateClockwise:
		playerData.IsRotatingClockwise = true
	case types.PlayerStopRotateClockwise:
		playerData.IsRotatingClockwise = false

	case types.PlayerStartRotateCounterClockwise:
		playerData.IsRotatingCounterClockwise = true
	case types.PlayerStopRotateCounterClockwise:
		playerData.IsRotatingCounterClockwise = false
	}
}

func (self *GameSimulation) fireBullet(player *donburi.Entry) *donburi.Entry {
	playerData := component.Player.Get(player)

	playerPosition := component.Position.Get(player)
	playerPosition.Forward(-3)

	bulletPosition := *playerPosition
	bulletPosition.Angle += math.Pi
	bulletPosition.Forward(-40)

	entity := self.ECS.World.Create(component.Bullet, component.Sprite, component.Position, component.Expirable)
	bullet := self.ECS.World.Entry(entity)

	component.Bullet.SetValue(
		bullet,
		component.BulletData{
			FiredBy: playerData.Id,
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
	component.Sprite.SetValue(
		bullet,
		assets.Bullet,
	)

	return bullet
}

func (self *GameSimulation) RespawnPlayer(player *donburi.Entry, newPosition component.PositionData) {
	playerData := component.Player.Get(player)
	playerData.Health = 100
	component.Position.SetValue(player, newPosition)
}

func (self *GameSimulation) SpawnPlayer(playerId types.PlayerId, position *component.PositionData, playerName string) *donburi.Entry {
	entity := self.ECS.World.Create(component.Player, component.Position, component.Animation, component.Sprite)
	player := self.ECS.World.Entry(entity)

	// Create and store player metadata
	playerData := component.PlayerData{
		Name:   playerName,
		Id:     playerId,
		Health: 100,
	}
	component.Player.SetValue(player, playerData)
	component.Position.SetValue(player, *position)
	component.Sprite.SetValue(player, getShipSprite(playerId))
	component.Animation.SetValue(player, component.NewAnimationData(assets.OrangeExhaustAnimation[0], 5))

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

func generateRandomFloat(min, max float64) float64 {
	return max*rand.Float64() + min
}

func GenerateRandomPlayerPosition() component.PositionData {
	return component.PositionData{
		X:     generateRandomFloat(ShipWidth, MapHeight-ShipHeight),
		Y:     generateRandomFloat(ShipHeight, MapHeight-ShipHeight),
		Angle: generateRandomFloat(0, 1),
	}
}
