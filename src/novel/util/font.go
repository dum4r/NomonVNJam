package util

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func init() {
	libraryFonts = make(map[rune]*truetype.Font)
	libraryFonts['$'], _ = truetype.Parse(fonts.MPlus1pRegular_ttf)
}

type TemplFont interface {
	Print() (string, rune)
}

var libraryFonts map[rune]*truetype.Font

func GetFont(decorator rune) *truetype.Font {
	if font, found := libraryFonts[decorator]; found {
		return font
	}
	return nil
}

func getFace(char rune, size float64) font.Face {
	return truetype.NewFace(libraryFonts[char], &truetype.Options{
		Size:    size,
		Hinting: font.HintingNone,
	})
}

// func DrawStringWrapped(s string, align Align, padding float64, img image.Image) {
// 	width, heigth := img.Bounds().Max.X, img.Bounds().Max.Y
// 	lines := WordWrap(s, width)

// 	// sync h formula with MeasureMultilineString
// 	h := float64(len(lines)) * dc.fontHeight * lineSpacing
// 	h -= (lineSpacing - 1) * dc.fontHeight

//		x -= ax * width
//		y -= ay * h
//		switch align {
//		case AlignLeft:
//			ax = 0
//		case AlignCenter:
//			ax = 0.5
//			x += width / 2
//		case AlignRight:
//			ax = 1
//			x += width
//		}
//		ay = 1
//		for _, line := range lines {
//			dc.DrawStringAnchored(line, x, y, ax, ay)
//			y += dc.fontHeight * lineSpacing
//		}
//	}
//
//	func PrintString(img image.Image, str string, color color.Color) image.Image {
//		size := float64(img.Bounds().Max.Y) * 0.9
//		w, h := MeasureString(removeDec(str), size, getFace('$', size))
//		x, y := .0, h*.85
//		drawString(str, size, x, y, color, img)
//		return img
//	}

// func minSize(str string, size, width, height float64) float64 {
// 	w, h := MeasureString(removeDec(str), size, getFace('$', size))
// 	if w < width && h < height {
// 		return size
// 	}
// 	return minSize(str, size-1., width, height)
// }

func minSize(str *string, width, height float64) (float64, float64, float64) {
	w, h := width, height
	size := height - 1
	for w >= width || h >= height {
		w, h = MeasureString(removeDec(*str), size, getFace('$', size))
		runtime.GC()
		size = size - 1
	}
	return size, w, h
}

func StringFlex(str *string, width, height float64, padding float64, clrFont color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	padding = padding - .88
	size, w, h := minSize(str, -width*padding, -height*padding)
	x, y := (width-w)/2, .42*(height+h)
	drawString(*str, size, x, y, clrFont, img)
	return img
}

func StringtoImageRect(str string, size float64, clrFont color.Color, clrBG color.RGBA) image.Image {
	w, h := MeasureString(removeDec(str), size, getFace('$', size))
	x, y := .0, h*.85
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h*1.1)))

	for x := 0; x < int(w); x++ {
		for y := 0; y < int(h*1.1); y++ {
			img.SetRGBA(x, y, clrBG)
		}
	}

	drawString(str, size, x, y, clrFont, img)
	return img
}

func drawString(str string, size, x, y float64, color color.Color, img *image.RGBA) {
	dot := fixp(x, y)
	previusChar := rune(-1)
	face := getFace('$', size)
	for _, char := range str {
		if previusChar >= 0 {
			dot.X += face.Kern(previusChar, char)
		}
		if _, a := libraryFonts[char]; a {
			face = getFace(char, size)
		} else {
			dr, mask, maskp, advance, ok := face.Glyph(dot, char)
			if !ok {
				continue
			}
			draw.DrawMask(img, dr, image.NewUniform(color), image.Point{}, mask, maskp, draw.Over)
			dot.X += advance
		}
		previusChar = char
	}
}

func MeasureString(s string, fontHeight float64, ff font.Face) (w, h float64) {
	a := font.MeasureString(ff, s)
	return float64(a >> 6), fontHeight
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{X: fix(x), Y: fix(y)}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

func removeDec(str string) string {
	for a, _ := range libraryFonts {
		str = strings.Replace(str, string(a), "", -1)
	}
	return str
}

func ImageToBytesSlow(img image.Image) []byte {
	size := img.Bounds().Size()
	w, h := size.X, size.Y
	bs := make([]byte, 4*w*h)

	dstImg := &image.RGBA{
		Pix:    bs,
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	draw.Draw(dstImg, image.Rect(0, 0, w, h), img, img.Bounds().Min, draw.Src)
	return bs
}
