package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// BIT ...
type BIT struct {
	BaseInstruction
}

// Execute ...
func (c *BIT) Execute(op []byte) (cycle int, err error) {
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
}
