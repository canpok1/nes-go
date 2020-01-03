package model

import (
	"fmt"
	"time"

	"github.com/canpok1/nes-go/pkg/log"
)

// Registers ...
type Registers struct {
	A  byte
	X  byte
	Y  byte
	S  byte
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

// updateZ ...
func (s *StatusRegister) updateZ(result byte) {
	s.Zero = (result == 0x00)
}

// updateN ...
func (s *StatusRegister) updateN(result byte) {
	s.Zero = ((result & 0x80) == 0x80)
}

// Run ...
func (c *CPU) Run() error {
	if err := c.interruptRESET(); err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Debug("CPU : %#v", c.String())
	for {
		log.Debug("===== cycle start =====")
		log.Debug("CPU : %#v", c.String())

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
		addr, err := c.makeAddress(ocp.AddressingMode)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// 命令を実行
		if err := c.exec(ocp.Mnemonic, addr); err != nil {
			return fmt.Errorf("%w", err)
		}

		time.Sleep(time.Second * 1)
	}
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (*OpcodeProp, error) {
	if p, ok := OpcodeProps[o]; ok {
		log.Debug("decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Debug("decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr Address
	var data byte
	var err error

	defer func() {
		if err != nil {
			log.Debug("fetch[addr=%#v] => error %#v", addr, err)
		} else {
			log.Debug("fetch[addr=%#v] => %#v", addr, data)
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

// makeAddress ...
func (c *CPU) makeAddress(mode AddressingMode) (*Address, error) {
	switch mode {
	case Accumulator:
		return nil, nil
	case Immediate:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		addr := Address(b)
		return &addr, nil
	case Absolute:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(h) << 8) | uint16(l))
		return &addr, nil
	case ZeroPage:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(l)
		return &addr, nil
	case IndexedZeroPageX:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(uint8(l) + uint8(c.registers.X))
		return &addr, nil
	case IndexedZeroPageY:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(uint8(l) + uint8(c.registers.Y))
		return &addr, nil
	case IndexedAbsoluteX:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.X))

		return &addr, nil
	case IndexedAbsoluteY:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.Y))
		return &addr, nil
	case Implied:
		return nil, nil
	case Relative:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(c.registers.PC + uint16(int8(b)))

		return &addr, nil
	case IndexedIndirect:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		dest := Address(uint8(b) + c.registers.X)

		l, err := c.bus.read(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		h, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(h) << 8) | uint16(l))
		return &addr, nil
	case IndirectIndexed:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		dest := Address(uint8(b) + c.registers.X)

		h, err := c.bus.read(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.Y))

		return &addr, nil
	case AbsoluteIndirect:
		f1, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		f2, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		dest := Address((uint16(f2) << 8) + uint16(f1))
		nextDest := dest + 1

		addrL, err := c.bus.read(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addrH, err := c.bus.read(nextDest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(addrH) << 8) + uint16(addrL))
		return &addr, nil
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
	log.Info("interrupt[reset]")

	beforeI := c.registers.P.Interrupt
	c.registers.P.Interrupt = true
	log.Debug("reset[Interrupt flag] %#v => %#v", beforeI, c.registers.P.Interrupt)

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
	log.Debug("reset[PC] %#v => %#v", beforePC, c.registers.PC)

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
func (c *CPU) exec(mne Mnemonic, addr *Address) error {
	var err error

	defer func() {
		if err != nil {
			log.Warn("exec %#v %#v => failed", mne, addr)
		} else if addr == nil {
			log.Info("exec %#v %#v => completed", mne, addr)
		} else {
			log.Info("exec %#v %#v => completed", mne, *addr)
		}
	}()

	switch mne {
	// TODO 実装 ADC Mnemonic = "ADC"
	// TODO 実装 SBC Mnemonic = "SBC"
	// TODO 実装 AND Mnemonic = "AND"
	// TODO 実装 ORA Mnemonic = "ORA"
	// TODO 実装 EOR Mnemonic = "EOR"
	// TODO 実装 ASL Mnemonic = "ASL"
	// TODO 実装 LSR Mnemonic = "LSR"
	// TODO 実装 ROL Mnemonic = "ROL"
	// TODO 実装 ROR Mnemonic = "ROR"
	// TODO 実装 BCC Mnemonic = "BCC"
	// TODO 実装 BCS Mnemonic = "BCS"
	// TODO 実装 BEQ Mnemonic = "BEQ"
	// TODO 実装 BNE Mnemonic = "BNE"
	case BNE:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}
		if !c.registers.P.Zero {
			c.registers.PC = uint16(*addr)
		}
	// TODO 実装 BVC Mnemonic = "BVC"
	// TODO 実装 BVS Mnemonic = "BVS"
	// TODO 実装 BPL Mnemonic = "BPL"
	// TODO 実装 BMI Mnemonic = "BMI"
	// TODO 実装 BIT Mnemonic = "BIT"
	case JMP:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}
		c.registers.PC = uint16(*addr)
	// TODO 実装 JSR Mnemonic = "JSR"
	// TODO 実装 RTS Mnemonic = "RTS"
	// TODO 実装 BRK Mnemonic = "BRK"
	// TODO 実装 RTI Mnemonic = "RTI"
	// TODO 実装 CMP Mnemonic = "CMP"
	// TODO 実装 CPX Mnemonic = "CPX"
	// TODO 実装 CPY Mnemonic = "CPY"
	// TODO 実装 INC Mnemonic = "INC"
	// TODO 実装 DEC Mnemonic = "DEC"
	case INX:
		c.registers.X = c.registers.X + 1
		c.registers.P.updateN(c.registers.X)
		c.registers.P.updateZ(c.registers.X)
	case DEX:
		c.registers.X = c.registers.X - 1
		c.registers.P.updateN(c.registers.X)
		c.registers.P.updateZ(c.registers.X)
	case INY:
		c.registers.Y = c.registers.Y + 1
		c.registers.P.updateN(c.registers.Y)
		c.registers.P.updateZ(c.registers.Y)
	case DEY:
		c.registers.Y = c.registers.Y - 1
		c.registers.P.updateN(c.registers.Y)
		c.registers.P.updateZ(c.registers.Y)
	// TODO 実装 CLC Mnemonic = "CLC"
	// TODO 実装 SEC Mnemonic = "SEC"
	// TODO 実装 CLI Mnemonic = "CLI"
	case SEI:
		c.registers.P.Interrupt = true
	// TODO 実装 CLD Mnemonic = "CLD"
	// TODO 実装 SED Mnemonic = "SED"
	// TODO 実装 CLV Mnemonic = "CLV"
	case LDA:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}

		var b byte
		b, err = c.bus.read(*addr)
		if err != nil {
			break
		}

		c.registers.A = b
		c.registers.P.updateN(c.registers.A)
		c.registers.P.updateZ(c.registers.A)
	case LDX:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}

		var b byte
		b, err = c.bus.read(*addr)
		if err != nil {
			break
		}

		c.registers.X = b
		c.registers.P.updateN(c.registers.X)
		c.registers.P.updateZ(c.registers.X)
	case LDY:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}

		var b byte
		b, err = c.bus.read(*addr)
		if err != nil {
			break
		}

		c.registers.Y = b
		c.registers.P.updateN(c.registers.Y)
		c.registers.P.updateZ(c.registers.Y)
	case STA:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			break
		}

		err = c.bus.write(*addr, c.registers.A)
		if err != nil {
			break
		}
	// TODO 実装 STX Mnemonic = "STX"
	// TODO 実装 STY Mnemonic = "STY"
	// TODO 実装 TAX Mnemonic = "TAX"
	// TODO 実装 TXA Mnemonic = "TXA"
	// TODO 実装 TAY Mnemonic = "TAY"
	// TODO 実装 TYA Mnemonic = "TYA"
	// TODO 実装 TSX Mnemonic = "TSX"
	case TXS:
		c.registers.S = c.registers.X
	// TODO 実装 PHA Mnemonic = "PHA"
	// TODO 実装 PLA Mnemonic = "PLA"
	// TODO 実装 PHP Mnemonic = "PHP"
	// TODO 実装 PLP Mnemonic = "PLP"
	// TODO 実装 NOP Mnemonic = "NOP"
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
	}

	return err
}
