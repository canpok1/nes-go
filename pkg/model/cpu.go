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
	log.Trace("CPU.update[A] %#v => %#v", old, r.a)
}

// updateX ...
func (r *CPURegisters) updateX(x byte) {
	old := r.x
	r.x = x
	log.Trace("CPU.update[X] %#v => %#v", old, r.x)
}

// updateY ...
func (r *CPURegisters) updateY(y byte) {
	old := r.y
	r.y = y
	log.Trace("CPU.update[Y] %#v => %#v", old, r.y)
}

// updateS ...
func (r *CPURegisters) updateS(s byte) {
	old := r.s
	r.s = s
	log.Trace("CPU.update[S] %#v => %#v", old, r.s)
}

// incrementPC ...
func (r *CPURegisters) incrementPC() {
	r.updatePC(r.pc + 1)
}

// updatePC ...
func (r *CPURegisters) updatePC(pc uint16) {
	old := r.pc
	r.pc = pc
	log.Trace("CPU.update[PC] %#v => %#v", old, r.pc)
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

// ToByte ...
func (s *CPUStatusRegister) ToByte() byte {
	var b byte = 0
	if s.negative {
		b = b + 0x80
	}
	if s.overflow {
		b = b + 0x40
	}
	if s.reserved {
		b = b + 0x20
	}
	if s.breakMode {
		b = b + 0x10
	}
	if s.decimalMode {
		b = b + 0x08
	}
	if s.interrupt {
		b = b + 0x04
	}
	if s.zero {
		b = b + 0x02
	}
	if s.carry {
		b = b + 0x01
	}
	return b
}

// updateAll ...
func (s *CPUStatusRegister) updateAll(b byte) {
	s.negative = (b & 0x80) == 0x80
	s.overflow = (b & 0x40) == 0x40
	s.reserved = (b & 0x20) == 0x20
	s.breakMode = (b & 0x10) == 0x10
	s.decimalMode = (b & 0x08) == 0x08
	s.interrupt = (b & 0x04) == 0x04
	s.zero = (b & 0x02) == 0x02
	s.carry = (b & 0x01) == 0x01
}

// updateN ...
func (s *CPUStatusRegister) updateN(result byte) {
	old := s.negative
	s.negative = ((result & 0x80) == 0x80)
	log.Trace("CPU.update[N] %#v => %#v", old, s.negative)
}

// updateV
func (s *CPUStatusRegister) updateV(result int16) {
	old := s.overflow
	s.overflow = (result < 0x7F) || (result > 0x80)
	log.Trace("CPU.update[V] %#v => %#v", old, s.overflow)
}

// updateI ...
func (s *CPUStatusRegister) updateI(i bool) {
	old := s.interrupt
	s.interrupt = i
	log.Trace("CPU.update[I] %#v => %#v", old, s.interrupt)
}

// updateZ ...
func (s *CPUStatusRegister) updateZ(result byte) {
	old := s.zero
	s.zero = (result == 0x00)
	log.Trace("CPU.update[Z] %#v => %#v", old, s.zero)
}

// updateC ...
func (s *CPUStatusRegister) updateC(result int16) {
	old := s.carry
	s.carry = result > 0xFF
	log.Trace("CPU.update[C] %#v => %#v", old, s.carry)
}

// CPUStack ...
type CPUStack struct {
	stack []byte
}

// NewCPUStack ...
func NewCPUStack() *CPUStack {
	return &CPUStack{[]byte{}}
}

// Push ...
func (s *CPUStack) Push(b byte) {
	s.stack = append(s.stack, b)
}

// Pop ...
func (s *CPUStack) Pop() byte {
	b := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-2]
	return b
}

// CPU ...
type CPU struct {
	registers   *CPURegisters
	bus         *Bus
	shouldReset bool
	stack       *CPUStack
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{
		registers:   NewCPURegisters(),
		shouldReset: true,
		stack:       NewCPUStack(),
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
	log.Trace("===== CPU RUN =====")
	log.Trace(c.String())

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
	op, err := c.fetchAsOperand(ocp.AddressingMode)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	// 命令を実行
	if err := c.exec(ocp.Mnemonic, ocp.AddressingMode, op); err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return ocp.Cycle, nil
}

// decodeOpcode ...
func decodeOpcode(o Opcode) (*OpcodeProp, error) {
	if p, ok := OpcodeProps[o]; ok {
		log.Trace("CPU.decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Trace("CPU.decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr Address
	var data byte
	var err error

	log.Trace("CPU.fetch ...")
	defer func() {
		if err != nil {
			log.Warn("CPU.fetch[addr=%#v] => error %#v", addr, err)
		} else {
			log.Trace("CPU.fetch[addr=%#v] => %#v", addr, data)
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

// fetchAsOperand ...
func (c *CPU) fetchAsOperand(mode AddressingMode) (op []byte, err error) {
	log.Trace("CPU.fetchAsOperand[%#v] ...", mode)
	defer func() {
		if err != nil {
			log.Warn("CPU.fetchAsOperand[%#v] => %#v", mode, err)
		} else {
			log.Trace("CPU.fetchAsOperand[%#v] => %#v", mode, op)
		}
	}()

	switch mode {
	case Accumulator:
		fallthrough
	case Implied:
		op = []byte{}
		return
	case Immediate:
		fallthrough
	case ZeroPage:
		fallthrough
	case IndexedZeroPageX:
		fallthrough
	case IndexedZeroPageY:
		fallthrough
	case Relative:
		var b byte
		b, err = c.fetch()
		if err != nil {
			return
		}
		op = []byte{b}
		return
	case Absolute:
		fallthrough
	case IndexedAbsoluteX:
		fallthrough
	case IndexedAbsoluteY:
		fallthrough
	case IndexedIndirect:
		fallthrough
	case IndirectIndexed:
		fallthrough
	case AbsoluteIndirect:
		var l byte
		l, err = c.fetch()
		if err != nil {
			return
		}

		var h byte
		h, err = c.fetch()
		if err != nil {
			return
		}

		op = []byte{l, h}
		return
	default:
		err = fmt.Errorf("failed to fetch operands, AddressingMode is unsupported; mode: %#v", mode)
		return
	}
}

// makeAddress ...
func (c *CPU) makeAddress(mode AddressingMode, op []byte) (addr Address, err error) {
	log.Trace("CPU.makeAddress[%#v][%#v] ...", mode, op)
	defer func() {
		if err != nil {
			log.Warn("CPU.makeAddress[%#v][%#v] => %#v", mode, op, err)
		} else {
			log.Trace("CPU.makeAddress[%#v][%#v] => %#v", mode, op, addr)
		}
	}()

	switch mode {
	case Absolute:
		l := op[0]
		h := op[1]
		addr = Address((uint16(h) << 8) | uint16(l))
		return
	case ZeroPage:
		l := op[0]
		addr = Address(l)
		return
	case IndexedZeroPageX:
		l := op[0]
		addr = Address(uint8(l) + uint8(c.registers.x))
		return
	case IndexedZeroPageY:
		l := op[0]
		addr = Address(uint8(l) + uint8(c.registers.y))
		return
	case IndexedAbsoluteX:
		l := op[0]
		h := op[1]
		addr = Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.x))
		return
	case IndexedAbsoluteY:
		l := op[0]
		h := op[1]
		addr = Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.y))
		return
	case Relative:
		b := op[0]
		addr = Address(c.registers.pc + uint16(int8(b)))
		return
	case IndexedIndirect:
		b := op[0]
		dest := Address(uint8(b) + c.registers.x)

		var l byte
		l, err = c.bus.readByCPU(dest)
		if err != nil {
			return
		}

		h := op[1]

		addr = Address((uint16(h) << 8) | uint16(l))
		return
	case IndirectIndexed:
		b := op[0]
		dest := Address(uint8(b) + c.registers.x)

		var h byte
		h, err = c.bus.readByCPU(dest)
		if err != nil {
			return
		}

		l := op[1]

		addr = Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.y))
		return
	case AbsoluteIndirect:
		f1 := op[0]
		f2 := op[1]

		dest := Address((uint16(f2) << 8) + uint16(f1))
		nextDest := dest + 1

		var addrL byte
		addrL, err = c.bus.readByCPU(dest)
		if err != nil {
			return
		}

		var addrH byte
		addrH, err = c.bus.readByCPU(nextDest)
		if err != nil {
			return
		}

		addr = Address((uint16(addrH) << 8) + uint16(addrL))
		return
	default:
		err = fmt.Errorf("failed to make address, AddressingMode is not supported; mode: %#v", mode)
		return
	}
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
func (c *CPU) exec(mne Mnemonic, mode AddressingMode, op []byte) (err error) {
	log.Debug("CPU.exec[%v][%v][%#v] ...", mne, mode, op)

	defer func() {
		if err != nil {
			log.Warn("CPU.exec[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("CPU.exec[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

	switch mne {
	case ADC:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := int16(int8(c.registers.a)) + int16(int8(b))
		if c.registers.p.carry {
			ans = ans + 1
		}

		c.registers.a = byte(ans & 0xFF)
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateV(ans)
		c.registers.p.updateZ(c.registers.a)
		c.registers.p.updateC(ans)
		return
	case SBC:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := int16(int8(c.registers.a)) - int16(int8(b))
		if !c.registers.p.carry {
			ans = ans - 1
		}

		c.registers.a = byte(ans & 0xFF)
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateV(ans)
		c.registers.p.updateZ(c.registers.a)
		c.registers.p.updateC(ans)
		return
	case AND:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = c.registers.a & b
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case ORA:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = c.registers.a | b
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case EOR:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = c.registers.a ^ b
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case ASL:
		var b byte
		if mode == Accumulator {
			b = c.registers.a
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = b << 1
		c.registers.p.updateN(c.registers.a)
		c.registers.p.carry = (b & 0x80) == 0x80
		return
	case LSR:
		var b byte
		if mode == Accumulator {
			b = c.registers.a
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = b >> 1
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		c.registers.p.carry = (b & 0x01) == 0x01
		return
	case ROL:
		var b byte
		if mode == Accumulator {
			b = c.registers.a
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = b << 1
		if c.registers.p.carry {
			c.registers.a = c.registers.a + 1
		}

		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		c.registers.p.carry = (b & 0x80) == 0x80
		return
	case ROR:
		var b byte
		if mode == Accumulator {
			b = c.registers.a
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.a = b >> 1
		if c.registers.p.carry {
			c.registers.a = c.registers.a + 0x80
		}

		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		c.registers.p.carry = (b & 0x01) == 0x01
		return
	case BCC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.p.carry {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BCS:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.p.carry {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BEQ:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.p.zero {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BNE:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.p.zero {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BVC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.p.carry {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BVS:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.p.carry {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BPL:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.p.negative {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BMI:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.p.negative {
			c.registers.updatePC(uint16(addr))
		}
		return
	case BIT:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.readByCPU(addr)
		if err != nil {
			return
		}

		c.registers.p.zero = (c.registers.a & b) == 0
		c.registers.p.negative = (b & 0x80) == 0x80
		c.registers.p.overflow = (b & 0x40) == 0x40
		return
	case JMP:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		c.registers.updatePC(uint16(addr))
		return
	case JSR:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.readByCPU(addr)
		if err != nil {
			return
		}

		c.stack.Push(byte((c.registers.pc & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.pc & 0x00FF))
		c.registers.pc = uint16(b)
		return
	case RTS:
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.pc = (uint16(h) << 8) + uint16(l) + 1
		return
	case BRK:
		c.registers.p.breakMode = true
		c.registers.incrementPC()
		c.stack.Push(byte((c.registers.pc & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.pc & 0x00FF))
		c.stack.Push(c.registers.p.ToByte())
		c.registers.p.interrupt = true

		var l, h byte
		l, err = c.bus.readByCPU(0xFFFE)
		if err != nil {
			return
		}
		h, err = c.bus.readByCPU(0xFFFF)
		if err != nil {
			return
		}
		c.registers.pc = (uint16(h) << 8) + uint16(l)
		return
	case RTI:
		c.registers.p.updateAll(c.stack.Pop())
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.pc = (uint16(h) << 8) + uint16(l)
		return
	case CMP:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.a - b
		c.registers.p.updateN(ans)
		c.registers.p.updateZ(ans)
		c.registers.p.carry = ans >= 0
		return
	case CPX:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.x - b
		c.registers.p.updateN(ans)
		c.registers.p.updateZ(ans)
		c.registers.p.carry = ans >= 0
		return
	case CPY:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.y - b
		c.registers.p.updateN(ans)
		c.registers.p.updateZ(ans)
		c.registers.p.carry = ans >= 0
		return
	case INC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.readByCPU(addr)
		if err != nil {
			return
		}

		ans := b + 1
		err = c.bus.writeByCPU(addr, ans)
		if err != nil {
			return
		}

		c.registers.p.updateN(ans)
		c.registers.p.updateZ(ans)
		return
	case DEC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.readByCPU(addr)
		if err != nil {
			return
		}

		ans := b - 1
		err = c.bus.writeByCPU(addr, ans)
		if err != nil {
			return
		}

		c.registers.p.updateN(ans)
		c.registers.p.updateZ(ans)
		return
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
	case CLC:
		c.registers.p.carry = false
		return
	case SEC:
		c.registers.p.carry = true
		return
	case CLI:
		c.registers.p.updateI(false)
		return
	case SEI:
		c.registers.p.updateI(true)
		return
	case CLD:
		c.registers.p.decimalMode = false
		return
	case SED:
		c.registers.p.decimalMode = true
		return
	case CLV:
		c.registers.p.overflow = false
		return
	case LDA:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.updateA(b)
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case LDX:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.updateX(b)
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case LDY:
		var b byte
		if mode == Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.readByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.updateY(b)
		c.registers.p.updateN(c.registers.y)
		c.registers.p.updateZ(c.registers.y)
		return
	case STA:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.writeByCPU(addr, c.registers.a)
		if err != nil {
			return
		}
		return
	case STX:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.writeByCPU(addr, c.registers.x)
		if err != nil {
			return
		}
		return
	case STY:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.writeByCPU(addr, c.registers.y)
		if err != nil {
			return
		}
		return
	case TAX:
		c.registers.x = c.registers.a
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case TXA:
		c.registers.a = c.registers.x
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case TAY:
		c.registers.y = c.registers.a
		c.registers.p.updateN(c.registers.y)
		c.registers.p.updateZ(c.registers.y)
		return
	case TYA:
		c.registers.a = c.registers.y
		c.registers.p.updateN(c.registers.a)
		c.registers.p.updateZ(c.registers.a)
		return
	case TSX:
		c.registers.x = c.registers.s
		c.registers.p.updateN(c.registers.x)
		c.registers.p.updateZ(c.registers.x)
		return
	case TXS:
		c.registers.updateS(c.registers.x)
		return
	case PHA:
		c.stack.Push(c.registers.a)
		return
	case PLA:
		c.registers.a = c.stack.Pop()
		return
	case PHP:
		c.stack.Push(c.registers.p.ToByte())
		return
	case PLP:
		c.registers.p.updateAll(c.stack.Pop())
		return
	case NOP:
		return
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}
