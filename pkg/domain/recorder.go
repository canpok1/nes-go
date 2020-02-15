package domain

import (
	"fmt"
	"nes-go/pkg/log"
	"strings"
)

// Recorder ...
type Recorder struct {
	PC             uint16
	FetchedValue   []byte
	Mnemonic       Mnemonic
	Documented     bool
	AddressingMode AddressingMode
	Address        []Address
	Data           *byte
	A              byte
	X              byte
	Y              byte
	P              byte
	SP             byte
	Dot            uint16
	Scanline       uint16
	Cycle          int
}

// makeFetchedValueString ...
func (e *Recorder) makeFetchedValueString() string {
	switch len(e.FetchedValue) {
	case 1:
		return fmt.Sprintf("%02X       ", e.FetchedValue[0])
	case 2:
		return fmt.Sprintf("%02X %02X    ", e.FetchedValue[0], e.FetchedValue[1])
	case 3:
		return fmt.Sprintf("%02X %02X %02X ", e.FetchedValue[0], e.FetchedValue[1], e.FetchedValue[2])
	default:
		return "         "
	}
}

// makeOperandString ...
func (e *Recorder) makeOperandString() string {
	var operand string
	switch e.AddressingMode {
	case Accumulator:
		operand = fmt.Sprintf("A")
	case Immediate:
		if e.Data != nil {
			operand = fmt.Sprintf("#$%02X", *e.Data)
		}
	case Absolute:
		var a0 string
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("$%04X", e.Address[0])
		}
		if e.Data == nil {
			operand = fmt.Sprintf("%s", a0)
		} else {
			operand = fmt.Sprintf("%s = %02X", a0, *e.Data)
		}
	case ZeroPage:
		var a0 string
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("$%02X", e.Address[0])
		}

		if e.Data == nil {
			operand = fmt.Sprintf("%s", a0)
		} else {
			operand = fmt.Sprintf("%s = %02X", a0, *e.Data)
		}
	case IndexedZeroPageX:
		var f1, a0, d string
		if len(e.FetchedValue) >= 2 {
			f1 = fmt.Sprintf("%02X", e.FetchedValue[1])
		}
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%02X", e.Address[0])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("$%s,X @ %s = %s", f1, a0, d)
	case IndexedZeroPageY:
		var f1, a0, d string
		if len(e.FetchedValue) >= 2 {
			f1 = fmt.Sprintf("%02X", e.FetchedValue[1])
		}
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%02X", e.Address[0])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("$%s,Y @ %s = %s", f1, a0, d)
	case IndexedAbsoluteX:
		var a0, a1, d string
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%04X", e.Address[0])
		}
		if len(e.Address) >= 2 {
			a1 = fmt.Sprintf("%04X", e.Address[1])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("$%s,X @ %s = %s", a0, a1, d)
	case IndexedAbsoluteY:
		var a0, a1, d string
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%04X", e.Address[0])
		}
		if len(e.Address) >= 2 {
			a1 = fmt.Sprintf("%04X", e.Address[1])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("$%s,Y @ %s = %s", a0, a1, d)
	case Implied:
	case Relative:
		operand = fmt.Sprintf("$%04X", e.Address[0])
	case IndexedIndirect:
		var f1, a0, a1, d string
		if len(e.FetchedValue) >= 2 {
			f1 = fmt.Sprintf("%02X", e.FetchedValue[1])
		}
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%02X", e.Address[0])
		}
		if len(e.Address) >= 2 {
			a1 = fmt.Sprintf("%04X", e.Address[1])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("($%s,X) @ %s = %s = %s", f1, a0, a1, d)
	case IndirectIndexed:
		var f1, a0, a1, d string
		if len(e.FetchedValue) >= 2 {
			f1 = fmt.Sprintf("%02X", e.FetchedValue[1])
		}
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%04X", e.Address[0])
		}
		if len(e.Address) >= 2 {
			a1 = fmt.Sprintf("%04X", e.Address[1])
		}
		if e.Data != nil {
			d = fmt.Sprintf("%02X", *e.Data)
		}
		operand = fmt.Sprintf("($%s),Y = %s @ %s = %s", f1, a0, a1, d)
	case AbsoluteIndirect:
		var a0, a1 string
		if len(e.Address) >= 1 {
			a0 = fmt.Sprintf("%04X", e.Address[0])
		}
		if len(e.Address) >= 2 {
			a1 = fmt.Sprintf("%04X", e.Address[1])
		}
		operand = fmt.Sprintf("($%s) = %s", a0, a1)
	default:
		operand = ""
	}

	s := operand + strings.Repeat(" ", 27)
	return s[0:27]
}

// String ...
func (e *Recorder) String() string {
	f := e.makeFetchedValueString()
	o := e.makeOperandString()

	m := " "
	if !e.Documented {
		m = "*"
	}

	return fmt.Sprintf("%04X  %v%v%v %v A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d", e.PC, f, m, e.Mnemonic, o, e.A, e.X, e.Y, e.P, e.SP, e.Dot, e.Scanline, e.Cycle)
}

// AddAddress ...
func (e *Recorder) AddAddress(addr Address) {
	log.Trace("Recorder.AddAddress[addr=%#v]", addr)
	e.Address = append(e.Address, addr)
}
