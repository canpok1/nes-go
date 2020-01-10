package impl

import (
	"fmt"

	"nes-go/pkg/domain"
	"nes-go/pkg/log"
	"nes-go/pkg/impl/cpu"
)

// Operand ...
type Operand struct {
	Data    *byte
	Address *domain.Address
}

// String ...
func (o *Operand) String() string {
	d := ""
	if o.Data != nil {
		d = fmt.Sprintf("%#v", *o.Data)
	}

	a := ""
	if o.Address != nil {
		a = fmt.Sprintf("%#v", *o.Address)
	}

	return fmt.Sprintf("{Data:%#v, Address:%#v}", d, a)
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
	registers   *cpu.Registers
	bus         *Bus
	shouldReset bool
	stack       *CPUStack
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{
		registers:   cpu.NewRegisters(),
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

	if !c.registers.P.InterruptDisable && c.registers.P.BreakMode {
		if err := c.InterruptBRK(); err != nil {
			return 0, fmt.Errorf("%w", err)
		}
		return 0, nil
	}

	if !c.registers.P.InterruptDisable && !c.registers.P.BreakMode {
		if err := c.InterruptIRQ(); err != nil {
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
func decodeOpcode(o domain.Opcode) (*domain.OpcodeProp, error) {
	if p, ok := domain.OpcodeProps[o]; ok {
		log.Trace("CPU.decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Trace("CPU.decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr domain.Address
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

	addr = domain.Address(c.registers.PC)
	data, err = c.bus.ReadByCPU(addr)
	if err != nil {
		return data, fmt.Errorf("failed to fetch; %w", err)
	}

	c.registers.IncrementPC()

	return data, nil
}

// fetchOpcode ...
func (c *CPU) fetchOpcode() (domain.Opcode, error) {
	v, err := c.fetch()
	if err != nil {
		return domain.ErrorOpcode, fmt.Errorf("%w", err)
	}
	return domain.Opcode(v), nil
}

// fetchAsOperand ...
func (c *CPU) fetchAsOperand(mode domain.AddressingMode) (op []byte, err error) {
	log.Trace("CPU.fetchAsOperand[%#v] ...", mode)
	defer func() {
		if err != nil {
			log.Warn("CPU.fetchAsOperand[%#v] => %#v", mode, err)
		} else {
			log.Trace("CPU.fetchAsOperand[%#v] => %#v", mode, op)
		}
	}()

	switch mode {
	case domain.Accumulator:
		fallthrough
	case domain.Implied:
		op = []byte{}
		return
	case domain.Immediate:
		fallthrough
	case domain.ZeroPage:
		fallthrough
	case domain.IndexedZeroPageX:
		fallthrough
	case domain.IndexedZeroPageY:
		fallthrough
	case domain.Relative:
		var b byte
		b, err = c.fetch()
		if err != nil {
			return
		}
		op = []byte{b}
		return
	case domain.Absolute:
		fallthrough
	case domain.IndexedAbsoluteX:
		fallthrough
	case domain.IndexedAbsoluteY:
		fallthrough
	case domain.IndexedIndirect:
		fallthrough
	case domain.IndirectIndexed:
		fallthrough
	case domain.AbsoluteIndirect:
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
func (c *CPU) makeAddress(mode domain.AddressingMode, op []byte) (addr domain.Address, err error) {
	log.Trace("CPU.makeAddress[%#v][%#v] ...", mode, op)
	defer func() {
		if err != nil {
			log.Warn("CPU.makeAddress[%#v][%#v] => %#v", mode, op, err)
		} else {
			log.Trace("CPU.makeAddress[%#v][%#v] => %#v", mode, op, addr)
		}
	}()

	switch mode {
	case domain.Absolute:
		l := op[0]
		h := op[1]
		addr = domain.Address((uint16(h) << 8) | uint16(l))
		return
	case domain.ZeroPage:
		l := op[0]
		addr = domain.Address(l)
		return
	case domain.IndexedZeroPageX:
		l := op[0]
		addr = domain.Address(uint8(l) + uint8(c.registers.X))
		return
	case domain.IndexedZeroPageY:
		l := op[0]
		addr = domain.Address(uint8(l) + uint8(c.registers.Y))
		return
	case domain.IndexedAbsoluteX:
		l := op[0]
		h := op[1]
		addr = domain.Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.X))
		return
	case domain.IndexedAbsoluteY:
		l := op[0]
		h := op[1]
		addr = domain.Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.Y))
		return
	case domain.Relative:
		b := op[0]
		addr = domain.Address(c.registers.PC + uint16(int8(b)))
		return
	case domain.IndexedIndirect:
		b := op[0]
		dest := domain.Address(uint8(b) + c.registers.X)

		var l byte
		l, err = c.bus.ReadByCPU(dest)
		if err != nil {
			return
		}

		h := op[1]

		addr = domain.Address((uint16(h) << 8) | uint16(l))
		return
	case domain.IndirectIndexed:
		b := op[0]
		dest := domain.Address(uint8(b) + c.registers.X)

		var h byte
		h, err = c.bus.ReadByCPU(dest)
		if err != nil {
			return
		}

		l := op[1]

		addr = domain.Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.Y))
		return
	case domain.AbsoluteIndirect:
		f1 := op[0]
		f2 := op[1]

		dest := domain.Address((uint16(f2) << 8) + uint16(f1))
		nextDest := dest + 1

		var addrL byte
		addrL, err = c.bus.ReadByCPU(dest)
		if err != nil {
			return
		}

		var addrH byte
		addrH, err = c.bus.ReadByCPU(nextDest)
		if err != nil {
			return
		}

		addr = domain.Address((uint16(addrH) << 8) + uint16(addrL))
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

	l, err := c.bus.ReadByCPU(0xFFFC)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	h, err := c.bus.ReadByCPU(0xFFFD)
	if err != nil {
		return fmt.Errorf("failed to reset; %w", err)
	}

	c.registers.UpdatePC((uint16(h) << 8) | uint16(l))

	c.shouldReset = false
	return nil
}

// InterruptBRK ...
func (c *CPU) InterruptBRK() error {
	log.Info("CPU.Interrupt[BRK] ...")

	c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
	c.stack.Push(byte(c.registers.PC & 0x00FF))
	c.stack.Push(c.registers.P.ToByte())
	c.registers.P.InterruptDisable = true

	l, err := c.bus.ReadByCPU(0xFFFE)
	if err != nil {
		return fmt.Errorf("failed to interrupt[BRK]; %w", err)
	}
	h, err := c.bus.ReadByCPU(0xFFFF)
	if err != nil {
		return fmt.Errorf("failed to interrupt[BRK]; %w", err)
	}
	c.registers.PC = (uint16(h) << 8) + uint16(l)
	return nil
}

// InterruptIRQ ...
func (c *CPU) InterruptIRQ() error {
	log.Info("CPU.Interrupt[IRQ] ...")

	c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
	c.stack.Push(byte(c.registers.PC & 0x00FF))
	c.stack.Push(c.registers.P.ToByte())
	c.registers.P.InterruptDisable = true

	l, err := c.bus.ReadByCPU(0xFFFE)
	if err != nil {
		return fmt.Errorf("failed to interrupt[IRQ]; %w", err)
	}
	h, err := c.bus.ReadByCPU(0xFFFF)
	if err != nil {
		return fmt.Errorf("failed to interrupt[IRQ]; %w", err)
	}
	c.registers.PC = (uint16(h) << 8) + uint16(l)
	return nil
}

// exec ...
func (c *CPU) exec(mne domain.Mnemonic, mode domain.AddressingMode, op []byte) (err error) {
	log.Debug("CPU.exec[%v][%v][%#v] ...", mne, mode, op)

	defer func() {
		if err != nil {
			log.Warn("CPU.exec[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("CPU.exec[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

	switch mne {
	case domain.ADC:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
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
	case domain.SBC:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
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
	case domain.AND:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.A = c.registers.A & b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.ORA:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.A = c.registers.A | b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.EOR:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.A = c.registers.A ^ b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.ASL:
		var b byte
		if mode == domain.Accumulator {
			b = c.registers.A
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.A = b << 1
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case domain.LSR:
		var b byte
		if mode == domain.Accumulator {
			b = c.registers.A
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		c.registers.A = b >> 1
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case domain.ROL:
		var b byte
		if mode == domain.Accumulator {
			b = c.registers.A
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
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
	case domain.ROR:
		var b byte
		if mode == domain.Accumulator {
			b = c.registers.A
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
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
	case domain.BCC:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BCS:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BEQ:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BNE:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BVC:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BVS:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BPL:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BMI:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case domain.BIT:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			return
		}

		c.registers.P.Zero = (c.registers.A & b) == 0
		c.registers.P.Negative = (b & 0x80) == 0x80
		c.registers.P.Overflow = (b & 0x40) == 0x40
		return
	case domain.JMP:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		c.registers.UpdatePC(uint16(addr))
		return
	case domain.JSR:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			return
		}

		c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.PC & 0x00FF))
		c.registers.PC = uint16(b)
		return
	case domain.RTS:
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l) + 1
		return
	case domain.BRK:
		c.registers.P.BreakMode = true
		c.registers.IncrementPC()
		return
	case domain.RTI:
		c.registers.P.UpdateAll(c.stack.Pop())
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l)
		return
	case domain.CMP:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.A - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
		return
	case domain.CPX:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.X - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
		return
	case domain.CPY:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}
		ans := c.registers.Y - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = ans >= 0
		return
	case domain.INC:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			return
		}

		ans := b + 1
		err = c.bus.WriteByCPU(addr, ans)
		if err != nil {
			return
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		return
	case domain.DEC:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			return
		}

		ans := b - 1
		err = c.bus.WriteByCPU(addr, ans)
		if err != nil {
			return
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		return
	case domain.INX:
		c.registers.UpdateX(c.registers.X + 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.DEX:
		c.registers.UpdateX(c.registers.X - 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.INY:
		c.registers.UpdateY(c.registers.Y + 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case domain.DEY:
		c.registers.UpdateY(c.registers.Y - 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case domain.CLC:
		c.registers.P.Carry = false
		return
	case domain.SEC:
		c.registers.P.Carry = true
		return
	case domain.CLI:
		c.registers.P.UpdateI(false)
		return
	case domain.SEI:
		c.registers.P.UpdateI(true)
		return
	case domain.CLD:
		c.registers.P.DecimalMode = false
		return
	case domain.SED:
		c.registers.P.DecimalMode = true
		return
	case domain.CLV:
		c.registers.P.Overflow = false
		return
	case domain.LDA:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.UpdateA(b)
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.LDX:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.UpdateX(b)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.LDY:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, err = c.makeAddress(mode, op); err != nil {
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				return
			}
		}

		c.registers.UpdateY(b)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case domain.STA:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.A)
		if err != nil {
			return
		}
		return
	case domain.STX:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.X)
		if err != nil {
			return
		}
		return
	case domain.STY:
		var addr domain.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.Y)
		if err != nil {
			return
		}
		return
	case domain.TAX:
		c.registers.X = c.registers.A
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.TXA:
		c.registers.A = c.registers.X
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.TAY:
		c.registers.Y = c.registers.A
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case domain.TYA:
		c.registers.A = c.registers.Y
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.TSX:
		c.registers.X = c.registers.S
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.TXS:
		c.registers.UpdateS(c.registers.X)
		return
	case domain.PHA:
		c.stack.Push(c.registers.A)
		return
	case domain.PLA:
		c.registers.A = c.stack.Pop()
		return
	case domain.PHP:
		c.stack.Push(c.registers.P.ToByte())
		return
	case domain.PLP:
		c.registers.P.UpdateAll(c.stack.Pop())
		return
	case domain.NOP:
		return
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}