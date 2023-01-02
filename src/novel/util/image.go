package util

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func GetImage(filePath string) image.Image {
	file, err := assets.Open(filePath)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func GetImageEbiten(urlIntroName string) *ebiten.Image {
	return ebiten.NewImageFromImage(GetImage(urlIntroName))
}

// Create a loop in the folder and internally the function handles the items found
// internal func returns true if there is an error
func ReadDir(nameFolder string, f func(bool, string) bool) (err error) {
	str := ""
	if dir, err := assets.ReadDir(nameFolder); err == nil {
		for _, file := range dir {
			isDir := file.IsDir()
			strName := file.Name()
			if f(isDir, strName) {
				str += fmt.Sprintf("ReadDir error %d (isDirr:%d - name:%d)\n", nameFolder, isDir, strName)
				err = errors.New(str)
			}
		}
	}
	return err
}

func ImageForBytes(file []byte) image.Image {
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}
	return img
}
