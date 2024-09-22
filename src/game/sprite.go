package game

import (
	"image"
	"log"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	image *ebiten.Image
}

type SpriteSheet struct {
	sprites []Sprite
}

type SpriteBuilder struct {
	image *ebiten.Image
	tiles []Tile
}

type Tile struct {
	// the horizontal offset of the tile.
	X int
	// the vertical offset of the tile.
	Y int
	// the width of the tile.
	Width int
	// the width of the tile.
	Height int
}

func (s SpriteBuilder) FromFile(file string) SpriteBuilder {
	image, _, err := ebitenutil.NewImageFromFile(file)
	if err != nil {
		log.Fatal(err)
	}

	s.image = image
	return s
}

type CreateTilesInput struct {
	x_start int
	y_start int
	width   int
	height  int
	x_count int
	y_count int
}

func (s SpriteBuilder) CreateTiles(input CreateTilesInput) SpriteBuilder {
	x_last := input.x_start + input.x_count*input.width
	y_last := input.y_start + input.y_count*input.height

	for x := input.x_start; x < x_last; x += input.width {
		for y := input.y_start; y < y_last; y += input.height {
			tile := Tile{X: x, Y: y, Width: input.width, Height: input.height}
			s.tiles = append(s.tiles, tile)
		}
	}
	return s
}

func (s SpriteBuilder) FilterTiles(tiles_to_filter []Tile) SpriteBuilder {
	s.tiles = slices.DeleteFunc(s.tiles, func(offset Tile) bool {
		return slices.Contains(tiles_to_filter, offset)
	})
	return s
}

func (s SpriteBuilder) BuildAsSpriteSheet() SpriteSheet {
	sprites := make([]Sprite, len(s.tiles))

	for i, tile := range s.tiles {
		rect := image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height)
		sprites[i] = Sprite{image: s.image.SubImage(rect).(*ebiten.Image)}
	}

	return SpriteSheet{sprites: sprites}
}

func (s SpriteBuilder) BuildAsBackgroundSprite(background_width, background_height, tile_width, tile_height int) Sprite {
	new_image := ebiten.NewImage(int(background_width), int(background_height))
	subimages := make([]*ebiten.Image, len(s.tiles))

	for i, tile := range s.tiles {
		rect := image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height)
		subimages[i] = s.image.SubImage(rect).(*ebiten.Image)
	}

	for x := 0; x < background_width; x += tile_width {
		for y := 0; y < background_height; y += tile_height {
			drawOptions := ebiten.DrawImageOptions{}
			drawOptions.GeoM.Translate(float64(x), float64(y))

			random_subimage := subimages[rand.Intn(len(s.tiles))]
			new_image.DrawImage(random_subimage, &drawOptions)
		}
	}

	return Sprite{image: new_image}
}

func NewSpriteBuilder() SpriteBuilder {
	return SpriteBuilder{
		image: nil,
	}
}
