package model

import (
	"fmt"
)

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
func (c *CPU) Run() error {
	for {
		// PC（プログラムカウンタ）からオペコードをフェッチ（PCをインクリメント）
		oc, err := c.fetchOpcode()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// 命令とアドレッシング・モードを判別
		ocp, err := decodeOpcode(oc)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// （必要であれば）オペランドをフェッチ（PCをインクリメント）
		ors, err := c.fetchOperands(ocp.AddressingMode)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// 命令を実行
		if err := c.exec(ocp.Mnemonic, ors...); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (*OpcodeProp, error) {
	if p, ok := OpcodeProps[o]; ok {
		return &p, nil
	}
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (uint32, error) {
	// TODO 実装
	return 0, nil
}

// fetchOpcode ...
func (c *CPU) fetchOpcode() (Opcode, error) {
	v, err := c.fetch()
	if err != nil {
		return ErrorOpcode, fmt.Errorf("%w", err)
	}
	return Opcode(v), nil
}

// fetchOperands ...
func (c *CPU) fetchOperands(mode AddressingMode) ([]Operand, error) {
	// TODO 実装
	return []Operand{}, nil
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
func (c *CPU) exec(mne Mnemonic, opr ...Operand) error {
	// TODO 実装
	// （必要であれば）演算対象となるアドレスを算出
	return nil
}
