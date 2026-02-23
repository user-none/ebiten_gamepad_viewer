package iosgamepadviewer

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/user-none/ebiten_gamepad_viewer/game"
)

func init() {
	mobile.SetGame(&game.Game{})
}

func Dummy() {}
