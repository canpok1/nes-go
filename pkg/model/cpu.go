package model

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
	"github.com/canpok1/nes-go/pkg/model/cpu"
)

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
	registers   *cpu.CPURegisters
	bus         *Bus
	shouldReset bool
	stack       *CPUStack
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{
		registers:   cpu.NewCPURegisters(),
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
		if err := c.InterruptRESET(); err != nil {
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

	addr = Address(c.registers.PC)
	data, err = c.bus.readByCPU(addr)
	if err != nil {
		return data, fmt.Errorf("failed to fetch; %w", err)
	}

	c.registers.IncrementPC()

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
		addr = Address(uint8(l) + uint8(c.registers.X))
		return
	case IndexedZeroPageY:
		l := op[0]
		addr = Address(uint8(l) + uint8(c.registers.Y))
		return
	case IndexedAbsoluteX:
		l := op[0]
		h := op[1]
		addr = Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.X))
		return
	case IndexedAbsoluteY:
		l := op[0]
		h := op[1]
		addr = Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.Y))
		return
	case Relative:
		b := op[0]
		addr = Address(c.registers.PC + uint16(int8(b)))
		return
	case IndexedIndirect:
		b := op[0]
		dest := Address(uint8(b) + c.registers.X)

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
		dest := Address(uint8(b) + c.registers.X)

		var h byte
		h, err = c.bus.readByCPU(dest)
		if err != nil {
			return
		}

		l := op[1]

		addr = Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.Y))
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

// InterruptNMI ...
func (c *CPU) InterruptNMI() {
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#Interrupt
}

// InterruptRESET ...
func (c *CPU) InterruptRESET() error {
	log.Info("CPU.Interrupt[RESET] ...")

	c.registers.P.UpdateI(true)

	l, err := c.bus.readByCPU(0xFFFC)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	h, err := c.bus.readByCPU(0xFFFD)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	c.registers.UpdatePC((uint16(h) << 8) | uint16(l))

	c.shouldReset = false
	return nil
}

// InterruptBRK ...
func (c *CPU) InterruptBRK() {
	log.Info("CPU.Interrupt[BRK] ...")
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#Interrupt
}

// InterruptIRQ ...
func (c *CPU) InterruptIRQ() {
	log.Info("CPU.Interrupt[IRQ] ...")
	// TODO 実装
	// http://pgate1.at-ninja.jp/NES_on_FPGA/nes_cpu.htm#Interrupt
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
		ans := int16(int8(c.registers.A)) + int16(int8(b))
		if c.registers.P.Carry {
			ans = ans + 1
		}

		c.registers.A = byte(ans & 0xFF)
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateV(ans)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.UpdateC(ans)
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
		ans := int16(int8(c.registers.A)) - int16(int8(b))
		if !c.registers.P.Carry {
			ans = ans - 1
		}

		c.registers.A = byte(ans & 0xFF)
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateV(ans)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.UpdateC(ans)
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
		c.registers.A = c.registers.A & b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
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
		c.registers.A = c.registers.A | b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
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
		c.registers.A = c.registers.A ^ b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case ASL:
		var b byte
		if mode == Accumulator {
			b = c.registers.A
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
		c.registers.A = b << 1
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case LSR:
		var b byte
		if mode == Accumulator {
			b = c.registers.A
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
		c.registers.A = b >> 1
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case ROL:
		var b byte
		if mode == Accumulator {
			b = c.registers.A
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
		c.registers.A = b << 1
		if c.registers.P.Carry {
			c.registers.A = c.registers.A + 1
		}

		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case ROR:
		var b byte
		if mode == Accumulator {
			b = c.registers.A
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
		c.registers.A = b >> 1
		if c.registers.P.Carry {
			c.registers.A = c.registers.A + 0x80
		}

		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case BCC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BCS:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BEQ:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BNE:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BVC:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BVS:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BPL:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case BMI:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
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

		c.registers.P.Zero = (c.registers.A & b) == 0
		c.registers.P.Negative = (b & 0x80) == 0x80
		c.registers.P.Overflow = (b & 0x40) == 0x40
		return
	case JMP:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		c.registers.UpdatePC(uint16(addr))
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

		c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.PC & 0x00FF))
		c.registers.PC = uint16(b)
		return
	case RTS:
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l) + 1
		return
	case BRK:
		c.registers.P.BreakMode = true
		c.registers.IncrementPC()
		c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.PC & 0x00FF))
		c.stack.Push(c.registers.P.ToByte())
		c.registers.P.Interrupt = true

		var l, h byte
		l, err = c.bus.readByCPU(0xFFFE)
		if err != nil {
			return
		}
		h, err = c.bus.readByCPU(0xFFFF)
		if err != nil {
			return
		}
		c.registers.PC = (uint16(h) << 8) + uint16(l)
		return
	case RTI:
		c.registers.P.UpdateAll(c.stack.Pop())
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l)
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
		ans := c.registers.A - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
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
		ans := c.registers.X - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
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
		ans := c.registers.Y - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
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

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
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

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		return
	case INX:
		c.registers.UpdateX(c.registers.X + 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case DEX:
		c.registers.UpdateX(c.registers.X - 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case INY:
		c.registers.UpdateY(c.registers.Y + 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case DEY:
		c.registers.UpdateY(c.registers.Y - 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case CLC:
		c.registers.P.Carry = false
		return
	case SEC:
		c.registers.P.Carry = true
		return
	case CLI:
		c.registers.P.UpdateI(false)
		return
	case SEI:
		c.registers.P.UpdateI(true)
		return
	case CLD:
		c.registers.P.DecimalMode = false
		return
	case SED:
		c.registers.P.DecimalMode = true
		return
	case CLV:
		c.registers.P.Overflow = false
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

		c.registers.UpdateA(b)
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
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

		c.registers.UpdateX(b)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
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

		c.registers.UpdateY(b)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case STA:
		var addr Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.writeByCPU(addr, c.registers.A)
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

		err = c.bus.writeByCPU(addr, c.registers.X)
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

		err = c.bus.writeByCPU(addr, c.registers.Y)
		if err != nil {
			return
		}
		return
	case TAX:
		c.registers.X = c.registers.A
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case TXA:
		c.registers.A = c.registers.X
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case TAY:
		c.registers.Y = c.registers.A
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case TYA:
		c.registers.A = c.registers.Y
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case TSX:
		c.registers.X = c.registers.S
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case TXS:
		c.registers.UpdateS(c.registers.X)
		return
	case PHA:
		c.stack.Push(c.registers.A)
		return
	case PLA:
		c.registers.A = c.stack.Pop()
		return
	case PHP:
		c.stack.Push(c.registers.P.ToByte())
		return
	case PLP:
		c.registers.P.UpdateAll(c.stack.Pop())
		return
	case NOP:
		return
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}
