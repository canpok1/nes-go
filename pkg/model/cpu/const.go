package cpu

import (
	"fmt"
	"nes-go/pkg/model"
)

// Opcode ...
type Opcode uint8

const (
	// ErrorOpcode ...
	ErrorOpcode Opcode = 0
)

// Operand ...
type Operand struct {
	Data    *byte
	Address *model.Address
}

// String ...
func (o *Operand) String() string {
	d := ""
	if o.Data != nil {
		d = fmt.Sprintf("%#v", *o.Data)
	}

	a := ""
	if o.Address != nil {
		a = fmt.Sprintf("%#v", *o.Address)
	}

	return fmt.Sprintf("{Data:%#v, Address:%#v}", d, a)
}

// Mnemonic ...
type Mnemonic string

const (
	// ADC ...
	ADC Mnemonic = "ADC"
	// SBC ...
	SBC Mnemonic = "SBC"
	// AND ...
	AND Mnemonic = "AND"
	// ORA ...
	ORA Mnemonic = "ORA"
	// EOR ...
	EOR Mnemonic = "EOR"
	// ASL ...
	ASL Mnemonic = "ASL"
	// LSR ...
	LSR Mnemonic = "LSR"
	// ROL ...
	ROL Mnemonic = "ROL"
	// ROR ...
	ROR Mnemonic = "ROR"
	// BCC ...
	BCC Mnemonic = "BCC"
	// BCS ...
	BCS Mnemonic = "BCS"
	// BEQ ...
	BEQ Mnemonic = "BEQ"
	// BNE ...
	BNE Mnemonic = "BNE"
	// BVC ...
	BVC Mnemonic = "BVC"
	// BVS ...
	BVS Mnemonic = "BVS"
	// BPL ...
	BPL Mnemonic = "BPL"
	// BMI ...
	BMI Mnemonic = "BMI"
	// BIT ...
	BIT Mnemonic = "BIT"
	// JMP ...
	JMP Mnemonic = "JMP"
	// JSR ...
	JSR Mnemonic = "JSR"
	// RTS ...
	RTS Mnemonic = "RTS"
	// BRK ...
	BRK Mnemonic = "BRK"
	// RTI ...
	RTI Mnemonic = "RTI"
	// CMP ...
	CMP Mnemonic = "CMP"
	// CPX ...
	CPX Mnemonic = "CPX"
	// CPY ...
	CPY Mnemonic = "CPY"
	// INC ...
	INC Mnemonic = "INC"
	// DEC ...
	DEC Mnemonic = "DEC"
	// INX ...
	INX Mnemonic = "INX"
	// DEX ...
	DEX Mnemonic = "DEX"
	// INY ...
	INY Mnemonic = "INY"
	// DEY ...
	DEY Mnemonic = "DEY"
	// CLC ...
	CLC Mnemonic = "CLC"
	// SEC ...
	SEC Mnemonic = "SEC"
	// CLI ...
	CLI Mnemonic = "CLI"
	// SEI ...
	SEI Mnemonic = "SEI"
	// CLD ...
	CLD Mnemonic = "CLD"
	// SED ...
	SED Mnemonic = "SED"
	// CLV ...
	CLV Mnemonic = "CLV"
	// LDA ...
	LDA Mnemonic = "LDA"
	// LDX ...
	LDX Mnemonic = "LDX"
	// LDY ...
	LDY Mnemonic = "LDY"
	// STA ...
	STA Mnemonic = "STA"
	// STX ...
	STX Mnemonic = "STX"
	// STY ...
	STY Mnemonic = "STY"
	// TAX ...
	TAX Mnemonic = "TAX"
	// TXA ...
	TXA Mnemonic = "TXA"
	// TAY ...
	TAY Mnemonic = "TAY"
	// TYA ...
	TYA Mnemonic = "TYA"
	// TSX ...
	TSX Mnemonic = "TSX"
	// TXS ...
	TXS Mnemonic = "TXS"
	// PHA ...
	PHA Mnemonic = "PHA"
	// PLA ...
	PLA Mnemonic = "PLA"
	// PHP ...
	PHP Mnemonic = "PHP"
	// PLP ...
	PLP Mnemonic = "PLP"
	// NOP ...
	NOP Mnemonic = "NOP"
)

// AddressingMode ...
// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm
type AddressingMode string

const (
	// Accumulator ...
	Accumulator AddressingMode = "Accumulator"
	// Immediate ...
	Immediate AddressingMode = "Immediate"
	// Absolute ...
	Absolute AddressingMode = "Absolute"
	// ZeroPage ...
	ZeroPage AddressingMode = "ZeroPage"
	// IndexedZeroPageX ...
	IndexedZeroPageX AddressingMode = "IndexedZeroPageX"
	// IndexedZeroPageY ...
	IndexedZeroPageY AddressingMode = "IndexedZeroPageY"
	// IndexedAbsoluteX ...
	IndexedAbsoluteX AddressingMode = "IndexedAbsoluteX"
	// IndexedAbsoluteY ...
	IndexedAbsoluteY AddressingMode = "IndexedAbsoluteY"
	// Implied ...
	Implied AddressingMode = "Implied"
	// Relative ...
	Relative AddressingMode = "Relative"
	// IndexedIndirect ...
	IndexedIndirect AddressingMode = "IndexedIndirect"
	// IndirectIndexed ...
	IndirectIndexed AddressingMode = "IndirectIndexed"
	// AbsoluteIndirect ...
	AbsoluteIndirect AddressingMode = "AbsoluteIndirect"
)

// OpcodeProp ...
type OpcodeProp struct {
	Mnemonic       Mnemonic
	AddressingMode AddressingMode
	Cycle          int
}

// OpcodeProps ...
var OpcodeProps = map[Opcode]OpcodeProp{
	0x69: OpcodeProp{ADC, Immediate, 2},
	0x65: OpcodeProp{ADC, ZeroPage, 3},
	0x75: OpcodeProp{ADC, IndexedZeroPageX, 4},
	0x6D: OpcodeProp{ADC, Absolute, 4},
	0x7D: OpcodeProp{ADC, IndexedAbsoluteX, 4},
	0x79: OpcodeProp{ADC, IndexedAbsoluteY, 4},
	0x61: OpcodeProp{ADC, IndexedIndirect, 6},
	0x71: OpcodeProp{ADC, IndirectIndexed, 5},
	0xE9: OpcodeProp{SBC, Immediate, 2},
	0xE5: OpcodeProp{SBC, ZeroPage, 3},
	0xF5: OpcodeProp{SBC, IndexedZeroPageX, 4},
	0xED: OpcodeProp{SBC, Absolute, 4},
	0xFD: OpcodeProp{SBC, IndexedAbsoluteX, 4},
	0xF9: OpcodeProp{SBC, IndexedAbsoluteY, 4},
	0xE1: OpcodeProp{SBC, IndexedIndirect, 6},
	0xF1: OpcodeProp{SBC, IndirectIndexed, 5},
	0x29: OpcodeProp{AND, Immediate, 2},
	0x25: OpcodeProp{AND, ZeroPage, 3},
	0x35: OpcodeProp{AND, IndexedZeroPageX, 4},
	0x2D: OpcodeProp{AND, Absolute, 4},
	0x3D: OpcodeProp{AND, IndexedAbsoluteX, 4},
	0x39: OpcodeProp{AND, IndexedAbsoluteY, 4},
	0x21: OpcodeProp{AND, IndexedIndirect, 6},
	0x31: OpcodeProp{AND, IndirectIndexed, 5},
	0x09: OpcodeProp{ORA, Immediate, 2},
	0x05: OpcodeProp{ORA, ZeroPage, 3},
	0x15: OpcodeProp{ORA, IndexedZeroPageX, 4},
	0x0D: OpcodeProp{ORA, Absolute, 4},
	0x1D: OpcodeProp{ORA, IndexedAbsoluteX, 4},
	0x19: OpcodeProp{ORA, IndexedAbsoluteY, 4},
	0x01: OpcodeProp{ORA, IndexedIndirect, 6},
	0x11: OpcodeProp{ORA, IndirectIndexed, 5},
	0x49: OpcodeProp{EOR, Immediate, 2},
	0x45: OpcodeProp{EOR, ZeroPage, 3},
	0x55: OpcodeProp{EOR, IndexedZeroPageX, 4},
	0x4D: OpcodeProp{EOR, Absolute, 4},
	0x5D: OpcodeProp{EOR, IndexedAbsoluteX, 4},
	0x59: OpcodeProp{EOR, IndexedAbsoluteY, 4},
	0x41: OpcodeProp{EOR, IndexedIndirect, 6},
	0x51: OpcodeProp{EOR, IndirectIndexed, 5},
	0x0A: OpcodeProp{ASL, Accumulator, 2},
	0x06: OpcodeProp{ASL, ZeroPage, 5},
	0x16: OpcodeProp{ASL, IndexedZeroPageX, 6},
	0x0E: OpcodeProp{ASL, Absolute, 6},
	0x1E: OpcodeProp{ASL, IndexedAbsoluteX, 6},
	0x4A: OpcodeProp{LSR, Accumulator, 2},
	0x46: OpcodeProp{LSR, ZeroPage, 5},
	0x56: OpcodeProp{LSR, IndexedZeroPageX, 6},
	0x4E: OpcodeProp{LSR, Absolute, 6},
	0x5E: OpcodeProp{LSR, IndexedAbsoluteX, 6},
	0x2A: OpcodeProp{ROL, Accumulator, 2},
	0x26: OpcodeProp{ROL, ZeroPage, 5},
	0x36: OpcodeProp{ROL, IndexedZeroPageX, 6},
	0x2E: OpcodeProp{ROL, Absolute, 6},
	0x3E: OpcodeProp{ROL, IndexedAbsoluteX, 6},
	0x6A: OpcodeProp{ROR, Accumulator, 2},
	0x66: OpcodeProp{ROR, ZeroPage, 5},
	0x76: OpcodeProp{ROR, IndexedZeroPageX, 6},
	0x6E: OpcodeProp{ROR, Absolute, 6},
	0x7E: OpcodeProp{ROR, IndexedAbsoluteX, 6},
	0x90: OpcodeProp{BCC, Relative, 2},
	0xB0: OpcodeProp{BCS, Relative, 2},
	0xF0: OpcodeProp{BEQ, Relative, 2},
	0xD0: OpcodeProp{BNE, Relative, 2},
	0x50: OpcodeProp{BVC, Relative, 2},
	0x70: OpcodeProp{BVS, Relative, 2},
	0x10: OpcodeProp{BPL, Relative, 2},
	0x30: OpcodeProp{BMI, Relative, 2},
	0x24: OpcodeProp{BIT, ZeroPage, 3},
	0x2C: OpcodeProp{BIT, Absolute, 4},
	0x4C: OpcodeProp{JMP, Absolute, 3},
	0x6C: OpcodeProp{JMP, AbsoluteIndirect, 5},
	0x20: OpcodeProp{JSR, Absolute, 6},
	0x60: OpcodeProp{RTS, Implied, 6},
	0x00: OpcodeProp{BRK, Implied, 7},
	0x40: OpcodeProp{RTI, Implied, 6},
	0xC9: OpcodeProp{CMP, Immediate, 2},
	0xC5: OpcodeProp{CMP, ZeroPage, 3},
	0xD5: OpcodeProp{CMP, IndexedZeroPageX, 4},
	0xCD: OpcodeProp{CMP, Absolute, 4},
	0xDD: OpcodeProp{CMP, IndexedAbsoluteX, 4},
	0xD9: OpcodeProp{CMP, IndexedAbsoluteY, 4},
	0xC1: OpcodeProp{CMP, IndexedIndirect, 6},
	0xD1: OpcodeProp{CMP, IndirectIndexed, 5},
	0xE0: OpcodeProp{CPX, Immediate, 2},
	0xE4: OpcodeProp{CPX, ZeroPage, 3},
	0xEC: OpcodeProp{CPX, Absolute, 4},
	0xC0: OpcodeProp{CPY, Immediate, 2},
	0xC4: OpcodeProp{CPY, ZeroPage, 3},
	0xCC: OpcodeProp{CPY, Absolute, 4},
	0xE6: OpcodeProp{INC, ZeroPage, 5},
	0xF6: OpcodeProp{INC, IndexedZeroPageX, 6},
	0xEE: OpcodeProp{INC, Absolute, 6},
	0xFE: OpcodeProp{INC, IndexedAbsoluteX, 6},
	0xC6: OpcodeProp{DEC, ZeroPage, 5},
	0xD6: OpcodeProp{DEC, IndexedZeroPageX, 6},
	0xCE: OpcodeProp{DEC, Absolute, 6},
	0xDE: OpcodeProp{DEC, IndexedAbsoluteX, 6},
	0xE8: OpcodeProp{INX, Implied, 2},
	0xCA: OpcodeProp{DEX, Implied, 2},
	0xC8: OpcodeProp{INY, Implied, 2},
	0x88: OpcodeProp{DEY, Implied, 2},
	0x18: OpcodeProp{CLC, Implied, 2},
	0x38: OpcodeProp{SEC, Implied, 2},
	0x58: OpcodeProp{CLI, Implied, 2},
	0x78: OpcodeProp{SEI, Implied, 2},
	0xD8: OpcodeProp{CLD, Implied, 2},
	0xF8: OpcodeProp{SED, Implied, 2},
	0xB8: OpcodeProp{CLV, Implied, 2},
	0xA9: OpcodeProp{LDA, Immediate, 2},
	0xA5: OpcodeProp{LDA, ZeroPage, 3},
	0xB5: OpcodeProp{LDA, IndexedZeroPageX, 4},
	0xAD: OpcodeProp{LDA, Absolute, 4},
	0xBD: OpcodeProp{LDA, IndexedAbsoluteX, 4},
	0xB9: OpcodeProp{LDA, IndexedAbsoluteY, 4},
	0xA1: OpcodeProp{LDA, IndexedIndirect, 6},
	0xB1: OpcodeProp{LDA, IndirectIndexed, 5},
	0xA2: OpcodeProp{LDX, Immediate, 2},
	0xA6: OpcodeProp{LDX, ZeroPage, 3},
	0xB6: OpcodeProp{LDX, IndexedZeroPageY, 4},
	0xAE: OpcodeProp{LDX, Absolute, 4},
	0xBE: OpcodeProp{LDX, IndexedAbsoluteY, 4},
	0xA0: OpcodeProp{LDY, Immediate, 2},
	0xA4: OpcodeProp{LDY, ZeroPage, 3},
	0xB4: OpcodeProp{LDY, IndexedZeroPageX, 4},
	0xAC: OpcodeProp{LDY, Absolute, 4},
	0xBC: OpcodeProp{LDY, IndexedAbsoluteX, 4},
	0x85: OpcodeProp{STA, ZeroPage, 3},
	0x95: OpcodeProp{STA, IndexedZeroPageX, 4},
	0x8D: OpcodeProp{STA, Absolute, 4},
	0x9D: OpcodeProp{STA, IndexedAbsoluteX, 4},
	0x99: OpcodeProp{STA, IndexedAbsoluteY, 4},
	0x81: OpcodeProp{STA, IndexedIndirect, 6},
	0x91: OpcodeProp{STA, IndirectIndexed, 5},
	0x86: OpcodeProp{STX, ZeroPage, 3},
	0x96: OpcodeProp{STX, IndexedZeroPageY, 4},
	0x8E: OpcodeProp{STX, Absolute, 4},
	0x84: OpcodeProp{STY, ZeroPage, 3},
	0x94: OpcodeProp{STY, IndexedZeroPageX, 4},
	0x8C: OpcodeProp{STY, Absolute, 4},
	0xAA: OpcodeProp{TAX, Implied, 2},
	0x8A: OpcodeProp{TXA, Implied, 2},
	0xA8: OpcodeProp{TAY, Implied, 2},
	0x98: OpcodeProp{TYA, Implied, 2},
	0x9A: OpcodeProp{TXS, Implied, 2},
	0xBA: OpcodeProp{TSX, Implied, 2},
	0x48: OpcodeProp{PHA, Implied, 3},
	0x68: OpcodeProp{PLA, Implied, 4},
	0x08: OpcodeProp{PHP, Implied, 3},
	0x28: OpcodeProp{PLP, Implied, 4},
	0xEA: OpcodeProp{NOP, Implied, 2},
}
