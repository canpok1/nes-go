package model

// Opcode ...
type Opcode uint8

const (
	ErrorOpcode Opcode = 0
)

// Address ...
type Address uint16

// Mnemonic ...
type Mnemonic string

const (
	ADC Mnemonic = "ADC"
	SBC Mnemonic = "SBC"
	AND Mnemonic = "AND"
	ORA Mnemonic = "ORA"
	EOR Mnemonic = "EOR"
	ASL Mnemonic = "ASL"
	LSR Mnemonic = "LSR"
	ROL Mnemonic = "ROL"
	ROR Mnemonic = "ROR"
	BCC Mnemonic = "BCC"
	BCS Mnemonic = "BCS"
	BEQ Mnemonic = "BEQ"
	BNE Mnemonic = "BNE"
	BVC Mnemonic = "BVC"
	BVS Mnemonic = "BVS"
	BPL Mnemonic = "BPL"
	BMI Mnemonic = "BMI"
	BIT Mnemonic = "BIT"
	JMP Mnemonic = "JMP"
	JSR Mnemonic = "JSR"
	RTS Mnemonic = "RTS"
	BRK Mnemonic = "BRK"
	RTI Mnemonic = "RTI"
	CMP Mnemonic = "CMP"
	CPX Mnemonic = "CPX"
	CPY Mnemonic = "CPY"
	INC Mnemonic = "INC"
	DEC Mnemonic = "DEC"
	INX Mnemonic = "INX"
	DEX Mnemonic = "DEX"
	INY Mnemonic = "INY"
	DEY Mnemonic = "DEY"
	CLC Mnemonic = "CLC"
	SEC Mnemonic = "SEC"
	CLI Mnemonic = "CLI"
	SEI Mnemonic = "SEI"
	CLD Mnemonic = "CLD"
	SED Mnemonic = "SED"
	CLV Mnemonic = "CLV"
	LDA Mnemonic = "LDA"
	LDX Mnemonic = "LDX"
	LDY Mnemonic = "LDY"
	STA Mnemonic = "STA"
	STX Mnemonic = "STX"
	STY Mnemonic = "STY"
	TAX Mnemonic = "TAX"
	TXA Mnemonic = "TXA"
	TAY Mnemonic = "TAY"
	TYA Mnemonic = "TYA"
	TSX Mnemonic = "TSX"
	TXS Mnemonic = "TXS"
	PHA Mnemonic = "PHA"
	PLA Mnemonic = "PLA"
	PHP Mnemonic = "PHP"
	PLP Mnemonic = "PLP"
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
}

// OpcodeProps ...
var OpcodeProps map[Opcode]OpcodeProp

func init() {
	OpcodeProps = map[Opcode]OpcodeProp{
		0x69: OpcodeProp{ADC, Immediate},
		0x65: OpcodeProp{ADC, ZeroPage},
		0x75: OpcodeProp{ADC, IndexedZeroPageX},
		0x6D: OpcodeProp{ADC, Absolute},
		0x7D: OpcodeProp{ADC, IndexedAbsoluteX},
		0x79: OpcodeProp{ADC, IndexedAbsoluteY},
		0x61: OpcodeProp{ADC, IndexedIndirect},
		0x71: OpcodeProp{ADC, IndirectIndexed},
		0xE9: OpcodeProp{SBC, Immediate},
		0xE5: OpcodeProp{SBC, ZeroPage},
		0xF5: OpcodeProp{SBC, IndexedZeroPageX},
		0xED: OpcodeProp{SBC, Absolute},
		0xFD: OpcodeProp{SBC, IndexedAbsoluteX},
		0xF9: OpcodeProp{SBC, IndexedAbsoluteY},
		0xE1: OpcodeProp{SBC, IndexedIndirect},
		0xF1: OpcodeProp{SBC, IndirectIndexed},
		0x29: OpcodeProp{AND, Immediate},
		0x25: OpcodeProp{AND, ZeroPage},
		0x35: OpcodeProp{AND, IndexedZeroPageX},
		0x2D: OpcodeProp{AND, Absolute},
		0x3D: OpcodeProp{AND, IndexedAbsoluteX},
		0x39: OpcodeProp{AND, IndexedAbsoluteY},
		0x21: OpcodeProp{AND, IndexedIndirect},
		0x31: OpcodeProp{AND, IndirectIndexed},
		0x09: OpcodeProp{ORA, Immediate},
		0x05: OpcodeProp{ORA, ZeroPage},
		0x15: OpcodeProp{ORA, IndexedZeroPageX},
		0x0D: OpcodeProp{ORA, Absolute},
		0x1D: OpcodeProp{ORA, IndexedAbsoluteX},
		0x19: OpcodeProp{ORA, IndexedAbsoluteY},
		0x01: OpcodeProp{ORA, IndexedIndirect},
		0x11: OpcodeProp{ORA, IndirectIndexed},
		0x49: OpcodeProp{EOR, Immediate},
		0x45: OpcodeProp{EOR, ZeroPage},
		0x55: OpcodeProp{EOR, IndexedZeroPageX},
		0x4D: OpcodeProp{EOR, Absolute},
		0x5D: OpcodeProp{EOR, IndexedAbsoluteX},
		0x59: OpcodeProp{EOR, IndexedAbsoluteY},
		0x41: OpcodeProp{EOR, IndexedIndirect},
		0x51: OpcodeProp{EOR, IndirectIndexed},
		0x0A: OpcodeProp{ASL, Accumulator},
		0x06: OpcodeProp{ASL, ZeroPage},
		0x16: OpcodeProp{ASL, IndexedZeroPageX},
		0x0E: OpcodeProp{ASL, Absolute},
		0x1E: OpcodeProp{ASL, IndexedAbsoluteX},
		0x4A: OpcodeProp{LSR, Accumulator},
		0x46: OpcodeProp{LSR, ZeroPage},
		0x56: OpcodeProp{LSR, IndexedZeroPageX},
		0x4E: OpcodeProp{LSR, Absolute},
		0x5E: OpcodeProp{LSR, IndexedAbsoluteX},
		0x2A: OpcodeProp{ROL, Accumulator},
		0x26: OpcodeProp{ROL, ZeroPage},
		0x36: OpcodeProp{ROL, IndexedZeroPageX},
		0x2E: OpcodeProp{ROL, Absolute},
		0x3E: OpcodeProp{ROL, IndexedAbsoluteX},
		0x6A: OpcodeProp{ROR, Accumulator},
		0x66: OpcodeProp{ROR, ZeroPage},
		0x76: OpcodeProp{ROR, IndexedZeroPageX},
		0x6E: OpcodeProp{ROR, Absolute},
		0x7E: OpcodeProp{ROR, IndexedAbsoluteX},
		0x90: OpcodeProp{BCC, Relative},
		0xB0: OpcodeProp{BCS, Relative},
		0xF0: OpcodeProp{BEQ, Relative},
		0xD0: OpcodeProp{BNE, Relative},
		0x50: OpcodeProp{BVC, Relative},
		0x70: OpcodeProp{BVS, Relative},
		0x10: OpcodeProp{BPL, Relative},
		0x30: OpcodeProp{BMI, Relative},
		0x24: OpcodeProp{BIT, ZeroPage},
		0x2C: OpcodeProp{BIT, Absolute},
		0x4C: OpcodeProp{JMP, Absolute},
		0x6C: OpcodeProp{JMP, AbsoluteIndirect},
		0x20: OpcodeProp{JSR, Absolute},
		0x60: OpcodeProp{RTS, Implied},
		0x00: OpcodeProp{BRK, Implied},
		0x40: OpcodeProp{RTI, Implied},
		0xC9: OpcodeProp{CMP, Immediate},
		0xC5: OpcodeProp{CMP, ZeroPage},
		0xD5: OpcodeProp{CMP, IndexedZeroPageX},
		0xCD: OpcodeProp{CMP, Absolute},
		0xDD: OpcodeProp{CMP, IndexedAbsoluteX},
		0xD9: OpcodeProp{CMP, IndexedAbsoluteY},
		0xC1: OpcodeProp{CMP, IndexedIndirect},
		0xD1: OpcodeProp{CMP, IndirectIndexed},
		0xE0: OpcodeProp{CPX, Immediate},
		0xE4: OpcodeProp{CPX, ZeroPage},
		0xEC: OpcodeProp{CPX, Absolute},
		0xC0: OpcodeProp{CPY, Immediate},
		0xC4: OpcodeProp{CPY, ZeroPage},
		0xCC: OpcodeProp{CPY, Absolute},
		0xE6: OpcodeProp{INC, ZeroPage},
		0xF6: OpcodeProp{INC, IndexedZeroPageX},
		0xEE: OpcodeProp{INC, Absolute},
		0xFE: OpcodeProp{INC, IndexedAbsoluteX},
		0xC6: OpcodeProp{DEC, ZeroPage},
		0xD6: OpcodeProp{DEC, IndexedZeroPageX},
		0xCE: OpcodeProp{DEC, Absolute},
		0xDE: OpcodeProp{DEC, IndexedAbsoluteX},
		0xE8: OpcodeProp{INX, Implied},
		0xCA: OpcodeProp{DEX, Implied},
		0xC8: OpcodeProp{INY, Implied},
		0x88: OpcodeProp{DEY, Implied},
		0x18: OpcodeProp{CLC, Implied},
		0x38: OpcodeProp{SEC, Implied},
		0x58: OpcodeProp{CLI, Implied},
		0x78: OpcodeProp{SEI, Implied},
		0xD8: OpcodeProp{CLD, Implied},
		0xF8: OpcodeProp{SED, Implied},
		0xB8: OpcodeProp{CLV, Implied},
		0xA9: OpcodeProp{LDA, Immediate},
		0xA5: OpcodeProp{LDA, ZeroPage},
		0xB5: OpcodeProp{LDA, IndexedZeroPageX},
		0xAD: OpcodeProp{LDA, Absolute},
		0xBD: OpcodeProp{LDA, IndexedAbsoluteX},
		0xB9: OpcodeProp{LDA, IndexedAbsoluteY},
		0xA1: OpcodeProp{LDA, IndexedIndirect},
		0xB1: OpcodeProp{LDA, IndirectIndexed},
		0xA2: OpcodeProp{LDX, Immediate},
		0xA6: OpcodeProp{LDX, ZeroPage},
		0xB6: OpcodeProp{LDX, IndexedZeroPageY},
		0xAE: OpcodeProp{LDX, Absolute},
		0xBE: OpcodeProp{LDX, IndexedAbsoluteY},
		0xA0: OpcodeProp{LDY, Immediate},
		0xA4: OpcodeProp{LDY, ZeroPage},
		0xB4: OpcodeProp{LDY, IndexedZeroPageX},
		0xAC: OpcodeProp{LDY, Absolute},
		0xBC: OpcodeProp{LDY, IndexedAbsoluteX},
		0x85: OpcodeProp{STA, ZeroPage},
		0x95: OpcodeProp{STA, IndexedZeroPageX},
		0x8D: OpcodeProp{STA, Absolute},
		0x9D: OpcodeProp{STA, IndexedAbsoluteX},
		0x99: OpcodeProp{STA, IndexedAbsoluteY},
		0x81: OpcodeProp{STA, IndexedIndirect},
		0x91: OpcodeProp{STA, IndirectIndexed},
		0x86: OpcodeProp{STX, ZeroPage},
		0x96: OpcodeProp{STX, IndexedZeroPageY},
		0x8E: OpcodeProp{STX, Absolute},
		0x84: OpcodeProp{STY, ZeroPage},
		0x94: OpcodeProp{STY, IndexedZeroPageX},
		0x8C: OpcodeProp{STY, Absolute},
		0xAA: OpcodeProp{TAX, Implied},
		0x8A: OpcodeProp{TXA, Implied},
		0xA8: OpcodeProp{TAY, Implied},
		0x98: OpcodeProp{TYA, Implied},
		0x9A: OpcodeProp{TXS, Implied},
		0xBA: OpcodeProp{TSX, Implied},
		0x48: OpcodeProp{PHA, Implied},
		0x68: OpcodeProp{PLA, Implied},
		0x08: OpcodeProp{PHP, Implied},
		0x28: OpcodeProp{PLP, Implied},
		0xEA: OpcodeProp{NOP, Implied},
	}
}
