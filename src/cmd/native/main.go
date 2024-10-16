package main

import (
	"log"
	"space-shooter"
	"space-shooter/config"
)

func main() {

	config := config.AppConfig{
		ScreenWidth:        1080,
		ScreenHeight:       720,
		ServerWebsocketURL: "ws://localhost:8080/events/ws",
	}

	app := app.NewApp(&config)
	if err := app.RunApp(); err != nil {
		log.Fatal(err)
	}
}
