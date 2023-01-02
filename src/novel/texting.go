package novel

import (
	"image"
	"image/color"
	"nomon/src/novel/util"
)

type Texting interface {
	load() bool
	recalc(w, h float64)
	image() *image.Image
}

type TextFlex struct {
	Text      *string
	Padding   float64
	ColorFont color.Color
	img       image.Image
}

func (s *TextFlex) load() bool {
	if s.Text == nil {
		panic("Error novel: pointer of Text is nil.")
	}
	return true
}

func (s *TextFlex) image() *image.Image {
	return &s.img
}

func (s *TextFlex) recalc(w, h float64) {
	s.img = util.StringFlex(s.Text, w, h, s.Padding, s.ColorFont)
}
