package main

import (
	"log"
	"space-shooter/app"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	width, height := 1280, 720
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	app := app.NewApp(width, height)

	if err := ebiten.RunGame(&app); err != nil {
		log.Fatal(err)
	}

}
