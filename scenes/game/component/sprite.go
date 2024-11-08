package component

import (
	"space-shooter/assets"

	"github.com/yohamta/donburi"
)

var Sprite = donburi.NewComponentType[assets.Sprite]()
