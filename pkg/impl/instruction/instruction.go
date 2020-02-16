package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/impl/component"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// Instruction ...
type Instruction interface {
	SetAllParams(org *BaseInstruction)
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
func (f *Factory) Make(ocp *domain.OpcodeProp) Instruction {
	params := BaseInstruction{
		registers: f.Registers,
		bus:       f.Bus,
		recorder:  f.Recorder,

		ocp: ocp,

		fetch:     f.Fetch,
		pushStack: f.PushStack,
		popStack:  f.PopStack,
	}

	var ins Instruction
	switch ocp.Mnemonic {
	case domain.ADC:
		ins = &ADC{}
	case domain.SBC:
		ins = &SBC{}
	case domain.AND:
		ins = &AND{}
	case domain.ORA:
		ins = &ORA{}
	case domain.EOR:
		ins = &EOR{}
	case domain.ASL:
		ins = &ASL{}
	case domain.LSR:
		ins = &LSR{}
	case domain.ROL:
		ins = &ROL{}
	case domain.ROR:
		ins = &ROR{}
	case domain.BCC:
		ins = &BCC{}
	case domain.BCS:
		ins = &BCS{}
	case domain.BEQ:
		ins = &BEQ{}
	case domain.BNE:
		ins = &BNE{}
	case domain.BVC:
		ins = &BVC{}
	case domain.BVS:
		ins = &BVS{}
	case domain.BPL:
		ins = &BPL{}
	case domain.BMI:
		ins = &BMI{}
	case domain.BIT:
		ins = &BIT{}
	case domain.JMP:
		ins = &JMP{}
	case domain.JSR:
		ins = &JSR{}
	case domain.RTS:
		ins = &RTS{}
	case domain.BRK:
		ins = &BRK{}
	case domain.RTI:
		ins = &RTI{}
	case domain.CMP:
		ins = &CMP{}
	case domain.CPX:
		ins = &CPX{}
	case domain.CPY:
		ins = &CPY{}
	case domain.INC:
		ins = &INC{}
	case domain.DEC:
		ins = &DEC{}
	case domain.INX:
		ins = &INX{}
	case domain.DEX:
		ins = &DEX{}
	case domain.INY:
		ins = &INY{}
	case domain.DEY:
		ins = &DEY{}
	case domain.CLC:
		ins = &CLC{}
	case domain.SEC:
		ins = &SEC{}
	case domain.CLI:
		ins = &CLI{}
	case domain.SEI:
		ins = &SEI{}
	case domain.CLD:
		ins = &CLD{}
	case domain.SED:
		ins = &SED{}
	case domain.CLV:
		ins = &CLV{}
	case domain.LDA:
		ins = &LDA{}
	case domain.LDX:
		ins = &LDX{}
	case domain.LDY:
		ins = &LDY{}
	case domain.STA:
		ins = &STA{}
	case domain.STX:
		ins = &STX{}
	case domain.STY:
		ins = &STY{}
	case domain.TAX:
		ins = &TAX{}
	case domain.TXA:
		ins = &TXA{}
	case domain.TAY:
		ins = &TAY{}
	case domain.TYA:
		ins = &TYA{}
	case domain.TSX:
		ins = &TSX{}
	case domain.TXS:
		ins = &TXS{}
	case domain.PHA:
		ins = &PHA{}
	case domain.PLA:
		ins = &PLA{}
	case domain.PHP:
		ins = &PHP{}
	case domain.PLP:
		ins = &PLP{}
	case domain.LAX:
		ins = &LAX{}
	case domain.SAX:
		ins = &SAX{}
	case domain.DCP:
		ins = &DCP{}
	case domain.ISC:
		fallthrough
	case domain.ISB:
		ins = &ISC{}
	case domain.SLO:
		ins = &SLO{}
	case domain.RLA:
		ins = &RLA{}
	case domain.SRE:
		ins = &SRE{}
	case domain.RRA:
		ins = &RRA{}
	case domain.NOP:
		ins = &NOP{}
	default:
		ins = &BaseInstruction{}
	}
	ins.SetAllParams(&params)
	return ins
}

// BaseInstruction ...
type BaseInstruction struct {
	registers *component.CPURegisters
	bus       domain.Bus

	ocp      *domain.OpcodeProp
	recorder *domain.Recorder

	fetch     func() (byte, error)
	pushStack func(byte) error
	popStack  func() (byte, error)
}

// SetAllParams ...
func (b *BaseInstruction) SetAllParams(org *BaseInstruction) {
	b.registers = org.registers
	b.bus = org.bus
	b.ocp = org.ocp
	b.recorder = org.recorder
	b.fetch = org.fetch
	b.pushStack = org.pushStack
	b.popStack = org.popStack
}

// FetchAsOperand ...
func (b *BaseInstruction) FetchAsOperand() (op []byte, err error) {
	switch b.ocp.AddressingMode {
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
		err = xerrors.Errorf("failed to fetch operands, AddressingMode is unsupported; mode: %#v", b.ocp.AddressingMode)
		return
	}
}

func (b *BaseInstruction) makeAddress(op []byte) (addr domain.Address, pageCrossed bool, err error) {
	log.Trace("BaseInstruction.makeAddress[%#v][%#v] ...", b.ocp.AddressingMode, op)
	defer func() {
		if err != nil {
			log.Warn("BaseInstruction.makeAddress[%#v][%#v] => %#v", b.ocp.AddressingMode, op, err)
		} else {
			log.Trace("BaseInstruction.makeAddress[%#v][%#v] => %#v", b.ocp.AddressingMode, op, addr)
		}
	}()

	switch b.ocp.AddressingMode {
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
		err = xerrors.Errorf("failed to make address, AddressingMode is not supported; mode: %#v", b.ocp.AddressingMode)
		return
	}
}

// Execute ...
func (c *BaseInstruction) Execute(op []byte) (cycle int, err error) {
	return 0, xerrors.Errorf("failed to exec, mnemonic is not supported; mnemonic: %#v", c.ocp.AddressingMode)
}
