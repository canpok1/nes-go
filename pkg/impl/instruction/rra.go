package instruction

import (
	"nes-go/pkg/domain"
	"nes-go/pkg/log"

	"golang.org/x/xerrors"
)

// RRA ...
type RRA struct {
	BaseInstruction
}

// Execute ...
func (c *RRA) Execute(op []byte) (cycle int, err error) {
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

	var b, ans1 byte
	if mode == domain.Accumulator {
		b = c.registers.A
		c.recorder.Data = &b

		ans1 = b >> 1
		if c.registers.P.Carry {
			ans1 = ans1 | 0x80
		}

		c.registers.A = ans1
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
		c.recorder.Data = &b

		ans1 = b >> 1
		if c.registers.P.Carry {
			ans1 = ans1 | 0x80
		}

		err = c.bus.WriteByCPU(addr, ans1)
		if err != nil {
			err = xerrors.Errorf(": %w", err)
			return
		}
	}
	c.registers.P.Carry = (b & 0x01) == 0x01

	ans2 := uint16(c.registers.A) + uint16(ans1)
	if c.registers.P.Carry {
		ans2 = ans2 + 1
	}

	beforeA := c.registers.A
	c.registers.A = byte(ans2 & 0xFF)
	c.registers.P.UpdateN(c.registers.A)
	if (b & 0x80) == 0x00 {
		c.registers.P.UpdateV(beforeA, c.registers.A)
	} else {
		c.registers.P.ClearV()
	}
	c.registers.P.UpdateZ(c.registers.A)
	c.registers.P.UpdateC(ans2)
	return
}
