package main

import (
	"testing"
)

// Tests the logic functions of the virtual machine

func TestEQU(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x12})
	m.Load([]byte{0x08}) // EQU
	m.Execute()

	expected := CreateStack([]byte{0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestEQUk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x88}) // EQUk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestEQU2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0xab, 0xcd, 0xef, 0x01})
	m.Load([]byte{0x28}) // EQU2
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestEQU2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0xab, 0xcd, 0xab, 0xcd})
	m.Load([]byte{0xa8}) // EQU2k
	m.Execute()

	expected := CreateStack([]byte{0xab, 0xcd, 0xab, 0xcd, 0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNEQ(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x12})
	m.Load([]byte{0x09}) // NEQ
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNEQk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x89}) // NEQk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNEQ2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0xab, 0xcd, 0xef, 0x01})
	m.Load([]byte{0x29}) // NEQ2
	m.Execute()

	expected := CreateStack([]byte{0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNEQ2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0xab, 0xcd, 0xab, 0xcd})
	m.Load([]byte{0xa9}) // NEQ2k
	m.Execute()

	expected := CreateStack([]byte{0xab, 0xcd, 0xab, 0xcd, 0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestGTH(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x0a}) // GTH
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestGTHk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x34, 0x12})
	m.Load([]byte{0x8a}) // GTHk
	m.Execute()

	expected := CreateStack([]byte{0x34, 0x12, 0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestGTH2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x34, 0x56, 0x12, 0x34})
	m.Load([]byte{0x2a}) // GTH2
	m.Execute()

	expected := CreateStack([]byte{0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestGTH2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x34, 0x56})
	m.Load([]byte{0xaa}) // GTH2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x34, 0x56, 0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestLTH(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01, 0x01})
	m.Load([]byte{0x0b}) // LTH
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestLTHk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01, 0x00})
	m.Load([]byte{0x8b}) // LTHk
	m.Execute()

	expected := CreateStack([]byte{0x01, 0x00, 0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestLTH2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x01, 0x00, 0x00})
	m.Load([]byte{0x2b}) // LTH2
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestLTH2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x01, 0x00, 0x00})
	m.Load([]byte{0xab}) // LTH2k
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x01, 0x00, 0x00, 0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestJMP(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x02})
	m.Load([]byte{0x0c}) // JMP
	m.Execute()

	expected := CreateStack([]byte{})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage+0x02 {
		t.Logf("Actual ProgramCounter: 0x%.4x", m.ProgramCounter)
		t.Logf("Expect ProgramCounter: 0x%.4x", ProgramStartPage+0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJMPk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x02})
	m.Load([]byte{0x8c}) // JMPk
	m.Execute()

	expected := CreateStack([]byte{0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage+0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage+0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJMP2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x02})
	m.Load([]byte{0x2c}) // JMP2
	m.Execute()

	expected := CreateStack([]byte{})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != 0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", 0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJMP2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x02})
	m.Load([]byte{0xac}) // JMP2k
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != 0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", 0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJCN(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x02})
	m.Load([]byte{0x0d}) // JCN
	m.Execute()

	expected := CreateStack([]byte{0x00})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage)
		t.Fatal("Program counters differed")
	}
}

func TestJCNk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01, 0x02})
	m.Load([]byte{0x8d}) // JCNk
	m.Execute()

	expected := CreateStack([]byte{0x01, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage+0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage+0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJCN2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01, 0x00, 0x02})
	m.Load([]byte{0x2d}) // JCN2
	m.Execute()

	expected := CreateStack([]byte{0x01})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != 0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", 0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJCN2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x00, 0x02})
	m.Load([]byte{0xad}) // JCN2k
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x00, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage)
		t.Fatal("Program counters differed")
	}
}

func TestJSR(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x02})
	m.Load([]byte{0x0e}) // JSR
	m.Execute()

	//fmt.Printf("%.4x\n", m.ProgramCounter)

	expected := CreateStack([]byte{})
	expectedReturn := CreateStack([]byte{byte(ProgramStartPage >> 8), byte(ProgramStartPage&0xff) + 1})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ReturnStack.Data != expectedReturn.Data {
		t.Logf("Actual: %v", m.ReturnStack.Data)
		t.Logf("Expect: %v", expectedReturn.Data)
		t.Fatal("Return stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage+0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage+0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJSRk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x02})
	m.Load([]byte{0x8e}) // JSRk
	m.Execute()

	expected := CreateStack([]byte{0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != ProgramStartPage+0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", ProgramStartPage+0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJSR2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x02})
	m.Load([]byte{0x2e}) // JSR2
	m.Execute()

	expected := CreateStack([]byte{})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != 0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", 0x02)
		t.Fatal("Program counters differed")
	}
}

func TestJSR2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x02})
	m.Load([]byte{0xae}) // JSR2k
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}

	if m.ProgramCounter != 0x02 {
		t.Logf("Expect ProgramCounter: %v", m.ProgramCounter)
		t.Logf("Actual ProgramCounter: %v", 0x02)
		t.Fatal("Program counters differed")
	}
}
