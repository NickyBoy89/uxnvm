package main

type UxnError byte

func (ue UxnError) Error() string {
	switch ue {
	case ErrUnknown:
		return "Error: Unknown"
	case ErrUnderflow:
		return "Error: Underflow"
	case ErrOverflow:
		return "Error: Overflow"
	case ErrDivByZero:
		return "Error: Divide by zero"
	}
	panic("Unknown error")
}

const (
	ErrUnknown UxnError = iota
	ErrUnderflow
	ErrOverflow
	ErrDivByZero
)

type Stack struct {
	Data    [254]byte
	Error   UxnError
	Pointer byte
}

func (s Stack) String() string {
	return HexPrint(s.Data[:s.Pointer])
}

func (s *Stack) Push8(x byte) {
	if s.Pointer == 0xff {
		s.Error = ErrOverflow
		panic("error")
	}
	s.Data[s.Pointer] = x
	s.Pointer++
}

func (s *Stack) Push16(x uint16) {
	j := s.Pointer
	if j >= 0xfe {
		s.Error = ErrOverflow
		panic(s.Error)
	}
	k := x
	s.Data[j] = byte(k >> 8)
	s.Data[j+1] = byte(k)
	s.Pointer = j + 2
}

func (s *Stack) Pop8(srcStackPtr *byte) byte {
	if *srcStackPtr == 0 {
		s.Error = ErrUnderflow
		panic(s.Error)
	}
	*srcStackPtr--
	return s.Data[s.Pointer]
}

func (s *Stack) Pop16(srcStackPtr *byte) uint16 {
	if *srcStackPtr <= 1 {
		s.Error = ErrUnderflow
		panic(s.Error)
	}
	o := uint16(s.Data[*srcStackPtr-1]) + uint16(s.Data[*srcStackPtr-2])<<8
	*srcStackPtr -= 2
	return o
}
