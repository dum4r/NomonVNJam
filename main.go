package main

import (
	"embed"
	"math/rand"
	"nomon/src/novel"
	"nomon/src/scenes"
	"nomon/src/scenes/lang"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed assets/*
	assetsFS embed.FS
	//go:embed icon.png
	iconFile []byte
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	nomon, err := novel.Start(
		&assetsFS,
		&iconFile,
		"Nomon",
		&scenes.Player,
		map[string]any{
			"menu": &lang.Menu,
		},
		novel.AddFont("IBM-Plex-Mono/IBMPlexMono-Regular.ttf", '$'),
		novel.AddFont("IBM-Plex-Mono/IBMPlexMono-Bold.ttf", '*'),
		novel.AddFont("IBM-Plex-Mono/IBMPlexMono-Thin.ttf", '|'),
	)
	if err != nil {
		panic(err)
	}

	scenes.FullScreen = &nomon.Fullscreen
	nomon.SetCurtain(&scenes.Curtain)
	novel.Quit = scenes.Curtain.Quit

	novel.AddScene("Menu", &scenes.Menu{})
	novel.AddScene("Config", &scenes.Config{})

	for name, scene := range scenes.Scenes {
		novel.AddScene(name, scene)
	}

	if err := ebiten.RunGame(nomon); err != nil && err != novel.Termination {
		panic(err)
	}
}
