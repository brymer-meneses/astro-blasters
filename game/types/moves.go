package types

type PlayerMove int64

const (
	PlayerIdle = iota

	PlayerStartForward
	PlayerStartRotateClockwise
	PlayerStartRotateCounterClockwise

	PlayerStopForward
	PlayerStopRotateClockwise
	PlayerStopRotateCounterClockwise

	PlayerStartFireBullet
	PlayerStopFireBullet
)
