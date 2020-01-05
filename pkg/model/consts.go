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

const (
	// ResolutionWidth ... 解像度(横)
	ResolutionWidth = 256
	// ResolutionHeight ... 解像度(縦)
	ResolutionHeight = 240
	// SpriteWidth ... スプライトサイズ（横）
	SpriteWidth = 8
	// SpriteHeight ... スプライトサイズ（横）
	SpriteHeight = 8
)

var colors = [][]byte{
	{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96},
	{0xA1, 0x00, 0x5E}, {0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00},
	{0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00}, {0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E},
	{0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05}, {0x05, 0x05, 0x05},
	{0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
	{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00},
	{0xC4, 0x62, 0x00}, {0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55},
	{0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21}, {0x09, 0x09, 0x09}, {0x09, 0x09, 0x09},
	{0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF}, {0xD4, 0x80, 0xFF},
	{0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
	{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4},
	{0x05, 0xFB, 0xFF}, {0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D},
	{0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF}, {0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB},
	{0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0}, {0xFF, 0xEF, 0xA6},
	{0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
	{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},
}
