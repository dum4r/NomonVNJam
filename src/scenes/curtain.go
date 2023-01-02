package scenes

import (
	"image/color"
	"nomon/src/novel"
	"nomon/src/scenes/lang"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	cursorX     int
	cursorY     int
	everTrue    bool = true
	Curtain     curtain
	FullScreen  *bool
	Resolutions = []string{
		"640x480",
		"800x600",
		"1024x720",
		"1024x768",
		"1280x1024",
		"1920x1080",
	}
)

type curtain struct {
	stepIntro bool

	// Quit View
	btnEquis *novel.Element
	StepQuit bool
	lblQuit  *novel.Element
	btnQuit  *novel.Element
	bgQuit   *novel.Element

	cursor    *novel.Element
	EnterFunc func() error
}

func (m *curtain) Destroy() { *m = curtain{} }
func (m *curtain) InitElmts() {
	m.cursor = novel.NewElm(&everTrue, &novel.GroupImgFor{Url: "cursor", Current: "n"},
		novel.OptionElm{
			Ratio:   true,
			Options: map[novel.Options]float64{novel.XScale: .25},
		})

	// lbl Quit
	m.lblQuit = newLabel(&m.StepQuit, &lang.Menu.QuestExit, color.White,
		novel.OptionElm{
			Aling:   novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .1, novel.XScaleScreen: .8, novel.YScaleScreen: .1},
		})

	m.btnQuit = novel.NewElm(&m.StepQuit, &novel.SolidImage{Url: "btnQuit.png"},
		novel.OptionElm{
			Ratio:   true,
			Aling:   novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .5, novel.XScaleScreen: .3, novel.YScaleScreen: .3},
		})
	m.btnQuit.AddInputMouse(inpututil.IsMouseButtonJustReleased, ebiten.MouseButtonLeft, m.quit)

	m.bgQuit = novel.NewElm(&m.StepQuit,
		&novel.SolidColor{Color: color.RGBA{0, 0, 0, 200}},
		novel.OptionElm{
			Ratio: false, Options: map[novel.Options]float64{novel.XScaleScreen: 1, novel.YScaleScreen: 1},
		})
	m.bgQuit.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, m.bgQuitFunc)

	m.btnEquis = novel.NewElm(FullScreen,
		&novel.GroupImgFor{Url: "curtain/btnEquis", Current: "n"},
		novel.OptionElm{
			Ratio: true, Aling: novel.Right, Options: map[novel.Options]float64{novel.YScaleScreen: .1, novel.XPositionScreen: 1},
		})
	m.btnEquis.Hover = func() {
		m.btnEquis.Resour.(*novel.GroupImgFor).Current = "h"
		m.btnEquis.AniEnding()
	}
	m.btnEquis.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, m.Quit)

	m.stepIntro = true
	novel.NewElm(&m.stepIntro,
		&novel.SolidImage{Url: "curtain/logobg.png"},
		novel.OptionElm{
			Ratio: true,
			Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.YScaleScreen:    .5,
				novel.XPositionScreen: .5,
				novel.YPositionScreen: .5,
			},
		}).ShowFunc(30, func(this *novel.Element) {
		this.Hide(3)
		novel.SetScene("Menu")
	})
}

func (m *curtain) Update() error {
	m.cursor.Resour.Renew()
	if !novel.BlnInput {
		m.cursor.Resour.(*novel.GroupImgFor).Current = "h"
	}
	if x, y := ebiten.CursorPosition(); cursorX != x || cursorY != y {
		cursorX = x
		cursorY = y
		m.cursor.SetPosition(float64(cursorX), float64(cursorY))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return m.Quit()
	}
	if m.StepQuit {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			m.cursor.SetPosition(m.bgQuit.GetPositionCodinate(.5, 0.2))
			m.EnterFunc = m.bgQuitFunc
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			m.cursor.SetPosition(m.btnQuit.GetPositionCodinate(.5, .5))
			m.EnterFunc = m.quit
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if m.EnterFunc != nil {
				return m.EnterFunc()
			}
		}
	}
	return nil
}

func (m *curtain) bgQuitFunc() error {
	novel.HideState(&m.StepQuit, 2)
	return nil
}
func (m *curtain) Quit() error {
	if !m.bgQuit.IsActive() {
		m.lblQuit.Show(3)
		m.btnQuit.Aclare(3)
		m.bgQuit.Aclare(5)
	}
	return nil
}
func (m *curtain) quit() error { return novel.Terminator() }

// Config Scene
type Config struct {
	lblConfig     *novel.Element
	btnBack       *novel.Element
	chckScreen    *novel.Element
	lblResolution *novel.Element
	slcScreen     *novel.Element
	lblLang       *novel.Element
	slcLang       *novel.Element
}

func (m *Config) Destroy() { *m = Config{} }
func (m *Config) InitElmts() {
	stepConfig := true
	newLabel(&stepConfig, &lang.Menu.Config, color.White,
		novel.OptionElm{
			Aling:   novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .1, novel.XScaleScreen: .8, novel.YScaleScreen: .08},
		})

	// FullScreen
	newLabel(&stepConfig, &lang.Menu.ChooseFullScreen, color.White,
		novel.OptionElm{
			Options: map[novel.Options]float64{novel.XPositionScreen: .1, novel.YPositionScreen: .15, novel.XScaleScreen: .4, novel.YScaleScreen: .045},
		})
	newCheck(&stepConfig, FullScreen, novel.UpdateWin,
		novel.OptionElm{Ratio: true, Aling: novel.CenterX,
			Options: map[novel.Options]float64{novel.XPositionScreen: .7, novel.YPositionScreen: .15, novel.YScaleScreen: .045},
		})

	newLabel(&stepConfig, &lang.Menu.ChooseResolution, color.White,
		novel.OptionElm{
			Options: map[novel.Options]float64{novel.XPositionScreen: .1, novel.YPositionScreen: .2, novel.XScaleScreen: .4, novel.YScaleScreen: .045},
		})

	newSelectFlex(&stepConfig, Resolutions, color.White,
		func(s string) {},
		novel.OptionElm{
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .2, novel.XScaleScreen: .4, novel.YScaleScreen: .045},
		},
	)
	// m.lblResolution = novel.NewElm(&stepConfig, false, novel.Center, baseNoCollition,
	// 	map[novel.Options]float64{novel.XPositionScreen: .25, novel.YPositionScreen: .3, novel.XScaleScreen: .5, novel.YScaleScreen: .1})
	// counter := 0.0
	// for str, point := range Resolutions {
	// 	name := str
	// 	width, height := point.X, point.Y
	// 	yScaleScreen := .45 / float64(len(Resolutions))
	// 	tempElm := novel.NewElm(&stepConfig, false, novel.Center, &novel.SolidColor{Color: color.RGBA{0, 0, 0, 255}},
	// 		map[novel.Options]float64{
	// 			novel.XPositionScreen: .25,
	// 			novel.XScaleScreen:    .3,
	// 			novel.YScaleScreen:    yScaleScreen,
	// 			novel.YPositionScreen: .5 + (yScaleScreen * counter),
	// 		})
	// 	tempElm.AddText(&novel.TextFlex{Text: &name, Padding: 0.2, ColorFont: color.White})
	// 	tempElm.Hover = func() { tempElm.Resour.(*novel.SolidColor).Color = color.RGBA{100, 100, 100, 70} }
	// 	tempElm.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
	// 		ebiten.SetWindowSize(width, height)
	// 		return nil
	// 	})
	// 	counter++
	// }
	// //btn Save
	// m.btnBack = novel.NewElm(&stepConfig, false, novel.Center,
	// 	&novel.SolidColor{Color: color.RGBA{100, 100, 100, 70}},
	// 	map[novel.Options]float64{
	// 		novel.XPositionScreen: .75,
	// 		novel.XScaleScreen:    .3,
	// 		novel.YScaleScreen:    .2,
	// 		novel.YPositionScreen: .15,
	// 	})
	// m.btnBack.AddText(&novel.TextFlex{Text: &lang.Menu.BtnConfig, Padding: 0.15, ColorFont: color.Black})
	// m.btnBack.Hover = func() { m.btnBack.Resour.(*novel.SolidColor).Color = color.RGBA{255, 255, 255, 70} }
	// m.btnBack.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
	// 	novel.HideState(&stepConfig, 2)
	// 	ControllerScenes("")
	// 	return nil
	// })

	// m.lblResolution.AddText(&novel.TextFlex{Text: &lang.Menu.ChooseResolution, ColorFont: color.White, Padding: .2})
	// m.lblLang = novel.NewElm(&stepConfig, false, novel.Center, baseNoCollition,
	// 	map[novel.Options]float64{novel.XPositionScreen: .75, novel.YPositionScreen: .3, novel.XScaleScreen: .5, novel.YScaleScreen: .1})
	// m.lblLang.AddText(&novel.TextFlex{Text: &lang.Menu.ChooseLang, ColorFont: color.White, Padding: .2})
	// totalLangs := float64(len(util.Languages()))
	// for i, leng := range util.Languages() {
	// 	nameLang := leng
	// 	yScaleScreen := .45 / totalLangs
	// 	tempElm := novel.NewElm(&stepConfig, false, novel.Center, &novel.SolidColor{Color: color.RGBA{0, 0, 0, 255}},
	// 		map[novel.Options]float64{
	// 			novel.XPositionScreen: .75,
	// 			novel.XScaleScreen:    .3,
	// 			novel.YScaleScreen:    yScaleScreen,
	// 			novel.YPositionScreen: .5 + (yScaleScreen * float64(i)),
	// 		})
	// 	tempElm.AddText(&novel.TextFlex{Text: &nameLang, Padding: 0.2, ColorFont: color.White})
	// 	tempElm.Hover = func() { tempElm.Resour.(*novel.SolidColor).Color = color.RGBA{100, 100, 100, 70} }
	// 	tempElm.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
	// 		go novel.PushLanguage(nameLang)
	// 		return nil
	// 	})
	// }
}
func (*Config) Update() error { return nil }
