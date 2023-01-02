package util

import (
	"embed"
	"errors"

	"github.com/golang/freetype/truetype"
)

const (
	slash        = "/"
	folderAssets = "assets"
	folderFonts  = folderAssets + slash + "fonts"
	formatFont   = "([a-z,0-9]+).ttf"
	folderLang   = folderAssets + slash + "lang"
	formatLang   = "([a-z,0-9]+).json"
)

var assets *embed.FS

func SetAssetsFS(newAssets *embed.FS, currenLang *string, tables map[string]any, fonts ...TemplFont) error {
	assets = newAssets
	if assets == nil {
		return errors.New("Resources Error: not find assetsFS")
	}

	if err := initLenguaje(currenLang, tables); err != nil {
		return err
	}

	for _, f := range fonts {
		url, decoration := f.Print()
		file, err := assets.ReadFile(folderFonts + slash + url)
		if err != nil {
			return err
		}
		libraryFonts[decoration], err = truetype.Parse(file)
		if err != nil {
			return err
		}
	}
	return nil
}
