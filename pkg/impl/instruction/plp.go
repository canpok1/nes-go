package instruction

import (
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// PLP ...
type PLP struct {
	BaseInstruction
}

// Execute ...
func (c *PLP) Execute(op []byte) (cycle int, err error) {
	mne := c.ocp.Mnemonic
	mode := c.ocp.AddressingMode
	cycle = c.ocp.Cycle

	log.Trace("BaseInstruction.Execute[%#x][%v][%v][%#v] ...", c.registers.PC, mne, mode, op)

	c.recorder.Mnemonic = mne
	c.recorder.Documented = c.ocp.Documented
	c.recorder.AddressingMode = mode

	defer func() {
		if err != nil {
			log.Warn("BaseInstruction.Execute[%v][%v][%#v] => %v", mne, mode, op, err)
		} else {
			log.Trace("BaseInstruction.Execute[%v][%v][%#v] => completed", mne, mode, op)
		}
	}()

	var b byte
	if b, err = c.popStack(); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	c.registers.P.UpdateAll(b)
	return
}
