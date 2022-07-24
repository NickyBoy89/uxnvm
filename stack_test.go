package main

import (
	"testing"
)

// Tests the stack functions of the virtual machine

func CreateStack(bytes []byte) Stack {
	var s Stack
	for i := range bytes {
		s.Data[i] = bytes[i]
		s.Pointer = byte(len(bytes))
	}
	return s
}

func TestLIT(t *testing.T) {
	var m Machine
	m.Load([]byte{0x80, 0x12}) // LIT 12
	m.Execute()

	expected := CreateStack([]byte{0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestLIT2(t *testing.T) {
	var m Machine
	m.Load([]byte{0xa0, 0xab, 0xcd}) // LIT2 ab cd
	m.Execute()

	expected := CreateStack([]byte{0xab, 0xcd})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestINC(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01})
	m.Load([]byte{0x01}) // INC
	m.Execute()

	expected := CreateStack([]byte{0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestINCk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x01})
	m.Load([]byte{0x81}) // INCk
	m.Execute()

	expected := CreateStack([]byte{0x01, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestINC2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x01})
	m.Load([]byte{0x21}) // INC2
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestINC2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x00, 0x01})
	m.Load([]byte{0xa1}) // INC2k
	m.Execute()

	expected := CreateStack([]byte{0x00, 0x01, 0x00, 0x02})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestPOP(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x02}) // POP
	m.Execute()

	expected := CreateStack([]byte{0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestPOPk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x82}) // POPk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestPOP2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x22}) // POP2
	m.Execute()

	expected := CreateStack([]byte{})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestPOP2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x82}) // POP2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNIP(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x03}) // NIP
	m.Execute()

	expected := CreateStack([]byte{0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNIPk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x83}) // NIPk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNIP2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0x23}) // NIP2
	m.Execute()

	expected := CreateStack([]byte{0x56, 0x78})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestNIP2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0xa3}) // NIP2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x56, 0x78})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestSWP(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x04}) // SWP
	m.Execute()

	expected := CreateStack([]byte{0x34, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestSWPk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x84}) // SWPk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x34, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestSWP2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0x24}) // SWP2
	m.Execute()

	expected := CreateStack([]byte{0x56, 0x78, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestSWP2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0xa4}) // SWP2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x56, 0x78, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
	}
}

func TestROT(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56})
	m.Load([]byte{0x05}) // ROT
	m.Execute()

	expected := CreateStack([]byte{0x34, 0x56, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestROTk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56})
	m.Load([]byte{0x85}) // ROTk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x34, 0x56, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestROT2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc})
	m.Load([]byte{0x25}) // ROT2
	m.Execute()

	expected := CreateStack([]byte{0x56, 0x78, 0x9a, 0xbc, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestROT2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc})
	m.Load([]byte{0xa5}) // ROT2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0x56, 0x78, 0x9a, 0xbc, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestDUP(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x06}) // DUP
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestDUPk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x86}) // DUPk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x34, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestDUP2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x26}) // DUP2
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestDUP2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0xa6}) // DUP2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x12, 0x34, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestOVR(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x07}) // OVR
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestOVRk(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34})
	m.Load([]byte{0x87}) // OVRk
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x12, 0x34, 0x12})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestOVR2(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0x27}) // OVR2
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}

func TestOVR2k(t *testing.T) {
	var m Machine
	m.WorkingStack = CreateStack([]byte{0x12, 0x34, 0x56, 0x78})
	m.Load([]byte{0xa7}) // OVR2k
	m.Execute()

	expected := CreateStack([]byte{0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x56, 0x78, 0x12, 0x34})

	if m.WorkingStack.Data != expected.Data {
		t.Logf("Actual: %v", m.WorkingStack.Data)
		t.Logf("Expect: %v", expected.Data)
		t.Fatal("Stacks differed")
	}
}
