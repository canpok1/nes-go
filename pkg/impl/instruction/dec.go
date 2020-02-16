package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// DEC ...
type DEC struct {
	BaseInstruction
}

// Execute ...
func (c *DEC) Execute(op []byte) (cycle int, err error) {
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
}
