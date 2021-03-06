package instruction

import (
	"nes-go/pkg/log"
)

// TAY ...
type TAY struct {
	BaseInstruction
}

// Execute ...
func (c *TAY) Execute(op []byte) (cycle int, err error) {
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

	c.registers.Y = c.registers.A
	c.registers.P.UpdateN(c.registers.Y)
	c.registers.P.UpdateZ(c.registers.Y)
	return
}
