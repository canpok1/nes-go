package model

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
)

// CPURegisters ...
type CPURegisters struct {
	a  byte
	x  byte
	y  byte
	s  byte
	p  *CPUStatusRegister
	pc uint16
}

// NewCPURegisters ...
func NewCPURegisters() *CPURegisters {
	// initialize as CPU power up state
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state
	return &CPURegisters{
		a:  0,
		x:  0,
		y:  0,
		s:  0xFD,
		p:  NewCPUStatusRegister(),
		pc: 0,
	}
}

// String ...
func (r *CPURegisters) String() string {
	return fmt.Sprintf(
		"{A:%#v, X:%#v, Y:%#v, S:%#v, P:%v, PC:%#v}",
		r.a,
		r.x,
		r.y,
		r.s,
		r.p.String(),
		r.pc,
	)
}

// updateA ...
func (r *CPURegisters) updateA(a byte) {
	old := r.a
	r.a = a
	log.Debug("CPU.update[A] %#v => %#v", old, r.a)
}

// updateX ...
func (r *CPURegisters) updateX(x byte) {
	old := r.x
	r.x = x
	log.Debug("CPU.update[X] %#v => %#v", old, r.x)
}

// updateY ...
func (r *CPURegisters) updateY(y byte) {
	old := r.y
	r.y = y
	log.Debug("CPU.update[Y] %#v => %#v", old, r.y)
}

// updateS ...
func (r *CPURegisters) updateS(s byte) {
	old := r.s
	r.s = s
	log.Debug("CPU.update[S] %#v => %#v", old, r.s)
}

// incrementPC ...
func (r *CPURegisters) incrementPC() {
	r.updatePC(r.pc + 1)
}

// updatePC ...
func (r *CPURegisters) updatePC(pc uint16) {
	old := r.pc
	r.pc = pc
	log.Debug("CPU.update[PC] %#v => %#v", old, r.pc)
}

// CPUStatusRegister ...
// https://qiita.com/bokuweb/items/1575337bef44ae82f4d3#%E3%82%B9%E3%83%86%E3%83%BC%E3%82%BF%E3%82%B9%E3%83%AC%E3%82%B8%E3%82%B9%E3%82%BF
type CPUStatusRegister struct {
	negative    bool // bit7	N	ネガティブ	演算結果のbit7が1の時にセット
	overflow    bool // bit6	V	オーバーフロー	P演算結果がオーバーフローを起こした時にセット
	reserved    bool // bit5	R	予約済み	常にセットされている
	breakMode   bool // bit4	B	ブレークモード	BRK発生時にセット、IRQ発生時にクリア
	decimalMode bool // bit3	D	デシマルモード	0:デフォルト、1:BCDモード (未実装)
	interrupt   bool // bit2	I	IRQ禁止	0:IRQ許可、1:IRQ禁止
	zero        bool // bit1	Z	ゼロ	演算結果が0の時にセット
	carry       bool // bit0	C	キャリー	キャリー発生時にセット
}

// NewCPUStatusRegister ...
func NewCPUStatusRegister() *CPUStatusRegister {
	return &CPUStatusRegister{
		negative:    false,
		overflow:    false,
		reserved:    true,
		breakMode:   true,
		decimalMode: false,
		interrupt:   true,
		zero:        false,
		carry:       false,
	}
}

// String ...
func (s *CPUStatusRegister) String() string {
	return fmt.Sprintf(
		"{N:%#v, V:%#v, R:%#v, B:%#v, D:%#v, I:%#v, Z:%#v, C:%#v}",
		s.negative,
		s.overflow,
		s.reserved,
		s.breakMode,
		s.decimalMode,
		s.interrupt,
		s.zero,
		s.carry,
	)
}

// updateN ...
func (s *CPUStatusRegister) updateN(result byte) {
	old := s.negative
	s.negative = ((result & 0x80) == 0x80)
	log.Debug("CPU.update[N] %#v => %#v", old, s.negative)
}

// updateI ...
func (s *CPUStatusRegister) updateI(i bool) {
	old := s.interrupt
	s.interrupt = i
	log.Debug("CPU.update[I] %#v => %#v", old, s.interrupt)
}

// updateZ ...
func (s *CPUStatusRegister) updateZ(result byte) {
	old := s.zero
	s.zero = (result == 0x00)
	log.Debug("CPU.update[Z] %#v => %#v", old, s.zero)
}

// CPU ...
type CPU struct {
	registers   *CPURegisters
	bus         *Bus
	shouldReset bool
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{
		registers:   NewCPURegisters(),
		shouldReset: true,
	}
}

// String ...
func (c *CPU) String() string {
	return fmt.Sprintf(
		"CPU Info\nregisters: %v\nshould reset: %v",
		c.registers.String(),
		c.shouldReset,
	)
}

// SetBus ...
func (c *CPU) SetBus(b *Bus) {
	c.bus = b
}

// Run ...
func (c *CPU) Run() (int, error) {
	log.Debug("===== CPU RUN =====")
	log.Debug(c.String())

	if c.bus == nil {
		return 0, fmt.Errorf("failed to run, bus is nil")
	}

	if c.shouldReset {
		if err := c.interruptRESET(); err != nil {
			return 0, fmt.Errorf("%w", err)
		}
		return 0, nil
	}

	// PC（プログラムカウンタ）からオペコードをフェッチ（PCをインクリメント）
	oc, err := c.fetchOpcode()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	// 命令とアドレッシング・モードを判別
	ocp, err := decodeOpcode(oc)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	// （必要であれば）オペランドをフェッチ（PCをインクリメント）
	addr, err := c.makeAddress(ocp.AddressingMode)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	// 命令を実行
	if err := c.exec(ocp.Mnemonic, addr); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return ocp.Cycle, nil
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (*OpcodeProp, error) {
	if p, ok := OpcodeProps[o]; ok {
		log.Debug("CPU.decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Debug("CPU.decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr Address
	var data byte
	var err error

	log.Debug("CPU.fetch ...")
	defer func() {
		if err != nil {
			log.Debug("CPU.fetch[addr=%#v] => error %#v", addr, err)
		} else {
			log.Debug("CPU.fetch[addr=%#v] => %#v", addr, data)
		}
	}()

	addr = Address(c.registers.pc)
	data, err = c.bus.readByCPU(addr)
	if err != nil {
		return data, fmt.Errorf("failed to fetch; %w", err)
	}

	c.registers.incrementPC()

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

		addr := Address(uint8(l) + uint8(c.registers.x))
		return &addr, nil
	case IndexedZeroPageY:
		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(uint8(l) + uint8(c.registers.y))
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

		addr := Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.x))

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

		addr := Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.y))
		return &addr, nil
	case Implied:
		return nil, nil
	case Relative:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address(c.registers.pc + uint16(int8(b)))

		return &addr, nil
	case IndexedIndirect:
		b, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}
		dest := Address(uint8(b) + c.registers.x)

		l, err := c.bus.readByCPU(dest)
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
		dest := Address(uint8(b) + c.registers.x)

		h, err := c.bus.readByCPU(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		l, err := c.fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.y))

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

		addrL, err := c.bus.readByCPU(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addrH, err := c.bus.readByCPU(nextDest)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch operands; %w", err)
		}

		addr := Address((uint16(addrH) << 8) + uint16(addrL))
		return &addr, nil
	}

	return nil, fmt.Errorf("failed to fetch operands, AddressingMode is unsupported; mode: %#v", mode)
}

// interruptNMI ...
func (c *CPU) interruptNMI() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// interruptRESET ...
func (c *CPU) interruptRESET() error {
	log.Info("CPU.interrupt[RESET] ...")

	c.registers.p.updateI(true)

	l, err := c.bus.readByCPU(0xFFFC)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	h, err := c.bus.readByCPU(0xFFFD)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	c.registers.updatePC((uint16(h) << 8) | uint16(l))

	c.shouldReset = false
	return nil
}

// interruptBRK ...
func (c *CPU) interruptBRK() {
	log.Info("CPU.interrupt[BRK] ...")
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// interruptIRQ ...
func (c *CPU) interruptIRQ() {
	log.Info("CPU.interrupt[IRQ] ...")
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#interrupt
}

// exec ...
func (c *CPU) exec(mne Mnemonic, addr *Address) (err error) {
	log.Debug("CPU.exec[%#v][%#v] ...", mne, addr)
	defer func() {
		addrStr := ""
		if addr != nil {
			addrStr = fmt.Sprintf("%#v", *addr)
		}

		if err != nil {
			log.Warn("CPU.exec[%v][%#v] => %v", mne, addrStr, err)
		} else {
			log.Info("CPU.exec[%v][%#v] => completed", mne, addrStr)
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
			return
		}
		if !c.registers.p.zero {
			c.registers.updatePC(uint16(*addr))
		}
		return
	// TODO 実装 BVC Mnemonic = "BVC"
	// TODO 実装 BVS Mnemonic = "BVS"
	// TODO 実装 BPL Mnemonic = "BPL"
	// TODO 実装 BMI Mnemonic = "BMI"
	// TODO 実装 BIT Mnemonic = "BIT"
	case JMP:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			return
		}
		c.registers.updatePC(uint16(*addr))
		return
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
		c.registers.updateX(c.registers.x + 1)
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case DEX:
		c.registers.updateX(c.registers.x - 1)
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case INY:
		c.registers.updateY(c.registers.y + 1)
		c.registers.p.updateN(c.registers.y)
		c.registers.p.updateZ(c.registers.y)
		return
	case DEY:
		c.registers.updateY(c.registers.y - 1)
		c.registers.p.updateN(c.registers.y)
		c.registers.p.updateZ(c.registers.y)
		return
	// TODO 実装 CLC Mnemonic = "CLC"
	// TODO 実装 SEC Mnemonic = "SEC"
	// TODO 実装 CLI Mnemonic = "CLI"
	case SEI:
		c.registers.p.updateI(true)
		return
	// TODO 実装 CLD Mnemonic = "CLD"
	// TODO 実装 SED Mnemonic = "SED"
	// TODO 実装 CLV Mnemonic = "CLV"
	case LDA:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			return
		}

		var b byte
		b, err = c.bus.readByCPU(*addr)
		if err != nil {
			return
		}

		c.registers.updateA(b)
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case LDX:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			return
		}

		var b byte
		b, err = c.bus.readByCPU(*addr)
		if err != nil {
			return
		}

		c.registers.updateX(b)
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case LDY:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			return
		}

		var b byte
		b, err = c.bus.readByCPU(*addr)
		if err != nil {
			return
		}

		c.registers.updateY(b)
		c.registers.p.updateN(c.registers.y)
		c.registers.p.updateZ(c.registers.y)
		return
	case STA:
		if addr == nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, address: %#v", mne, addr)
			return
		}

		err = c.bus.writeByCPU(*addr, c.registers.a)
		if err != nil {
			return
		}
		return
	// TODO 実装 STX Mnemonic = "STX"
	// TODO 実装 STY Mnemonic = "STY"
	// TODO 実装 TAX Mnemonic = "TAX"
	// TODO 実装 TXA Mnemonic = "TXA"
	// TODO 実装 TAY Mnemonic = "TAY"
	// TODO 実装 TYA Mnemonic = "TYA"
	// TODO 実装 TSX Mnemonic = "TSX"
	case TXS:
		c.registers.updateS(c.registers.x)
		return
	// TODO 実装 PHA Mnemonic = "PHA"
	// TODO 実装 PLA Mnemonic = "PLA"
	// TODO 実装 PHP Mnemonic = "PHP"
	// TODO 実装 PLP Mnemonic = "PLP"
	// TODO 実装 NOP Mnemonic = "NOP"
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}
