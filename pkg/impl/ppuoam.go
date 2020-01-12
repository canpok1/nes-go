package impl

import "nes-go/pkg/domain"

// PPUOAM ...
type PPUOAM struct {
	store            []byte
	lastWriteAddress uint8
}

// NewPPUOAM ...
func NewPPUOAM() *PPUOAM {
	p := PPUOAM{
		store:            make([]byte, 256),
		lastWriteAddress: 0,
	}
	return &p
}

// Write ...
func (p *PPUOAM) Write(oamaddr uint8, b byte) {
	p.store[oamaddr] = b
	p.lastWriteAddress = oamaddr
}

// Each ...
func (p *PPUOAM) Each(exec func(domain.Sprite) error) error {
	s := domain.Sprite{}
	for i := uint8(0); i <= p.lastWriteAddress; i++ {
		b := p.store[i]
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
