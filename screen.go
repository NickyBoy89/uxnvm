package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScreenDevice controls access to the Uxn machine's screen
// Reference: https://wiki.xxiivv.com/site/varvara.html#screen
var ScreenDevice = Device{
	ReadByte: func(d *Device, port byte) byte {
		switch port {
		case 0x2:
			//return uxn_screen.width >> 8
		case 0x3:
			//return uxn_screen.width
		case 0x4:
			//return uxn_screen.height >> 8
		case 0x5:
			//return uxn_screen.height
		}
		return d.Data[port]
	},
	WriteByte: func(d *Device, port byte) {
		switch port {
		case 0xe: // Write a pixel to the screen
			panic("Pixel")
		default:
			panic(fmt.Sprintf("Unhandled screen port: %.2x\n", port))
		}
	},
}

const (
	ScreenWidth  = 64 * 8
	ScreenHeight = 40 * 8
)

type UxnScreen struct {
}

func (us *UxnScreen) Update() error {
	return nil
}

func (us *UxnScreen) Draw(screen *ebiten.Image) {
}

func (us *UxnScreen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
