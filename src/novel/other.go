package novel

import "github.com/hajimehoshi/ebiten/v2"

type inputMouseInter func(ebiten.MouseButton) bool

type inputMouse struct {
	input  inputMouseInter
	button ebiten.MouseButton
	action func() error
}

type alings uint8

const (
	// options for the position according to the size of the image
	Default alings = iota
	CenterX
	CenterY
	Botton
	Right
	Deep
	Center
	BottonCenterX
	RightCenterY
)

type Options uint8

const (
	XPosition Options = iota
	YPosition
	XPositionScreen
	YPositionScreen
	XScale
	YScale
	XScaleScreen
	YScaleScreen
	XMaxImage
	YMaxImage
	XMinImage
	YMinImage
)

func AddFont(url string, dec rune) *templFont { return &templFont{url: url, dec: dec} }

type templFont struct {
	url string
	dec rune
}

func (t *templFont) Print() (string, rune) { return t.url, t.dec }
