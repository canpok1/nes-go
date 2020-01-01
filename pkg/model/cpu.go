package model

import (
	"fmt"
	"log"
	"time"
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

// String ...
func (r *Registers) String() string {
	return fmt.Sprintf(
		"{A:%#v, X:%#v, Y:%#v, S:%#v, P:%#v, PC:%#v}",
		r.A,
		r.X,
		r.Y,
		r.S,
		r.P,
		r.PC,
	)

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
	wram []byte
	wramMirror []byte
	ppuRegister []byte
	ppuRegisterMirror []byte
	io []byte
	exrom []byte
	exram []byte
	programRom *PRGROM
}

// NewCPUBus ...
func NewCPUBus(p *PRGROM) *CPUBus {
	return &CPUBus{
		wram: make([]byte, 0x0800),
		ppuRegister: make([]byte, 0x0008),
		io: make([]byte, 0x0020),
		exrom: make([]byte, 0x1FE0),
		exram: make([]byte, 0x2000),
		programRom: p,
	}
}

// CPU ...
type CPU struct {
	registers *Registers
	bus       *CPUBus
}

// NewCPU ...
func NewCPU(p *ROM) *CPU {
	return &CPU{
		registers: NewRegisters(),
		bus:       NewCPUBus(p.prgrom),
	}
}

// String ...
func (c *CPU) String() string {
	return fmt.Sprintf("registers %v", c.registers.String())
}

// read ...
func (c *CPUBus) read(addr uint32) (byte, error) {
	var data byte 
	var err error
	defer func() {
		if err != nil {
			log.Printf("CPUBus.read[addr=%#v] => %#v", addr, err)
		} else {
			log.Printf("CPUBus.read[addr=%#v] => %#v", addr, data)
		}
	}()

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		data = c.wram[addr]
		return data, nil
	} 

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		data = c.wram[addr - 0x0800]
		return data, nil
	} 

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		data = c.ppuRegister[addr - 0x2000]
		return data, nil
	} 

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		data = c.ppuRegister[addr - 0x2008]
		return data, nil
	} 

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		data = c.io[addr - 0x4000]
		return data, nil
	} 

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		data = c.exrom[addr - 0x4000]
		return data, nil
	} 

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		data = c.exram[addr - 0x6000]
		return data, nil
	} 

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		r := *c.programRom
		if len(r) <= 0x4000 {
			data = r[addr - 0xC000]
		} else {
			data = r[addr - 0x8000]
		}
		return data, nil
	} 

	return 0, fmt.Errorf("failed read, addr out of range; addr: %#v", addr)
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

	log.Printf("CPUBus.write[addr=%#v] <= %#v", addr, data)
}

// Run ...
func (c *CPU) Run() error {
	if err := c.interruptRESET(); err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Printf("CPU : %#v", c.String())
	for {
		log.Printf("===== cycle start =====")
		log.Printf("CPU : %#v", c.String())

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

		time.Sleep(time.Second * 1)
	}
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (*OpcodeProp, error) {
	if p, ok := OpcodeProps[o]; ok {
		log.Printf("decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Printf("decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (uint32, error) {
	var addr uint32
	var data uint32
	var err error

	defer func() {
		if err != nil {
			log.Printf("fetch[addr=%#v] => error %#v", addr, err)
		} else {
			log.Printf("fetch[addr=%#v] => %#v", addr, data)
		}
	}()

	// TODO 実装
	return data, nil
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
func (c *CPU) interruptRESET() error {
	log.Printf("interrupt[reset]")

	beforeI := c.registers.P.Interrupt
	c.registers.P.Interrupt = true
	log.Printf("reset[Interrupt flag] %#v => %#v", beforeI, c.registers.P.Interrupt)

	l, err := c.bus.read(0xFFFC)
	if err != nil {
		return fmt.Errorf("failed reset; %w", err)
	}

	h, err := c.bus.read(0xFFFD)
	if err != nil {
		return fmt.Errorf("failed reset; %w", err)
	}

	beforePC := c.registers.PC
	c.registers.PC = uint16((h << 4) | l)
	log.Printf("reset[PC] %#v => %#v", beforePC, c.registers.PC)

	return nil
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
	log.Printf("exec %#v %#v", mne, opr)
	// TODO 実装
	// （必要であれば）演算対象となるアドレスを算出
	return nil
}
