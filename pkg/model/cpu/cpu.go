package cpu

import (
	"fmt"

	"nes-go/pkg/log"
	"nes-go/pkg/model"
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
	registers   *Registers
	bus         *model.Bus
	shouldReset bool
	stack       *CPUStack
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{
		registers:   NewRegisters(),
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
func (c *CPU) SetBus(b *model.Bus) {
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
func decodeOpcode(o model.Opcode) (*model.OpcodeProp, error) {
	if p, ok := model.OpcodeProps[o]; ok {
		log.Trace("CPU.decode[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Trace("CPU.decode[%#v] => not found", o)
	return nil, fmt.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr model.Address
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

	addr = model.Address(c.registers.PC)
	data, err = c.bus.ReadByCPU(addr)
	if err != nil {
		return data, fmt.Errorf("failed to fetch; %w", err)
	}

	c.registers.IncrementPC()

	return data, nil
}

// fetchOpcode ...
func (c *CPU) fetchOpcode() (model.Opcode, error) {
	v, err := c.fetch()
	if err != nil {
		return model.ErrorOpcode, fmt.Errorf("%w", err)
	}
	return model.Opcode(v), nil
}

// fetchAsOperand ...
func (c *CPU) fetchAsOperand(mode model.AddressingMode) (op []byte, err error) {
	log.Trace("CPU.fetchAsOperand[%#v] ...", mode)
	defer func() {
		if err != nil {
			log.Warn("CPU.fetchAsOperand[%#v] => %#v", mode, err)
		} else {
			log.Trace("CPU.fetchAsOperand[%#v] => %#v", mode, op)
		}
	}()

	switch mode {
	case model.Accumulator:
		fallthrough
	case model.Implied:
		op = []byte{}
		return
	case model.Immediate:
		fallthrough
	case model.ZeroPage:
		fallthrough
	case model.IndexedZeroPageX:
		fallthrough
	case model.IndexedZeroPageY:
		fallthrough
	case model.Relative:
		var b byte
		b, err = c.fetch()
		if err != nil {
			return
		}
		op = []byte{b}
		return
	case model.Absolute:
		fallthrough
	case model.IndexedAbsoluteX:
		fallthrough
	case model.IndexedAbsoluteY:
		fallthrough
	case model.IndexedIndirect:
		fallthrough
	case model.IndirectIndexed:
		fallthrough
	case model.AbsoluteIndirect:
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
func (c *CPU) makeAddress(mode model.AddressingMode, op []byte) (addr model.Address, err error) {
	log.Trace("CPU.makeAddress[%#v][%#v] ...", mode, op)
	defer func() {
		if err != nil {
			log.Warn("CPU.makeAddress[%#v][%#v] => %#v", mode, op, err)
		} else {
			log.Trace("CPU.makeAddress[%#v][%#v] => %#v", mode, op, addr)
		}
	}()

	switch mode {
	case model.Absolute:
		l := op[0]
		h := op[1]
		addr = model.Address((uint16(h) << 8) | uint16(l))
		return
	case model.ZeroPage:
		l := op[0]
		addr = model.Address(l)
		return
	case model.IndexedZeroPageX:
		l := op[0]
		addr = model.Address(uint8(l) + uint8(c.registers.X))
		return
	case model.IndexedZeroPageY:
		l := op[0]
		addr = model.Address(uint8(l) + uint8(c.registers.Y))
		return
	case model.IndexedAbsoluteX:
		l := op[0]
		h := op[1]
		addr = model.Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.X))
		return
	case model.IndexedAbsoluteY:
		l := op[0]
		h := op[1]
		addr = model.Address(((uint16(h) << 8) | uint16(l)) + uint16(c.registers.Y))
		return
	case model.Relative:
		b := op[0]
		addr = model.Address(c.registers.PC + uint16(int8(b)))
		return
	case model.IndexedIndirect:
		b := op[0]
		dest := model.Address(uint8(b) + c.registers.X)

		var l byte
		l, err = c.bus.ReadByCPU(dest)
		if err != nil {
			return
		}

		h := op[1]

		addr = model.Address((uint16(h) << 8) | uint16(l))
		return
	case model.IndirectIndexed:
		b := op[0]
		dest := model.Address(uint8(b) + c.registers.X)

		var h byte
		h, err = c.bus.ReadByCPU(dest)
		if err != nil {
			return
		}

		l := op[1]

		addr = model.Address((uint16(h) << 8) + uint16(l) + uint16(c.registers.Y))
		return
	case model.AbsoluteIndirect:
		f1 := op[0]
		f2 := op[1]

		dest := model.Address((uint16(f2) << 8) + uint16(f1))
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

		addr = model.Address((uint16(addrH) << 8) + uint16(addrL))
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
func (c *CPU) exec(mne model.Mnemonic, mode model.AddressingMode, op []byte) (err error) {
	log.Debug("CPU.exec[%v][%v][%#v] ...", mne, mode, op)

	defer func() {
		if err != nil {
			log.Warn("CPU.exec[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("CPU.exec[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

	switch mne {
	case model.ADC:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.SBC:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.AND:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.ORA:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.EOR:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.ASL:
		var b byte
		if mode == model.Accumulator {
			b = c.registers.A
		} else {
			var addr model.Address
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
	case model.LSR:
		var b byte
		if mode == model.Accumulator {
			b = c.registers.A
		} else {
			var addr model.Address
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
	case model.ROL:
		var b byte
		if mode == model.Accumulator {
			b = c.registers.A
		} else {
			var addr model.Address
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
	case model.ROR:
		var b byte
		if mode == model.Accumulator {
			b = c.registers.A
		} else {
			var addr model.Address
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
	case model.BCC:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BCS:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BEQ:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BNE:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BVC:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BVS:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BPL:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if !c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BMI:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		if c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
		}
		return
	case model.BIT:
		var addr model.Address
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
	case model.JMP:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			return
		}
		c.registers.UpdatePC(uint16(addr))
		return
	case model.JSR:
		var addr model.Address
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
	case model.RTS:
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l) + 1
		return
	case model.BRK:
		c.registers.P.BreakMode = true
		c.registers.IncrementPC()
		c.stack.Push(byte((c.registers.PC & 0xFF00) >> 8))
		c.stack.Push(byte(c.registers.PC & 0x00FF))
		c.stack.Push(c.registers.P.ToByte())
		c.registers.P.Interrupt = true

		var l, h byte
		l, err = c.bus.ReadByCPU(0xFFFE)
		if err != nil {
			return
		}
		h, err = c.bus.ReadByCPU(0xFFFF)
		if err != nil {
			return
		}
		c.registers.PC = (uint16(h) << 8) + uint16(l)
		return
	case model.RTI:
		c.registers.P.UpdateAll(c.stack.Pop())
		l := c.stack.Pop()
		h := c.stack.Pop()
		c.registers.PC = (uint16(h) << 8) + uint16(l)
		return
	case model.CMP:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.CPX:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.CPY:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.INC:
		var addr model.Address
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
	case model.DEC:
		var addr model.Address
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
	case model.INX:
		c.registers.UpdateX(c.registers.X + 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case model.DEX:
		c.registers.UpdateX(c.registers.X - 1)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case model.INY:
		c.registers.UpdateY(c.registers.Y + 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case model.DEY:
		c.registers.UpdateY(c.registers.Y - 1)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case model.CLC:
		c.registers.P.Carry = false
		return
	case model.SEC:
		c.registers.P.Carry = true
		return
	case model.CLI:
		c.registers.P.UpdateI(false)
		return
	case model.SEI:
		c.registers.P.UpdateI(true)
		return
	case model.CLD:
		c.registers.P.DecimalMode = false
		return
	case model.SED:
		c.registers.P.DecimalMode = true
		return
	case model.CLV:
		c.registers.P.Overflow = false
		return
	case model.LDA:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.LDX:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.LDY:
		var b byte
		if mode == model.Immediate {
			if len(op) < 1 {
				err = fmt.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr model.Address
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
	case model.STA:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.A)
		if err != nil {
			return
		}
		return
	case model.STX:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.X)
		if err != nil {
			return
		}
		return
	case model.STY:
		var addr model.Address
		if addr, err = c.makeAddress(mode, op); err != nil {
			err = fmt.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		err = c.bus.WriteByCPU(addr, c.registers.Y)
		if err != nil {
			return
		}
		return
	case model.TAX:
		c.registers.X = c.registers.A
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case model.TXA:
		c.registers.A = c.registers.X
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case model.TAY:
		c.registers.Y = c.registers.A
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case model.TYA:
		c.registers.A = c.registers.Y
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case model.TSX:
		c.registers.X = c.registers.S
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case model.TXS:
		c.registers.UpdateS(c.registers.X)
		return
	case model.PHA:
		c.stack.Push(c.registers.A)
		return
	case model.PLA:
		c.registers.A = c.stack.Pop()
		return
	case model.PHP:
		c.stack.Push(c.registers.P.ToByte())
		return
	case model.PLP:
		c.registers.P.UpdateAll(c.stack.Pop())
		return
	case model.NOP:
		return
	default:
		err = fmt.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}
