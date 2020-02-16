package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// LDY ...
type LDY struct {
	BaseInstruction
}

// Execute ...
func (c *LDY) Execute(op []byte) (cycle int, err error) {
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

	var b byte
	if mode == domain.Immediate {
		if len(op) < 1 {
			err = xerrors.Errorf("data is nil; mnemonic: %#v, op: %#v", mne, op)
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

		if c.ocp.AddressingMode == domain.IndexedAbsoluteX && pageCrossed {
			cycle++
		}
	}
	c.recorder.Data = &b

	c.registers.UpdateY(b)
	c.registers.P.UpdateN(c.registers.Y)
	c.registers.P.UpdateZ(c.registers.Y)
	return
}
