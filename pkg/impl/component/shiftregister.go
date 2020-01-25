package component

// ShiftRegister16bit ...
type ShiftRegister16bit struct {
	high byte
	low  byte
}

// GetLow ...
func (s *ShiftRegister16bit) GetLow() byte {
	return s.low
}

// SetHigh ...
func (s *ShiftRegister16bit) SetHigh(d byte) {
	s.high = d
}

// Shift ...
func (s *ShiftRegister16bit) Shift() {
	s.low = (s.low >> 1) | ((s.high & 0x01) << 7)
	s.high = s.high >> 1
}

