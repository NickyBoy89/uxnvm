package main

import (
	"fmt"
	"strings"
)

func HexPrint(arr []byte) string {
	var result strings.Builder
	result.WriteRune('[')
	for index, item := range arr {
		fmt.Fprintf(&result, "%.2x", item)
		if index < len(arr)-1 {
			result.WriteRune(' ')
		}
	}
	result.WriteRune(']')
	return result.String()
}

// Page of memory where the program starts executing
const ProgramStartPage uint16 = 0x100

type Uxn struct {
	// The stacks that the machine uses
	WorkingStack, ReturnStack Stack
	// Swapped during return mode
	Src, Dst *Stack
	// A list of external devices that the machine can access
	Devices [16]Device
	// 64k of memory
	Memory [65536]byte
	// The current element in memory
	ProgramCounter uint16
	// Whether the program should continue executing
	Halted bool
}

func (u *Uxn) Poke8(at uint16, data byte) {
	u.Memory[at] = data
}

func (u *Uxn) Poke16(at uint16, data uint16) {
	u.Memory[at] = byte(data >> 8)
	u.Memory[at+1] = byte(data)
}

func (u *Uxn) Peek16(at uint16) uint16 {
	return uint16(u.Memory[at])<<8 + uint16(u.Memory[at+1])
}

func (u *Uxn) Peek8(at uint16) byte {
	return u.Memory[at]
}

func (u *Uxn) Warp8(x byte) {
	s := int8(x)
	if s < 0 {
		u.ProgramCounter -= uint16(-s)
	} else {
		u.ProgramCounter += uint16(s)
	}
}

func (u *Uxn) Warp16(x uint16) {
	u.ProgramCounter = x
}

// Execute takes a single byte from the where the Program Counter is pointing in
// memory and executes it
func (u *Uxn) Execute() {
	instr := u.Memory[u.ProgramCounter]
	u.ProgramCounter++
	// Return Mode
	if instr&0x40 != 0 {
		u.Src = &u.ReturnStack
		u.Dst = &u.WorkingStack
	} else {
		u.Src = &u.WorkingStack
		u.Dst = &u.ReturnStack
	}

	// Short Mode
	shortMode := instr&0x20 != 0

	// A pointer to the current source stack
	var srcStackPtr *byte

	// Keep Mode
	if instr&0x80 != 0 {
		temp := u.Src.Pointer
		srcStackPtr = &temp
	} else {
		srcStackPtr = &u.Src.Pointer
	}

	// Get the top 5 bytes of the instruction
	switch instr & 0x1f {
	/* Stack */
	case 0x00: /* LIT */
		if shortMode {
			a := u.Peek16(u.ProgramCounter)
			u.Src.Push16(a)
		} else {
			a := u.Peek8(u.ProgramCounter)
			u.Src.Push8(a)
		}
		u.ProgramCounter++
		if shortMode {
			u.ProgramCounter++
		}
	case 0x01: /* INC */
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(a + 1)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(a + 1)
		}
	case 0x02: /* POP */
		if shortMode {
			u.Src.Pop16(srcStackPtr)
		} else {
			u.Src.Pop8(srcStackPtr)
		}
	case 0x03: /* NIP */
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Src.Pop16(srcStackPtr)
			u.Src.Push16(a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Src.Pop8(srcStackPtr)
			u.Src.Push8(a)
		}
	case 0x04: // SWP\
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(a)
			u.Src.Push16(b)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(a)
			u.Src.Push8(b)
		}
	case 0x05: // ROT
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			c := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b)
			u.Src.Push16(a)
			u.Src.Push16(c)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			c := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b)
			u.Src.Push8(a)
			u.Src.Push8(c)
		}
	case 0x06: // DUP
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(a)
			u.Src.Push16(a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(a)
			u.Src.Push8(a)
		}
	case 0x07: // OVR
		if shortMode {
			//file_instance
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b)
			u.Src.Push16(a)
			u.Src.Push16(b)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b)
			u.Src.Push8(a)
			u.Src.Push8(b)
		}
	/* Logic */
	case 0x08: // EQU
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			if b == a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b == a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		}
	case 0x09: // NEQ
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			if b != a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b != a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		}
	case 0x0a: // GTH
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			if b > a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b > a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		}
	case 0x0b: // LTH
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			if b < a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b < a {
				u.Src.Push8(0x01)
			} else {
				u.Src.Push8(0x00)
			}
		}
	case 0x0c: // JMP
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Warp16(a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Warp8(a)
		}
	case 0x0d: // JCN
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b != 0x00 {
				u.Warp16(a)
			}
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if b != 0x00 {
				u.Warp8(a)
			}
		}
	case 0x0e: // JSR
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Dst.Push16(u.ProgramCounter)
			u.Warp16(a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Dst.Push16(u.ProgramCounter)
			u.Warp8(a)
		}
	case 0x0f: // STH
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			u.Dst.Push16(a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			u.Dst.Push8(a)
		}
		/* Memory */
	case 0x10: // LDZ
		a := uint16(u.Src.Pop8(srcStackPtr))
		if shortMode {
			b := u.Peek16(a)
			u.Src.Push16(b)
		} else {
			b := u.Peek8(a)
			u.Src.Push8(b)
		}
	case 0x11: // STZ
		a := uint16(u.Src.Pop8(srcStackPtr))
		if shortMode {
			b := u.Src.Pop16(srcStackPtr)
			u.Poke16(a, b)
		} else {
			b := u.Src.Pop8(srcStackPtr)
			u.Poke8(a, b)
		}
	case 0x12: // LDR
		a := int8(u.Src.Pop8(srcStackPtr))
		if shortMode {
			var b uint16
			if a < 0 {
				b = u.Peek16(uint16(*srcStackPtr - byte(-a)))
			} else {
				b = u.Peek16(uint16(*srcStackPtr + byte(a)))
			}
			u.Src.Push16(b)
		} else {
			var b byte
			if a < 0 {
				b = u.Peek8(uint16(*srcStackPtr - byte(-a)))
			} else {
				b = u.Peek8(uint16(*srcStackPtr + byte(a)))
			}
			u.Src.Push8(b)
		}
	case 0x13: // STR
		a := int8(u.Src.Pop8(srcStackPtr))
		c := *srcStackPtr
		if a < 0 {
			c -= byte(-a)
		} else {
			c += byte(a)
		}
		if shortMode {
			b := u.Src.Pop16(srcStackPtr)
			u.Poke16(uint16(c), b)
		} else {
			b := u.Src.Pop8(srcStackPtr)
			u.Poke8(uint16(c), b)
		}
	case 0x14: // LDA
		a := u.Src.Pop16(srcStackPtr)
		if shortMode {
			b := u.Peek16(a)
			u.Src.Push16(b)
		} else {
			b := u.Peek8(a)
			u.Src.Push8(b)
		}
	case 0x15: // STA
		a := u.Src.Pop16(srcStackPtr)
		if shortMode {
			b := u.Src.Pop16(srcStackPtr)
			u.Poke16(a, b)
		} else {
			b := u.Src.Pop8(srcStackPtr)
			u.Poke8(a, b)
		}
	case 0x16: // DEI
		deviceIndex := u.Src.Pop8(srcStackPtr)
		if u.Devices[deviceIndex>>4].u == nil {
			panic(fmt.Sprintf("Device does not exist: %.2x", deviceIndex))
		}
		if shortMode {
			b := u.Devices[deviceIndex>>4].DeviceRead16(deviceIndex)
			u.Src.Push16(b)
		} else {
			b := u.Devices[deviceIndex>>4].DeviceRead8(deviceIndex)
			u.Src.Push8(b)
		}
	case 0x17: // DEO
		deviceIndex := u.Src.Pop8(srcStackPtr)
		if u.Devices[deviceIndex>>4].u == nil {
			panic(fmt.Sprintf("Device does not exist: %.2x", deviceIndex))
		}
		if shortMode {
			b := u.Src.Pop16(srcStackPtr)
			u.Devices[deviceIndex>>4].DeviceWrite16(deviceIndex, b)
		} else {
			b := u.Src.Pop8(srcStackPtr)
			u.Devices[deviceIndex>>4].DeviceWrite8(deviceIndex, b)
		}
	// Arithmetic
	case 0x18: // ADD
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b + a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b + a)
		}
	case 0x19: // SUB
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b - a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b - a)
		}
	case 0x1a: // MUL
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b * a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b * a)
		}
	case 0x1b: // DIV
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			if a == 0 {
				panic(ErrDivByZero)
			}
			u.Src.Push16(b / a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			if a == 0 {
				panic(ErrDivByZero)
			}
			u.Src.Push8(b / a)
		}
	case 0x1c: // AND
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b & a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b & a)
		}
	case 0x1d: // ORA
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b | a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b | a)
		}
	case 0x1e: // EOR
		if shortMode {
			a := u.Src.Pop16(srcStackPtr)
			b := u.Src.Pop16(srcStackPtr)
			u.Src.Push16(b ^ a)
		} else {
			a := u.Src.Pop8(srcStackPtr)
			b := u.Src.Pop8(srcStackPtr)
			u.Src.Push8(b ^ a)
		}
	case 0x1f: // SFT
		a := u.Src.Pop8(srcStackPtr)
		if shortMode {
			b := u.Src.Pop16(srcStackPtr)
			c := b >> (a & 0x0f) << ((a & 0xf0) >> 4)
			u.Src.Push16(c)
		} else {
			b := u.Src.Pop8(srcStackPtr)
			c := b >> (a & 0x0f) << ((a & 0xf0) >> 4)
			u.Src.Push8(c)
		}
	default:
		panic(fmt.Sprintf("Unhandled instruction %.2x", instr&0x1f))
	}
}

// AddDevice links a device to a `uxn` virtual machine at the given port
func (u *Uxn) AddDevice(port byte, device Device) {
	u.Devices[port] = device
	u.Devices[port].u = u
}

// Load takes in a `uxn` rom and loads it into memory to be executed
func (u *Uxn) Load(rom []byte) {
	for offset := 0; offset < len(rom); offset++ {
		u.Memory[ProgramStartPage+uint16(offset)] = rom[offset]
	}
	u.ProgramCounter = ProgramStartPage
}
