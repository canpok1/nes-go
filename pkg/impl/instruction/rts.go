package instruction

import (
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// RTS ...
type RTS struct {
	BaseInstruction
}

// Execute ...
func (c *RTS) Execute(op []byte) (cycle int, err error) {
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

	var l, h byte
	if l, err = c.popStack(); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	if h, err = c.popStack(); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	c.registers.PC = (uint16(h) << 8) | uint16(l)

	// 6502のバグ:インクリメントしたもの復帰アドレスとする
	c.registers.PC = c.registers.PC + 1

	return
}
