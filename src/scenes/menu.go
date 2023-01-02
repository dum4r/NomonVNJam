package scenes

import (
	"image/color"
	"nomon/src/novel"
	"nomon/src/scenes/lang"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Menu struct {
	stepTwo bool
	btnPlay *novel.Element

	stepThird bool
	btnConfig *novel.Element
	btnExtr   *novel.Element
	btnCrdts  *novel.Element

	list []any
}

func (m *Menu) InitElmts() {
	stepFirst := true
	// Tittle Label
	novel.NewElm(&stepFirst, &novel.Oneframe{Url: "menu/tittle"},
		novel.OptionElm{
			Ratio: true, Aling: novel.Center,
			Options: map[novel.Options]float64{novel.YPositionScreen: .2, novel.YScaleScreen: .2, novel.XPositionScreen: .5},
		}).AniEndingFunc(func() { m.stepTwo = true })

	// Label version
	novel.NewElm(&m.stepTwo,
		&novel.SimpleString{Text: &lang.Version, Size: 10, ColorFont: color.Black, ColorBG: color.RGBA{0, 0, 0, 30}},
		novel.OptionElm{Options: map[novel.Options]float64{novel.YPositionScreen: .38, novel.XPositionScreen: .6}}).Show(30)

	// First button Creditos
	m.btnCrdts = novel.NewElm(&m.stepThird,
		&novel.SolidColor{Color: color.RGBA{100, 100, 100, 70}},
		novel.OptionElm{Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.XPositionScreen: .5,
				novel.YPositionScreen: .6,
				novel.XScaleScreen:    .3,
				novel.YScaleScreen:    .09,
				novel.YMaxImage:       30,
				novel.YPosition:       -60,
			},
		})
	m.btnCrdts.AddText(&novel.TextFlex{Text: &lang.Menu.Credits, Padding: 0.15, ColorFont: color.Black})
	m.btnCrdts.Hover = func() { m.btnCrdts.Resour.(*novel.SolidColor).Color = color.RGBA{255, 255, 255, 70} }
	m.btnCrdts.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
		return nil
	})

	// Second button Config
	m.btnConfig = novel.NewElm(&m.stepThird,
		&novel.SolidColor{Color: color.RGBA{100, 100, 100, 70}},
		novel.OptionElm{
			Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.XPositionScreen: .5,
				novel.YPositionScreen: .6,
				novel.XScaleScreen:    .3,
				novel.YScaleScreen:    .09,
				novel.YMaxImage:       30,
				novel.YPosition:       -30,
			},
		})
	m.btnConfig.AddText(&novel.TextFlex{Text: &lang.Menu.Config, Padding: 0.15, ColorFont: color.Black})
	m.btnConfig.Hover = func() { m.btnConfig.Resour.(*novel.SolidColor).Color = color.RGBA{255, 255, 255, 255} }
	m.btnConfig.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
		return novel.SetScene("Config")
	})

	// Third button Extra
	m.btnExtr = novel.NewElm(&m.stepThird,
		&novel.SolidColor{Color: color.RGBA{100, 100, 100, 70}},
		novel.OptionElm{
			Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.XPositionScreen: .5,
				novel.YPositionScreen: .6,
				novel.XScaleScreen:    .3,
				novel.YScaleScreen:    .09,
				novel.YMaxImage:       30,
			},
		})
	m.btnExtr.AddText(&novel.TextFlex{Text: &lang.Menu.Extra, Padding: 0.15, ColorFont: color.Black})
	m.btnExtr.Hover = func() { m.btnExtr.Resour.(*novel.SolidColor).Color = color.RGBA{255, 255, 255, 70} }
	m.btnExtr.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
		println("Click 3")
		return nil
	})

	// Last Button Play
	m.btnPlay = novel.NewElm(&m.stepTwo,
		&novel.GroupImgFor{Url: "menu/btnPlay", Current: "i"},
		novel.OptionElm{
			Ratio: true, Aling: novel.BottonCenterX,
			Options: map[novel.Options]float64{
				novel.XPositionScreen: .5,
				novel.YPositionScreen: 1,
				novel.XScaleScreen:    .35,
			},
		})
	m.btnPlay.AniEndingFunc(func() {
		m.btnPlay.Resour.(*novel.GroupImgFor).SetCurrent("n")
		m.btnPlay.Recalc()
		m.btnPlay.AddText(&novel.TextFlex{Text: &lang.Menu.Star, Padding: 0.3, ColorFont: color.White})
		novel.ShowState(&m.stepThird, 6)

		m.btnPlay.AniInfinite()
		m.btnPlay.Hover = func() { m.btnPlay.Resour.(*novel.GroupImgFor).Current = "h" }
		m.btnPlay.AddInputMouse(inpututil.IsMouseButtonJustPressed, ebiten.MouseButtonLeft, func() error {
			return novel.SetScene("wellcome")
		})
	})

	// BackGround
	novel.NewElm(&stepFirst,
		&novel.Oneframe{Url: "menu/background"},
		novel.OptionElm{
			Ratio: true, Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.YScaleScreen:    1,
				novel.XScaleScreen:    1,
				novel.YPositionScreen: 0.5,
				novel.XPositionScreen: 0.5,
			},
		}).AniInfinite()

	// list := []struct {
	// 	e *novel.Element
	// 	f func() error
	// }{}

}
func (m *Menu) Update() error {
	if !Curtain.StepQuit {
		// if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		// 	Curtain.cursor.SetPosition(m.bgQuit.GetPositionCodinate(.5, 0.2))
		// 	m.EnterFunc = m.bgQuitFunc
		// 	return nil
		// }
		// if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// 	if m.EnterFunc != nil {
		// 		return m.EnterFunc()
		// 	}
		// }
	}
	return nil
}
func (m *Menu) Destroy() { *m = Menu{} }
