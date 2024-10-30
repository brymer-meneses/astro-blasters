package main

import (
	"flag"
	"fmt"
	"log"
	"space-shooter"
	"space-shooter/config"
)

func main() {
	var port = flag.Int("port", 8080, "Port of the IP Address")
	var address = flag.String("address", "localhost", "IP Address")

	flag.Parse()

	url := fmt.Sprintf("ws://%s:%d/events/ws", *address, *port)

	config := config.AppConfig{
		ScreenWidth:        1080,
		ScreenHeight:       720,
		ServerWebsocketURL: url,
	}

	app := app.NewApp(&config)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
