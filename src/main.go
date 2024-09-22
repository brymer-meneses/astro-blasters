package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"space-shooter/game"
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
