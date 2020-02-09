package domain

import (
	"fmt"
)

// Recorder ...
type Recorder struct {
	PC             uint16
	FetchedValue   []byte
	Mnemonic       Mnemonic
	AddressingMode AddressingMode
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
	var v1, v2 byte

	if len(e.FetchedValue) >= 2 {
		v1 = e.FetchedValue[1]
	}
	if len(e.FetchedValue) >= 3 {
		v2 = e.FetchedValue[2]
	}

	switch e.AddressingMode {
	case Absolute:
		return fmt.Sprintf("$%02X%02X                      ", v2, v1)
	default:
		return ""
	}
}

// String ...
func (e *Recorder) String() string {
	fetchedValue := e.makeFetchedValueString()
	operand := e.makeOperandString()

	return fmt.Sprintf("%04X  %v %v %v A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d", e.PC, fetchedValue, e.Mnemonic, operand, e.A, e.X, e.Y, e.P, e.SP, e.Dot, e.Scanline, e.Cycle)
}
