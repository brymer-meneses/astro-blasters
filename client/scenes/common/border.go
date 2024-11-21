package common

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Border struct {
	Image *ebiten.Image
}

func NewBorder(width int, height int) *Border {
	border := ebiten.NewImage(int(width), int(height))

	// Create a new Border instance
	b := &Border{
		Image: border,
	}

	// Tile dimensions
	// tileWidth := 16
	// tileHeight := 16

	// Define tile regions for borders within ui.png
	// rects := []assets.TileIndex{
	// {X: 1, Y: 0}, // blue tile
	// {X: 0, Y: 0}, // purple tile
	// {X: 1, Y: 3}, // grey tile

	// {X: 0, Y: 1}, // white border tile
	// {X: 0, Y: 2}, // black border tile (unnoticeable)
	// {X: 0, Y: 3}, // gold tile
	// {X: 0, Y: 4}, // gold border tile
	// {X: 1, Y: 4}, // grey border tile
	// }

	return b
}
