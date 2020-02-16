package impl

import (
	"fmt"

	"nes-go/pkg/domain"
	"nes-go/pkg/impl/component"
	"nes-go/pkg/impl/instruction"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
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

// CPU ...
type CPU struct {
	registers   *component.CPURegisters
	bus         domain.Bus
	shouldReset bool
	shouldNMI   bool

	beforeNMIActive bool

	firstPC    *uint16
	executeLog *domain.Recorder

	iFactory *instruction.Factory
}

// NewCPU ...
func NewCPU(pc *uint16) domain.CPU {
	c := CPU{
		registers:       component.NewCPURegisters(),
		shouldReset:     true,
		shouldNMI:       false,
		beforeNMIActive: false,
		firstPC:         pc,
		executeLog:      &domain.Recorder{},
	}

	f := instruction.Factory{
		Registers: c.registers,
		Bus:       c.bus,
		Fetch:     c.fetch,
		PushStack: c.pushStack,
		PopStack:  c.popStack,
	}
	c.iFactory = &f

	return &c
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
func (c *CPU) SetBus(b domain.Bus) {
	c.bus = b
	c.iFactory.Bus = b
}

// SetRecorder ...
func (c *CPU) SetRecorder(r *domain.Recorder) {
	c.executeLog = r
	c.iFactory.Recorder = r
}

// Run ...
func (c *CPU) Run() (int, error) {
	log.Trace("===== CPU RUN =====")
	log.Trace(c.String())

	c.executeLog.PC = c.registers.PC
	c.executeLog.FetchedValue = nil
	c.executeLog.Mnemonic = domain.NOP
	c.executeLog.Documented = false
	c.executeLog.AddressingMode = domain.Implied
	c.executeLog.Address = nil
	c.executeLog.Data = nil
	c.executeLog.A = c.registers.A
	c.executeLog.X = c.registers.X
	c.executeLog.Y = c.registers.Y
	c.executeLog.P = c.registers.P.ToByte()
	c.executeLog.SP = c.registers.S

	if c.shouldReset {
		if err := c.interruptRESET(); err != nil {
			return 0, xerrors.Errorf(": %w", err)
		}
		return 7, nil
	}

	if c.shouldNMI {
		if err := c.InterruptNMI(); err != nil {
			return 0, xerrors.Errorf(": %w", err)
		}
		return 0, nil
	}

	if !c.registers.P.InterruptDisable && c.registers.P.BreakMode {
		if err := c.interruptBRK(); err != nil {
			return 0, xerrors.Errorf(": %w", err)
		}
		return 0, nil
	}

	if !c.registers.P.InterruptDisable && !c.registers.P.BreakMode {
		if err := c.interruptIRQ(); err != nil {
			return 0, xerrors.Errorf(": %w", err)
		}
		return 0, nil
	}

	// PC（プログラムカウンタ）からオペコードをフェッチ（PCをインクリメント）
	oc, err := c.fetchOpcode()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	// 命令とアドレッシング・モードを判別
	ocp, err := decodeOpcode(oc)
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	instruction := c.iFactory.Make(ocp)

	// （必要であれば）オペランドをフェッチ（PCをインクリメント）
	op, err := instruction.FetchAsOperand()
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	// 命令を実行
	if cycle, err := instruction.Execute(op); err != nil {
		return 0, xerrors.Errorf(": %w", err)
	} else {
		return cycle, nil
	}
}

// decodeOpcode ...
func decodeOpcode(o domain.Opcode) (*domain.OpcodeProp, error) {
	if p, ok := domain.OpcodeProps[o]; ok {
		log.Trace("begin[opcode=%#v] => %#v", o, p)
		return &p, nil
	}
	log.Trace("begin[%#v] => not found", o)
	return nil, xerrors.Errorf("opcode is not support; opcode: %#v", o)
}

// fetch ...
func (c *CPU) fetch() (byte, error) {
	var addr domain.Address
	var data byte
	var err error

	log.Trace("begin ...")
	defer func() {
		if err != nil {
			log.Warn("end[addr=%#v] => error %#v", addr, err)
		} else {
			log.Trace("end[addr=%#v] => %#v", addr, data)
		}
	}()

	addr = domain.Address(c.registers.PC)
	data, err = c.bus.ReadByCPU(addr)
	if err != nil {
		return data, xerrors.Errorf("failed to fetch: %w", err)
	}

	c.registers.IncrementPC()

	c.executeLog.FetchedValue = append(c.executeLog.FetchedValue, data)

	return data, nil
}

// fetchOpcode ...
func (c *CPU) fetchOpcode() (domain.Opcode, error) {
	v, err := c.fetch()
	if err != nil {
		return domain.ErrorOpcode, xerrors.Errorf(": %w", err)
	}
	return domain.Opcode(v), nil
}

// InterruptNMI ...
func (c *CPU) InterruptNMI() error {
	log.Trace("begin[Interrupt NMI] ...")
	defer log.Trace("end[Interrupt NMI]")

	c.registers.P.BreakMode = false
	if err := c.pushStack(byte((c.registers.PC & 0xFF00) >> 8)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(byte(c.registers.PC & 0x00FF)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(c.registers.P.ToByte()); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	c.registers.P.InterruptDisable = true

	l, err := c.bus.ReadByCPU(0xFFFA)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	h, err := c.bus.ReadByCPU(0xFFFB)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	c.registers.PC = (uint16(h) << 8) + uint16(l)

	c.shouldNMI = false
	return nil
}

// interruptRESET ...
func (c *CPU) interruptRESET() error {
	log.Trace("begin[Interrupt RESET] ...")
	defer log.Trace("end[Interrupt RESET]")

	c.registers.P.UpdateI(true)

	if c.firstPC != nil {
		c.registers.UpdatePC(*c.firstPC)
	} else {
		l, err := c.bus.ReadByCPU(0xFFFC)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		h, err := c.bus.ReadByCPU(0xFFFD)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		c.registers.UpdatePC((uint16(h) << 8) | uint16(l))
	}

	c.shouldReset = false
	return nil
}

// interruptBRK ...
func (c *CPU) interruptBRK() error {
	log.Trace("begin[Interrupt BRK] ...")
	defer log.Trace("end[Interrupt BRK]")

	if err := c.pushStack(byte((c.registers.PC & 0xFF00) >> 8)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(byte(c.registers.PC & 0x00FF)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(c.registers.P.ToByte()); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	c.registers.P.InterruptDisable = true

	l, err := c.bus.ReadByCPU(0xFFFE)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	h, err := c.bus.ReadByCPU(0xFFFF)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	c.registers.PC = (uint16(h) << 8) + uint16(l)
	return nil
}

// interruptIRQ ...
func (c *CPU) interruptIRQ() error {
	log.Trace("begin[Interrupt IRQ] ...")
	defer log.Trace("end[Interrupt IRQ]")

	if err := c.pushStack(byte((c.registers.PC & 0xFF00) >> 8)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(byte(c.registers.PC & 0x00FF)); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.pushStack(c.registers.P.ToByte()); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	c.registers.P.InterruptDisable = true

	l, err := c.bus.ReadByCPU(0xFFFE)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	h, err := c.bus.ReadByCPU(0xFFFF)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	c.registers.PC = (uint16(h) << 8) + uint16(l)
	return nil
}

// ReceiveNMI ...
func (c *CPU) ReceiveNMI(active bool) {
	log.Trace("begin[%v] ...", active)
	defer log.Trace("end[%v]", active)
	if !c.beforeNMIActive && active {
		// activeが　false => true となったときにNMI割り込みを発生
		c.shouldNMI = true
	}
	c.beforeNMIActive = active
}

// pushStack ...
func (c *CPU) pushStack(b byte) error {
	addr := domain.Address(uint16(0x0100) | uint16(c.registers.S))
	err := c.bus.WriteByCPU(addr, b)
	if err != nil {
		err = xerrors.Errorf(": %w", err)
	}
	c.registers.S--
	log.Trace("CPU.pushStack[%2X] => [addr=%v]", b, addr)
	return err
}

// popStack ...
func (c *CPU) popStack() (byte, error) {
	c.registers.S++
	addr := domain.Address(uint16(0x0100) | uint16(c.registers.S))
	b, err := c.bus.ReadByCPU(addr)
	if err != nil {
		err = xerrors.Errorf(": %w", err)
	}
	log.Trace("CPU.popStack[%2X] <= [addr=%v]", b, addr)
	return b, err
}
