package main

import (
	"log"
	"space-shooter/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	width, height := 1280, 720
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Space Shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := game.NewGame(width, height)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}

}
