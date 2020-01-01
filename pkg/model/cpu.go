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
func (c *CPUBus) read(addr Address) (byte, error) {
	var data byte 
	var err error
	var target string
	defer func() {
		if err != nil {
			log.Printf("CPUBus.read[addr=%#v] => %#v", addr, err)
		} else {
			log.Printf("CPUBus.read[addr=%#v][%v] => %#v", addr, target, data)
		}
	}()

	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		data = c.wram[addr]
		return data, nil
	} 

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		data = c.wram[addr - 0x0800]
		return data, nil
	} 

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		data = c.ppuRegister[addr - 0x2000]
		return data, nil
	} 

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		data = c.ppuRegister[addr - 0x2008]
		return data, nil
	} 

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		data = c.io[addr - 0x4000]
		return data, nil
	} 

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		target = "EX ROM"
		data = c.exrom[addr - 0x4000]
		return data, nil
	} 

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		data = c.exram[addr - 0x6000]
		return data, nil
	} 

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		target = "PRG-ROM"
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
func (c *CPUBus) write(addr uint32, data byte) error {
	var err error
	var target string
	defer func() {
		if err != nil {
			log.Printf("CPUBus.write[addr=%#v] => %#v", addr, err)
		} else {
			log.Printf("CPUBus.write[addr=%#v][%v] <= %#v", addr, target, data)
		}
	}()
	// 0x0000～0x07FF	0x0800	WRAM
	if addr >= 0x0000 && addr <= 0x07FF {
		target = "WRAM"
		c.wram[addr] = data
		return nil
	} 

	// 0x0800～0x1FFF	-	WRAMのミラー
	if addr >= 0x0800 && addr <= 0x1FFF {
		target = "WRAM Mirror"
		c.wram[addr - 0x0800] = data
		return nil
	} 

	// 0x2000～0x2007	0x0008	PPU レジスタ
	if addr >= 0x2000 && addr <= 0x2007 {
		target = "PPU Register"
		c.ppuRegister[addr - 0x2000] = data
		return nil
	} 

	// 0x2008～0x3FFF	-	PPUレジスタのミラー
	if addr >= 0x2008 && addr <= 0x3FFF {
		target = "PPU Register Mirror"
		c.ppuRegister[addr - 0x2008] = data
		return nil
	} 

	// 0x4000～0x401F	0x0020	APU I/O、PAD
	if addr >= 0x4000 && addr <= 0x401F {
		target = "APU I/O, PAD"
		c.io[addr - 0x4000] = data
		return nil
	} 

	// 0x4020～0x5FFF	0x1FE0	拡張ROM
	if addr >= 0x4000 && addr <= 0x401F {
		return fmt.Errorf("failed write, cannot write EX ROM; addr: %#v", addr)
	} 

	// 0x6000～0x7FFF	0x2000	拡張RAM
	if addr >= 0x6000 && addr <= 0x7FFF {
		target = "EX RAM"
		c.exram[addr - 0x6000] = data
		return nil
	} 

	// 0x8000～0xBFFF	0x4000	PRG-ROM
	// 0xC000～0xFFFF	0x4000	PRG-ROM
	if addr >= 0x8000 && addr <= 0xFFFF {
		return fmt.Errorf("failed write, cannot write PRG-ROM; addr: %#v", addr)
	} 

	return fmt.Errorf("failed write, addr out of range; addr: %#v", addr)
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
func (c *CPU) fetch() (byte, error) {
	var addr Address
	var data byte
	var err error

	defer func() {
		if err != nil {
			log.Printf("fetch[addr=%#v] => error %#v", addr, err)
		} else {
			log.Printf("fetch[addr=%#v] => %#v", addr, data)
		}
	}()

	addr = Address(c.registers.PC)
	data, err = c.bus.read(addr)
	if err != nil {
		return data, fmt.Errorf("failed fetch; %w", err)
	}

	c.registers.PC = c.registers.PC + 1

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
	switch mode {
	case Accumulator:
		return []Operand{}, nil
	case Immediate:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		return []Operand{Operand(b)}, nil
	case Absolute:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand((uint16(h) << 8) | uint16(l))

		return []Operand{op}, nil
	case ZeroPage:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(l)

		return []Operand{op}, nil
	case IndexedZeroPageX:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(uint8(l) + uint8(c.registers.X))

		return []Operand{op}, nil
	case IndexedZeroPageY:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(uint8(l) + uint8(c.registers.Y))

		return []Operand{op}, nil
	case IndexedAbsoluteX:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.X))

		return []Operand{op}, nil
	case IndexedAbsoluteY:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.Y))

		return []Operand{op}, nil
	case Implied:
		return []Operand{}, nil
	case Relative:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand(c.registers.PC + uint16(int8(b)))

		return []Operand{op}, nil
	case IndexedIndirect:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		addr := Address(uint8(b) + c.registers.X)

		l, err := c.bus.read(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand((uint16(h) << 8) | uint16(l))

		return []Operand{op}, nil
	case IndirectIndexed:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		addr := Address(uint8(b) + c.registers.X)

		h, err := c.bus.read(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand((uint16(h) << 8) + uint16(l) + uint16(c.registers.Y))

		return []Operand{op}, nil
	case AbsoluteIndirect:
		f1, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		f2, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(f2) << 8) + uint16(f1))
		nextAddr := addr + 1

		opl, err := c.bus.read(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		oph, err := c.bus.read(nextAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		op := Operand((uint16(oph) << 8) + uint16(opl))

		return []Operand{op}, nil
	}

	return nil, fmt.Errorf("failed fetch operands, AddressingMode is unsupported; mode: %#v", mode)
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
	c.registers.PC = (uint16(h) << 8) | uint16(l)
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
