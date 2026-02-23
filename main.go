package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/user-none/ebiten_gamepad_viewer/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Controller Test")
	if err := ebiten.RunGame(&game.Game{}); err != nil {
		log.Fatal(err)
	}
}
