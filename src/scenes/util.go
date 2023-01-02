package scenes

import (
	"image/color"
	"nomon/src/novel"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var baseNoCollition = &novel.SolidColor{Color: color.RGBA{0, 0, 0, 0}}

var (
	textIzq string = " << "
	textDer string = " >> "
)

// this select not consider options elements max min scale, olny XScaleScreen
func newSelectFlex(state *bool, array []string, color color.Color, f func(string), op novel.OptionElm) *novel.Element {
	if scale, ok := op.Options[novel.XScaleScreen]; ok {

		index := 0
		text := array[index]
		xscreen, yscreen := novel.DimesionWin()

		lbl := newLabel(state, &text, color, op)

		checkIndex := func() {
			if index < 0 {
				index = 0
				return
			}
			if top := len(array) - 1; index > top {
				index = top
				return
			}
			text = array[index]
			spl := strings.Split(text, "x")
			w, _ := strconv.Atoi(spl[0])
			h, _ := strconv.Atoi(spl[1])
			ebiten.SetWindowSize(w, h)
		}

		op1 := op
		op1.Options = make(map[novel.Options]float64)
		// Copy from the original map to the target map
		for key, value := range op.Options {
			op1.Options[key] = value
		}
		h := (yscreen * op.Options[novel.YScaleScreen]) * 2
		op1.Options[novel.XScaleScreen] = h / xscreen

		newSimpleBtn(state, &textIzq, color, op1,
			func() error {
				defer checkIndex()
				index--
				return nil
			})

		op2 := op1
		op2.Options = make(map[novel.Options]float64)
		// Copy from the original map to the target map
		for key, value := range op1.Options {
			op2.Options[key] = value
		}

		op2.Options[novel.XPositionScreen] += ((xscreen * scale) - h) / xscreen
		newSimpleBtn(state, &textDer, color, op2,
			func() error {
				defer checkIndex()
				index++
				return nil
			})
		return lbl
	}
	panic("Error New Select not option XScaleScreen")
}

func newSimpleBtn(state *bool, text *string, colorFont color.Color, op novel.OptionElm, clic func() error) *novel.Element {
	btn := novel.NewElm(state,
		&novel.SimpleString{Text: text, Size: 50, ColorFont: colorFont, ColorBG: color.RGBA{100, 100, 100, 255}},
		op,
	)
	btn.AniInfinite()
	btn.Hover = func() {
		btn.Resour.(*novel.SimpleString).ColorBG = color.RGBA{150, 150, 150, 255}
	}
	btn.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, clic)
	return btn
}

func newLabel(state *bool, text *string, color color.Color, op novel.OptionElm) *novel.Element {
	lbl := novel.NewElm(state, baseNoCollition, op)
	lbl.AddText(&novel.TextFlex{Text: text, ColorFont: color})
	return lbl
}
func newCheck(state *bool, value *bool, f func(), op novel.OptionElm) {
	current := "0"
	if *value {
		current = "1"
	}
	check := novel.NewElm(state, &novel.GroupImgFor{Url: "curtain/check", Current: current}, op)
	check.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
		check.Resour.(*novel.GroupImgFor).SetCurrent("a")
		*value = !*value
		if f != nil {
			f()
		}
		if *value {
			check.AniEndingFunc(func() { check.Resour.(*novel.GroupImgFor).SetCurrent("1") })
			return nil
		}
		check.AniEndingFuncRever(func() { check.Resour.(*novel.GroupImgFor).SetCurrent("0") })
		return nil
	})
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
