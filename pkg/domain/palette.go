package domain

// Palette ...
type Palette []byte

// NewPalette ...
func NewPalette() *Palette {
	p := Palette(make([]byte, 4))
	return &p
}

// GetColor ...
func (p *Palette) GetColor(no uint8) (uint8, uint8, uint8) {
	index := (*p)[no]
	c := colors[index]
	return c[0], c[1], c[2]
}
