# Ebiten Gamepad Viewer

A gamepad/controller test utility built with [Ebiten](https://ebitengine.org/). Displays real-time input state for connected controllers including buttons, D-pad, analog sticks, and trigger values. Also supports vibration/rumble testing.

## Features

- Visual display of all standard gamepad inputs
- Analog stick position with X/Y axis values
- Trigger values (L2/R2) displayed numerically
- Vibration testing with multiple presets (strong, weak, both, stop)
- Button press/release logging to console
- Supports desktop (macOS, Windows, Linux) and iOS

## Requirements

- Go 1.25+
- For iOS: Xcode and ebitenmobile

## Building

### Desktop

```
make desktop
```

The binary is output to `build/ebiten_gamepad_viewer`.

### iOS

```
make ios
```

This produces `build/Iosgamepadviewer.xcframework` which is used by the Xcode project in `ios/xcode/`.

To build the iOS app:

1. Copy `ios/xcode/Signing.xcconfig.template` to `ios/xcode/Signing.xcconfig`
2. Edit `Signing.xcconfig` with your Apple development team ID
3. Open `ios/xcode/EbitenGamepadViewer.xcodeproj` in Xcode and build

### Both

```
make all
```

### Clean

```
make clean
```

## Usage

Run the desktop binary or launch the iOS app and connect a controller. The viewer will display:

- **D-pad** - directional button state
- **Face buttons** - A, B, X, Y state
- **Shoulder buttons** - L1, R1, L2, R2 with trigger values
- **Center buttons** - Select, Start, L3, R3
- **Analog sticks** - position within range with axis values
- **Vibration controls** - click/tap the preset buttons at the bottom to test rumble
