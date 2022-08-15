package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	input, err := os.ReadFile("hello.rom")
	if err != nil {
		panic(err)
	}

	var uxn Uxn

	uxn.Devices[0x0] = Device{
		DeviceInput: func(d *Device, port byte) byte {
			switch port {
			case 0x2:
				return d.u.WorkingStack.Pointer
			case 0x3:
				return d.u.ReturnStack.Pointer
			default:
				return d.Data[port]
			}
		},
		DeviceOut: func(d *Device, port byte) {
			switch port {
			case 0x2:
				d.u.WorkingStack.Pointer = d.Data[port]
			case 0x3:
				d.u.ReturnStack.Pointer = d.Data[port]
			case 0xe:
				panic("system_inspect")
				//system_inspect(d.u)
			default:
				//panic("system_deo_special")
				//system_deo_special(d, port)
			}
		},
	}

	uxn.Devices[0x1] = Device{
		DeviceInput: NilDei,
		DeviceOut: func(d *Device, port byte) {
			var out io.Writer
			switch port {
			case 0x8:
				out = os.Stdout
			case 0x9:
				out = os.Stderr
			}

			if out != nil {
				out.Write([]byte{d.Data[port]})
			}
		},
	}

	uxn.Load(input)
	for {
		uxn.Execute()
		fmt.Println(uxn.WorkingStack)
	}
}
