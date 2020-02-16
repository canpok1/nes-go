package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/impl/component"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// Instruction ...
type Instruction interface {
	FetchAsOperand() (op []byte, err error)
	Execute(op []byte) (cycle int, err error)
}

// Factory ...
type Factory struct {
	Registers *component.CPURegisters
	Bus       domain.Bus
	Recorder  *domain.Recorder
	Fetch     func() (byte, error)
	PushStack func(byte) error
	PopStack  func() (byte, error)
}

// Make ...
func (f *Factory) Make(opc *domain.OpcodeProp) Instruction {
	return &BaseInstruction{
		registers: f.Registers,
		bus:       f.Bus,
		recorder:  f.Recorder,

		opc: opc,

		fetch:     f.Fetch,
		pushStack: f.PushStack,
		popStack:  f.PopStack,
	}
}

// BaseInstruction ...
type BaseInstruction struct {
	registers *component.CPURegisters
	bus       domain.Bus

	opc      *domain.OpcodeProp
	recorder *domain.Recorder

	fetch     func() (byte, error)
	pushStack func(byte) error
	popStack  func() (byte, error)
}

// fetchAsOperand ...
func (b *BaseInstruction) FetchAsOperand() (op []byte, err error) {
	switch b.opc.AddressingMode {
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
		fallthrough
	case domain.IndexedIndirect:
		fallthrough
	case domain.IndirectIndexed:
		var d byte
		d, err = b.fetch()
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		op = []byte{d}
		return
	case domain.Absolute:
		fallthrough
	case domain.IndexedAbsoluteX:
		fallthrough
	case domain.IndexedAbsoluteY:
		fallthrough
	case domain.AbsoluteIndirect:
		var l byte
		l, err = b.fetch()
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var h byte
		h, err = b.fetch()
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		op = []byte{l, h}
		return
	default:
		err = xerrors.Errorf("failed to fetch operands, AddressingMode is unsupported; mode: %#v", b.opc.AddressingMode)
		return
	}
}

func (b *BaseInstruction) makeAddress(op []byte) (addr domain.Address, pageCrossed bool, err error) {
	log.Trace("BaseInstruction.makeAddress[%#v][%#v] ...", b.opc.AddressingMode, op)
	defer func() {
		if err != nil {
			log.Warn("BaseInstruction.makeAddress[%#v][%#v] => %#v", b.opc.AddressingMode, op, err)
		} else {
			log.Trace("BaseInstruction.makeAddress[%#v][%#v] => %#v", b.opc.AddressingMode, op, addr)
		}
	}()

	switch b.opc.AddressingMode {
	case domain.Absolute:
		l := op[0]
		h := op[1]
		addr = domain.Address((uint16(h) << 8) | uint16(l))
		b.recorder.AddAddress(addr)
		return
	case domain.ZeroPage:
		l := op[0]
		addr = domain.Address(l)
		b.recorder.AddAddress(addr)
		return
	case domain.IndexedZeroPageX:
		l := op[0]
		addr = domain.Address(uint8(l) + uint8(b.registers.X))
		b.recorder.AddAddress(addr)
		return
	case domain.IndexedZeroPageY:
		l := op[0]
		addr = domain.Address(uint8(l) + uint8(b.registers.Y))
		b.recorder.AddAddress(addr)
		return
	case domain.IndexedAbsoluteX:
		l := op[0]
		h := op[1]

		addr = domain.Address(((uint16(h) << 8) | uint16(l)))
		b.recorder.AddAddress(addr)
		addr = domain.Address(uint16(addr) + uint16(b.registers.X))
		b.recorder.AddAddress(addr)

		pageCrossed = (uint16(addr) & 0xFF00) != (uint16(h) << 8)

		return
	case domain.IndexedAbsoluteY:
		l := op[0]
		h := op[1]

		addr = domain.Address(((uint16(h) << 8) | uint16(l)))
		b.recorder.AddAddress(addr)
		addr = domain.Address(uint16(addr) + uint16(b.registers.Y))
		b.recorder.AddAddress(addr)

		pageCrossed = (uint16(addr) & 0xFF00) != (uint16(h) << 8)

		return
	case domain.Relative:
		d := op[0]
		addr = domain.Address(b.registers.PC + uint16(int8(d)))
		b.recorder.AddAddress(addr)

		pageCrossed = (uint16(addr) & 0xFF00) != (b.registers.PC & 0xFF00)

		return
	case domain.IndexedIndirect:
		d := op[0]
		destL := domain.Address(d + b.registers.X)
		destH := domain.Address(d + b.registers.X + 1)
		b.recorder.AddAddress(destL)

		var l byte
		l, err = b.bus.ReadByCPU(destL)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var h byte
		h, err = b.bus.ReadByCPU(destH)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		addr = domain.Address((uint16(h) << 8) | uint16(l))
		b.recorder.AddAddress(addr)

		pageCrossed = (uint16(addr) & 0xFF00) != (uint16(h) << 8)

		return
	case domain.IndirectIndexed:
		d := op[0]
		destL := domain.Address(d)
		destH := domain.Address(d + 1)

		var h byte
		h, err = b.bus.ReadByCPU(destH)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var l byte
		l, err = b.bus.ReadByCPU(destL)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		addr = domain.Address((uint16(h) << 8) + uint16(l))
		b.recorder.AddAddress(addr)

		addr = domain.Address(uint16(addr) + uint16(b.registers.Y))
		b.recorder.AddAddress(addr)

		pageCrossed = (uint16(addr) & 0xFF00) != (uint16(h) << 8)

		return
	case domain.AbsoluteIndirect:
		f1 := op[0]
		f2 := op[1]

		destL := domain.Address((uint16(f2) << 8) | uint16(f1))

		// 6502のバグ：上位バイト取得先アドレスの上位8ビットはインクリメントの影響を受けない
		destH := domain.Address((uint16(f2) << 8) | uint16(f1+1))

		b.recorder.AddAddress(destL)

		var addrL byte
		addrL, err = b.bus.ReadByCPU(destL)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var addrH byte
		addrH, err = b.bus.ReadByCPU(destH)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		addr = domain.Address((uint16(addrH) << 8) | uint16(addrL))
		b.recorder.AddAddress(addr)
		return
	default:
		err = xerrors.Errorf("failed to make address, AddressingMode is not supported; mode: %#v", b.opc.AddressingMode)
		return
	}
}

// Execute ...
func (c *BaseInstruction) Execute(op []byte) (cycle int, err error) {
	mne := c.opc.Mnemonic
	mode := c.opc.AddressingMode
	cycle = c.opc.Cycle

	log.Trace("BaseInstruction.Execute[%#x][%v][%v][%#v] ...", c.registers.PC, mne, mode, op)

	c.recorder.Mnemonic = mne
	c.recorder.Documented = c.opc.Documented
	c.recorder.AddressingMode = mode

	defer func() {
		if err != nil {
			log.Warn("BaseInstruction.Execute[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("BaseInstruction.Execute[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

	switch mne {
	case domain.ADC:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		ans := uint16(c.registers.A) + uint16(b)
		if c.registers.P.Carry {
			ans = ans + 1
		}

		beforeA := c.registers.A
		c.registers.A = byte(ans & 0xFF)
		c.registers.P.UpdateN(c.registers.A)
		if (b & 0x80) == 0x00 {
			c.registers.P.UpdateV(beforeA, c.registers.A)
		} else {
			c.registers.P.ClearV()
		}
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.UpdateC(ans)
		return
	case domain.SBC:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		ans := uint16(c.registers.A) - uint16(b)
		if !c.registers.P.Carry {
			ans = ans - 1
		}

		beforeA := c.registers.A
		c.registers.A = byte(ans & 0x00FF)
		c.registers.P.UpdateN(c.registers.A)
		if (beforeA & 0x80) == (b & 0x80) {
			c.registers.P.ClearV()
		} else {
			c.registers.P.UpdateV(beforeA, c.registers.A)
		}
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (ans & 0xFF00) == 0x0000
		return
	case domain.AND:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		c.registers.A = c.registers.A & b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.ORA:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		c.registers.A = c.registers.A | b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.EOR:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		c.registers.A = c.registers.A ^ b
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.ASL:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b << 1
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b << 1
			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case domain.LSR:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b >> 1
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b >> 1
			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case domain.ROL:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b << 1
			if c.registers.P.Carry {
				ans = ans + 1
			}
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b << 1
			if c.registers.P.Carry {
				ans = ans + 1
			}

			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case domain.ROR:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b >> 1
			if c.registers.P.Carry {
				ans = ans | 0x80
			}

			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b >> 1
			if c.registers.P.Carry {
				ans = ans | 0x80
			}

			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case domain.BCC:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if !c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BCS:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if c.registers.P.Carry {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BEQ:
		var addr domain.Address
		var pageCrossed bool
		if addr, pageCrossed, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
			cycle++
			if pageCrossed {
				cycle++
			}
		}

		return
	case domain.BNE:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if !c.registers.P.Zero {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BVC:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if !c.registers.P.Overflow {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BVS:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if c.registers.P.Overflow {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BPL:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if !c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BMI:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		if c.registers.P.Negative {
			c.registers.UpdatePC(uint16(addr))
			cycle++
		}
		return
	case domain.BIT:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		c.registers.P.Zero = (c.registers.A & b) == 0
		c.registers.P.Negative = (b & 0x80) == 0x80
		c.registers.P.Overflow = (b & 0x40) == 0x40
		return
	case domain.JMP:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		c.registers.UpdatePC(uint16(addr))
		return
	case domain.JSR:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		// 6502のバグ:1つ前のアドレスを格納
		pc := c.registers.PC - 1

		if err = c.pushStack(byte((pc & 0xFF00) >> 8)); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		if err = c.pushStack(byte(pc & 0x00FF)); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.PC = uint16(addr)
		return
	case domain.RTS:
		var l, h byte
		if l, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		if h, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.PC = (uint16(h) << 8) | uint16(l)

		// 6502のバグ:インクリメントしたもの復帰アドレスとする
		c.registers.PC = c.registers.PC + 1

		return
	case domain.BRK:
		c.registers.P.BreakMode = true
		c.registers.IncrementPC()
		return
	case domain.RTI:
		var b byte
		if b, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.P.UpdateAll(b)

		var l, h byte
		if l, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		if h, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.PC = (uint16(h) << 8) | uint16(l)
		return
	case domain.CMP:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		ans := c.registers.A - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = c.registers.A >= b
		return
	case domain.CPX:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		ans := c.registers.X - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = c.registers.X >= b
		return
	case domain.CPY:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.recorder.Data = &b
		ans := c.registers.Y - b
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = c.registers.Y >= b
		return
	case domain.INC:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		ans := b + 1
		err = c.bus.WriteByCPU(addr, ans)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		return
	case domain.DEC:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		ans := b - 1
		err = c.bus.WriteByCPU(addr, ans)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
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
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var (
				addr        domain.Address
				pageCrossed bool
			)
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			switch c.opc.AddressingMode {
			case domain.IndexedAbsoluteX:
				fallthrough
			case domain.IndexedAbsoluteY:
				fallthrough
			case domain.IndirectIndexed:
				if pageCrossed {
					cycle++
				}
			}
		}
		c.recorder.Data = &b

		c.registers.UpdateA(b)
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		return
	case domain.LDX:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			var pageCrossed bool
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			if c.opc.AddressingMode == domain.IndexedAbsoluteY && pageCrossed {
				cycle++
			}
		}
		c.recorder.Data = &b

		c.registers.UpdateX(b)
		c.registers.P.UpdateN(c.registers.X)
		c.registers.P.UpdateZ(c.registers.X)
		return
	case domain.LDY:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			var pageCrossed bool
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			if c.opc.AddressingMode == domain.IndexedAbsoluteX && pageCrossed {
				cycle++
			}
		}
		c.recorder.Data = &b

		c.registers.UpdateY(b)
		c.registers.P.UpdateN(c.registers.Y)
		c.registers.P.UpdateZ(c.registers.Y)
		return
	case domain.STA:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		err = c.bus.WriteByCPU(addr, c.registers.A)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		return
	case domain.STX:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		err = c.bus.WriteByCPU(addr, c.registers.X)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		return
	case domain.STY:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf("failed to exec, address is nil; mnemonic: %#v, op: %#v", mne, op)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		err = c.bus.WriteByCPU(addr, c.registers.Y)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
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
		if err = c.pushStack(c.registers.A); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		return
	case domain.PLA:
		if c.registers.A, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.UpdateN(c.registers.A)
		return
	case domain.PHP:
		p := c.registers.P.ToByte()

		// 6502のバグ：スタックに格納するフラグはBフラグがセットされた状態になる
		p = (p & 0xEF) | 0x10

		if err = c.pushStack(p); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		return
	case domain.PLP:
		var b byte
		if b, err = c.popStack(); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.registers.P.UpdateAll(b)
		return
	case domain.LAX:
		var b byte
		if mode == domain.Immediate {
			if len(op) < 1 {
				err = xerrors.Errorf("failed to exec, data is nil; mnemonic: %#v, op: %#v", mne, op)
				return
			}
			b = op[0]
		} else {
			var addr domain.Address
			var pageCrossed bool
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			if c.opc.AddressingMode == domain.IndirectIndexed && pageCrossed {
				cycle++
			}
		}
		c.recorder.Data = &b

		c.registers.UpdateA(b)
		c.registers.UpdateX(b)

		c.registers.P.UpdateN(b)
		c.registers.P.UpdateZ(b)
		return
	case domain.SAX:
		switch mode {
		case domain.Absolute:
			fallthrough
		case domain.ZeroPage:
			fallthrough
		case domain.IndexedZeroPageX:
			fallthrough
		case domain.IndexedZeroPageY:
			fallthrough
		case domain.IndexedAbsoluteX:
			fallthrough
		case domain.IndexedAbsoluteY:
			fallthrough
		case domain.Relative:
			fallthrough
		case domain.IndexedIndirect:
			fallthrough
		case domain.IndirectIndexed:
			fallthrough
		case domain.AbsoluteIndirect:
			var addr domain.Address
			var pageCrossed bool
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			var b byte
			b, err = c.bus.ReadByRecorder(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			switch c.opc.AddressingMode {
			case domain.IndexedAbsoluteX:
				if pageCrossed {
					cycle++
				}
			}

			ans := c.registers.A & c.registers.X
			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		return
	case domain.DCP:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		var b byte
		b, err = c.bus.ReadByRecorder(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		ans := b - 1
		err = c.bus.WriteByCPU(addr, ans)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		ans = c.registers.A - ans
		c.registers.P.UpdateN(ans)
		c.registers.P.UpdateZ(ans)
		c.registers.P.Carry = c.registers.A >= b
		return
	case domain.ISB:
		fallthrough
	case domain.ISC:
		var addr domain.Address
		if addr, _, err = c.makeAddress(op); err != nil {
			return
		}

		var b byte
		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
		c.recorder.Data = &b

		ans1 := b + 1
		err = c.bus.WriteByCPU(addr, ans1)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		ans2 := uint16(c.registers.A) - uint16(ans1)
		if !c.registers.P.Carry {
			ans2 = ans2 - 1
		}

		beforeA := c.registers.A
		c.registers.A = byte(ans2 & 0x00FF)
		c.registers.P.UpdateN(c.registers.A)
		if (beforeA & 0x80) == (b & 0x80) {
			c.registers.P.ClearV()
		} else {
			c.registers.P.UpdateV(beforeA, c.registers.A)
		}
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (ans2 & 0xFF00) == 0x0000
		return
	case domain.SLO:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b << 1
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b << 1
			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.registers.A = c.registers.A | ans
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case domain.RLA:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b << 1
			if c.registers.P.Carry {
				ans = ans + 1
			}
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b << 1
			if c.registers.P.Carry {
				ans = ans + 1
			}

			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		c.registers.A = c.registers.A & ans
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x80) == 0x80
		return
	case domain.SRE:
		var b, ans byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans = b >> 1
			c.registers.A = ans
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans = b >> 1
			err = c.bus.WriteByCPU(addr, ans)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}

		c.registers.A = c.registers.A ^ ans
		c.registers.P.UpdateN(c.registers.A)
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.Carry = (b & 0x01) == 0x01
		return
	case domain.RRA:
		var b, ans1 byte
		if mode == domain.Accumulator {
			b = c.registers.A
			c.recorder.Data = &b

			ans1 = b >> 1
			if c.registers.P.Carry {
				ans1 = ans1 | 0x80
			}

			c.registers.A = ans1
		} else {
			var addr domain.Address
			if addr, _, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			b, err = c.bus.ReadByCPU(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			ans1 = b >> 1
			if c.registers.P.Carry {
				ans1 = ans1 | 0x80
			}

			err = c.bus.WriteByCPU(addr, ans1)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
		}
		c.registers.P.Carry = (b & 0x01) == 0x01

		ans2 := uint16(c.registers.A) + uint16(ans1)
		if c.registers.P.Carry {
			ans2 = ans2 + 1
		}

		beforeA := c.registers.A
		c.registers.A = byte(ans2 & 0xFF)
		c.registers.P.UpdateN(c.registers.A)
		if (b & 0x80) == 0x00 {
			c.registers.P.UpdateV(beforeA, c.registers.A)
		} else {
			c.registers.P.ClearV()
		}
		c.registers.P.UpdateZ(c.registers.A)
		c.registers.P.UpdateC(ans2)
		return
	case domain.STP:
		fallthrough
	case domain.NOP:
		switch mode {
		case domain.Absolute:
			fallthrough
		case domain.ZeroPage:
			fallthrough
		case domain.IndexedZeroPageX:
			fallthrough
		case domain.IndexedZeroPageY:
			fallthrough
		case domain.IndexedAbsoluteX:
			fallthrough
		case domain.IndexedAbsoluteY:
			fallthrough
		case domain.Relative:
			fallthrough
		case domain.IndexedIndirect:
			fallthrough
		case domain.IndirectIndexed:
			fallthrough
		case domain.AbsoluteIndirect:
			var addr domain.Address
			var pageCrossed bool
			if addr, pageCrossed, err = c.makeAddress(op); err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}

			var b byte
			b, err = c.bus.ReadByRecorder(addr)
			if err != nil {
				err = xerrors.Errorf(": %w", err)
				return
			}
			c.recorder.Data = &b

			switch c.opc.AddressingMode {
			case domain.IndexedAbsoluteX:
				if pageCrossed {
					cycle++
				}
			}

		case domain.Immediate:
			b := op[0]
			c.recorder.Data = &b
		}

		return
	default:
		err = xerrors.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", mne)
		return
	}
}
