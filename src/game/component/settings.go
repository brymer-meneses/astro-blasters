package component

import "github.com/yohamta/donburi"

type SettingsData struct {
	PlayerId PlayerId
}

var Settings = donburi.NewComponentType[SettingsData]()
