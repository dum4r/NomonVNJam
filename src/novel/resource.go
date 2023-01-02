package novel

import (
	"image"
	"image/color"
	"nomon/src/novel/util"
	"regexp"
	"sort"
	"strings"
)

const (
	slash           string = "/"
	assetsFolder    string = "assets" + slash
	resourcesFolder string = assetsFolder + "resources" + slash
	regexPng        string = "([a-z,0-9]+).png"
)

type Resource interface {
	Load() bool
	Image() *image.Image
	Renew()
	Destroy()
}

type SolidImage struct {
	Url string
	img image.Image
}

func (s *SolidImage) Load() bool {
	s.img = util.GetImage(resourcesFolder + s.Url)
	return true
}

func (s *SolidImage) Image() *image.Image { return &s.img }
func (s *SolidImage) Destroy()            { *s = SolidImage{} }
func (s *SolidImage) Renew()              {}

// ++++++++++++++
type SolidColor struct {
	Color  color.Color
	bColor color.Color
}

func (s *SolidColor) Load() bool {
	s.bColor = s.Color
	return s.Color != nil
}

func (s *SolidColor) Image() *image.Image {
	img := s.createImg()
	return &img
}

func (s *SolidColor) createImg() image.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{20, 20}})
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			img.Set(x, y, s.Color)
		}
	}

	// img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	// img.Set(0, 0, s.Color)
	return img
}

func (s *SolidColor) Destroy() { *s = SolidColor{} }
func (s *SolidColor) Renew()   { s.Color = s.bColor }

// =====================

type Oneframe struct {
	Url   string
	imgs  []image.Image
	index int
}

func (g *Oneframe) Load() bool {
	if err := util.ReadDir(resourcesFolder+g.Url, func(isDir bool, nameFile string) bool {
		if isDir {
			return true
		} else {
			if match, err := regexp.MatchString(regexPng, nameFile); match && err == nil {
				img := util.GetImage(resourcesFolder + g.Url + slash + nameFile)
				g.imgs = append(g.imgs, img)
			} else {
				panic(err)
			}
		}
		return false
	}); err != nil {
		return false
	}
	return true
}

func (g *Oneframe) Image() *image.Image { return &g.imgs[g.index] }

func (s *Oneframe) Destroy() { *s = Oneframe{} }
func (s *Oneframe) Renew()   {}

// ====================

type SimpleString struct {
	Text, bText           *string
	Size, bSize           float64
	ColorFont, bColorFont color.Color
	ColorBG, bColorBG     color.RGBA
}

func (s *SimpleString) Load() bool {
	if s.Text == nil {
		panic("Error novel: pointer of String is nil.")
	}
	s.bText = s.Text
	s.bSize = s.Size
	s.bColorFont = s.ColorFont
	s.bColorBG = s.ColorBG
	return true
}

func (s *SimpleString) Image() *image.Image {
	text := (*s.Text)
	if text == "" {
		text = " "
	}
	x := util.StringtoImageRect(text, s.Size, s.ColorFont, s.ColorBG)
	return &x
}

func (s *SimpleString) Destroy() { *s = SimpleString{} }
func (s *SimpleString) Renew() {
	s.Text = s.bText
	s.Size = s.bSize
	s.ColorFont = s.bColorFont
	s.ColorBG = s.bColorBG
}

type GroupImgFor struct {
	Url               string
	Current, bCurrent string
	imgs              map[string][]image.Image
	index             int
}

func (g *GroupImgFor) SetCurrent(c string) {
	g.index = 0
	g.bCurrent, g.Current = c, c
}

func (g *GroupImgFor) Load() bool {
	g.bCurrent = g.Current
	g.imgs = make(map[string][]image.Image)
	if err := g.readDir(resourcesFolder+g.Url, regexPng); err != nil {
		println(err)
		return false
	}
	return true
}

func (g *GroupImgFor) Image() *image.Image { return &g.imgs[g.Current][g.index] }

func (g *GroupImgFor) readDir(url, regeStr string) error {
	var group []image.Image
	var keys []string
	index := strings.Split(url, slash)
	defer func() {
		if len(group) != 0 {
			sort.Slice(group[:], func(i, j int) bool {
				return keys[i] < keys[j]
			})
			g.imgs[index[len(index)-1]] = group
		}
	}()
	return util.ReadDir(url, func(isDir bool, nameFile string) bool {
		if isDir {
			g.readDir(url+slash+nameFile, regeStr)
		} else {
			if match, err := regexp.MatchString(regeStr, nameFile); match && err == nil {
				img := util.GetImageEbiten(url + slash + nameFile)
				group = append(group, img)
				keys = append(keys, nameFile)
			} else {
				panic(err)
			}
		}
		return false
	})
}

func (g *GroupImgFor) Renew() {
	g.Current = g.bCurrent
	g.index = 0
}

func (g *GroupImgFor) Destroy() { *g = GroupImgFor{} }
