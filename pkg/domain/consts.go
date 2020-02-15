package domain

const (
	// ResolutionWidth ... 解像度(横)
	ResolutionWidth = 256
	// ResolutionHeight ... 解像度(縦)
	ResolutionHeight = 240
	// SpriteWidth ... スプライトサイズ（横）
	SpriteWidth = 8
	// SpriteHeight ... スプライトサイズ（横）
	SpriteHeight = 8

	// TileWidth ... タイルサイズ（横）
	TileWidth = 8
	// TileHeight ... タイルサイズ（縦）
	TileHeight = 8

	// NameTableWidth ... ネームテーブルサイズ（横）
	NameTableWidth = ResolutionWidth / TileWidth
	// NameTableHeight ... ネームテーブルサイズ（縦）
	NameTableHeight = ResolutionHeight / TileHeight
	// NameTableBaseAddress ... ネームテーブルの基準アドレス(ネームテーブル0の開始アドレス)
	NameTableBaseAddress = uint16(0x2000)

	// AttributeTableBaseAddress ... 属性テーブルの基準アドレス(属性テーブル0の開始アドレス)
	AttributeTableBaseAddress = NameTableBaseAddress + NameTableWidth*NameTableHeight
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

// Opcode ...
type Opcode uint8

const (
	// ErrorOpcode ...
	ErrorOpcode Opcode = 0
)

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
	// STP ...
	STP Mnemonic = "STP"
	// SLO ...
	SLO Mnemonic = "SLO"
	// ANC ...
	ANC Mnemonic = "ANC"
	// RLA ...
	RLA Mnemonic = "RLA"
	// SRE ...
	SRE Mnemonic = "SRE"
	// RRA ...
	RRA Mnemonic = "RRA"
	// ARR ...
	ARR Mnemonic = "ARR"
	// SAX ...
	SAX Mnemonic = "SAX"
	// XAA ...
	XAA Mnemonic = "XAA"
	// AHX ...
	AHX Mnemonic = "AHX"
	// TAS ...
	TAS Mnemonic = "TAS"
	// SHX ...
	SHX Mnemonic = "SHX"
	// SHY ...
	SHY Mnemonic = "SHY"
	// LAX ...
	LAX Mnemonic = "LAX"
	// LAS ...
	LAS Mnemonic = "LAS"
	// DCP ...
	DCP Mnemonic = "DCP"
	// AXS ...
	AXS Mnemonic = "AXS"
	// ISC ...
	ISC Mnemonic = "ISC"
	// ISB ...
	ISB Mnemonic = "ISB"
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
	Documented     bool
}

// OpcodeProps ...
var OpcodeProps = map[Opcode]OpcodeProp{
	0x00: OpcodeProp{BRK, Implied, 7, true},
	0x01: OpcodeProp{ORA, IndexedIndirect, 6, true},
	0x02: OpcodeProp{STP, Implied, 4, false},
	0x03: OpcodeProp{SLO, IndexedIndirect, 8, false},
	0x04: OpcodeProp{NOP, ZeroPage, 3, false},
	0x05: OpcodeProp{ORA, ZeroPage, 3, true},
	0x06: OpcodeProp{ASL, ZeroPage, 5, true},
	0x07: OpcodeProp{SLO, ZeroPage, 5, false},
	0x08: OpcodeProp{PHP, Implied, 3, true},
	0x09: OpcodeProp{ORA, Immediate, 2, true},
	0x0A: OpcodeProp{ASL, Accumulator, 2, true},
	0x0B: OpcodeProp{ANC, Implied, 2, false},
	0x0C: OpcodeProp{NOP, Absolute, 4, false},
	0x0D: OpcodeProp{ORA, Absolute, 4, true},
	0x0E: OpcodeProp{ASL, Absolute, 6, true},
	0x0F: OpcodeProp{SLO, Absolute, 6, false},
	0x10: OpcodeProp{BPL, Relative, 2, true},
	0x11: OpcodeProp{ORA, IndirectIndexed, 5, true},
	0x12: OpcodeProp{STP, Implied, 4, false},
	0x13: OpcodeProp{SLO, IndirectIndexed, 8, false},
	0x14: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0x15: OpcodeProp{ORA, IndexedZeroPageX, 4, true},
	0x16: OpcodeProp{ASL, IndexedZeroPageX, 6, true},
	0x17: OpcodeProp{SLO, IndexedZeroPageX, 6, false},
	0x18: OpcodeProp{CLC, Implied, 2, true},
	0x19: OpcodeProp{ORA, IndexedAbsoluteY, 4, true},
	0x1A: OpcodeProp{NOP, Implied, 2, false},
	0x1B: OpcodeProp{SLO, IndexedAbsoluteY, 7, false},
	0x1C: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0x1D: OpcodeProp{ORA, IndexedAbsoluteX, 4, true},
	0x1E: OpcodeProp{ASL, IndexedAbsoluteX, 7, true},
	0x1F: OpcodeProp{SLO, IndexedAbsoluteX, 7, false},
	0x20: OpcodeProp{JSR, Absolute, 6, true},
	0x21: OpcodeProp{AND, IndexedIndirect, 6, true},
	0x22: OpcodeProp{STP, Implied, 4, false},
	0x23: OpcodeProp{RLA, IndexedIndirect, 8, false},
	0x24: OpcodeProp{BIT, ZeroPage, 3, true},
	0x25: OpcodeProp{AND, ZeroPage, 3, true},
	0x26: OpcodeProp{ROL, ZeroPage, 5, true},
	0x27: OpcodeProp{RLA, ZeroPage, 5, false},
	0x28: OpcodeProp{PLP, Implied, 4, true},
	0x29: OpcodeProp{AND, Immediate, 2, true},
	0x2A: OpcodeProp{ROL, Accumulator, 2, true},
	0x2B: OpcodeProp{ANC, Immediate, 4, false},
	0x2C: OpcodeProp{BIT, Absolute, 4, true},
	0x2D: OpcodeProp{AND, Absolute, 4, true},
	0x2E: OpcodeProp{ROL, Absolute, 6, true},
	0x2F: OpcodeProp{RLA, Absolute, 6, false},
	0x30: OpcodeProp{BMI, Relative, 2, true},
	0x31: OpcodeProp{AND, IndirectIndexed, 5, true},
	0x32: OpcodeProp{STP, Implied, 4, false},
	0x33: OpcodeProp{RLA, IndirectIndexed, 8, false},
	0x34: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0x35: OpcodeProp{AND, IndexedZeroPageX, 4, true},
	0x36: OpcodeProp{ROL, IndexedZeroPageX, 6, true},
	0x37: OpcodeProp{RLA, IndexedZeroPageX, 6, false},
	0x38: OpcodeProp{SEC, Implied, 2, true},
	0x39: OpcodeProp{AND, IndexedAbsoluteY, 4, true},
	0x3A: OpcodeProp{NOP, Implied, 2, false},
	0x3B: OpcodeProp{RLA, IndexedAbsoluteY, 7, false},
	0x3C: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0x3D: OpcodeProp{AND, IndexedAbsoluteX, 4, true},
	0x3E: OpcodeProp{ROL, IndexedAbsoluteX, 7, true},
	0x3F: OpcodeProp{RLA, IndexedAbsoluteX, 7, false},
	0x40: OpcodeProp{RTI, Implied, 6, true},
	0x41: OpcodeProp{EOR, IndexedIndirect, 6, true},
	0x42: OpcodeProp{STP, Implied, 4, false},
	0x43: OpcodeProp{SRE, IndexedIndirect, 8, false},
	0x44: OpcodeProp{NOP, ZeroPage, 3, false},
	0x45: OpcodeProp{EOR, ZeroPage, 3, true},
	0x46: OpcodeProp{LSR, ZeroPage, 5, true},
	0x47: OpcodeProp{SRE, ZeroPage, 5, false},
	0x48: OpcodeProp{PHA, Implied, 3, true},
	0x49: OpcodeProp{EOR, Immediate, 2, true},
	0x4A: OpcodeProp{LSR, Accumulator, 2, true},
	0x4B: OpcodeProp{LSR, Immediate, 5, false},
	0x4C: OpcodeProp{JMP, Absolute, 3, true},
	0x4D: OpcodeProp{EOR, Absolute, 4, true},
	0x4E: OpcodeProp{LSR, Absolute, 6, true},
	0x4F: OpcodeProp{SRE, Absolute, 6, false},
	0x50: OpcodeProp{BVC, Relative, 2, true},
	0x51: OpcodeProp{EOR, IndirectIndexed, 5, true},
	0x52: OpcodeProp{STP, Implied, 4, false},
	0x53: OpcodeProp{SRE, IndirectIndexed, 8, false},
	0x54: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0x55: OpcodeProp{EOR, IndexedZeroPageX, 4, true},
	0x56: OpcodeProp{LSR, IndexedZeroPageX, 6, true},
	0x57: OpcodeProp{SRE, IndexedZeroPageX, 6, false},
	0x58: OpcodeProp{CLI, Implied, 2, true},
	0x59: OpcodeProp{EOR, IndexedAbsoluteY, 4, true},
	0x5A: OpcodeProp{NOP, Implied, 2, false},
	0x5B: OpcodeProp{SRE, IndexedAbsoluteY, 7, false},
	0x5C: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0x5D: OpcodeProp{EOR, IndexedAbsoluteX, 4, true},
	0x5E: OpcodeProp{LSR, IndexedAbsoluteX, 7, true},
	0x5F: OpcodeProp{SRE, IndexedAbsoluteX, 7, false},
	0x60: OpcodeProp{RTS, Implied, 6, true},
	0x61: OpcodeProp{ADC, IndexedIndirect, 6, true},
	0x62: OpcodeProp{STP, Implied, 4, false},
	0x63: OpcodeProp{RRA, IndexedIndirect, 8, false},
	0x64: OpcodeProp{NOP, ZeroPage, 3, false},
	0x65: OpcodeProp{ADC, ZeroPage, 3, true},
	0x66: OpcodeProp{ROR, ZeroPage, 5, true},
	0x67: OpcodeProp{RRA, ZeroPage, 5, false},
	0x68: OpcodeProp{PLA, Implied, 4, true},
	0x69: OpcodeProp{ADC, Immediate, 2, true},
	0x6A: OpcodeProp{ROR, Accumulator, 2, true},
	0x6B: OpcodeProp{ARR, Immediate, 4, false},
	0x6C: OpcodeProp{JMP, AbsoluteIndirect, 5, true},
	0x6D: OpcodeProp{ADC, Absolute, 4, true},
	0x6E: OpcodeProp{ROR, Absolute, 6, true},
	0x6F: OpcodeProp{RRA, Absolute, 6, false},
	0x70: OpcodeProp{BVS, Relative, 2, true},
	0x71: OpcodeProp{ADC, IndirectIndexed, 5, true},
	0x72: OpcodeProp{STP, Implied, 4, false},
	0x73: OpcodeProp{RRA, IndirectIndexed, 8, false},
	0x74: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0x75: OpcodeProp{ADC, IndexedZeroPageX, 4, true},
	0x76: OpcodeProp{ROR, IndexedZeroPageX, 6, true},
	0x77: OpcodeProp{RRA, IndexedZeroPageX, 6, false},
	0x78: OpcodeProp{SEI, Implied, 2, true},
	0x79: OpcodeProp{ADC, IndexedAbsoluteY, 4, true},
	0x7A: OpcodeProp{NOP, Implied, 2, false},
	0x7B: OpcodeProp{RRA, IndexedAbsoluteY, 7, false},
	0x7C: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0x7D: OpcodeProp{ADC, IndexedAbsoluteX, 4, true},
	0x7E: OpcodeProp{ROR, IndexedAbsoluteX, 7, true},
	0x7F: OpcodeProp{RRA, IndexedAbsoluteX, 7, false},
	0x80: OpcodeProp{NOP, Immediate, 2, false},
	0x81: OpcodeProp{STA, IndexedIndirect, 6, true},
	0x82: OpcodeProp{NOP, Immediate, 2, false},
	0x83: OpcodeProp{SAX, IndexedIndirect, 6, false},
	0x84: OpcodeProp{STY, ZeroPage, 3, true},
	0x85: OpcodeProp{STA, ZeroPage, 3, true},
	0x86: OpcodeProp{STX, ZeroPage, 3, true},
	0x87: OpcodeProp{SAX, ZeroPage, 3, false},
	0x88: OpcodeProp{DEY, Implied, 2, true},
	0x89: OpcodeProp{NOP, Immediate, 2, false},
	0x8A: OpcodeProp{TXA, Implied, 2, true},
	0x8B: OpcodeProp{XAA, Immediate, 4, false},
	0x8C: OpcodeProp{STY, Absolute, 4, true},
	0x8D: OpcodeProp{STA, Absolute, 4, true},
	0x8E: OpcodeProp{STX, Absolute, 4, true},
	0x8F: OpcodeProp{SAX, Absolute, 4, false},
	0x90: OpcodeProp{BCC, Relative, 2, true},
	0x91: OpcodeProp{STA, IndirectIndexed, 6, true},
	0x92: OpcodeProp{STP, Implied, 4, false},
	0x93: OpcodeProp{AHX, IndirectIndexed, 4, false},
	0x94: OpcodeProp{STY, IndexedZeroPageX, 4, true},
	0x95: OpcodeProp{STA, IndexedZeroPageX, 4, true},
	0x96: OpcodeProp{STX, IndexedZeroPageY, 4, true},
	0x97: OpcodeProp{SAX, IndexedZeroPageY, 4, false},
	0x98: OpcodeProp{TYA, Implied, 2, true},
	0x99: OpcodeProp{STA, IndexedAbsoluteY, 5, true},
	0x9A: OpcodeProp{TXS, Implied, 2, true},
	0x9B: OpcodeProp{TAS, IndexedAbsoluteY, 4, false},
	0x9C: OpcodeProp{SHY, IndexedAbsoluteX, 4, false},
	0x9D: OpcodeProp{STA, IndexedAbsoluteX, 5, true},
	0x9E: OpcodeProp{SHX, IndexedAbsoluteY, 4, false},
	0x9F: OpcodeProp{AHX, IndexedAbsoluteY, 4, false},
	0xA0: OpcodeProp{LDY, Immediate, 2, true},
	0xA1: OpcodeProp{LDA, IndexedIndirect, 6, true},
	0xA2: OpcodeProp{LDX, Immediate, 2, true},
	0xA3: OpcodeProp{LAX, IndexedIndirect, 6, false},
	0xA4: OpcodeProp{LDY, ZeroPage, 3, true},
	0xA5: OpcodeProp{LDA, ZeroPage, 3, true},
	0xA6: OpcodeProp{LDX, ZeroPage, 3, true},
	0xA7: OpcodeProp{LAX, ZeroPage, 3, false},
	0xA8: OpcodeProp{TAY, Implied, 2, true},
	0xA9: OpcodeProp{LDA, Immediate, 2, true},
	0xAA: OpcodeProp{TAX, Implied, 2, true},
	0xAB: OpcodeProp{LAX, Immediate, 4, false},
	0xAC: OpcodeProp{LDY, Absolute, 4, true},
	0xAD: OpcodeProp{LDA, Absolute, 4, true},
	0xAE: OpcodeProp{LDX, Absolute, 4, true},
	0xAF: OpcodeProp{LAX, Absolute, 4, false},
	0xB0: OpcodeProp{BCS, Relative, 2, true},
	0xB1: OpcodeProp{LDA, IndirectIndexed, 5, true},
	0xB2: OpcodeProp{STP, Implied, 4, false},
	0xB3: OpcodeProp{LAX, IndirectIndexed, 5, false},
	0xB4: OpcodeProp{LDY, IndexedZeroPageX, 4, true},
	0xB5: OpcodeProp{LDA, IndexedZeroPageX, 4, true},
	0xB6: OpcodeProp{LDX, IndexedZeroPageY, 4, true},
	0xB7: OpcodeProp{LAX, IndexedZeroPageY, 4, false},
	0xB8: OpcodeProp{CLV, Implied, 2, true},
	0xB9: OpcodeProp{LDA, IndexedAbsoluteY, 4, true},
	0xBA: OpcodeProp{TSX, Implied, 2, true},
	0xBB: OpcodeProp{LAS, IndexedAbsoluteY, 4, false},
	0xBC: OpcodeProp{LDY, IndexedAbsoluteX, 4, true},
	0xBD: OpcodeProp{LDA, IndexedAbsoluteX, 4, true},
	0xBE: OpcodeProp{LDX, IndexedAbsoluteY, 4, true},
	0xBF: OpcodeProp{LAX, IndexedAbsoluteY, 4, false},
	0xC0: OpcodeProp{CPY, Immediate, 2, true},
	0xC1: OpcodeProp{CMP, IndexedIndirect, 6, true},
	0xC2: OpcodeProp{NOP, Immediate, 2, false},
	0xC3: OpcodeProp{DCP, IndexedIndirect, 8, false},
	0xC4: OpcodeProp{CPY, ZeroPage, 3, true},
	0xC5: OpcodeProp{CMP, ZeroPage, 3, true},
	0xC6: OpcodeProp{DEC, ZeroPage, 5, true},
	0xC7: OpcodeProp{DCP, ZeroPage, 5, false},
	0xC8: OpcodeProp{INY, Implied, 2, true},
	0xC9: OpcodeProp{CMP, Immediate, 2, true},
	0xCA: OpcodeProp{DEX, Implied, 2, true},
	0xCB: OpcodeProp{AXS, Immediate, 4, false},
	0xCC: OpcodeProp{CPY, Absolute, 4, true},
	0xCD: OpcodeProp{CMP, Absolute, 4, true},
	0xCE: OpcodeProp{DEC, Absolute, 6, true},
	0xCF: OpcodeProp{DCP, Absolute, 6, false},
	0xD0: OpcodeProp{BNE, Relative, 2, true},
	0xD1: OpcodeProp{CMP, IndirectIndexed, 5, true},
	0xD2: OpcodeProp{STP, Implied, 4, false},
	0xD3: OpcodeProp{DCP, IndirectIndexed, 8, false},
	0xD4: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0xD5: OpcodeProp{CMP, IndexedZeroPageX, 4, true},
	0xD6: OpcodeProp{DEC, IndexedZeroPageX, 6, true},
	0xD7: OpcodeProp{DCP, IndexedZeroPageX, 6, false},
	0xD8: OpcodeProp{CLD, Implied, 2, true},
	0xD9: OpcodeProp{CMP, IndexedAbsoluteY, 4, true},
	0xDA: OpcodeProp{NOP, Implied, 2, false},
	0xDB: OpcodeProp{DCP, IndexedAbsoluteY, 7, false},
	0xDC: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0xDD: OpcodeProp{CMP, IndexedAbsoluteX, 4, true},
	0xDE: OpcodeProp{DEC, IndexedAbsoluteX, 7, true},
	0xDF: OpcodeProp{DCP, IndexedAbsoluteX, 7, false},
	0xE0: OpcodeProp{CPX, Immediate, 2, true},
	0xE1: OpcodeProp{SBC, IndexedIndirect, 6, true},
	0xE2: OpcodeProp{NOP, Immediate, 2, false},
	0xE3: OpcodeProp{ISB, IndexedIndirect, 8, false},
	0xE4: OpcodeProp{CPX, ZeroPage, 3, true},
	0xE5: OpcodeProp{SBC, ZeroPage, 3, true},
	0xE6: OpcodeProp{INC, ZeroPage, 5, true},
	0xE7: OpcodeProp{ISB, ZeroPage, 5, false},
	0xE8: OpcodeProp{INX, Implied, 2, true},
	0xE9: OpcodeProp{SBC, Immediate, 2, true},
	0xEA: OpcodeProp{NOP, Implied, 2, true},
	0xEB: OpcodeProp{SBC, Immediate, 2, false},
	0xEC: OpcodeProp{CPX, Absolute, 4, true},
	0xED: OpcodeProp{SBC, Absolute, 4, true},
	0xEE: OpcodeProp{INC, Absolute, 6, true},
	0xEF: OpcodeProp{ISB, Absolute, 6, false},
	0xF0: OpcodeProp{BEQ, Relative, 2, true},
	0xF1: OpcodeProp{SBC, IndirectIndexed, 5, true},
	0xF2: OpcodeProp{STP, Implied, 4, false},
	0xF3: OpcodeProp{ISB, IndirectIndexed, 8, false},
	0xF4: OpcodeProp{NOP, IndexedZeroPageX, 4, false},
	0xF5: OpcodeProp{SBC, IndexedZeroPageX, 4, true},
	0xF6: OpcodeProp{INC, IndexedZeroPageX, 6, true},
	0xF7: OpcodeProp{ISB, IndexedZeroPageX, 6, false},
	0xF8: OpcodeProp{SED, Implied, 2, true},
	0xF9: OpcodeProp{SBC, IndexedAbsoluteY, 4, true},
	0xFA: OpcodeProp{NOP, Implied, 2, false},
	0xFB: OpcodeProp{ISB, IndexedAbsoluteY, 7, false},
	0xFC: OpcodeProp{NOP, IndexedAbsoluteX, 4, false},
	0xFD: OpcodeProp{SBC, IndexedAbsoluteX, 4, true},
	0xFE: OpcodeProp{INC, IndexedAbsoluteX, 7, true},
	0xFF: OpcodeProp{ISB, IndexedAbsoluteX, 7, false},
}
