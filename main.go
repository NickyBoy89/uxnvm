package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Error: Need to specify an input rom, `command [rom-name.rom]`")
	}
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var uxn Uxn

	uxn.AddDevice(0x0, SystemDevice)  // System
	uxn.AddDevice(0x1, ConsoleDevice) // Console
	uxn.AddDevice(0xa, FileDevice)    // File

	uxn.Load(input)

	for !uxn.Halted {
		uxn.Execute()
		//fmt.Println(uxn.WorkingStack)
	}
}
