package main

import (
	"os"
)

func main() {
	input, err := os.ReadFile("hello.rom")
	if err != nil {
		panic(err)
	}

	var uxn Machine

	uxn.Load(input)
	for {
		uxn.Execute()
	}
}
