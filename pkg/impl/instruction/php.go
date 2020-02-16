package instruction

import (
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// PHP ...
type PHP struct {
	BaseInstruction
}

// Execute ...
func (c *PHP) Execute(op []byte) (cycle int, err error) {
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

	p := c.registers.P.ToByte()

	// 6502のバグ：スタックに格納するフラグはBフラグがセットされた状態になる
	p = (p & 0xEF) | 0x10

	if err = c.pushStack(p); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	return
}
