package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Error: Need to specify an input rom, `command [rom-name.rom]`")
	}

	// Load the rom from disk
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Create the Uxn virtual machine
	var uxn Uxn

	uxn.AddDevice(0x0, SystemDevice)  // System
	uxn.AddDevice(0x1, ConsoleDevice) // Console
	uxn.AddDevice(0x2, DummyDevice)   // Screen
	uxn.AddDevice(0x3, DummyDevice)   // Audio
	uxn.AddDevice(0x4, DummyDevice)   // Audio
	uxn.AddDevice(0x5, DummyDevice)   // Audio
	uxn.AddDevice(0x6, DummyDevice)   // Audio
	uxn.AddDevice(0x7, DummyDevice)   // MIDI
	uxn.AddDevice(0x8, DummyDevice)   // Controller
	uxn.AddDevice(0x9, DummyDevice)   // Mouse
	uxn.AddDevice(0xa, DummyDevice)   // File
	uxn.AddDevice(0xb, DummyDevice)   // File
	uxn.AddDevice(0xc, DummyDevice)   // Datetime
	uxn.AddDevice(0xd, DummyDevice)   // Empty
	uxn.AddDevice(0xe, DummyDevice)   // Reserved
	uxn.AddDevice(0xf, DummyDevice)   // Reserved

	// Load the rom into the create Uxn virtual machine
	uxn.Load(input)

	// Execute the instructions one at a time
	for !uxn.Halted {
		uxn.Execute()
	}
}
