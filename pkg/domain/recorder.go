package domain

import (
	"fmt"
	"strings"
)

// Recorder ...
type Recorder struct {
	PC             uint16
	FetchedValue   []byte
	Mnemonic       Mnemonic
	AddressingMode AddressingMode
	Address        Address
	Data           byte
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
	case Immediate:
		operand = fmt.Sprintf("#$%02X", e.Data)
	case Absolute:
		operand = fmt.Sprintf("$%04X", e.Address)
	case ZeroPage:
		operand = fmt.Sprintf("$%02X = %02X", e.Address, e.Data)
	case IndexedZeroPageX:
	case IndexedZeroPageY:
	case IndexedAbsoluteX:
	case IndexedAbsoluteY:
	case Implied:
	case Relative:
		operand = fmt.Sprintf("$%04X", e.Address)
	case IndexedIndirect:
	case IndirectIndexed:
	case AbsoluteIndirect:
	default:
		operand = ""
	}

	s := operand + strings.Repeat(" ", 27)
	return s[0:27]
}

// String ...
func (e *Recorder) String() string {
	fetchedValue := e.makeFetchedValueString()
	operand := e.makeOperandString()

	return fmt.Sprintf("%04X  %v %v %v A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d", e.PC, fetchedValue, e.Mnemonic, operand, e.A, e.X, e.Y, e.P, e.SP, e.Dot, e.Scanline, e.Cycle)
}
