package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// JSR ...
type JSR struct {
	BaseInstruction
}

// Execute ...
func (c *JSR) Execute(op []byte) (cycle int, err error) {
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

	var addr domain.Address
	if addr, _, err = c.makeAddress(op); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}

	// 6502のバグ:1つ前のアドレスを格納
	pc := c.registers.PC - 1

	if err = c.pushStack(byte((pc & 0xFF00) >> 8)); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	if err = c.pushStack(byte(pc & 0x00FF)); err != nil {
		err = xerrors.Errorf(": %w", err)
		return
	}
	c.registers.PC = uint16(addr)
	return
}
