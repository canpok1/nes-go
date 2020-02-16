package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// SAX ...
type SAX struct {
	BaseInstruction
}

// Execute ...
func (c *SAX) Execute(op []byte) (cycle int, err error) {
	mne := c.ocp.Mnemonic
	mode := c.ocp.AddressingMode
	cycle = c.ocp.Cycle

	log.Trace("begin[%#x][%v][%v][%#v] ...", c.registers.PC, mne, mode, op)

	c.recorder.Mnemonic = mne
	c.recorder.Documented = c.ocp.Documented
	c.recorder.AddressingMode = mode

	defer func() {
		if err != nil {
			log.Warn("end[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("end[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

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

		switch c.ocp.AddressingMode {
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
}
