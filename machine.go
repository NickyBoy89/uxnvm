package main

import (
	"errors"
	"fmt"
)

var StackOverflowError = errors.New("Stack overflow!")
var StackUnderflowError = errors.New("Stack underflow!")

// Page of memory where the program starts executing
const ProgramStartPage = 0x100

type Stack struct {
	Data    [8]byte
	Error   byte
	Pointer byte
}

func (s Stack) String() string {
	res := "["
	for i := 0; i < len(s.Data); i++ {
		res += fmt.Sprintf("%.2x", s.Data[i])
		if i+1 < len(s.Data) {
			res += " "
		}
	}
	return res + "]"
}

func (s *Stack) Peek16() uint16 {
	return uint16(s.Data[s.Pointer-1]) + uint16(s.Data[s.Pointer-2])
}

func (s *Stack) Push16(v uint16) {
	s.Data[s.Pointer-2] = byte(v >> 8)
	s.Data[s.Pointer-1] = byte(v & 0xff)
}

// Swap swaps the two elements at index n1 and n2 respectively
func (s *Stack) Swap(n1, n2 byte) {
	s.Data[n1], s.Data[n2] = s.Data[n2], s.Data[n1]
}

type Machine struct {
	// Stacks
	WorkingStack Stack
	ReturnStack  Stack
	IO           struct {
		Devices [256]byte
	}
	Memory         [65536]byte
	ProgramCounter uint16
}

func (m *Machine) PeekMem16() uint16 {
	return uint16(m.Memory[m.ProgramCounter+1])<<8 + uint16(m.Memory[m.ProgramCounter+2])
}

func (m *Machine) Load(rom []byte) {
	for offset := 0; offset < len(rom); offset++ {
		m.Memory[ProgramStartPage+offset] = rom[offset]
	}
	m.ProgramCounter = ProgramStartPage
}

func (m *Machine) Execute() {
	op := m.Memory[m.ProgramCounter]

	fmt.Printf("Executing instr: %.2x\n", op)

	shortMode := op&0x20 != 0
	shortModeInt := op & 0x20 >> 5
	returnMode := op&0x40 != 0
	keepMode := op&0x80 != 0
	// keepModeInt is `1` when keep mode is active, `0` otherwise
	keepModeInt := op & 0x80 >> 7

	_, _, _ = keepMode, returnMode, shortMode

	// Mask the opcode to the last five bits to extract the opcode
	switch op & 0x1f {
	// Stack instructions
	case 0x00: // LIT
		fmt.Println("Pushing literal to the stack")
		if shortMode {
			v := m.PeekMem16()
			m.ProgramCounter += 2
			m.WorkingStack.Pointer += 2
			m.WorkingStack.Push16(v)
			fmt.Println(m.WorkingStack)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.Memory[m.ProgramCounter]
			m.WorkingStack.Pointer++
			m.ProgramCounter++
		}
	case 0x01: // INC
		if shortMode {
			v := m.WorkingStack.Peek16() + 1
			m.ProgramCounter += 2
			if keepMode {
				m.WorkingStack.Pointer += 2
			}
			m.WorkingStack.Push16(v)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer-(1-keepModeInt)] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] + 1
		}
	case 0x02: // POP
		if !keepMode {
			m.WorkingStack.Pointer -= shortModeInt + 1
		}
	case 0x03: // NIP
		if shortMode {
			v := m.WorkingStack.Peek16()
			if keepMode {
				m.WorkingStack.Pointer += 2
			} else {
				m.WorkingStack.Pointer -= 2
			}
			m.WorkingStack.Push16(v)
		} else {
			v := m.WorkingStack.Data[m.WorkingStack.Pointer-1]
			if keepMode {
				m.WorkingStack.Pointer++
			} else {
				m.WorkingStack.Pointer--
			}
			m.WorkingStack.Data[m.WorkingStack.Pointer] = v
		}
	case 0x04: // SWP
		if keepMode {
			if shortMode {
				m.WorkingStack.Pointer += 4
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-3]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-4]
			} else {
				m.WorkingStack.Pointer++
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-2]
			}
		}
		if shortMode {
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-3)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-4)
		} else {
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-2)
		}
	case 0x05: // ROT
		if keepMode {
			if shortMode {
				m.WorkingStack.Pointer += 6
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-7]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-8]
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = m.WorkingStack.Data[m.WorkingStack.Pointer-9]
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = m.WorkingStack.Data[m.WorkingStack.Pointer-10]
				m.WorkingStack.Data[m.WorkingStack.Pointer-5] = m.WorkingStack.Data[m.WorkingStack.Pointer-11]
				m.WorkingStack.Data[m.WorkingStack.Pointer-6] = m.WorkingStack.Data[m.WorkingStack.Pointer-12]
			} else {
				m.WorkingStack.Pointer += 3
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-4]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-5]
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = m.WorkingStack.Data[m.WorkingStack.Pointer-6]
			}
		}
		if shortMode {
			// First swap
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-3)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-4)
			// Second swap
			m.WorkingStack.Swap(m.WorkingStack.Pointer-3, m.WorkingStack.Pointer-5)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-4, m.WorkingStack.Pointer-6)
		} else {
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-2)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-3)
		}
	case 0x06: // DUP
		if shortMode {
			v := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 2
			m.WorkingStack.Push16(v)
			if keepMode {
				m.WorkingStack.Pointer += 2
				m.WorkingStack.Push16(v)
			}
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1]
			m.WorkingStack.Pointer++
			if keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1]
				m.WorkingStack.Pointer++
			}
		}
	case 0x07: // OVR
		if keepMode {
			if shortMode {
				v := m.WorkingStack.Peek16()
				m.WorkingStack.Pointer += 2
				m.WorkingStack.Push16(v)
			} else {
				m.WorkingStack.Pointer++
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-2]
			}
		}
		if shortMode {
			m.WorkingStack.Pointer += 2
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-5]
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-6]
		} else {
			m.WorkingStack.Pointer++
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-3]
		}
	// Logic instructions
	case 0x08: // EQU
		var res byte
		if shortMode &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] == m.WorkingStack.Data[m.WorkingStack.Pointer-3] &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] == m.WorkingStack.Data[m.WorkingStack.Pointer-4] {
			res = 0x01
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-1] == m.WorkingStack.Data[m.WorkingStack.Pointer-2] {
			res = 0x01
		}
		if keepMode {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = res
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] = res
			m.WorkingStack.Pointer--
		}
	case 0x09: // NEQ
		var res byte
		if shortMode &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] != m.WorkingStack.Data[m.WorkingStack.Pointer-3] &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] != m.WorkingStack.Data[m.WorkingStack.Pointer-4] {
			res = 0x01
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-1] != m.WorkingStack.Data[m.WorkingStack.Pointer-2] {
			res = 0x01
		}
		if keepMode {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = res
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] = res
			m.WorkingStack.Pointer--
		}
	case 0x0a: // GTH
		var res byte
		if shortMode &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] < m.WorkingStack.Data[m.WorkingStack.Pointer-3] &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] < m.WorkingStack.Data[m.WorkingStack.Pointer-4] {
			res = 0x01
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-1] < m.WorkingStack.Data[m.WorkingStack.Pointer-2] {
			res = 0x01
		}
		if keepMode {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = res
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] = res
			m.WorkingStack.Pointer--
		}
	case 0x0b: // LTH
		var res byte
		if shortMode &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] > m.WorkingStack.Data[m.WorkingStack.Pointer-3] &&
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] > m.WorkingStack.Data[m.WorkingStack.Pointer-4] {
			res = 0x01
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-1] > m.WorkingStack.Data[m.WorkingStack.Pointer-2] {
			res = 0x01
		}
		if keepMode {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = res
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer-2] = res
			m.WorkingStack.Pointer--
		}
	case 0x0c: //	JMP
		if shortMode {
			m.ProgramCounter = m.WorkingStack.Peek16()
		} else {
			m.ProgramCounter += uint16(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
		}
	case 0x0d: // JCN
		if shortMode && m.WorkingStack.Data[m.WorkingStack.Pointer-3] != 00 {
			m.ProgramCounter = m.WorkingStack.Peek16()
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-2] != 00 {
			m.ProgramCounter += uint16(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
		}
	case 0x0e: // JSR
		m.ReturnStack.Push16(m.ProgramCounter)
		m.ReturnStack.Pointer += 2
		if shortMode {
			m.ProgramCounter = m.WorkingStack.Peek16()
		} else {
			m.ProgramCounter += uint16(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
		}
	case 0x0f: // STH
		m.ReturnStack.Data[m.ReturnStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1]
		m.ReturnStack.Pointer++
	// Memory instructions
	case 0x10: // LDZ
		m.WorkingStack.Data[m.WorkingStack.Pointer] = m.Memory[m.WorkingStack.Data[m.WorkingStack.Pointer-1]]
		m.WorkingStack.Pointer++
	case 0x11: // STZ
		m.Memory[m.WorkingStack.Data[m.WorkingStack.Pointer-1]] = m.WorkingStack.Data[m.WorkingStack.Pointer-2]
	case 0x12: // LDR
		panic("UNIMPLEMENTED")
		//m.WorkingStack.Data[m.WorkingStack.Pointer] = m.Memory[m.ProgramCounter+m.WorkingStack.Data[m.WorkingStack.Pointer-1]]
		//m.WorkingStack.Pointer++
	case 0x13: // STR
		panic("UNIMPLEMENTED")
		//m.Memory[m.ProgramCounter + m.WorkingStack.Data[m.WorkingStack.Pointer-1]
	case 0x14: // LDA
		fmt.Println("Load aboslute")
		ind := m.WorkingStack.Peek16()
		if !keepMode {
			if shortMode {
				m.WorkingStack.Pointer += 2
			} else {
				m.WorkingStack.Pointer++
			}
		}
		if shortMode {
			m.WorkingStack.Push16(uint16(m.Memory[ind])<<8 + uint16(m.Memory[ind+1]))
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.Memory[ind]
		}
	case 0x15: // STA
		panic("UNIMPLEMENTED")
	case 0x16: // DEI
		panic("UNIMPLEMENTED")
	case 0x17: // DEO
		panic("UNIMPLEMENTED")
	// Arithmetic instructions
	case 0x18: // ADD
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a + b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] + m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x19: // SUB
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a - b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] - m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1a: // MUL
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a * b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] * m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1b: // DIV
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a / b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] / m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1c: // AND
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a & b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] & m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1d: // ORA
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a | b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] | m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1e: // EOR
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a ^ b)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] ^ m.WorkingStack.Data[m.WorkingStack.Pointer-2]
		}
	case 0x1f: // SFT
		if shortMode {
			a := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer -= 2
			b := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(a >> (b & 0xff) << (b * 0xff00))
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-2] >> (m.WorkingStack.Data[m.WorkingStack.Pointer-1] & 0xf) << (m.WorkingStack.Data[m.WorkingStack.Pointer-1] & 0xf0)
		}
	default:
		panic(fmt.Sprintf("Unknown opcode: HEX: 0x%x, BIN: %8b", op&0x1f, op))
	}
}

func Describe(op byte) {
	fmt.Printf("Opcode: %x\n", op)
	if op&0x01 != 0 {
		fmt.Println("Keep")
	}
	if op&0x02 != 0 {
		fmt.Println("Return")
	}
	if op&0x04 != 0 {
		fmt.Println("Short")
	}
}
