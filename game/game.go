package game

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

var (
	colorBackground   = color.RGBA{0x1a, 0x1a, 0x2e, 0xff}
	colorInactive     = color.RGBA{0x55, 0x55, 0x55, 0xff}
	colorActive       = color.RGBA{0x00, 0xdd, 0x00, 0xff}
	colorStickBg      = color.RGBA{0x33, 0x33, 0x33, 0xff}
	colorStickDot     = color.RGBA{0xff, 0xdd, 0x00, 0xff}
	colorVibrate      = color.RGBA{0x33, 0x66, 0xcc, 0xff}
	colorVibrateHover = color.RGBA{0x44, 0x88, 0xff, 0xff}
	colorStickOutline = color.RGBA{0x66, 0x66, 0x66, 0xff}
)

type vibratePreset struct {
	label  string
	strong float64
	weak   float64
}

var vibratePresets = []vibratePreset{
	{"Strong Only", 1.0, 0},
	{"Weak Only", 0, 1.0},
	{"Both 50%", 0.5, 0.5},
	{"Both 100%", 1.0, 1.0},
	{"Stop", 0, 0},
}

var standardButtonNames = map[ebiten.StandardGamepadButton]string{
	ebiten.StandardGamepadButtonRightBottom:      "A",
	ebiten.StandardGamepadButtonRightRight:       "B",
	ebiten.StandardGamepadButtonRightLeft:        "X",
	ebiten.StandardGamepadButtonRightTop:         "Y",
	ebiten.StandardGamepadButtonFrontTopLeft:     "L1",
	ebiten.StandardGamepadButtonFrontBottomLeft:  "L2",
	ebiten.StandardGamepadButtonFrontTopRight:    "R1",
	ebiten.StandardGamepadButtonFrontBottomRight: "R2",
	ebiten.StandardGamepadButtonCenterLeft:       "Select",
	ebiten.StandardGamepadButtonCenterRight:      "Start",
	ebiten.StandardGamepadButtonLeftStick:        "L3",
	ebiten.StandardGamepadButtonRightStick:       "R3",
	ebiten.StandardGamepadButtonLeftTop:          "Up",
	ebiten.StandardGamepadButtonLeftBottom:       "Down",
	ebiten.StandardGamepadButtonLeftLeft:         "Left",
	ebiten.StandardGamepadButtonLeftRight:        "Right",
}

type Game struct {
	gamepadID    ebiten.GamepadID
	connected    bool
	gamepadIDBuf []ebiten.GamepadID
	touchIDBuf   []ebiten.TouchID
	loggedAxes   []float64
}

func (g *Game) Update() error {
	g.gamepadIDBuf = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIDBuf[:0])
	for _, id := range g.gamepadIDBuf {
		if !g.connected {
			g.gamepadID = id
			g.connected = true
			name := ebiten.GamepadName(id)
			sdlID := ebiten.GamepadSDLID(id)
			standard := ebiten.IsStandardGamepadLayoutAvailable(id)
			log.Printf("gamepad connected: id=%d name=%q sdl_id=%s standard=%v buttons=%d axes=%d",
				id, name, sdlID, standard, ebiten.GamepadButtonCount(id), ebiten.GamepadAxisCount(id))
			g.loggedAxes = nil
			break
		}
	}

	if g.connected && inpututil.IsGamepadJustDisconnected(g.gamepadID) {
		log.Printf("gamepad disconnected: id=%d", g.gamepadID)
		g.connected = false
	}

	if g.connected {
		g.logInputs()
		g.checkVibrateInput()
	}

	return nil
}

func (g *Game) checkVibrateInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		g.handleVibrate(cx, cy)
	}

	g.touchIDBuf = inpututil.AppendJustPressedTouchIDs(g.touchIDBuf[:0])
	for _, id := range g.touchIDBuf {
		tx, ty := ebiten.TouchPosition(id)
		g.handleVibrate(tx, ty)
	}
}

func (g *Game) logInputs() {
	if ebiten.IsStandardGamepadLayoutAvailable(g.gamepadID) {
		for btn := ebiten.StandardGamepadButton(0); btn <= ebiten.StandardGamepadButtonMax; btn++ {
			name := standardButtonNames[btn]
			if name == "" {
				name = fmt.Sprintf("button(%d)", btn)
			}
			if inpututil.IsStandardGamepadButtonJustPressed(g.gamepadID, btn) {
				val := ebiten.StandardGamepadButtonValue(g.gamepadID, btn)
				log.Printf("pressed:  %s (value=%.2f)", name, val)
			}
			if inpututil.IsStandardGamepadButtonJustReleased(g.gamepadID, btn) {
				log.Printf("released: %s", name)
			}
		}
	}
	g.logRawInputs()
}

func (g *Game) handleVibrate(cx, cy int) {
	if g.handleVibrateRow(cx, cy, vibratePresets, 505, math.MaxInt64) {
		return
	}
	g.handleVibrateRow(cx, cy, vibratePresets[:len(vibratePresets)-1], 545, 3*time.Second)
}

func (g *Game) handleVibrateRow(cx, cy int, presets []vibratePreset, y int, duration time.Duration) bool {
	btnW := 120
	btnH := 30
	totalW := len(presets)*btnW + (len(presets)-1)*10
	startX := (ScreenWidth - totalW) / 2

	for i, p := range presets {
		bx := startX + i*(btnW+10)
		if cx >= bx && cx <= bx+btnW && cy >= y && cy <= y+btnH {
			d := duration
			if p.strong == 0 && p.weak == 0 {
				d = 0
			}
			log.Printf("vibrate: %s strong=%.2f weak=%.2f duration=%s", p.label, p.strong, p.weak, d)
			ebiten.VibrateGamepad(g.gamepadID, &ebiten.VibrateGamepadOptions{
				Duration:        d,
				StrongMagnitude: p.strong,
				WeakMagnitude:   p.weak,
			})
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	g.drawHeader(screen)

	if !g.connected {
		return
	}

	if !ebiten.IsStandardGamepadLayoutAvailable(g.gamepadID) {
		g.drawRawView(screen)
		g.drawVibrateButtons(screen)
		return
	}

	g.drawDPad(screen)
	g.drawFaceButtons(screen)
	g.drawShoulderRow(screen)
	g.drawStickButtons(screen)
	g.drawAnalogSticks(screen)
	g.drawVibrateButtons(screen)
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	if !g.connected {
		ebitenutil.DebugPrintAt(screen, "No Gamepad Connected - Please connect a controller", 200, 10)
		return
	}

	name := ebiten.GamepadName(g.gamepadID)
	if ebiten.IsStandardGamepadLayoutAvailable(g.gamepadID) {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gamepad: %s (Standard Layout)", name), 10, 10)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Gamepad: %s (Non-Standard Layout - raw inputs)", name), 10, 10)
	}
}

func (g *Game) drawDPad(screen *ebiten.Image) {
	baseX := 120
	baseY := 130
	bw := float32(50)
	bh := float32(35)
	gap := float32(5)

	drawButton(screen, float32(baseX)-bw/2, float32(baseY)-bh-gap, bw, bh, "Up",
		g.isPressed(ebiten.StandardGamepadButtonLeftTop))
	drawButton(screen, float32(baseX)-bw/2, float32(baseY)+gap, bw, bh, "Down",
		g.isPressed(ebiten.StandardGamepadButtonLeftBottom))
	drawButton(screen, float32(baseX)-bw-gap-bw/2, float32(baseY)-bh/2, bw, bh, "Left",
		g.isPressed(ebiten.StandardGamepadButtonLeftLeft))
	drawButton(screen, float32(baseX)+gap-bw/2+bw, float32(baseY)-bh/2, bw, bh, "Right",
		g.isPressed(ebiten.StandardGamepadButtonLeftRight))
}

func (g *Game) drawFaceButtons(screen *ebiten.Image) {
	baseX := float32(650)
	baseY := float32(130)
	r := float32(22)
	gap := float32(50)

	drawCircleButton(screen, baseX, baseY-gap, r, "Y",
		g.isPressed(ebiten.StandardGamepadButtonRightTop))
	drawCircleButton(screen, baseX, baseY+gap, r, "A",
		g.isPressed(ebiten.StandardGamepadButtonRightBottom))
	drawCircleButton(screen, baseX-gap, baseY, r, "X",
		g.isPressed(ebiten.StandardGamepadButtonRightLeft))
	drawCircleButton(screen, baseX+gap, baseY, r, "B",
		g.isPressed(ebiten.StandardGamepadButtonRightRight))
}

func (g *Game) drawShoulderRow(screen *ebiten.Image) {
	y := float32(250)
	bw := float32(55)
	bh := float32(30)

	drawButton(screen, 60, y, bw, bh, "L1",
		g.isPressed(ebiten.StandardGamepadButtonFrontTopLeft))
	l2val := ebiten.StandardGamepadButtonValue(g.gamepadID, ebiten.StandardGamepadButtonFrontBottomLeft)
	drawButton(screen, 125, y, bw, bh, fmt.Sprintf("L2 %.2f", l2val),
		g.isPressed(ebiten.StandardGamepadButtonFrontBottomLeft))

	drawButton(screen, 300, y, 60, bh, "Select",
		g.isPressed(ebiten.StandardGamepadButtonCenterLeft))
	drawButton(screen, 370, y, 60, bh, "Start",
		g.isPressed(ebiten.StandardGamepadButtonCenterRight))

	drawButton(screen, 620, y, bw, bh, "R1",
		g.isPressed(ebiten.StandardGamepadButtonFrontTopRight))
	r2val := ebiten.StandardGamepadButtonValue(g.gamepadID, ebiten.StandardGamepadButtonFrontBottomRight)
	drawButton(screen, 685, y, bw, bh, fmt.Sprintf("R2 %.2f", r2val),
		g.isPressed(ebiten.StandardGamepadButtonFrontBottomRight))
}

func (g *Game) drawStickButtons(screen *ebiten.Image) {
	bw := float32(40)
	bh := float32(28)
	y := float32(400) - bh/2

	drawButton(screen, 355, y, bw, bh, "L3",
		g.isPressed(ebiten.StandardGamepadButtonLeftStick))
	drawButton(screen, 405, y, bw, bh, "R3",
		g.isPressed(ebiten.StandardGamepadButtonRightStick))
}

func (g *Game) drawAnalogSticks(screen *ebiten.Image) {
	stickR := float32(60)

	lx := float32(250)
	ly := float32(400)
	drawStick(screen, lx, ly, stickR, "Left Stick",
		ebiten.StandardGamepadAxisValue(g.gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal),
		ebiten.StandardGamepadAxisValue(g.gamepadID, ebiten.StandardGamepadAxisLeftStickVertical))

	rx := float32(550)
	ry := float32(400)
	drawStick(screen, rx, ry, stickR, "Right Stick",
		ebiten.StandardGamepadAxisValue(g.gamepadID, ebiten.StandardGamepadAxisRightStickHorizontal),
		ebiten.StandardGamepadAxisValue(g.gamepadID, ebiten.StandardGamepadAxisRightStickVertical))
}

func (g *Game) drawVibrateButtons(screen *ebiten.Image) {
	g.drawVibrateRow(screen, vibratePresets, 505, "")
	g.drawVibrateRow(screen, vibratePresets[:len(vibratePresets)-1], 545, " 3s")
}

func (g *Game) drawVibrateRow(screen *ebiten.Image, presets []vibratePreset, y int, suffix string) {
	btnW := 120
	btnH := 30
	totalW := len(presets)*btnW + (len(presets)-1)*10
	startX := (ScreenWidth - totalW) / 2

	cx, cy := ebiten.CursorPosition()

	for i, p := range presets {
		bx := startX + i*(btnW+10)
		hovered := cx >= bx && cx <= bx+btnW && cy >= y && cy <= y+btnH

		clr := colorVibrate
		if hovered {
			clr = colorVibrateHover
		}

		vector.FillRect(screen, float32(bx), float32(y), float32(btnW), float32(btnH), clr, false)

		text := p.label + suffix
		labelX := bx + (btnW-len(text)*6)/2
		labelY := y + 8
		ebitenutil.DebugPrintAt(screen, text, labelX, labelY)
	}
}

func (g *Game) isPressed(btn ebiten.StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(g.gamepadID, btn)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func drawButton(screen *ebiten.Image, x, y, w, h float32, label string, active bool) {
	clr := colorInactive
	if active {
		clr = colorActive
	}
	vector.FillRect(screen, x, y, w, h, clr, false)

	labelX := int(x) + (int(w)-len(label)*6)/2
	labelY := int(y) + (int(h)-12)/2
	ebitenutil.DebugPrintAt(screen, label, labelX, labelY)
}

func drawCircleButton(screen *ebiten.Image, cx, cy, r float32, label string, active bool) {
	clr := colorInactive
	if active {
		clr = colorActive
	}
	vector.FillCircle(screen, cx, cy, r, clr, true)

	labelX := int(cx) - len(label)*3
	labelY := int(cy) - 6
	ebitenutil.DebugPrintAt(screen, label, labelX, labelY)
}

func drawStick(screen *ebiten.Image, cx, cy, r float32, label string, axisX, axisY float64) {
	vector.FillCircle(screen, cx, cy, r, colorStickBg, true)
	vector.StrokeCircle(screen, cx, cy, r, 2, colorStickOutline, true)

	dotX := cx + float32(axisX)*r*0.9
	dotY := cy + float32(axisY)*r*0.9
	vector.FillCircle(screen, dotX, dotY, 8, colorStickDot, true)

	ebitenutil.DebugPrintAt(screen, label, int(cx)-len(label)*3, int(cy-r)-18)

	valStr := fmt.Sprintf("X:%+.2f Y:%+.2f", axisX, axisY)
	ebitenutil.DebugPrintAt(screen, valStr, int(cx)-len(valStr)*3, int(cy+r)+8)
}
