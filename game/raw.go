package game

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	axisLogThreshold = 0.1

	rawGridX         = 30
	rawButtonGridY   = 60
	rawButtonsPerRow = 8
	rawButtonW       = 85
	rawButtonH       = 26
	rawButtonGap     = 10
	maxRawButtons    = 40

	rawAxesPerRow = 2
	rawAxisColW   = 380
	rawAxisRowH   = 26
	rawAxisBarW   = 260
	rawAxisBarH   = 14
	maxRawAxes    = 10
)

// axisLogNeeded reports whether an axis value has moved far enough from the
// last logged value to be worth logging again.
func axisLogNeeded(prev, cur float64) bool {
	return math.Abs(cur-prev) >= axisLogThreshold
}

// rawButtonPos returns the top-left corner of raw button i in the grid.
func rawButtonPos(i int) (float32, float32) {
	col := i % rawButtonsPerRow
	row := i / rawButtonsPerRow
	x := rawGridX + col*(rawButtonW+rawButtonGap)
	y := rawButtonGridY + row*(rawButtonH+rawButtonGap)
	return float32(x), float32(y)
}

// rawAxisPos returns the top-left corner of raw axis row i, starting below a
// button grid holding buttonCount buttons.
func rawAxisPos(i, buttonCount int) (float32, float32) {
	if buttonCount > maxRawButtons {
		buttonCount = maxRawButtons
	}
	btnRows := (buttonCount + rawButtonsPerRow - 1) / rawButtonsPerRow
	baseY := rawButtonGridY + btnRows*(rawButtonH+rawButtonGap) + 20

	col := i % rawAxesPerRow
	row := i / rawAxesPerRow
	x := rawGridX + col*rawAxisColW
	y := baseY + row*rawAxisRowH
	return float32(x), float32(y)
}

// logRawInputs logs raw button presses/releases and axis movement to the
// console. This works for any gamepad, with or without a standard layout
// mapping, and the indices match those used in mapping strings (bN/aN).
// D-pads exposed as hats have no raw API and will not appear here.
func (g *Game) logRawInputs() {
	id := g.gamepadID

	for b := 0; b < ebiten.GamepadButtonCount(id); b++ {
		btn := ebiten.GamepadButton(b)
		if inpututil.IsGamepadButtonJustPressed(id, btn) {
			log.Printf("raw pressed:  b%d", b)
		}
		if inpututil.IsGamepadButtonJustReleased(id, btn) {
			log.Printf("raw released: b%d", b)
		}
	}

	axisCount := ebiten.GamepadAxisCount(id)
	if len(g.loggedAxes) != axisCount {
		g.loggedAxes = make([]float64, axisCount)
	}
	for a := 0; a < axisCount; a++ {
		v := ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(a))
		if axisLogNeeded(g.loggedAxes[a], v) {
			log.Printf("raw axis: a%d=%+.2f", a, v)
			g.loggedAxes[a] = v
		}
	}
}

// drawRawView draws every raw button and axis of the connected gamepad so
// controllers without a standard layout mapping can still be exercised.
func (g *Game) drawRawView(screen *ebiten.Image) {
	id := g.gamepadID
	buttonCount := ebiten.GamepadButtonCount(id)
	axisCount := ebiten.GamepadAxisCount(id)

	ebitenutil.DebugPrintAt(screen,
		fmt.Sprintf("Buttons: %d  Axes: %d (raw indices match mapping strings)", buttonCount, axisCount),
		rawGridX, 32)

	shownButtons := buttonCount
	if shownButtons > maxRawButtons {
		shownButtons = maxRawButtons
	}
	for b := 0; b < shownButtons; b++ {
		x, y := rawButtonPos(b)
		drawButton(screen, x, y, rawButtonW, rawButtonH, fmt.Sprintf("b%d", b),
			ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(b)))
	}
	if buttonCount > maxRawButtons {
		x, y := rawAxisPos(0, buttonCount)
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("(showing %d of %d buttons)", maxRawButtons, buttonCount), int(x), int(y)-16)
	}

	shownAxes := axisCount
	if shownAxes > maxRawAxes {
		shownAxes = maxRawAxes
	}
	for a := 0; a < shownAxes; a++ {
		x, y := rawAxisPos(a, buttonCount)
		v := ebiten.GamepadAxisValue(id, ebiten.GamepadAxisType(a))
		drawAxisBar(screen, x, y, v, fmt.Sprintf("a%d %+.2f", a, v))
	}
}

// drawAxisBar draws a labeled horizontal bar for an axis value in [-1, 1].
func drawAxisBar(screen *ebiten.Image, x, y float32, v float64, label string) {
	ebitenutil.DebugPrintAt(screen, label, int(x), int(y))

	barX := x + 90
	vector.FillRect(screen, barX, y, rawAxisBarW, rawAxisBarH, colorStickBg, false)
	vector.FillRect(screen, barX+rawAxisBarW/2-1, y, 2, rawAxisBarH, colorStickOutline, false)

	clamped := math.Max(-1, math.Min(1, v))
	markerX := barX + rawAxisBarW/2 + float32(clamped)*(rawAxisBarW/2-3)
	vector.FillRect(screen, markerX-3, y, 6, rawAxisBarH, colorStickDot, false)
}
