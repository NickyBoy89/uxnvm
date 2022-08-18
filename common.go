package main

import (
	"io"
	"os"
)

// DummyDevice only exists for protyping purposes, and will always return `0`
// when read from, as well as not do anything when written to
var DummyDevice = Device{
	ReadByte: func(d *Device, port byte) byte {
		return 0
	},
	WriteByte: func(d *Device, port byte) {
	},
}

// The SystemDevice controls the execution of the Uxn system
// Reference: https://wiki.xxiivv.com/site/varvara.html#system
var SystemDevice = Device{
	ReadByte: func(d *Device, port byte) byte {
		switch port {
		case 0x2:
			return d.u.WorkingStack.Pointer
		case 0x3:
			return d.u.ReturnStack.Pointer
		default:
			return d.Data[port]
		}
	},
	WriteByte: func(d *Device, port byte) {
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

// The ConsoleDevice controls input and output from the host
// Reference: https://wiki.xxiivv.com/site/varvara.html#console
var ConsoleDevice = Device{
	ReadByte: func(d *Device, port byte) byte {
		panic("Tried to read from unimplemented device")
	},
	WriteByte: func(d *Device, port byte) {
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
