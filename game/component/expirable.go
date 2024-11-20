package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type ExpirableData struct {
	ExpiresWhen time.Time
}

func NewExpirable(duration time.Duration) ExpirableData {
	return ExpirableData{ExpiresWhen: time.Now().Add(duration)}
}

var Expirable = donburi.NewComponentType[ExpirableData]()
