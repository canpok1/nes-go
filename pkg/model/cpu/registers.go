package cpu

import (
	"fmt"

	"github.com/canpok1/nes-go/pkg/log"
)

// CPURegisters ...
type CPURegisters struct {
	A  byte
	X  byte
	Y  byte
	S  byte
	P  *CPUStatusRegister
	PC uint16
}

// NewCPURegisters ...
func NewCPURegisters() *CPURegisters {
	// initialize as CPU power up state
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state
	return &CPURegisters{
		A:  0,
		X:  0,
		Y:  0,
		S:  0xFD,
		P:  NewCPUStatusRegister(),
		PC: 0,
	}
}

// String ...
func (r *CPURegisters) String() string {
	return fmt.Sprintf(
		"{A:%#v, X:%#v, Y:%#v, S:%#v, P:%v, PC:%#v}",
		r.A,
		r.X,
		r.Y,
		r.S,
		r.P.String(),
		r.PC,
	)
}

// UpdateA ...
func (r *CPURegisters) UpdateA(a byte) {
	old := r.A
	r.A = a
	log.Trace("CPU.update[A] %#v => %#v", old, r.A)
}

// UpdateX ...
func (r *CPURegisters) UpdateX(x byte) {
	old := r.X
	r.X = x
	log.Trace("CPU.update[X] %#v => %#v", old, r.X)
}

// UpdateY ...
func (r *CPURegisters) UpdateY(y byte) {
	old := r.Y
	r.Y = y
	log.Trace("CPU.update[Y] %#v => %#v", old, r.Y)
}

// UpdateS ...
func (r *CPURegisters) UpdateS(s byte) {
	old := r.S
	r.S = s
	log.Trace("CPU.update[S] %#v => %#v", old, r.S)
}

// IncrementPC ...
func (r *CPURegisters) IncrementPC() {
	r.UpdatePC(r.PC + 1)
}

// UpdatePC ...
func (r *CPURegisters) UpdatePC(pc uint16) {
	old := r.PC
	r.PC = pc
	log.Trace("CPU.update[PC] %#v => %#v", old, r.PC)
}

// CPUStatusRegister ...
// https://qiita.com/bokuweb/items/1575337bef44ae82f4d3#%E3%82%B9%E3%83%86%E3%83%BC%E3%82%BF%E3%82%B9%E3%83%AC%E3%82%B8%E3%82%B9%E3%82%BF
type CPUStatusRegister struct {
	Negative    bool // bit7	N	ネガティブ	演算結果のbit7が1の時にセット
	Overflow    bool // bit6	V	オーバーフロー	P演算結果がオーバーフローを起こした時にセット
	Reserved    bool // bit5	R	予約済み	常にセットされている
	BreakMode   bool // bit4	B	ブレークモード	BRK発生時にセット、IRQ発生時にクリア
	DecimalMode bool // bit3	D	デシマルモード	0:デフォルト、1:BCDモード (未実装)
	Interrupt   bool // bit2	I	IRQ禁止	0:IRQ許可、1:IRQ禁止
	Zero        bool // bit1	Z	ゼロ	演算結果が0の時にセット
	Carry       bool // bit0	C	キャリー	キャリー発生時にセット
}

// NewCPUStatusRegister ...
func NewCPUStatusRegister() *CPUStatusRegister {
	return &CPUStatusRegister{
		Negative:    false,
		Overflow:    false,
		Reserved:    true,
		BreakMode:   true,
		DecimalMode: false,
		Interrupt:   true,
		Zero:        false,
		Carry:       false,
	}
}

// String ...
func (s *CPUStatusRegister) String() string {
	return fmt.Sprintf(
		"{N:%#v, V:%#v, R:%#v, B:%#v, D:%#v, I:%#v, Z:%#v, C:%#v}",
		s.Negative,
		s.Overflow,
		s.Reserved,
		s.BreakMode,
		s.DecimalMode,
		s.Interrupt,
		s.Zero,
		s.Carry,
	)
}

// ToByte ...
func (s *CPUStatusRegister) ToByte() byte {
	var b byte = 0
	if s.Negative {
		b = b + 0x80
	}
	if s.Overflow {
		b = b + 0x40
	}
	if s.Reserved {
		b = b + 0x20
	}
	if s.BreakMode {
		b = b + 0x10
	}
	if s.DecimalMode {
		b = b + 0x08
	}
	if s.Interrupt {
		b = b + 0x04
	}
	if s.Zero {
		b = b + 0x02
	}
	if s.Carry {
		b = b + 0x01
	}
	return b
}

// UpdateAll ...
func (s *CPUStatusRegister) UpdateAll(b byte) {
	s.Negative = (b & 0x80) == 0x80
	s.Overflow = (b & 0x40) == 0x40
	s.Reserved = (b & 0x20) == 0x20
	s.BreakMode = (b & 0x10) == 0x10
	s.DecimalMode = (b & 0x08) == 0x08
	s.Interrupt = (b & 0x04) == 0x04
	s.Zero = (b & 0x02) == 0x02
	s.Carry = (b & 0x01) == 0x01
}

// UpdateN ...
func (s *CPUStatusRegister) UpdateN(result byte) {
	old := s.Negative
	s.Negative = ((result & 0x80) == 0x80)
	log.Trace("CPU.update[N] %#v => %#v", old, s.Negative)
}

// UpdateV ...
func (s *CPUStatusRegister) UpdateV(result int16) {
	old := s.Overflow
	s.Overflow = (result < 0x7F) || (result > 0x80)
	log.Trace("CPU.update[V] %#v => %#v", old, s.Overflow)
}

// UpdateI ...
func (s *CPUStatusRegister) UpdateI(i bool) {
	old := s.Interrupt
	s.Interrupt = i
	log.Trace("CPU.update[I] %#v => %#v", old, s.Interrupt)
}

// UpdateZ ...
func (s *CPUStatusRegister) UpdateZ(result byte) {
	old := s.Zero
	s.Zero = (result == 0x00)
	log.Trace("CPU.update[Z] %#v => %#v", old, s.Zero)
}

// UpdateC ...
func (s *CPUStatusRegister) UpdateC(result int16) {
	old := s.Carry
	s.Carry = result > 0xFF
	log.Trace("CPU.update[C] %#v => %#v", old, s.Carry)
}
