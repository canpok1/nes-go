package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// SBC ...
type SBC struct {
	BaseInstruction
}

// Execute ...
func (c *SBC) Execute(op []byte) (cycle int, err error) {
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
		if addr, _, err = c.makeAddress(op); err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}

		b, err = c.bus.ReadByCPU(addr)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
	}
	c.recorder.Data = &b
	ans := uint16(c.registers.A) - uint16(b)
	if !c.registers.P.Carry {
		ans = ans - 1
	}

	beforeA := c.registers.A
	c.registers.A = byte(ans & 0x00FF)
	c.registers.P.UpdateN(c.registers.A)
	if (beforeA & 0x80) == (b & 0x80) {
		c.registers.P.ClearV()
	} else {
		c.registers.P.UpdateV(beforeA, c.registers.A)
	}
	c.registers.P.UpdateZ(c.registers.A)
	c.registers.P.Carry = (ans & 0xFF00) == 0x0000
	return
}
