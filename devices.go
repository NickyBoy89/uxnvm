package main

import (
	"io"
	"os"
)

var SystemDevice = Device{
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
		case 0xe: // Prints the contents of the stacks
			panic("system_inspect")
			//system_inspect(d.u)
		case 0xf: // Halts the program
			d.u.Halted = true
		default:
			//panic("system_deo_special")
			//system_deo_special(d, port)
		}
	},
}

var ConsoleDevice = Device{
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

var FileDevice = Device{
	DeviceInput: func(d *Device, port byte) byte {
		return 0
	},
	DeviceOut: func(d *Device, port byte) {
	},
}
