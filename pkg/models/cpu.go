package model

// Registers ...
type Registers struct {
	A  uint8
	X  uint8
	Y  uint8
	S  uint8
	P  *StatusRegister
	PC uint16
}

// NewRegisters ...
func NewRegisters() *Registers {
	return &Registers{
		P: NewStatusRegister(),
	}
}

// StatusRegister ...
// https://qiita.com/bokuweb/items/1575337bef44ae82f4d3#%E3%82%B9%E3%83%86%E3%83%BC%E3%82%BF%E3%82%B9%E3%83%AC%E3%82%B8%E3%82%B9%E3%82%BF
type StatusRegister struct {
	Negative  bool // bit7	N	ネガティブ	演算結果のbit7が1の時にセット
	Overflow  bool // bit6	V	オーバーフロー	P演算結果がオーバーフローを起こした時にセット
	Reserved  bool // bit5	R	予約済み	常にセットされている
	Break     bool // bit4	B	ブレークモード	BRK発生時にセット、IRQ発生時にクリア
	Decimal   bool // bit3	D	デシマルモード	0:デフォルト、1:BCDモード (未実装)
	Interrupt bool // bit2	I	IRQ禁止	0:IRQ許可、1:IRQ禁止
	Zero      bool // bit1	Z	ゼロ	演算結果が0の時にセット
	Carry     bool // bit0	C	キャリー	キャリー発生時にセット
}

// NewStatusRegister ...
func NewStatusRegister() *StatusRegister {
	return &StatusRegister{
		Negative:  false,
		Overflow:  false,
		Reserved:  true,
		Break:     false,
		Decimal:   false,
		Interrupt: false,
		Zero:      false,
		Carry:     false,
	}
}

// CPUBus ...
type CPUBus struct {
	programRom *PRGROM
}

// NewCPUBus ...
func NewCPUBus(p *PRGROM) *CPUBus {
	return &CPUBus{
		programRom: p,
	}
}

// CPU ...
type CPU struct {
	registers *Registers
	bus       *CPUBus
}

// NewCPU ...
func NewCPU(p *PRGROM) *CPU {
	return &CPU{
		registers: NewRegisters(),
		bus:       NewCPUBus(p),
	}
}

// read ...
func (c *CPUBus) read(addr uint32) byte {
	// TODO 実装
	// 0x0000～0x07FF	0x0800	WRAM
	// 0x0800～0x1FFF	-	WRAMのミラー
	// 0x2000～0x2007	0x0008	PPU レジスタ
	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	// 0x4000～0x401F	0x0020	APU I/O、PAD
	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	// 0x6000～0x7FFF	0x2000	拡張RAM
	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM

	return 0
}

// write ...
func (c *CPUBus) write(addr uint32, data byte) {
	// TODO 実装
	// 0x0000～0x07FF	0x0800	WRAM
	// 0x0800～0x1FFF	-	WRAMのミラー
	// 0x2000～0x2007	0x0008	PPU レジスタ
	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	// 0x4000～0x401F	0x0020	APU I/O、PAD
	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	// 0x6000～0x7FFF	0x2000	拡張RAM
	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
}

// Run ...
func (c *CPU) Run() {
	// TODO 実装
	// https://qiita.com/bokuweb/items/1575337bef44ae82f4d3#%E5%AE%9F%E8%A3%85%E3%82%A4%E3%83%A1%E3%83%BC%E3%82%B8

	// PC（プログラムカウンタ）からオペコードをフェッチ（PCをインクリメント）
	oc := c.fetchOpcode()

	// 命令とアドレッシング・モードを判別
	mne, mode, _ := decodeOpcode(oc)

	// （必要であれば）オペランドをフェッチ（PCをインクリメント）
	ors := c.fetchOperands(mode)

	// 命令を実行
	c.exec(mne, ors...)
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (Mnemonic, AddressingMode, int) {
	// TODO 実装
	return "", "", 0
}

// fetch ...
func (c *CPU) fetch() uint32 {
	// TODO 実装
	return 0
}

// fetchOpcode ...
func (c *CPU) fetchOpcode() Opcode {
	return Opcode(c.fetch())
}

// fetchOperands ...
func (c *CPU) fetchOperands(mode AddressingMode) []Operand {
	// TODO 実装
	return []Operand{}
}

// interruptNMI ...
func (c *CPU) interruptNMI() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// interruptRESET ...
func (c *CPU) interruptRESET() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// interruptBRK ...
func (c *CPU) interruptBRK() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// interruptIRQ ...
func (c *CPU) interruptIRQ() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// exec ...
func (c *CPU) exec(mne Mnemonic, opr ...Operand) {
	// TODO 実装
	// （必要であれば）演算対象となるアドレスを算出
}
