package assets

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Background Sprite
var Ships Sprite

var FontNarrow *text.GoTextFaceSource

func init() {
	Background = mustLoadImageFromBytes(backgrounds, 128, 256)
	Ships = mustLoadImageFromBytes(ships, 8, 8)

	FontNarrow = mustLoadFontFromBytes(munroNarrow)
}

//go:embed MunroFont/munro-narrow.ttf
var munroNarrow []byte

//go:embed SpaceShooterAssetPack/BackGrounds.png
var backgrounds []byte

//go:embed SpaceShooterAssetPack/Ships.png
var ships []byte

func mustLoadImageFromBytes(data []byte, width, height int) Sprite {
	image, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return Sprite{image, width, height}
}

func mustLoadFontFromBytes(data []byte) *text.GoTextFaceSource {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return fontSource
}
