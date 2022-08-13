package main

import (
	"errors"
	"fmt"
)

var StackOverflowError = errors.New("Stack overflow!")
var StackUnderflowError = errors.New("Stack underflow!")

// Page of memory where the program starts executing
const ProgramStartPage uint16 = 0x100

type Stack struct {
	Data    [254]byte
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
	return uint16(s.Data[s.Pointer-2])<<8 + uint16(s.Data[s.Pointer-1])
}

func (s *Stack) Push16(v uint16) {
	s.Data[s.Pointer-2] = byte(v >> 8)
	s.Data[s.Pointer-1] = byte(v & 0xff)
}

// Swap swaps the two elements at index n1 and n2 respectively
func (s *Stack) Swap(n1, n2 byte) {
	s.Data[n1], s.Data[n2] = s.Data[n2], s.Data[n1]
}

// ZeroFrom zeroes the section of stack from the start index to the end, inclusive
func (s *Stack) ZeroFrom(start, end byte) {
	for i := start; i < end+1; i++ {
		s.Data[i] = 0x00
	}
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
	return uint16(m.Memory[m.ProgramCounter])<<8 + uint16(m.Memory[m.ProgramCounter+1])
}

func (m *Machine) Load(rom []byte) {
	for offset := 0; offset < len(rom); offset++ {
		m.Memory[ProgramStartPage+uint16(offset)] = rom[offset]
	}
	m.ProgramCounter = ProgramStartPage
}

func (m *Machine) Execute() {
	op := m.Memory[m.ProgramCounter]
	m.ProgramCounter++

	//fmt.Printf("Executing instr: %.2x\n", op)

	shortMode := op&0x20 != 0
	shortModeInt := op & 0x20 >> 5
	returnMode := op&0x40 != 0
	keepMode := op&0x80 != 0

	_, _, _ = keepMode, returnMode, shortMode

	// Mask the opcode to the last five bits to extract the opcode
	switch op & 0x1f {
	// Stack instructions
	case 0x00: // LIT
		if shortMode {
			v := m.PeekMem16()
			m.ProgramCounter += 2
			m.WorkingStack.Pointer += 2
			m.WorkingStack.Push16(v)
			m.ProgramCounter += 1
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
			if keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-1] + 1
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1]++
			}
		}
	case 0x02: // POP
		if !keepMode {
			// NOTE: This zeroing might not be necessary
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
			if shortMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
			}
			m.WorkingStack.Pointer -= shortModeInt + 1
		}
	case 0x03: // NIP
		if shortMode {
			v := m.WorkingStack.Peek16()
			if keepMode {
				m.WorkingStack.Pointer += 2
				m.WorkingStack.Push16(v)
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = m.WorkingStack.Data[m.WorkingStack.Pointer-2]
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = m.WorkingStack.Data[m.WorkingStack.Pointer-1]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Pointer -= 2
			}
		} else {
			v := m.WorkingStack.Data[m.WorkingStack.Pointer-1]
			if keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer] = v
				m.WorkingStack.Pointer++
			} else {
				m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-2)
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Pointer--
			}
		}
	case 0x04: // SWP
		if keepMode {
			if shortMode {
				m.WorkingStack.Pointer += 4
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-3]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-4]
			} else {
				m.WorkingStack.Pointer += 2
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-3]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-4]
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
			// Second swap
			m.WorkingStack.Swap(m.WorkingStack.Pointer-3, m.WorkingStack.Pointer-5)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-4, m.WorkingStack.Pointer-6)
			// First swap
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-3)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-4)
		} else {
			m.WorkingStack.Swap(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-3)
			m.WorkingStack.Swap(m.WorkingStack.Pointer-1, m.WorkingStack.Pointer-2)
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
				m.WorkingStack.Pointer += 4
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-5]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-6]
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = m.WorkingStack.Data[m.WorkingStack.Pointer-7]
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = m.WorkingStack.Data[m.WorkingStack.Pointer-8]
			} else {
				m.WorkingStack.Pointer += 2
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = m.WorkingStack.Data[m.WorkingStack.Pointer-3]
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = m.WorkingStack.Data[m.WorkingStack.Pointer-4]
			}
		}

		if shortMode {
			m.WorkingStack.Pointer -= 2
			v := m.WorkingStack.Peek16()
			m.WorkingStack.Pointer += 4
			m.WorkingStack.Push16(v)
		} else {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = m.WorkingStack.Data[m.WorkingStack.Pointer-2]
			m.WorkingStack.Pointer++
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
			m.WorkingStack.Pointer++
		} else {
			if shortMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = 0x00
				m.WorkingStack.Pointer -= 3
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Pointer--
			}
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = res
		}
	case 0x09: // NEQ
		var res byte
		if shortMode {
			if m.WorkingStack.Data[m.WorkingStack.Pointer-1] != m.WorkingStack.Data[m.WorkingStack.Pointer-3] ||
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] != m.WorkingStack.Data[m.WorkingStack.Pointer-4] {
				res = 0x01
			}
		} else if m.WorkingStack.Data[m.WorkingStack.Pointer-1] != m.WorkingStack.Data[m.WorkingStack.Pointer-2] {
			res = 0x01
		}
		if keepMode {
			m.WorkingStack.Data[m.WorkingStack.Pointer] = res
			m.WorkingStack.Pointer++
		} else {
			if shortMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = 0x00
				m.WorkingStack.Pointer -= 3
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Pointer--
			}
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = res
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
			if shortMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = 0x00
				m.WorkingStack.Pointer -= 3
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Pointer--
			}
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = res
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
			if shortMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-3] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-4] = 0x00
				m.WorkingStack.Pointer -= 3
			} else {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Data[m.WorkingStack.Pointer-2] = 0x00
				m.WorkingStack.Pointer--
			}
			m.WorkingStack.Data[m.WorkingStack.Pointer-1] = res
		}
	case 0x0c: //	JMP
		if shortMode {
			m.ProgramCounter = m.WorkingStack.Peek16()
			if !keepMode {
				m.WorkingStack.ZeroFrom(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-1)
				m.WorkingStack.Pointer -= 2
			}
		} else {
			s := int8(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
			if s < 0 {
				m.ProgramCounter -= uint16(-s)
			} else {
				m.ProgramCounter += uint16(s)
			}
			if !keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Pointer--
			}
		}
		//m.ProgramCounter--
	case 0x0d: // JCN
		if shortMode {
			if m.WorkingStack.Data[m.WorkingStack.Pointer-3] != 0x00 {
				m.ProgramCounter = m.WorkingStack.Peek16()
			}
			if !keepMode {
				m.WorkingStack.ZeroFrom(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-1)
				m.WorkingStack.Pointer -= 2
			}
		} else {
			if m.WorkingStack.Data[m.WorkingStack.Pointer-2] != 0x00 {
				s := int8(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
				if s < 0 {
					m.ProgramCounter -= uint16(-s)
				} else {
					m.ProgramCounter += uint16(s)
				}
			}
			if !keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Pointer--
			}
		}
	case 0x0e: // JSR
		fmt.Printf("%.4x\n", m.ProgramCounter)
		m.ReturnStack.Pointer += 2
		m.ReturnStack.Push16(m.ProgramCounter)

		if shortMode {
			m.ProgramCounter = m.WorkingStack.Peek16()
			if !keepMode {
				m.WorkingStack.ZeroFrom(m.WorkingStack.Pointer-2, m.WorkingStack.Pointer-1)
				m.WorkingStack.Pointer -= 2
			}
		} else {
			s := int8(m.WorkingStack.Data[m.WorkingStack.Pointer-1])
			if s < 0 {
				m.ProgramCounter -= uint16(-s)
			} else {
				m.ProgramCounter += uint16(s)
			}
			if !keepMode {
				m.WorkingStack.Data[m.WorkingStack.Pointer-1] = 0x00
				m.WorkingStack.Pointer--
			}
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
