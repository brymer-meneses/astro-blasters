package main

import (
	"log"
	"space-shooter"
	"space-shooter/config"
)

func main() {
	app := app.NewApp(config.AppConfig{ScreenWidth: 1080, ScreenHeight: 720})
	if err := app.RunApp(); err != nil {
		log.Fatal(err)
	}
}
