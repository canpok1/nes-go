package impl

// Sprite ...
type Sprite struct {
	Y         byte
	TileIndex byte
	Attribute byte
	X         byte
}

// PPUOAM ...
type PPUOAM []byte

// NewPPUOAM ...
func NewPPUOAM() *PPUOAM {
	p := PPUOAM(make([]byte, 256))
	return &p
}

// Write ...
func (p *PPUOAM) Write(oamaddr uint8, b byte) {
	(*p)[oamaddr] = b
}

// Each ...
func (p *PPUOAM) Each(exec func(Sprite) error) error {
	s := Sprite{}
	for i, b := range *p {
		offset := i % 4
		switch offset {
		case 0:
			s.Y = b
		case 1:
			s.TileIndex = b
		case 2:
			s.Attribute = b
		case 3:
			s.X = b
			if err := exec(s); err != nil {
				return err
			}
		}
	}
	return nil
}
