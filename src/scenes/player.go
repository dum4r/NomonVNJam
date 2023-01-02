package scenes

import (
	"image/color"
	"nomon/src/novel"
	"nomon/src/scenes/lang"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type player struct {
	Name string
}

var (
	Player player
	Scenes = map[string]novel.Scene1{
		"wellcome": &WellcomeScene{},
	}
)

func ControllerScenes(currentScene string) {
	println("== Controller Scenes ==")
	if currentScene == "" {
		if Player.Name == "" {
			novel.SetScene("wellcome")
			println("No Name")
			return
		}
		println("Config widt name")
		return
	}
}

// Scene Wellcome
type WellcomeScene struct {
	// Name View
	stepName bool
	Name     *novel.Element
	runes    []rune
	text     string
}

func (m *WellcomeScene) Destroy() { *m = WellcomeScene{} }
func (m *WellcomeScene) InitElmts() {
	blnWell := true
	// Nomon Label
	novel.NewElm(&blnWell, &novel.SolidImage{Url: "menu/tittle/99.png"},
		novel.OptionElm{
			Ratio: true, Aling: novel.Center,
			Options: map[novel.Options]float64{
				novel.XScaleScreen:    .8,
				novel.XPositionScreen: .5,
				novel.YPositionScreen: .35,
			}}).Show(40)

	// Version label
	novel.NewElm(&blnWell,
		&novel.SimpleString{Text: &lang.Version, Size: 20, ColorFont: color.White, ColorBG: color.RGBA{0, 0, 0, 0}},
		novel.OptionElm{
			Ratio: true, Aling: novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .7},
		}).ShowFunc(60,
		func(this *novel.Element) {
			this.HideFunc(16, func() { novel.ShowState(&m.stepName, 8) })
		})

	// Name View
	newLabel(&m.stepName, &lang.Menu.WriteName, color.White,
		novel.OptionElm{
			Aling:   novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .1, novel.XScaleScreen: .9, novel.YScaleScreen: .1},
		})
	m.Name = newLabel(&m.stepName, &m.text, color.White,
		novel.OptionElm{
			Aling:   novel.Center,
			Options: map[novel.Options]float64{novel.XPositionScreen: .5, novel.YPositionScreen: .5, novel.XScaleScreen: .5, novel.YScaleScreen: .2},
		})
}

func (m *WellcomeScene) Update() error {
	if m.stepName {
		// If the backspace key is pressed, remove one character.
		if repeatingKeyPressed(ebiten.KeyBackspace) && len(m.text) >= 1 {
			m.text = m.text[:len(m.text)-1]
			m.Name.Recalc()
			return nil
		}

		m.runes = ebiten.AppendInputChars(m.runes[:0])
		if len(m.text) < 14 && len(m.runes) > 0 {
			m.text += string(m.runes)
			m.Name.Recalc()
			return nil
		}

		if name := strings.TrimSpace(m.text); name != "" &&
			(inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter)) {
			Player.Name = name
			m.stepName = false
			novel.SetScene("Menu")
			return nil
		}
	}
	return nil
}
