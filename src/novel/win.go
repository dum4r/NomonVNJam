package novel

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"nomon/src/novel/util"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	chao       = "Chaou"
	stepSecond = 60
	stepWorld  = time.Second / stepSecond
	fileConfig = "config.json"
)

var (
	Quit                            func() error
	win                             window
	BlnInput, BlnTick, BlnElmsHover = true, false, true
	Termination                     = errors.New(chao)
)

// Window is admin, scene iterator and interface to Ebitengine
type window struct {
	blnWin bool
	// context     *fonts.Context
	windowTitle  string
	Lang         string
	Position     image.Point
	Dimesion     image.Point
	Fullscreen   bool
	TPS          int
	Ticks        time.Duration
	TimeAutoSave time.Duration
	Obj          any

	windowResizable bool

	icon []image.Image

	currentScene string
	curtain      string
	scenes       map[string]*sceneWin
	elements     []*Element

	t1 time.Time
	t2 time.Time
}

// Start window, settings and global var AssetsFiles
func Start(assetsFS *embed.FS, icon *[]byte, title string, Obj any, tables map[string]any, arrs ...util.TemplFont) (*window, error) {
	win.scenes = make(map[string]*sceneWin)
	win.windowTitle = title
	win.windowResizable = false
	win.Obj = Obj

	if err := win.loadConfigs(); err != nil {
		return nil, err
	}

	if err := util.SetAssetsFS(assetsFS, &win.Lang, tables, arrs...); err != nil {
		return nil, err
	}
	arrs = nil

	win.icon = []image.Image{util.ImageForBytes(*icon)}
	UpdateWin()

	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum) // force Maximun FPS
	ebiten.SetCursorMode(ebiten.CursorModeHidden)    // Not cursor
	ebiten.SetScreenClearedEveryFrame(true)
	// ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetInitFocused(true)

	return &win, nil
}

// load Config if file-config not exist, create the file with data predefined
func (w *window) loadConfigs() error {
	if _, err := os.Stat(fileConfig); os.IsNotExist(err) {
		if _, err := os.Create(fileConfig); err != nil {
			return err
		}
		w.Dimesion = image.Point{1366, 768}
		w.Lang = "Español"
		w.TPS = 60
		w.Ticks = 15
		w.TimeAutoSave = 1

		file, _ := json.MarshalIndent(w, "", " ")
		if err := os.WriteFile(fileConfig, file, 0644); err != nil {
			return err
		}
		return nil
	} else {
		dat, err := os.ReadFile(fileConfig)
		if err != nil {
			return err
		}
		if len(dat) < 1 {
			panic("Error novel: no data in file loading config")
		}
		if err := json.Unmarshal([]byte(dat), &win); err != nil {
			print("Error novel: corrupted file")
			panic(err)
		}
	}
	return nil
}

// Push settings to fileConfig
func (w *window) saveSettings() error {
	file, _ := json.MarshalIndent(w, "", " ")
	if err := os.WriteFile(fileConfig, file, 0644); err == nil {
		return err
	}
	return nil
}

func (w *window) changeScene(newScene string) {
	if w.scenes[w.currentScene] != nil {
		// Delete elmts of previous scene
		for index, end := 0, len(w.elements); index < end; {
			elm := w.elements[index]
			if elm.scene == w.scenes[w.currentScene] {
				w.elements = append(w.elements[:index], w.elements[index+1:]...) // delete Element from array window
				elm.Destroy()
				end--
			} else {
				index++
			}
		}
		go (*w.scenes[w.currentScene]).Destroy()
	}
	runtime.GC()
	(*w.scenes[newScene]).init()
	w.currentScene = newScene
}

func (w *window) Update() error {
	win.blnWin = true
	if ebiten.IsWindowBeingClosed() {
		if Quit == nil {
			return Terminator()
		}
		return Quit()
	}
	if w.scenes[w.currentScene] != nil {
		if (*w.scenes[w.currentScene]).IsDone() {
			if err := (*w.scenes[w.currentScene]).Update(); err != nil {
				return err
			}
		}
	}

	BlnInput, BlnTick, BlnElmsHover = true, false, true
	x, y := ebiten.CursorPosition()
	if time.Now().After(w.t1) {
		w.t1 = time.Now().Add(time.Second / w.Ticks)
		BlnTick = true
	}

	for _, o := range w.elements {
		if !w.blnWin {
			return nil // Return Update() if change scene while execute loop elemets
		}
		if o.In(x, y, BlnElmsHover) && *o.state {
			if BlnInput && len(o.mouseInputs) > 0 {
				for _, i := range o.mouseInputs {
					if i.input(i.button) {
						if err := i.action(); err != nil {
							return err
						}
					}
				}
				BlnInput = false
			}
			BlnElmsHover = false
		}

		if nmax := len(o.animations); nmax > 0 && BlnTick && *o.state {
			if !o.animations[nmax-1]() {
				o.animations = append(o.animations[:nmax-1], o.animations[nmax:]...)
			}
		}
	}

	if w.curtain != "" {
		if err := (*w.scenes[w.curtain]).Update(); err != nil {
			return err
		}
	}

	if time.Now().After(w.t2) {
		w.t2 = time.Now().Add(time.Minute * w.TimeAutoSave)
		go w.saveSettings()
	}

	return nil
}

// Draw text option for translate
var (
	opimg = &ebiten.NewImageFromImageOptions{
		Unmanaged:      true,
		PreserveBounds: true,
	}
	optex = ebiten.DrawImageOptions{}
)

func (w *window) Draw(screen *ebiten.Image) {
	img := ebiten.NewImageWithOptions(screen.Bounds(), &ebiten.NewImageOptions{})
	for i := len(w.elements) - 1; i >= 0; i-- {
		if elm := w.elements[i]; *elm.state {
			eimg := ebiten.NewImageFromImageWithOptions(*elm.Resour.Image(), opimg)
			img.DrawImage(eimg, elm.ImageOptions)
			if text := elm.text; text != nil {
				eimg = ebiten.NewImageFromImageWithOptions(*text.image(), opimg)
				optex.GeoM.Translate(float64(elm.px), float64(elm.py))
				optex.ColorM = elm.ImageOptions.ColorM
				img.DrawImage(eimg, &optex)
				optex.GeoM.Reset()
				optex.ColorM.Reset()
			}
			eimg.Dispose()
		}
	}
	screen.DrawImage(img, &optex)
	img.Dispose()
	msg := fmt.Sprint("FPS=", ebiten.ActualFPS(), " TPS=", ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (w *window) Layout(width, height int) (int, int) {
	win.Position.X, win.Position.Y = ebiten.WindowPosition() // position
	if w.Dimesion.X != width || w.Dimesion.Y != height {
		w.Dimesion.X, w.Dimesion.Y = width, height
		for _, elm := range w.elements {
			elm.Recalc()
		}
	}
	return w.Dimesion.X, w.Dimesion.Y
}

func (w *window) SetCurtain(c Scene1) {
	AddScene("Curtain", c)
	win.curtain = "Curtain"
	(*w.scenes[w.curtain]).init()
}

func PushLanguage(lang string) {
	if !win.blnWin {
		return
	}
	win.Lang = lang
	util.PushLanguage()
	for _, e := range win.elements {
		if e.text != nil {
			go e.Recalc()
		}
	}
}

func Terminator() error {
	if win.blnWin {
		win.saveSettings()
		return Termination
	}
	return nil
}

func SetScene(nameScene string) error {
	if win.blnWin {
		win.blnWin = false
		win.changeScene(nameScene)
	}
	return nil
}

// Update settings of the var windows to Ebitengine
func UpdateWin() {
	w := &win
	wFull, hFull := ebiten.ScreenSizeInFullscreen()
	ebiten.SetWindowPosition(w.Position.X, w.Position.Y)
	if win.Fullscreen {
		ebiten.SetWindowSize(wFull, hFull)
		ebiten.SetWindowDecorated(false)
	} else {
		if wFull < win.Dimesion.X {
			win.Dimesion.X = wFull
		}
		if hFull < win.Dimesion.Y {
			win.Dimesion.Y = hFull
		}
		ebiten.SetWindowSize(win.Dimesion.X, win.Dimesion.Y)
		ebiten.SetWindowDecorated(true)
	}
	ebiten.SetFullscreen(win.Fullscreen)
	ebiten.SetWindowTitle(win.windowTitle)
	if w.windowResizable {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	} else {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	}
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetWindowIcon(w.icon)
	ebiten.SetTPS(win.TPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
}
func DimesionWin() (float64, float64) { return float64(win.Dimesion.X), float64(win.Dimesion.Y) }

func init() {
	// allow directX if passed as program flag
	for _, arg := range os.Args {
		if arg == "--directX" {
			return
		}
	}

	// set openGL as the graphics backend otherwise
	err := os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	if err != nil {
		panic(err)
	}
}
