package novel

import (
	"runtime"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Element struct {
	scene *sceneWin
	state State

	mouseInputs  []inputMouse
	Hover        func()
	hovered      bool
	ratio        bool
	options      map[Options]float64
	aling        alings
	sx, sy       float64
	px, py       int
	ImageOptions *ebiten.DrawImageOptions

	Resour     Resource
	text       Texting
	animations []animation
}

func (elm *Element) AddInputMouse(input inputMouseInter, button ebiten.MouseButton, action func() error) {
	if input == nil || action == nil {
		return
	}
	elm.mouseInputs = append(elm.mouseInputs, inputMouse{input, button, action})
}

func (elm *Element) GetPosition() (int, int) { return elm.px, elm.py }
func (elm *Element) GetPositionCodinate(x float64, y float64) (float64, float64) {
	xWidth, yWidth := float64((*elm.Resour.Image()).Bounds().Max.X), float64((*elm.Resour.Image()).Bounds().Max.Y)
	return float64(elm.px) + (elm.sx * xWidth * x), float64(elm.py) + (elm.sy * yWidth * y)
}

func (elm *Element) SetPosition(x, y float64) {
	elm.options[XPosition] = x
	elm.options[YPosition] = y
	elm.Recalc()
}

// Func In() Scaling the position and measures of img for calculate if the curso hover it,
// this is why it is not recommended to use very small base images
func (elm *Element) In(x, y int, blnElmsHover bool) bool {
	_, _, _, a := (*elm.Resour.Image()).At(int(float64(x-elm.px)/elm.sx), int(float64(y-elm.py)/elm.sy)).RGBA()

	blnIn := a > 0
	if elm.Hover != nil {
		if blnIn && !elm.hovered && blnElmsHover {
			elm.hovered = true
			elm.Hover()
		}
		if !blnIn && elm.hovered {
			elm.hovered = false
			elm.Resour.Renew()
		}
	}
	return blnIn
}

func (elm *Element) Recalc() {
	elm.sx, elm.sy = 1, 1
	elm.px, elm.py = 0, 0
	xscreen, yscreen := DimesionWin()
	xWidth, yWidth := float64((*elm.Resour.Image()).Bounds().Max.X), float64((*elm.Resour.Image()).Bounds().Max.Y)

	keys := []int{}
	for k := range elm.options {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, intName := range keys {
		name := Options(intName)
		switch name {
		case XPosition:
			elm.px += int(elm.options[name])
		case YPosition:
			elm.py += int(elm.options[name])
		case XPositionScreen:
			elm.px += int(xscreen * elm.options[name])
		case YPositionScreen:
			elm.py += int(yscreen * elm.options[name])
		case XScale:
			elm.sx = elm.options[name] * xWidth
		case YScale:
			elm.sy = elm.options[name] * yWidth
		case XScaleScreen:
			elm.sx = (xscreen * elm.options[name]) / xWidth
		case YScaleScreen:
			elm.sy = (yscreen * elm.options[name]) / yWidth
		case XMaxImage:
			calc(&elm.sx, xWidth*elm.sx > elm.options[name], xWidth, elm.options[name])
		case YMaxImage:
			calc(&elm.sy, yWidth*elm.sy > elm.options[name], yWidth, elm.options[name])
		case XMinImage:
			calc(&elm.sx, xWidth*elm.sx < elm.options[name], xWidth, elm.options[name])
		case YMinImage:
			calc(&elm.sy, yWidth*elm.sy < elm.options[name], yWidth, elm.options[name])
		}
	}
	keys = nil

	if elm.ratio {
		if elm.sx > elm.sy {
			elm.sy = elm.sx
			goto FINISH
		}
		elm.sx = elm.sy
	}
	goto FINISH

FINISH:
	CENTERX := func() { elm.px -= int((xWidth * elm.sx) / 2) }
	CENTERY := func() { elm.py -= int((yWidth * elm.sy) / 2) }
	BOTTON := func() { elm.py -= int((yWidth * elm.sy)) }
	RIGHT := func() { elm.px -= int((xWidth * elm.sx)) }
	switch elm.aling {
	case CenterX:
		CENTERX()
	case CenterY:
		CENTERY()
	case Botton:
		BOTTON()
	case Right:
		RIGHT()
	case Deep:
		RIGHT()
		BOTTON()
	case Center:
		CENTERX()
		CENTERY()
	case BottonCenterX:
		CENTERX()
		BOTTON()
	case RightCenterY:
		CENTERY()
		RIGHT()
	}

	elm.ImageOptions.GeoM.Reset()
	elm.ImageOptions.GeoM.Scale(elm.sx, elm.sy)
	elm.ImageOptions.GeoM.Translate(float64(elm.px), float64(elm.py))
	if elm.text != nil {
		go elm.text.recalc(elm.sx*xWidth, elm.sy*yWidth)
	}
	runtime.GC()
}

func (elm *Element) AddText(text Texting) {
	elm.text = text
	if elm.text.load() {
		xWidth, yWidth := float64((*elm.Resour.Image()).Bounds().Max.X), float64((*elm.Resour.Image()).Bounds().Max.Y)
		elm.text.recalc(elm.sx*xWidth, elm.sy*yWidth)
	} else {
		elm.text = nil
	}
}

func (e *Element) Destroy() {
	e.Resour = nil
	e.Resour = nil
	e.mouseInputs = []inputMouse{}
	e.options = map[Options]float64{}
	e.animations = []animation{}
	*e = Element{}
}

// func NewElm(state *bool, ratio bool, aling alings, resour Resource, options map[Options]float64) *Element {
type OptionElm struct {
	Ratio   bool
	Aling   alings
	Options map[Options]float64
}

func NewElm(state *bool, resour Resource, op OptionElm) *Element {
	elm := Element{}
	if currentScene == nil || state == nil {
		panic("Error novel: should not create elements out a scene, usage InitElmts()")
	}
	elm.scene = currentScene
	elm.state, elm.ratio = state, op.Ratio
	elm.aling = op.Aling
	elm.Resour = resour
	elm.options = op.Options
	if !elm.Resour.Load() {
		panic("Error novel: component can't get the resources")
	}
	elm.ImageOptions = &ebiten.DrawImageOptions{}
	elm.Recalc()
	win.elements = append(win.elements, &elm)
	return &elm
}

func calc(size *float64, bln bool, width, value float64) {
	if bln {
		*size = value / width
	}
}
