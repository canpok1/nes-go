package instruction

import (
	"nes-go/pkg/log"
)

// INX ...
type INX struct {
	BaseInstruction
}

// Execute ...
func (c *INX) Execute(op []byte) (cycle int, err error) {
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

	c.registers.UpdateX(c.registers.X + 1)
	c.registers.P.UpdateN(c.registers.X)
	c.registers.P.UpdateZ(c.registers.X)
	return
}
