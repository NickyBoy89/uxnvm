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
	Data    [254]byte
	Error   byte
	Pointer byte
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

func (m *Machine) Load(rom []byte) {
	for offset := 0; offset < len(rom); offset++ {
		m.Memory[ProgramStartPage+offset] = rom[offset]
	}
	m.ProgramCounter = ProgramStartPage
}

func (m *Machine) Execute() {
	var op byte
	shortMode := op&0x20 != 0
	returnMode := op&0x40 != 0
	keepMode := op&0x80 != 0

	_, _, _ = keepMode, returnMode, shortMode

	// Mask the opcode to the last five bits to extract the opcode
	switch op & 0x1f {
	// Stack instructions
	case 0x00: // LIT
	case 0x01: // INC
	case 0x02: // POP
	case 0x03: // NIP
	case 0x04: // SWP
	case 0x05: // ROT
	case 0x06: // DUP
	case 0x07: // OVR
	// Logic instructions
	case 0x08: // EQU
	case 0x09: // NEQ
	case 0x0a: // GTH
	case 0x0b: // LTH
	case 0x0c: //	JMP
	case 0x0d: // JCN
	case 0x0e: // JSR
	case 0x0f: // STH
	// Memory instructions
	case 0x10: // LDZ
	case 0x11: // STZ
	case 0x12: // LDR
	case 0x13: // STR
	case 0x14: // LDA
	case 0x15: // STA
	case 0x16: // DEI
	case 0x17: // DEO
	// Arithmetic instructions
	case 0x18: // ADD
	case 0x19: // SUB
	case 0x1a: // MUL
	case 0x1b: // DIV
	case 0x1c: // ADD
	case 0x1d: // ORA
	case 0x1e: // EOR
	case 0x1f: // SFT
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
