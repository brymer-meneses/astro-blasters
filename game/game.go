package game

import (
	"log"
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

const (
	PlayerDamagePerHit  = 5
	PlayerMovementSpeed = 5
	PlayerRotationSpeed = 5

	BulletSpeed = 20

	MapWidth  = 4096
	MapHeight = 4096

	ShipWidth  = 32
	ShipHeight = 32
)

type GameSimulation struct {
	ECS             *ecs.ECS
	OnBulletCollide func(player *donburi.Entry, bullet *donburi.Entry)
	OnBulletFire    func(player *donburi.Entry)
}

func NewGameSimulation() *GameSimulation {
	return &GameSimulation{
		ECS:             ecs.NewECS(donburi.NewWorld()),
		OnBulletCollide: func(player *donburi.Entry, bullet *donburi.Entry) {},
		OnBulletFire:    func(player *donburi.Entry) {},
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
		futureBulletPosition.Forward(-BulletSpeed)

		didCollide := false
		var collidedPlayer *donburi.Entry

		for player := range donburi.NewQuery(filter.Contains(component.Player)).Iter(self.ECS.World) {
			playerData := component.Player.Get(player)
			isDamageable := playerData.IsAlive && playerData.IsConnected

			if isDamageable && component.Position.Get(player).IntersectsWith(&futureBulletPosition, 20) {
				didCollide = true
				collidedPlayer = player
			}
		}

		if !didCollide {
			component.Position.SetValue(bullet, futureBulletPosition)
			continue
		}

		if self.OnBulletCollide != nil {
			self.OnBulletCollide(collidedPlayer, bullet)
		}

		self.spawnExplosion(&futureBulletPosition)
		self.ECS.World.Remove(bullet.Entity())
	}

	for player := range donburi.NewQuery(filter.Contains(component.Player)).Iter(self.ECS.World) {
		playerData := component.Player.Get(player)

		if playerData.IsFiringBullet {
			self.OnBulletFire(player)
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

func (self *GameSimulation) UpdatePlayerHealth(playerId types.PlayerId, health float64) {
	player := self.FindCorrespondingPlayer(playerId)
	playerData := component.Player.Get(player)
	playerData.Health = health
}

func (self *GameSimulation) RegisterPlayerDisconnection(player *donburi.Entry) {
	playerData := component.Player.Get(player)
	playerData.IsConnected = false
}

func (self *GameSimulation) RegisterPlayerDeath(victim, killer *donburi.Entry) {
	killerData := component.Player.Get(killer)
	killerData.Score += 10

	victimData := component.Player.Get(victim)
	victimData.Score /= 2
	victimData.IsFiringBullet = false
	victimData.IsRotatingClockwise = false
	victimData.IsMovingForward = false
	victimData.IsRotatingCounterClockwise = false
	victimData.IsAlive = false
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

func (self *GameSimulation) RegisterPlayerFire(player *donburi.Entry) {
	playerPosition := component.Position.Get(player)

	bullet1 := *playerPosition
	bullet1.Angle += math.Pi
	bullet1.X -= 15 * math.Cos(bullet1.Angle)
	bullet1.Y -= 15 * math.Sin(bullet1.Angle)
	bullet1.Forward(-40)

	bullet2 := *playerPosition
	bullet2.Angle += math.Pi
	bullet2.X += 15 * math.Cos(bullet2.Angle)
	bullet2.Y += 15 * math.Sin(bullet2.Angle)
	bullet2.Forward(-40)

	self.FireBullet(player, bullet1)
	self.FireBullet(player, bullet2)
}

func (self *GameSimulation) FireBullet(player *donburi.Entry, bulletPosition component.PositionData) *donburi.Entry {
	playerData := component.Player.Get(player)

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
	playerData.IsAlive = true
	component.Position.SetValue(player, newPosition)
}

func (self *GameSimulation) CreatePlayer(playerId types.PlayerId, position *component.PositionData, playerName string, IsConnected bool) *donburi.Entry {
	entity := self.ECS.World.Create(component.Player, component.Position, component.Animation, component.Sprite)
	player := self.ECS.World.Entry(entity)

	playerData := component.PlayerData{
		Name:        playerName,
		Id:          playerId,
		Health:      100,
		IsAlive:     true,
		IsConnected: IsConnected,
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

func GenerateRandomPlayerPosition() component.PositionData {
	return component.PositionData{
		X:     generateRandomFloat(ShipWidth, 0.80*MapWidth),
		Y:     generateRandomFloat(ShipHeight, 0.80*MapHeight),
		Angle: generateRandomFloat(0, 1),
	}
}

func getShipSprite(playerId types.PlayerId) *ebiten.Image {
	i := int(playerId)
	return assets.Ships.GetTile(assets.TileIndex{X: 1, Y: i % 5})
}

func generateRandomFloat(min, max float64) float64 {
	return max*rand.Float64() + min
}
